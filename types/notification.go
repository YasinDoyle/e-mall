package types

type NotificationListReq struct {
	BasePage
	UnreadOnly bool `json:"unread_only" form:"unread_only"`
}

type NotificationMarkReadReq struct {
	IDs []uint `json:"ids" form:"ids"`
}

type NotificationUnreadResp struct {
	UnreadCount int64 `json:"unread_count"`
}

type NotificationResp struct {
	ID            uint   `json:"id"`
	RecipientType string `json:"recipient_type"`
	RecipientID   uint   `json:"recipient_id"`
	Scene         string `json:"scene"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Payload       string `json:"payload"`
	Read          bool   `json:"read"`
	CreatedAt     int64  `json:"created_at"`
}
