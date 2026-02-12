# Logstash 日志收集集成说明

## 功能概述

项目已成功集成 Logstash 日志收集功能，实现了以下特性：

- ✅ TCP 实时日志传输
- ✅ 自动重连机制
- ✅ JSON 格式日志处理
- ✅ Logstash 过滤和转换
- ✅ 发送到 Elasticsearch
- ✅ Kibana 可视化

## 架构流程

```
┌─────────────┐      TCP:5000       ┌──────────────┐      ┌──────────────────┐      ┌─────────┐
│  Hi-Go 应用  │ ──────────────────> │   Logstash   │ ───> │  Elasticsearch   │ ───> │ Kibana  │
│             │   JSON logs         │              │      │                  │      │         │
└─────────────┘                     └──────────────┘      └──────────────────┘      └─────────┘
                                           │
                                           ├─ 过滤
                                           ├─ 转换
                                           └─ 增强
```

## 配置说明

### 1. 应用配置

在配置文件中（`configs/dev.yaml`, `configs/prod.yaml` 等），添加 Logstash 配置：

```yaml
# Elasticsearch 配置（Logstash 会将日志发送到这里）
elasticsearch:
  enabled: true
  addrs:
    - http://localhost:9200
  username: ""
  password: ""
  index: "hi-go-logs"
  max_retry: 3

# Logstash 日志收集配置
logstash:
  enabled: true          # 是否启用 Logstash 日志收集
  host: localhost        # Logstash 服务器地址
  port: 5000            # Logstash TCP 端口
  protocol: tcp         # 协议：tcp 或 udp
  timeout: 5            # 连接超时（秒）
  reconnect: true       # 是否自动重连
  buffer_size: 8192     # 缓冲区大小（字节）
```

### 2. Logstash 配置

#### logstash.conf

位置：`configs/logstash/logstash.conf`

```ruby
input {
  # TCP 输入 - 接收应用通过 TCP 发送的 JSON 日志
  tcp {
    port => 5000
    codec => json_lines
  }
}

filter {
  # 处理时间字段
  if [ts] {
    date {
      match => ["ts", "ISO8601"]
      target => "@timestamp"
      remove_field => ["ts"]
    }
  }

  # 添加应用标识
  mutate {
    add_field => {
      "app" => "hi-go"
    }
  }

  # 转换级别为大写
  if [level] {
    mutate {
      uppercase => ["level"]
    }
  }
}

output {
  # 输出到 Elasticsearch
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "hi-go-logs-%{+YYYY.MM.dd}"  # 按日期创建索引
  }
}
```

#### logstash.yml

位置：`configs/logstash/logstash.yml`

```yaml
http.host: "0.0.0.0"
http.port: 9600
log.level: info
path.logs: /var/log/logstash
```

### 3. 不同环境配置

#### 开发环境 (dev.yaml)
```yaml
logstash:
  enabled: false  # 默认关闭，可手动开启测试
  host: localhost
  port: 5000
  protocol: tcp
  timeout: 5
  reconnect: true
  buffer_size: 8192
```

#### UAT/生产环境 (uat.yaml, prod.yaml)
```yaml
logstash:
  enabled: true  # 建议启用
  host: prod-logstash-server
  port: 5000
  protocol: tcp
  timeout: 5
  reconnect: true
  buffer_size: 8192
```

## 快速开始

### 1. 启动 ELK Stack

使用提供的启动脚本：

```bash
./scripts/start-elk.sh
```

或手动启动：

```bash
docker-compose -f docker-compose-elk.yml up -d
```

这将启动：
- **Elasticsearch** - 端口 9200
- **Logstash** - 端口 5000 (TCP), 9600 (API)
- **Kibana** - 端口 5601

### 2. 配置应用

修改 `configs/dev.yaml`：

```yaml
logstash:
  enabled: true  # 启用 Logstash
  host: localhost
  port: 5000
```

### 3. 启动应用

```bash
GO_ENV=dev go run main.go
```

### 4. 验证日志流

#### 检查 Logstash 日志
```bash
docker logs hi-go-logstash -f
```

#### 检查 Elasticsearch 中的日志
```bash
# 查看索引
curl http://localhost:9200/_cat/indices?v

# 查看日志
curl http://localhost:9200/hi-go-logs-*/_search?pretty

# 搜索最新10条日志
curl -X POST http://localhost:9200/hi-go-logs-*/_search?pretty \
  -H 'Content-Type: application/json' \
  -d '{"size":10,"sort":[{"@timestamp":"desc"}]}'
```

#### 在 Kibana 中查看
访问 http://localhost:5601 并创建索引模式 `hi-go-logs-*`

## 特性详解

### 1. TCP 实时传输

日志通过 TCP 连接实时发送到 Logstash：
- 低延迟
- 可靠传输
- 保持连接

### 2. 自动重连

