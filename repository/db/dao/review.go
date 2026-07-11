package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type ReviewDao struct {
	*gorm.DB
}

func NewReviewDao(ctx context.Context) *ReviewDao {
	return &ReviewDao{NewDBClient(ctx)}
}

// CreateReview 创建评价
func (d *ReviewDao) CreateReview(r *model.Review) error {
	return d.DB.Create(r).Error
}

// HasReviewed 检查用户是否已对该订单评价过
func (d *ReviewDao) HasReviewed(userID, orderID uint) bool {
	var count int64
	d.DB.Model(&model.Review{}).Where("user_id = ? AND order_id = ?", userID, orderID).Count(&count)
	return count > 0
}

// ListByProduct 查询商品评价列表（分页）
func (d *ReviewDao) ListByProduct(productID uint, page, size int) (reviews []*model.Review, total int64, err error) {
	db := d.DB.Model(&model.Review{}).Where("product_id = ?", productID)
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("created_at DESC").
		Offset((page - 1) * size).
		Limit(size).
		Find(&reviews).Error
	return
}

// DeleteReview 管理员删除评价
func (d *ReviewDao) DeleteReview(id uint) error {
	return d.DB.Delete(&model.Review{}, id).Error
}
