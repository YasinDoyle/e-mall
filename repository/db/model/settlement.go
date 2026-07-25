package model

import "gorm.io/gorm"

const (
	SettlementStatusPending   = "pending"
	SettlementStatusGenerated = "generated"
	SettlementStatusPaid      = "paid"
	SettlementStatusRefunded  = "refunded"
)

type Settlement struct {
	gorm.Model
	SellerID         uint    `gorm:"not null;index" json:"seller_id"`
	OrderID          uint    `gorm:"not null;uniqueIndex" json:"order_id"`
	OrderNum         uint64  `gorm:"not null;index" json:"order_num"`
	GrossAmount      float64 `gorm:"not null" json:"gross_amount"`
	CommissionRate   float64 `gorm:"not null" json:"commission_rate"`
	CommissionAmount float64 `gorm:"not null" json:"commission_amount"`
	SettlementAmount float64 `gorm:"not null" json:"settlement_amount"`
	Status           string  `gorm:"size:32;not null;default:'pending';index" json:"status"`
	PaidAt           *int64  `json:"paid_at"`
}
