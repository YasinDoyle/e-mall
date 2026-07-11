package model

import "gorm.io/gorm"

// Review 商品评价
type Review struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	ProductID uint   `gorm:"not null;index"`
	OrderID   uint   `gorm:"not null"`
	Rating    uint   `gorm:"not null;default:5"` // 1-5 星
	Content   string `gorm:"size:1000"`
	Images    string `gorm:"size:2000"` // 多张图片 URL，逗号分隔
	UserName  string `gorm:"size:100"`
	UserAvatar string `gorm:"size:500"`
}
