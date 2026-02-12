# Elasticsearch 日志集成说明

## 功能概述

项目已成功集成 Elasticsearch 日志功能，实现了以下特性：

- ✅ 日志自动写入 Elasticsearch
- ✅ 批量写入，提高性能
- ✅ 异步刷新机制
- ✅ 多环境配置支持
- ✅ 与文件日志和控制台日志并行工作

## 配置说明

### 1. Elasticsearch 配置项

在配置文件中（`configs/dev.yaml`, `configs/prod.yaml` 等），添加以下配置：

```yaml
elasticsearch:
  enabled: true  # 是否启用 Elasticsearch 日志（true/false）
  addrs:
    - http://localhost:9200  # Elasticsearch 集群地址（可配置多个）
    - http://es-node2:9200
  username: "elastic"      # 用户名（可选，如果ES开启了安全认证）
  password: "password"     # 密码（可选）
  index: "hi-go-logs"      # 索引名称
  max_retry: 3             # 最大重试次数
```

### 2. 不同环境的配置

#### 开发环境 (dev.yaml)
```yaml
elasticsearch:
  enabled: false  # 开发环境默认关闭
  addrs:
    - http://localhost:9200
  username: ""
  password: ""
  index: "hi-go-logs-dev"
  max_retry: 3
```

#### 测试环境 (test.yaml)
```yaml
elasticsearch:
  enabled: false  # 测试环境可选
  addrs:
    - http://localhost:9200
  index: "hi-go-logs-test"
```

#### UAT环境 (uat.yaml)
```yaml
elasticsearch:
  enabled: true  # UAT环境建议启用
  addrs:
    - http://uat-es-server:9200
  username: "es_user"
  password: "es_password"
  index: "hi-go-logs-uat"
```

#### 生产环境 (prod.yaml)
```yaml
elasticsearch:
  enabled: true  # 生产环境建议启用
  addrs:
    - http://prod-es-server:9200
  username: "es_user"
  password: "es_password"  # 建议使用环境变量
  index: "hi-go-logs-prod"
```

## 使用方法

### 1. 启动 Elasticsearch

使用 Docker 快速启动（用于开发测试）：

```bash
docker run -d \
  --name elasticsearch \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  docker.elastic.co/elasticsearch/elasticsearch:8.11.0
```

### 2. 启动应用

修改配置文件启用 Elasticsearch：

```yaml
# configs/dev.yaml
elasticsearch:
  enabled: true  # 改为 true
  addrs:
    - http://localhost:9200
  index: "hi-go-logs-dev"
```

然后启动应用：

```bash
# 使用开发环境配置
GO_ENV=dev go run main.go

# 或使用其他环境
GO_ENV=prod go run main.go
```

### 3. 查看日志

在 Elasticsearch 中查看日志：

```bash
# 查看所有日志
curl http://localhost:9200/hi-go-logs-dev/_search?pretty

# 查看最新10条日志
curl http://localhost:9200/hi-go-logs-dev/_search?pretty -H 'Content-Type: application/json' -d '
{
  "size": 10,
  "sort": [{"@timestamp": "desc"}]
}'

# 搜索特定级别的日志
curl http://localhost:9200/hi-go-logs-dev/_search?pretty -H 'Content-Type: application/json' -d '
{
  "query": {
    "match": {"level": "error"}
  }
}'
```

## 日志数据结构

写入 Elasticsearch 的日志包含以下字段：

```json
{
  "@timestamp": "2026-02-12T10:30:45+08:00",
  "level": "info",
  "logger": "hi-go",
  "caller": "main.go:123",
  "message": "应用启动中...",
  "fields": {
    "env": "dev",
    "custom_field": "custom_value"
  }
}
```

## 性能优化

### 批量写入

日志采用批量写入机制，默认配置：

- **批量大小**：100条日志
- **刷新间隔**：5秒

这意味着：
1. 当缓冲区累积100条日志时，自动批量写入 ES
2. 如果日志量较少，每5秒自动刷新一次
3. 应用退出时会自动刷新剩余日志

### 异步写入

日志写入是异步的，不会阻塞主业务逻辑：
- 控制台日志：实时输出
- 文件日志：实时写入
- Elasticsearch日志：批量异步写入

