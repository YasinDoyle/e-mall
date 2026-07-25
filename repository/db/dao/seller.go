package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type SellerDao struct {
	*gorm.DB
}

func NewSellerDao(ctx context.Context) *SellerDao {
	return &SellerDao{NewDBClient(ctx)}
}

func NewSellerDaoByDB(db *gorm.DB) *SellerDao {
	return &SellerDao{db}
}

func (dao *SellerDao) CreateSellerProfile(profile *model.SellerProfile) error {
	return dao.DB.Create(profile).Error
}

func (dao *SellerDao) GetSellerProfileByUserID(userID uint) (*model.SellerProfile, error) {
	var profile model.SellerProfile
	err := dao.DB.Preload("User").
		Where("user_id = ?", userID).
		First(&profile).Error
	return &profile, err
}

func (dao *SellerDao) GetSellerProfileByID(id uint) (*model.SellerProfile, error) {
	var profile model.SellerProfile
	err := dao.DB.Preload("User").
		Where("id = ?", id).
		First(&profile).Error
	return &profile, err
}

func (dao *SellerDao) UpdateSellerApplication(userID uint, shopName, description string) error {
	return dao.DB.Model(&model.SellerProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"shop_name":     shopName,
			"description":   description,
			"status":        consts.SellerStatusPending,
			"reject_reason": "",
			"approved_at":   nil,
		}).Error
}

func (dao *SellerDao) ListSellerProfiles(page, size int, status *uint) (profiles []*model.SellerProfile, total int64, err error) {
	countDB := dao.DB.Model(&model.SellerProfile{})
	if status != nil {
		countDB = countDB.Where("status = ?", *status)
	}
	if err = countDB.Count(&total).Error; err != nil {
		return
	}
	listDB := dao.DB.Model(&model.SellerProfile{}).Preload("User")
	if status != nil {
		listDB = listDB.Where("status = ?", *status)
	}
	err = listDB.Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&profiles).Error
	return
}

func (dao *SellerDao) AuditSellerProfile(id uint, status uint, rejectReason string) error {
	updates := map[string]interface{}{
		"status":        status,
		"reject_reason": rejectReason,
		"approved_at":   nil,
	}
	if status == consts.SellerStatusApproved {
		now := time.Now()
		updates["reject_reason"] = ""
		updates["approved_at"] = &now
	}
	result := dao.DB.Model(&model.SellerProfile{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商家资料不存在")
	}
	return nil
}
