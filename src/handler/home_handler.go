package handler

import (
	"hi-go/src/model"
	"hi-go/src/service"

	"github.com/gin-gonic/gin"
)

// HomeHandler 首页处理器
type HomeHandler struct {
	homeService *service.HomeService
}

// NewHomeHandler 创建首页处理器实例
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{
		homeService: service.NewHomeService(),
	}
}

// List 获取首页列表
func (h *HomeHandler) List(c *gin.Context) {
	var req model.HomeListRequest

	// 1. 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层获取列表
	resp, err := h.homeService.GetList(&req)
	if err != nil {
		model.ServerError(c, "获取列表失败: "+err.Error())
		return
	}

	// 3. 返回成功响应
	model.Success(c, resp)
}

// Create 创建模拟数据
func (h *HomeHandler) Create(c *gin.Context) {
	// 调用服务层创建30条模拟数据
	if err := h.homeService.CreateMockData(); err != nil {
		model.ServerError(c, "创建数据失败: "+err.Error())
		return
	}

	model.SuccessWithMessage(c, "成功创建30条模拟数据", nil)
}
