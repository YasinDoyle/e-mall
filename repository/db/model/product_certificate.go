package model

import "gorm.io/gorm"

type ProductCertificate struct {
	gorm.Model
	ProductID       uint   `gorm:"not null;index" json:"product_id"`
	CertificateType string `gorm:"size:80;index" json:"certificate_type"`
	Name            string `gorm:"size:120" json:"name"`
	FilePath        string `gorm:"size:255" json:"file_path"`
}
