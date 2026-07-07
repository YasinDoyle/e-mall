package types

type AdminIDReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

type AdminListReq struct {
	PageNum  int `json:"page_num" form:"page_num"`
	PageSize int `json:"page_size" form:"page_size"`
}

// ===== 分类 =====

type AdminCategoryReq struct {
	CategoryName string `json:"category_name" binding:"required"`
}

type AdminCategoryUpdateReq struct {
	ID           uint   `json:"id" binding:"required"`
	CategoryName string `json:"category_name" binding:"required"`
}

// ===== 轮播图 =====

type AdminCarouselReq struct {
	ImgPath   string `json:"img_path" binding:"required"`
	ProductID uint   `json:"product_id"`
}

// ===== 公告 =====

type AdminNoticeReq struct {
	Text string `json:"text" binding:"required"`
}

type AdminNoticeUpdateReq struct {
	ID   uint   `json:"id" binding:"required"`
	Text string `json:"text" binding:"required"`
}

// ===== 用户管理 =====

type AdminUserBanReq struct {
	ID     uint `json:"id" binding:"required"`
	Banned bool `json:"banned"`
}

// ===== 商品管理 =====

type AdminProductListReq struct {
	PageNum     int   `json:"page_num" form:"page_num"`
	PageSize    int   `json:"page_size" form:"page_size"`
	AuditStatus *uint `json:"audit_status" form:"audit_status"`
}

type AdminProductAuditReq struct {
	ID          uint `json:"id" binding:"required"`
	AuditStatus uint `json:"audit_status"`
}
