# 构建阶段
FROM golang:1.23-alpine AS builder

# 安装必要的构建工具
RUN apk add --no-cache git gcc musl-dev

# 设置工作目录
WORKDIR /build

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o hi-go main.go

# 运行阶段
FROM alpine:latest

# 安装 ca-certificates 用于 HTTPS 请求
RUN apk --no-cache add ca-certificates tzdata

# 设置时区为上海
ENV TZ=Asia/Shanghai

# 创建应用目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /build/hi-go .

# 复制配置文件
COPY --from=builder /build/configs ./configs

# 复制 Swagger 文档（如果存在）
COPY --from=builder /build/docs ./docs

# 创建日志目录
RUN mkdir -p logs

# 暴露端口
EXPOSE 8000

# 设置环境变量
ENV GO_ENV=dev

# 运行应用
CMD ["./hi-go"]
