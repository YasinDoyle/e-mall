package dao

import (
	"strings"
	"testing"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestBuildGenerateCompletedForOrderQueryExcludesActiveAfterSales(t *testing.T) {
	db := orderLogDryRunDB(t)

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return buildGenerateCompletedForOrderQuery(tx, 42).
			Update("status", model.SettlementStatusGenerated)
	})
	if !strings.Contains(sql, "order_id = 42") {
		t.Fatalf("generate order settlement should scope by order id, got %s", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "NOT EXISTS") || !strings.Contains(sql, "after_sale") {
		t.Fatalf("generate order settlement should exclude active after-sales, got %s", sql)
	}
	if !strings.Contains(sql, "status NOT IN") {
		t.Fatalf("generate order settlement should keep terminal after-sales settleable, got %s", sql)
	}
}

func TestBuildSettlementReadyOrderQueryExcludesActiveAfterSales(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildSettlementReadyOrderQuery(db, 42, 7).
		First(&model.Order{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "o.id = ?") || !strings.Contains(sql, "o.boss_id = ?") {
		t.Fatalf("ready order query should scope by order and seller, got %s", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "NOT EXISTS") || !strings.Contains(sql, "after_sale") {
		t.Fatalf("ready order query should exclude active after-sales, got %s", sql)
	}
}

func TestBuildSettlementReadyOrderQueryRequiresReceivedOrder(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildSettlementReadyOrderQuery(db, 42, 7).
		First(&model.Order{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "o.type = ?") {
		t.Fatalf("ready order query must require received order status, got %s", sql)
	}
}

func TestBuildSettlementMarkPaidReadyOrderQueryScopesSettlementOrder(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildSettlementMarkPaidReadyOrderQuery(db, &model.Settlement{OrderID: 42, SellerID: 7}).
		First(&model.Order{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "o.id = ?") || !strings.Contains(sql, "o.boss_id = ?") {
		t.Fatalf("mark-paid ready query must scope by settlement order and seller, got %s", sql)
	}
	if !strings.Contains(sql, "o.type = ?") || !strings.Contains(strings.ToUpper(sql), "NOT EXISTS") {
		t.Fatalf("mark-paid ready query must require received order and no active after-sale, got %s", sql)
	}
}

func TestBuildGenerateCompletedForSellerOrderSubqueryExcludesActiveAfterSales(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildGenerateCompletedForSellerOrderSubquery(db, 7).
		Find(&[]struct{}{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "o.boss_id = ?") {
		t.Fatalf("seller order subquery should scope by seller, got %s", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "NOT EXISTS") || !strings.Contains(sql, "after_sale") {
		t.Fatalf("seller order subquery should exclude active after-sales, got %s", sql)
	}
}

func TestBuildMarkRefundedByOrderIDQueryIncludesPaidSettlements(t *testing.T) {
	db := orderLogDryRunDB(t)

	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return buildMarkRefundedByOrderIDQuery(tx, 42).
			Update("status", model.SettlementStatusRefunded)
	})
	if !strings.Contains(sql, "'paid'") {
		t.Fatalf("refunded settlement query must include paid settlements, got %s", sql)
	}
}
