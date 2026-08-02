package aftersale

import (
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

var allowedTransitions = map[string]map[string]struct{}{
	consts.AfterSaleStatusRequested: {
		consts.AfterSaleStatusSellerApproved: {},
		consts.AfterSaleStatusSellerRejected: {},
		consts.AfterSaleStatusClosed:         {},
	},
	consts.AfterSaleStatusSellerRejected: {
		consts.AfterSaleStatusPlatformIntervening: {},
		consts.AfterSaleStatusClosed:              {},
	},
	consts.AfterSaleStatusSellerApproved: {
		consts.AfterSaleStatusRefunded: {},
		consts.AfterSaleStatusClosed:   {},
	},
	consts.AfterSaleStatusPlatformIntervening: {
		consts.AfterSaleStatusRefunded: {},
		consts.AfterSaleStatusClosed:   {},
	},
}

func EnsureTransition(from, to string) error {
	if from == to {
		return nil
	}
	if next, ok := allowedTransitions[from]; ok {
		if _, ok = next[to]; ok {
			return nil
		}
	}
	return e.NewBusinessError(e.ErrorAfterSaleStatusInvalid)
}
