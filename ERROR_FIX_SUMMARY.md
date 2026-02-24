# 错误修复总结

## 错误原因

你遇到的错误：
```
invalid character '<' looking for beginning of value
```

**根本原因**：和风天气 API 返回了 **HTTP 400 Bad Request**（HTML 错误页面），而不是预期的 JSON 数据。

通过测试发现：
```bash
curl "https://devapi.qweather.com/v7/city/lookup?location=北京&key=HE2409021346141011"
# 返回：<!doctype html><html>... HTTP Status 400 – Bad Request
```

说明 **API Key 无效或已过期**。

## 已完成的改进

### 1. 增强错误处理 ✅

现在会显示详细的错误信息：
- HTTP 状态码
- API 响应内容预览（前 200 字符）
- 更明确的错误提示

修改的文件：
- [weather_skill.go](../src/service/aiservice/weather_skill.go)

### 2. URL 编码支持 ✅

添加了对中文城市名的 URL 编码：
```go
url.QueryEscape(location)  // 将 "北京" 编码为 "%E5%8C%97%E4%BA%AC"
```

### 3. 创建故障排除文档 ✅

- [WEATHER_API_TROUBLESHOOTING.md](WEATHER_API_TROUBLESHOOTING.md) - 详细的问题诊断和解决方案

## 你需要做的事情

### 📌 必须：更新 API Key

1. **获取新的 API Key**
   - 访问：https://dev.qweather.com/
   - 注册/登录账号
   - 创建应用并获取 **Web API** Key（不是 SDK Key）

2. **更新配置文件**
   
   编辑 `configs/dev.yaml`：
   ```yaml
   ai:
     skills:
       weather:
         api_key: "在这里填入你的新 API Key"  # ← 替换这里！
   ```

3. **验证 API Key**
   
   使用新的 key 测试：
   ```bash
   # 替换 YOUR_NEW_KEY
   curl "https://devapi.qweather.com/v7/weather/now?location=101010100&key=YOUR_NEW_KEY"
   ```
   
   成功的响应应该是 JSON 格式：
   ```json
   {
     "code": "200",
     "now": {
       "temp": "20",
       ...
     }
   }
   ```

4. **重启服务**
   ```bash
   go run main.go
   ```

5. **测试功能**
   ```bash
   curl -X POST http://localhost:8000/api/ai/chat \
     -H "Content-Type: application/json" \
     -d '{"message": "北京今天天气怎么样？"}'
   ```

## 现在的优势

即使 API 再次出现问题，你会看到更清晰的错误信息：

**之前**：
```
invalid character '<' looking for beginning of value
```

**现在**：
```
qweather city lookup returned status 400, response: <!doctype html><html lang="en"><head><title>HTTP Status 400...
```

这样能立即知道是 API 返回了错误页面，而不是代码问题。

## 其他选项

如果无法获取有效的和风天气 API Key，可以：

**选项 1**：使用其他天气服务
- OpenWeather API
- AccuWeather API

**选项 2**：禁用天气功能（测试其他 AI 功能）
```yaml
ai:
  skills:
    weather:
      enabled: false
```

## 需要帮助？

查看完整的故障排除文档：[WEATHER_API_TROUBLESHOOTING.md](WEATHER_API_TROUBLESHOOTING.md)