## 高可用配置

### 1. 集群配置

配置多个 ES 节点实现高可用：

```yaml
elasticsearch:
  enabled: true
  addrs:
    - http://es-node1:9200
    - http://es-node2:9200
    - http://es-node3:9200
  max_retry: 5  # 增加重试次数
```

### 2. 连接失败处理

如果 Elasticsearch 连接失败：
- 应用**仍然会正常启动**
- 日志会继续输出到控制台和文件
- 错误信息会输出到 stderr

## 使用 Kibana 可视化

### 1. 启动 Kibana

```bash
docker run -d \
  --name kibana \
  --link elasticsearch:elasticsearch \
  -p 5601:5601 \
  docker.elastic.co/kibana/kibana:8.11.0
```

### 2. 访问 Kibana

访问 http://localhost:5601

### 3. 创建索引模式

1. 打开 Kibana
2. 进入 Management → Stack Management → Index Patterns
3. 创建索引模式：`hi-go-logs-*`
4. 选择时间字段：`@timestamp`
5. 进入 Discover 查看日志

## 故障排查

### 问题1：Elasticsearch 连接失败

**错误信息**：
```
Failed to initialize Elasticsearch writer: failed to connect to elasticsearch
```

**解决方法**：
1. 检查 ES 是否正常运行：`curl http://localhost:9200`
2. 检查配置文件中的地址是否正确
3. 检查网络连接和防火墙设置

### 问题2：日志没有写入

**可能原因**：
1. `enabled: false` - 检查配置是否启用
2. ES 连接失败 - 查看启动日志
3. 缓冲区未满且刷新间隔未到 - 等待5秒或重启应用

### 问题3：日志写入延迟

这是正常现象，因为采用批量写入：
- 最多延迟5秒（刷新间隔）
- 如需立即写入，可调整配置（不推荐，影响性能）

## 最佳实践

### 1. 环境区分

为不同环境使用不同的索引名称：
- 开发环境：`hi-go-logs-dev`
- 测试环境：`hi-go-logs-test`
- UAT环境：`hi-go-logs-uat`
- 生产环境：`hi-go-logs-prod`

### 2. 使用环境变量

生产环境的敏感信息建议使用环境变量：

```bash
export ES_USERNAME=your_username
export ES_PASSWORD=your_password
```

然后在代码中读取环境变量覆盖配置。

### 3. 索引生命周期管理

建议配置 ILM（Index Lifecycle Management）：
- 设置索引滚动策略（如每天创建新索引）
- 设置数据保留期（如保留30天）
- 自动删除旧数据

### 4. 监控告警

配置 Kibana 告警规则：
- 错误日志数量超过阈值
- 特定错误类型出现
- 日志写入速率异常

## 代码示例

日志使用方式不变，自动同时写入多个目标：

```go
// 结构化日志
logger.Info("用户登录", 
    zap.String("username", "admin"),
    zap.String("ip", "192.168.1.1"))

// 格式化日志  
logger.Infof("处理订单: %s", orderId)

// 错误日志
logger.Error("数据库连接失败", 
    zap.Error(err),
    zap.String("database", "users"))
```

这些日志会自动：
1. 输出到控制台（彩色格式）
2. 写入日志文件（JSON格式）
3. 批量写入 Elasticsearch（JSON格式）

## 依赖包

项目添加了以下依赖：

```
github.com/elastic/go-elasticsearch/v8 v8.19.3
github.com/elastic/elastic-transport-go/v8 v8.8.0
```

## 总结

Elasticsearch 日志集成为项目提供：

✅ **集中化日志管理** - 所有日志统一存储在 ES  
✅ **强大的搜索能力** - 快速检索和分析日志  
✅ **可视化分析** - 通过 Kibana 查看日志趋势  
✅ **高性能** - 批量异步写入，不影响业务  
✅ **高可用** - 支持 ES 集群，自动重试  
✅ **灵活配置** - 不同环境独立配置

---

如有问题，请参考：
- [Elasticsearch 官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)
- [go-elasticsearch GitHub](https://github.com/elastic/go-elasticsearch)
