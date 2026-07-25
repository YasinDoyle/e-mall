package types

type SellerApplyReq struct {
	ShopName    string `json:"shop_name" binding:"required"`
	Description string `json:"description"`
}

type SellerProfileResp struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	UserName     string `json:"user_name"`
	ShopName     string `json:"shop_name"`
	Description  string `json:"description"`
	Status       uint   `json:"status"`
	StatusText   string `json:"status_text"`
	RejectReason string `json:"reject_reason"`
	ApprovedAt   int64  `json:"approved_at"`
	CreatedAt    int64  `json:"created_at"`
}
