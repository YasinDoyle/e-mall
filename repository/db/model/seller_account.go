package model

import "gorm.io/gorm"

type SellerAccount struct {
	gorm.Model
	SellerID         uint          `gorm:"not null;uniqueIndex" json:"seller_id"`
	Seller           SellerProfile `gorm:"foreignKey:SellerID;references:UserID" json:"seller"`
	AvailableBalance float64       `gorm:"not null;default:0" json:"available_balance"`
	FrozenBalance    float64       `gorm:"not null;default:0" json:"frozen_balance"`
	TotalIncome      float64       `gorm:"not null;default:0" json:"total_income"`
	TotalWithdrawn   float64       `gorm:"not null;default:0" json:"total_withdrawn"`
}
