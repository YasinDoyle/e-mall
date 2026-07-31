package orderstate

import (
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

var allowedOrderTransitions = map[uint]map[uint]struct{}{
	consts.OrderTypeUnPaid: {
		consts.OrderTypePendingShipping: {},
		consts.OrderTypeCanceled:        {},
	},
	consts.OrderTypePendingShipping: {
		consts.OrderTypeShipping:        {},
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeShipping: {
		consts.OrderTypeReceipt:         {},
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeReceipt: {
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeRefundRequested: {
		consts.OrderTypeRefunded:        {},
		consts.OrderTypePendingShipping: {},
		consts.OrderTypeShipping:        {},
		consts.OrderTypeReceipt:         {},
	},
}

func EnsureOrderStatusTransition(from, to uint) error {
	if from == to {
		return nil
	}
	if next, ok := allowedOrderTransitions[from]; ok {
		if _, ok = next[to]; ok {
			return nil
		}
	}
	return e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
}

func OrderStatusIsTerminal(status uint) bool {
	return status == consts.OrderTypeCanceled || status == consts.OrderTypeRefunded
}
