package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

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

func (dao *SettlementDao) GetByOrderIDForUpdate(orderID uint) (*model.Settlement, error) {
	var settlement model.Settlement
	err := dao.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_id = ?", orderID).
		First(&settlement).Error
	return &settlement, err
}

func (dao *SettlementDao) GenerateCompletedForSeller(sellerID uint) (int64, error) {
	result := dao.DB.Model(&model.Settlement{}).
		Where("seller_id = ? AND status = ?", sellerID, model.SettlementStatusPending).
		Where("order_id IN (?)", buildGenerateCompletedForSellerOrderSubquery(dao.DB, sellerID)).
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
		if err := buildSettlementReadyOrderQuery(tx, settlement.OrderID, settlement.SellerID).
			First(&model.Order{}).Error; err != nil {
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
	return buildGenerateCompletedForOrderQuery(dao.DB, orderID).
		Update("status", model.SettlementStatusGenerated).Error
}

func (dao *SettlementDao) ListPaidSettlements() ([]*model.Settlement, error) {
	settlements := make([]*model.Settlement, 0)
	err := dao.DB.
		Where("status = ?", model.SettlementStatusPaid).
		Order("id ASC").
		Find(&settlements).Error
	return settlements, err
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
	account, accountErr := dao.getSellerAccountSummary(sellerID)
	if accountErr != nil && !errors.Is(accountErr, gorm.ErrRecordNotFound) {
		return nil, accountErr
	}
	resp.AvailableAmount = account.AvailableBalance
	for _, row := range rows {
		switch row.Status {
		case model.SettlementStatusPending:
			resp.PendingAmount = row.Amount
		case model.SettlementStatusGenerated:
			resp.GeneratedAmount = row.Amount
		case model.SettlementStatusPaid:
			resp.PaidAmount = row.Amount
		case model.SettlementStatusRefunded:
			resp.RefundedAmount = row.Amount
		}
	}
	return resp, nil
}

func (dao *SettlementDao) getSellerAccountSummary(sellerID uint) (model.SellerAccount, error) {
	var account model.SellerAccount
	err := buildSellerAccountSummaryQuery(dao.DB, sellerID).First(&account).Error
	return account, err
}

func buildSellerAccountSummaryQuery(db *gorm.DB, sellerID uint) *gorm.DB {
	return db.Table("seller_account").Where("seller_id = ?", sellerID)
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
		if err := buildSettlementMarkPaidReadyOrderQuery(tx, &settlement).
			First(&model.Order{}).Error; err != nil {
			return err
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
	return buildMarkRefundedByOrderIDQuery(dao.DB, orderID).
		Update("status", model.SettlementStatusRefunded).Error
}

func buildGenerateCompletedForOrderQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.Settlement{}).
		Where("order_id = ? AND status = ?", orderID, model.SettlementStatusPending).
		Where("NOT EXISTS (?)", activeAfterSaleSubquery(db, "settlement.order_id"))
}

func buildSettlementReadyOrderQuery(db *gorm.DB, orderID, sellerID uint) *gorm.DB {
	return db.Table("`order` AS o").
		Where("o.id = ? AND o.boss_id = ? AND o.type = ? AND o.refund_status = ?",
			orderID,
			sellerID,
			consts.OrderTypeReceipt,
			consts.OrderRefundStatusNone,
		).
		Where("NOT EXISTS (?)", activeAfterSaleSubquery(db, "o.id"))
}

func buildSettlementMarkPaidReadyOrderQuery(db *gorm.DB, settlement *model.Settlement) *gorm.DB {
	if settlement == nil {
		return buildSettlementReadyOrderQuery(db, 0, 0)
	}
	return buildSettlementReadyOrderQuery(db, settlement.OrderID, settlement.SellerID)
}

func buildGenerateCompletedForSellerOrderSubquery(db *gorm.DB, sellerID uint) *gorm.DB {
	return db.Table("`order` AS o").
		Select("o.id").
		Where("o.boss_id = ? AND o.type = ? AND o.refund_status = ?",
			sellerID,
			consts.OrderTypeReceipt,
			consts.OrderRefundStatusNone,
		).
		Where("NOT EXISTS (?)", activeAfterSaleSubquery(db, "o.id"))
}

func buildMarkRefundedByOrderIDQuery(db *gorm.DB, orderID uint) *gorm.DB {
	return db.Model(&model.Settlement{}).
		Where("order_id = ? AND status IN ?", orderID, []string{
			model.SettlementStatusPending,
			model.SettlementStatusGenerated,
			model.SettlementStatusPaid,
		})
}

func activeAfterSaleSubquery(db *gorm.DB, orderIDColumn string) *gorm.DB {
	return db.Session(&gorm.Session{NewDB: true}).
		Table("after_sale AS af").
		Select("1").
		Where("af.order_id = "+orderIDColumn).
		Where("af.deleted_at IS NULL").
		Where("af.status NOT IN ?", []string{
			consts.AfterSaleStatusRefunded,
			consts.AfterSaleStatusClosed,
		})
}
