# SSE工具调用问题修复说明

## 问题描述

之前的SSE实现在处理工具调用时存在以下问题：

1. **响应不完整**：工具调用后没有继续生成AI的最终回答
2. **工具调用片段化**：DeepSeek的流式API返回的工具调用是增量式的，被错误地当作多个独立的工具调用
3. **调试信息泄露**：向客户端发送了tool_calls事件，但格式不正确且用户不需要看到

## 修复内容

### 1. Service层：实现多轮对话逻辑

**文件**：`src/service/aiservice/service.go`

**改进**：
- `ChatStream` 方法现在支持完整的多轮对话流程
- 检测到工具调用时，自动执行工具并继续请求AI生成回答
- 最多支持5轮对话（防止无限循环）

**流程**：
```
用户问题 
  ↓
AI流式响应（包含工具调用）
  ↓
执行工具（如查询天气）
  ↓
AI流式响应（基于工具结果的最终回答）← 这个之前缺失
  ↓
完成
```

### 2. DeepSeek客户端：正确累积增量工具调用

**文件**：`src/service/aiservice/deepseek_client.go`

**问题**：
DeepSeek的流式API返回的tool_calls是增量式的：
```json
// 第1个chunk
{"tool_calls": [{"index": 0, "id": "call_xxx", "type": "function", "function": {"name": "get_weather"}}]}

// 第2个chunk  
{"tool_calls": [{"index": 0, "function": {"arguments": "{"}}]}

// 第3个chunk
{"tool_calls": [{"index": 0, "function": {"arguments": "\"location\""}}]}

// 第4个chunk
{"tool_calls": [{"index": 0, "function": {"arguments": ":"}}]}

// ... 更多chunks
```

**改进**：
- 新增 `StreamToolCallDelta` 结构体处理增量数据
- 使用map按index累积工具调用
- 合并 `function.arguments` 的所有增量
- 只在完成时发送最终的完整工具调用

### 3. Handler层：简化SSE响应

**文件**：`src/handler/ai_handler.go`

**改进**：
- 移除了 `tool_calls` 事件发送（用户不需要看到内部工具调用过程）
- 只发送 `message` 事件（内容增量）和 `done` 事件（完成标志）
- 保持简洁的用户体验

## 测试验证

### 测试场景1：简单问答（无工具调用）

**请求**：
```bash
curl -X POST http://localhost:8080/api/ai/chat2 \
  -H "Content-Type: application/json" \
  -d '{"message": "你好"}' -N
```

**预期响应**：
```
event:message
data:{"content": "你好"}

event:message
data:{"content": "！"}

event:message
data:{"content": "有什么"}

event:message
data:{"content": "可以"}

event:message
data:{"content": "帮"}

event:message
data:{"content": "您"}

event:message
data:{"content": "的"}

event:message
data:{"content": "吗"}

event:message
data:{"content": "？"}

event:done
data:{"finish_reason": "stop"}
```

### 测试场景2：需要工具调用（天气查询）

**请求**：
```bash
curl -X POST http://localhost:8080/api/ai/chat2 \
  -H "Content-Type: application/json" \
  -d '{"message": "上海的天气适合跑步吗？"}' -N
```

**预期响应流程**：
1. **阶段1**：AI表示要查询天气
   ```
   event:message
   data:{"content": "我来"}
   
   event:message
   data:{"content": "帮您"}
   
   event:message
   data:{"content": "查询"}
   
   event:message
   data:{"content": "上海的"}
   
   event:message
   data:{"content": "天气"}
   ```

2. **内部处理**：执行 `get_weather("上海")` 工具调用（用户不可见）

3. **阶段2**：AI基于天气数据生成回答
   ```
   event:message
   data:{"content": "根据"}
   
   event:message
   data:{"content": "查询"}
   
   event:message
   data:{"content": "，"}
   
   event:message
   data:{"content": "上海"}
   
   event:message
   data:{"content": "今天"}
   
   event:message
   data:{"content": "天气"}
   
   event:message
   data:{"content": "晴朗"}
   
   event:message
   data:{"content": "，"}
   
   event:message
   data:{"content": "温度"}
   
   event:message
   data:{"content": "15"}
   
   event:message
   data:{"content": "°C"}
   
   event:message
   data:{"content": "，"}
   
   event:message
   data:{"content": "湿度"}
   
   event:message
   data:{"content": "60%"}
   
   event:message
   data:{"content": "，"}
   
   event:message
   data:{"content": "非常"}
   
   event:message
   data:{"content": "适合"}
   
   event:message
   data:{"content": "跑步"}
   
   event:message
   data:{"content": "！"}
   
   event:done
   data:{"finish_reason": "stop"}
   ```

## 技术细节

### 增量工具调用合并算法

```go
// 使用map按index累积
toolCallsMap := make(map[int]*ToolCall)

for _, deltaTC := range delta.ToolCalls {
    idx := deltaTC.Index
    if _, exists := toolCallsMap[idx]; !exists {
        // 首次出现，初始化
        toolCallsMap[idx] = &ToolCall{
            ID:   deltaTC.ID,
            Type: deltaTC.Type,
            Function: FunctionCall{
                Name:      deltaTC.Function.Name,
                Arguments: deltaTC.Function.Arguments,
            },
        }
    } else {
        // 累积arguments
        toolCallsMap[idx].Function.Arguments += deltaTC.Function.Arguments
    }
}
```

### 多轮对话循环

```go
maxRounds := 5
for round := 0; round < maxRounds; round++ {
    // 1. 调用AI流式接口
    streamChan, err := s.client.ChatStream(ctx, messages, tools)
    
    // 2. 读取并转发流式响应
    for resp := range streamChan { ... }
    
    // 3. 如果有工具调用，执行并继续
    if len(accumulatedToolCalls) > 0 {
        toolResults, _ := s.executeToolCalls(ctx, accumulatedToolCalls)
        messages = append(messages, toolResults...)
        continue // 下一轮
    }
    
    // 4. 无工具调用，完成
    return
}
```

## 验证清单

- [x] 编译成功
- [x] 简单问答正常工作
- [x] 工具调用能完整执行并返回最终回答
- [x] SSE响应格式正确
- [x] 不再向客户端发送tool_calls调试信息
- [x] 网页测试客户端正常显示打字机效果

## 后续优化建议

1. **进度提示**：可以考虑在工具执行时发送一个 `event:status` 提示用户"正在查询天气..."
2. **超时处理**：为工具执行添加超时保护
3. **错误重试**：工具调用失败时可以重试或优雅降级
4. **并行工具调用**：如果AI同时调用多个工具，可以并行执行

## 总结

修复后的SSE实现：
- ✅ **完整性**：多轮对话流程完整，工具调用后继续生成回答
- ✅ **正确性**：正确处理DeepSeek流式API的增量工具调用
- ✅ **简洁性**：客户端只看到内容流，不看到内部实现细节
- ✅ **用户体验**：流畅的打字机效果，无中断
