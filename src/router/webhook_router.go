package router

import (
	"hi-go/src/handler"
	"hi-go/src/middleware"

	"github.com/gin-gonic/gin"
)

// 设置 Webhook 模块路由
func SetupWebhookRoutes(r *gin.RouterGroup) {
	// 创建处理器实例
	webhookHandler := handler.NewWebhookHandler()

	// Webhook 模块路由组
	webhook := r.Group("/webhook")
	// 需要 JWT 认证的路由
	webhook.Use(middleware.JWTAuth())
	{
		// 创建 webhook
		webhook.POST("/create", webhookHandler.Create)
		// 更新 webhook
		webhook.POST("/update", webhookHandler.Update)
		// 删除 webhook
		webhook.DELETE("/delete/:id", webhookHandler.Delete)
		// 根据 ID 获取 webhook 详情
		webhook.GET("/detail/:id", webhookHandler.GetByID)
		// 获取 webhook 列表
		webhook.GET("/list", webhookHandler.GetList)
		// 生成签名
		webhook.POST("/sign", webhookHandler.Sign)
	}

	// 公开的 webhook 回调路由（不需要认证）
	webhook.POST("/callback/:secret", webhookHandler.Callback)
}
