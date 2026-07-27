package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

const defaultCommissionRate = 0.05

type SettlementSrv struct{}

var settlementSrvIns *SettlementSrv
var settlementSrvOnce sync.Once

func GetSettlementSrv() *SettlementSrv {
	settlementSrvOnce.Do(func() {
		settlementSrvIns = &SettlementSrv{}
	})
	return settlementSrvIns
}

type settlementCalculation struct {
	GrossAmount      float64
	CommissionAmount float64
	SettlementAmount float64
}

func calculateOrderSettlement(unitPrice float64, num int, commissionRate float64) (*settlementCalculation, error) {
	if unitPrice <= 0 || num <= 0 {
		return nil, e.NewBusinessError(e.ErrorSettlementInvalidAmount)
	}
	if commissionRate < 0 || commissionRate >= 1 {
		return nil, e.NewBusinessError(e.ErrorSettlementInvalidRate)
	}
	gross := roundMoney(unitPrice * float64(num))
	commission := roundMoney(gross * commissionRate)
	return &settlementCalculation{
		GrossAmount:      gross,
		CommissionAmount: commission,
		SettlementAmount: roundMoney(gross - commission),
	}, nil
}

func roundMoney(amount float64) float64 {
	return math.Round(amount*100) / 100
}

func (s *SettlementSrv) HandleOrderPaid(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	rate, err := dao.NewSettlementDaoByDB(tx).ActiveCommissionRate()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		rate = defaultCommissionRate
	} else if err != nil {
		return err
	}
	calc, err := calculateOrderSettlement(order.Money, order.Num, rate)
	if err != nil {
		return err
	}

	flowDao := dao.NewAccountFlowDaoByDB(tx)
	flows := []*model.AccountFlow{
		newAccountFlow(order, order.UserID, order.BossID, model.AccountFlowTypeBuyerPay, "out", calc.GrossAmount, "买家支付"),
		newAccountFlow(order, order.BossID, order.BossID, model.AccountFlowTypeSellerPending, "in", calc.SettlementAmount, "卖家待结算收入"),
		newAccountFlow(order, 0, order.BossID, model.AccountFlowTypePlatformCommission, "in", calc.CommissionAmount, "平台佣金"),
	}
	for _, flow := range flows {
		if err = flowDao.Create(flow); err != nil {
			return err
		}
	}

	return dao.NewSettlementDaoByDB(tx).CreatePending(&model.Settlement{
		SellerID:         order.BossID,
		OrderID:          order.ID,
		OrderNum:         order.OrderNum,
		GrossAmount:      calc.GrossAmount,
		CommissionRate:   rate,
		CommissionAmount: calc.CommissionAmount,
		SettlementAmount: calc.SettlementAmount,
		Status:           model.SettlementStatusPending,
	})
}

func (s *SettlementSrv) HandleOrderRefunded(ctx context.Context, tx *gorm.DB, order *model.Order) error {
	if err := dao.NewSettlementDaoByDB(tx).MarkRefundedByOrderID(order.ID); err != nil {
		return err
	}
	amount := roundMoney(order.Money * float64(order.Num))
	return dao.NewAccountFlowDaoByDB(tx).Create(
		newAccountFlow(order, order.UserID, order.BossID, model.AccountFlowTypeRefund, "in", amount, "订单退款"),
	)
}

