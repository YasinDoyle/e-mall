package dao

import (
	"strings"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestBuildAdminOrderTrendQueryDoesNotFilterByOrderStatus(t *testing.T) {
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
		var rows []AdminOrderTrendRow
		return buildAdminOrderTrendQuery(tx, time.Date(2026, 7, 20, 0, 0, 0, 0, time.Local), time.Date(2026, 7, 27, 0, 0, 0, 0, time.Local)).Scan(&rows)
	})

	if strings.Contains(sql, "type IN") {
		t.Fatalf("trend query should not filter order status, got %s", sql)
	}
	if !strings.Contains(sql, "DATE_FORMAT(paid_at, '%Y-%m-%d') AS date") {
		t.Fatalf("trend query should aggregate by formatted paid_at date, got %s", sql)
	}
}

func TestBuildPlatformRevenueQueryUsesPlatformCommissionFlows(t *testing.T) {
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
		var total float64
		return buildPlatformRevenueQuery(tx).Scan(&total)
	})

	if !strings.Contains(sql, "FROM `account_flows`") {
		t.Fatalf("platform revenue query should read account_flows table, got %s", sql)
	}
	if !strings.Contains(sql, "flow_type = 'platform_commission'") {
		t.Fatalf("platform revenue query should filter platform commission flow, got %s", sql)
	}
	if !strings.Contains(sql, "direction = 'in'") {
		t.Fatalf("platform revenue query should filter inflows, got %s", sql)
	}
}
