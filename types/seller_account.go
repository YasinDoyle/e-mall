package types

type SellerAccountSummaryResp struct {
	SellerID         uint    `json:"seller_id"`
	AvailableBalance float64 `json:"available_balance"`
	FrozenBalance    float64 `json:"frozen_balance"`
	TotalIncome      float64 `json:"total_income"`
	TotalWithdrawn   float64 `json:"total_withdrawn"`
}

type AdminSellerAccountBackfillResp struct {
	SellerCount     int64   `json:"seller_count"`
	SettlementCount int64   `json:"settlement_count"`
	Amount          float64 `json:"amount"`
}

type SellerWithdrawApplyReq struct {
	Amount       float64 `json:"amount" form:"amount" binding:"required"`
	PayeeName    string  `json:"payee_name" form:"payee_name" binding:"required"`
	PayeeAccount string  `json:"payee_account" form:"payee_account" binding:"required"`
	PayeeChannel string  `json:"payee_channel" form:"payee_channel"`
}

type SellerWithdrawListReq struct {
	BasePage
	SellerID uint   `json:"seller_id" form:"seller_id"`
	Status   string `json:"status" form:"status"`
}

type AdminSellerWithdrawListReq = SellerWithdrawListReq

type SellerWithdrawAuditReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	Reason string `json:"reason" form:"reason"`
}

type AdminSellerWithdrawAuditReq = SellerWithdrawAuditReq

type SellerWithdrawPaidReq struct {
	ID     uint   `json:"id" form:"id" binding:"required"`
	Status string `json:"status" form:"status" binding:"required"`
	Reason string `json:"reason" form:"reason"`
}

type AdminSellerWithdrawPaidReq = SellerWithdrawPaidReq

type SellerWithdrawResp struct {
	ID                uint    `json:"id"`
	SellerID          uint    `json:"seller_id"`
	UserName          string  `json:"user_name"`
	NickName          string  `json:"nick_name"`
	ShopName          string  `json:"shop_name"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	StatusText        string  `json:"status_text"`
	PayeeName         string  `json:"payee_name"`
	PayeeAccount      string  `json:"payee_account"`
	PayeeChannel      string  `json:"payee_channel"`
	AuditReason       string  `json:"audit_reason"`
	AuditOperatorID   uint    `json:"audit_operator_id"`
	AuditOperatorName string  `json:"audit_operator_name"`
	PaidOperatorID    uint    `json:"paid_operator_id"`
	PaidOperatorName  string  `json:"paid_operator_name"`
	CreatedAt         int64   `json:"created_at"`
	AuditedAt         int64   `json:"audited_at"`
	PaidAt            int64   `json:"paid_at"`
}
