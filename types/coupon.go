package types

import "time"

// ===== 管理员创建优惠券 =====

type AdminCouponCreateReq struct {
	Name       string    `json:"name" binding:"required"`
	CouponType uint      `json:"coupon_type" binding:"required,oneof=1 2"`
	Discount   float64   `json:"discount" binding:"required"`
	MinAmount  float64   `json:"min_amount"`
	Stock      int       `json:"stock"` // -1=无限
	ExpireAt   time.Time `json:"expire_at" binding:"required"`
}

type AdminCouponOfflineReq struct {
	ID uint `json:"id" binding:"required"`
}

type CouponResp struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	CouponType uint      `json:"coupon_type"`
	Discount   float64   `json:"discount"`
	MinAmount  float64   `json:"min_amount"`
	Stock      int       `json:"stock"`
	ExpireAt   time.Time `json:"expire_at"`
	CreatedAt  int64     `json:"created_at"`
}

// ===== 用户领券 =====

type CouponClaimReq struct {
	CouponID uint `json:"coupon_id" binding:"required"`
}

// ===== 用户查看自己的优惠券 =====

type UserCouponResp struct {
	ID         uint      `json:"id"` // UserCoupon.ID
	CouponID   uint      `json:"coupon_id"`
	Name       string    `json:"name"`
	CouponType uint      `json:"coupon_type"`
	Discount   float64   `json:"discount"`
	MinAmount  float64   `json:"min_amount"`
	ExpireAt   time.Time `json:"expire_at"`
	Status     uint      `json:"status"` // 0:未使用 1:已使用
}
