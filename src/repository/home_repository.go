package repository

import (
	"hi-go/src/model"
	"hi-go/src/utils/mysql"
)

// HomeRepository 首页数据访问层
type HomeRepository struct{}

// NewHomeRepository 创建首页仓储实例
func NewHomeRepository() *HomeRepository {
	return &HomeRepository{}
}

// List 获取首页列表（分页）
func (r *HomeRepository) List(page, pageSize int) ([]model.Home, int64, error) {
	var homes []model.Home
	var total int64

	// 构建查询
	query := mysql.Database.Model(&model.Home{}).
		Where("status = ?", 1).    // 只查询启用的
		Order("sort ASC, id DESC") // 按排序和ID倒序

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&homes).Error; err != nil {
		return nil, 0, err
	}

	return homes, total, nil
}

// 根据ID查找首页内容
func (r *HomeRepository) FindByID(id int64) (*model.Home, error) {
	var home model.Home
	err := mysql.Database.First(&home, id).Error
	if err != nil {
		return nil, err
	}
	return &home, nil
}

// 创建首页内容
func (r *HomeRepository) Create(home *model.Home) error {
	return mysql.Database.Create(home).Error
}

// 批量创建首页内容
func (r *HomeRepository) BatchCreate(homes []model.Home) error {
	return mysql.Database.Create(&homes).Error
}

// Update 更新首页内容
func (r *HomeRepository) Update(id int64, updates map[string]interface{}) error {
	return mysql.Database.Model(&model.Home{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除首页内容
func (r *HomeRepository) Delete(id int64) error {
	return mysql.Database.Delete(&model.Home{}, id).Error
}
