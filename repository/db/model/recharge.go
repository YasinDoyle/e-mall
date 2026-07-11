package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	RechargeChannelWechat = "wechat"
	RechargeChannelAlipay = "alipay"

	RechargeStatusPending  = "pending"
	RechargeStatusPaid     = "paid"
	RechargeStatusCredited = "credited"
	RechargeStatusFailed   = "failed"

	RechargeRefundNone       = "none"
	RechargeRefundProcessing = "processing"
	RechargeRefundSuccess    = "success"
	RechargeRefundFailed     = "failed"
)

type RechargeOrder struct {
	gorm.Model
	OrderNum        string `gorm:"size:64;not null;uniqueIndex"`
	UserID          uint   `gorm:"not null;index"`
	Channel         string `gorm:"size:20;not null;index"`
	Amount          float64
	Status          string `gorm:"size:20;not null;default:'pending';index"`
	ProviderTradeNo string `gorm:"size:128;index"`
	PaidAt          *time.Time
	RefundNo        string `gorm:"size:64;index"`
	RefundAmount    float64
	RefundStatus    string `gorm:"size:20;not null;default:'none';index"`
	RefundReason    string `gorm:"size:255"`
	RefundedAt      *time.Time
}
