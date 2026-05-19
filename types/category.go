package types

type ListCategoryReq struct {
}

type ListCategoryResp struct {
	ID           uint   `json:"id"`
	CategoryName string `json:"category_name"`
	CreatedAt    int64  `json:"created_at"`
}
