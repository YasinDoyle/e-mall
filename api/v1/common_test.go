package v1

import (
	"net/http"
	"net/http/httptest"
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

func TestErrorResponseLocalizesBusinessErrorMessageFromHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	ctx.Request.Header.Set("Accept-Language", "en-US,en;q=0.9")

	err := e.NewBusinessError(e.ErrorSellerNotApproved)
	resp := ErrorResponse(ctx, err)

	if resp.Status != e.ErrorSellerNotApproved {
		t.Fatalf("expected status %d, got %d", e.ErrorSellerNotApproved, resp.Status)
	}
	if resp.Msg != e.GetMsgByLocale(e.ErrorSellerNotApproved, "en-US") {
		t.Fatalf("expected localized msg %q, got %q", e.GetMsgByLocale(e.ErrorSellerNotApproved, "en-US"), resp.Msg)
	}
	if resp.MsgKey != e.GetMsgKey(e.ErrorSellerNotApproved) {
		t.Fatalf("expected msg key %q, got %q", e.GetMsgKey(e.ErrorSellerNotApproved), resp.MsgKey)
	}
	if resp.Data != resp.Msg {
		t.Fatalf("expected data to match localized msg %q, got %v", resp.Msg, resp.Data)
	}
	if resp.Error != resp.Msg {
		t.Fatalf("expected error to match localized msg %q, got %q", resp.Msg, resp.Error)
	}
}
