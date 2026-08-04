package dao

import (
	"strings"
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestOrderLogisticsModelMigratesWithSingularTable(t *testing.T) {
	db := orderLogisticsDryRunDB(t)

	stmt := db.Session(&gorm.Session{DryRun: true}).
		Create(&model.OrderLogistics{
			OrderID:     1,
			OrderNum:    1001,
			NodeType:    "manual_shipped",
			Description: "SF SF123",
			OccurredAt:  123456,
		}).Statement

	if stmt.Table != "order_logistics" {
		t.Fatalf("expected order_logistics table, got %q", stmt.Table)
	}
}

func TestOrderLogisticsListByOrderIDOrdersChronologically(t *testing.T) {
	db := orderLogisticsDryRunDB(t)

	stmt := buildOrderLogisticsListByOrderIDQuery(db, 42).
		Find(&[]model.OrderLogistics{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "order_id = ?") {
		t.Fatalf("order logistics list should scope by order id, got %s", sql)
	}
	if !strings.Contains(sql, "ORDER BY occurred_at ASC,id ASC") &&
		!strings.Contains(sql, "ORDER BY occurred_at ASC, id ASC") {
		t.Fatalf("order logistics list should be chronological and stable, got %s", sql)
	}
}

func orderLogisticsDryRunDB(t *testing.T) *gorm.DB {
	t.Helper()

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

	return db
}
