package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"

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

	// 创建处理器实例
	authHandler := handler.NewAuthHandler()

	// 公开路由组
	api := r.Group("/api")
	{
		// 认证相关
		api.POST("/login", authHandler.Login)       // 登录
		api.POST("/register", authHandler.Register) // 注册

		// 需要认证的路由
		authorized := api.Group("")
		authorized.Use(middleware.JWTAuth()) // JWT 认证中间件
		{
			authorized.GET("/profile", authHandler.GetProfile) // 获取个人信息

			// 管理员路由
			admin := authorized.Group("/admin")
			admin.Use(middleware.RoleAuth("admin")) // 角色权限检查
			{
				// 这里添加管理员专属接口
				admin.GET("/users", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "管理员用户列表"})
				})
			}
		}
	}

	return r
}
