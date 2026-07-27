package dao

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestBuildSellerAccountSummaryQueryUsesSellerAccountTable(t *testing.T) {
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
		var account struct{}
		return buildSellerAccountSummaryQuery(tx, 42).First(&account)
	})

	if !strings.Contains(sql, "FROM `seller_account`") {
		t.Fatalf("seller account summary should read seller_account table, got %s", sql)
	}
	if !strings.Contains(sql, "seller_id = 42") {
		t.Fatalf("seller account summary should filter by seller_id, got %s", sql)
	}
}
