package dao

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type SellerWithdrawDao struct {
	*gorm.DB
}

func NewSellerWithdrawDao(ctx context.Context) *SellerWithdrawDao {
	return &SellerWithdrawDao{NewDBClient(ctx)}
}

func NewSellerWithdrawDaoByDB(db *gorm.DB) *SellerWithdrawDao {
	return &SellerWithdrawDao{db}
}

func buildListSellerWithdrawQuery(db *gorm.DB, req *types.SellerWithdrawListReq) *gorm.DB {
	query := db.Table("seller_withdraw AS sw").
		Select("sw.*").
		Joins("LEFT JOIN seller_profile AS sp ON sw.seller_id = sp.user_id").
		Joins("LEFT JOIN user AS u ON sp.user_id = u.id").
		Where("sw.deleted_at IS NULL")
	if req != nil {
		if req.SellerID > 0 {
			query = query.Where("sw.seller_id = ?", req.SellerID)
		}
		if req.Status != "" {
			query = query.Where("sw.status = ?", req.Status)
		}
		if req.PageNum <= 0 {
			req.PageNum = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = consts.BasePageSize
		}
	}
	return query
}

func (dao *SellerWithdrawDao) Create(withdraw *model.SellerWithdraw) error {
	return dao.DB.Create(withdraw).Error
}

func (dao *SellerWithdrawDao) GetByID(id uint) (*model.SellerWithdraw, error) {
	var withdraw model.SellerWithdraw
	err := dao.DB.Preload("Seller").
		Preload("Seller.User").
		Where("id = ?", id).
		First(&withdraw).Error
	return &withdraw, err
}

func (dao *SellerWithdrawDao) List(req *types.SellerWithdrawListReq) ([]*model.SellerWithdraw, int64, error) {
	query := dao.DB.Model(&model.SellerWithdraw{}).
		Preload("Seller").
		Preload("Seller.User")
	if req != nil {
		if req.SellerID > 0 {
			query = query.Where("seller_id = ?", req.SellerID)
		}
		if req.Status != "" {
			query = query.Where("status = ?", req.Status)
		}
		if req.PageNum <= 0 {
			req.PageNum = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = consts.BasePageSize
		}
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	list := make([]*model.SellerWithdraw, 0)
	if err := query.Order("created_at DESC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (dao *SellerWithdrawDao) ListByRelatedTypeAndID(relatedType string, relatedID uint) ([]*model.AccountFlow, error) {
	flows := make([]*model.AccountFlow, 0)
	err := dao.DB.Where("related_type = ? AND related_id = ?", relatedType, relatedID).
		Order("created_at ASC").
		Find(&flows).Error
	return flows, err
}

func (dao *SellerWithdrawDao) UpdateStatus(id uint, fromStatus, toStatus string, updates map[string]interface{}) (*model.SellerWithdraw, error) {
	var withdraw model.SellerWithdraw
	err := dao.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).First(&withdraw).Error; err != nil {
			return err
		}
		if withdraw.Status != fromStatus {
			return gorm.ErrRecordNotFound
		}
		updates["status"] = toStatus
		if err := tx.Model(&model.SellerWithdraw{}).
			Where("id = ? AND status = ?", id, fromStatus).
			Updates(updates).Error; err != nil {
			return err
		}
		if v, ok := updates["audited_at"]; ok {
			if t, ok := v.(*time.Time); ok {
				withdraw.AuditedAt = t
			}
		}
		if v, ok := updates["paid_at"]; ok {
			if t, ok := v.(*time.Time); ok {
				withdraw.PaidAt = t
			}
		}
		withdraw.Status = toStatus
		if reason, ok := updates["audit_reason"].(string); ok {
			withdraw.AuditReason = reason
		}
		return nil
	})
	return &withdraw, err
}
