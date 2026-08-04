package model

import "gorm.io/gorm"

type OrderLogistics struct {
	gorm.Model
	OrderID     uint   `gorm:"not null;index"`
	OrderNum    uint64 `gorm:"not null;index"`
	NodeType    string `gorm:"size:32;not null;index"`
	Description string `gorm:"size:255;not null"`
	OccurredAt  int64  `gorm:"not null;index"`
}
