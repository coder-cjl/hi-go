package handler

import (
	"hi-go/src/model"
	"hi-go/src/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器实例
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: service.NewAuthService(),
	}
}

// Login 登录接口
// @Summary 用户登录
// @Description 通过用户名和密码登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body model.LoginRequest true "登录参数"
// @Success 200 {object} model.Resp{data=model.LoginResponse}
// @Router /api/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层登录
	resp, err := h.authService.Login(&req)
	if err != nil {
		model.Unauthorized(c, err.Error())
		return
	}

	// 3. 返回成功响应
	model.Success(c, resp)
}

// Register 注册接口
// @Summary 用户注册
// @Description 注册新用户
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body model.RegisterRequest true "注册参数"
// @Success 200 {object} model.Resp{data=model.User}
// @Router /api/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	// 1. 绑定并验证请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 调用服务层注册
	user, err := h.authService.Register(&req)
	if err != nil {
		model.ParamError(c, err.Error())
		return
	}

	// 3. 返回成功响应
	model.SuccessWithMessage(c, "注册成功", user)
}

// GetProfile 获取用户信息
// @Summary 获取当前用户信息
// @Description 获取已登录用户的个人信息
// @Tags 用户
// @Produce json
// @Success 200 {object} model.Resp{data=model.User}
// @Router /api/profile [get]
// @Security Bearer
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// 从上下文获取用户ID（中间件已设置）
	userID, exists := c.Get("userID")
	if !exists {
		model.Unauthorized(c, "未授权")
		return
	}

	// 将 string 类型的 userID 转换为 uint
	userIDStr, ok := userID.(string)
	if !ok {
		model.ParamError(c, "用户ID格式错误")
		return
	}

	userIDUint, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		model.ParamError(c, "用户ID无效")
		return
	}

	// 获取用户信息
	user, err := h.authService.GetUserByID(uint(userIDUint))
	if err != nil {
		model.NotFound(c, err.Error())
		return
	}

	model.Success(c, user)
}
