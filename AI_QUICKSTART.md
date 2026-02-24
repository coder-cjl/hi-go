# AI Skills 快速开始指南

## 前置条件

1. Go 1.24+ 已安装
2. Redis 已安装并运行
3. MySQL 已安装并运行

## 第一步：获取 API 密钥

### 1. DeepSeek API Key

1. 访问 https://platform.deepseek.com/
2. 注册账号并登录
3. 创建API Key并复制

### 2. 和风天气 API Key

1. 访问 https://dev.qweather.com/
2. 注册账号并登录
3. 创建应用并选择"Web API"
4. 选择"免费订阅"（每天1000次请求）
5. 复制API Key

## 第二步：配置项目

编辑 `configs/dev.yaml` 文件，填入你的 API 密钥：

```yaml
ai:
  enabled: true
  provider: deepseek
  model: deepseek-chat
  api_key: "sk-xxxxxxxxxxxxxx"  # 替换为你的DeepSeek API Key
  base_url: https://api.deepseek.com
  timeout: 30
  max_tokens: 4000
  temperature: 0.7
  system_prompt: "你是一个智能助手，可以帮助用户查询天气等信息。请根据用户的问题，使用合适的工具来获取信息并给出回答。"
  
  skills:
    weather:
      enabled: true
      provider: qweather
      api_key: "xxxxxxxxxxxxxx"  # 替换为你的和风天气API Key
      base_url: https://devapi.qweather.com/v7
      timeout: 10
      cache_ttl: 1800
      max_retries: 3
```

## 第三步：启动服务

```bash
# 方式1：直接运行
go run main.go

# 方式2：编译后运行
go build -o hi-go .
./hi-go

# 方式3：使用提供的脚本
./run.sh
```

服务将在 `http://localhost:8000` 启动。

## 第四步：测试功能

### 方式1：使用 curl 命令

```bash
curl -X POST http://localhost:8000/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "北京今天天气怎么样？"}'
```

### 方式2：使用测试脚本

```bash
./test_ai_skills.sh
```

### 方式3：使用 Postman 或其他 API 工具

- **URL**: `http://localhost:8000/api/ai/chat`
- **Method**: `POST`
- **Headers**: `Content-Type: application/json`
- **Body**:
  ```json
  {
    "message": "深圳今天天气怎么样？"
  }
  ```

## 示例对话

### 查询单个城市

**请求**:
```json
{
  "message": "上海今天天气如何？"
}
```

**响应**:
```json
{
  "trace_id": "...",
  "code": 0,
  "message": "success",
  "data": {
    "reply": "上海今天天气多云，温度18°C，体感温度17°C，湿度60%，风速12km/h。"
  }
}
```

### 对比多个城市

**请求**:
```json
{
  "message": "比较一下北京和深圳今天的天气"
}
```

**响应**:
AI会自动调用两次天气查询，然后对比结果给出回答。

### 咨询建议

**请求**:
```json
{
  "message": "我明天要去广州出差，需要带伞吗？"
}
```

**响应**:
AI会查询广州天气，并根据天气状况给出建议。

## 常见问题

### Q: 提示 "deepseek api error"
A: 检查你的 DeepSeek API Key 是否正确，以及是否有足够的额度。

### Q: 提示 "qweather api error"
A: 检查你的和风天气 API Key 是否正确，订阅计划是否有效。

### Q: 查询城市失败
A: 确认城市名称使用中文，如："北京"、"上海"、"深圳"等。

### Q: Redis 连接失败
A: 确保 Redis 服务正在运行，并检查 `configs/dev.yaml` 中的 Redis 配置。

### Q: 如何禁用 AI 功能
A: 在 `configs/dev.yaml` 中设置：
```yaml
ai:
  enabled: false
```

## 架构说明

```
用户请求 → AI Handler → AI Service → DeepSeek API
                              ↓
                        Tool Calling
                              ↓
                      Weather Skill → 和风天气API
                              ↓
                         Redis Cache
```

1. 用户发送问题到 `/api/ai/chat`
2. AI Service 将问题发送给 DeepSeek
3. DeepSeek 识别需要调用天气查询工具
4. Weather Skill 调用和风天气API（优先使用缓存）
5. 将天气数据返回给 DeepSeek
6. DeepSeek 基于天气数据生成友好的回答
7. 返回给用户

## 下一步

- 查看完整文档：[docs/AI_SKILLS.md](docs/AI_SKILLS.md)
- 了解如何扩展新的 Skills
- 配置生产环境

## 支持

如有问题，请查看日志输出或提交 Issue。
