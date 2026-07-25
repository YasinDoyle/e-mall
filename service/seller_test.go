package service

import (
	"strings"
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestValidateSellerApplyReq(t *testing.T) {
	req := &types.SellerApplyReq{
		ShopName:    "程记数码",
		Description: "主营手机、电脑和配件",
	}
	if err := validateSellerApplyReq(req); err != nil {
		t.Fatalf("expected valid seller application to pass, got %v", err)
	}
}

func TestValidateSellerApplyReqRejectsBlankShopName(t *testing.T) {
	req := &types.SellerApplyReq{
		ShopName:    "   ",
		Description: "主营手机、电脑和配件",
	}
	err := validateSellerApplyReq(req)
	if err == nil {
		t.Fatal("expected blank shop name to fail")
	}
	assertBusinessCode(t, err, e.ErrorSellerShopNameRequired)
}

func TestValidateAdminSellerAuditReqRequiresRejectReason(t *testing.T) {
	req := &types.AdminSellerAuditReq{
		ID:     1,
		Status: consts.SellerStatusRejected,
	}
	err := validateAdminSellerAuditReq(req)
	if err == nil {
		t.Fatal("expected rejected seller audit without reason to fail")
	}
	if !strings.Contains(err.Error(), "拒绝原因") {
		t.Fatalf("expected error to mention reject reason, got %v", err)
	}
	assertBusinessCode(t, err, e.ErrorSellerRejectReasonMissing)
}

func assertBusinessCode(t *testing.T, err error, want int) {
	t.Helper()
	code, ok := e.CodeFromError(err)
	if !ok {
		t.Fatalf("expected business error code %d, got plain error %v", want, err)
	}
	if code != want {
		t.Fatalf("expected business error code %d, got %d", want, code)
	}
}