func (s *SettlementSrv) AdminList(ctx context.Context, req *types.AdminSettlementListReq) (interface{}, error) {
	normalizeSettlementPage(req)
	list, total, err := dao.NewSettlementDao(ctx).List(req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp := make([]*types.AdminSettlementResp, 0, len(list))
	for _, settlement := range list {
		resp = append(resp, buildSettlementResp(settlement))
	}
	return &types.DataListResp{Item: resp, Total: total}, nil
}

func (s *SettlementSrv) AdminGenerate(ctx context.Context, req *types.AdminSettlementGenerateReq) (interface{}, error) {
	profile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(req.SellerID)
	if err != nil || profile == nil || !profile.IsApproved() {
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.LogrusObj.Error(err)
			return nil, err
		}
		return nil, e.NewBusinessError(e.ErrorSettlementSellerInvalid)
	}
	count, err := dao.NewSettlementDao(ctx).GenerateCompletedForSeller(req.SellerID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return &types.AdminSettlementGenerateResp{SellerID: req.SellerID, Count: count}, nil
}

func (s *SettlementSrv) AdminGenerateOne(ctx context.Context, req *types.AdminSettlementGenerateOneReq) (interface{}, error) {
	settlement, err := dao.NewSettlementDao(ctx).GenerateCompletedByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSettlementStatusInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSettlementResp(settlement), nil
}

func (s *SettlementSrv) AdminMarkPaid(ctx context.Context, req *types.AdminSettlementPayReq) (interface{}, error) {
	var settlement *model.Settlement
	err := dao.NewSettlementDao(ctx).Transaction(func(tx *gorm.DB) error {
		var txErr error
		settlement, txErr = dao.NewSettlementDaoByDB(tx).MarkPaid(req.ID)
		if txErr != nil {
			return txErr
		}
		return GetSellerAccountSrv().CreditSettlement(tx, settlement)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSettlementStatusInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSettlementResp(settlement), nil
}

func (s *SettlementSrv) AdminDetail(ctx context.Context, req *types.AdminSettlementDetailReq) (interface{}, error) {
	settlement, err := dao.NewSettlementDao(ctx).GetByID(req.ID)
	if err != nil {
		return nil, err
	}
	flows, err := dao.NewAccountFlowDao(ctx).ListByOrderID(settlement.OrderID)
	if err != nil {
		return nil, err
	}
	return &types.DataListResp{Item: flows, Total: int64(len(flows))}, nil
}

func (s *SettlementSrv) SellerSummary(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	profile, err := dao.NewSellerDao(ctx).GetSellerProfileByUserID(u.Id)
	if err != nil || profile == nil || !profile.IsApproved() {
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.LogrusObj.Error(err)
			return nil, err
		}
		return nil, e.NewBusinessError(e.ErrorSellerNotApproved)
	}
	resp, err := dao.NewSettlementDao(ctx).SellerSummary(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return resp, nil
}

func normalizeSettlementPage(req *types.AdminSettlementListReq) {
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
}

func newAccountFlow(order *model.Order, userID, sellerID uint, flowType, direction string, amount float64, remark string) *model.AccountFlow {
	return &model.AccountFlow{
		FlowNo:      fmt.Sprintf("%s-%d-%d", flowType, order.ID, time.Now().UnixNano()),
		OrderID:     order.ID,
		OrderNum:    order.OrderNum,
		UserID:      userID,
		SellerID:    sellerID,
		RelatedType: "order",
		RelatedID:   order.ID,
		FlowType:    flowType,
		Direction:   direction,
		Amount:      amount,
		Remark:      remark,
	}
}

func newAccountFlowFromSettlement(settlement *model.Settlement, flowType, direction string, amount float64, remark string) *model.AccountFlow {
	return &model.AccountFlow{
		FlowNo:      fmt.Sprintf("%s-%d-%d", flowType, settlement.ID, time.Now().UnixNano()),
		OrderID:     settlement.OrderID,
		OrderNum:    settlement.OrderNum,
		UserID:      settlement.SellerID,
		SellerID:    settlement.SellerID,
		RelatedType: "settlement",
		RelatedID:   settlement.ID,
		FlowType:    flowType,
		Direction:   direction,
		Amount:      amount,
		Remark:      remark,
	}
}

func buildSettlementResp(settlement *model.Settlement) *types.AdminSettlementResp {
	var paidAt int64
	if settlement.PaidAt != nil {
		paidAt = *settlement.PaidAt
	}
	return &types.AdminSettlementResp{
		ID:               settlement.ID,
		SellerID:         settlement.SellerID,
		OrderID:          settlement.OrderID,
		OrderNum:         settlement.OrderNum,
		GrossAmount:      settlement.GrossAmount,
		CommissionRate:   settlement.CommissionRate,
		CommissionAmount: settlement.CommissionAmount,
		SettlementAmount: settlement.SettlementAmount,
		Status:           settlement.Status,
		PaidAt:           paidAt,
		CreatedAt:        settlement.CreatedAt.Unix(),
	}
}
