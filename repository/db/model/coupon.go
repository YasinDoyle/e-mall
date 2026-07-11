package model

import (
	"time"

	"gorm.io/gorm"
)

// CouponType 优惠券类型
const (
	CouponTypeFixed   = 1 // 固定减免金额
	CouponTypePercent = 2 // 折扣（百分比）
)

// Coupon 优惠券模板（管理员创建）
type Coupon struct {
	gorm.Model
	Name       string    `gorm:"size:100;not null"`
	CouponType uint      `gorm:"default:1"`     // 1:固定减免 2:折扣
	Discount   float64   `gorm:"not null"`      // 固定减免金额 或 折扣系数（如 0.9 表示九折）
	MinAmount  float64   `gorm:"default:0"`     // 使用门槛（订单满多少才可用）
	Stock      int       `gorm:"not null"`      // 发行库存，-1 表示无限
	ExpireAt   time.Time `gorm:"not null"`      // 过期时间
}
