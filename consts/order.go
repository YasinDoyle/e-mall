package consts

const (
	OrderTypeUnPaid = iota + 1
	OrderTypePendingShipping
	OrderTypeShipping
	OrderTypeReceipt
	OrderTypeRefundRequested
	OrderTypeRefunded
)

const (
	OrderTypeCanceled = 7
)

const (
	OrderActionCreate        = "create"
	OrderActionPay           = "pay"
	OrderActionCancel        = "cancel"
	OrderActionShip          = "ship"
	OrderActionReceive       = "receive"
	OrderActionRefundRequest = "refund_request"
	OrderActionRefundApprove = "refund_approve"
	OrderActionRefundReject  = "refund_reject"
	OrderActionAfterSale     = "after_sale"
)

const (
	OrderPaymentChannelBalance = "balance"
	OrderPaymentChannelWechat  = "wechat"
	OrderPaymentChannelAlipay  = "alipay"

	OrderPaymentStatusPending = "pending"
	OrderPaymentStatusPaid    = "paid"
	OrderPaymentStatusFailed  = "failed"
	OrderPaymentStatusClosed  = "closed"
)

const OrderTimeKey = "OrderTime"

const OrderPaidQueue = "rabbitmq-order-paid-queue"

const RechargePaidQueue = "rabbitmq-recharge-paid-queue"

const (
	OrderRefundStatusNone = iota
	OrderRefundStatusRequested
	OrderRefundStatusRefunded
)

var OrderTypeMap = map[int]string{
	OrderTypeUnPaid:          "未支付",
	OrderTypePendingShipping: "已支付，待发货",
	OrderTypeShipping:        "已发货，待收货",
	OrderTypeReceipt:         "已收货，交易成功",
	OrderTypeRefundRequested: "退款申请中",
	OrderTypeRefunded:        "已退款",
	OrderTypeCanceled:        "已取消",
}
