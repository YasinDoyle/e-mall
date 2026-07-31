package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/domain/orderstate"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

type OrderUsecase struct{}

const (
	orderLogisticsNodeManualShipped  = "manual_shipped"
	orderLogisticsNodeManualReceived = "manual_received"
	shipmentInfoMaxLen               = 64
)

func NewOrderUsecase() *OrderUsecase {
	return &OrderUsecase{}
}

func buildOrderLog(orderID uint, orderNum uint64, action string, fromType, toType uint, operatorType string, operatorID uint, remark string) (*model.OrderLog, error) {
	if err := orderstate.EnsureOrderStatusTransition(fromType, toType); err != nil {
		return nil, err
	}
	return &model.OrderLog{
		OrderID:      orderID,
		OrderNum:     orderNum,
		Action:       action,
		FromType:     fromType,
		ToType:       toType,
		OperatorType: operatorType,
		OperatorID:   operatorID,
		Remark:       remark,
	}, nil
}

func (u *OrderUsecase) CancelUnpaid(ctx context.Context, orderID uint) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	var orderNum uint64
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(orderID)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userInfo.Id || order.Type != consts.OrderTypeUnPaid {
			return gorm.ErrRecordNotFound
		}
		orderNum = order.OrderNum

		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionCancel, order.Type, consts.OrderTypeCanceled, "buyer", userInfo.Id, "cancel unpaid order")
		if txErr != nil {
			return txErr
		}
		if txErr = orderDao.CancelUnpaidOrderByUser(order.ID, userInfo.Id, time.Now()); txErr != nil {
			return txErr
		}
		return dao.NewOrderLogDaoByDB(tx).Create(orderLog)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	return map[string]interface{}{
		"order_id":  orderID,
		"order_num": orderNum,
		"type":      consts.OrderTypeCanceled,
	}, nil
}

func (u *OrderUsecase) Ship(ctx context.Context, req *types.OrderShipReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	req.LogisticsCompany = strings.TrimSpace(req.LogisticsCompany)
	req.TrackingNo = strings.TrimSpace(req.TrackingNo)
	if req.LogisticsCompany == "" || req.TrackingNo == "" ||
		len(req.LogisticsCompany) > shipmentInfoMaxLen || len(req.TrackingNo) > shipmentInfoMaxLen {
		return nil, e.NewBusinessError(e.InvalidParams)
	}

	var shipped event.OrderShipped
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.BossID != userInfo.Id || order.Type != consts.OrderTypePendingShipping {
			return gorm.ErrRecordNotFound
		}

		shippedAt := time.Now()
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionShip, order.Type, consts.OrderTypeShipping, "seller", userInfo.Id, req.TrackingNo)
		if txErr != nil {
			return txErr
		}
		if txErr = orderDao.UpdateOrderShippingByBoss(order.ID, userInfo.Id, req.LogisticsCompany, req.TrackingNo, shippedAt); txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogisticsDaoByDB(tx).Create(&model.OrderLogistics{
			OrderID:     order.ID,
			OrderNum:    order.OrderNum,
			NodeType:    orderLogisticsNodeManualShipped,
			Description: req.LogisticsCompany + " " + req.TrackingNo,
			OccurredAt:  shippedAt.Unix(),
		}); txErr != nil {
			return txErr
		}
		shipped = event.OrderShipped{
			OrderID:    order.ID,
			OrderNum:   order.OrderNum,
			BuyerID:    order.UserID,
			TrackingNo: req.TrackingNo,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	event.Publish(ctx, shipped)
	return nil, nil
}

func (u *OrderUsecase) Receive(ctx context.Context, req *types.OrderReceiveReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	var received event.OrderReceived
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userInfo.Id || order.Type != consts.OrderTypeShipping {
			return gorm.ErrRecordNotFound
		}

		receivedAt := time.Now()
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionReceive, order.Type, consts.OrderTypeReceipt, "buyer", userInfo.Id, "confirm receipt")
		if txErr != nil {
			return txErr
		}
		if txErr = orderDao.UpdateOrderReceivedByUser(order.ID, userInfo.Id, receivedAt); txErr != nil {
			return txErr
		}
		if txErr = dao.NewSettlementDaoByDB(tx).GenerateCompletedForOrder(order.ID); txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogisticsDaoByDB(tx).Create(&model.OrderLogistics{
			OrderID:     order.ID,
			OrderNum:    order.OrderNum,
			NodeType:    orderLogisticsNodeManualReceived,
			Description: "buyer confirmed receipt",
			OccurredAt:  receivedAt.Unix(),
		}); txErr != nil {
			return txErr
		}
		received = event.OrderReceived{
			OrderID:  order.ID,
			OrderNum: order.OrderNum,
			SellerID: order.BossID,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	event.Publish(ctx, received)
	return nil, nil
}

