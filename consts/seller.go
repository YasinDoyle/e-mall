package consts

const (
	SellerStatusPending  = 0
	SellerStatusApproved = 1
	SellerStatusRejected = 2
	SellerStatusBanned   = 3
)

var SellerStatusMap = map[uint]string{
	SellerStatusPending:  "待审核",
	SellerStatusApproved: "已通过",
	SellerStatusRejected: "已拒绝",
	SellerStatusBanned:   "已封禁",
}
