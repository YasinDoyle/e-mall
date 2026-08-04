package aftersale

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestEnsureTransitionAllowsP2Flow(t *testing.T) {
	tests := [][2]string{
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusSellerApproved},
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusSellerRejected},
		{consts.AfterSaleStatusSellerRejected, consts.AfterSaleStatusPlatformIntervening},
		{consts.AfterSaleStatusSellerApproved, consts.AfterSaleStatusRefunded},
		{consts.AfterSaleStatusPlatformIntervening, consts.AfterSaleStatusRefunded},
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusClosed},
	}
	for _, tt := range tests {
		if err := EnsureTransition(tt[0], tt[1]); err != nil {
			t.Fatalf("expected %s -> %s to pass: %v", tt[0], tt[1], err)
		}
	}
}

func TestEnsureTransitionRejectsIllegalFlow(t *testing.T) {
	tests := [][2]string{
		{consts.AfterSaleStatusRefunded, consts.AfterSaleStatusRequested},
		{consts.AfterSaleStatusSellerApproved, consts.AfterSaleStatusClosed},
	}
	for _, tt := range tests {
		err := EnsureTransition(tt[0], tt[1])
		if e.BusinessCode(err) != e.ErrorAfterSaleStatusInvalid {
			t.Fatalf("expected after-sale status error for %s -> %s, got %d (%v)", tt[0], tt[1], e.BusinessCode(err), err)
		}
	}
}
