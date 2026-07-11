package types

type OrderServiceReq struct {
	OrderId   uint `form:"order_id" json:"order_id"`
	ProductID uint `form:"product_id" json:"product_id"`
	Num       uint `form:"num" json:"num"`
	AddressID uint `form:"address_id" json:"address_id"`
	Money     int  `form:"money" json:"money"`
	BossID    uint `form:"boss_id" json:"boss_id"`
	UserID    uint `form:"user_id" json:"user_id"`
	OrderNum  uint `form:"order_num" json:"order_num"`
	Type      int  `form:"type" json:"type"`
	*BasePage
}

type OrderCreateReq struct {
	OrderId   uint `form:"order_id" json:"order_id"`
	ProductID uint `form:"product_id" json:"product_id"`
	Num       uint `form:"num" json:"num"`
	AddressID uint `form:"address_id" json:"address_id"`
	Money     int  `form:"money" json:"money"`
	BossID    uint `form:"boss_id" json:"boss_id"`
	UserID    uint `form:"user_id" json:"user_id"`
	OrderNum  uint `form:"order_num" json:"order_num"`
	Type      int  `form:"type" json:"type"`
	CouponID  uint `form:"coupon_id" json:"coupon_id"` // 可选：使用优惠券
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

type OrderShowReq struct {
	OrderId uint `json:"order_id" form:"order_id"`
}

type OrderDeleteReq struct {
	OrderId uint `json:"order_id" form:"order_id"`
}

type OrderShipReq struct {
	OrderId    uint   `json:"order_id" form:"order_id"`
	TrackingNo string `json:"tracking_no" form:"tracking_no"`
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
	ID            uint    `json:"id"`
	OrderNum      uint64  `json:"order_num"`
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
	UserID        uint    `json:"user_id"`
	ProductID     uint    `json:"product_id"`
	BossID        uint    `json:"boss_id"`
	Num           uint    `json:"num"`
	AddressName   string  `json:"address_name"`
	AddressPhone  string  `json:"address_phone"`
	Address       string  `json:"address"`
	Type          uint    `json:"type"`
	Money         float64 `json:"money"`
	RefundStatus  int     `json:"refund_status"`
	RefundReason  string  `json:"refund_reason"`
	TrackingNo    string  `json:"tracking_no"`
	Name          string  `json:"name"`
	ImgPath       string  `json:"img_path"`
	DiscountPrice string  `json:"discount_price"`
}
