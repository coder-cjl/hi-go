package middleware

import (
	"bytes"
	"encoding/json"
	"hi-go/src/utils/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseWriter 包装器，用于捕获响应内容
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 读取请求体
		var requestBody string
		if c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = string(bodyBytes)
				// 恢复请求体，以便后续处理器可以读取
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// 包装 ResponseWriter 以捕获响应内容
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 记录请求日志
		duration := time.Since(startTime)

		// 获取响应内容
		responseBody := blw.body.String()

		// 尝试格式化 JSON 响应（如果是 JSON）
		var responseJSON interface{}
		if err := json.Unmarshal(blw.body.Bytes(), &responseJSON); err == nil {
			// 是有效的 JSON，记录格式化后的
			responseBody = blw.body.String()
		}

		// 构建日志字段
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 添加请求体（如果有）
		if requestBody != "" {
			fields = append(fields, zap.String("request_body", requestBody))
		}

		// 添加响应体（限制长度，避免日志过大）
		if len(responseBody) > 0 {
			maxLen := 1000 // 最大记录1000字符
			if len(responseBody) > maxLen {
				fields = append(fields, zap.String("response_body", responseBody[:maxLen]+"..."))
				fields = append(fields, zap.Int("response_body_size", len(responseBody)))
			} else {
				fields = append(fields, zap.String("response_body", responseBody))
			}
		}

		// 添加错误信息（如果有）
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		logger.Info("HTTP 请求", fields...)
	}
}
