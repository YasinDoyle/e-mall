package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cast"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var MoneySrvIns *MoneySrv
var MoneySrvOnce sync.Once

type MoneySrv struct {
}

const moneyViewAuthTTL = 10 * time.Minute

func moneyViewAuthKey(userID uint) string {
	return fmt.Sprintf("money:view-auth:%d", userID)
}

func GetMoneySrv() *MoneySrv {
	MoneySrvOnce.Do(func() {
		MoneySrvIns = &MoneySrv{}
	})
	return MoneySrvIns
}

// MoneyShow 展示用户的金额
func (s *MoneySrv) MoneyShow(ctx context.Context, req *types.MoneyShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user, err := dao.NewUserDao(ctx).GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if !user.HasPayKey() {
		err = errors.New("请先设置支付密码")
		return
	}
	if req.Key == "" {
		ok := hasRecentMoneyViewAuth(ctx, user.ID)
		if !ok {
			err = errors.New("请输入支付密码")
			return
		}
	} else if !user.CheckPayKey(req.Key) {
		err = errors.New("支付密码错误")
		return
	} else {
		rememberMoneyViewAuth(ctx, user.ID)
	}
	money, err := user.DecryptMoney()
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.MoneyShowResp{
		UserID:    user.ID,
		UserName:  user.UserName,
		UserMoney: cast.ToString(money),
	}

	return
}

func hasRecentMoneyViewAuth(ctx context.Context, userID uint) bool {
	if cache.RedisClient == nil {
		return false
	}
	ok, err := cache.RedisClient.Exists(ctx, moneyViewAuthKey(userID)).Result()
	return err == nil && ok > 0
}

func rememberMoneyViewAuth(ctx context.Context, userID uint) {
	if cache.RedisClient == nil {
		return
	}
	if err := cache.RedisClient.Set(ctx, moneyViewAuthKey(userID), "1", moneyViewAuthTTL).Err(); err != nil {
		log.LogrusObj.Error(err)
	}
}

func (s *MoneySrv) SetPayKey(ctx context.Context, req *types.MoneySetPayKeyReq) (resp interface{}, err error) {
	if req.Key == "" || len(req.Key) != consts.EncryptMoneyKeyLength {
		return nil, errors.New("支付密码必须是6位")
	}
	if req.Key != req.KeyConfirm {
		return nil, errors.New("两次支付密码输入不一致")
	}

	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user, err := dao.NewUserDao(ctx).GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if user.HasPayKey() {
		return nil, errors.New("支付密码已设置")
	}
	if err = user.SetInitialMoneyWithPayKey(req.Key); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if err = dao.NewUserDao(ctx).UpdateUserById(u.Id, user); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return "设置成功", nil
}
