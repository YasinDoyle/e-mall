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

func NewFlashSaleDaoByDB(db *gorm.DB) *FlashSaleDao {
	return &FlashSaleDao{db}
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

func (dao *FlashSaleDao) GetByProductID(productID uint) (resp *model.FlashSale, err error) {
	err = dao.Model(&model.FlashSale{}).
		Where("product_id = ? AND num > 0", productID).
		First(&resp).Error

	return
}

func (dao *FlashSaleDao) CreateAsyncOrder(record *model.FlashSale2MQ) error {
	return dao.Model(&model.FlashSale2MQ{}).Create(&record).Error
}

func (dao *FlashSaleDao) HasAsyncOrder(flashSaleID, userID uint) (bool, error) {
	var count int64
	err := dao.Model(&model.FlashSale2MQ{}).
		Where("flash_sale_id = ? AND user_id = ?", flashSaleID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
