package service

import "testing"

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
