package service

import (
	"fmt"
	"time"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

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
