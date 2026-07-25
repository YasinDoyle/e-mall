package v1

import (
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestErrorResponsePreservesBusinessErrorCode(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	err := e.NewBusinessError(e.ErrorSellerNotApproved)

	resp := ErrorResponse(ctx, err)

	if resp.Status != e.ErrorSellerNotApproved {
		t.Fatalf("expected status %d, got %d", e.ErrorSellerNotApproved, resp.Status)
	}
	if resp.Msg != e.GetMsg(e.ErrorSellerNotApproved) {
		t.Fatalf("expected msg %q, got %q", e.GetMsg(e.ErrorSellerNotApproved), resp.Msg)
	}
	if resp.Data != e.GetMsg(e.ErrorSellerNotApproved) {
		t.Fatalf("expected data %q, got %v", e.GetMsg(e.ErrorSellerNotApproved), resp.Data)
	}
	if resp.Error != e.GetMsg(e.ErrorSellerNotApproved) {
		t.Fatalf("expected error %q, got %q", e.GetMsg(e.ErrorSellerNotApproved), resp.Error)
	}
}
