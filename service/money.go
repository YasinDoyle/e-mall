package service

import (
	"context"
	"sync"

	"github.com/spf13/cast"

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
	money, err := user.DecryptMoney(req.Key)
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
