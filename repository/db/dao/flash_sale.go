package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type FlashSaleDao struct {
	*gorm.DB
}

func NewFlashSaleDao(ctx context.Context) *FlashSaleDao {
	return &FlashSaleDao{NewDBClient(ctx)}
}

func (dao *FlashSaleDao) Create(in *model.FlashSale) error {
	return dao.Model(&model.FlashSale{}).Create(&in).Error
}

func (dao *FlashSaleDao) BatchCreate(in []*model.FlashSale) error {
	return dao.Model(&model.FlashSale{}).
		CreateInBatches(&in, consts.ProductBatchCreate).Error
}

func (dao *FlashSaleDao) CreateByList(in []*model.FlashSale) error {
	return dao.Model(&model.FlashSale{}).Create(&in).Error
}

func (dao *FlashSaleDao) ListFlashSales() (resp []*model.FlashSale, err error) {
	err = dao.Model(&model.FlashSale{}).
		Where("num > 0").Find(&resp).Error

	return
}
