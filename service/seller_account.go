package service

import (
	"context"
	"errors"
	"strings"
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

type SellerAccountSrv struct{}

var sellerAccountSrvIns *SellerAccountSrv
var sellerAccountSrvOnce sync.Once

func GetSellerAccountSrv() *SellerAccountSrv {
	sellerAccountSrvOnce.Do(func() {
		sellerAccountSrvIns = &SellerAccountSrv{}
	})
	return sellerAccountSrvIns
}

func (s *SellerAccountSrv) Summary(ctx context.Context) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
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
	account, err := dao.NewSellerAccountDao(ctx).GetOrCreateBySellerID(profile.UserID)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSellerAccountSummaryResp(account), nil
}

func (s *SellerAccountSrv) WithdrawList(ctx context.Context, req *types.SellerWithdrawListReq) (interface{}, error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
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
	if req == nil {
		req = &types.SellerWithdrawListReq{}
	}
	req.SellerID = profile.UserID
	normalizeSellerWithdrawPage(req)
	list, total, err := dao.NewSellerWithdrawDao(ctx).List(req)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp := make([]*types.SellerWithdrawResp, 0, len(list))
	for _, item := range list {
		resp = append(resp, buildSellerWithdrawResp(item))
	}
	return &types.DataListResp{Item: resp, Total: total}, nil
}

