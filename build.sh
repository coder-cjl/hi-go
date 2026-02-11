#!/bin/bash

# Hi-Go 应用构建脚本
# 用法: ./build.sh [env]
# env: dev, test, uat, prod (默认: prod)

ENV=${1:-prod}
OUTPUT="hi-go-${ENV}"

echo "========================================="
echo "  Hi-Go Application Builder"
echo "========================================="
echo "Environment: $ENV"
echo "Output: $OUTPUT"
echo "========================================="

# 检查配置文件是否存在
if [ ! -f "configs/${ENV}.yaml" ]; then
    echo "错误: 配置文件 configs/${ENV}.yaml 不存在"
    exit 1
fi

# 构建应用
echo "正在构建应用..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $OUTPUT main.go

if [ $? -eq 0 ]; then
    echo "========================================="
    echo "构建成功！"
    echo "输出文件: $OUTPUT"
    echo "运行方式: GO_ENV=$ENV ./$OUTPUT"
    echo "========================================="
else
    echo "构建失败"
    exit 1
fi
