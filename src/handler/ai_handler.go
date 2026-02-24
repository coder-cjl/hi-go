package handler

import (
	"hi-go/src/model"
	"hi-go/src/service/aiservice"

	"github.com/gin-gonic/gin"
)

// AIHandler AI处理器
type AIHandler struct {
	aiService *aiservice.Service
}

// NewAIHandler 创建AI处理器
func NewAIHandler(aiService *aiservice.Service) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

// ChatRequest 对话请求
type ChatRequest struct {
	Message string `json:"message" binding:"required" example:"北京今天天气怎么样？"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	Reply string `json:"reply" example:"北京今天天气晴朗，温度20°C..."`
}

// Chat 对话接口
// @Summary AI对话
// @Description 与AI进行对话，支持天气查询等技能
// @Tags AI
// @Accept json
// @Produce json
// @Param request body ChatRequest true "对话内容"
// @Success 200 {object} model.Response{data=ChatResponse}
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/ai/chat [post]
func (h *AIHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.aiService.Chat(c.Request.Context(), req.Message)
	if err != nil {
		model.ServerError(c, "AI服务调用失败: "+err.Error())
		return
	}

	model.Success(c, ChatResponse{
		Reply: result,
	})
}
