package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/YasinDoyle/e-mall/consts"
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
	d := dao.DB.Table("`order` AS o").
		Where("o.user_id = ?", uid).
		Where("o.deleted_at IS NULL").
		Where("o.buyer_deleted = ?", false)
	if req.Type != 0 {
		d = d.Where("o.type = ?", req.Type)
	}
	if err = d.Count(&count).Error; err != nil {
		return
	}

	err = buildListOrderByConditionQuery(dao.DB, uid, req).
		Find(&r).Error

	return
}

func buildListOrderByConditionQuery(db *gorm.DB, uid uint, req *types.OrderListReq) *gorm.DB {
	query := db.Table("`order` AS o").
		Joins("LEFT JOIN product as p ON p.id = o.product_id").
		Joins("LEFT JOIN address as a ON a.id = o.address_id").
		Where("o.user_id = ?", uid).
		Where("o.deleted_at IS NULL").
		Where("o.buyer_deleted = ?", false)
	if req.Type != 0 {
		query = query.Where("o.type = ?", req.Type)
	}

	return query.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).Order("o.created_at DESC").
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"o.money AS money," +
			"o.refund_status AS refund_status," +
			"o.refund_reason AS refund_reason," +
			"o.payment_channel AS payment_channel," +
			"o.logistics_company AS logistics_company," +
			"o.tracking_no AS tracking_no," +
			"IFNULL(UNIX_TIMESTAMP(o.shipped_at), 0) AS shipped_at," +
			"IFNULL(UNIX_TIMESTAMP(o.received_at), 0) AS received_at," +
			"IFNULL(UNIX_TIMESTAMP(o.canceled_at), 0) AS canceled_at," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address")
}

func (dao *OrderDao) ListOrderByBoss(bossID uint, req *types.SellerOrderListReq) (r []*types.OrderListResp, count int64, err error) {
	d := dao.DB.Table("`order` AS o").
		Where("o.boss_id = ?", bossID).
		Where("o.deleted_at IS NULL")
	if req.Type != 0 {
		d = d.Where("o.type = ?", req.Type)
	}
	if err = d.Count(&count).Error; err != nil {
		return
	}

	err = buildListOrderByBossQuery(dao.DB, bossID, req).
		Find(&r).Error

	return
}

func buildListOrderByBossQuery(db *gorm.DB, bossID uint, req *types.SellerOrderListReq) *gorm.DB {
	query := db.Table("`order` AS o").
		Joins("LEFT JOIN product as p ON p.id = o.product_id").
		Joins("LEFT JOIN address as a ON a.id = o.address_id").
		Joins("LEFT JOIN settlement as s ON s.order_id = o.id").
		Where("o.boss_id = ?", bossID).
		Where("o.deleted_at IS NULL")
	if req.Type != 0 {
		query = query.Where("o.type = ?", req.Type)
	}

	return query.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).Order("o.created_at DESC").
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"o.money AS money," +
			"o.refund_status AS refund_status," +
			"o.refund_reason AS refund_reason," +
			"o.payment_channel AS payment_channel," +
			"o.logistics_company AS logistics_company," +
			"o.tracking_no AS tracking_no," +
			"IFNULL(UNIX_TIMESTAMP(o.shipped_at), 0) AS shipped_at," +
			"IFNULL(UNIX_TIMESTAMP(o.received_at), 0) AS received_at," +
			"IFNULL(UNIX_TIMESTAMP(o.canceled_at), 0) AS canceled_at," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address," +
			"s.gross_amount AS gross_amount," +
			"s.commission_amount AS commission_amount," +
			"s.settlement_amount AS settlement_amount," +
			"s.status AS settlement_status")
}

func (dao *OrderDao) ListOrdersAdmin(req *types.AdminOrderListReq) (r []*types.OrderListResp, count int64, err error) {
	applyFilters := func(db *gorm.DB) *gorm.DB {
		if req.Type != 0 {
			db = db.Where("o.`type` = ?", req.Type)
		}
		if req.RefundStatus != nil {
			db = db.Where("o.`refund_status` = ?", *req.RefundStatus)
		}
		return db
	}

	countDB := applyFilters(dao.DB.Table("`order` AS o"))
	if err = countDB.Count(&count).Error; err != nil {
		return
	}

	query := dao.DB.
		Table("`order` AS o").
		Joins("LEFT JOIN product AS p ON p.id = o.product_id").
		Joins("LEFT JOIN address AS a ON a.id = o.address_id")

	query = applyFilters(query)

	err = query.
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Order("o.created_at DESC").
		Select(`
			o.id AS id,
			o.order_num AS order_num,
			UNIX_TIMESTAMP(o.created_at) AS created_at,
			UNIX_TIMESTAMP(o.updated_at) AS updated_at,
			o.user_id AS user_id,
			o.product_id AS product_id,
			o.boss_id AS boss_id,
			o.num AS num,
			o.type AS type,
			o.money AS money,
			o.refund_status AS refund_status,
			o.refund_reason AS refund_reason,
			o.payment_channel AS payment_channel,
			o.logistics_company AS logistics_company,
			o.tracking_no AS tracking_no,
			IFNULL(UNIX_TIMESTAMP(o.shipped_at), 0) AS shipped_at,
			IFNULL(UNIX_TIMESTAMP(o.received_at), 0) AS received_at,
			IFNULL(UNIX_TIMESTAMP(o.canceled_at), 0) AS canceled_at,
			p.name AS name,
			p.discount_price AS discount_price,
			p.img_path AS img_path,
			a.name AS address_name,
			a.phone AS address_phone,
			a.address AS address
		`).
		Scan(&r).Error

	return
}

