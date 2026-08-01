package types

type OrderServiceReq struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	ProductID uint    `form:"product_id" json:"product_id"`
	Num       uint    `form:"num" json:"num"`
	AddressID uint    `form:"address_id" json:"address_id"`
	Money     float64 `form:"money" json:"money"`
	BossID    uint    `form:"boss_id" json:"boss_id"`
	UserID    uint    `form:"user_id" json:"user_id"`
	OrderNum  uint    `form:"order_num" json:"order_num"`
	Type      int     `form:"type" json:"type"`
	*BasePage
}

type OrderCreateReq struct {
	OrderId   uint    `form:"order_id" json:"order_id"`
	ProductID uint    `form:"product_id" json:"product_id"`
	Num       uint    `form:"num" json:"num"`
	AddressID uint    `form:"address_id" json:"address_id"`
	Money     float64 `form:"money" json:"money"`
	BossID    uint    `form:"boss_id" json:"boss_id"`
	UserID    uint    `form:"user_id" json:"user_id"`
	OrderNum  uint    `form:"order_num" json:"order_num"`
	Type      int     `form:"type" json:"type"`
	CouponID  uint    `form:"coupon_id" json:"coupon_id"` // 可选：使用优惠券
}

type OrderCreateResp struct {
	ID       uint    `json:"id"`
	OrderNum uint64  `json:"order_num"`
	Money    float64 `json:"money"`
	CouponID uint    `json:"coupon_id"`
}

type OrderListReq struct {
	Type int `form:"type" json:"type"`
	BasePage
}

type SellerOrderListReq struct {
	Type int `form:"type" json:"type"`
	BasePage
}

type AdminOrderListReq struct {
	Type         int  `form:"type" json:"type"`
	RefundStatus *int `form:"refund_status" json:"refund_status"`
	BasePage
}

type OrderShowReq struct {
	OrderId uint `json:"order_id" form:"order_id"`
}

type OrderDeleteReq struct {
	OrderId uint `json:"order_id" form:"order_id"`
}

type OrderShipReq struct {
	OrderId          uint   `json:"order_id" form:"order_id"`
	LogisticsCompany string `json:"logistics_company" form:"logistics_company"`
	TrackingNo       string `json:"tracking_no" form:"tracking_no"`
}

type ShipmentInfo struct {
	LogisticsCompany string `json:"logistics_company"`
	TrackingNo       string `json:"tracking_no"`
}

type OrderReceiveReq struct {
	OrderId uint `json:"order_id" form:"order_id"`
}

type OrderRefundRequestReq struct {
	OrderId uint   `json:"order_id" form:"order_id" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

type AdminOrderRefundApproveReq struct {
	OrderId uint   `json:"order_id" form:"order_id" binding:"required"`
	Key     string `json:"key" form:"key" binding:"required"`
}

type OrderRefundResp struct {
	OrderId      uint    `json:"order_id"`
	OrderNum     uint64  `json:"order_num"`
	RefundAmount float64 `json:"refund_amount"`
	RefundStatus int     `json:"refund_status"`
	Type         uint    `json:"type"`
}

type OrderListResp struct {
	ID               uint    `json:"id"`
	OrderNum         uint64  `json:"order_num"`
	CreatedAt        int64   `json:"created_at"`
	UpdatedAt        int64   `json:"updated_at"`
	UserID           uint    `json:"user_id"`
	ProductID        uint    `json:"product_id"`
	BossID           uint    `json:"boss_id"`
	Num              uint    `json:"num"`
	AddressName      string  `json:"address_name"`
	AddressPhone     string  `json:"address_phone"`
	Address          string  `json:"address"`
	Type             uint    `json:"type"`
	Money            float64 `json:"money"`
	RefundStatus     int     `json:"refund_status"`
	RefundReason     string  `json:"refund_reason"`
	PaymentChannel   string  `json:"payment_channel"`
	LogisticsCompany string  `json:"logistics_company"`
	TrackingNo       string  `json:"tracking_no"`
	ShippedAt        int64   `json:"shipped_at"`
	ReceivedAt       int64   `json:"received_at"`
	CanceledAt       int64   `json:"canceled_at"`
	Name             string  `json:"name"`
	ImgPath          string  `json:"img_path"`
	DiscountPrice    string  `json:"discount_price"`
	GrossAmount      float64 `json:"gross_amount"`
	CommissionAmount float64 `json:"commission_amount"`
	SettlementAmount float64 `json:"settlement_amount"`
	SettlementStatus string  `json:"settlement_status"`
}

// OrderLogResp is an order fulfillment operation audit record.
type OrderLogResp struct {
	ID           uint   `json:"id"`
	OrderID      uint   `json:"order_id"`
	OrderNum     uint64 `json:"order_num"`
	Action       string `json:"action"`
	FromType     uint   `json:"from_type"`
	ToType       uint   `json:"to_type"`
	OperatorType string `json:"operator_type"`
	OperatorID   uint   `json:"operator_id"`
	Remark       string `json:"remark"`
	CreatedAt    int64  `json:"created_at"`
}

type OrderLogisticsResp struct {
	ID          uint   `json:"id"`
	OrderID     uint   `json:"order_id"`
	OrderNum    uint64 `json:"order_num"`
	NodeType    string `json:"node_type"`
	Description string `json:"description"`
	OccurredAt  int64  `json:"occurred_at"`
}
