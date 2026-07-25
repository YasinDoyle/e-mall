package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
)

type SellerProfile struct {
	gorm.Model
	UserID       uint   `gorm:"not null;uniqueIndex"`
	User         User   `gorm:"foreignKey:UserID"`
	ShopName     string `gorm:"size:255;not null;index"`
	Description  string `gorm:"size:1000"`
	Status       uint   `gorm:"not null;default:0;index"`
	RejectReason string `gorm:"size:500"`
	ApprovedAt   *time.Time
}

func (s *SellerProfile) IsApproved() bool {
	return s != nil && s.Status == consts.SellerStatusApproved
}
