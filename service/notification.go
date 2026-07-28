package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	notificationhub "github.com/YasinDoyle/e-mall/domain/notification"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var NotificationSrvIns *NotificationSrv
var NotificationSrvOnce sync.Once

type NotificationSrv struct{}

func GetNotificationSrv() *NotificationSrv {
	NotificationSrvOnce.Do(func() { NotificationSrvIns = &NotificationSrv{} })
	return NotificationSrvIns
}

func (s *NotificationSrv) UserList(ctx context.Context, req *types.NotificationListReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.list(ctx, model.NotificationRecipientUser, u.Id, req)
}

func (s *NotificationSrv) AdminList(ctx context.Context, req *types.NotificationListReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.list(ctx, model.NotificationRecipientAdmin, u.Id, req)
}

func (s *NotificationSrv) UserUnreadCount(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.unreadCount(ctx, model.NotificationRecipientUser, u.Id)
}

func (s *NotificationSrv) AdminUnreadCount(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.unreadCount(ctx, model.NotificationRecipientAdmin, u.Id)
}

func (s *NotificationSrv) UserMarkRead(ctx context.Context, req *types.NotificationMarkReadReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.markRead(ctx, model.NotificationRecipientUser, u.Id, req)
}

func (s *NotificationSrv) AdminMarkRead(ctx context.Context, req *types.NotificationMarkReadReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	return s.markRead(ctx, model.NotificationRecipientAdmin, u.Id, req)
}

func (s *NotificationSrv) UserMarkAllRead(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	if err = dao.NewNotificationDao(ctx).MarkAllRead(model.NotificationRecipientUser, u.Id); err != nil {
		return nil, err
	}
	return "操作成功", nil
}

func (s *NotificationSrv) AdminMarkAllRead(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	if err = dao.NewNotificationDao(ctx).MarkAllRead(model.NotificationRecipientAdmin, u.Id); err != nil {
		return nil, err
	}
	return "操作成功", nil
}

func (s *NotificationSrv) StreamUnreadCount(c *gin.Context, recipientType string) {
	u, err := ctl.GetUserInfo(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": e.ErrorAuthCheckTokenFail, "msg": e.GetMsg(e.ErrorAuthCheckTokenFail)})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	subscription := notificationhub.Subscribe(recipientType, u.Id)
	defer subscription.Close()
	sendUnreadCountEvent(c, recipientType, u.Id)
	for {
		select {
		case <-c.Request.Context().Done():
			return
		case <-subscription.C:
			sendUnreadCountEvent(c, recipientType, u.Id)
		case <-ticker.C:
			sendUnreadCountEvent(c, recipientType, u.Id)
		}
	}
}

func (s *NotificationSrv) NotifyUser(ctx context.Context, recipientID uint, scene, title, content string, payload any) {
	if recipientID == 0 {
		return
	}
	notification := buildNotification(model.NotificationRecipientUser, recipientID, scene, title, content, payload)
	if err := dao.NewNotificationDao(ctx).Create(notification); err != nil {
		log.LogrusObj.Errorf("create user notification failed: %v", err)
		return
	}
	notificationhub.Publish(model.NotificationRecipientUser, recipientID)
}

func (s *NotificationSrv) NotifyUserByDB(tx *gorm.DB, recipientID uint, scene, title, content string, payload any) {
	if recipientID == 0 {
		return
	}
	notification := buildNotification(model.NotificationRecipientUser, recipientID, scene, title, content, payload)
	if err := dao.NewNotificationDaoByDB(tx).Create(notification); err != nil {
		log.LogrusObj.Errorf("create user notification failed: %v", err)
		return
	}
}

func (s *NotificationSrv) NotifyAdmins(ctx context.Context, scene, title, content string, payload any) {
	admins, err := dao.NewUserDao(ctx).ListAdminUsers()
	if err != nil {
		log.LogrusObj.Errorf("list admin users for notification failed: %v", err)
		return
	}
	notifications := make([]*model.Notification, 0, len(admins))
	for _, admin := range admins {
		notifications = append(notifications, buildNotification(model.NotificationRecipientAdmin, admin.ID, scene, title, content, payload))
	}
	if err = dao.NewNotificationDao(ctx).BatchCreate(notifications); err != nil {
		log.LogrusObj.Errorf("create admin notifications failed: %v", err)
		return
	}
	for _, admin := range admins {
		notificationhub.Publish(model.NotificationRecipientAdmin, admin.ID)
	}
}

func (s *NotificationSrv) NotifyAdminsByDB(tx *gorm.DB, scene, title, content string, payload any) {
	var admins []*model.User
	if err := tx.Model(&model.User{}).Where("is_admin = ?", true).Find(&admins).Error; err != nil {
		log.LogrusObj.Errorf("list admin users for notification failed: %v", err)
		return
	}
	notifications := make([]*model.Notification, 0, len(admins))
	for _, admin := range admins {
		notifications = append(notifications, buildNotification(model.NotificationRecipientAdmin, admin.ID, scene, title, content, payload))
	}
	if err := dao.NewNotificationDaoByDB(tx).BatchCreate(notifications); err != nil {
		log.LogrusObj.Errorf("create admin notifications failed: %v", err)
		return
	}
}

func (s *NotificationSrv) list(ctx context.Context, recipientType string, recipientID uint, req *types.NotificationListReq) (interface{}, error) {
	normalizeNotificationPage(req)
	notifications, total, err := dao.NewNotificationDao(ctx).List(recipientType, recipientID, req.UnreadOnly, req.PageNum, req.PageSize)
	if err != nil {
		return nil, err
	}
	list := make([]*types.NotificationResp, 0, len(notifications))
	for _, notification := range notifications {
		list = append(list, buildNotificationResp(notification))
	}
	return &types.DataListResp{Item: list, Total: total}, nil
}

func (s *NotificationSrv) unreadCount(ctx context.Context, recipientType string, recipientID uint) (interface{}, error) {
	total, err := dao.NewNotificationDao(ctx).CountUnread(recipientType, recipientID)
	if err != nil {
		return nil, err
	}
	return &types.NotificationUnreadResp{UnreadCount: total}, nil
}

func (s *NotificationSrv) markRead(ctx context.Context, recipientType string, recipientID uint, req *types.NotificationMarkReadReq) (interface{}, error) {
	if err := validateNotificationMarkReadReq(req); err != nil {
		return nil, err
	}
	if err := dao.NewNotificationDao(ctx).MarkRead(recipientType, recipientID, req.IDs); err != nil {
		return nil, err
	}
	return "操作成功", nil
}

func (s *NotificationSrv) markAllRead(ctx context.Context, recipientType string, recipientID uint) (interface{}, error) {
	if err := dao.NewNotificationDao(ctx).MarkAllRead(recipientType, recipientID); err != nil {
		return nil, err
	}
	return "操作成功", nil
}

func normalizeNotificationPage(req *types.NotificationListReq) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
}