func (s *SellerAccountSrv) WithdrawApply(ctx context.Context, req *types.SellerWithdrawApplyReq) (interface{}, error) {
	if err := validateSellerWithdrawApplyReq(req); err != nil {
		return nil, err
	}
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
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

	var withdraw *model.SellerWithdraw
	err = dao.NewSellerWithdrawDao(ctx).DB.Transaction(func(tx *gorm.DB) error {
		accountDao := dao.NewSellerAccountDaoByDB(tx)
		account, txErr := accountDao.GetOrCreateBySellerID(profile.UserID)
		if txErr != nil {
			return txErr
		}
		if txErr = applySellerWithdrawFreeze(account, req.Amount); txErr != nil {
			return txErr
		}
		if txErr = accountDao.Save(account); txErr != nil {
			return txErr
		}

		withdraw = &model.SellerWithdraw{
			SellerID:     profile.UserID,
			Amount:       roundMoney(req.Amount),
			Status:       model.SellerWithdrawStatusPending,
			PayeeName:    strings.TrimSpace(req.PayeeName),
			PayeeAccount: strings.TrimSpace(req.PayeeAccount),
			PayeeChannel: defaultSellerWithdrawChannel(req.PayeeChannel),
		}
		if txErr = dao.NewSellerWithdrawDaoByDB(tx).Create(withdraw); txErr != nil {
			return txErr
		}
		flow := newSellerRelatedAccountFlow(profile.UserID, "withdraw", withdraw.ID, model.AccountFlowTypeSellerWithdrawFreeze, "out", withdraw.Amount, "提现冻结")
		return dao.NewAccountFlowDaoByDB(tx).Create(flow)
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSellerWithdrawResp(withdraw), nil
}

func (s *SellerAccountSrv) AdminWithdrawList(ctx context.Context, req *types.AdminSellerWithdrawListReq) (interface{}, error) {
	if req == nil {
		req = &types.AdminSellerWithdrawListReq{}
	}
	normalizeSellerWithdrawPage((*types.SellerWithdrawListReq)(req))
	list, total, err := dao.NewSellerWithdrawDao(ctx).List((*types.SellerWithdrawListReq)(req))
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	resp := make([]*types.SellerWithdrawResp, 0, len(list))
	for _, item := range list {
		resp = append(resp, buildSellerWithdrawResp(item))
	}
	return &types.DataListResp{Item: resp, Total: total}, nil
}

func (s *SellerAccountSrv) AdminWithdrawAudit(ctx context.Context, req *types.AdminSellerWithdrawAuditReq) (interface{}, error) {
	if err := validateAdminSellerWithdrawAuditReq(req); err != nil {
		return nil, err
	}
	adminUser, err := dao.NewUserDao(ctx).GetUserById(userIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	if req.Status == model.SellerWithdrawStatusRejected && strings.TrimSpace(req.Reason) == "" {
		return nil, e.NewBusinessError(e.ErrorSellerWithdrawReasonMissing)
	}

	var withdraw *model.SellerWithdraw
	err = dao.NewSellerWithdrawDao(ctx).DB.Transaction(func(tx *gorm.DB) error {
		withdrawDao := dao.NewSellerWithdrawDaoByDB(tx)
		accountDao := dao.NewSellerAccountDaoByDB(tx)
		withdraw, err = withdrawDao.GetByID(req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return e.NewBusinessError(e.ErrorSellerWithdrawNotFound)
			}
			return err
		}
		if !withdraw.IsPending() {
			return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
		}
		now := time.Now()
		updates := map[string]interface{}{
			"audited_at":          &now,
			"audit_operator_id":   adminUser.ID,
			"audit_operator_name": adminUser.UserName,
		}
		if req.Status == model.SellerWithdrawStatusApproved {
			updates["status"] = model.SellerWithdrawStatusApproved
			withdraw.Status = model.SellerWithdrawStatusApproved
		} else {
			account, txErr := accountDao.GetOrCreateBySellerID(withdraw.SellerID)
			if txErr != nil {
				return txErr
			}
			if txErr = applySellerWithdrawRejected(account, withdraw.Amount); txErr != nil {
				return txErr
			}
			if txErr = accountDao.Save(account); txErr != nil {
				return txErr
			}
			updates["status"] = model.SellerWithdrawStatusRejected
			updates["audit_reason"] = strings.TrimSpace(req.Reason)
			withdraw.Status = model.SellerWithdrawStatusRejected
			withdraw.AuditReason = strings.TrimSpace(req.Reason)
			if txErr := dao.NewAccountFlowDaoByDB(tx).Create(
				newSellerRelatedAccountFlow(withdraw.SellerID, "withdraw", withdraw.ID, model.AccountFlowTypeSellerWithdrawUnfreeze, "in", withdraw.Amount, "提现拒绝解冻"),
			); txErr != nil {
				return txErr
			}
		}
		withdraw.AuditedAt = &now
		withdraw.AuditOperatorID = adminUser.ID
		withdraw.AuditOperatorName = adminUser.UserName
		return tx.Model(&model.SellerWithdraw{}).
			Where("id = ? AND status = ?", withdraw.ID, model.SellerWithdrawStatusPending).
			Updates(updates).Error
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSellerWithdrawResp(withdraw), nil
}

func (s *SellerAccountSrv) AdminWithdrawPaid(ctx context.Context, req *types.AdminSellerWithdrawPaidReq) (interface{}, error) {
	if err := validateAdminSellerWithdrawPaidReq(req); err != nil {
		return nil, err
	}
	adminUser, err := dao.NewUserDao(ctx).GetUserById(userIDFromContext(ctx))
	if err != nil {
		return nil, err
	}
	var withdraw *model.SellerWithdraw
	err = dao.NewSellerWithdrawDao(ctx).DB.Transaction(func(tx *gorm.DB) error {
		withdrawDao := dao.NewSellerWithdrawDaoByDB(tx)
		accountDao := dao.NewSellerAccountDaoByDB(tx)
		withdraw, err = withdrawDao.GetByID(req.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return e.NewBusinessError(e.ErrorSellerWithdrawNotFound)
			}
			return err
		}
		if withdraw.Status != model.SellerWithdrawStatusApproved {
			return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
		}
		now := time.Now()
		updates := map[string]interface{}{
			"paid_at":            &now,
			"paid_operator_id":   adminUser.ID,
			"paid_operator_name": adminUser.UserName,
		}
		account, txErr := accountDao.GetOrCreateBySellerID(withdraw.SellerID)
		if txErr != nil {
			return txErr
		}
		if req.Status == model.SellerWithdrawStatusPaid {
			if txErr = applySellerWithdrawPaid(account, withdraw.Amount); txErr != nil {
				return txErr
			}
			updates["status"] = model.SellerWithdrawStatusPaid
			withdraw.Status = model.SellerWithdrawStatusPaid
			if txErr = accountDao.Save(account); txErr != nil {
				return txErr
			}
			if txErr = dao.NewAccountFlowDaoByDB(tx).Create(
				newSellerRelatedAccountFlow(withdraw.SellerID, "withdraw", withdraw.ID, model.AccountFlowTypeSellerWithdrawPaid, "out", withdraw.Amount, "提现打款"),
			); txErr != nil {
				return txErr
			}
		} else if req.Status == model.SellerWithdrawStatusFailed {
			if txErr = applySellerWithdrawRejected(account, withdraw.Amount); txErr != nil {
				return txErr
			}
			updates["status"] = model.SellerWithdrawStatusFailed
			updates["audit_reason"] = strings.TrimSpace(req.Reason)
			withdraw.Status = model.SellerWithdrawStatusFailed
			withdraw.AuditReason = strings.TrimSpace(req.Reason)
			if txErr = accountDao.Save(account); txErr != nil {
				return txErr
			}
			if txErr = dao.NewAccountFlowDaoByDB(tx).Create(
				newSellerRelatedAccountFlow(withdraw.SellerID, "withdraw", withdraw.ID, model.AccountFlowTypeSellerWithdrawUnfreeze, "in", withdraw.Amount, "提现失败解冻"),
			); txErr != nil {
				return txErr
			}
		} else {
			return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
		}
		withdraw.PaidAt = &now
		withdraw.PaidOperatorID = adminUser.ID
		withdraw.PaidOperatorName = adminUser.UserName
		return tx.Model(&model.SellerWithdraw{}).
			Where("id = ? AND status = ?", withdraw.ID, model.SellerWithdrawStatusApproved).
			Updates(updates).Error
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return buildSellerWithdrawResp(withdraw), nil
}

func (s *SellerAccountSrv) AdminWithdrawDetail(ctx context.Context, req *types.AdminIDReq) (interface{}, error) {
	withdraw, err := dao.NewSellerWithdrawDao(ctx).GetByID(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSellerWithdrawNotFound)
		}
		return nil, err
	}
	flows, err := dao.NewAccountFlowDao(ctx).ListByRelatedTypeAndID("withdraw", withdraw.ID)
	if err != nil {
		return nil, err
	}
	return &types.DataListResp{Item: flows, Total: int64(len(flows))}, nil
}

func (s *SellerAccountSrv) CreditSettlement(tx *gorm.DB, settlement *model.Settlement) error {
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

func (s *SellerAccountSrv) AdminBackfillPaidSettlementCredits(ctx context.Context) (interface{}, error) {
	settlements, err := dao.NewSettlementDao(ctx).ListPaidSettlements()
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	creditedIDs, err := dao.NewAccountFlowDao(ctx).ListRelatedIDsByTypeAndFlowType(
		"settlement",
		model.AccountFlowTypeSellerSettlementCredit,
	)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	creditedSet := make(map[uint]struct{}, len(creditedIDs))
	for _, id := range creditedIDs {
		creditedSet[id] = struct{}{}
	}
	plan := buildSellerSettlementCreditBackfillPlan(settlements, creditedSet)
	if len(plan) == 0 {
		return &types.AdminSellerAccountBackfillResp{}, nil
	}

	resp := &types.AdminSellerAccountBackfillResp{}
	err = dao.NewSellerAccountDao(ctx).DB.Transaction(func(tx *gorm.DB) error {
		accountDao := dao.NewSellerAccountDaoByDB(tx)
		flowDao := dao.NewAccountFlowDaoByDB(tx)
		for sellerID, missingSettlements := range plan {
			account, txErr := accountDao.GetOrCreateBySellerID(sellerID)
			if txErr != nil {
				return txErr
			}
			var sellerAmount float64
			for _, settlement := range missingSettlements {
				sellerAmount += settlement.SettlementAmount
			}
			if txErr = applySellerSettlementCredit(account, sellerAmount); txErr != nil {
				return txErr
			}
			if txErr = accountDao.Save(account); txErr != nil {
				return txErr
			}
			for _, settlement := range missingSettlements {
				if txErr = flowDao.Create(
					newSellerRelatedAccountFlow(
						settlement.SellerID,
						"settlement",
						settlement.ID,
						model.AccountFlowTypeSellerSettlementCredit,
						"in",
						settlement.SettlementAmount,
						"结算回填",
					),
				); txErr != nil {
					return txErr
				}
			}
			resp.SellerCount++
			resp.SettlementCount += int64(len(missingSettlements))
			resp.Amount = roundMoney(resp.Amount + sellerAmount)
		}
		return nil
	})
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	return resp, nil
}

func validateSellerWithdrawApplyReq(req *types.SellerWithdrawApplyReq) error {
	if req == nil {
		return e.NewBusinessError(e.ErrorSellerWithdrawAmountInvalid)
	}
	if req.Amount <= 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawAmountInvalid)
	}
	if strings.TrimSpace(req.PayeeName) == "" || strings.TrimSpace(req.PayeeAccount) == "" {
		return e.NewBusinessError(e.ErrorSellerWithdrawPayeeRequired)
	}
	return nil
}

func validateAdminSellerWithdrawAuditReq(req *types.AdminSellerWithdrawAuditReq) error {
	if req == nil || req.ID == 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	switch req.Status {
	case model.SellerWithdrawStatusApproved, model.SellerWithdrawStatusRejected:
	default:
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	if req.Status == model.SellerWithdrawStatusRejected && strings.TrimSpace(req.Reason) == "" {
		return e.NewBusinessError(e.ErrorSellerWithdrawReasonMissing)
	}
	return nil
}

func validateAdminSellerWithdrawPaidReq(req *types.AdminSellerWithdrawPaidReq) error {
	if req == nil || req.ID == 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	switch req.Status {
	case model.SellerWithdrawStatusPaid, model.SellerWithdrawStatusFailed:
	default:
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	if req.Status == model.SellerWithdrawStatusFailed && strings.TrimSpace(req.Reason) == "" {
		return e.NewBusinessError(e.ErrorSellerWithdrawReasonMissing)
	}
	return nil
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

func applySellerWithdrawFreeze(account *model.SellerAccount, amount float64) error {
	if account == nil || amount <= 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawAmountInvalid)
	}
	if account.AvailableBalance < amount {
		return e.NewBusinessError(e.ErrorSellerWithdrawInsufficientBalance)
	}
	account.AvailableBalance = roundMoney(account.AvailableBalance - amount)
	account.FrozenBalance = roundMoney(account.FrozenBalance + amount)
	return nil
}

func applySellerWithdrawPaid(account *model.SellerAccount, amount float64) error {
	if account == nil || amount <= 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	if account.FrozenBalance < amount {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	account.FrozenBalance = roundMoney(account.FrozenBalance - amount)
	account.TotalWithdrawn = roundMoney(account.TotalWithdrawn + amount)
	return nil
}

func applySellerWithdrawRejected(account *model.SellerAccount, amount float64) error {
	if account == nil || amount <= 0 {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	if account.FrozenBalance < amount {
		return e.NewBusinessError(e.ErrorSellerWithdrawStatusInvalid)
	}
	account.AvailableBalance = roundMoney(account.AvailableBalance + amount)
	account.FrozenBalance = roundMoney(account.FrozenBalance - amount)
	return nil
}

func buildSellerAccountSummaryResp(account *model.SellerAccount) *types.SellerAccountSummaryResp {
	if account == nil {
		return &types.SellerAccountSummaryResp{}
	}
	return &types.SellerAccountSummaryResp{
		SellerID:         account.SellerID,
		AvailableBalance: account.AvailableBalance,
		FrozenBalance:    account.FrozenBalance,
		TotalIncome:      account.TotalIncome,
		TotalWithdrawn:   account.TotalWithdrawn,
	}
}

func buildSellerWithdrawResp(withdraw *model.SellerWithdraw) *types.SellerWithdrawResp {
	if withdraw == nil {
		return nil
	}
	var auditedAt, paidAt int64
	if withdraw.AuditedAt != nil {
		auditedAt = withdraw.AuditedAt.Unix()
	}
	if withdraw.PaidAt != nil {
		paidAt = withdraw.PaidAt.Unix()
	}
	return &types.SellerWithdrawResp{
		ID:                withdraw.ID,
		SellerID:          withdraw.SellerID,
		UserName:          withdraw.Seller.User.UserName,
		NickName:          withdraw.Seller.User.NickName,
		ShopName:          withdraw.Seller.ShopName,
		Amount:            withdraw.Amount,
		Status:            withdraw.Status,
		StatusText:        sellerWithdrawStatusText(withdraw.Status),
		PayeeName:         withdraw.PayeeName,
		PayeeAccount:      withdraw.PayeeAccount,
		PayeeChannel:      withdraw.PayeeChannel,
		AuditReason:       withdraw.AuditReason,
		AuditOperatorID:   withdraw.AuditOperatorID,
		AuditOperatorName: withdraw.AuditOperatorName,
		PaidOperatorID:    withdraw.PaidOperatorID,
		PaidOperatorName:  withdraw.PaidOperatorName,
		CreatedAt:         withdraw.CreatedAt.Unix(),
		AuditedAt:         auditedAt,
		PaidAt:            paidAt,
	}
}

func buildSellerSettlementCreditBackfillPlan(settlements []*model.Settlement, creditedIDs map[uint]struct{}) map[uint][]*model.Settlement {
	plan := make(map[uint][]*model.Settlement)
	for _, settlement := range settlements {
		if settlement == nil || settlement.Status != model.SettlementStatusPaid {
			continue
		}
		if _, ok := creditedIDs[settlement.ID]; ok {
			continue
		}
		plan[settlement.SellerID] = append(plan[settlement.SellerID], settlement)
	}
	return plan
}

func sellerWithdrawStatusText(status string) string {
	switch status {
	case model.SellerWithdrawStatusPending:
		return "待审核"
	case model.SellerWithdrawStatusApproved:
		return "已通过"
	case model.SellerWithdrawStatusRejected:
		return "已拒绝"
	case model.SellerWithdrawStatusPaid:
		return "已打款"
	case model.SellerWithdrawStatusFailed:
		return "打款失败"
	default:
		return "未知"
	}
}

func normalizeSellerWithdrawPage(req *types.SellerWithdrawListReq) {
	if req.PageSize == 0 {
		req.PageSize = consts.BasePageSize
	}
	if req.PageNum == 0 {
		req.PageNum = 1
	}
}

func defaultSellerWithdrawChannel(channel string) string {
	channel = strings.TrimSpace(channel)
	if channel == "" {
		return "manual"
	}
	return channel
}

func userIDFromContext(ctx context.Context) uint {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return 0
	}
	return u.Id
}
