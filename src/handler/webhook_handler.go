package handler

import (
	"encoding/json"
	"fmt"
	"hi-go/src/model"
	"hi-go/src/service"
	"hi-go/src/utils/logger"

	"github.com/gin-gonic/gin"
)

// WebhookHandler Webhook 处理器
type WebhookHandler struct {
	webhookService *service.WebhookService
}

// NewWebhookHandler 创建 webhook 处理器实例
func NewWebhookHandler() *WebhookHandler {
	return &WebhookHandler{
		webhookService: service.NewWebhookService(),
	}
}

// Create 创建 webhook
// @Summary      创建 webhook
// @Description  创建一个新的 webhook 配置，返回唯一的 webhook URL
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        request  body      model.WebhookCreateRequest  true  "创建请求参数"
// @Success      200      {object}  model.Response{data=model.WebhookResponse}  "创建成功"
// @Failure      400      {object}  model.Response  "参数错误"
// @Failure      500      {object}  model.Response  "服务器错误"
// @Router       /webhook/create [post]
func (h *WebhookHandler) Create(c *gin.Context) {
	var req model.WebhookCreateRequest

	// 1. 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 获取用户 ID
	userID := getUserID(c)

	// 3. 调用服务层创建
	resp, err := h.webhookService.Create(&req, userID)
	if err != nil {
		model.ServerError(c, "创建失败: "+err.Error())
		return
	}

	// 4. 返回成功响应
	model.Success(c, resp)
}

// Update 更新 webhook
// @Summary      更新 webhook
// @Description  根据 ID 更新 webhook 配置信息
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        request  body      model.WebhookUpdateRequest  true  "更新请求参数"
// @Success      200      {object}  model.Response{data=model.WebhookResponse}  "更新成功"
// @Failure      400      {object}  model.Response  "参数错误"
// @Failure      404      {object}  model.Response  "记录不存在"
// @Failure      500      {object}  model.Response  "服务器错误"
// @Router       /webhook/update [post]
func (h *WebhookHandler) Update(c *gin.Context) {
	var req model.WebhookUpdateRequest

	// 1. 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 获取用户 ID
	userID := getUserID(c)

	// 3. 调用服务层更新
	resp, err := h.webhookService.Update(&req, userID)
	if err != nil {
		model.ParamError(c, err.Error())
		return
	}

	// 4. 返回成功响应
	model.Success(c, resp)
}

// Delete 删除 webhook
// @Summary      删除 webhook
// @Description  根据 ID 删除 webhook 配置
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Webhook ID"
// @Success      200  {object}  model.Response  "删除成功"
// @Failure      400  {object}  model.Response  "参数错误"
// @Failure      404  {object}  model.Response  "记录不存在"
// @Failure      500  {object}  model.Response  "服务器错误"
// @Router       /webhook/delete/:id [delete]
func (h *WebhookHandler) Delete(c *gin.Context) {
	// 1. 获取 ID 参数
	id, err := getInt64Param(c, "id")
	if err != nil {
		model.ParamError(c, "无效的 ID")
		return
	}

	// 2. 获取用户 ID
	userID := getUserID(c)

	// 3. 调用服务层删除
	if err := h.webhookService.Delete(id, userID); err != nil {
		model.ParamError(c, err.Error())
		return
	}

	// 4. 返回成功响应
	model.SuccessWithMessage(c, "删除成功", nil)
}

// GetByID 根据 ID 获取 webhook 详情
// @Summary      获取 webhook 详情
// @Description  根据 ID 获取 webhook 的详细信息
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        id   path      int64  true  "Webhook ID"
// @Success      200  {object}  model.Response{data=model.WebhookResponse}  "查询成功"
// @Failure      400  {object}  model.Response  "参数错误"
// @Failure      404  {object}  model.Response  "记录不存在"
// @Failure      500  {object}  model.Response  "服务器错误"
// @Router       /webhook/detail/:id [get]
func (h *WebhookHandler) GetByID(c *gin.Context) {
	// 1. 获取 ID 参数
	id, err := getInt64Param(c, "id")
	if err != nil {
		model.ParamError(c, "无效的 ID")
		return
	}

	// 2. 获取用户 ID
	userID := getUserID(c)

	// 3. 调用服务层查询
	resp, err := h.webhookService.GetByID(id, userID)
	if err != nil {
		model.ParamError(c, err.Error())
		return
	}

	// 4. 返回成功响应
	model.Success(c, resp)
}

