package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"

	"github.com/gin-gonic/gin"
)

// 设置首页模块路由
func SetupHomeRoutes(r *gin.RouterGroup) {
	// 创建处理器实例
	homeHandler := handler.NewHomeHandler()

	// 首页模块路由组
	home := r.Group("/home")
	home.Use(middleware.JWTAuth())
	{
		// 获取首页列表（不需要认证）
		home.GET("/list", homeHandler.List)
		// 创建模拟数据（不需要认证）
		home.POST("/create", homeHandler.Create)
		// 更新首页内容（不需要认证）
		home.PUT("/:id", homeHandler.Update)
	}
}
