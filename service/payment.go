package service

import (
	"context"
	"sync"

	"github.com/YasinDoyle/e-mall/application"
	"github.com/YasinDoyle/e-mall/types"
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
	return application.NewPaymentUsecase().PayDown(ctx, req)
}
