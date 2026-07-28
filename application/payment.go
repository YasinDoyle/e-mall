package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/rabbitmq"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

type PaymentUsecase struct{}

func NewPaymentUsecase() *PaymentUsecase {
	return &PaymentUsecase{}
}

func (u *PaymentUsecase) PayDown(ctx context.Context, req *types.PaymentDownReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var paidEvent *types.OrderPaidEvent
	var domainEvent event.OrderPaid
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		userID := userInfo.Id

		order, txErr := dao.NewOrderDaoByDB(tx).GetOrderById(req.OrderId, userID)
		if txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}
		if order.Type != consts.OrderTypeUnPaid {
			return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
		}
		if txErr = ensureNotBuyingOwnProduct(order.UserID, order.BossID); txErr != nil {
			return txErr
		}

		paidAt := time.Now()
		if txErr = dao.NewOrderDaoByDB(tx).UpdateOrderPaidById(req.OrderId, userID, paidAt); txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
			}
			log.LogrusObj.Error(txErr)
			return txErr
		}
		order.Type = consts.OrderTypePendingShipping
		order.PaidAt = &paidAt

		totalAmount := order.Money * float64(order.Num)
		paidEvent = &types.OrderPaidEvent{
			OrderID:     order.ID,
			OrderNum:    order.OrderNum,
			UserID:      order.UserID,
			BossID:      order.BossID,
			ProductID:   order.ProductID,
			Num:         order.Num,
			TotalAmount: totalAmount,
			PaidAt:      paidAt,
		}
		domainEvent = event.OrderPaid{
			OrderID:  order.ID,
			OrderNum: order.OrderNum,
			BuyerID:  order.UserID,
			SellerID: order.BossID,
		}

		userDao := dao.NewUserDaoByDB(tx)
		buyer, txErr := userDao.GetUserById(userID)
		if txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}
		if !buyer.HasPayKey() {
			return e.NewBusinessError(e.ErrorPaymentPayKeyRequired)
		}

		balance, txErr := buyer.DecryptMoney(req.Key)
		if txErr != nil {
			log.LogrusObj.Error(txErr)
			return e.NewBusinessError(e.ErrorPaymentPayKeyInvalid)
		}
		if balance-totalAmount < 0 {
			return e.NewBusinessError(e.ErrorPaymentBalanceInsufficient)
		}

		buyer.Money = fmt.Sprintf("%f", balance-totalAmount)
		buyer.Money, txErr = buyer.EncryptMoney(req.Key)
		if txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}
		if txErr = userDao.UpdateUserById(userID, buyer); txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}

		productDao := dao.NewProductDaoByDB(tx)
		if _, txErr = productDao.GetProductById(order.ProductID); txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}
		if txErr = productDao.DecreaseStock(order.ProductID, order.Num); txErr != nil {
			if errors.Is(txErr, gorm.ErrRecordNotFound) {
				return e.NewBusinessError(e.ErrorPaymentStockInsufficient)
			}
			log.LogrusObj.Error(txErr)
			return txErr
		}

		return handleOrderPaid(tx, order)
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	if paidEvent != nil {
		if zremErr := cache.RedisClient.ZRem(ctx, consts.OrderTimeKey, fmt.Sprintf("%d", paidEvent.OrderNum)).Err(); zremErr != nil {
			log.LogrusObj.Error(zremErr)
		}
		if publishErr := rabbitmq.PublishJSON(ctx, consts.OrderPaidQueue, paidEvent); publishErr != nil {
			log.LogrusObj.Error(publishErr)
		}
		event.Publish(ctx, domainEvent)
	}
	return nil, nil
}
