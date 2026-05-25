package types

type FlashSaleImportReq struct {
}

type FlashSaleReq struct {
	FlashSaleId uint   `json:"flash_sale_id" form:"flash_sale_id"`
	ProductId   uint   `json:"product_id" form:"product_id"`
	BossId      uint   `json:"boss_id" form:"boss_id"`
	AddressId   uint   `json:"address_id" form:"address_id"`
	Key         string `json:"key" form:"key"`
}

type ListFlashSaleReq struct {
	PageSize int64 `json:"page_size" form:"page_size"`
	PageNum  int64 `json:"page_num" form:"page_num"`
}

type GetFlashSaleReq struct {
	ProductId uint `json:"product_id" form:"product_id"`
}
