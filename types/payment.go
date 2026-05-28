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
