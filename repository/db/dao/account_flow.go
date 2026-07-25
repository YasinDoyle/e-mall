package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type AccountFlowDao struct {
	*gorm.DB
}

func NewAccountFlowDao(ctx context.Context) *AccountFlowDao {
	return &AccountFlowDao{NewDBClient(ctx)}
}

func NewAccountFlowDaoByDB(db *gorm.DB) *AccountFlowDao {
	return &AccountFlowDao{db}
}

func (dao *AccountFlowDao) Create(flow *model.AccountFlow) error {
	return dao.DB.Create(flow).Error
}

func (dao *AccountFlowDao) ListByOrderID(orderID uint) ([]*model.AccountFlow, error) {
	flows := make([]*model.AccountFlow, 0)
	err := dao.DB.Where("order_id = ?", orderID).Order("created_at ASC").Find(&flows).Error
	return flows, err
}
