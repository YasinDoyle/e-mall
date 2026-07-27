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
	ProductID uint   `json:"product_id" binding:"required,gt=0"`
}

type AdminCarouselUpdateReq struct {
	ID        uint   `json:"id" binding:"required"`
	ImgPath   string `json:"img_path" binding:"required"`
	ProductID uint   `json:"product_id" binding:"required,gt=0"`
}

type AdminUploadResp struct {
	URL string `json:"url"`
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

type AdminUserResp struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nick_name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Avatar    string `json:"avatar"`
	IsAdmin   bool   `json:"is_admin"`
	CreatedAt int64  `json:"created_at"`
}

// ===== 商家管理 =====

type AdminSellerListReq struct {
	PageNum  int   `json:"page_num" form:"page_num"`
	PageSize int   `json:"page_size" form:"page_size"`
	Status   *uint `json:"status" form:"status"`
}

type AdminSellerAuditReq struct {
	ID           uint   `json:"id" binding:"required"`
	Status       uint   `json:"status"`
	RejectReason string `json:"reject_reason"`
}

type AdminSellerResp struct {
	ID           uint   `json:"id"`
	UserID       uint   `json:"user_id"`
	UserName     string `json:"user_name"`
	NickName     string `json:"nick_name"`
	Email        string `json:"email"`
	ShopName     string `json:"shop_name"`
	Description  string `json:"description"`
	Status       uint   `json:"status"`
	StatusText   string `json:"status_text"`
	RejectReason string `json:"reject_reason"`
	ApprovedAt   int64  `json:"approved_at"`
	CreatedAt    int64  `json:"created_at"`
}

// ===== 商品管理 =====

type AdminProductListReq struct {
	PageNum     int   `json:"page_num" form:"page_num"`
	PageSize    int   `json:"page_size" form:"page_size"`
	AuditStatus *uint `json:"audit_status" form:"audit_status"`
	Status      *uint `json:"status" form:"status"`
}

type AdminProductAuditReq struct {
	ID          uint `json:"id" binding:"required"`
	AuditStatus uint `json:"audit_status" form:"audit_status"`
	Status      uint `json:"status" form:"status"`
}

type AdminProductResp struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CategoryID    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	OnSale        bool   `json:"on_sale"`
	Num           int    `json:"num"`
	BossID        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	BossAvatar    string `json:"boss_avatar"`
	AuditStatus   uint   `json:"audit_status"`
	Status        uint   `json:"status"`
	CreatedAt     int64  `json:"created_at"`
}

// ===== 统计 =====

type AdminStatsOrdersReq struct {
	StartDate string `json:"start_date" form:"start_date"`
	EndDate   string `json:"end_date" form:"end_date"`
}

type AdminStatsOverviewResp struct {
	TodayOrders     int64   `json:"today_orders"`
	TotalSales      float64 `json:"total_sales"`
	PlatformRevenue float64 `json:"platform_revenue"`
	RegisteredUsers int64   `json:"registered_users"`
}

type AdminStatsOrdersResp struct {
	Dates        []string  `json:"dates"`
	OrderCounts  []int64   `json:"order_counts"`
	SalesAmounts []float64 `json:"sales_amounts"`
}
