package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/repository/rabbitmq"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var PaymentSrvIns *PaymentSrv
var PaymentSrvOnce sync.Once

type PaymentSrv struct {
}

func GetPaymentSrv() *PaymentSrv {
	PaymentSrvOnce.Do(func() {
		PaymentSrvIns = &PaymentSrv{}
	})
	return PaymentSrvIns
}

// TODO 目前买家和卖家的支付密码要一致，这个后续优化一下。。

// PayDown 支付操作
func (s *PaymentSrv) PayDown(ctx context.Context, req *types.PaymentDownReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	var paidEvent *types.OrderPaidEvent
	err = dao.NewOrderDao(ctx).Transaction(func(tx *gorm.DB) error {
		uId := u.Id

		payment, err := dao.NewOrderDaoByDB(tx).GetOrderById(req.OrderId, uId)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		if payment.Type != consts.OrderTypeUnPaid {
			return errors.New("订单已支付或状态不允许支付")
		}

		paidAt := time.Now()
		err = dao.NewOrderDaoByDB(tx).UpdateOrderPaidById(req.OrderId, uId, paidAt)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("订单已支付或状态不允许支付")
			}
			log.LogrusObj.Error(err)
			return err
		}
		payment.Type = consts.OrderTypePendingShipping
		payment.PaidAt = &paidAt

		money := payment.Money
		num := payment.Num
		money = money * float64(num)
		paidEvent = &types.OrderPaidEvent{
			OrderID:     payment.ID,
			OrderNum:    payment.OrderNum,
			UserID:      payment.UserID,
			BossID:      payment.BossID,
			ProductID:   payment.ProductID,
			Num:         payment.Num,
			TotalAmount: money,
			PaidAt:      paidAt,
		}

		userDao := dao.NewUserDaoByDB(tx)
		user, err := userDao.GetUserById(uId)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		// 对钱进行解密。减去订单。再进行加密。
		moneyFloat, err := user.DecryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		if moneyFloat-money < 0.0 { // 金额不足进行回滚
			log.LogrusObj.Error(err)
			return errors.New("金币不足")
		}

		finMoney := fmt.Sprintf("%f", moneyFloat-money)
		user.Money = finMoney
		user.Money, err = user.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		err = userDao.UpdateUserById(uId, user)
		if err != nil { // 更新用户金额失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		boss, err := userDao.GetUserById(uint(req.BossID))
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		moneyFloat, _ = boss.DecryptMoney(req.Key)
		finMoney = fmt.Sprintf("%f", moneyFloat+money)
		boss.Money = finMoney
		boss.Money, err = boss.EncryptMoney(req.Key)
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}

		err = userDao.UpdateUserById(uint(req.BossID), boss)
		if err != nil { // 更新boss金额失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		productDao := dao.NewProductDaoByDB(tx)
		product, err := productDao.GetProductById(uint(req.ProductID))
		if err != nil {
			log.LogrusObj.Error(err)
			return err
		}
		err = productDao.DecreaseStock(uint(req.ProductID), num)
		if err != nil { // 更新商品数量减少失败，回滚
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("库存不足")
			}
			log.LogrusObj.Error(err)
			return err
		}

		productUser := model.Product{
			Name:          product.Name,
			CategoryID:    product.CategoryID,
			Title:         product.Title,
			Info:          product.Info,
			ImgPath:       product.ImgPath,
			Price:         product.Price,
			DiscountPrice: product.DiscountPrice,
			Num:           num,
			OnSale:        false,
			BossID:        uId,
			BossName:      user.UserName,
			BossAvatar:    user.Avatar,
		}

		err = productDao.CreateProduct(&productUser)
		if err != nil { // 买完商品后创建成了自己的商品失败。订单失败，回滚
			log.LogrusObj.Error(err)
			return err
		}

		return nil

	})

	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if paidEvent != nil {
		if zremErr := cache.RedisClient.ZRem(ctx, OrderTimeKey, fmt.Sprintf("%d", paidEvent.OrderNum)).Err(); zremErr != nil {
			log.LogrusObj.Error(zremErr)
		}
	}
	if paidEvent != nil {
		if publishErr := rabbitmq.PublishJSON(ctx, consts.OrderPaidQueue, paidEvent); publishErr != nil {
			log.LogrusObj.Error(publishErr)
		}
	}

	return
}
