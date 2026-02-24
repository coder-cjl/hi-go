#!/bin/bash

echo "=== 测试AI天气查询功能 ==="
echo ""

# 测试查询北京天气
echo "测试: 北京今天天气怎么样？"
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "北京今天天气怎么样？"}' \
  2>/dev/null | jq '.'

echo ""
echo "---"
echo ""

# 测试查询上海天气
echo "测试: 上海的天气如何？"
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "上海的天气如何？"}' \
  2>/dev/null | jq '.'

echo ""
echo "---"
echo ""

# 测试查询深圳天气
echo "测试: 帮我查一下深圳的天气"
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "帮我查一下深圳的天气"}' \
  2>/dev/null | jq '.'
