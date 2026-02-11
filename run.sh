#!/bin/bash

# Hi-Go 应用启动脚本
# 用法: ./run.sh [env]
# env: dev, test, uat, prod (默认: dev)

ENV=${1:-dev}

echo "========================================="
echo "  Hi-Go Application Starter"
echo "========================================="
echo "Environment: $ENV"
echo "========================================="

# 检查配置文件是否存在
if [ ! -f "configs/${ENV}.yaml" ]; then
    echo "错误: 配置文件 configs/${ENV}.yaml 不存在"
    exit 1
fi

# 设置环境变量
export GO_ENV=$ENV

# 运行应用
echo "正在启动应用..."
go run main.go
