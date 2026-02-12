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
// @Summary      获取首页内容列表
// @Description  分页获取首页内容列表，包括标题、描述、图片等信息
// @Tags         首页模块
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "页码（默认1）"
// @Param        page_size  query     int  false  "每页数量（默认20，最大100）"
// @Success      200        {object}  model.Response{data=model.HomeListDataResponse}  "获取成功"
// @Failure      400        {object}  model.Response  "参数错误"
// @Failure      500        {object}  model.Response  "服务器错误"
// @Router       /home/list [get]
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
// @Summary      创建模拟首页数据
// @Description  自动生成30条模拟首页数据，用于测试
// @Tags         首页模块
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response  "创建成功"
// @Failure      500  {object}  model.Response  "服务器错误"
// @Router       /home/create [post]
func (h *HomeHandler) Create(c *gin.Context) {
	// 调用服务层创建30条模拟数据
	if err := h.homeService.CreateMockData(); err != nil {
		model.ServerError(c, "创建数据失败: "+err.Error())
		return
	}

	model.SuccessWithMessage(c, "成功创建30条模拟数据", nil)
}

// Update 更新首页内容
// @Summary      更新首页内容
// @Description  根据ID更新首页内容信息
// @Tags         首页模块
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true  "首页内容ID"
// @Param        request  body      model.HomeUpdateRequest    true  "更新请求参数"
// @Success      200      {object}  model.Response  "更新成功"
// @Failure      400      {object}  model.Response  "参数错误"
// @Failure      404      {object}  model.Response  "记录不存在"
// @Failure      500      {object}  model.Response  "服务器错误"
// @Router       /home/update [post]
func (h *HomeHandler) Update(c *gin.Context) {
	// 2. 绑定请求参数
	var req model.HomeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 3. 调用服务层更新
	if err := h.homeService.Update(&req); err != nil {
		model.ServerError(c, "更新失败: "+err.Error())
		return
	}

	model.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除首页内容
// @Summary      删除首页内容
// @Description  根据ID删除首页内容信息
// @Tags         首页模块
// @Accept       json
// @Produce      json
// @Param        request  body      model.HomeDeleteRequest  true  "删除请求参数"
// @Success      200      {object}  model.Response  "删除成功"
// @Failure      400      {object}  model.Response  "参数错误"
// @Failure      404      {object}  model.Response  "记录不存在"
// @Failure      500      {object}  model.Response  "服务器错误"
// @Router       /home/delete [delete]
func (h *HomeHandler) Delete(c *gin.Context) {
	// 1. 绑定请求参数
	var req model.HomeDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层删除
	if err := h.homeService.Delete(&req); err != nil {
		model.ServerError(c, "删除失败: "+err.Error())
		return
	}

	model.SuccessWithMessage(c, "删除成功", nil)
}

// Search 搜索首页内容
// @Summary      搜索首页内容
// @Description  根据关键词搜索首页标题或描述
// @Tags         首页模块
// @Accept       json
// @Produce      json
// @Param        keyword   query     string  false  "搜索关键词"
// @Param        page      query     int     false  "页码（默认1）"
// @Param        page_size query     int     false  "每页数量（默认20，最大100）"
// @Success      200       {object}  model.Response{data=model.HomeListDataResponse}  "搜索成功"
// @Failure      400       {object}  model.Response  "参数错误"
// @Failure      500       {object}  model.Response  "服务器错误"
// @Router       /home/search [get]
func (h *HomeHandler) Search(c *gin.Context) {
	var req model.HomeSearchRequest

	// 1. 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层搜索
	resp, err := h.homeService.Search(&req)
	if err != nil {
		model.ServerError(c, "搜索失败: "+err.Error())
		return
	}

	// 3. 返回成功响应
	model.Success(c, resp)
}
