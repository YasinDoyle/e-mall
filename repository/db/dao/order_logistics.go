package dao

import (
	"context"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/gorm"
)

type OrderLogisticsDao struct {
	*gorm.DB
}

func NewOrderLogisticsDao(ctx context.Context) *OrderLogisticsDao {
	return &OrderLogisticsDao{NewDBClient(ctx)}
}

func NewOrderLogisticsDaoByDB(db *gorm.DB) *OrderLogisticsDao {
	return &OrderLogisticsDao{db}
}

func (dao *OrderLogisticsDao) Create(node *model.OrderLogistics) error {
	return dao.DB.Create(node).Error
}

func (dao *OrderLogisticsDao) ListByOrderID(orderID uint) ([]*model.OrderLogistics, error) {
	var nodes []*model.OrderLogistics
	err := buildOrderLogisticsListByOrderIDQuery(dao.DB, orderID).Find(&nodes).Error
	return nodes, err
}

func buildOrderLogisticsListByOrderIDQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.OrderLogistics{}).
		Where("order_id = ?", orderID).
		Order("occurred_at ASC").
		Order("id ASC")
}
