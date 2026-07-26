package dao

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type SettlementDao struct {
	*gorm.DB
}

func NewSettlementDao(ctx context.Context) *SettlementDao {
	return &SettlementDao{NewDBClient(ctx)}
}

func NewSettlementDaoByDB(db *gorm.DB) *SettlementDao {
	return &SettlementDao{db}
}

func (dao *SettlementDao) ActiveCommissionRate() (float64, error) {
	var config model.CommissionConfig
	err := dao.DB.Where("enabled = ?", true).Order("id DESC").First(&config).Error
	if err != nil {
		return 0, err
	}
	return config.Rate, nil
}

func (dao *SettlementDao) CreatePending(settlement *model.Settlement) error {
	return dao.DB.Create(settlement).Error
}

func (dao *SettlementDao) List(req *types.AdminSettlementListReq) ([]*model.Settlement, int64, error) {
	db := dao.DB.Model(&model.Settlement{})
	if req.SellerID > 0 {
		db = db.Where("seller_id = ?", req.SellerID)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	list := make([]*model.Settlement, 0)
	err := db.Order("created_at DESC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error
	return list, total, err
}

func (dao *SettlementDao) GetByID(id uint) (*model.Settlement, error) {
	var settlement model.Settlement
	err := dao.DB.Where("id = ?", id).First(&settlement).Error
	return &settlement, err
}

func (dao *SettlementDao) GenerateCompletedForSeller(sellerID uint) (int64, error) {
	result := dao.DB.Model(&model.Settlement{}).
		Where("seller_id = ? AND status = ?", sellerID, model.SettlementStatusPending).
		Where("order_id IN (?)",
			dao.DB.Model(&model.Order{}).
				Select("id").
				Where("boss_id = ? AND type = ? AND refund_status = ?", sellerID, consts.OrderTypeReceipt, consts.OrderRefundStatusNone),
		).
		Update("status", model.SettlementStatusGenerated)
	return result.RowsAffected, result.Error
}

func (dao *SettlementDao) GenerateCompletedByID(id uint) (*model.Settlement, error) {
	var settlement model.Settlement
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&settlement).Error; err != nil {
			return err
		}
		if settlement.Status != model.SettlementStatusPending {
			return gorm.ErrRecordNotFound
		}
		var order model.Order
		if err := tx.Where("id = ? AND boss_id = ? AND type = ? AND refund_status = ?",
			settlement.OrderID,
			settlement.SellerID,
			consts.OrderTypeReceipt,
			consts.OrderRefundStatusNone,
		).First(&order).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Settlement{}).
			Where("id = ? AND status = ?", id, model.SettlementStatusPending).
			Update("status", model.SettlementStatusGenerated).Error; err != nil {
			return err
		}
		settlement.Status = model.SettlementStatusGenerated
		return nil
	})
	return &settlement, err
}

func (dao *SettlementDao) GenerateCompletedForOrder(orderID uint) error {
	return dao.DB.Model(&model.Settlement{}).
		Where("order_id = ? AND status = ?", orderID, model.SettlementStatusPending).
		Update("status", model.SettlementStatusGenerated).Error
}

func (dao *SettlementDao) SellerSummary(sellerID uint) (*types.SellerSettlementSummaryResp, error) {
	var rows []struct {
		Status string
		Amount float64
	}
	err := dao.DB.Model(&model.Settlement{}).
		Select("status, COALESCE(SUM(settlement_amount), 0) AS amount").
		Where("seller_id = ?", sellerID).
		Group("status").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	resp := &types.SellerSettlementSummaryResp{}
	for _, row := range rows {
		switch row.Status {
		case model.SettlementStatusPending:
			resp.PendingAmount = row.Amount
		case model.SettlementStatusGenerated:
			resp.GeneratedAmount = row.Amount
		case model.SettlementStatusPaid:
			resp.PaidAmount = row.Amount
			resp.AvailableAmount = row.Amount
		case model.SettlementStatusRefunded:
			resp.RefundedAmount = row.Amount
		}
	}
	return resp, nil
}

func (dao *SettlementDao) MarkPaid(id uint) (*model.Settlement, error) {
	var settlement model.Settlement
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&settlement).Error; err != nil {
			return err
		}
		if settlement.Status != model.SettlementStatusGenerated {
			return gorm.ErrRecordNotFound
		}
		now := time.Now().Unix()
		if err := tx.Model(&model.Settlement{}).
			Where("id = ? AND status = ?", id, model.SettlementStatusGenerated).
			Updates(map[string]interface{}{
				"status":  model.SettlementStatusPaid,
				"paid_at": &now,
			}).Error; err != nil {
			return err
		}
		settlement.Status = model.SettlementStatusPaid
		settlement.PaidAt = &now
		return nil
	})
	return &settlement, err
}

func (dao *SettlementDao) MarkRefundedByOrderID(orderID uint) error {
	return dao.DB.Model(&model.Settlement{}).
		Where("order_id = ? AND status IN ?", orderID, []string{
			model.SettlementStatusPending,
			model.SettlementStatusGenerated,
		}).
		Update("status", model.SettlementStatusRefunded).Error
}
