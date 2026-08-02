package model

import "gorm.io/gorm"

type AfterSale struct {
	gorm.Model
	OrderID      uint    `gorm:"not null;index" json:"order_id"`
	OrderNum     uint64  `gorm:"not null;index" json:"order_num"`
	BuyerID      uint    `gorm:"not null;index" json:"buyer_id"`
	SellerID     uint    `gorm:"not null;index" json:"seller_id"`
	Type         string  `gorm:"size:32;not null;index" json:"type"`
	Status       string  `gorm:"size:32;not null;index" json:"status"`
	Reason       string  `gorm:"size:255;not null" json:"reason"`
	RefundAmount float64 `gorm:"not null" json:"refund_amount"`
	SellerReason string  `gorm:"size:255" json:"seller_reason"`
	PlatformNote string  `gorm:"size:255" json:"platform_note"`
	RefundedAt   *int64  `json:"refunded_at"`
	ClosedAt     *int64  `json:"closed_at"`
}
