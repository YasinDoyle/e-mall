package model

import "gorm.io/gorm"

// CouponUsed 用户优惠券使用状态
const (
	UserCouponUnused = 0 // 未使用
	UserCouponUsed   = 1 // 已使用
)

// UserCoupon 用户领券记录
type UserCoupon struct {
	gorm.Model
	UserID   uint `gorm:"not null;uniqueIndex:idx_user_coupon_once"`
	CouponID uint `gorm:"not null;uniqueIndex:idx_user_coupon_once"`
	Status   uint `gorm:"default:0"` // 0:未使用 1:已使用
	OrderID  uint `gorm:"default:0"` // 核销时记录订单 ID
}
