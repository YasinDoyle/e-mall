package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{NewDBClient(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

// GetProductById 通过 id 获取product
func (dao *ProductDao) GetProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("id=?", id).First(&product).Error

	return
}

// ShowProductById 通过 id 获取product
func (dao *ProductDao) ShowProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("id=?", id).First(&product).Error
	return
}

// ListProductByCondition 获取商品列表
func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page types.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where(condition).
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).
		Find(&products).Error

	return
}

// CreateProduct 创建商品
func (dao *ProductDao) CreateProduct(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).Create(&product).Error
}

// CountProductByCondition 根据情况获取商品的数量
func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where(condition).Count(&total).Error

	return
}

// DeleteProduct 删除商品
func (dao *ProductDao) DeleteProduct(pId, uId uint) error {
	return dao.DB.Model(&model.Product{}).
		Where("id=? AND boss_id=?", pId, uId).Delete(&model.Product{}).Error
}

// UpdateProduct 更新商品
func (dao *ProductDao) UpdateProduct(pId uint, product *model.Product) error {
	return dao.DB.Model(&model.Product{}).
		Where("id=?", pId).Updates(&product).Error
}

func (dao *ProductDao) DecreaseStock(pId uint, count int) error {
	result := dao.DB.Model(&model.Product{}).
		Where("id = ? AND num >= ?", pId, count).
		UpdateColumn("num", gorm.Expr("num - ?", count))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// SearchProduct 搜索商品
func (dao *ProductDao) SearchProduct(info string, page types.BasePage) (products []*model.Product, count int64, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).
		Find(&products).Error

	if err != nil {
		return
	}

	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Count(&count).
		Error

	return
}
