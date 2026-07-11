package model

type FlashSale struct {
	Id         uint    `gorm:"primarykey" json:"id"`
	ProductId  uint    `gorm:"not null" json:"product_id"`
	BossId     uint    `gorm:"not null" json:"boss_id"`
	Title      string  `json:"title"`
	Money      float64 `json:"money"`
	Num        int     `gorm:"not null" json:"num"`
	CustomId   uint    `json:"custom_id"`
	CustomName string  `json:"custom_name"`
}

func (FlashSale) TableName() string {
	return "skill_products"
}

type FlashSale2MQ struct {
	FlashSaleId uint    `json:"flash_sale_id"`
	ProductId   uint    `json:"product_id"`
	BossId      uint    `json:"boss_id"`
	UserId      uint    `json:"user_id"`
	Money       float64 `json:"money"`
	AddressId   uint    `json:"address_id"`
	Key         string  `json:"key"`
}

func (FlashSale2MQ) TableName() string {
	return "skill_product2_mqs"
}
