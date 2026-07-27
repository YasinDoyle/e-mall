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

func (dao *AccountFlowDao) ListByRelatedTypeAndID(relatedType string, relatedID uint) ([]*model.AccountFlow, error) {
	flows := make([]*model.AccountFlow, 0)
	err := dao.DB.Where("related_type = ? AND related_id = ?", relatedType, relatedID).
		Order("created_at ASC").
		Find(&flows).Error
	return flows, err
}

func (dao *AccountFlowDao) ListRelatedIDsByTypeAndFlowType(relatedType, flowType string) ([]uint, error) {
	relatedIDs := make([]uint, 0)
	err := dao.DB.Model(&model.AccountFlow{}).
		Distinct().
		Where("related_type = ? AND flow_type = ?", relatedType, flowType).
		Pluck("related_id", &relatedIDs).Error
	return relatedIDs, err
}
