package service

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestEnsureSellerProfileApproved(t *testing.T) {
	profile := &model.SellerProfile{Status: consts.SellerStatusApproved}
	if err := ensureSellerProfileApproved(profile); err != nil {
		t.Fatalf("expected approved seller profile to pass, got %v", err)
	}
}

func TestEnsureSellerProfileApprovedRejectsPending(t *testing.T) {
	profile := &model.SellerProfile{Status: consts.SellerStatusPending}
	err := ensureSellerProfileApproved(profile)
	if err == nil {
		t.Fatal("expected pending seller profile to be rejected")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerNotApproved)
}

func TestEnsureSellerCanChangeSaleStatusRequiresApprovedSellerWhenEnabling(t *testing.T) {
	profile := &model.SellerProfile{Status: consts.SellerStatusPending}
	err := ensureSellerCanChangeSaleStatus(profile, true)
	if err == nil {
		t.Fatal("expected enabling sale with pending seller profile to be rejected")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerNotApproved)
}

func TestEnsureSellerCanChangeSaleStatusAllowsDisablingWithoutSellerProfile(t *testing.T) {
	if err := ensureSellerCanChangeSaleStatus(nil, false); err != nil {
		t.Fatalf("expected disabling sale without seller profile to pass, got %v", err)
	}
}

func TestEnsureSellerCanEnableTrading(t *testing.T) {
	user := &model.User{}
	if err := user.SetInitialMoneyWithPayKey("123456"); err != nil {
		t.Fatalf("expected pay key setup to pass, got %v", err)
	}
	if err := ensureSellerCanEnableTrading(user, true); err != nil {
		t.Fatalf("expected seller with pay key to enable trading, got %v", err)
	}
}

func TestEnsureSellerCanEnableTradingRejectsMissingPayKey(t *testing.T) {
	user := &model.User{}
	err := ensureSellerCanEnableTrading(user, true)
	if err == nil {
		t.Fatal("expected seller without pay key to be rejected when enabling trading")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerPayKeyRequired)
}

func TestEnsureSellerCanEnableTradingAllowsDisableWithoutPayKey(t *testing.T) {
	user := &model.User{}
	if err := ensureSellerCanEnableTrading(user, false); err != nil {
		t.Fatalf("expected disabling trading without pay key to pass, got %v", err)
	}
}

func TestEnsureNotBuyingOwnProductRejectsSelfPurchase(t *testing.T) {
	err := ensureNotBuyingOwnProduct(1, 1)
	if err == nil {
		t.Fatal("expected self purchase to be rejected")
	}
	assertServiceBusinessCode(t, err, e.ErrorOrderSelfPurchaseForbidden)
}

func TestEnsureNotBuyingOwnProductAllowsDifferentBuyerAndSeller(t *testing.T) {
	if err := ensureNotBuyingOwnProduct(1, 2); err != nil {
		t.Fatalf("expected different buyer and seller to pass, got %v", err)
	}
}

func assertServiceBusinessCode(t *testing.T, err error, want int) {
	t.Helper()
	code, ok := e.CodeFromError(err)
	if !ok {
		t.Fatalf("expected business error code %d, got plain error %v", want, err)
	}
	if code != want {
		t.Fatalf("expected business error code %d, got %d", want, code)
	}
}
