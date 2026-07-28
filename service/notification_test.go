package service

import (
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestNormalizeNotificationPageDefaults(t *testing.T) {
	req := &types.NotificationListReq{}
	normalizeNotificationPage(req)
	if req.PageNum != 1 {
		t.Fatalf("expected default page num 1, got %d", req.PageNum)
	}
	if req.PageSize != 10 {
		t.Fatalf("expected default page size 10, got %d", req.PageSize)
	}
}

func TestValidateNotificationMarkReadReqRejectsZeroID(t *testing.T) {
	err := validateNotificationMarkReadReq(&types.NotificationMarkReadReq{IDs: []uint{1, 0}})
	if err == nil {
		t.Fatal("expected zero id to fail")
	}
	assertServiceBusinessCode(t, err, e.InvalidParams)
}

func TestValidateNotificationMarkReadReqRejectsEmptyIDs(t *testing.T) {
	err := validateNotificationMarkReadReq(&types.NotificationMarkReadReq{})
	if err == nil {
		t.Fatal("expected empty ids to fail")
	}
	assertServiceBusinessCode(t, err, e.InvalidParams)
}

func TestBuildNotificationRespIncludesCreatedAt(t *testing.T) {
	notification := &model.Notification{
		RecipientType: model.NotificationRecipientUser,
		RecipientID:   3,
		Scene:         model.NotificationSceneSellerAudit,
		Title:         "审核通过",
		Content:       "店铺已通过审核",
		Payload:       `{"seller_id":3}`,
	}
	notification.ID = 9

	resp := buildNotificationResp(notification)
	if resp.ID != 9 || resp.Scene != model.NotificationSceneSellerAudit {
		t.Fatalf("unexpected notification resp: %+v", resp)
	}
}
