package model

import "gorm.io/gorm"

// OrderLog records immutable business operation history for order fulfillment
// status transitions, such as pay, ship, receive, cancel, and refund approval.
type OrderLog struct {
	gorm.Model
	OrderID      uint   `gorm:"not null;index" json:"order_id"`
	OrderNum     uint64 `gorm:"not null;index" json:"order_num"`
	Action       string `gorm:"size:32;not null;index" json:"action"`
	FromType     uint   `gorm:"not null" json:"from_type"`
	ToType       uint   `gorm:"not null" json:"to_type"`
	OperatorType string `gorm:"size:32;not null;index" json:"operator_type"`
	OperatorID   uint   `gorm:"not null;index" json:"operator_id"`
	Remark       string `gorm:"size:255" json:"remark"`
}
