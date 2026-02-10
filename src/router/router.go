package router

import (
	"hi-go/src/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup 设置所有路由
func Setup() *gin.Engine {
	// 创建 Gin 引擎
	r := gin.New()

	// 全局中间件
	r.Use(middleware.Recovery()) // 错误恢复
	r.Use(middleware.TransID())  // 事务ID生成
	r.Use(middleware.Logger())   // 日志记录
	r.Use(middleware.CORS())     // 跨域处理

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务运行正常",
		})
	})

	// API 路由组
	api := r.Group("/api")
	{
		// 注册各模块路由
		SetupUserRoutes(api) // 用户模块
	}

	return r
}
