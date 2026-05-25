package model

type FlashSale struct {
	Id         uint `gorm:"primarykey"`
	ProductId  uint `gorm:"not null"`
	BossId     uint `gorm:"not null"`
	Title      string
	Money      float64
	Num        int `gorm:"not null"`
	CustomId   uint
	CustomName string
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