func validateNotificationMarkReadReq(req *types.NotificationMarkReadReq) error {
	if req == nil {
		return e.NewBusinessError(e.InvalidParams)
	}
	if len(req.IDs) == 0 {
		return e.NewBusinessError(e.InvalidParams)
	}
	for _, id := range req.IDs {
		if id == 0 {
			return e.NewBusinessError(e.InvalidParams)
		}
	}
	return nil
}

func buildNotification(recipientType string, recipientID uint, scene, title, content string, payload any) *model.Notification {
	return &model.Notification{
		RecipientType: recipientType,
		RecipientID:   recipientID,
		Scene:         strings.TrimSpace(scene),
		Title:         strings.TrimSpace(title),
		Content:       strings.TrimSpace(content),
		Payload:       marshalNotificationPayload(payload),
	}
}

func buildNotificationResp(notification *model.Notification) *types.NotificationResp {
	if notification == nil {
		return &types.NotificationResp{}
	}
	return &types.NotificationResp{
		ID:            notification.ID,
		RecipientType: notification.RecipientType,
		RecipientID:   notification.RecipientID,
		Scene:         notification.Scene,
		Title:         notification.Title,
		Content:       notification.Content,
		Payload:       notification.Payload,
		Read:          notification.Read,
		CreatedAt:     notification.CreatedAt.Unix(),
	}
}

func marshalNotificationPayload(payload any) string {
	if payload == nil {
		return ""
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

func sendUnreadCountEvent(c *gin.Context, recipientType string, recipientID uint) {
	total, err := dao.NewNotificationDao(c.Request.Context()).CountUnread(recipientType, recipientID)
	if err != nil {
		log.LogrusObj.Errorf("stream unread notification count failed: %v", err)
		return
	}
	c.SSEvent("unread", map[string]int64{"unread_count": total})
	c.Writer.Flush()
}

func notificationTitle(scene, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return fallback
	}
	return fmt.Sprintf("业务通知：%s", scene)
}
