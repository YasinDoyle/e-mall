package model

import "gorm.io/gorm"

const (
	NotificationRecipientUser  = "user"
	NotificationRecipientAdmin = "admin"

	NotificationSceneSellerAudit   = "seller_audit"
	NotificationSceneProductAudit  = "product_audit"
	NotificationSceneOrderPaid     = "order_paid"
	NotificationSceneOrderShipped  = "order_shipped"
	NotificationSceneOrderRefunded = "order_refunded"
	NotificationSceneSettlement    = "settlement"
	NotificationSceneWithdraw      = "withdraw"
)

type Notification struct {
	gorm.Model
	RecipientType string `gorm:"size:32;not null;index:idx_notification_recipient_read" json:"recipient_type"`
	RecipientID   uint   `gorm:"not null;index:idx_notification_recipient_read" json:"recipient_id"`
	Scene         string `gorm:"size:64;not null;index" json:"scene"`
	Title         string `gorm:"size:128;not null" json:"title"`
	Content       string `gorm:"size:500;not null" json:"content"`
	Payload       string `gorm:"type:text" json:"payload"`
	Read          bool   `gorm:"not null;default:false;index:idx_notification_recipient_read" json:"read"`
}
