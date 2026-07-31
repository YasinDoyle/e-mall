package dao

import (
	"context"
	"errors"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/gorm"
)

type OrderLogDao struct {
	*gorm.DB
}

func NewOrderLogDao(ctx context.Context) *OrderLogDao {
	return &OrderLogDao{NewDBClient(ctx)}
}

func NewOrderLogDaoByDB(db *gorm.DB) *OrderLogDao {
	return &OrderLogDao{db}
}

func (dao *OrderLogDao) Create(log *model.OrderLog) error {
	return dao.DB.Create(log).Error
}

func (dao *OrderLogDao) ListByOrderID(orderID uint) ([]*model.OrderLog, error) {
	var logs []*model.OrderLog
	err := buildOrderLogListByOrderIDQuery(dao.DB, orderID).Find(&logs).Error
	return logs, err
}

func (dao *OrderLogDao) HasLogsForOrder(orderID uint) (bool, error) {
	var log model.OrderLog
	err := buildOrderLogHasLogsForOrderQuery(dao.DB, orderID).Take(&log).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func buildOrderLogListByOrderIDQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.OrderLog{}).
		Where("order_id = ?", orderID).
		Order("created_at ASC").
		Order("id ASC")
}

func buildOrderLogHasLogsForOrderQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.OrderLog{}).
		Select("id").
		Where("order_id = ?", orderID).
		Limit(1)
}
