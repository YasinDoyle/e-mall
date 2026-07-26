package dao

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/types"
)

func TestBuildListOrderByConditionQueryUsesOrderAliasForSoftDelete(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open dry run db: %v", err)
	}

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var orders []*types.OrderListResp
		return buildListOrderByConditionQuery(tx, 1, &types.OrderListReq{
			BasePage: types.BasePage{PageNum: 1, PageSize: 10},
		}).Find(&orders)
	})

	if strings.Contains(sql, "`order`.`deleted_at`") {
		t.Fatalf("query should not use original table name after alias, got %s", sql)
	}
	if !strings.Contains(sql, "o.deleted_at IS NULL") {
		t.Fatalf("query should filter soft deletes through alias, got %s", sql)
	}
	if !strings.Contains(sql, "FROM `order` AS o") {
		t.Fatalf("query should use explicit order alias, got %s", sql)
	}
	if !strings.Contains(sql, "o.buyer_deleted = false") {
		t.Fatalf("buyer query should hide buyer-deleted orders, got %s", sql)
	}
}

func TestBuildListOrderByBossQueryUsesOrderAliasForSoftDelete(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("open dry run db: %v", err)
	}

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		var orders []*types.OrderListResp
		return buildListOrderByBossQuery(tx, 1, &types.SellerOrderListReq{
			BasePage: types.BasePage{PageNum: 1, PageSize: 10},
		}).Find(&orders)
	})

	if strings.Contains(sql, "`order`.`deleted_at`") {
		t.Fatalf("query should not use original table name after alias, got %s", sql)
	}
	if !strings.Contains(sql, "o.deleted_at IS NULL") {
		t.Fatalf("query should filter soft deletes through alias, got %s", sql)
	}
	if !strings.Contains(sql, "FROM `order` AS o") {
		t.Fatalf("query should use explicit order alias, got %s", sql)
	}
	if !strings.Contains(sql, "o.boss_id = 1") {
		t.Fatalf("query should filter seller orders by boss id, got %s", sql)
	}
	if strings.Contains(sql, "buyer_deleted") {
		t.Fatalf("seller query should not hide orders deleted by buyer, got %s", sql)
	}
}
