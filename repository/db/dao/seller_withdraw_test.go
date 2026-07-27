package dao

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/types"
)

func TestBuildListSellerWithdrawBySellerIDQueryFiltersBySellerAndStatus(t *testing.T) {
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
		var list []*types.SellerWithdrawResp
		return buildListSellerWithdrawQuery(tx, &types.SellerWithdrawListReq{
			BasePage: types.BasePage{PageNum: 1, PageSize: 10},
			Status:   "approved",
			SellerID: 42,
		}).Find(&list)
	})

	if !strings.Contains(sql, "FROM seller_withdraw AS sw") {
		t.Fatalf("expected explicit seller_withdraw alias, got %s", sql)
	}
	if !strings.Contains(sql, "sw.seller_id = 42") {
		t.Fatalf("expected seller filter, got %s", sql)
	}
	if !strings.Contains(sql, "sw.status = 'approved'") {
		t.Fatalf("expected status filter, got %s", sql)
	}
}
