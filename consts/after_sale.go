package consts

const (
	AfterSaleTypeRefundOnly   = "refund_only"
	AfterSaleTypeReturnRefund = "return_refund"

	AfterSaleStatusRequested           = "requested"
	AfterSaleStatusSellerApproved      = "seller_approved"
	AfterSaleStatusSellerRejected      = "seller_rejected"
	AfterSaleStatusPlatformIntervening = "platform_intervening"
	AfterSaleStatusRefunded            = "refunded"
	AfterSaleStatusClosed              = "closed"

	AfterSaleActionApprove   = "approve"
	AfterSaleActionReject    = "reject"
	AfterSaleActionIntervene = "intervene"
	AfterSaleActionRefund    = "refund"
	AfterSaleActionClose     = "close"
)
