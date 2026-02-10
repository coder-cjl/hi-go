package middleware

import (
	"hi-go/src/model"

	"github.com/gin-gonic/gin"
)

// TransID 为每个请求生成唯一的追踪ID
func TransID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从请求头获取 TraceID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			// 如果请求头中没有，则生成新的
			traceID = model.GenerateTransID()
		}

		// 将 TraceID 存储到 Context 中
		c.Set(model.TraceIDKey, traceID)

		// 将 TraceID 设置到响应头中，方便客户端追踪
		c.Header("X-Trace-ID", traceID)

		c.Next()
	}
}
