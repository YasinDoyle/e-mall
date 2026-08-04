package service

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestCalculateOrderSettlementUsesGlobalCommissionRate(t *testing.T) {
	result, err := calculateOrderSettlement(200, 2, 0.05)
	if err != nil {
		t.Fatalf("expected valid settlement calculation, got %v", err)
	}

	if result.GrossAmount != 400 {
		t.Fatalf("expected gross amount 400, got %.2f", result.GrossAmount)
	}
	if result.CommissionAmount != 20 {
		t.Fatalf("expected commission amount 20, got %.2f", result.CommissionAmount)
	}
	if result.SettlementAmount != 380 {
		t.Fatalf("expected settlement amount 380, got %.2f", result.SettlementAmount)
	}
}

func TestCalculateOrderSettlementRejectsInvalidRate(t *testing.T) {
	_, err := calculateOrderSettlement(100, 1, 1.2)
	if err == nil {
		t.Fatal("expected invalid commission rate to fail")
	}
}

func TestCanMarkSettlementPaidRequiresReceivedOrderAndNoActiveAfterSale(t *testing.T) {
	settlement := &model.Settlement{Status: model.SettlementStatusGenerated}
	readyOrder := &model.Order{Type: consts.OrderTypeReceipt, RefundStatus: consts.OrderRefundStatusNone}
	if !canMarkSettlementPaid(settlement, readyOrder, false) {
		t.Fatal("expected received order without active after-sale to be payable")
	}
	if canMarkSettlementPaid(settlement, &model.Order{Type: consts.OrderTypePendingShipping, RefundStatus: consts.OrderRefundStatusNone}, false) {
		t.Fatal("expected paid but unshipped order to be blocked")
	}
	if canMarkSettlementPaid(settlement, readyOrder, true) {
		t.Fatal("expected active after-sale to block settlement payout")
	}
	if canMarkSettlementPaid(&model.Settlement{Status: model.SettlementStatusPending}, readyOrder, false) {
		t.Fatal("expected pending settlement to be blocked")
	}
}
