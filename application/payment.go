package application

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/repository/rabbitmq"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

type PaymentUsecase struct{}

const orderPaymentNoPrefix = "OP"

type orderPaidArtifacts struct {
	paidEvent   *types.OrderPaidEvent
	domainEvent event.OrderPaid
}

func NewPaymentUsecase() *PaymentUsecase {
	return &PaymentUsecase{}
}

func IsOrderPaymentNo(paymentNo string) bool {
	if !strings.HasPrefix(paymentNo, orderPaymentNoPrefix) || len(paymentNo) == len(orderPaymentNoPrefix) {
		return false
	}
	for _, r := range paymentNo[len(orderPaymentNoPrefix):] {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func (u *PaymentUsecase) PayDown(ctx context.Context, req *types.PaymentDownReq) (interface{}, error) {
	return u.PayOrderByBalance(ctx, req)
}

func (u *PaymentUsecase) PayOrderByBalance(ctx context.Context, req *types.PaymentDownReq) (interface{}, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var artifacts orderPaidArtifacts
	var resp *types.OrderPaymentResp
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		userID := userInfo.Id

		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(req.OrderId)
		if txErr != nil {
			log.LogrusObj.Error(txErr)
			return txErr
		}
		if order.UserID != userID || order.Type != consts.OrderTypeUnPaid {
			return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
		}
		if txErr = ensureNotBuyingOwnProduct(order.UserID, order.BossID); txErr != nil {
			return txErr
		}

		paidAt := time.Now()
		totalAmount := order.Money * float64(order.Num)
		paymentDao := dao.NewOrderPaymentDaoByDB(tx)
		pending, txErr := paymentDao.GetPendingByOrderID(order.ID)
		if txErr = ensureNoPendingExternalPaymentForBalance(pending, txErr); txErr != nil {
			return txErr
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

		payment := newOrderPayment(order, consts.OrderPaymentChannelBalance, totalAmount)
		if txErr = paymentDao.ClosePendingByOrderID(order.ID); txErr != nil {
			return txErr
		}
		if txErr = paymentDao.Create(payment); txErr != nil {
			return txErr
		}
		if _, txErr = paymentDao.MarkPaid(payment.PaymentNo, "balance", paidAt); txErr != nil {
			return txErr
		}

		if artifacts, txErr = completeOrderPaidInTx(tx, order, consts.OrderPaymentChannelBalance, paidAt, userID); txErr != nil {
			return txErr
		}
		resp = buildOrderPaymentResp(payment)
		resp.Status = consts.OrderPaymentStatusPaid
		resp.PaidAt = paidAt.Unix()
		resp.OrderStatus = consts.OrderTypePendingShipping
		return nil
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	publishOrderPaidArtifacts(ctx, artifacts)
	return resp, nil
}

func (u *PaymentUsecase) CreateOrderGatewayPayment(ctx context.Context, orderID uint, channel string) (*types.OrderPaymentResp, error) {
	if channel != consts.OrderPaymentChannelWechat && channel != consts.OrderPaymentChannelAlipay {
		return nil, e.NewBusinessError(e.InvalidParams)
	}
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	var resp *types.OrderPaymentResp
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(orderID)
		if txErr != nil {
			return txErr
		}
		if order.UserID != userInfo.Id || order.Type != consts.OrderTypeUnPaid {
			return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
		}
		if txErr = ensureNotBuyingOwnProduct(order.UserID, order.BossID); txErr != nil {
			return txErr
		}

		paymentDao := dao.NewOrderPaymentDaoByDB(tx)
		pending, txErr := paymentDao.GetPendingByOrderID(order.ID)
		reusable, createNew, txErr := resolvePendingOrderGatewayPayment(pending, txErr, channel)
		if txErr != nil {
			return txErr
		}
		if !createNew {
			resp = buildOrderPaymentResp(reusable)
			return nil
		}

		payment := newOrderPayment(order, channel, order.Money*float64(order.Num))
		if txErr = paymentDao.Create(payment); txErr != nil {
			return txErr
		}
		resp = buildOrderPaymentResp(payment)
		return nil
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return resp, nil
}

func (u *PaymentUsecase) MarkOrderPaymentFailed(ctx context.Context, paymentNo string) error {
	if paymentNo == "" {
		return e.NewBusinessError(e.InvalidParams)
	}
	return dao.NewOrderPaymentDao(ctx).MarkFailed(paymentNo)
}

func (u *PaymentUsecase) OrderPaymentStatus(ctx context.Context, paymentNo string) (*types.OrderPaymentResp, error) {
	userInfo, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}
	payment, err := dao.NewOrderPaymentDao(ctx).GetByPaymentNo(paymentNo)
	if err != nil {
		return nil, err
	}
	if payment.UserID != userInfo.Id {
		return nil, gorm.ErrRecordNotFound
	}
	resp := buildOrderPaymentResp(payment)
	order, err := dao.NewOrderDao(ctx).GetOrderByID(payment.OrderID)
	if err == nil {
		resp.OrderStatus = order.Type
	}
	return resp, nil
}

func (u *PaymentUsecase) HandleOrderPaymentCallback(ctx context.Context, req *types.OrderPaymentCallbackReq) error {
	if req == nil || req.PaymentNo == "" {
		return e.NewBusinessError(e.InvalidParams)
	}
	var artifacts orderPaidArtifacts
	err := dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		paymentDao := dao.NewOrderPaymentDaoByDB(tx)
		payment, txErr := paymentDao.GetByPaymentNoForUpdate(req.PaymentNo)
		if txErr != nil {
			return txErr
		}
		if txErr = validateOrderPaymentCallback(payment, req.PaidAmount, req.ExpectedChannel); txErr != nil {
			return txErr
		}
		if payment.Status == consts.OrderPaymentStatusPaid {
			return nil
		}

		orderDao := dao.NewOrderDaoByDB(tx)
		order, txErr := orderDao.GetOrderByIdForUpdate(payment.OrderID)
		if txErr != nil {
			return txErr
		}
		if order.ID != payment.OrderID || order.UserID != payment.UserID || order.Type != consts.OrderTypeUnPaid {
			return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
		}
		if roundMoney(order.Money*float64(order.Num)) != roundMoney(payment.Amount) {
			return e.NewBusinessError(e.ErrorPaymentAmountMismatch)
		}

		paidAt := req.PaidAt
		if paidAt.IsZero() {
			paidAt = time.Now()
		}
		fresh, txErr := paymentDao.MarkPaid(payment.PaymentNo, req.ProviderTradeNo, paidAt)
		if txErr != nil {
			return txErr
		}
		if !fresh {
			return nil
		}

		artifacts, txErr = completeOrderPaidInTx(tx, order, payment.Channel, paidAt, payment.UserID)
		if txErr != nil {
			return txErr
		}
		return paymentDao.CloseOtherPendingByOrderID(order.ID, payment.PaymentNo)
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return err
	}
	publishOrderPaidArtifacts(ctx, artifacts)
	return nil
}

func validateOrderPaymentCallback(payment *model.OrderPayment, paidAmount float64, expectedChannel string) error {
	if payment == nil {
		return e.NewBusinessError(e.InvalidParams)
	}
	if expectedChannel != "" && payment.Channel != expectedChannel {
		return e.NewBusinessError(e.InvalidParams)
	}
	if math.Abs(roundMoney(payment.Amount)-roundMoney(paidAmount)) > 0.001 {
		return e.NewBusinessError(e.ErrorPaymentAmountMismatch)
	}
	if payment.Status == consts.OrderPaymentStatusPaid || payment.Status == consts.OrderPaymentStatusPending {
		return nil
	}
	return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
}

func resolvePendingOrderGatewayPayment(pending *model.OrderPayment, lookupErr error, requestedChannel string) (*model.OrderPayment, bool, error) {
	if errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return nil, true, nil
	}
	if lookupErr != nil {
		return nil, false, lookupErr
	}
	if pending == nil {
		return nil, true, nil
	}
	if pending.Channel != requestedChannel {
		return nil, false, e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
	}
	return pending, false, nil
}

