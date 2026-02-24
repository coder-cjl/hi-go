package router

import (
	"hi-go/src/handler"

	"github.com/gin-gonic/gin"
)

// RegisterAIRoutes 注册AI相关路由
func RegisterAIRoutes(r *gin.RouterGroup, aiHandler *handler.AIHandler) {
	if aiHandler == nil {
		return
	}

	ai := r.Group("/ai")
	{
		ai.POST("/chat", aiHandler.Chat)
	}
}
