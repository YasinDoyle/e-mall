package dao

import (
	"strings"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/types"
)

func TestBuildListFavoriteByUserIdQueryUsesFavoriteAliasForSoftDelete(t *testing.T) {
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
		var favorites []*types.FavoriteListResp
		return buildListFavoriteByUserIdQuery(tx, 1, 10, 1).Find(&favorites)
	})

	if strings.Contains(sql, "`favorite`.`deleted_at`") {
		t.Fatalf("query should not use original table name after alias, got %s", sql)
	}
	if !strings.Contains(sql, "f.deleted_at IS NULL") {
		t.Fatalf("query should filter soft deletes through alias, got %s", sql)
	}
}
