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
		// 获取首页列表
		home.GET("/list", homeHandler.List)
		// 创建模拟数据
		home.POST("/create", homeHandler.Create)
		// 更新首页内容
		home.POST("/update", homeHandler.Update)
		// 删除首页内容
		home.DELETE("/delete", homeHandler.Delete)
		// 搜索首页内容
		home.GET("/search", homeHandler.Search)
	}
}