如果 Logstash 连接断开：
- 自动尝试重连
- 日志不会丢失（继续写入文件/控制台）
- 重连成功后恢复发送

### 3. JSON 格式

所有日志都以 JSON 格式发送：

```json
{
  "level": "INFO",
  "ts": "2026-02-12T10:30:45+08:00",
  "logger": "hi-go",
  "caller": "main.go:123",
  "msg": "应用启动中...",
  "env": "dev"
}
```

### 4. Logstash 处理

Logstash 可以对日志进行：
- **过滤** - 根据条件过滤日志
- **转换** - 修改字段格式
- **增强** - 添加额外字段
- **路由** - 发送到不同的输出

### 5. 索引管理

按日期创建索引：
- `hi-go-logs-2026.02.12`
- `hi-go-logs-2026.02.13`
- ...

便于：
- 按日期查询
- 定期清理旧数据
- 优化存储

## 高级配置

### 1. 自定义 Logstash 过滤器

在 `configs/logstash/logstash.conf` 中添加过滤器：

```ruby
filter {
  # GeoIP 处理（如果日志中有 IP 地址）
  if [client_ip] {
    geoip {
      source => "client_ip"
      target => "geoip"
    }
  }

  # User-Agent 解析
  if [user_agent] {
    useragent {
      source => "user_agent"
      target => "ua"
    }
  }

  # 删除不需要的字段
  mutate {
    remove_field => ["host", "port"]
  }

  # 添加自定义字段
  mutate {
    add_field => {
      "environment" => "production"
      "app_version" => "1.0.0"
    }
  }
}
```

### 2. 多输出配置

除了 Elasticsearch，还可以输出到其他地方：

```ruby
output {
  # 输出到 Elasticsearch
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "hi-go-logs-%{+YYYY.MM.dd}"
  }

  # 同时输出到文件（备份）
  file {
    path => "/var/log/logstash/hi-go-%{+YYYY-MM-dd}.log"
    codec => line { format => "%{message}" }
  }

  # 输出到 Kafka（可选）
  # kafka {
  #   bootstrap_servers => "kafka:9092"
  #   topic_id => "hi-go-logs"
  # }
}
```

### 3. 条件路由

根据日志级别路由到不同索引：

```ruby
output {
  if [level] == "ERROR" {
    elasticsearch {
      hosts => ["elasticsearch:9200"]
      index => "hi-go-errors-%{+YYYY.MM.dd}"
    }
  } else {
    elasticsearch {
      hosts => ["elasticsearch:9200"]
      index => "hi-go-logs-%{+YYYY.MM.dd}"
    }
  }
}
```

## 性能优化

### 1. 批量处理

Logstash 默认使用批量处理提高性能。可以在 `logstash.yml` 中调整：

```yaml
pipeline.batch.size: 125    # 批量大小
pipeline.batch.delay: 50    # 批量延迟（毫秒）
pipeline.workers: 2         # 工作线程数
```

### 2. 缓冲区大小

应用端可以调整缓冲区大小：

```yaml
logstash:
  buffer_size: 16384  # 增大缓冲区
```

### 3. 持久化队列

在 Logstash 中启用持久化队列（防止数据丢失）：

```yaml
# logstash.yml
queue.type: persisted
queue.max_bytes: 1gb
```

## Docker 环境配置

### 1. 应用在 Docker 中运行

如果应用也在 Docker 容器中，需要使用正确的网络配置：

```yaml
# docker-compose.yml
services:
  app:
    image: hi-go:latest
    networks:
      - elk
    environment:
      - LOGSTASH_HOST=logstash  # 使用服务名称
      - LOGSTASH_PORT=5000

networks:
  elk:
    external: true  # 使用 ELK Stack 的网络
```

配置文件：
```yaml
logstash:
  enabled: true
  host: logstash  # 容器名称
  port: 5000
```

### 2. 宿主机应用连接 Docker Logstash

```yaml
logstash:
  enabled: true
  host: localhost  # 宿主机地址
  port: 5000       # 映射的端口
```

## 故障排查

### 问题1：无法连接到 Logstash

**错误信息：**
```
Warning: initial connection to logstash failed: dial tcp localhost:5000: connect: connection refused
```

**解决方法：**
1. 检查 Logstash 是否运行：`docker ps | grep logstash`
2. 检查端口是否开放：`netstat -an | grep 5000`
3. 检查防火墙设置
4. 验证配置中的 host 和 port 是否正确

### 问题2：日志没有到达 Elasticsearch

**检查步骤：**

1. **检查 Logstash 是否接收到日志：**
```bash
docker logs hi-go-logstash -f
```

2. **检查 Logstash 管道状态：**
```bash
curl http://localhost:9600/_node/stats/pipelines?pretty
```

3. **检查 Logstash 配置是否正确：**
```bash
docker exec hi-go-logstash /usr/share/logstash/bin/logstash \
  --config.test_and_exit \
  -f /usr/share/logstash/pipeline/logstash.conf
```

