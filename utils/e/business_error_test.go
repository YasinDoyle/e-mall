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
