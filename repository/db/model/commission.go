package model

import "gorm.io/gorm"

type CommissionConfig struct {
	gorm.Model
	Name    string  `gorm:"size:64;uniqueIndex;not null"`
	Rate    float64 `gorm:"not null"`
	Enabled bool    `gorm:"default:true;index"`
}
