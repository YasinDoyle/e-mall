package service

import (
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestEnsureSellerCanEnableTrading(t *testing.T) {
	user := &model.User{PayKeySet: true}
	if err := ensureSellerCanEnableTrading(user, true); err != nil {
		t.Fatalf("expected seller with pay key to enable trading, got %v", err)
	}
}

func TestEnsureSellerCanEnableTradingRejectsMissingPayKey(t *testing.T) {
	user := &model.User{}
	if err := ensureSellerCanEnableTrading(user, true); err == nil {
		t.Fatal("expected seller without pay key to be rejected when enabling trading")
	}
}

func TestEnsureSellerCanEnableTradingAllowsDisableWithoutPayKey(t *testing.T) {
	user := &model.User{}
	if err := ensureSellerCanEnableTrading(user, false); err != nil {
		t.Fatalf("expected disabling trading without pay key to pass, got %v", err)
	}
}
