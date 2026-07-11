package consts

const (
	OrderTypeUnPaid = iota + 1
	OrderTypePendingShipping
	OrderTypeShipping
	OrderTypeReceipt
	OrderTypeRefundRequested
	OrderTypeRefunded
)

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
}
