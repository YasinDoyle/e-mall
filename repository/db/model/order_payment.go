package model

import (
	"time"

	"gorm.io/gorm"
)

type OrderPayment struct {
	gorm.Model
	OrderID         uint    `gorm:"not null;index"`
	OrderNum        uint64  `gorm:"not null;index"`
	PaymentNo       string  `gorm:"size:64;not null;uniqueIndex"`
	UserID          uint    `gorm:"not null;index"`
	Channel         string  `gorm:"size:20;not null;index"`
	Amount          float64 `gorm:"not null"`
	Status          string  `gorm:"size:20;not null;default:'pending';index"`
	ProviderTradeNo string  `gorm:"size:128;index"`
	PaidAt          *time.Time
	ClosedAt        *time.Time
}
