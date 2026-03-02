package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hi-go/src/model"
	"hi-go/src/repository"
	"hi-go/src/utils/logger"
	"hi-go/src/utils/snowflake"
)

// WebhookService Webhook 业务逻辑层
type WebhookService struct {
	webhookRepo *repository.WebhookRepository
}

// 创建 WebhookService 实例
func NewWebhookService() *WebhookService {
	return &WebhookService{
		webhookRepo: repository.NewWebhookRepository(),
	}
}

// Create 创建 webhook（返回包含 secret 的响应）
func (s *WebhookService) Create(req *model.WebhookCreateRequest, userID int64) (*model.WebhookResponseWithSecret, error) {
	// 生成唯一 secret
	secret, err := generateSecret()
	if err != nil {
		return nil, fmt.Errorf("生成密钥失败: %v", err)
	}

	// 构建 webhook 对象
	webhook := &model.Webhook{
		ID:          snowflake.MustGenerate(),
		Name:        req.Name,
		CallbackURL: req.CallbackURL,
		Event:       req.Event,
		Secret:      secret,
		Enabled:     1,
		UserID:      userID,
	}

	// 如果传入了 enabled 参数
	if req.Enabled != nil {
		if *req.Enabled {
			webhook.Enabled = 1
		} else {
			webhook.Enabled = 0
		}
	}

	// 保存到数据库
	if err := s.webhookRepo.Create(webhook); err != nil {
		return nil, fmt.Errorf("创建 webhook 失败: %v", err)
	}

	return webhook.ToResponseWithSecret(), nil
}

// Update 更新 webhook
func (s *WebhookService) Update(req *model.WebhookUpdateRequest, userID int64) (*model.WebhookResponse, error) {
	// 1. 检查 webhook 是否存在
	webhook, err := s.webhookRepo.FindByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("webhook 不存在")
	}

	// 2. 检查权限（只能修改自己的 webhook）
	if webhook.UserID != userID {
		return nil, fmt.Errorf("无权限修改此 webhook")
	}

	// 3. 构建更新数据
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.CallbackURL != "" {
		updates["callback_url"] = req.CallbackURL
	}
	if req.Event != "" {
		updates["event"] = req.Event
	}
	if req.Enabled != nil {
		if *req.Enabled {
			updates["enabled"] = 1
		} else {
			updates["enabled"] = 0
		}
	}

	// 4. 执行更新
	if len(updates) == 0 {
		return nil, fmt.Errorf("没有需要更新的字段")
	}

	// 直接更新
	webhook.Name = req.Name
	webhook.CallbackURL = req.CallbackURL
	webhook.Event = req.Event
	if req.Enabled != nil {
		if *req.Enabled {
			webhook.Enabled = 1
		} else {
			webhook.Enabled = 0
		}
	}

	if err := s.webhookRepo.Update(webhook); err != nil {
		return nil, fmt.Errorf("更新 webhook 失败: %v", err)
	}

	return webhook.ToResponse(), nil
}

// Delete 删除 webhook
func (s *WebhookService) Delete(id int64, userID int64) error {
	// 1. 检查 webhook 是否存在
	webhook, err := s.webhookRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("webhook 不存在")
	}

	// 2. 检查权限（只能删除自己的 webhook）
	if webhook.UserID != userID {
		return fmt.Errorf("无权限删除此 webhook")
	}

	// 3. 执行删除
	return s.webhookRepo.Delete(id)
}

// GetByID 根据 ID 获取 webhook
func (s *WebhookService) GetByID(id int64, userID int64) (*model.WebhookResponse, error) {
	webhook, err := s.webhookRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("webhook 不存在")
	}

	// 检查权限
	if webhook.UserID != userID {
		return nil, fmt.Errorf("无权限访问此 webhook")
	}

	return webhook.ToResponse(), nil
}

// GetList 获取用户的 webhook 列表
func (s *WebhookService) GetList(userID int64) ([]*model.WebhookResponse, error) {
	webhooks, err := s.webhookRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("获取 webhook 列表失败: %v", err)
	}

	// 转换为响应结构
	result := make([]*model.WebhookResponse, 0, len(webhooks))
	for _, w := range webhooks {
		result = append(result, w.ToResponse())
	}

	return result, nil
}

// VerifyCallback 验证 webhook 回调签名
func (s *WebhookService) VerifyCallback(secret string, signature string, body []byte) bool {
	// 根据 secret 查找 webhook
	webhook, err := s.webhookRepo.FindBySecret(secret)
	if err != nil {
		return false
	}

	// 如果 webhook 未启用，也返回 false
	if webhook.Enabled != 1 {
		return false
	}

	logger.Debug(string(body))

	// 计算签名
	expectedSignature := generateSignature(body, webhook.Secret)

	// 使用恒定时间比较防止时序攻击
	return hmac.Equal([]byte(expectedSignature), []byte(signature))
}

// generateSecret 生成随机密钥
func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateSignature 生成 HMAC-SHA256 签名
func generateSignature(body []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

// Sign 生成签名（供 API 调用）
func (s *WebhookService) Sign(req *model.WebhookSignRequest, userID int64) (*model.WebhookSignResponse, error) {
	// 1. 检查 webhook 是否存在
	webhook, err := s.webhookRepo.FindByID(req.ID)
	if err != nil {
		return nil, fmt.Errorf("webhook 不存在")
	}

	// 2. 检查权限
	if webhook.UserID != userID {
		return nil, fmt.Errorf("无权限访问此 webhook")
	}

	logger.Debug(req.Body)

	// 3. 生成签名
	signature := generateSignature([]byte(req.Body), webhook.Secret)

	return &model.WebhookSignResponse{
		Signature: signature,
		Method:    "HMAC-SHA256",
		Header:    "X-Webhook-Signature",
		Secret:    webhook.Secret,
	}, nil
}
