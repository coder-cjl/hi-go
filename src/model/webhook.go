package model

import (
	"time"

	"gorm.io/gorm"
)

// Webhook Webhook 模型
type Webhook struct {
	ID          int64          `gorm:"primarykey;autoIncrement:false" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	CallbackURL string         `gorm:"type:varchar(500);not null" json:"callback_url"`
	Event       string         `gorm:"type:varchar(100);not null" json:"event"`
	Secret      string         `gorm:"type:varchar(255);not null" json:"-"`
	Enabled     int            `gorm:"type:tinyint;default:1;comment:'1-启用 0-禁用'" json:"enabled"`
	UserID      int64          `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// 指定表名
func (Webhook) TableName() string {
	return "webhooks"
}

// WebhookCreateRequest 创建 webhook 请求
type WebhookCreateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	CallbackURL string `json:"callback_url" binding:"required,url"`
	Event       string `json:"event" binding:"required,min=1,max=100"`
	Enabled     *bool  `json:"enabled"`
}

// WebhookUpdateRequest 更新 webhook 请求
type WebhookUpdateRequest struct {
	ID          int64  `json:"id" binding:"required"`
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	CallbackURL string `json:"callback_url" binding:"omitempty,url"`
	Event       string `json:"event" binding:"omitempty,min=1,max=100"`
	Enabled     *bool  `json:"enabled"`
}

// WebhookResponse webhook 响应（不包含敏感信息）
type WebhookResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
	Event       string `json:"event"`
	Enabled     int    `json:"enabled"`
	UserID      int64  `json:"user_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ToResponse 转换为响应结构
func (w *Webhook) ToResponse() *WebhookResponse {
	return &WebhookResponse{
		ID:          w.ID,
		Name:        w.Name,
		CallbackURL: w.CallbackURL,
		Event:       w.Event,
		Enabled:     w.Enabled,
		UserID:      w.UserID,
		CreatedAt:   w.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   w.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// WebhookResponseWithSecret 包含 secret 的响应（仅创建时返回一次）
type WebhookResponseWithSecret struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	CallbackURL string `json:"callback_url"`
	Event       string `json:"event"`
	Enabled     int    `json:"enabled"`
	UserID      int64  `json:"user_id"`
	Secret      string `json:"secret"` // 仅创建时返回，后续不再显示
	CallbackURLFull string `json:"callback_url_full"` // 完整的回调地址
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ToResponseWithSecret 转换为包含 secret 的响应结构
func (w *Webhook) ToResponseWithSecret() *WebhookResponseWithSecret {
	return &WebhookResponseWithSecret{
		ID:                w.ID,
		Name:              w.Name,
		CallbackURL:       w.CallbackURL,
		Event:             w.Event,
		Enabled:           w.Enabled,
		UserID:            w.UserID,
		Secret:            w.Secret,
		CallbackURLFull:   "/api/webhook/callback/" + w.Secret,
		CreatedAt:         w.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         w.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

// WebhookSignRequest 生成签名请求
type WebhookSignRequest struct {
	ID   int64  `json:"id" binding:"required"`
	Body string `json:"body" binding:"required"` // 需要签名的请求体
}

// WebhookSignResponse 签名响应
type WebhookSignResponse struct {
	Signature string `json:"signature"`
	Method    string `json:"method"`
	Header    string `json:"header"`
	Secret    string `json:"secret"` // 返回 secret 用于调试
}
