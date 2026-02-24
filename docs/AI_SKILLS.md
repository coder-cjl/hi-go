# AI Skills 功能使用说明

## 功能概述

本项目集成了 AI Skills 功能，支持通过 DeepSeek AI 调用天气查询等技能。AI 会根据用户的问题自动选择合适的工具来获取信息并给出回答。

## 架构设计

### 核心组件

1. **Skill 接口** (`src/service/aiservice/skill.go`)
   - 定义统一的技能接口
   - 可扩展的技能注册表

2. **天气查询技能** (`src/service/aiservice/weather_skill.go`)
   - 支持和风天气 API
   - 自动缓存查询结果（30分钟）
   - 支持中国主要城市

3. **AI 服务** (`src/service/aiservice/service.go`)
   - 处理对话逻辑
   - 自动执行工具调用
   - 支持多轮对话

4. **DeepSeek 客户端** (`src/service/aiservice/deepseek_client.go`)
   - 符合 OpenAI Function Calling 标准
   - 支持工具调用（Tools）

## 配置说明

在 `configs/dev.yaml` 中添加以下配置：

```yaml
ai:
  enabled: true                    # 是否启用AI功能
  provider: deepseek               # AI提供商
  model: deepseek-chat             # 模型名称
  api_key: your-deepseek-api-key   # DeepSeek API密钥
  base_url: https://api.deepseek.com
  timeout: 30
  max_tokens: 4000
  temperature: 0.7
  system_prompt: "你是一个智能助手，可以帮助用户查询天气等信息。请根据用户的问题，使用合适的工具来获取信息并给出回答。"
  
  skills:
    weather:
      enabled: true                    # 是否启用天气查询功能
      provider: qweather               # 天气服务提供商
      api_key: your-qweather-api-key   # 和风天气 API Key
      base_url: https://devapi.qweather.com/v7
      timeout: 10
      cache_ttl: 1800                  # 缓存30分钟
      max_retries: 3
```

### 获取 API 密钥

1. **DeepSeek API Key**
   - 访问：https://platform.deepseek.com/
   - 注册并创建 API Key

2. **和风天气 API Key**
   - 访问：https://dev.qweather.com/
   - 注册并创建项目获取 API Key
   - 开发版免费，支持1000次/天的请求

## API 使用

### 发送对话请求

**请求**

```bash
POST /api/ai/chat
Content-Type: application/json

{
  "message": "北京今天天气怎么样？"
}
```

**响应**

```json
{
  "trace_id": "550e8400-e29b-41d4-a716-446655440000",
  "code": 0,
  "message": "success",
  "data": {
    "reply": "北京今天天气晴朗，温度20°C，体感温度19°C，湿度45%，风速15km/h。适合外出活动！"
  }
}
```

## 使用示例

### 示例 1：查询天气

```bash
curl -X POST http://localhost:8000/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "深圳今天天气怎么样？"}'
```

### 示例 2：多城市对比

```bash
curl -X POST http://localhost:8000/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "比较一下北京和上海今天的天气"}'
```

### 示例 3：询问建议

```bash
curl -X POST http://localhost:8000/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "我明天要去广州出差，需要注意什么？"}}'
```

## 代码结构

```
src/
├── config/
│   └── types.go                    # 添加了 AI 相关配置类型
├── handler/
│   └── ai_handler.go               # AI 对话 HTTP 处理器
├── router/
│   ├── ai_router.go                # AI 路由注册
│   └── router.go                   # 主路由（已集成 AI 路由）
├── service/
│   └── aiservice/
│       ├── skill.go                # Skill 接口定义
│       ├── weather_skill.go        # 天气查询技能实现
│       ├── service.go              # AI 服务核心逻辑
│       ├── deepseek_client.go      # DeepSeek API 客户端
│       └── init.go                 # AI 服务初始化
└── utils/
    └── cache/
        ├── cache.go                # 缓存接口
        └── redis_cache.go          # Redis 缓存实现
```

## 扩展新技能

要添加新的 Skill（如翻译、汇率查询等），只需：

1. 实现 `Skill` 接口
2. 在配置中添加对应的Config结构
3. 在 `initialize/ai.go` 中注册新技能

示例：

```go
// 1. 实现 Skill 接口
type TranslationSkill struct {
    config config.TranslationSkillConfig
}

func (t *TranslationSkill) Name() string {
    return "translate"
}

func (t *TranslationSkill) Description() string {
    return "翻译文本到指定语言"
}

func (t *TranslationSkill) Parameters() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "text": map[string]string{
                "type": "string",
                "description": "要翻译的文本",
            },
            "target_lang": map[string]string{
                "type": "string",
                "description": "目标语言，如：en, zh, ja",
            },
        },
        "required": []string{"text", "target_lang"},
    }
}

func (t *TranslationSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
    // 实现翻译逻辑
    return translationResult, nil
}

func (t *TranslationSkill) IsEnabled() bool {
    return t.config.Enabled
}

// 2. 在 initialize/ai.go 中注册
if config.Config.AI.Skills.Translation.Enabled {
    translationSkill := NewTranslationSkill(config.Config.AI.Skills.Translation)
    registry.Register(translationSkill)
}
```

## 注意事项

1. **API 密钥安全**：请勿将 API 密钥提交到代码仓库，建议使用环境变量
2. **请求限制**：注意 API 提供商的请求频率限制
3. **缓存策略**：天气数据缓存30分钟，城市ID永久缓存
4. **错误处理**：API 调用失败会返回友好的错误提示

## 性能优化

- ✅ Redis 缓存减少重复API调用
- ✅ 城市ID缓存避免重复查询
- ✅ 30分钟天气数据缓存
- ✅ 支持并发请求

## 未来扩展

- [ ] 支持更多天气数据提供商（OpenWeather等）
- [ ] 添加更多 Skills（翻译、搜索、计算等）
- [ ] 支持流式响应（SSE）
- [ ] 添加AI对话历史记录
- [ ] 支持多模型切换
