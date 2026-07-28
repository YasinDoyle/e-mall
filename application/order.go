package application

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

type OrderUsecase struct{}

func NewOrderUsecase() *OrderUsecase {
	return &OrderUsecase{}
}

func (u *OrderUsecase) AdminRefundApprove(ctx context.Context, req *types.AdminOrderRefundApproveReq) (interface{}, error) {
	refundAmount := float64(0)
	var orderNum uint64
	var buyerID uint
	var sellerID uint

	err := dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			return txErr
		}
		if order.Type != consts.OrderTypeRefundRequested || order.RefundStatus != consts.OrderRefundStatusRequested {
			return e.NewBusinessError(e.ErrorRefundStatusInvalid)
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

		return orderDao.MarkOrderRefunded(order.ID)
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
