package application

import (
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
)

const defaultCommissionRate = 0.05

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

func handleOrderPaid(tx *gorm.DB, order *model.Order) error {
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
		newOrderAccountFlow(order, order.UserID, order.BossID, model.AccountFlowTypeBuyerPay, "out", calc.GrossAmount, "买家支付"),
		newOrderAccountFlow(order, order.BossID, order.BossID, model.AccountFlowTypeSellerPending, "in", calc.SettlementAmount, "卖家待结算收入"),
		newOrderAccountFlow(order, 0, order.BossID, model.AccountFlowTypePlatformCommission, "in", calc.CommissionAmount, "平台佣金"),
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

func handleOrderRefunded(tx *gorm.DB, order *model.Order) error {
	settlementDao := dao.NewSettlementDaoByDB(tx)
	settlement, err := settlementDao.GetByOrderIDForUpdate(order.ID)
	if err != nil {
		return err
	}
	if settlement.Status == model.SettlementStatusRefunded {
		return e.NewBusinessError(e.ErrorRefundStatusInvalid)
	}
	if settlement.Status != model.SettlementStatusPending && settlement.Status != model.SettlementStatusGenerated {
		return e.NewBusinessError(e.ErrorRefundStatusInvalid)
	}

	flows := buildRefundAccountFlows(order, roundMoney(order.Money*float64(order.Num)))
	if len(flows) != 3 {
		return e.NewBusinessError(e.ErrorDatabase)
	}
	if settlement.SettlementAmount > 0 {
		flows[1].Amount = settlement.SettlementAmount
		flows[1].Remark = "卖家待结算冲正"
	}
	if settlement.CommissionAmount > 0 {
		flows[2].Amount = settlement.CommissionAmount
		flows[2].Remark = "平台佣金冲正"
	}
	flowDao := dao.NewAccountFlowDaoByDB(tx)
	for _, flow := range flows {
		if err := flowDao.Create(flow); err != nil {
			return err
		}
	}
	if err := settlementDao.MarkRefundedByOrderID(order.ID); err != nil {
		return err
	}
	return nil
}

func creditSettlement(tx *gorm.DB, settlement *model.Settlement) error {
	if tx == nil || settlement == nil {
		return nil
	}
	accountDao := dao.NewSellerAccountDaoByDB(tx)
	account, err := accountDao.GetOrCreateBySellerID(settlement.SellerID)
	if err != nil {
		return err
	}
	if err = applySellerSettlementCredit(account, settlement.SettlementAmount); err != nil {
		return err
	}
	if err = accountDao.Save(account); err != nil {
		return err
	}
	return dao.NewAccountFlowDaoByDB(tx).Create(
		newSellerRelatedAccountFlow(settlement.SellerID, "settlement", settlement.ID, model.AccountFlowTypeSellerSettlementCredit, "in", settlement.SettlementAmount, "结算入账"),
	)
}

func applySellerSettlementCredit(account *model.SellerAccount, amount float64) error {
	if account == nil {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	if amount <= 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawAmountInvalid)
	}
	account.AvailableBalance = roundMoney(account.AvailableBalance + amount)
	account.TotalIncome = roundMoney(account.TotalIncome + amount)
	return nil
}

func ensureNotBuyingOwnProduct(userID, bossID uint) error {
	if userID != 0 && userID == bossID {
		return e.NewBusinessError(e.ErrorOrderSelfPurchaseForbidden)
	}
	return nil
}

func newOrderAccountFlow(order *model.Order, userID, sellerID uint, flowType, direction string, amount float64, remark string) *model.AccountFlow {
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

func buildRefundAccountFlows(order *model.Order, amount float64) []*model.AccountFlow {
	return []*model.AccountFlow{
		newStableOrderAccountFlow(order, order.UserID, order.BossID, model.AccountFlowTypeRefund, "in", amount, "订单退款", "buyer"),
		newStableOrderAccountFlow(order, order.BossID, order.BossID, model.AccountFlowTypeSellerPending, "out", 0, "卖家待结算冲正", "seller"),
		newStableOrderAccountFlow(order, 0, order.BossID, model.AccountFlowTypePlatformCommission, "out", 0, "平台佣金冲正", "platform"),
	}
}

func newStableOrderAccountFlow(order *model.Order, userID, sellerID uint, flowType, direction string, amount float64, remark, suffix string) *model.AccountFlow {
	return &model.AccountFlow{
		FlowNo:      fmt.Sprintf("%s-%d-%s", flowType, order.ID, suffix),
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

func newSellerRelatedAccountFlow(sellerID uint, relatedType string, relatedID uint, flowType, direction string, amount float64, remark string) *model.AccountFlow {
	return &model.AccountFlow{
		FlowNo:      fmt.Sprintf("%s-%d-%d", flowType, relatedID, time.Now().UnixNano()),
		UserID:      sellerID,
		SellerID:    sellerID,
		RelatedType: relatedType,
		RelatedID:   relatedID,
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
