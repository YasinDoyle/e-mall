package model

import "gorm.io/gorm"

const (
	AccountFlowTypeBuyerPay               = "buyer_pay"
	AccountFlowTypeSellerPending          = "seller_pending"
	AccountFlowTypeSellerSettlementCredit = "seller_settlement_credit"
	AccountFlowTypePlatformCommission     = "platform_commission"
	AccountFlowTypeSettlementPaid         = "settlement_paid"
	AccountFlowTypeRefund                 = "refund"
	AccountFlowTypeSellerWithdrawFreeze   = "seller_withdraw_freeze"
	AccountFlowTypeSellerWithdrawPaid     = "seller_withdraw_paid"
	AccountFlowTypeSellerWithdrawUnfreeze = "seller_withdraw_unfreeze"
	AccountFlowTypeSellerManualAdjustment = "seller_manual_adjustment"
)

type AccountFlow struct {
	gorm.Model
	FlowNo      string  `gorm:"size:64;uniqueIndex;not null" json:"flow_no"`
	OrderID     uint    `gorm:"index" json:"order_id"`
	OrderNum    uint64  `gorm:"index" json:"order_num"`
	UserID      uint    `gorm:"index" json:"user_id"`
	SellerID    uint    `gorm:"index" json:"seller_id"`
	RelatedType string  `gorm:"size:32;index" json:"related_type"`
	RelatedID   uint    `gorm:"index" json:"related_id"`
	FlowType    string  `gorm:"size:32;not null;index" json:"flow_type"`
	Direction   string  `gorm:"size:16;not null" json:"direction"`
	Amount      float64 `gorm:"not null" json:"amount"`
	Remark      string  `gorm:"size:255" json:"remark"`
}
