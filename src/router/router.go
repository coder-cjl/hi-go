package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"
	"hi-go/src/service/aiservice"

	"github.com/gin-gonic/gin"
)

// 设置所有路由
func Setup() *gin.Engine {
	// 创建 Gin 引擎
	r := gin.New()

	// 全局中间件
	// 错误恢复
	r.Use(middleware.Recovery())
	// 事务ID生成
	r.Use(middleware.TransID())
	// 日志记录
	r.Use(middleware.Logger())
	// 跨域处理
	r.Use(middleware.CORS())

	// 设置 API 文档路由
	SetupDocsRoutes(r)

	// API 路由组
	api := r.Group("/api")
	{
		// 用户模块路由
		SetupUserRoutes(api)
		// 首页模块路由
		SetupHomeRoutes(api)
		// AI模块路由（如果AI服务已启用）
		if aiservice.GlobalService != nil {
			aiHandler := handler.NewAIHandler(aiservice.GlobalService)
			RegisterAIRoutes(api, aiHandler)
		}
	}

	return r
}
