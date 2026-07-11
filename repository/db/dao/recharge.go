package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type RechargeDao struct {
	*gorm.DB
}

var (
	ErrRechargeNotFound       = errors.New("recharge order not found")
	ErrRechargeStatusConflict = errors.New("recharge order status conflict")
)

func NewRechargeDao(ctx context.Context) *RechargeDao {
	return &RechargeDao{NewDBClient(ctx)}
}

func NewRechargeDaoByDB(db *gorm.DB) *RechargeDao {
	return &RechargeDao{db}
}

func (d *RechargeDao) Create(order *model.RechargeOrder) error {
	return d.DB.Create(order).Error
}

func (d *RechargeDao) GetByOrderNum(orderNum string) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	err := d.DB.Where("order_num = ?", orderNum).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRechargeNotFound
	}
	return &order, err
}

func (d *RechargeDao) MarkFailed(orderNum string) error {
	return d.DB.Model(&model.RechargeOrder{}).
		Where("order_num = ? AND status = ?", orderNum, model.RechargeStatusPending).
		Update("status", model.RechargeStatusFailed).Error
}

func (d *RechargeDao) MarkPaid(orderNum, providerTradeNo string, paidAt time.Time) (*model.RechargeOrder, bool, error) {
	var order *model.RechargeOrder
	freshPaid := false
	err := d.DB.Transaction(func(tx *gorm.DB) error {
		var current model.RechargeOrder
		if err := tx.Where("order_num = ?", orderNum).First(&current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRechargeNotFound
			}
			return err
		}

		order = &current
		if current.Status != model.RechargeStatusPending {
			return nil
		}

		result := tx.Model(&model.RechargeOrder{}).
			Where("id = ? AND status = ?", current.ID, model.RechargeStatusPending).
			Updates(map[string]interface{}{
				"status":            model.RechargeStatusPaid,
				"provider_trade_no": providerTradeNo,
				"paid_at":           paidAt,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrRechargeStatusConflict
		}

		current.Status = model.RechargeStatusPaid
		current.ProviderTradeNo = providerTradeNo
		current.PaidAt = &paidAt
		order = &current
		freshPaid = true
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return order, freshPaid, nil
}

func (d *RechargeDao) MarkUserPaidCredited(userID uint) error {
	return d.DB.Model(&model.RechargeOrder{}).
		Where("user_id = ? AND status = ? AND refund_status IN ?", userID, model.RechargeStatusPaid, []string{model.RechargeRefundNone, model.RechargeRefundFailed}).
		Update("status", model.RechargeStatusCredited).Error
}

func (d *RechargeDao) ListPendingCreditOrdersByUser(userID uint) ([]model.RechargeOrder, error) {
	var orders []model.RechargeOrder
	err := d.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ? AND status = ? AND refund_status IN ?", userID, model.RechargeStatusPaid, []string{model.RechargeRefundNone, model.RechargeRefundFailed}).
		Find(&orders).Error
	return orders, err
}

func (d *RechargeDao) MarkOrdersCredited(userID uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	result := d.DB.Model(&model.RechargeOrder{}).
		Where("user_id = ? AND id IN ? AND status = ? AND refund_status IN ?", userID, ids, model.RechargeStatusPaid, []string{model.RechargeRefundNone, model.RechargeRefundFailed}).
		Update("status", model.RechargeStatusCredited)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != int64(len(ids)) {
		return ErrRechargeStatusConflict
	}
	return nil
}

func (d *RechargeDao) SumPendingCreditByUser(userID uint) (float64, error) {
	var total float64
	err := d.DB.Model(&model.RechargeOrder{}).
		Where("user_id = ? AND status = ? AND refund_status IN ?", userID, model.RechargeStatusPaid, []string{model.RechargeRefundNone, model.RechargeRefundFailed}).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (d *RechargeDao) MarkRefundProcessing(orderNum, refundNo, reason string, amount float64) error {
	result := d.DB.Model(&model.RechargeOrder{}).
		Where("order_num = ? AND status = ? AND refund_status IN ?", orderNum, model.RechargeStatusPaid, []string{model.RechargeRefundNone, model.RechargeRefundFailed}).
		Updates(map[string]interface{}{
			"refund_no":     refundNo,
			"refund_amount": amount,
			"refund_reason": reason,
			"refund_status": model.RechargeRefundProcessing,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRechargeStatusConflict
	}
	return nil
}

func (d *RechargeDao) MarkRefundResult(orderNum, refundNo, refundStatus string, refundedAt *time.Time) error {
	updates := map[string]interface{}{
		"refund_no":     refundNo,
		"refund_status": refundStatus,
	}
	if refundedAt != nil {
		updates["refunded_at"] = *refundedAt
	}

	return d.DB.Model(&model.RechargeOrder{}).
		Where("order_num = ?", orderNum).
		Updates(updates).Error
}
