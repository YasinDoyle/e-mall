package dao

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type OrderDao struct {
	*gorm.DB
}

func NewOrderDao(ctx context.Context) *OrderDao {
	return &OrderDao{
		NewDBClient(ctx),
	}
}

func NewOrderDaoByDB(db *gorm.DB) *OrderDao {
	return &OrderDao{db}
}

// CreateOrder 创建订单
func (dao *OrderDao) CreateOrder(order *model.Order) error {
	return dao.DB.Create(&order).Error
}

// ListOrderByCondition 获取订单List

func (dao *OrderDao) ListOrderByCondition(uid uint, req *types.OrderListReq) (r []*types.OrderListResp, count int64, err error) {

	d := dao.DB.Model(&model.Order{}).Where("user_id = ?", uid)
	if req.Type != 0 {
		d = d.Where("type = ?", req.Type)
	}
	d.Count(&count)
	db := dao.DB.Model(&model.Order{}).
		Joins("As o LEFT JOIN product as p ON p.id = o.product_id").
		Joins("LEFT JOIN address as a ON a.id = o.address_id").
		Where("o.user_id = ?", uid)

	db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).Order("created_at DESC").
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address").
		Find(&r)

	return
}

func (dao *OrderDao) GetOrderById(id, uId uint) (r *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ?", id, uId).
		First(&r).Error

	return
}

// ShowOrderById 获取订单详情
func (dao *OrderDao) ShowOrderById(id, uId uint) (r *types.OrderListResp, err error) {
	err = dao.DB.Model(&model.Order{}).
		Joins("AS o LEFT JOIN product AS p ON p.id = o.product_id").
		Joins("LEFT JOIN address AS a ON a.id = o.address_id").
		Where("o.id = ? AND o.user_id = ?", id, uId).
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address").
		Find(&r).Error

	return
}

// DeleteOrderById 获取订单详情
func (dao *OrderDao) DeleteOrderById(id, uId uint) error {
	return dao.DB.Model(&model.Order{}).
		Where("id=? AND user_id = ?", id, uId).
		Delete(&model.Order{}).Error
}

// UpdateOrderById 更新订单详情
func (dao *OrderDao) UpdateOrderById(id, uId uint, order *model.Order) error {
	return dao.DB.Where("id = ? AND user_id = ?", id, uId).
		Updates(order).Error
}

func (dao *OrderDao) UpdateOrderPaidById(id, uId uint, paidAt time.Time) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type = ?", id, uId, 1).
		Updates(map[string]interface{}{
			"type":    2,
			"paid_at": paidAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) UpdateOrderTypeByBoss(id, bossId, fromType, toType uint) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND boss_id = ? AND type = ?", id, bossId, fromType).
		Update("type", toType)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) UpdateOrderTypeByUser(id, uId, fromType, toType uint) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type = ?", id, uId, fromType).
		Update("type", toType)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) DeleteUnpaidOrderByOrderNum(orderNum uint64) error {
	return dao.DB.Where("order_num = ? AND type = ?", orderNum, 1).
		Delete(&model.Order{}).Error
}
