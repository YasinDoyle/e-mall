package orderstate

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestCanTransitionOrderStatusAllowsP2HappyPath(t *testing.T) {
	tests := []struct {
		name string
		from uint
		to   uint
	}{
		{"same_state_unpaid", consts.OrderTypeUnPaid, consts.OrderTypeUnPaid},
		{"pay", consts.OrderTypeUnPaid, consts.OrderTypePendingShipping},
		{"cancel_unpaid", consts.OrderTypeUnPaid, consts.OrderTypeCanceled},
		{"ship", consts.OrderTypePendingShipping, consts.OrderTypeShipping},
		{"request_refund_before_ship", consts.OrderTypePendingShipping, consts.OrderTypeRefundRequested},
		{"receive", consts.OrderTypeShipping, consts.OrderTypeReceipt},
		{"request_refund_after_ship", consts.OrderTypeShipping, consts.OrderTypeRefundRequested},
		{"request_refund_after_receipt", consts.OrderTypeReceipt, consts.OrderTypeRefundRequested},
		{"refund", consts.OrderTypeRefundRequested, consts.OrderTypeRefunded},
		{"refund_reject_to_pending_shipping", consts.OrderTypeRefundRequested, consts.OrderTypePendingShipping},
		{"refund_reject_to_shipping", consts.OrderTypeRefundRequested, consts.OrderTypeShipping},
		{"refund_reject_to_receipt", consts.OrderTypeRefundRequested, consts.OrderTypeReceipt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsureOrderStatusTransition(tt.from, tt.to); err != nil {
				t.Fatalf("expected transition %d -> %d to pass, got %v", tt.from, tt.to, err)
			}
		})
	}
}

func TestCanTransitionOrderStatusRejectsIllegalTransitions(t *testing.T) {
	tests := []struct {
		name string
		from uint
		to   uint
	}{
		{"ship_unpaid", consts.OrderTypeUnPaid, consts.OrderTypeShipping},
		{"pay_shipped", consts.OrderTypeShipping, consts.OrderTypePendingShipping},
		{"cancel_paid", consts.OrderTypePendingShipping, consts.OrderTypeCanceled},
		{"refund_completed_without_after_sale", consts.OrderTypeReceipt, consts.OrderTypeRefunded},
		{"leave_refunded", consts.OrderTypeRefunded, consts.OrderTypeReceipt},
		{"leave_canceled", consts.OrderTypeCanceled, consts.OrderTypePendingShipping},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureOrderStatusTransition(tt.from, tt.to)
			if err == nil {
				t.Fatalf("expected transition %d -> %d to fail", tt.from, tt.to)
			}
			if code := e.BusinessCode(err); code != e.ErrorOrderStatusTransitionInvalid {
				t.Fatalf("expected ErrorOrderStatusTransitionInvalid, got %d (%v)", code, err)
			}
		})
	}
}

func TestOrderStatusIsTerminal(t *testing.T) {
	if !OrderStatusIsTerminal(consts.OrderTypeCanceled) || !OrderStatusIsTerminal(consts.OrderTypeRefunded) {
		t.Fatal("canceled and refunded must be terminal")
	}
	if OrderStatusIsTerminal(consts.OrderTypeShipping) {
		t.Fatal("shipping must not be terminal")
	}
}
