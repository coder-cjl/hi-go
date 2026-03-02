package repository

import (
	"hi-go/src/model"
	"hi-go/src/utils/mysql"
)

// WebhookRepository Webhook 数据访问层
type WebhookRepository struct{}

// 创建 WebhookRepository 实例
func NewWebhookRepository() *WebhookRepository {
	return &WebhookRepository{}
}

// Create 创建 webhook
func (r *WebhookRepository) Create(webhook *model.Webhook) error {
	return mysql.Database.Create(webhook).Error
}

// Update 更新 webhook
func (r *WebhookRepository) Update(webhook *model.Webhook) error {
	return mysql.Database.Save(webhook).Error
}

// Delete 删除 webhook（软删除）
func (r *WebhookRepository) Delete(id int64) error {
	return mysql.Database.Delete(&model.Webhook{}, "id = ?", id).Error
}

// FindByID 根据 ID 查找 webhook
func (r *WebhookRepository) FindByID(id int64) (*model.Webhook, error) {
	var webhook model.Webhook
	err := mysql.Database.Where("id = ?", id).First(&webhook).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// FindBySecret 根据 secret 查找 webhook（用于回调验证）
func (r *WebhookRepository) FindBySecret(secret string) (*model.Webhook, error) {
	var webhook model.Webhook
	err := mysql.Database.Where("secret = ? AND enabled = ?", secret, 1).First(&webhook).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// FindByUserID 根据用户 ID 查找 webhook 列表
func (r *WebhookRepository) FindByUserID(userID int64) ([]*model.Webhook, error) {
	var webhooks []*model.Webhook
	err := mysql.Database.Where("user_id = ?", userID).Find(&webhooks).Error
	return webhooks, err
}

// ExistsByID 检查 webhook 是否存在
func (r *WebhookRepository) ExistsByID(id int64) (bool, error) {
	var count int64
	err := mysql.Database.Model(&model.Webhook{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
