package application

import (
	"testing"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
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
