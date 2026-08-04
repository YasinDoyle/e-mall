package application

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestBuildOrderLogRejectsIllegalTransition(t *testing.T) {
	_, err := buildOrderLog(1, 1001, consts.OrderActionShip, consts.OrderTypeUnPaid, consts.OrderTypeShipping, "seller", 2, "ship unpaid")
	if err == nil {
		t.Fatal("expected illegal transition to fail")
	}
	if e.BusinessCode(err) != e.ErrorOrderStatusTransitionInvalid {
		t.Fatalf("expected transition business code, got %d", e.BusinessCode(err))
	}
}

func TestBuildOrderLogCapturesOperator(t *testing.T) {
	log, err := buildOrderLog(1, 1001, consts.OrderActionPay, consts.OrderTypeUnPaid, consts.OrderTypePendingShipping, "buyer", 7, "balance")
	if err != nil {
		t.Fatal(err)
	}
	if log.OperatorType != "buyer" || log.OperatorID != 7 || log.Action != consts.OrderActionPay {
		t.Fatalf("unexpected log: %+v", log)
	}
}

func TestBuildOrderLogCapturesAdminOperator(t *testing.T) {
	log, err := buildOrderLog(1, 1001, consts.OrderActionRefundApprove, consts.OrderTypeRefundRequested, consts.OrderTypeRefunded, "admin", 9, "approve refund")
	if err != nil {
		t.Fatal(err)
	}
	if log.OperatorType != "admin" || log.OperatorID != 9 {
		t.Fatalf("expected admin operator to be captured, got %+v", log)
	}
}
