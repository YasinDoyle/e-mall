package e

import (
	"errors"
	"testing"
)

func TestNewBusinessErrorCarriesCodeAndMessage(t *testing.T) {
	err := NewBusinessError(ErrorSellerNotApproved)

	code, ok := CodeFromError(err)
	if !ok {
		t.Fatal("expected business error code to be discoverable")
	}
	if code != ErrorSellerNotApproved {
		t.Fatalf("expected code %d, got %d", ErrorSellerNotApproved, code)
	}
	if err.Error() != GetMsg(ErrorSellerNotApproved) {
		t.Fatalf("expected default message %q, got %q", GetMsg(ErrorSellerNotApproved), err.Error())
	}
}

func TestCodeFromErrorUnwrapsWrappedBusinessError(t *testing.T) {
	err := NewBusinessError(ErrorSellerAuditPending)
	wrapped := errors.Join(errors.New("outer"), err)

	code, ok := CodeFromError(wrapped)
	if !ok {
		t.Fatal("expected wrapped business error code to be discoverable")
	}
	if code != ErrorSellerAuditPending {
		t.Fatalf("expected code %d, got %d", ErrorSellerAuditPending, code)
	}
}

func TestPaymentAndRefundBusinessErrorsHaveStableCodes(t *testing.T) {
	cases := []struct {
		name string
		code int
		msg  string
	}{
		{"order not payable", ErrorOrderPayStatusInvalid, "订单已支付或状态不允许支付"},
		{"pay key required", ErrorPaymentPayKeyRequired, "请先设置支付密码"},
		{"pay key invalid", ErrorPaymentPayKeyInvalid, "支付密码错误"},
		{"balance insufficient", ErrorPaymentBalanceInsufficient, "金币不足"},
		{"stock insufficient", ErrorPaymentStockInsufficient, "库存不足"},
		{"refund status invalid", ErrorRefundStatusInvalid, "订单状态不允许退款审批"},
		{"refund amount invalid", ErrorRefundAmountInvalid, "退款金额不合法"},
		{"refund not found", ErrorRefundNotFound, "退款申请不存在"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewBusinessError(tc.code)
			code, ok := CodeFromError(err)
			if !ok {
				t.Fatal("expected business error code to be discoverable")
			}
			if code != tc.code {
				t.Fatalf("expected code %d, got %d", tc.code, code)
			}
			if err.Error() != tc.msg {
				t.Fatalf("expected message %q, got %q", tc.msg, err.Error())
			}
		})
	}
}
