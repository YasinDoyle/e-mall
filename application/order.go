package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	aftersalestate "github.com/YasinDoyle/e-mall/domain/aftersale"
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
	afterSaleReasonMaxLen            = 255
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
		if afterSale, txErr := dao.NewAfterSaleDaoByDB(tx).GetByOrderIDForUpdate(order.ID); txErr == nil {
			if afterSale.Status == consts.AfterSaleStatusRefunded {
				return e.NewBusinessError(e.ErrorRefundStatusInvalid)
			}
		} else if !errors.Is(txErr, gorm.ErrRecordNotFound) {
			return txErr
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

		paymentChannel := order.PaymentChannel
		if paymentChannel == "" {
			paymentChannel = consts.OrderPaymentChannelBalance
		}
		if paymentChannel == consts.OrderPaymentChannelBalance {
			if txErr = refundBuyerBalanceInTx(tx, order.UserID, refundAmount); txErr != nil {
				return txErr
			}
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

func (u *OrderUsecase) RequestAfterSale(ctx context.Context, req *types.AfterSaleRequestReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	req.Type = strings.TrimSpace(req.Type)
	req.Reason = strings.TrimSpace(req.Reason)
	if err := validateAfterSaleRequest(req.Type, req.Reason); err != nil {
		return nil, err
	}

	var resp *types.AfterSaleListResp
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		order, txErr := dao.NewOrderDaoByDB(tx).GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userInfo.Id {
			return gorm.ErrRecordNotFound
		}
		if order.RefundStatus != consts.OrderRefundStatusNone {
			return e.NewBusinessError(e.ErrorAfterSaleStatusInvalid)
		}
		if !isAfterSaleOrderTypeAllowed(order.Type, req.Type) {
			return e.NewBusinessError(e.ErrorAfterSaleStatusInvalid)
		}

		hasActive, txErr := dao.NewAfterSaleDaoByDB(tx).HasActiveByOrderID(order.ID)
		if txErr != nil {
			return txErr
		}
		if hasActive {
			return e.NewBusinessError(e.ErrorAfterSaleStatusInvalid)
		}

		afterSale := &model.AfterSale{
			OrderID:      order.ID,
			OrderNum:     order.OrderNum,
			BuyerID:      order.UserID,
			SellerID:     order.BossID,
			Type:         req.Type,
			Status:       consts.AfterSaleStatusRequested,
			Reason:       req.Reason,
			RefundAmount: order.Money * float64(order.Num),
		}
		if txErr = dao.NewAfterSaleDaoByDB(tx).Create(afterSale); txErr != nil {
			return txErr
		}
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionAfterSale, order.Type, order.Type, "buyer", userInfo.Id, req.Type+": "+req.Reason)
		if txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		resp = buildAfterSaleResp(afterSale)
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorAfterSaleNotFound)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	return resp, nil
}

func (u *OrderUsecase) SellerHandleAfterSale(ctx context.Context, req *types.SellerAfterSaleHandleReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	req.Action = strings.TrimSpace(req.Action)
	req.Reason = strings.TrimSpace(req.Reason)
	if req.Action != consts.AfterSaleActionApprove && req.Action != consts.AfterSaleActionReject {
		return nil, e.NewBusinessError(e.ErrorAfterSaleActionInvalid)
	}
	if req.Action == consts.AfterSaleActionReject && req.Reason == "" {
		return nil, e.NewBusinessError(e.InvalidParams)
	}
	if len(req.Reason) > afterSaleReasonMaxLen {
		return nil, e.NewBusinessError(e.InvalidParams)
	}

	var resp *types.AfterSaleListResp
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		afterSale, txErr := dao.NewAfterSaleDaoByDB(tx).GetByIDForUpdate(req.AfterSaleID)
		if txErr != nil {
			return txErr
		}
		if afterSale.SellerID != userInfo.Id {
			return gorm.ErrRecordNotFound
		}

		nextStatus, txErr := sellerAfterSaleNextStatus(req.Action)
		if txErr != nil {
			return txErr
		}
		if txErr = aftersalestate.EnsureTransition(afterSale.Status, nextStatus); txErr != nil {
			return txErr
		}

		updates := map[string]interface{}{
			"status": nextStatus,
		}
		if req.Reason != "" {
			updates["seller_reason"] = req.Reason
		}
		if txErr = dao.NewAfterSaleDaoByDB(tx).UpdateByID(afterSale.ID, updates); txErr != nil {
			return txErr
		}
		order, txErr := dao.NewOrderDaoByDB(tx).GetOrderByID(afterSale.OrderID)
		if txErr != nil {
			return txErr
		}
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionAfterSale, order.Type, order.Type, "seller", userInfo.Id, req.Action)
		if txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		afterSale.Status = nextStatus
		afterSale.UpdatedAt = time.Now()
		if req.Reason != "" {
			afterSale.SellerReason = req.Reason
		}
		resp = buildAfterSaleResp(afterSale)
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorAfterSaleNotFound)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	return resp, nil
}

