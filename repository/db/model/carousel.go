package model

import "gorm.io/gorm"

type Carousel struct {
	gorm.Model
	ImgPath   string
	ProductID uint `gorm:"not null"`
}
