package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 设置用户模块路由
func SetupUserRoutes(r *gin.RouterGroup) {
	// 创建处理器实例
	authHandler := handler.NewAuthHandler()

	// 用户模块路由组
	user := r.Group("/user")
	{
		// 公开接口（不需要认证）
		user.POST("/login", authHandler.Login)       // 登录
		user.POST("/register", authHandler.Register) // 注册

		// 需要认证的接口
		auth := user.Group("")
		auth.Use(middleware.JWTAuth()) // JWT 认证中间件
		{
			auth.GET("/profile", authHandler.GetProfile)       // 获取个人信息
			auth.PUT("/profile", authHandler.UpdateProfile)    // 更新个人信息（预留）
			auth.POST("/change-password", authHandler.ChangePassword) // 修改密码（预留）
		}

		// 管理员接口
		admin := user.Group("/admin")
		admin.Use(middleware.JWTAuth())          // JWT 认证
		admin.Use(middleware.RoleAuth("admin")) // 角色权限检查
		{
			admin.GET("/list", authHandler.GetUserList)       // 获取用户列表（预留）
			admin.GET("/:id", authHandler.GetUserByID)        // 获取指定用户信息（预留）
			admin.PUT("/:id", authHandler.UpdateUser)         // 更新用户信息（预留）
			admin.DELETE("/:id", authHandler.DeleteUser)      // 删除用户（预留）
			admin.PUT("/:id/status", authHandler.UpdateUserStatus) // 更新用户状态（预留）
		}
	}
}
