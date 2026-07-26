package model

import (
	"time"

	"gorm.io/gorm"
)

// Order 订单信息
type Order struct {
	gorm.Model
	UserID       uint   `gorm:"not null"`
	ProductID    uint   `gorm:"not null"`
	BossID       uint   `gorm:"not null"`
	AddressID    uint   `gorm:"not null"`
	Num          int    // 数量
	OrderNum     uint64 // 订单号
	Type         uint   // 1 未支付  2 已支付
	Money        float64
	PaidAt       *time.Time
	RefundStatus int    `gorm:"default:0;index"` // 0 无退款 1 申请中 2 已退款
	RefundReason string `gorm:"size:255"`
	TrackingNo   string `gorm:"size:64"`
	BuyerDeleted bool   `gorm:"default:false;index"`
}
