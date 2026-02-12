#!/bin/bash

# ELK Stack 快速启动脚本
# 用于启动 Elasticsearch + Logstash + Kibana

set -e

echo "======================================"
echo "  启动 ELK Stack"
echo "  Elasticsearch + Logstash + Kibana"
echo "======================================"
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 启动服务
echo "🚀 正在启动 ELK Stack..."
docker-compose -f docker-compose-elk.yml up -d

echo ""
echo "⏳ 等待服务启动..."
sleep 10

# 检查 Elasticsearch 状态
echo ""
echo "📊 检查 Elasticsearch 状态..."
for i in {1..30}; do
    if curl -s http://localhost:9200 > /dev/null 2>&1; then
        echo "✅ Elasticsearch 已就绪"
        break
    fi
    echo "   等待 Elasticsearch 启动... ($i/30)"
    sleep 2
done

# 检查 Logstash 状态
echo ""
echo "📊 检查 Logstash 状态..."
for i in {1..30}; do
    if curl -s http://localhost:9600 > /dev/null 2>&1; then
        echo "✅ Logstash 已就绪"
        break
    fi
    echo "   等待 Logstash 启动... ($i/30)"
    sleep 2
done

# 检查 Kibana 状态
echo ""
echo "📊 检查 Kibana 状态..."
for i in {1..60}; do
    if curl -s http://localhost:5601/api/status > /dev/null 2>&1; then
        echo "✅ Kibana 已就绪"
        break
    fi
    echo "   等待 Kibana 启动... ($i/60)"
    sleep 2
done

echo ""
echo "======================================"
echo "  ✅ ELK Stack 启动成功！"
echo "======================================"
echo ""
echo "服务访问地址："
echo "  - Elasticsearch: http://localhost:9200"
echo "  - Logstash API:  http://localhost:9600"
echo "  - Kibana:        http://localhost:5601"
echo ""
echo "Logstash 接收地址："
echo "  - TCP 端口:      localhost:5000"
echo ""
echo "快速操作："
echo "  - 查看服务状态： docker-compose -f docker-compose-elk.yml ps"
echo "  - 查看 Logstash 日志： docker logs hi-go-logstash -f"
echo "  - 查看所有日志： docker-compose -f docker-compose-elk.yml logs -f"
echo "  - 停止服务：     docker-compose -f docker-compose-elk.yml down"
echo "  - 删除数据：     docker-compose -f docker-compose-elk.yml down -v"
echo ""
echo "现在可以启动应用并启用 Logstash："
echo "  1. 修改 configs/dev.yaml 中 logstash.enabled: true"
echo "  2. 运行: GO_ENV=dev go run main.go"
echo ""
