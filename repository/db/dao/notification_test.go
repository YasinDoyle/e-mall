package dao

import (
	"strings"
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func notificationDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open dry run db: %v", err)
	}
	return db
}

func TestBuildNotificationListQueryScopesRecipientAndUnread(t *testing.T) {
	db := notificationDryRunDB(t)
	stmt := buildNotificationListQuery(db, model.NotificationRecipientUser, 7, true, 1, 20).
		Find(&[]model.Notification{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "recipient_type = ?") || !strings.Contains(sql, "recipient_id = ?") {
		t.Fatalf("notification list should scope recipient, got %s", sql)
	}
	if !strings.Contains(sql, "`read` = ?") && !strings.Contains(sql, "read = ?") {
		t.Fatalf("notification list should filter unread records, got %s", sql)
	}
}

func TestBuildNotificationUnreadCountQueryScopesRecipient(t *testing.T) {
	db := notificationDryRunDB(t)
	var count int64
	stmt := buildNotificationUnreadCountQuery(db, model.NotificationRecipientAdmin, 1).
		Count(&count).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "recipient_type = ?") || !strings.Contains(sql, "recipient_id = ?") {
		t.Fatalf("unread count should scope recipient, got %s", sql)
	}
	if !strings.Contains(sql, "`read` = ?") && !strings.Contains(sql, "read = ?") {
		t.Fatalf("unread count should only count unread records, got %s", sql)
	}
}

func TestBuildNotificationMarkReadQueryScopesRecipient(t *testing.T) {
	db := notificationDryRunDB(t)
	stmt := buildNotificationMarkReadQuery(db, model.NotificationRecipientUser, 9, []uint{1, 2}).
		Find(&[]model.Notification{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "recipient_type = ?") || !strings.Contains(sql, "recipient_id = ?") {
		t.Fatalf("mark read should scope recipient, got %s", sql)
	}
	if !strings.Contains(sql, "id IN") {
		t.Fatalf("mark read should target requested ids, got %s", sql)
	}
}
