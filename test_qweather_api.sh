#!/bin/bash

# 测试和风天气 API 是否正常工作

API_KEY="HE2409021346141011"
BASE_URL="https://devapi.qweather.com/v7"
CITY="北京"

echo "======================================"
echo "测试和风天气 API"
echo "======================================"
echo ""

# 测试1: 查询城市ID
echo "1. 查询城市ID: ${CITY}"
CITY_RESPONSE=$(curl -s "${BASE_URL}/city/lookup?location=${CITY}&key=${API_KEY}")
echo "响应: ${CITY_RESPONSE}"
echo ""

# 提取城市ID (使用jq工具，如果没有就显示原始响应)
if command -v jq &> /dev/null; then
    CITY_ID=$(echo "${CITY_RESPONSE}" | jq -r '.location[0].id // empty')
    CODE=$(echo "${CITY_RESPONSE}" | jq -r '.code // empty')
    
    echo "状态码: ${CODE}"
    echo "城市ID: ${CITY_ID}"
    echo ""
    
    if [ -z "${CITY_ID}" ] || [ "${CODE}" != "200" ]; then
        echo "❌ 城市查询失败！"
        echo "可能的原因："
        echo "1. API Key 无效或过期"
        echo "2. API 请求次数已用完"
        echo "3. API URL 错误"
        exit 1
    fi
    
    echo "✅ 城市查询成功"
    echo ""
    
    # 测试2: 查询实时天气
    echo "2. 查询实时天气"
    WEATHER_RESPONSE=$(curl -s "${BASE_URL}/weather/now?location=${CITY_ID}&key=${API_KEY}")
    echo "响应: ${WEATHER_RESPONSE}"
    echo ""
    
    WEATHER_CODE=$(echo "${WEATHER_RESPONSE}" | jq -r '.code // empty')
    TEMP=$(echo "${WEATHER_RESPONSE}" | jq -r '.now.temp // empty')
    TEXT=$(echo "${WEATHER_RESPONSE}" | jq -r '.now.text // empty')
    
    echo "状态码: ${WEATHER_CODE}"
    echo "温度: ${TEMP}°C"
    echo "天气: ${TEXT}"
    echo ""
    
    if [ -z "${TEMP}" ] || [ "${WEATHER_CODE}" != "200" ]; then
        echo "❌ 天气查询失败！"
        exit 1
    fi
    
    echo "✅ 天气查询成功"
else
    echo "提示: 安装 jq 工具可以看到更详细的解析结果"
    echo "macOS: brew install jq"
    echo "Linux: apt-get install jq"
fi

echo ""
echo "======================================"
echo "测试完成"
echo "======================================"
