package types

type AdminSettlementListReq struct {
	PageNum  int    `json:"page_num" form:"page_num"`
	PageSize int    `json:"page_size" form:"page_size"`
	SellerID uint   `json:"seller_id" form:"seller_id"`
	Status   string `json:"status" form:"status"`
}

type AdminSettlementGenerateReq struct {
	SellerID uint `json:"seller_id" binding:"required"`
}

type AdminSettlementGenerateOneReq struct {
	ID uint `json:"id" binding:"required"`
}

type AdminSettlementPayReq struct {
	ID uint `json:"id" binding:"required"`
}

type AdminSettlementDetailReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type SellerSettlementSummaryResp struct {
	AvailableAmount float64 `json:"available_amount"`
	PendingAmount   float64 `json:"pending_amount"`
	GeneratedAmount float64 `json:"generated_amount"`
	PaidAmount      float64 `json:"paid_amount"`
	RefundedAmount  float64 `json:"refunded_amount"`
}

type AdminSettlementResp struct {
	ID                uint    `json:"id"`
	SellerID          uint    `json:"seller_id"`
	OrderID           uint    `json:"order_id"`
	OrderNum          uint64  `json:"order_num"`
	GrossAmount       float64 `json:"gross_amount"`
	CommissionRate    float64 `json:"commission_rate"`
	CommissionAmount  float64 `json:"commission_amount"`
	SettlementAmount  float64 `json:"settlement_amount"`
	Status            string  `json:"status"`
	OrderType         uint    `json:"order_type"`
	OrderRefundStatus int     `json:"order_refund_status"`
	CanMarkPaid       bool    `json:"can_mark_paid"`
	PaidAt            int64   `json:"paid_at"`
	CreatedAt         int64   `json:"created_at"`
}

type AdminSettlementGenerateResp struct {
	SellerID uint  `json:"seller_id"`
	Count    int64 `json:"count"`
}
