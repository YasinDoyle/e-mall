package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/application"
	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	util "github.com/YasinDoyle/e-mall/utils/log"
)

const orderTimeoutScanInterval = 5 * time.Second

var OrderSrvIns *OrderSrv
var OrderSrvOnce sync.Once

type OrderSrv struct {
}

func GetOrderSrv() *OrderSrv {
	OrderSrvOnce.Do(func() {
		OrderSrvIns = &OrderSrv{}
	})
	return OrderSrvIns
}

func (s *OrderSrv) OrderCreate(ctx context.Context, req *types.OrderCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}

	product, err := dao.NewProductDao(ctx).GetProductById(req.ProductID)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	if err = ensureNotBuyingOwnProduct(u.Id, product.BossID); err != nil {
		return nil, err
	}

	originalMoney := req.Money
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressByAid(req.AddressID, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productNum := strconv.Itoa(int(req.ProductID))
	userNum := strconv.Itoa(int(u.Id))
	number = number + productNum + userNum
	orderNum, _ := strconv.ParseUint(number, 10, 64)

	orderDao := dao.NewOrderDao(ctx)
	var order *model.Order
	err = orderDao.Transaction(func(tx *gorm.DB) error {
		finalMoney := originalMoney
		var userCouponID uint
		if req.CouponID > 0 {
			finalMoney, userCouponID, err = calcDiscount(dao.NewCouponDaoByDB(tx), u.Id, req.CouponID, originalMoney)
			if err != nil {
				return err
			}
		}

		order = &model.Order{
			UserID:    u.Id,
			ProductID: req.ProductID,
			BossID:    product.BossID,
			AddressID: address.ID,
			Num:       int(req.Num),
			Money:     finalMoney,
			Type:      consts.OrderTypeUnPaid,
			OrderNum:  orderNum,
		}
		if err = dao.NewOrderDaoByDB(tx).CreateOrder(order); err != nil {
			return err
		}

		if userCouponID > 0 {
			if err = dao.NewCouponDaoByDB(tx).UseCoupon(userCouponID, order.ID); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, dao.ErrCouponUseFailed) {
			return nil, errors.New("优惠券不存在或已使用")
		}
		util.LogrusObj.Error(err)
		return nil, err
	}

	// 订单号存入Redis中，设置过期时间
	data := redis.Z{
		Score:  float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		Member: orderNum,
	}
	cache.RedisClient.ZAdd(cache.RedisContext, consts.OrderTimeKey, data)

	resp = &types.OrderCreateResp{
		ID:       order.ID,
		OrderNum: order.OrderNum,
		Money:    order.Money,
		CouponID: req.CouponID,
	}

	return
}

func ensureNotBuyingOwnProduct(userID, bossID uint) error {
	if userID != 0 && userID == bossID {
		return e.NewBusinessError(e.ErrorOrderSelfPurchaseForbidden)
	}
	return nil
}

func StartOrderTimeoutWorker(ctx context.Context) {
	ticker := time.NewTicker(orderTimeoutScanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			consumeExpiredOrders(ctx)
		}
	}
}

func consumeExpiredOrders(ctx context.Context) {
	orderNums, err := cache.RedisClient.ZRangeByScore(ctx, consts.OrderTimeKey, &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: 50,
	}).Result()
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	if len(orderNums) == 0 {
		return
	}

	orderDao := dao.NewOrderDao(ctx)
	for _, orderNumStr := range orderNums {
		orderNum, parseErr := strconv.ParseUint(orderNumStr, 10, 64)
		if parseErr != nil {
			util.LogrusObj.Error(parseErr)
			_ = cache.RedisClient.ZRem(ctx, consts.OrderTimeKey, orderNumStr).Err()
			continue
		}

		if err = orderDao.DeleteUnpaidOrderByOrderNum(orderNum); err != nil {
			util.LogrusObj.Error(err)
			continue
		}

		if err = cache.RedisClient.ZRem(ctx, consts.OrderTimeKey, orderNumStr).Err(); err != nil {
			util.LogrusObj.Error(err)
		}
	}
}

func (s *OrderSrv) OrderList(ctx context.Context, req *types.OrderListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	orders, total, err := dao.NewOrderDao(ctx).ListOrderByCondition(u.Id, req)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	for i := range orders {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			orders[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}

	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}

	return
}

func (s *OrderSrv) AdminOrderList(ctx context.Context, req *types.AdminOrderListReq) (resp interface{}, err error) {
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}

	orders, total, err := dao.NewOrderDao(ctx).ListOrdersAdmin(req)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	for i := range orders {
		if conf.Config.System.UploadModel == consts.UploadModelLocal && orders[i].ImgPath != "" {
			orders[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}

	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}
	return
}

func (s *OrderSrv) SellerOrderList(ctx context.Context, req *types.SellerOrderListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}

	sellerProfile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(u.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSellerNotApproved)
		}
		util.LogrusObj.Error(err)
		return nil, err
	}
	if err = ensureSellerProfileApproved(sellerProfile); err != nil {
		return nil, err
	}

	orders, total, err := dao.NewOrderDao(ctx).ListOrderByBoss(u.Id, req)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	for i := range orders {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			orders[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}

	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}
	return
}

func (s *OrderSrv) OrderShow(ctx context.Context, req *types.OrderShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	order, err := dao.NewOrderDao(ctx).ShowOrderById(req.OrderId, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		order.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + order.ImgPath
	}

	resp = order

	return
}

func (s *OrderSrv) OrderDelete(ctx context.Context, req *types.OrderDeleteReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewOrderDao(ctx).DeleteOrderById(req.OrderId, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}

	return
}

func (s *OrderSrv) OrderShip(ctx context.Context, req *types.OrderShipReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().Ship(ctx, req)
}

func (s *OrderSrv) OrderRefundRequest(ctx context.Context, req *types.OrderRefundRequestReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().RefundRequest(ctx, req)
}

func (s *OrderSrv) AdminOrderRefundApprove(ctx context.Context, req *types.AdminOrderRefundApproveReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().AdminRefundApprove(ctx, req)
}

func (s *OrderSrv) OrderReceive(ctx context.Context, req *types.OrderReceiveReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().Receive(ctx, req)
}

func (s *OrderSrv) OrderCancel(ctx context.Context, req *types.OrderDeleteReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().CancelUnpaid(ctx, req.OrderId)
}

func (s *OrderSrv) OrderLogs(ctx context.Context, req *types.OrderShowReq) (resp interface{}, err error) {
	return application.NewOrderUsecase().Logs(ctx, req.OrderId)
}
