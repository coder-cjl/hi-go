# 天气 API 配置和故障排除

## 问题诊断

根据错误日志：
```
invalid character '<' looking for beginning of value
```

这表明和风天气 API 返回的是 HTML 错误页面而不是 JSON。常见原因：

## 1. API Key 无效或过期

### 获取新的 API Key

1. 访问 https://dev.qweather.com/
2. 注册/登录账号
3. 进入"控制台" → "应用管理"
4. 创建新应用或查看现有应用
5. 复制 **Web API** 的 Key（不是 SDK Key）
6. 注意：免费订阅每天限制 1000 次请求

### 验证 API Key

```bash
# 将 YOUR_API_KEY 替换为实际的 key
curl "https://devapi.qweather.com/v7/weather/now?location=101010100&key=YOUR_API_KEY"
```

正确的响应应该是：
```json
{
  "code": "200",
  "updateTime": "...",
  "now": {
    "temp": "20",
    ...
  }
}
```

错误的响应（HTML）：
```html
<!doctype html>...
```

## 2. API 访问限制

和风天气免费版限制：
- **每天 1000 次调用**
- 超过限制会返回错误

检查你的使用量：
1. 登录 https://dev.qweather.com/
2. 查看"控制台" → "使用统计"

## 3. 配置正确的 API

### 开发版 vs 商业版

- **开发版**（免费）: `https://devapi.qweather.com/v7`
- **商业版**（付费）: `https://api.qweather.com/v7`

确保在配置文件中使用正确的 URL。

## 4. 更新配置

编辑 `configs/dev.yaml`:

```yaml
ai:
  skills:
    weather:
      enabled: true
      provider: qweather
      api_key: "在这里填入你的有效 API Key"  # ← 重要！
      base_url: https://devapi.qweather.com/v7  # 开发版
      timeout: 10
      cache_ttl: 1800
      max_retries: 3
```

## 5. 测试 API

重启服务后测试：

```bash
curl -X POST http://localhost:8000/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "北京今天天气怎么样？"}'
```

## 6. 查看日志

如果仍然有问题，日志现在会显示更详细的错误信息：

- HTTP 状态码
- API 响应的前 200 字符
- 具体的错误代码

这样可以更容易定位问题。

## 常见错误代码

和风天气 API 错误代码：

| Code | 说明 |
|------|------|
| 200 | 成功 |
| 400 | 请求错误，可能是参数错误 |
| 401 | 认证失败，Key 无效 |
| 402 | 超过访问次数 |
| 403 | 无访问权限 |
| 404 | 数据不存在 |
| 429 | 请求次数超限 |
| 500 | 服务器错误 |

## 替代方案

如果和风天气不可用，可以：

1. **使用其他天气 API**
   - OpenWeather API
   - AccuWeather API
   
2. **禁用天气功能**
   ```yaml
   ai:
     skills:
       weather:
         enabled: false  # 禁用天气查询
   ```

3. **实现 mock 数据**（仅用于测试）

## 下一步

1. 获取有效的和风天气 API Key
2. 更新 `configs/dev.yaml` 中的 `api_key`
3. 重启服务
4. 再次测试

如果问题仍然存在，请查看服务日志获取详细错误信息。
