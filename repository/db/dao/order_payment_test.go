package dao

import (
	"strings"
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
	"gorm.io/gorm"
)

func TestOrderPaymentModelMigratesWithSingularTable(t *testing.T) {
	db := orderLogDryRunDB(t)

	if err := db.AutoMigrate(&model.OrderPayment{}); err != nil {
		t.Fatalf("auto migrate order payment: %v", err)
	}

	stmt := db.Session(&gorm.Session{DryRun: true}).
		Model(&model.OrderPayment{}).
		Where("payment_no = ?", "OP123").
		Find(&[]model.OrderPayment{}).Statement

	if stmt.Table != "order_payment" {
		t.Fatalf("expected singular table order_payment, got %q", stmt.Table)
	}
}

func TestBuildPendingOrderPaymentByOrderIDQueryScopesByPendingStatus(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildPendingOrderPaymentByOrderIDQuery(db, 42).
		Take(&model.OrderPayment{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "order_id = ?") {
		t.Fatalf("pending payment query should scope by order id, got %s", sql)
	}
	if !strings.Contains(sql, "status = ?") {
		t.Fatalf("pending payment query should scope by status, got %s", sql)
	}
}

func TestBuildOrderPaymentByPaymentNoForUpdateQueryScopesAndLocks(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildOrderPaymentByPaymentNoForUpdateQuery(db, "OP123").
		Take(&model.OrderPayment{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "payment_no = ?") {
		t.Fatalf("for-update payment query should scope by payment_no, got %s", sql)
	}
	if !strings.Contains(strings.ToUpper(sql), "FOR UPDATE") {
		t.Fatalf("for-update payment query should lock selected row, got %s", sql)
	}
}

func TestBuildCloseOtherPendingOrderPaymentsQueryExcludesPaidPayment(t *testing.T) {
	db := orderLogDryRunDB(t)

	stmt := buildCloseOtherPendingOrderPaymentsQuery(db, 42, "OP-paid").
		Find(&[]model.OrderPayment{}).Statement

	sql := stmt.SQL.String()
	if !strings.Contains(sql, "order_id = ?") {
		t.Fatalf("close other pending query should scope by order id, got %s", sql)
	}
	if !strings.Contains(sql, "status = ?") {
		t.Fatalf("close other pending query should scope by pending status, got %s", sql)
	}
	if !strings.Contains(sql, "payment_no <> ?") {
		t.Fatalf("close other pending query should exclude the paid payment, got %s", sql)
	}
}
