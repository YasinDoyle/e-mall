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
		{"self purchase forbidden", ErrorOrderSelfPurchaseForbidden, "不能购买自己发布的商品"},
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

func TestGetMsgByLocaleFallsBackToZhCNAndKeepsCodeStable(t *testing.T) {
	code := ErrorSellerNotApproved

	if got := GetMsgByLocale(code, "zh-CN"); got != GetMsg(code) {
		t.Fatalf("expected zh-CN message %q, got %q", GetMsg(code), got)
	}
	if got := GetMsgByLocale(code, "fr-FR"); got != GetMsg(code) {
		t.Fatalf("expected unsupported locale to fall back to %q, got %q", GetMsg(code), got)
	}

	err := NewBusinessError(code)
	discoveredCode, ok := CodeFromError(err)
	if !ok {
		t.Fatal("expected business error code to remain discoverable")
	}
	if discoveredCode != code {
		t.Fatalf("expected code %d, got %d", code, discoveredCode)
	}
}

func TestGetMsgKeyReturnsStableBusinessErrorKey(t *testing.T) {
	key := GetMsgKey(ErrorSellerNotApproved)
	if key != "seller.not_approved" {
		t.Fatalf("expected seller.not_approved key, got %q", key)
	}
	if fallback := GetMsgKey(0); fallback != "common.error" {
		t.Fatalf("expected common.error fallback key, got %q", fallback)
	}
}

func TestGetMsgByLocaleReturnsEnglishMessageAndStableKey(t *testing.T) {
	code := ErrorSellerNotApproved

	if got := GetMsgByLocale(code, "en-US"); got != "Please complete seller onboarding and pass review first" {
		t.Fatalf("expected english message, got %q", got)
	}
	if key := GetMsgKey(code); key != "seller.not_approved" {
		t.Fatalf("expected stable msg key, got %q", key)
	}
	if fallback := GetMsgByLocale(0, "en-US"); fallback != "fail" {
		t.Fatalf("expected english fallback message, got %q", fallback)
	}
}
