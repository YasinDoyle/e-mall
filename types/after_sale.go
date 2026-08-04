package types

type AfterSaleRequestReq struct {
	OrderId uint   `json:"order_id" form:"order_id" binding:"required"`
	Type    string `json:"type" form:"type" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

type SellerAfterSaleHandleReq struct {
	AfterSaleID uint   `json:"after_sale_id" form:"after_sale_id" binding:"required"`
	Action      string `json:"action" form:"action" binding:"required"`
	Reason      string `json:"reason" form:"reason"`
}

type AdminAfterSaleHandleReq struct {
	AfterSaleID uint   `json:"after_sale_id" form:"after_sale_id" binding:"required"`
	Action      string `json:"action" form:"action" binding:"required"`
	Note        string `json:"note" form:"note"`
}

type AfterSaleListReq struct {
	OrderID uint   `form:"order_id" json:"order_id"`
	Status  string `form:"status" json:"status"`
	Type    string `form:"type" json:"type"`
	BasePage
}

type AfterSaleResp struct {
	ID           uint    `json:"id"`
	OrderID      uint    `json:"order_id"`
	OrderNum     uint64  `json:"order_num"`
	BuyerID      uint    `json:"buyer_id"`
	SellerID     uint    `json:"seller_id"`
	Type         string  `json:"type"`
	Status       string  `json:"status"`
	Reason       string  `json:"reason"`
	RefundAmount float64 `json:"refund_amount"`
	SellerReason string  `json:"seller_reason"`
	PlatformNote string  `json:"platform_note"`
	CreatedAt    int64   `json:"created_at"`
	UpdatedAt    int64   `json:"updated_at"`
	RefundedAt   int64   `json:"refunded_at"`
	ClosedAt     int64   `json:"closed_at"`
}

type AfterSaleListResp = AfterSaleResp
