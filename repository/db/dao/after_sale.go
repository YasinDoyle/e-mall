package dao

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type AfterSaleDao struct {
	*gorm.DB
}

func NewAfterSaleDao(ctx context.Context) *AfterSaleDao {
	return &AfterSaleDao{NewDBClient(ctx)}
}

func NewAfterSaleDaoByDB(db *gorm.DB) *AfterSaleDao {
	return &AfterSaleDao{db}
}

func (dao *AfterSaleDao) Create(afterSale *model.AfterSale) error {
	return dao.DB.Create(afterSale).Error
}

func (dao *AfterSaleDao) GetByID(id uint) (*model.AfterSale, error) {
	var afterSale model.AfterSale
	err := dao.DB.Where("id = ?", id).First(&afterSale).Error
	return &afterSale, err
}

func (dao *AfterSaleDao) GetByIDForUpdate(id uint) (*model.AfterSale, error) {
	var afterSale model.AfterSale
	err := dao.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&afterSale).Error
	return &afterSale, err
}

func (dao *AfterSaleDao) GetByOrderIDForUpdate(orderID uint) (*model.AfterSale, error) {
	var afterSale model.AfterSale
	err := dao.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_id = ?", orderID).
		Order("id DESC").
		First(&afterSale).Error
	return &afterSale, err
}

func (dao *AfterSaleDao) TransitionStatus(id uint, fromStatus, toStatus string, updates map[string]interface{}) error {
	if updates == nil {
		updates = map[string]interface{}{}
	}
	updates["status"] = toStatus
	result := dao.DB.Model(&model.AfterSale{}).
		Where("id = ? AND status = ?", id, fromStatus).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (dao *AfterSaleDao) UpdateByID(id uint, updates map[string]interface{}) error {
	return dao.DB.Model(&model.AfterSale{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (dao *AfterSaleDao) HasActiveByOrderID(orderID uint) (bool, error) {
	var afterSale model.AfterSale
	err := dao.DB.Select("id").
		Where("order_id = ?", orderID).
		Order("id DESC").
		Take(&afterSale).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (dao *AfterSaleDao) ListByBuyer(buyerID uint, req *types.AfterSaleListReq) ([]*model.AfterSale, int64, error) {
	return dao.list(func(db *gorm.DB) *gorm.DB {
		return db.Where("buyer_id = ?", buyerID)
	}, req)
}

func (dao *AfterSaleDao) ListBySeller(sellerID uint, req *types.AfterSaleListReq) ([]*model.AfterSale, int64, error) {
	return dao.list(func(db *gorm.DB) *gorm.DB {
		return db.Where("seller_id = ?", sellerID)
	}, req)
}

func (dao *AfterSaleDao) ListAdmin(req *types.AfterSaleListReq) ([]*model.AfterSale, int64, error) {
	return dao.list(nil, req)
}

func (dao *AfterSaleDao) list(baseFilter func(*gorm.DB) *gorm.DB, req *types.AfterSaleListReq) ([]*model.AfterSale, int64, error) {
	query := buildAfterSaleListQuery(dao.DB, req)
	if baseFilter != nil {
		query = baseFilter(query)
	}
	countQuery := query.Session(&gorm.Session{})
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	list := make([]*model.AfterSale, 0)
	if err := query.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func buildAfterSaleListQuery(db *gorm.DB, req *types.AfterSaleListReq) *gorm.DB {
	query := db.Model(&model.AfterSale{})
	if req.OrderID > 0 {
		query = query.Where("order_id = ?", req.OrderID)
	}
	if req.Status != "" {
		query = query.Where("status = ?", strings.TrimSpace(req.Status))
	}
	if req.Type != "" {
		query = query.Where("type = ?", strings.TrimSpace(req.Type))
	}
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	return query.Order("created_at DESC")
}
