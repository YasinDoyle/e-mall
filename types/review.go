package types

// ReviewCreateReq 创建评价请求
type ReviewCreateReq struct {
	ProductID uint   `json:"product_id" binding:"required"`
	OrderID   uint   `json:"order_id" binding:"required"`
	Rating    uint   `json:"rating" binding:"required,min=1,max=5"`
	Content   string `json:"content"`
	Images    string `json:"images"` // 逗号分隔的图片 URL
}

// ReviewListReq 查询商品评价列表
type ReviewListReq struct {
	ProductID uint `form:"product_id" binding:"required"`
	PageNum   int  `form:"page_num"`
	PageSize  int  `form:"page_size"`
}

// ReviewAdminDeleteReq 管理员删除评价
type ReviewAdminDeleteReq struct {
	ID uint `json:"id" binding:"required"`
}

// ReviewResp 评价响应
type ReviewResp struct {
	ID         uint   `json:"id"`
	UserName   string `json:"user_name"`
	UserAvatar string `json:"user_avatar"`
	Rating     uint   `json:"rating"`
	Content    string `json:"content"`
	Images     string `json:"images"`
	CreatedAt  int64  `json:"created_at"`
}

type ReviewImageUploadResp struct {
	URL string `json:"url"`
}
