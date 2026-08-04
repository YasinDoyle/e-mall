package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type OrderPaymentDao struct {
	*gorm.DB
}

var ErrOrderPaymentStatusConflict = errors.New("order payment status conflict")

func NewOrderPaymentDao(ctx context.Context) *OrderPaymentDao {
	return &OrderPaymentDao{NewDBClient(ctx)}
}

func NewOrderPaymentDaoByDB(db *gorm.DB) *OrderPaymentDao {
	return &OrderPaymentDao{db}
}

func (dao *OrderPaymentDao) Create(payment *model.OrderPayment) error {
	return dao.DB.Create(payment).Error
}

func (dao *OrderPaymentDao) GetByPaymentNo(paymentNo string) (*model.OrderPayment, error) {
	var payment model.OrderPayment
	err := dao.DB.Where("payment_no = ?", paymentNo).First(&payment).Error
	return &payment, err
}

func (dao *OrderPaymentDao) GetByPaymentNoForUpdate(paymentNo string) (*model.OrderPayment, error) {
	var payment model.OrderPayment
	err := buildOrderPaymentByPaymentNoForUpdateQuery(dao.DB, paymentNo).First(&payment).Error
	return &payment, err
}

func (dao *OrderPaymentDao) GetPendingByOrderID(orderID uint) (*model.OrderPayment, error) {
	var payment model.OrderPayment
	err := buildPendingOrderPaymentByOrderIDQuery(dao.DB, orderID).Take(&payment).Error
	return &payment, err
}

func (dao *OrderPaymentDao) MarkPaid(paymentNo, providerTradeNo string, paidAt time.Time) (fresh bool, err error) {
	return markOrderPaymentPaid(dao.DB, paymentNo, providerTradeNo, paidAt)
}

func (dao *OrderPaymentDao) MarkFailed(paymentNo string) error {
	return dao.DB.Model(&model.OrderPayment{}).
		Where("payment_no = ? AND status = ?", paymentNo, consts.OrderPaymentStatusPending).
		Update("status", consts.OrderPaymentStatusFailed).Error
}

func (dao *OrderPaymentDao) ClosePendingByOrderID(orderID uint) error {
	now := time.Now()
	return dao.DB.Model(&model.OrderPayment{}).
		Where("order_id = ? AND status = ?", orderID, consts.OrderPaymentStatusPending).
		Updates(map[string]interface{}{
			"status":    consts.OrderPaymentStatusClosed,
			"closed_at": now,
		}).Error
}

func (dao *OrderPaymentDao) CloseOtherPendingByOrderID(orderID uint, keepPaymentNo string) error {
	now := time.Now()
	return buildCloseOtherPendingOrderPaymentsQuery(dao.DB, orderID, keepPaymentNo).
		Updates(map[string]interface{}{
			"status":    consts.OrderPaymentStatusClosed,
			"closed_at": now,
		}).Error
}

func buildPendingOrderPaymentByOrderIDQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.OrderPayment{}).
		Where("order_id = ?", orderID).
		Where("status = ?", consts.OrderPaymentStatusPending)
}

func buildOrderPaymentByPaymentNoForUpdateQuery(db *gorm.DB, paymentNo string) *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Model(&model.OrderPayment{}).
		Where("payment_no = ?", paymentNo)
}

func buildCloseOtherPendingOrderPaymentsQuery(db *gorm.DB, orderID uint, keepPaymentNo string) *gorm.DB {
	return db.Model(&model.OrderPayment{}).
		Where("order_id = ?", orderID).
		Where("status = ?", consts.OrderPaymentStatusPending).
		Where("payment_no <> ?", keepPaymentNo)
}

func markOrderPaymentPaid(db *gorm.DB, paymentNo, providerTradeNo string, paidAt time.Time) (fresh bool, err error) {
	var payment model.OrderPayment
	err = db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("payment_no = ?", paymentNo).
		First(&payment).Error
	if err != nil {
		return false, err
	}
	if payment.Status == consts.OrderPaymentStatusPaid {
		return false, nil
	}
	if payment.Status != consts.OrderPaymentStatusPending {
		return false, ErrOrderPaymentStatusConflict
	}

	result := db.Model(&model.OrderPayment{}).
		Where("id = ? AND status = ?", payment.ID, consts.OrderPaymentStatusPending).
		Updates(map[string]interface{}{
			"status":            consts.OrderPaymentStatusPaid,
			"provider_trade_no": providerTradeNo,
			"paid_at":           paidAt,
		})
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, ErrOrderPaymentStatusConflict
	}
	return true, nil
}