func ensureNoPendingExternalPaymentForBalance(pending *model.OrderPayment, lookupErr error) error {
	if errors.Is(lookupErr, gorm.ErrRecordNotFound) {
		return nil
	}
	if lookupErr != nil {
		return lookupErr
	}
	if pending != nil && pending.Channel != consts.OrderPaymentChannelBalance {
		return e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
	}
	return nil
}

func newOrderPayment(order *model.Order, channel string, amount float64) *model.OrderPayment {
	return &model.OrderPayment{
		OrderID:   order.ID,
		OrderNum:  order.OrderNum,
		PaymentNo: genOrderPaymentNo(),
		UserID:    order.UserID,
		Channel:   channel,
		Amount:    roundMoney(amount),
		Status:    consts.OrderPaymentStatusPending,
	}
}

func genOrderPaymentNo() string {
	return fmt.Sprintf("%s%d", orderPaymentNoPrefix, time.Now().UnixNano())
}

func buildOrderPaymentResp(payment *model.OrderPayment) *types.OrderPaymentResp {
	resp := &types.OrderPaymentResp{
		OrderID:   payment.OrderID,
		OrderNum:  payment.OrderNum,
		PaymentNo: payment.PaymentNo,
		Channel:   payment.Channel,
		Amount:    payment.Amount,
		Status:    payment.Status,
	}
	if payment.PaidAt != nil {
		resp.PaidAt = payment.PaidAt.Unix()
	}
	if payment.ClosedAt != nil {
		resp.ClosedAt = payment.ClosedAt.Unix()
	}
	return resp
}

