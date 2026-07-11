package dao

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type CouponDao struct {
	*gorm.DB
}

var (
	ErrCouponAlreadyClaimed = errors.New("coupon already claimed")
	ErrCouponUnavailable    = errors.New("coupon unavailable")
	ErrCouponUseFailed      = errors.New("coupon use failed")
)

func NewCouponDao(ctx context.Context) *CouponDao {
	return &CouponDao{NewDBClient(ctx)}
}

func NewCouponDaoByDB(db *gorm.DB) *CouponDao {
	return &CouponDao{db}
}

// CreateCoupon 管理员创建优惠券
func (d *CouponDao) CreateCoupon(c *model.Coupon) error {
	return d.DB.Create(c).Error
}

// GetCouponByID 获取优惠券信息
func (d *CouponDao) GetCouponByID(id uint) (*model.Coupon, error) {
	var coupon model.Coupon
	err := d.DB.First(&coupon, id).Error
	return &coupon, err
}

// ListCoupons 列出所有未过期且有库存的优惠券（用于展示）
func (d *CouponDao) ListCoupons() (coupons []*model.Coupon, err error) {
	err = d.DB.Where("expire_at > ? AND (stock > 0 OR stock = -1)", time.Now()).Find(&coupons).Error
	return
}

func (d *CouponDao) ListCouponsAdmin() (coupons []*model.Coupon, err error) {
	err = d.DB.Order("created_at DESC").Find(&coupons).Error
	return
}

func (d *CouponDao) OfflineCoupon(id uint) error {
	return d.DB.Model(&model.Coupon{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"stock":     0,
			"expire_at": time.Now(),
		}).Error
}

// DecrStock 扣减库存（stock=-1 时无限制）
func (d *CouponDao) DecrStock(id uint) error {
	result := d.DB.Model(&model.Coupon{}).
		Where("id = ? AND (stock > 0 OR stock = -1) AND expire_at > ?", id, time.Now()).
		UpdateColumn("stock", gorm.Expr("CASE WHEN stock = -1 THEN stock ELSE stock - 1 END"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// ClaimCoupon 用户领券
func (d *CouponDao) ClaimCoupon(userID, couponID uint) error {
	uc := &model.UserCoupon{UserID: userID, CouponID: couponID}
	return d.DB.Create(uc).Error
}

// ClaimCouponWithStock 在一个事务里扣券库存并创建用户领券记录。
func (d *CouponDao) ClaimCouponWithStock(userID, couponID uint) error {
	return d.DB.Transaction(func(tx *gorm.DB) error {
		var count int64
		if err := tx.Model(&model.UserCoupon{}).
			Where("user_id = ? AND coupon_id = ?", userID, couponID).
			Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return ErrCouponAlreadyClaimed
		}

		result := tx.Model(&model.Coupon{}).
			Where("id = ? AND (stock > 0 OR stock = -1) AND expire_at > ?", couponID, time.Now()).
			UpdateColumn("stock", gorm.Expr("CASE WHEN stock = -1 THEN stock ELSE stock - 1 END"))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrCouponUnavailable
		}

		if err := tx.Create(&model.UserCoupon{UserID: userID, CouponID: couponID}).Error; err != nil {
			return err
		}

		return nil
	})
}

// HasClaimed 检查用户是否已领过该券
func (d *CouponDao) HasClaimed(userID, couponID uint) bool {
	var count int64
	d.DB.Model(&model.UserCoupon{}).
		Where("user_id = ? AND coupon_id = ?", userID, couponID).Count(&count)
	return count > 0
}

// ListUserCoupons 查询用户的优惠券列表
func (d *CouponDao) ListUserCoupons(userID uint) (list []*model.UserCoupon, err error) {
	err = d.DB.Where("user_id = ?", userID).Find(&list).Error
	return
}

// GetUserCoupon 获取用户持有的某张未使用优惠券
func (d *CouponDao) GetUserCoupon(userID, couponID uint) (*model.UserCoupon, error) {
	var uc model.UserCoupon
	err := d.DB.Where("user_id = ? AND coupon_id = ? AND status = ?",
		userID, couponID, model.UserCouponUnused).First(&uc).Error
	return &uc, err
}

// UseCoupon 核销优惠券（事务内调用）
func (d *CouponDao) UseCoupon(userCouponID, orderID uint) error {
	result := d.DB.Model(&model.UserCoupon{}).
		Where("id = ? AND status = ?", userCouponID, model.UserCouponUnused).
		Updates(map[string]interface{}{
			"status":   model.UserCouponUsed,
			"order_id": orderID,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrCouponUseFailed
	}

	return nil
}
