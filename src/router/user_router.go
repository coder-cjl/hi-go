package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"

	"github.com/gin-gonic/gin"
)

// 设置用户模块路由
func SetupUserRoutes(r *gin.RouterGroup) {
	// 创建处理器实例
	authHandler := handler.NewAuthHandler()

	// 用户模块路由组
	user := r.Group("/user")
	{
		// 公开接口（不需要认证）
		// 登录
		user.POST("/login", authHandler.Login)
		// 注册
		user.POST("/register", authHandler.Register)

		// 需要认证的接口
		// JWT 认证中间件
		auth := user.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 获取个人信息
			auth.GET("/profile", authHandler.GetProfile)
		}

		// 管理员接口
		admin := user.Group("/admin")
		// JWT 认证
		admin.Use(middleware.JWTAuth())
		// 角色权限认证
		admin.Use(middleware.RoleAuth("admin"))
	}
}