func (u *OrderUsecase) AdminHandleAfterSale(ctx context.Context, req *types.AdminAfterSaleHandleReq) (interface{}, error) {
	adminInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	req.Action = strings.TrimSpace(req.Action)
	req.Note = strings.TrimSpace(req.Note)
	if len(req.Note) > afterSaleReasonMaxLen {
		return nil, e.NewBusinessError(e.InvalidParams)
	}

	var resp *types.AfterSaleListResp
	var refundedEvent *event.OrderRefunded
	var closedEvent *event.AfterSaleClosed
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		admin, txErr := dao.NewUserDaoByDB(tx).GetUserById(adminInfo.Id)
		if txErr != nil {
			return txErr
		}
		if !admin.IsAdmin {
			return e.NewBusinessError(e.ErrorAuthInsufficientAuthority)
		}

		afterSale, txErr := dao.NewAfterSaleDaoByDB(tx).GetByIDForUpdate(req.AfterSaleID)
		if txErr != nil {
			return txErr
		}

		nextStatus, txErr := adminAfterSaleNextStatus(req.Action)
		if txErr != nil {
			return txErr
		}
		if txErr = aftersalestate.EnsureTransition(afterSale.Status, nextStatus); txErr != nil {
			return txErr
		}

		updates := map[string]interface{}{
			"status": nextStatus,
		}
		switch nextStatus {
		case consts.AfterSaleStatusPlatformIntervening:
			if req.Note != "" {
				updates["platform_note"] = req.Note
			}
		case consts.AfterSaleStatusRefunded:
			now := time.Now().Unix()
			updates["refunded_at"] = &now
			afterSale.RefundedAt = &now
			if req.Note != "" {
				updates["platform_note"] = req.Note
			}
			order, refundAmount, txErr := completeAfterSaleRefundInTx(tx, afterSale, adminInfo.Id)
			if txErr != nil {
				return txErr
			}
			refundedEvent = &event.OrderRefunded{
				OrderID:      order.ID,
				OrderNum:     order.OrderNum,
				BuyerID:      order.UserID,
				SellerID:     order.BossID,
				RefundAmount: refundAmount,
			}
		case consts.AfterSaleStatusClosed:
			if req.Note == "" {
				return e.NewBusinessError(e.InvalidParams)
			}
			now := time.Now().Unix()
			updates["closed_at"] = &now
			afterSale.ClosedAt = &now
			updates["platform_note"] = req.Note
		}
		if txErr = dao.NewAfterSaleDaoByDB(tx).UpdateByID(afterSale.ID, updates); txErr != nil {
			return txErr
		}
		order, txErr := dao.NewOrderDaoByDB(tx).GetOrderByID(afterSale.OrderID)
		if txErr != nil {
			return txErr
		}
		if nextStatus == consts.AfterSaleStatusClosed {
			closedEvent = &event.AfterSaleClosed{
				OrderID:  order.ID,
				OrderNum: order.OrderNum,
				BuyerID:  order.UserID,
				SellerID: order.BossID,
				Note:     req.Note,
			}
		}
		orderLog, txErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionAfterSale, order.Type, order.Type, "admin", adminInfo.Id, req.Action)
		if txErr != nil {
			return txErr
		}
		if txErr = dao.NewOrderLogDaoByDB(tx).Create(orderLog); txErr != nil {
			return txErr
		}
		afterSale.Status = nextStatus
		afterSale.UpdatedAt = time.Now()
		if note, ok := updates["platform_note"]; ok {
			afterSale.PlatformNote = note.(string)
		}
		resp = buildAfterSaleResp(afterSale)
		return nil
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorAfterSaleNotFound)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}

	if refundedEvent != nil {
		event.Publish(ctx, *refundedEvent)
	}
	if closedEvent != nil {
		event.Publish(ctx, *closedEvent)
	}
	return resp, nil
}

