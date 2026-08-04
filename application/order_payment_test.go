package application

import (
	"errors"
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/utils/e"
	"gorm.io/gorm"
)

func TestValidateOrderPaymentCallbackRejectsAmountMismatch(t *testing.T) {
	payment := &model.OrderPayment{OrderID: 1, Amount: 88.66, Status: consts.OrderPaymentStatusPending, Channel: consts.OrderPaymentChannelWechat}
	if err := validateOrderPaymentCallback(payment, 88.65, consts.OrderPaymentChannelWechat); e.BusinessCode(err) != e.ErrorPaymentAmountMismatch {
		t.Fatalf("expected amount mismatch code, got %d (%v)", e.BusinessCode(err), err)
	}
}

func TestValidateOrderPaymentCallbackIsIdempotentForPaid(t *testing.T) {
	payment := &model.OrderPayment{OrderID: 1, Amount: 88.66, Status: consts.OrderPaymentStatusPaid, Channel: consts.OrderPaymentChannelWechat}
	if err := validateOrderPaymentCallback(payment, 88.66, consts.OrderPaymentChannelWechat); err != nil {
		t.Fatal(err)
	}
}

func TestValidateOrderPaymentCallbackRejectsChannelMismatch(t *testing.T) {
	payment := &model.OrderPayment{OrderID: 1, Amount: 88.66, Status: consts.OrderPaymentStatusPending, Channel: consts.OrderPaymentChannelWechat}
	if err := validateOrderPaymentCallback(payment, 88.66, consts.OrderPaymentChannelAlipay); e.BusinessCode(err) != e.InvalidParams {
		t.Fatalf("expected invalid params code for channel mismatch, got %d (%v)", e.BusinessCode(err), err)
	}
}

func TestIsOrderPaymentNoRequiresOPFollowedByDigits(t *testing.T) {
	cases := []struct {
		name      string
		paymentNo string
		want      bool
	}{
		{name: "valid product payment", paymentNo: "OP123", want: true},
		{name: "recharge order", paymentNo: "R123", want: false},
		{name: "bare prefix", paymentNo: "OP", want: false},
		{name: "malformed suffix", paymentNo: "OPABC", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsOrderPaymentNo(tc.paymentNo); got != tc.want {
				t.Fatalf("IsOrderPaymentNo(%q) = %v, want %v", tc.paymentNo, got, tc.want)
			}
		})
	}
}

func TestPendingOrderPaymentDecisionCreatesWhenNoPendingPayment(t *testing.T) {
	payment, createNew, err := resolvePendingOrderGatewayPayment(nil, gorm.ErrRecordNotFound, consts.OrderPaymentChannelWechat)
	if err != nil {
		t.Fatal(err)
	}
	if payment != nil || !createNew {
		t.Fatalf("expected no reusable payment and createNew=true, got payment=%v createNew=%v", payment, createNew)
	}
}

func TestPendingOrderPaymentDecisionReusesSameChannelPendingPayment(t *testing.T) {
	pending := &model.OrderPayment{PaymentNo: "OP123", Channel: consts.OrderPaymentChannelWechat, Status: consts.OrderPaymentStatusPending}

	payment, createNew, err := resolvePendingOrderGatewayPayment(pending, nil, consts.OrderPaymentChannelWechat)
	if err != nil {
		t.Fatal(err)
	}
	if payment != pending || createNew {
		t.Fatalf("expected same-channel pending payment to be reused, got payment=%v createNew=%v", payment, createNew)
	}
}

func TestPendingOrderPaymentDecisionRejectsDifferentChannelPendingPayment(t *testing.T) {
	pending := &model.OrderPayment{PaymentNo: "OP123", Channel: consts.OrderPaymentChannelWechat, Status: consts.OrderPaymentStatusPending}

	_, _, err := resolvePendingOrderGatewayPayment(pending, nil, consts.OrderPaymentChannelAlipay)
	if e.BusinessCode(err) != e.ErrorOrderPayStatusInvalid {
		t.Fatalf("expected order pay status invalid for cross-channel pending payment, got %d (%v)", e.BusinessCode(err), err)
	}
}

func TestPendingOrderPaymentDecisionPropagatesLookupError(t *testing.T) {
	lookupErr := errors.New("db failed")

	_, _, err := resolvePendingOrderGatewayPayment(nil, lookupErr, consts.OrderPaymentChannelWechat)
	if !errors.Is(err, lookupErr) {
		t.Fatalf("expected lookup error to propagate, got %v", err)
	}
}

func TestBalancePaymentRejectsExistingExternalPendingPayment(t *testing.T) {
	pending := &model.OrderPayment{PaymentNo: "OP123", Channel: consts.OrderPaymentChannelWechat, Status: consts.OrderPaymentStatusPending}

	err := ensureNoPendingExternalPaymentForBalance(pending, nil)
	if e.BusinessCode(err) != e.ErrorOrderPayStatusInvalid {
		t.Fatalf("expected order pay status invalid when external payment is pending, got %d (%v)", e.BusinessCode(err), err)
	}
}

func TestBalancePaymentAllowsMissingPendingPayment(t *testing.T) {
	if err := ensureNoPendingExternalPaymentForBalance(nil, gorm.ErrRecordNotFound); err != nil {
		t.Fatal(err)
	}
}

func TestBalancePaymentPendingLookupErrorPropagates(t *testing.T) {
	lookupErr := errors.New("db failed")

	err := ensureNoPendingExternalPaymentForBalance(nil, lookupErr)
	if !errors.Is(err, lookupErr) {
		t.Fatalf("expected lookup error to propagate, got %v", err)
	}
}
