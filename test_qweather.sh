#!/bin/bash

# 和风天气API测试脚本
API_KEY="1e9f39dd5a4d4a299064818401b0652b"
BASE_URL="https://devapi.qweather.com/v7"

echo "=== 测试和风天气API ==="
echo ""

# 测试1: 使用中文城市名（URL编码）
echo "测试1: 使用URL编码的中文城市名 '北京'"
LOCATION=$(printf '%s' '北京' | jq -sRr @uri)
echo "URL: ${BASE_URL}/weather/now?location=${LOCATION}&key=${API_KEY}"
curl -s "${BASE_URL}/weather/now?location=${LOCATION}&key=${API_KEY}" | jq '.'
echo ""

# 测试2: 使用拼音
echo "测试2: 使用拼音 'beijing'"
curl -s "${BASE_URL}/weather/now?location=beijing&key=${API_KEY}" | jq '.'
echo ""

# 测试3: 使用城市ID (北京的ID通常是101010100)
echo "测试3: 使用城市ID '101010100'"
curl -s "${BASE_URL}/weather/now?location=101010100&key=${API_KEY}" | jq '.'
echo ""

# 测试4: 验证API key
echo "测试4: 测试API key有效性（使用错误的key）"
curl -s "${BASE_URL}/weather/now?location=beijing&key=invalid_key_test" | jq '.'
echo ""