4. **检查 Elasticsearch 是否运行：**
```bash
curl http://localhost:9200
```

### 问题3：连接频繁断开重连

**可能原因：**
- 网络不稳定
- Logstash 负载过高
- 超时设置过短

**解决方法：**
```yaml
logstash:
  timeout: 30  # 增加超时时间
  reconnect: true  # 确保启用自动重连
```

### 问题4：日志格式不正确

在 Logstash 中开启调试输出：

```ruby
output {
  stdout {
    codec => rubydebug
  }
}
```

重启 Logstash 并查看输出：
```bash
docker logs hi-go-logstash -f
```

## 监控和维护

### 1. Logstash 健康检查

```bash
# API 健康检查
curl http://localhost:9600

# 节点信息
curl http://localhost:9600/_node?pretty

# 管道统计
curl http://localhost:9600/_node/stats/pipelines?pretty
```

### 2. 性能监控

查看关键指标：
```bash
curl http://localhost:9600/_node/stats/jvm?pretty
curl http://localhost:9600/_node/stats/process?pretty
```

### 3. 日志轮转

Elasticsearch 索引按日期自动创建，可以设置 ILM 策略自动删除旧索引：

```json
PUT _ilm/policy/hi-go-logs-policy
{
  "policy": {
    "phases": {
      "hot": {
        "actions": {
          "rollover": {
            "max_age": "1d"
          }
        }
      },
      "delete": {
        "min_age": "30d",
        "actions": {
          "delete": {}
        }
      }
    }
  }
}
```

## 最佳实践

### 1. 环境隔离

为不同环境使用不同的索引前缀：
- 开发：`hi-go-dev-logs-*`
- 测试：`hi-go-test-logs-*`
- 生产：`hi-go-prod-logs-*`

修改 Logstash 配置：
```ruby
output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "hi-go-%{environment}-logs-%{+YYYY.MM.dd}"
  }
}
```

### 2. 结构化日志

应用中使用结构化日志：

```go
logger.Info("用户登录",
    zap.String("username", "admin"),
    zap.String("ip", "192.168.1.1"),
    zap.Int("user_id", 123))
```

在 Logstash 中更容易处理和搜索。

### 3. 敏感数据处理

在 Logstash 中过滤敏感信息：

```ruby
filter {
  # 移除密码字段
  mutate {
    remove_field => ["password", "token", "secret"]
  }

  # 脱敏处理
  if [credit_card] {
    mutate {
      gsub => ["credit_card", "\d{12}", "************"]
    }
  }
}
```

### 4. 错误告警

配合 Elasticsearch Watcher 或 Kibana Alert 设置告警：
- 错误日志数量超过阈值
- 特定错误类型出现
- 关键业务日志缺失

## 与 Elasticsearch Writer 对比

项目同时支持两种方式将日志发送到 Elasticsearch：

### 直接写入 ES（elasticsearch writer）
```
应用 → Elasticsearch
```
- ✅ 简单直接
- ✅ 低延迟
- ❌ 无法处理和转换日志
- ❌ 单点故障

### 通过 Logstash（推荐）
```
应用 → Logstash → Elasticsearch
```
- ✅ 强大的日志处理能力
- ✅ 可以过滤、转换、增强日志
- ✅ 支持多输出
- ✅ 缓冲机制，防止 ES 过载
- ❌ 多一跳延迟
- ❌ 需要额外的 Logstash 服务

**建议：**
- 开发环境：可以直接用 ES writer 快速调试
- 生产环境：使用 Logstash 获得更好的可靠性和灵活性

## 相关资源

- [Logstash 官方文档](https://www.elastic.co/guide/en/logs tash/current/index.html)
- [Logstash 配置示例](https://www.elastic.co/guide/en/logstash/current/config-examples.html)
- [Logstash 过滤器插件](https://www.elastic.co/guide/en/logstash/current/filter-plugins.html)
- [Elasticsearch 官方文档](https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html)

## 总结

Logstash 集成为项目提供：

✅ **集中日志收集** - 统一收集所有应用日志  
✅ **强大的处理能力** - 过滤、转换、增强  
✅ **可靠传输** - TCP + 自动重连  
✅ **灵活输出** - 支持多种目标  
✅ **易于扩展** - 丰富的插件生态  
✅ **生产就绪** - 稳定可靠

---

**快速命令：**
```bash
# 启动 ELK Stack
./scripts/start-elk.sh

# 启用 Logstash（修改配置）
# configs/dev.yaml -> logstash.enabled: true

# 启动应用
GO_ENV=dev go run main.go

# 查看日志
docker logs hi-go-logstash -f

# 查看 ES 索引
curl http://localhost:9200/_cat/indices?v

# 访问 Kibana
open http://localhost:5601
```
