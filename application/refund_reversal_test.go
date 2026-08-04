package application

import (
	"testing"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestBuildRefundAccountFlowsAreStablePerOrder(t *testing.T) {
	order := &model.Order{Model: gorm.Model{ID: 1}, OrderNum: 9001, UserID: 7, BossID: 8}
	first := buildRefundAccountFlows(order, 66.50)
	second := buildRefundAccountFlows(order, 66.50)
	if len(first) != 3 || len(second) != 3 {
		t.Fatalf("expected 3 refund reversal flows, got %d and %d", len(first), len(second))
	}
	seen := map[string]bool{}
	for i, flow := range first {
		if seen[flow.FlowNo] {
			t.Fatalf("duplicate flow no %s", flow.FlowNo)
		}
		seen[flow.FlowNo] = true
		if flow.FlowNo != second[i].FlowNo {
			t.Fatalf("expected stable flow no at %d, got %s and %s", i, flow.FlowNo, second[i].FlowNo)
		}
	}
}

func TestApplySellerSettlementRefundDebitDebitsAvailableAndIncome(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 90,
		FrozenBalance:    5,
		TotalIncome:      100,
		TotalWithdrawn:   10,
	}

	if err := applySellerSettlementRefundDebit(account, 30); err != nil {
		t.Fatalf("expected refund debit to pass, got %v", err)
	}

	if account.AvailableBalance != 60 {
		t.Fatalf("expected available balance 60, got %.2f", account.AvailableBalance)
	}
	if account.FrozenBalance != 5 {
		t.Fatalf("expected frozen balance to remain 5, got %.2f", account.FrozenBalance)
	}
	if account.TotalIncome != 70 {
		t.Fatalf("expected total income 70, got %.2f", account.TotalIncome)
	}
	if account.TotalWithdrawn != 10 {
		t.Fatalf("expected total withdrawn to remain 10, got %.2f", account.TotalWithdrawn)
	}
}

func TestApplySellerSettlementRefundDebitAllowsFutureBalanceOffset(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 10,
		TotalIncome:      10,
	}

	if err := applySellerSettlementRefundDebit(account, 30); err != nil {
		t.Fatalf("expected platform refund debit to allow negative available balance, got %v", err)
	}
	if account.AvailableBalance != -20 {
		t.Fatalf("expected available balance -20, got %.2f", account.AvailableBalance)
	}
	if account.TotalIncome != -20 {
		t.Fatalf("expected total income -20, got %.2f", account.TotalIncome)
	}
}

func TestApplySellerSettlementRefundDebitRejectsInvalidAmount(t *testing.T) {
	err := applySellerSettlementRefundDebit(&model.SellerAccount{}, 0)
	if e.BusinessCode(err) != e.ErrorSellerWithdrawAmountInvalid {
		t.Fatalf("expected invalid amount error, got %d (%v)", e.BusinessCode(err), err)
	}
}
