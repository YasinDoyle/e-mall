package service

import (
	"testing"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestValidateSellerWithdrawApplyReqRejectsInvalidAmount(t *testing.T) {
	req := &types.SellerWithdrawApplyReq{
		Amount:       0,
		PayeeName:    "程小红",
		PayeeAccount: "6222000000000000",
	}
	err := validateSellerWithdrawApplyReq(req)
	if err == nil {
		t.Fatal("expected zero amount to fail")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerWithdrawAmountInvalid)
}

func TestValidateSellerWithdrawApplyReqRejectsMissingPayeeAccount(t *testing.T) {
	req := &types.SellerWithdrawApplyReq{
		Amount:       100,
		PayeeName:    "程小红",
		PayeeAccount: "  ",
	}
	err := validateSellerWithdrawApplyReq(req)
	if err == nil {
		t.Fatal("expected blank payee account to fail")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerWithdrawPayeeRequired)
}

func TestApplySellerWithdrawFreezeMovesAvailableToFrozen(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 100,
		FrozenBalance:    5,
		TotalIncome:      100,
		TotalWithdrawn:   0,
	}
	if err := applySellerWithdrawFreeze(account, 30); err != nil {
		t.Fatalf("expected freeze to pass, got %v", err)
	}
	if account.AvailableBalance != 70 {
		t.Fatalf("expected available balance 70, got %.2f", account.AvailableBalance)
	}
	if account.FrozenBalance != 35 {
		t.Fatalf("expected frozen balance 35, got %.2f", account.FrozenBalance)
	}
}

func TestApplySellerWithdrawFreezeRejectsInsufficientBalance(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 20,
		FrozenBalance:    0,
		TotalIncome:      20,
		TotalWithdrawn:   0,
	}
	err := applySellerWithdrawFreeze(account, 30)
	if err == nil {
		t.Fatal("expected insufficient balance to fail")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerWithdrawInsufficientBalance)
}

func TestApplySellerWithdrawPaidConsumesFrozenBalance(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 70,
		FrozenBalance:    30,
		TotalIncome:      100,
		TotalWithdrawn:   0,
	}
	if err := applySellerWithdrawPaid(account, 30); err != nil {
		t.Fatalf("expected paid transition to pass, got %v", err)
	}
	if account.AvailableBalance != 70 {
		t.Fatalf("expected available balance 70, got %.2f", account.AvailableBalance)
	}
	if account.FrozenBalance != 0 {
		t.Fatalf("expected frozen balance 0, got %.2f", account.FrozenBalance)
	}
	if account.TotalWithdrawn != 30 {
		t.Fatalf("expected total withdrawn 30, got %.2f", account.TotalWithdrawn)
	}
}

func TestApplySellerWithdrawPaidRejectsMissingFrozenBalance(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 70,
		FrozenBalance:    0,
		TotalIncome:      100,
		TotalWithdrawn:   0,
	}
	err := applySellerWithdrawPaid(account, 30)
	if err == nil {
		t.Fatal("expected repeated payout to fail")
	}
	assertServiceBusinessCode(t, err, e.ErrorSellerWithdrawStatusInvalid)
}

func TestApplySellerWithdrawRejectedRestoresAvailableBalance(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 70,
		FrozenBalance:    30,
		TotalIncome:      100,
		TotalWithdrawn:   0,
	}
	if err := applySellerWithdrawRejected(account, 30); err != nil {
		t.Fatalf("expected reject transition to pass, got %v", err)
	}
	if account.AvailableBalance != 100 {
		t.Fatalf("expected available balance 100, got %.2f", account.AvailableBalance)
	}
	if account.FrozenBalance != 0 {
		t.Fatalf("expected frozen balance 0, got %.2f", account.FrozenBalance)
	}
}

func TestCreditSellerSettlementAmountCreatesSellerIncome(t *testing.T) {
	account := &model.SellerAccount{
		AvailableBalance: 10,
		FrozenBalance:    0,
		TotalIncome:      10,
		TotalWithdrawn:   0,
	}
	if err := applySellerSettlementCredit(account, 80); err != nil {
		t.Fatalf("expected settlement credit to pass, got %v", err)
	}
	if account.AvailableBalance != 90 {
		t.Fatalf("expected available balance 90, got %.2f", account.AvailableBalance)
	}
	if account.TotalIncome != 90 {
		t.Fatalf("expected total income 90, got %.2f", account.TotalIncome)
	}
}

func TestBuildSellerSettlementCreditBackfillPlanSkipsAlreadyCreditedSettlements(t *testing.T) {
	settlements := []*model.Settlement{
		{Model: gorm.Model{ID: 1}, SellerID: 7, SettlementAmount: 94.81, Status: model.SettlementStatusPaid},
		{Model: gorm.Model{ID: 2}, SellerID: 7, SettlementAmount: 5.19, Status: model.SettlementStatusPaid},
		{Model: gorm.Model{ID: 3}, SellerID: 8, SettlementAmount: 10, Status: model.SettlementStatusPaid},
	}
	credited := map[uint]struct{}{
		2: {},
	}

	plan := buildSellerSettlementCreditBackfillPlan(settlements, credited)

	if len(plan) != 2 {
		t.Fatalf("expected 2 sellers in backfill plan, got %d", len(plan))
	}
	if len(plan[7]) != 1 || plan[7][0].ID != 1 {
		t.Fatalf("expected seller 7 to backfill settlement 1 only, got %+v", plan[7])
	}
	if len(plan[8]) != 1 || plan[8][0].ID != 3 {
		t.Fatalf("expected seller 8 to backfill settlement 3, got %+v", plan[8])
	}
}
