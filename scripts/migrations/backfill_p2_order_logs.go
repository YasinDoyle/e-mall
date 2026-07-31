package main

import (
	"context"
	"log"

	config "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const backfillBatchSize = 100
const backfillCurrentStateAction = "backfill_current_state"

func main() {
	config.InitConfig()
	dao.InitMysql()

	ctx := context.Background()
	db := dao.NewOrderDao(ctx).DB

	var inserted int
	var skipped int
	var orders []model.Order

	err := db.Model(&model.Order{}).
		Order("id ASC").
		FindInBatches(&orders, backfillBatchSize, func(tx *gorm.DB, batch int) error {
			for i := range orders {
				order := orders[i]
				if err := tx.Transaction(func(orderTx *gorm.DB) error {
					var lockedOrder model.Order
					if err := orderTx.Clauses(clause.Locking{Strength: "UPDATE"}).
						First(&lockedOrder, order.ID).Error; err != nil {
						return err
					}

					orderLogDao := dao.NewOrderLogDaoByDB(orderTx)
					hasLogs, err := orderLogDao.HasLogsForOrder(lockedOrder.ID)
					if err != nil {
						return err
					}
					if hasLogs {
						skipped++
						return nil
					}

					action := backfillCurrentStateAction
					if lockedOrder.Type == consts.OrderTypeUnPaid {
						action = consts.OrderActionCreate
					}
					if err := orderLogDao.Create(&model.OrderLog{
						OrderID:      lockedOrder.ID,
						OrderNum:     lockedOrder.OrderNum,
						Action:       action,
						FromType:     0,
						ToType:       lockedOrder.Type,
						OperatorType: "system",
						OperatorID:   0,
						Remark:       "backfill order initial state",
					}); err != nil {
						return err
					}
					inserted++
					return nil
				}); err != nil {
					return err
				}
			}
			return nil
		}).Error
	if err != nil {
		log.Fatalf("backfill P2 order logs failed: %v", err)
	}

	log.Printf("backfill P2 order logs complete: inserted=%d skipped=%d", inserted, skipped)
}
