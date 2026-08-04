package dao

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestAfterSaleModelMigratesWithSingularTable(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		t.Fatalf("open dry-run db: %v", err)
	}

	if err := db.AutoMigrate(&model.AfterSale{}); err != nil {
		t.Fatalf("auto migrate after sale: %v", err)
	}

	stmt := db.Session(&gorm.Session{DryRun: true}).
		Create(&model.AfterSale{
			OrderID:      42,
			OrderNum:     1001,
			BuyerID:      7,
			SellerID:     8,
			Type:         consts.AfterSaleTypeRefundOnly,
			Status:       consts.AfterSaleStatusRequested,
			Reason:       "item damaged",
			RefundAmount: 88.5,
		}).Statement

	if stmt.Table != "after_sale" {
		t.Fatalf("expected singular table after_sale, got %q", stmt.Table)
	}
}
