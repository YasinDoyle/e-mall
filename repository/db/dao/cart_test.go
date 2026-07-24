package dao

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/types"
)

func TestBuildListCartByUserIdQueryUsesCartAliasForSoftDelete(t *testing.T) {
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
		var carts []*types.CartResp
		return buildListCartByUserIdQuery(tx, 1).Find(&carts)
	})

	if strings.Contains(sql, "`cart`.`deleted_at`") {
		t.Fatalf("query should not use original table name after alias, got %s", sql)
	}
	if !strings.Contains(sql, "c.deleted_at IS NULL") {
		t.Fatalf("query should filter soft deletes through alias, got %s", sql)
	}
}
