#!/bin/bash
# Hi-Go Docker 镜像构建脚本
# 用法: ./docker.sh [env]
# env: dev, test, uat, prod (默认: prod)
ENV=${1:-dev}
IMAGE_NAME="hi-go:${ENV}"

echo "========================================="
echo "  Hi-Go Docker Image Builder"
echo "========================================="
echo "Environment: $ENV"
echo "Image Name: $IMAGE_NAME"
echo "=========================================" 

# 检查配置文件是否存在
if [ ! -f "configs/${ENV}.yaml" ]; then
    echo "错误: 配置文件 configs/${ENV}.yaml 不存在"
    exit 1
fi

# 构建 Docker 镜像
echo "正在构建 Docker 镜像..."
docker build --build-arg GO_ENV=$ENV -t $IMAGE_NAME .
if [ $? -eq 0 ]; then
    echo "========================================="
    echo "Docker 镜像构建成功！"
    echo "镜像名称: $IMAGE_NAME"
    echo "运行方式: docker run -e GO_ENV=$ENV -p 8080:8080 $IMAGE_NAME"
    echo "========================================="
else
    echo "Docker 镜像构建失败"
    exit 1
fi