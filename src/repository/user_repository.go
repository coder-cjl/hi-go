package repository

import (
	"hi-go/src/model"
	"hi-go/src/utils/mysql"
)

// UserRepository 用户数据访问层
type UserRepository struct{}

// NewUserRepository 创建用户仓储实例
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := mysql.Database.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据ID查找用户
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := mysql.Database.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	return mysql.Database.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return mysql.Database.Save(user).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := mysql.Database.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := mysql.Database.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
