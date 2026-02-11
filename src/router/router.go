package router

import (
	"hi-go/src/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "hi-go/docs" // 引入生成的docs包
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

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API 路由组
	api := r.Group("/api")
	{
		// 用户模块路由
		SetupUserRoutes(api)
		// 首页模块路由
		SetupHomeRoutes(api)
	}

	return r
}
