package dao

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type AdminDao struct {
	*gorm.DB
}

type AdminOrderTrendRow struct {
	Date        string
	OrderCount  int64
	SalesAmount float64
}

func NewAdminDao(ctx context.Context) *AdminDao {
	return &AdminDao{NewDBClient(ctx)}
}

// ===== 分类 =====

func (d *AdminDao) CreateCategory(name string) error {
	return d.DB.Create(&model.Category{CategoryName: name}).Error
}

func (d *AdminDao) UpdateCategory(id uint, name string) error {
	return d.DB.Model(&model.Category{}).Where("id = ?", id).Update("category_name", name).Error
}

func (d *AdminDao) DeleteCategory(id uint) error {
	return d.DB.Delete(&model.Category{}, id).Error
}

// ===== 轮播图 =====

func (d *AdminDao) CreateCarousel(imgPath string, productID uint) error {
	return d.DB.Create(&model.Carousel{ImgPath: imgPath, ProductID: productID}).Error
}

func (d *AdminDao) UpdateCarousel(id uint, imgPath string, productID uint) error {
	return d.DB.Model(&model.Carousel{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"img_path":   imgPath,
			"product_id": productID,
		}).Error
}

func (d *AdminDao) DeleteCarousel(id uint) error {
	return d.DB.Delete(&model.Carousel{}, id).Error
}

// ===== 公告 =====

func (d *AdminDao) ListNotice() (r []*model.Notice, err error) {
	err = d.DB.Find(&r).Error
	return
}

func (d *AdminDao) CreateNotice(text string) error {
	return d.DB.Create(&model.Notice{Text: text}).Error
}

func (d *AdminDao) UpdateNotice(id uint, text string) error {
	return d.DB.Model(&model.Notice{}).Where("id = ?", id).Update("text", text).Error
}

func (d *AdminDao) DeleteNotice(id uint) error {
	return d.DB.Delete(&model.Notice{}, id).Error
}

// ===== 用户管理 =====

func (d *AdminDao) ListUsers(page, size int) (users []*model.User, total int64, err error) {
	db := d.DB.Model(&model.User{})
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Offset((page - 1) * size).Limit(size).Find(&users).Error
	return
}

func (d *AdminDao) SetUserBan(id uint, banned bool) error {
	status := model.Active
	if banned {
		status = "banned"
	}
	return d.DB.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

// ===== 商品管理 =====

func (d *AdminDao) ListProductsAdmin(page, size int, auditStatus *uint) (products []*model.Product, total int64, err error) {
	db := d.DB.Model(&model.Product{})
	if auditStatus != nil {
		db = db.Where("audit_status = ?", *auditStatus)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Offset((page - 1) * size).Limit(size).Find(&products).Error
	return
}

func (d *AdminDao) AuditProduct(id, auditStatus uint) error {
	onSale := auditStatus == consts.ProductAuditApproved
	return d.DB.Model(&model.Product{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"audit_status": auditStatus,
			"on_sale":      onSale,
		}).Error
}

func (d *AdminDao) DeleteProduct(id uint) error {
	return d.DB.Delete(&model.Product{}, id).Error
}

// ===== 统计 =====

func (d *AdminDao) CountTodayOrders(start, end time.Time) (int64, error) {
	var total int64
	err := d.DB.Model(&model.Order{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&total).Error
	return total, err
}

func (d *AdminDao) SumTotalSales() (float64, error) {
	var total float64
	err := d.DB.Model(&model.Order{}).
		Where("type IN ?", []uint{
			consts.OrderTypePendingShipping,
			consts.OrderTypeShipping,
			consts.OrderTypeReceipt,
		}).
		Select("COALESCE(SUM(money * num), 0)").
		Scan(&total).Error
	return total, err
}

func (d *AdminDao) CountRegisteredUsers() (int64, error) {
	var total int64
	err := d.DB.Model(&model.User{}).Count(&total).Error
	return total, err
}

func (d *AdminDao) OrderTrend(start, end time.Time) ([]AdminOrderTrendRow, error) {
	rows := make([]AdminOrderTrendRow, 0)
	err := d.DB.Model(&model.Order{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Where("type IN ?", []uint{
			consts.OrderTypePendingShipping,
			consts.OrderTypeShipping,
			consts.OrderTypeReceipt,
		}).
		Select("DATE(created_at) AS date, COUNT(*) AS order_count, COALESCE(SUM(money * num), 0) AS sales_amount").
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&rows).Error
	return rows, err
}
