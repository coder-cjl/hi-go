package handler

import (
	"fmt"
	"hi-go/src/model"
	"hi-go/src/service/aiservice"
	"io"

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

// Chat2 流式对话接口 (SSE)
// @Summary AI流式对话
// @Description 与AI进行流式对话，使用SSE返回实时响应
// @Tags AI
// @Accept json
// @Produce text/event-stream
// @Param request body ChatRequest true "对话内容"
// @Success 200 {string} string "SSE流"
// @Failure 400 {object} model.Response
// @Failure 500 {object} model.Response
// @Router /api/v1/ai/chat2 [post]
func (h *AIHandler) Chat2(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "请求参数错误: "+err.Error())
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")
	c.Header("X-Accel-Buffering", "no") // 禁用nginx缓冲

	// 获取流式响应通道
	streamChan, err := h.aiService.ChatStream(c.Request.Context(), req.Message)
	if err != nil {
		// SSE错误格式
		c.SSEvent("error", fmt.Sprintf(`{"message": "%s"}`, err.Error()))
		c.Writer.Flush()
		return
	}

	// 向客户端发送流式响应
	c.Stream(func(w io.Writer) bool {
		if resp, ok := <-streamChan; ok {
			// 检查是否有错误
			if resp.Error != nil {
				c.SSEvent("error", fmt.Sprintf(`{"message": "%s"}`, resp.Error.Error()))
				return false
			}

			// 发送内容增量
			if resp.Content != "" {
				c.SSEvent("message", fmt.Sprintf(`{"content": "%s"}`, resp.Content))
			}

			// 发送工具调用（如果有）
			if len(resp.ToolCalls) > 0 {
				c.SSEvent("tool_calls", fmt.Sprintf(`{"tool_calls": %v}`, resp.ToolCalls))
			}

			// 如果完成，发送完成事件
			if resp.FinishReason != "" {
				c.SSEvent("done", fmt.Sprintf(`{"finish_reason": "%s"}`, resp.FinishReason))
				return false
			}

			c.Writer.Flush()
			return true
		}
		return false
	})
}
