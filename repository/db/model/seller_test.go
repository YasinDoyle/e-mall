package model

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
)

func TestSellerProfileIsApproved(t *testing.T) {
	seller := &SellerProfile{Status: consts.SellerStatusApproved}
	if !seller.IsApproved() {
		t.Fatal("approved seller should be allowed to operate seller features")
	}
}

func TestSellerProfileIsApprovedRejectsPending(t *testing.T) {
	seller := &SellerProfile{Status: consts.SellerStatusPending}
	if seller.IsApproved() {
		t.Fatal("pending seller should not be allowed to operate seller features")
	}
}
