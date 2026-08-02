package types

import "time"

type PaymentDownReq struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	Money     float64 `form:"money" json:"money"`
	OrderNo   string  `form:"orderNo" json:"orderNo"`
	ProductID int     `form:"product_id" json:"product_id"`
	PayTime   string  `form:"payTime" json:"payTime" `
	Sign      string  `form:"sign" json:"sign" `
	BossID    int     `form:"boss_id" json:"boss_id"`
	BossName  string  `form:"boss_name" json:"boss_name"`
	Num       int     `form:"num" json:"num"`
	Key       string  `form:"key" json:"key"`
}

type OrderGatewayPayReq struct {
	OrderID uint `form:"order_id" json:"order_id" binding:"required"`
}

type OrderPaymentStatusReq struct {
	PaymentNo string `form:"payment_no" json:"payment_no" binding:"required"`
}

type OrderPaymentResp struct {
	OrderID     uint    `json:"order_id"`
	OrderNum    uint64  `json:"order_num"`
	PaymentNo   string  `json:"payment_no"`
	Channel     string  `json:"channel"`
	Amount      float64 `json:"amount"`
	Status      string  `json:"status"`
	QRCodeURL   string  `json:"qr_code_url,omitempty"`
	PayURL      string  `json:"pay_url,omitempty"`
	PaidAt      int64   `json:"paid_at,omitempty"`
	ClosedAt    int64   `json:"closed_at,omitempty"`
	OrderStatus uint    `json:"order_status,omitempty"`
}

type OrderPaymentCallbackReq struct {
	PaymentNo       string
	ProviderTradeNo string
	ExpectedChannel string
	PaidAt          time.Time
	PaidAmount      float64
}

type OrderPaidEvent struct {
	OrderID     uint      `json:"order_id"`
	OrderNum    uint64    `json:"order_num"`
	UserID      uint      `json:"user_id"`
	BossID      uint      `json:"boss_id"`
	ProductID   uint      `json:"product_id"`
	Num         int       `json:"num"`
	TotalAmount float64   `json:"total_amount"`
	PaidAt      time.Time `json:"paid_at"`
}

type RechargePaidEvent struct {
	OrderNum        string    `json:"order_num"`
	UserID          uint      `json:"user_id"`
	Channel         string    `json:"channel"`
	Amount          float64   `json:"amount"`
	ProviderTradeNo string    `json:"provider_trade_no"`
	PaidAt          time.Time `json:"paid_at"`
}
