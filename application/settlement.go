package application

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/domain/event"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
	"github.com/YasinDoyle/e-mall/utils/log"
)

type SettlementUsecase struct{}

func NewSettlementUsecase() *SettlementUsecase {
	return &SettlementUsecase{}
}

func (u *SettlementUsecase) AdminMarkPaid(ctx context.Context, req *types.AdminSettlementPayReq) (interface{}, error) {
	var settlement *model.Settlement
	err := dao.NewSettlementDao(ctx).Transaction(func(tx *gorm.DB) error {
		var txErr error
		settlement, txErr = dao.NewSettlementDaoByDB(tx).MarkPaid(req.ID)
		if txErr != nil {
			return txErr
		}
		return creditSettlement(tx, settlement)
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.NewBusinessError(e.ErrorSettlementStatusInvalid)
		}
		log.LogrusObj.Error(err)
		return nil, err
	}
	event.Publish(ctx, event.SettlementPaid{Settlement: settlement})
	return buildSettlementResp(settlement), nil
}
