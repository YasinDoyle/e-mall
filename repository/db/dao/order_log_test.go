package dao

import (
	"strings"
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func TestOrderLogModelMigratesWithSingularTable(t *testing.T) {
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

	if err := db.AutoMigrate(&model.OrderLog{}); err != nil {
		t.Fatalf("auto migrate order log: %v", err)
	}

	stmt := db.Session(&gorm.Session{DryRun: true}).
		Model(&model.OrderLog{}).
		Where("order_id = ?", uint(42)).
		Order("created_at ASC").
		Order("id ASC").
		Find(&[]model.OrderLog{}).Statement

	if stmt.Table != "order_log" {
		t.Fatalf("expected singular table order_log, got %q", stmt.Table)
	}
}

func TestBuildOrderLogListByOrderIDQueryOrdersChronologically(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildOrderLogListByOrderIDQuery(db, 42).
		Find(&[]model.OrderLog{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "order_id = ?") {
		t.Fatalf("order log list should scope by order id, got %s", sql)
	}
	if !strings.Contains(sql, "ORDER BY created_at ASC,id ASC") &&
		!strings.Contains(sql, "ORDER BY created_at ASC, id ASC") {
		t.Fatalf("order log list should be chronological and stable, got %s", sql)
	}
}

func TestBuildOrderLogHasLogsForOrderQueryUsesExistenceProbe(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildOrderLogHasLogsForOrderQuery(db, 42).
		Take(&model.OrderLog{}).Statement

	sql := stmt.SQL.String()
	if strings.Contains(strings.ToUpper(sql), "COUNT") {
		t.Fatalf("has logs query should use an existence probe, got %s", sql)
	}
	if !strings.Contains(sql, "order_id = ?") {
		t.Fatalf("has logs query should scope by order id, got %s", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "LIMIT") || len(stmt.Vars) == 0 || stmt.Vars[len(stmt.Vars)-1] != 1 {
		t.Fatalf("has logs query should limit to one row, got sql=%s vars=%v", sql, stmt.Vars)
	}
}

func orderLogDryRunDB(t *testing.T) *gorm.DB {
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
