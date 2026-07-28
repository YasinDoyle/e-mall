package dao

import (
	"context"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/gorm"
)

type NotificationDao struct {
	*gorm.DB
}

func NewNotificationDao(ctx context.Context) *NotificationDao {
	return &NotificationDao{NewDBClient(ctx)}
}

func NewNotificationDaoByDB(db *gorm.DB) *NotificationDao {
	return &NotificationDao{db}
}

func (dao *NotificationDao) Create(notification *model.Notification) error {
	return dao.DB.Create(notification).Error
}

func (dao *NotificationDao) BatchCreate(notifications []*model.Notification) error {
	if len(notifications) == 0 {
		return nil
	}
	return dao.DB.Create(&notifications).Error
}

func (dao *NotificationDao) List(recipientType string, recipientID uint, unreadOnly bool, page, size int) (notifications []*model.Notification, total int64, err error) {
	countQuery := buildNotificationListQuery(dao.DB, recipientType, recipientID, unreadOnly, 0, 0)
	if err = countQuery.Count(&total).Error; err != nil {
		return
	}

	err = buildNotificationListQuery(dao.DB, recipientType, recipientID, unreadOnly, page, size).
		Find(&notifications).Error
	return
}

func (dao *NotificationDao) CountUnread(recipientType string, recipientID uint) (int64, error) {
	var total int64
	err := buildNotificationUnreadCountQuery(dao.DB, recipientType, recipientID).Count(&total).Error
	return total, err
}

func (dao *NotificationDao) MarkRead(recipientType string, recipientID uint, ids []uint) error {
	query := buildNotificationMarkReadQuery(dao.DB, recipientType, recipientID, ids)
	return query.Update("read", true).Error
}

func (dao *NotificationDao) MarkAllRead(recipientType string, recipientID uint) error {
	return dao.DB.Model(&model.Notification{}).
		Where("recipient_type = ? AND recipient_id = ? AND `read` = ?", recipientType, recipientID, false).
		Update("read", true).Error
}

func buildNotificationListQuery(db *gorm.DB, recipientType string, recipientID uint, unreadOnly bool, page, size int) *gorm.DB {
	query := db.Model(&model.Notification{}).
		Where("recipient_type = ? AND recipient_id = ?", recipientType, recipientID)
	if unreadOnly {
		query = query.Where("`read` = ?", false)
	}
	query = query.Order("created_at DESC")
	if page > 0 && size > 0 {
		query = query.Offset((page - 1) * size).Limit(size)
	}
	return query
}

func buildNotificationUnreadCountQuery(db *gorm.DB, recipientType string, recipientID uint) *gorm.DB {
	return db.Model(&model.Notification{}).
		Where("recipient_type = ? AND recipient_id = ? AND `read` = ?", recipientType, recipientID, false)
}

func buildNotificationMarkReadQuery(db *gorm.DB, recipientType string, recipientID uint, ids []uint) *gorm.DB {
	query := db.Model(&model.Notification{}).
		Where("recipient_type = ? AND recipient_id = ?", recipientType, recipientID)
	if len(ids) > 0 {
		query = query.Where("id IN ?", ids)
	}
	return query
}