func completeAfterSaleRefundInTx(tx *gorm.DB, afterSale *model.AfterSale, operatorID uint) (*model.Order, float64, error) {
	orderDao := dao.NewOrderDaoByDB(tx)
	order, err := orderDao.GetOrderByIdForUpdate(afterSale.OrderID)
	if err != nil {
		return nil, 0, err
	}
	if order.Type == consts.OrderTypeRefunded || order.RefundStatus == consts.OrderRefundStatusRefunded {
		return nil, 0, e.NewBusinessError(e.ErrorRefundStatusInvalid)
	}
	if order.Type != consts.OrderTypeRefundRequested || order.RefundStatus != consts.OrderRefundStatusRequested {
		orderLog, logErr := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionAfterSale, order.Type, consts.OrderTypeRefundRequested, "admin", operatorID, afterSale.Type+": "+afterSale.Reason)
		if logErr != nil {
			return nil, 0, logErr
		}
		if err = orderDao.RequestRefundForAfterSale(order.ID, afterSale.Reason); err != nil {
			return nil, 0, err
		}
		if err = dao.NewOrderLogDaoByDB(tx).Create(orderLog); err != nil {
			return nil, 0, err
		}
		order.Type = consts.OrderTypeRefundRequested
		order.RefundStatus = consts.OrderRefundStatusRequested
		order.RefundReason = afterSale.Reason
	}

	refundAmount := order.Money * float64(order.Num)
	if refundAmount <= 0 {
		return nil, 0, e.NewBusinessError(e.ErrorRefundAmountInvalid)
	}
	paymentChannel := order.PaymentChannel
	if paymentChannel == "" {
		paymentChannel = consts.OrderPaymentChannelBalance
	}
	if paymentChannel == consts.OrderPaymentChannelBalance {
		if err = refundBuyerBalanceInTx(tx, order.UserID, refundAmount); err != nil {
			return nil, 0, err
		}
	}
	if err = handleOrderRefunded(tx, order); err != nil {
		return nil, 0, err
	}
	orderLog, err := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionRefundApprove, order.Type, consts.OrderTypeRefunded, "admin", operatorID, "after-sale refund")
	if err != nil {
		return nil, 0, err
	}
	if err = orderDao.MarkOrderRefunded(order.ID); err != nil {
		return nil, 0, err
	}
	if err = dao.NewOrderLogDaoByDB(tx).Create(orderLog); err != nil {
		return nil, 0, err
	}
	order.Type = consts.OrderTypeRefunded
	order.RefundStatus = consts.OrderRefundStatusRefunded
	return order, refundAmount, nil
}

func refundBuyerBalanceInTx(tx *gorm.DB, buyerID uint, refundAmount float64) error {
	userDao := dao.NewUserDaoByDB(tx)
	buyer, err := userDao.GetUserById(buyerID)
	if err != nil {
		return err
	}
	if !buyer.HasPayKey() {
		return e.NewBusinessError(e.ErrorPaymentPayKeyRequired)
	}
	buyerMoney, err := buyer.DecryptMoney()
	if err != nil {
		return err
	}
	buyer.Money = fmt.Sprintf("%f", buyerMoney+refundAmount)
	buyer.Money, err = buyer.EncryptMoney()
	if err != nil {
		return err
	}
	return userDao.UpdateUserById(buyerID, buyer)
}

