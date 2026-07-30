package event

import "github.com/YasinDoyle/e-mall/repository/db/model"

type SellerApplied struct {
	SellerID uint
	ShopName string
}

type SellerAuditChanged struct {
	SellerID     uint
	Status       uint
	RejectReason string
}

type ProductSubmitted struct {
	Product *model.Product
}

type ProductChanged struct {
	Product *model.Product
}

type ProductDeleted struct {
	ProductID uint
}

type ProductAuditChanged struct {
	ProductID   uint
	SellerID    uint
	ProductName string
	AuditStatus uint
}

type OrderPaid struct {
	OrderID  uint
	OrderNum uint64
	BuyerID  uint
	SellerID uint
}

type OrderShipped struct {
	OrderID    uint
	OrderNum   uint64
	BuyerID    uint
	TrackingNo string
}

type RefundRequested struct {
	OrderID  uint
	OrderNum uint64
	SellerID uint
}

type OrderRefunded struct {
	OrderID      uint
	OrderNum     uint64
	BuyerID      uint
	SellerID     uint
	RefundAmount float64
}

type OrderReceived struct {
	OrderID  uint
	OrderNum uint64
	SellerID uint
}

type SettlementGenerated struct {
	SellerID     uint
	Count        int64
	SettlementID uint
	OrderID      uint
	Amount       float64
}

type SettlementPaid struct {
	Settlement *model.Settlement
}

type WithdrawApplied struct {
	Withdraw *model.SellerWithdraw
	ShopName string
}

type WithdrawAuditChanged struct {
	WithdrawID uint
	SellerID   uint
	Amount     float64
	Status     string
	Reason     string
}

type WithdrawPaidStatusChanged struct {
	WithdrawID uint
	SellerID   uint
	Amount     float64
	Status     string
	Reason     string
}