func (dao *OrderDao) GetOrderById(id, uId uint) (r *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ?", id, uId).
		First(&r).Error

	return
}

func (dao *OrderDao) GetOrderByIdForUpdate(id uint) (r *model.Order, err error) {
	err = dao.DB.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		First(&r).Error
	return
}

func (dao *OrderDao) GetOrderByID(id uint) (r *model.Order, err error) {
	err = dao.DB.Model(&model.Order{}).
		Where("id = ?", id).
		First(&r).Error
	return
}

// ShowOrderById 获取订单详情
func (dao *OrderDao) ShowOrderById(id, uId uint) (r *types.OrderListResp, err error) {
	r = &types.OrderListResp{}
	err = dao.DB.Table("`order` AS o").
		Joins("LEFT JOIN product AS p ON p.id = o.product_id").
		Joins("LEFT JOIN address AS a ON a.id = o.address_id").
		Where("o.id = ? AND o.user_id = ?", id, uId).
		Where("o.deleted_at IS NULL").
		Where("o.buyer_deleted = ?", false).
		Select("o.id AS id," +
			"o.order_num AS order_num," +
			"UNIX_TIMESTAMP(o.created_at) AS created_at," +
			"UNIX_TIMESTAMP(o.updated_at) AS updated_at," +
			"o.user_id AS user_id," +
			"o.product_id AS product_id," +
			"o.boss_id AS boss_id," +
			"o.num AS num," +
			"o.type AS type," +
			"o.money AS money," +
			"o.refund_status AS refund_status," +
			"o.refund_reason AS refund_reason," +
			"o.payment_channel AS payment_channel," +
			"o.logistics_company AS logistics_company," +
			"o.tracking_no AS tracking_no," +
			"IFNULL(UNIX_TIMESTAMP(o.shipped_at), 0) AS shipped_at," +
			"IFNULL(UNIX_TIMESTAMP(o.received_at), 0) AS received_at," +
			"IFNULL(UNIX_TIMESTAMP(o.canceled_at), 0) AS canceled_at," +
			"p.name AS name," +
			"p.discount_price AS discount_price," +
			"p.img_path AS img_path," +
			"a.name AS address_name," +
			"a.phone AS address_phone," +
			"a.address AS address").
		Take(r).Error

	return
}

// DeleteOrderById 获取订单详情
func (dao *OrderDao) DeleteOrderById(id, uId uint) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id=? AND user_id = ?", id, uId).
		Update("buyer_deleted", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateOrderById 更新订单详情
func (dao *OrderDao) UpdateOrderById(id, uId uint, order *model.Order) error {
	return dao.DB.Where("id = ? AND user_id = ?", id, uId).
		Updates(order).Error
}

func (dao *OrderDao) UpdateOrderPaidById(id, uId uint, paidAt time.Time) error {
	return dao.UpdateOrderPaidByChannel(id, uId, paidAt, consts.OrderPaymentChannelBalance)
}

func (dao *OrderDao) UpdateOrderPaidByChannel(id, uId uint, paidAt time.Time, channel string) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type = ?", id, uId, consts.OrderTypeUnPaid).
		Updates(map[string]interface{}{
			"type":            2,
			"paid_at":         paidAt,
			"payment_channel": channel,
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

func (dao *OrderDao) UpdateOrderShippingByBoss(id, bossId uint, logisticsCompany, trackingNo string, shippedAt time.Time) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND boss_id = ? AND type = ?", id, bossId, consts.OrderTypePendingShipping).
		Updates(map[string]interface{}{
			"type":              consts.OrderTypeShipping,
			"logistics_company": logisticsCompany,
			"tracking_no":       trackingNo,
			"shipped_at":        shippedAt,
		})
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

func (dao *OrderDao) UpdateOrderReceivedByUser(id, uId uint, receivedAt time.Time) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type = ?", id, uId, consts.OrderTypeShipping).
		Updates(map[string]interface{}{
			"type":        consts.OrderTypeReceipt,
			"received_at": receivedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) CancelUnpaidOrderByUser(id, uId uint, canceledAt time.Time) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type = ?", id, uId, consts.OrderTypeUnPaid).
		Updates(map[string]interface{}{
			"type":        consts.OrderTypeCanceled,
			"canceled_at": canceledAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) RequestRefundByUser(id, uId uint, reason string) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND user_id = ? AND type IN ?", id, uId, []uint{
			consts.OrderTypePendingShipping,
			consts.OrderTypeShipping,
			consts.OrderTypeReceipt,
		}).
		Where("refund_status = ?", consts.OrderRefundStatusNone).
		Updates(map[string]interface{}{
			"type":          consts.OrderTypeRefundRequested,
			"refund_status": consts.OrderRefundStatusRequested,
			"refund_reason": reason,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (dao *OrderDao) MarkOrderRefunded(id uint) error {
	result := dao.DB.Model(&model.Order{}).
		Where("id = ? AND type = ? AND refund_status = ?", id, consts.OrderTypeRefundRequested, consts.OrderRefundStatusRequested).
		Updates(map[string]interface{}{
			"type":          consts.OrderTypeRefunded,
			"refund_status": consts.OrderRefundStatusRefunded,
		})
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