// GetList 获取 webhook 列表
// @Summary      获取 webhook 列表
// @Description  获取当前用户的所有 webhook 配置列表
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.Response{data=[]model.WebhookResponse}  "查询成功"
// @Failure      500  {object}  model.Response  "服务器错误"
// @Router       /webhook/list [get]
func (h *WebhookHandler) GetList(c *gin.Context) {
	// 1. 获取用户 ID
	userID := getUserID(c)

	// 2. 调用服务层查询列表
	resp, err := h.webhookService.GetList(userID)
	if err != nil {
		model.ServerError(c, "获取列表失败: "+err.Error())
		return
	}

	// 3. 返回成功响应
	model.Success(c, resp)
}

// Callback 接收 webhook 回调
// @Summary      接收 webhook 回调
// @Description  外部服务调用此接口发送 webhook 事件，使用 X-Webhook-Signature 头进行签名验证
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        secret    path      string  true  "Webhook Secret"
// @Param        X-Webhook-Signature  header   string  true  "HMAC-SHA256 签名"
// @Success      200      {object}  model.Response  "接收成功"
// @Failure      401      {object}  model.Response  "签名验证失败"
// @Failure      404      {object}  model.Response  "webhook 不存在或已禁用"
// @Router       /webhook/callback/:secret [post]
func (h *WebhookHandler) Callback(c *gin.Context) {
	// 1. 获取 secret 参数
	secret := c.Param("secret")
	if secret == "" {
		model.ParamError(c, "缺少 secret 参数")
		return
	}

	// 2. 获取签名
	signature := c.GetHeader("X-Webhook-Signature")
	if signature == "" {
		model.ParamError(c, "缺少 X-Webhook-Signature 签名头")
		return
	}

	// 3. 获取请求体
	body, err := c.GetRawData()
	if err != nil {
		model.ParamError(c, "读取请求体失败")
		return
	}

	// 4. 验证签名
	if !h.webhookService.VerifyCallback(secret, signature, body) {
		model.Unauthorized(c, "签名验证失败")
		return
	}

	var data map[string]interface{}
	json.Unmarshal(body, &data)

	logger.Debug(string(body))

	switch data["event"] {
	case "cjl-pay-success":
		logger.Infof("收到支付成功事件，订单ID: %v", data["data"].(map[string]interface{})["order_id"])
	default:
		logger.Infof("收到未知事件: %v", data["event"])
	}

	// 5. 返回成功响应
	model.SuccessWithMessage(c, "接收成功", nil)
}

// Sign 生成签名
// @Summary      生成 webhook 签名
// @Description  根据 webhook ID 和请求体生成 HMAC-SHA256 签名，用于测试回调
// @Tags         Webhook 模块
// @Accept       json
// @Produce      json
// @Param        request  body      model.WebhookSignRequest  true  "签名请求参数"
// @Success      200      {object}  model.Response{data=model.WebhookSignResponse}  "签名生成成功"
// @Failure      400      {object}  model.Response  "参数错误"
// @Failure      404      {object}  model.Response  "webhook 不存在"
// @Router       /webhook/sign [post]
func (h *WebhookHandler) Sign(c *gin.Context) {
	var req model.WebhookSignRequest

	// 1. 绑定请求参数
	if err := c.ShouldBindJSON(&req); err != nil {
		model.ParamError(c, "参数错误: "+err.Error())
		return
	}

	// 2. 获取用户 ID
	userID := getUserID(c)

	// 3. 调用服务层生成签名
	resp, err := h.webhookService.Sign(&req, userID)
	if err != nil {
		model.ParamError(c, err.Error())
		return
	}

	// 4. 返回成功响应
	model.Success(c, resp)
}

// getUserID 从 gin.Context 中获取用户 ID
func getUserID(c *gin.Context) int64 {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	// userID 是字符串类型，转换为 int64
	var id int64
	fmt.Sscanf(userID.(string), "%d", &id)
	return id
}

// getInt64Param 从路径参数中获取 int64
func getInt64Param(c *gin.Context, name string) (int64, error) {
	val := c.Param(name)
	if val == "" {
		return 0, nil
	}
	var id int64
	_, err := fmt.Sscanf(val, "%d", &id)
	return id, err
}
