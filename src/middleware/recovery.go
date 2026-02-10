package middleware

import (
	"hi-go/src/model"
	"hi-go/src/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery 错误恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录 panic 日志
				logger.Error("发生 panic",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				// 返回错误响应
				model.ServerError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()
	}
}
