#!/bin/bash

echo "🧪 测试SSE工具调用修复"
echo ""
echo "测试场景：上海的天气适合跑步吗？"
echo "预期：AI会先查询天气，然后基于查询结果给出建议"
echo ""
echo "按Ctrl+C停止"
echo "======================================"
echo ""

curl -X POST http://localhost:8080/api/ai/chat2 \
  -H "Content-Type: application/json" \
  -d '{"message": "上海的天气适合跑步吗？"}' \
  -N

echo ""
echo ""
echo "======================================"
echo "✅ 如果看到了：1) AI说要查询天气 2) 完整的天气分析建议，说明修复成功！"
echo "❌ 如果只看到 AI说要查询天气就结束了，说明还有问题"
