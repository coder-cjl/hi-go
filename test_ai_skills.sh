#!/bin/bash

# AI Skills 功能测试脚本

BASE_URL="http://localhost:8000"

echo "======================================"
echo "AI Skills 功能测试"
echo "======================================"
echo ""

# 测试1: 查询单个城市天气
echo "测试1: 查询北京天气"
curl -X POST "${BASE_URL}/api/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{"message": "北京今天天气怎么样？"}' \
  | jq '.'

echo ""
echo "======================================"
echo ""

# 测试2: 查询多个城市
echo "测试2: 对比北京和上海天气"
curl -X POST "${BASE_URL}/api/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{"message": "比较一下北京和上海今天的天气"}' \
  | jq '.'

echo ""
echo "======================================"
echo ""

# 测试3: 咨询建议
echo "测试3: 咨询出行建议"
curl -X POST "${BASE_URL}/api/ai/chat" \
  -H "Content-Type: application/json" \
  -d '{"message": "明天去深圳出差，需要注意什么？"}' \
  | jq '.'

echo ""
echo "======================================"
echo "测试完成"
echo "======================================"