func completeOrderPaidInTx(tx *gorm.DB, order *model.Order, channel string, paidAt time.Time, operatorID uint) (orderPaidArtifacts, error) {
	orderLog, err := buildOrderLog(order.ID, order.OrderNum, consts.OrderActionPay, order.Type, consts.OrderTypePendingShipping, "buyer", operatorID, channel)
	if err != nil {
		return orderPaidArtifacts{}, err
	}

	orderDao := dao.NewOrderDaoByDB(tx)
	if err = orderDao.UpdateOrderPaidByChannel(order.ID, order.UserID, paidAt, channel); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orderPaidArtifacts{}, e.NewBusinessError(e.ErrorOrderPayStatusInvalid)
		}
		return orderPaidArtifacts{}, err
	}
	order.Type = consts.OrderTypePendingShipping
	order.PaidAt = &paidAt
	order.PaymentChannel = channel

	productDao := dao.NewProductDaoByDB(tx)
	if _, err = productDao.GetProductById(order.ProductID); err != nil {
		log.LogrusObj.Error(err)
		return orderPaidArtifacts{}, err
	}
	if err = productDao.DecreaseStock(order.ProductID, order.Num); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return orderPaidArtifacts{}, e.NewBusinessError(e.ErrorPaymentStockInsufficient)
		}
		return orderPaidArtifacts{}, err
	}

	if err = dao.NewOrderLogDaoByDB(tx).Create(orderLog); err != nil {
		return orderPaidArtifacts{}, err
	}
	if err = handleOrderPaid(tx, order); err != nil {
		return orderPaidArtifacts{}, err
	}

	totalAmount := roundMoney(order.Money * float64(order.Num))
	return orderPaidArtifacts{
		paidEvent: &types.OrderPaidEvent{
			OrderID:     order.ID,
			OrderNum:    order.OrderNum,
			UserID:      order.UserID,
			BossID:      order.BossID,
			ProductID:   order.ProductID,
			Num:         order.Num,
			TotalAmount: totalAmount,
			PaidAt:      paidAt,
		},
		domainEvent: event.OrderPaid{
			OrderID:  order.ID,
			OrderNum: order.OrderNum,
			BuyerID:  order.UserID,
			SellerID: order.BossID,
		},
	}, nil
}

func publishOrderPaidArtifacts(ctx context.Context, artifacts orderPaidArtifacts) {
	if artifacts.paidEvent == nil {
		return
	}
	if zremErr := cache.RedisClient.ZRem(ctx, consts.OrderTimeKey, fmt.Sprintf("%d", artifacts.paidEvent.OrderNum)).Err(); zremErr != nil {
		log.LogrusObj.Error(zremErr)
	}
	if publishErr := rabbitmq.PublishJSON(ctx, consts.OrderPaidQueue, artifacts.paidEvent); publishErr != nil {
		log.LogrusObj.Error(publishErr)
	}
	event.Publish(ctx, artifacts.domainEvent)
}
