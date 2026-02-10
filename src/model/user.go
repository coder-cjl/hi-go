package model

import (
	"time"

	"gorm.io/gorm"
)

// 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 表示不返回密码
	Email     string         `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string         `gorm:"type:varchar(255)" json:"avatar"`
	Status    int            `gorm:"type:tinyint;default:1;comment:'1-正常 0-禁用'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 指定表名
func (User) TableName() string {
	return "users"
}

// 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"` // 必填，长度3-50
	Password string `json:"password" binding:"required,min=6"`        // 必填，最小6位
}

// 登录响应
type LoginResponse struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 过期时间（秒）
}

// 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Nickname string `json:"nickname" binding:"omitempty,max=50"`
}
