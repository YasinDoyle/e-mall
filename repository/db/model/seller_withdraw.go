package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	SellerWithdrawStatusPending  = "pending"
	SellerWithdrawStatusApproved = "approved"
	SellerWithdrawStatusRejected = "rejected"
	SellerWithdrawStatusPaid     = "paid"
	SellerWithdrawStatusFailed   = "failed"
)

type SellerWithdraw struct {
	gorm.Model
	SellerID          uint          `gorm:"not null;index" json:"seller_id"`
	Seller            SellerProfile `gorm:"foreignKey:SellerID;references:UserID" json:"seller"`
	Amount            float64       `gorm:"not null" json:"amount"`
	Status            string        `gorm:"size:32;not null;default:'pending';index" json:"status"`
	PayeeName         string        `gorm:"size:64;not null" json:"payee_name"`
	PayeeAccount      string        `gorm:"size:128;not null" json:"payee_account"`
	PayeeChannel      string        `gorm:"size:32;not null;default:'manual'" json:"payee_channel"`
	AuditReason       string        `gorm:"size:255" json:"audit_reason"`
	AuditOperatorID   uint          `gorm:"index" json:"audit_operator_id"`
	AuditOperatorName string        `gorm:"size:64" json:"audit_operator_name"`
	PaidOperatorID    uint          `gorm:"index" json:"paid_operator_id"`
	PaidOperatorName  string        `gorm:"size:64" json:"paid_operator_name"`
	AuditedAt         *time.Time    `json:"audited_at"`
	PaidAt            *time.Time    `json:"paid_at"`
}

func (w *SellerWithdraw) IsPending() bool {
	return w != nil && w.Status == SellerWithdrawStatusPending
}
