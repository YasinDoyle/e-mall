package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type ProductCertificateDao struct {
	*gorm.DB
}

func NewProductCertificateDao(ctx context.Context) *ProductCertificateDao {
	return &ProductCertificateDao{NewDBClient(ctx)}
}

func NewProductCertificateDaoByDB(db *gorm.DB) *ProductCertificateDao {
	return &ProductCertificateDao{db}
}

func (dao *ProductCertificateDao) CreateProductCertificate(certificate *model.ProductCertificate) error {
	return dao.DB.Model(&model.ProductCertificate{}).Create(&certificate).Error
}

func (dao *ProductCertificateDao) ListByProductID(productID uint) (certificates []*model.ProductCertificate, err error) {
	err = dao.DB.Model(&model.ProductCertificate{}).
		Where("product_id = ?", productID).
		Order("id ASC").
		Find(&certificates).Error
	return
}

func (dao *ProductCertificateDao) DeleteByProductID(productID uint) error {
	return dao.DB.Where("product_id = ?", productID).Delete(&model.ProductCertificate{}).Error
}
