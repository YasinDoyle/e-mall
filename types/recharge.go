package types

// RechargeReq 发起充值请求
type RechargeReq struct {
	Amount float64 `json:"amount" binding:"required,gt=0"` // 充值金额（元）
}

// RechargeResp 充值发起响应
type RechargeResp struct {
	OrderNum  string `json:"order_num"`             // 充值订单号
	QRCodeURL string `json:"qr_code_url,omitempty"` // 微信：二维码 URL
	PayURL    string `json:"pay_url,omitempty"`     // 支付宝：跳转链接
}

// RechargeStatusResp 充值状态查询响应
type RechargeStatusResp struct {
	OrderNum      string  `json:"order_num"`
	Status        string  `json:"status"` // pending / paid / credited / failed
	Amount        float64 `json:"amount"`
	Channel       string  `json:"channel,omitempty"`
	PendingCredit float64 `json:"pending_credit,omitempty"`
	RefundStatus  string  `json:"refund_status,omitempty"`
}

type RechargeRefundReq struct {
	OrderNum string  `json:"order_num" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Reason   string  `json:"reason"`
}

type RechargeRefundResp struct {
	OrderNum     string  `json:"order_num"`
	RefundNo     string  `json:"refund_no"`
	Amount       float64 `json:"amount"`
	RefundStatus string  `json:"refund_status"`
	ProviderID   string  `json:"provider_id,omitempty"`
}