func (u *OrderUsecase) RefundRequest(ctx context.Context, req *types.OrderRefundRequestReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	var refundRequested event.RefundRequested
	var resp *types.OrderRefundResp
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userInfo.Id || order.RefundStatus != consts.OrderRefundStatusNone {
			return gorm.ErrRecordNotFound
		}

		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionRefundRequest, order.Type, consts.OrderTypeRefundRequested, "buyer", userInfo.Id, req.Reason)
		if txErr != nil {
			return txErr
		}
		if txErr = orderDao.RequestRefundByUser(order.ID, userInfo.Id, req.Reason); txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		refundRequested = event.RefundRequested{
			OrderID:  order.ID,
			OrderNum: order.OrderNum,
			SellerID: order.BossID,
		}
		resp = &types.OrderRefundResp{
			OrderId:      order.ID,
			OrderNum:     order.OrderNum,
			RefundAmount: order.Money * float64(order.Num),
			RefundStatus: consts.OrderRefundStatusRequested,
			Type:         consts.OrderTypeRefundRequested,
		}
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	event.Publish(ctx, refundRequested)
	return resp, nil
}

func (u *OrderUsecase) Logs(ctx context.Context, orderID uint) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	order, err := dao.NewOrderDao(ctx).GetOrderByID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	if order.UserID != userInfo.Id && order.BossID != userInfo.Id {
		return nil, gorm.ErrRecordNotFound
	}

	logs, err := dao.NewOrderLogDao(ctx).ListByOrderID(orderID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp := make([]*types.OrderLogResp, 0, len(logs))
	for _, item := range logs {
		resp = append(resp, &types.OrderLogResp{
			ID:           item.ID,
			OrderID:      item.OrderID,
			OrderNum:     item.OrderNum,
			Action:       item.Action,
			FromType:     item.FromType,
			ToType:       item.ToType,
			OperatorType: item.OperatorType,
			OperatorID:   item.OperatorID,
			Remark:       item.Remark,
			CreatedAt:    item.CreatedAt.Unix(),
		})
	}
	return resp, nil
}

func (u *OrderUsecase) AdminRefundApprove(ctx context.Context, req *types.AdminOrderRefundApproveReq) (interface{}, error) {
	adminInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	refundAmount := float64(0)
	var orderNum uint64
	var buyerID uint
	var sellerID uint

	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		admin, txErr := dao.NewUserDaoByDB(tx).GetUserById(adminInfo.Id)
		if txErr != nil {
			return txErr
		}
		if !admin.IsAdmin {
			return e.NewBusinessError(e.ErrorAuthInsufficientAuthority)
		}

		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.Type != consts.OrderTypeRefundRequested || order.RefundStatus != consts.OrderRefundStatusRequested {
			return e.NewBusinessError(e.ErrorRefundStatusInvalid)
		}
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionRefundApprove, order.Type, consts.OrderTypeRefunded, "admin", adminInfo.Id, "approve refund")
		if txErr != nil {
			return txErr
		}

		refundAmount = order.Money * float64(order.Num)
		if refundAmount <= 0 {
			return e.NewBusinessError(e.ErrorRefundAmountInvalid)
		}
		orderNum = order.OrderNum
		buyerID = order.UserID
		sellerID = order.BossID

		userDao := dao.NewUserDaoByDB(tx)
		buyer, txErr := userDao.GetUserById(order.UserID)
		if txErr != nil {
			return txErr
		}
		if !buyer.HasPayKey() {
			return e.NewBusinessError(e.ErrorPaymentPayKeyRequired)
		}

		buyerMoney, txErr := buyer.DecryptMoney(req.Key)
		if txErr != nil {
			return e.NewBusinessError(e.ErrorPaymentPayKeyInvalid)
		}

		buyer.Money = fmt.Sprintf("%f", buyerMoney+refundAmount)
		buyer.Money, txErr = buyer.EncryptMoney(req.Key)
		if txErr != nil {
			return txErr
		}
		if txErr = userDao.UpdateUserById(order.UserID, buyer); txErr != nil {
			return txErr
		}

		if txErr = handleOrderRefunded(tx, order); txErr != nil {
			return txErr
		}

		if txErr = orderDao.MarkOrderRefunded(order.ID); txErr != nil {
			return txErr
		}
		return dao.NewOrderLogDaoByDB(tx).Create(orderLog)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorRefundNotFound)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	event.Publish(ctx, event.OrderRefunded{
		OrderID:      req.OrderId,
		OrderNum:     orderNum,
		BuyerID:      buyerID,
		SellerID:     sellerID,
		RefundAmount: refundAmount,
	})
	return &types.OrderRefundResp{
		OrderId:      req.OrderId,
		OrderNum:     orderNum,
		RefundAmount: refundAmount,
		RefundStatus: consts.OrderRefundStatusRefunded,
		Type:         consts.OrderTypeRefunded,
	}, nil
}
