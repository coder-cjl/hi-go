### swg 本地查看
http://localhost:8000/swagger/index.html

## 日志系统

### Logstash 日志收集

项目集成了 Logstash 用于集中日志收集和处理。

**快速启动 ELK Stack：**

```bash
# 使用启动脚本
./scripts/start-elk.sh

# 或使用 docker-compose
docker-compose -f docker-compose-elk.yml up -d
```

**访问地址：**
- Elasticsearch: http://localhost:9200
- Logstash API: http://localhost:9600
- Kibana: http://localhost:5601

**配置说明：**

在 `configs/dev.yaml` 中配置：

```yaml
# 启用 Elasticsearch（Logstash 会将日志发送到 ES）
elasticsearch:
  enabled: true
  addrs:
    - http://localhost:9200
  index: "hi-go-logs"

# 启用 Logstash
logstash:
  enabled: true
  host: localhost
  port: 5000
  protocol: tcp
  timeout: 5
  reconnect: true
  buffer_size: 8192
```

**功能特性：**
- ✅ TCP 实时日志传输
- ✅ 自动重连机制
- ✅ JSON 格式日志
- ✅ Logstash 过滤和处理
- ✅ 发送到 Elasticsearch
- ✅ Kibana 可视化

**日志流程：**
```
应用 → Logstash (TCP:5000) → Elasticsearch → Kibana
```
