package service

import (
	"context"
	"errors"
	"sync"

	"github.com/spf13/cast"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
)

var MoneySrvIns *MoneySrv
var MoneySrvOnce sync.Once

type MoneySrv struct {
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
	money, err := user.DecryptMoney(req.Key)
	if err != nil {
		err = errors.New("支付密码错误")
		return
	}
	resp = &types.MoneyShowResp{
		UserID:    user.ID,
		UserName:  user.UserName,
		UserMoney: cast.ToString(money),
	}

	return
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
