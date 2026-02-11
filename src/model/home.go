package model

import "time"

// Home 首页内容模型
type Home struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"type:varchar(200);not null" json:"title"` // 标题
	Description string    `gorm:"type:varchar(500)" json:"description"`    // 描述
	ImageURL    string    `gorm:"type:varchar(500)" json:"image_url"`      // 图片URL
	Link        string    `gorm:"type:varchar(500)" json:"link"`           // 链接
	Sort        int       `gorm:"default:0" json:"sort"`                   // 排序（越小越靠前）
	Status      int       `gorm:"default:1" json:"status"`                 // 状态：1-启用 0-禁用
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`        // 创建时间
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`        // 更新时间
}

// TableName 指定表名
func (Home) TableName() string {
	return "home"
}

// HomeListRequest 首页列表请求
type HomeListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`      // 页码
	PageSize int `form:"page_size" binding:"omitempty,min=1"` // 每页数量
}

// HomeListData 首页列表data字段
type HomeListData struct {
	List  []Home `json:"list"`  // 列表数据
	Total int64  `json:"total"` // 总数
}

// HomeUpdateRequest 首页更新请求
type HomeUpdateRequest struct {
	Title       string `json:"title" binding:"omitempty,max=200"`       // 标题
	Description string `json:"description" binding:"omitempty,max=500"` // 描述
	ImageURL    string `json:"image_url" binding:"omitempty,max=500"`   // 图片URL
	Link        string `json:"link" binding:"omitempty,max=500"`        // 链接
	Sort        *int   `json:"sort" binding:"omitempty,min=0"`          // 排序
	Status      *int   `json:"status" binding:"omitempty,oneof=0 1"`    // 状态
}