func (u *OrderUsecase) ListBuyerAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	normalizeAfterSaleListReq(req)

	items, total, err := dao.NewAfterSaleDao(ctx).ListByBuyer(userInfo.Id, req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return &types.DataListResp{Item: buildAfterSaleRespList(items), Total: total}, nil
}

func (u *OrderUsecase) ListSellerAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	normalizeAfterSaleListReq(req)

	items, total, err := dao.NewAfterSaleDao(ctx).ListBySeller(userInfo.Id, req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return &types.DataListResp{Item: buildAfterSaleRespList(items), Total: total}, nil
}

func (u *OrderUsecase) ListAdminAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error) {
	adminInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	normalizeAfterSaleListReq(req)

	admin, err := dao.NewUserDao(ctx).GetUserById(adminInfo.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if !admin.IsAdmin {
		return nil, e.NewBusinessError(e.ErrorAuthInsufficientAuthority)
	}

	items, total, err := dao.NewAfterSaleDao(ctx).ListAdmin(req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return &types.DataListResp{Item: buildAfterSaleRespList(items), Total: total}, nil
}

func validateAfterSaleRequest(afterSaleType, reason string) error {
	if reason == "" || len(reason) > afterSaleReasonMaxLen {
		return e.NewBusinessError(e.InvalidParams)
	}
	switch afterSaleType {
	case consts.AfterSaleTypeRefundOnly, consts.AfterSaleTypeReturnRefund:
		return nil
	default:
		return e.NewBusinessError(e.InvalidParams)
	}
}

func isAfterSaleOrderTypeAllowed(orderType uint, afterSaleType string) bool {
	switch orderType {
	case consts.OrderTypePendingShipping, consts.OrderTypeShipping:
		return true
	case consts.OrderTypeReceipt:
		return afterSaleType == consts.AfterSaleTypeReturnRefund || afterSaleType == consts.AfterSaleTypeRefundOnly
	default:
		return false
	}
}

func sellerAfterSaleNextStatus(action string) (string, error) {
	switch action {
	case consts.AfterSaleActionApprove:
		return consts.AfterSaleStatusSellerApproved, nil
	case consts.AfterSaleActionReject:
		return consts.AfterSaleStatusSellerRejected, nil
	default:
		return "", e.NewBusinessError(e.ErrorAfterSaleActionInvalid)
	}
}

func adminAfterSaleNextStatus(action string) (string, error) {
	switch action {
	case consts.AfterSaleActionIntervene:
		return consts.AfterSaleStatusPlatformIntervening, nil
	case consts.AfterSaleActionRefund:
		return consts.AfterSaleStatusRefunded, nil
	case consts.AfterSaleActionClose:
		return consts.AfterSaleStatusClosed, nil
	default:
		return "", e.NewBusinessError(e.ErrorAfterSaleActionInvalid)
	}
}

func buildAfterSaleResp(afterSale *model.AfterSale) *types.AfterSaleListResp {
	resp := &types.AfterSaleListResp{
		ID:           afterSale.ID,
		OrderID:      afterSale.OrderID,
		OrderNum:     afterSale.OrderNum,
		BuyerID:      afterSale.BuyerID,
		SellerID:     afterSale.SellerID,
		Type:         afterSale.Type,
		Status:       afterSale.Status,
		Reason:       afterSale.Reason,
		RefundAmount: afterSale.RefundAmount,
		SellerReason: afterSale.SellerReason,
		PlatformNote: afterSale.PlatformNote,
		CreatedAt:    afterSale.CreatedAt.Unix(),
		UpdatedAt:    afterSale.UpdatedAt.Unix(),
	}
	if afterSale.RefundedAt != nil {
		resp.RefundedAt = *afterSale.RefundedAt
	}
	if afterSale.ClosedAt != nil {
		resp.ClosedAt = *afterSale.ClosedAt
	}
	return resp
}

func buildAfterSaleRespList(items []*model.AfterSale) []*types.AfterSaleListResp {
	resp := make([]*types.AfterSaleListResp, 0, len(items))
	for _, item := range items {
		resp = append(resp, buildAfterSaleResp(item))
	}
	return resp
}

func normalizeAfterSaleListReq(req *types.AfterSaleListReq) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = consts.BasePageSize
	}
}
