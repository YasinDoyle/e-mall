# P2 Order Fulfillment Loop Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the P2 fulfillment loop so orders move through payment, shipment, receipt, after-sale, refund reversal, and settlement with explicit state transitions and user-visible audit trails.

**Architecture:** Keep the existing `order.type` database column for compatibility with current APIs, but stop scattering raw numeric checks by introducing pure state transition helpers and application-layer use cases. Add operation logs, structured local logistics nodes, order payment records, and independent after-sale records; strong-consistency money and order transitions stay in `application`, while notifications remain `domain/event` side effects after transactions commit.

**Tech Stack:** Go, Gin, GORM, MySQL, Vue 3, TypeScript, Element Plus, Pinia, Vue I18n.

---

## Done Definition

**已有能力：**

- 余额支付 can move an unpaid order to paid, deduct buyer balance, decrease stock, emit `OrderPaid`, and create P1 settlement/account flow records.
- Seller center can ship an order with a logistics tracking number.
- Buyer order detail can confirm receipt and display a simple local timeline.
- Buyer can request refund; admin can approve refund; refund approval credits buyer balance and marks P1 settlement refunded.
- Notification and application/domain-event boundaries already exist.

**目标能力：**

- All order status changes go through one state machine and create immutable operation logs.
- One normal order can complete: create -> pay -> seller ship -> buyer receive -> settlement generate/pay.
- One after-sale order can complete: pay -> buyer request after-sale -> seller handle or platform intervene -> refund -> account-flow and settlement reversal.
- Buyer, seller, and admin UIs expose the relevant state, action, operation log, local logistics nodes, after-sale status, and refund result.
- Product orders support `balance`, `wechat`, and `alipay` payment channels through order payment records; gateway callbacks validate amount/status and are idempotent before emitting `OrderPaid`.

**非目标：**

- Real logistics provider tracking API is not part of P2; it is tracked in `.github/prompts/plan-eMallFullDevelopmentPlan.prompt.md` as `P2-Follow-up`.
- Risk scoring, risk-order backend, automatic payment-channel reconciliation, reliable outbox, dead-letter queues, and full observability are A1/A4 scope.
- Historical data is not silently normalized in runtime business code. Old orders get a one-time backfill script for operation logs.

**验收标准：**

- 正向主流程：buyer A creates an order, pays with balance, seller ships with company + tracking number, buyer sees timeline and confirms receipt, admin generates and pays settlement, seller account is credited.
- 隔离/边界流程：seller B cannot ship seller C's order; buyer A cannot view buyer B's logs or after-sale records; duplicate payment callback or duplicate refund approval does not duplicate account flow, stock deduction, refund, or settlement reversal.
- 回归流程：existing orders without operation logs remain queryable; after running the one-time backfill script, each existing order has at least an initial/current-state log; no runtime code reads old hidden localStorage keys or silently migrates old order data.

---

## File Structure

- `consts/order.go`: semantic order status aliases, payment channel/status constants, order action constants.
- `consts/after_sale.go`: after-sale type, reason, status, and action constants.
- `domain/orderstate/state.go`: pure order state machine validation and log action mapping.
- `domain/orderstate/state_test.go`: pure tests for legal/illegal transitions and terminal states.
- `domain/aftersale/state.go`: pure after-sale state machine validation.
- `domain/aftersale/state_test.go`: pure tests for after-sale transitions.
- `repository/db/model/order.go`: add `PaymentChannel`, `LogisticsCompany`, `ShippedAt`, `ReceivedAt`, `CanceledAt`.
- `repository/db/model/order_log.go`: immutable operation log model.
- `repository/db/model/order_logistics.go`: local logistics timeline node model.
- `repository/db/model/order_payment.go`: order payment record for balance/wechat/alipay product-order payments.
- `repository/db/model/after_sale.go`: independent after-sale workflow model.
- `repository/db/dao/migrate.go`: AutoMigrate new models and order fields.
- `repository/db/dao/order.go`: add focused transition persistence helpers and list/detail joins for new fields.
- `repository/db/dao/order_log.go`: create/list/backfill log DAO.
- `repository/db/dao/order_logistics.go`: create/list logistics node DAO.
- `repository/db/dao/order_payment.go`: create/get/mark-paid/mark-failed DAO with idempotent updates.
- `repository/db/dao/after_sale.go`: create/list/transition DAO.
- `application/order.go`: order create/cancel/ship/receive/after-sale/refund use cases with transactions.
- `application/payment.go`: balance product-order payment and gateway callback orchestration.
- `application/finance.go`: refund reversal helpers remain transaction-scoped and idempotent.
- `api/v1/order.go`: route handlers for cancel, logs, logistics, after-sale actions.
- `api/v1/payment_gateway.go`: product-order wechat/alipay create and notify routing.
- `types/order.go`: new request/response DTO fields.
- `types/after_sale.go`: after-sale DTOs.
- `routes/router.go`: user, seller, and admin P2 routes.
- `scripts/migrations/backfill_p2_order_logs.go`: one-time order log backfill for old data.
- `web/src/api/order.ts`: new buyer/seller order and after-sale APIs.
- `web/src/views/user/OrderListView.vue`: cancel/refund states and actions.
- `web/src/views/user/OrderDetailView.vue`: logs, local logistics nodes, after-sale panel.
- `web/src/views/checkout/PaymentView.vue`: balance/wechat/alipay product-order payment entry.
- `web/src/views/seller/SellerOrderListView.vue`: logistics company + tracking number ship form and after-sale handling.
- `web-admin/src/api/index.ts`: admin after-sale and order operation APIs.
- `web-admin/src/views/order/OrderView.vue`: order detail/logs and platform intervention/refund actions.
- `web/src/locales/zh-CN.ts`, `web/src/locales/en-US.ts`: user-facing P2 text.
- `web-admin/src/locales/zh-CN.ts`, `web-admin/src/locales/en-US.ts`: admin P2 text.

---

## Task 1: Order State Machine Foundation

**Files:**
- Modify: `consts/order.go`
- Create: `domain/orderstate/state.go`
- Test: `domain/orderstate/state_test.go`

- [ ] **Step 1: Write failing state-machine tests**

Add table tests in `domain/orderstate/state_test.go`:

```go
package orderstate

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestCanTransitionOrderStatusAllowsP2HappyPath(t *testing.T) {
	tests := []struct {
		name string
		from uint
		to   uint
	}{
		{"pay", consts.OrderTypeUnPaid, consts.OrderTypePendingShipping},
		{"ship", consts.OrderTypePendingShipping, consts.OrderTypeShipping},
		{"receive", consts.OrderTypeShipping, consts.OrderTypeReceipt},
		{"request_refund_before_ship", consts.OrderTypePendingShipping, consts.OrderTypeRefundRequested},
		{"request_refund_after_ship", consts.OrderTypeShipping, consts.OrderTypeRefundRequested},
		{"refund", consts.OrderTypeRefundRequested, consts.OrderTypeRefunded},
		{"cancel_unpaid", consts.OrderTypeUnPaid, consts.OrderTypeCanceled},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsureOrderStatusTransition(tt.from, tt.to); err != nil {
				t.Fatalf("expected transition %d -> %d to pass, got %v", tt.from, tt.to, err)
			}
		})
	}
}

func TestCanTransitionOrderStatusRejectsIllegalTransitions(t *testing.T) {
	tests := []struct {
		name string
		from uint
		to   uint
	}{
		{"ship_unpaid", consts.OrderTypeUnPaid, consts.OrderTypeShipping},
		{"pay_shipped", consts.OrderTypeShipping, consts.OrderTypePendingShipping},
		{"cancel_paid", consts.OrderTypePendingShipping, consts.OrderTypeCanceled},
		{"refund_completed_without_after_sale", consts.OrderTypeReceipt, consts.OrderTypeRefunded},
		{"leave_refunded", consts.OrderTypeRefunded, consts.OrderTypeReceipt},
		{"leave_canceled", consts.OrderTypeCanceled, consts.OrderTypePendingShipping},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := EnsureOrderStatusTransition(tt.from, tt.to)
			if err == nil {
				t.Fatalf("expected transition %d -> %d to fail", tt.from, tt.to)
			}
			if code := e.BusinessCode(err); code != e.ErrorOrderStatusTransitionInvalid {
				t.Fatalf("expected ErrorOrderStatusTransitionInvalid, got %d (%v)", code, err)
			}
		})
	}
}

func TestOrderStatusIsTerminal(t *testing.T) {
	if !OrderStatusIsTerminal(consts.OrderTypeCanceled) || !OrderStatusIsTerminal(consts.OrderTypeRefunded) {
		t.Fatal("canceled and refunded must be terminal")
	}
	if OrderStatusIsTerminal(consts.OrderTypeShipping) {
		t.Fatal("shipping must not be terminal")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./domain/orderstate -run 'TestCanTransitionOrderStatus|TestOrderStatusIsTerminal' -count=1
```

Expected: FAIL with undefined `OrderTypeCanceled`, `EnsureOrderStatusTransition`, or `ErrorOrderStatusTransitionInvalid`.

- [ ] **Step 3: Add semantic constants and state helpers**

In `consts/order.go`, keep existing numeric values and append:

```go
const (
	OrderTypeCanceled = 7
)

const (
	OrderActionCreate        = "create"
	OrderActionPay           = "pay"
	OrderActionCancel        = "cancel"
	OrderActionShip          = "ship"
	OrderActionReceive       = "receive"
	OrderActionRefundRequest = "refund_request"
	OrderActionRefundApprove = "refund_approve"
	OrderActionRefundReject  = "refund_reject"
	OrderActionAfterSale     = "after_sale"
)

const (
	OrderPaymentChannelBalance = "balance"
	OrderPaymentChannelWechat  = "wechat"
	OrderPaymentChannelAlipay  = "alipay"

	OrderPaymentStatusPending = "pending"
	OrderPaymentStatusPaid    = "paid"
	OrderPaymentStatusFailed  = "failed"
	OrderPaymentStatusClosed  = "closed"
)
```

Add `domain/orderstate/state.go`:

```go
package orderstate

import (
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

var allowedOrderTransitions = map[uint]map[uint]struct{}{
	consts.OrderTypeUnPaid: {
		consts.OrderTypePendingShipping: {},
		consts.OrderTypeCanceled:        {},
	},
	consts.OrderTypePendingShipping: {
		consts.OrderTypeShipping:         {},
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeShipping: {
		consts.OrderTypeReceipt:          {},
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeReceipt: {
		consts.OrderTypeRefundRequested: {},
	},
	consts.OrderTypeRefundRequested: {
		consts.OrderTypeRefunded:         {},
		consts.OrderTypePendingShipping: {},
		consts.OrderTypeShipping:         {},
		consts.OrderTypeReceipt:          {},
	},
}

func EnsureOrderStatusTransition(from, to uint) error {
	if from == to {
		return nil
	}
	if next, ok := allowedOrderTransitions[from]; ok {
		if _, ok = next[to]; ok {
			return nil
		}
	}
	return e.NewBusinessError(e.ErrorOrderStatusTransitionInvalid)
}

func OrderStatusIsTerminal(status uint) bool {
	return status == consts.OrderTypeCanceled || status == consts.OrderTypeRefunded
}
```

Add `ErrorOrderStatusTransitionInvalid` to `utils/e/code.go`, `utils/e/msg.go`, and locale message maps with stable `msg_key`.

- [ ] **Step 4: Run test to verify it passes**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./domain/orderstate -run 'TestCanTransitionOrderStatus|TestOrderStatusIsTerminal' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add consts/order.go domain/orderstate/state.go domain/orderstate/state_test.go utils/e/code.go utils/e/msg.go config/locales
git commit -m "feat: add order state machine foundation"
```

---

## Task 2: Order Operation Logs and Backfill

**Files:**
- Create: `repository/db/model/order_log.go`
- Create: `repository/db/dao/order_log.go`
- Modify: `repository/db/dao/migrate.go`
- Modify: `types/order.go`
- Create: `scripts/migrations/backfill_p2_order_logs.go`
- Test: `repository/db/dao/order_log_test.go`

- [ ] **Step 1: Write failing model/DAO tests**

Add `repository/db/dao/order_log_test.go`:

```go
package dao

import (
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestOrderLogModelMigratesWithSingularTable(t *testing.T) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local",
		SkipInitializeWithVersion: true,
	}), &gorm.Config{DryRun: true, DisableAutomaticPing: true})
	if err != nil {
		t.Fatal(err)
	}
	stmt := db.Session(&gorm.Session{DryRun: true}).Create(&model.OrderLog{
		OrderID: 1,
		OrderNum: 1001,
		Action: "pay",
		FromType: 1,
		ToType: 2,
		OperatorType: "buyer",
		OperatorID: 7,
		Remark: "balance payment",
	}).Statement
	if stmt.Table != "order_log" {
		t.Fatalf("expected order_log table, got %q", stmt.Table)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./repository/db/dao -run TestOrderLogModelMigratesWithSingularTable -count=1
```

Expected: FAIL because `model.OrderLog` does not exist.

- [ ] **Step 3: Add model, DAO, DTO, migration, and backfill script**

Create `repository/db/model/order_log.go`:

```go
package model

import "gorm.io/gorm"

type OrderLog struct {
	gorm.Model
	OrderID      uint   `gorm:"not null;index" json:"order_id"`
	OrderNum     uint64 `gorm:"not null;index" json:"order_num"`
	Action       string `gorm:"size:32;not null;index" json:"action"`
	FromType     uint   `gorm:"not null" json:"from_type"`
	ToType       uint   `gorm:"not null" json:"to_type"`
	OperatorType string `gorm:"size:32;not null;index" json:"operator_type"`
	OperatorID   uint   `gorm:"not null;index" json:"operator_id"`
	Remark       string `gorm:"size:255" json:"remark"`
}
```

Create `repository/db/dao/order_log.go`:

```go
package dao

import (
	"context"

	"gorm.io/gorm"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

type OrderLogDao struct{ *gorm.DB }

func NewOrderLogDao(ctx context.Context) *OrderLogDao { return &OrderLogDao{NewDBClient(ctx)} }
func NewOrderLogDaoByDB(db *gorm.DB) *OrderLogDao    { return &OrderLogDao{db} }

func (dao *OrderLogDao) Create(log *model.OrderLog) error {
	return dao.DB.Create(log).Error
}

func (dao *OrderLogDao) ListByOrderID(orderID uint) ([]*model.OrderLog, error) {
	list := make([]*model.OrderLog, 0)
	err := dao.DB.Where("order_id = ?", orderID).Order("created_at ASC, id ASC").Find(&list).Error
	return list, err
}
```

Modify `repository/db/dao/migrate.go` to AutoMigrate `&model.OrderLog{}`.

Add DTOs to `types/order.go`:

```go
type OrderLogResp struct {
	ID           uint   `json:"id"`
	OrderID      uint   `json:"order_id"`
	OrderNum     uint64 `json:"order_num"`
	Action       string `json:"action"`
	FromType     uint   `json:"from_type"`
	ToType       uint   `json:"to_type"`
	OperatorType string `json:"operator_type"`
	OperatorID   uint   `json:"operator_id"`
	Remark       string `json:"remark"`
	CreatedAt    int64  `json:"created_at"`
}
```

Create `scripts/migrations/backfill_p2_order_logs.go` as a one-time command that opens the project DB config, scans existing orders, and inserts one `create` or `backfill_current_state` log only when no logs exist for the order. Use `OrderLogDao` and a unique in-code check, not runtime compatibility logic.

- [ ] **Step 4: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./repository/db/dao -run TestOrderLogModelMigratesWithSingularTable -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add repository/db/model/order_log.go repository/db/dao/order_log.go repository/db/dao/order_log_test.go repository/db/dao/migrate.go types/order.go scripts/migrations/backfill_p2_order_logs.go
git commit -m "feat: add order operation logs"
```

---

## Task 3: Application-Layer Order Transitions

**Files:**
- Modify: `application/order.go`
- Modify: `application/payment.go`
- Modify: `service/order.go`
- Modify: `repository/db/dao/order.go`
- Modify: `api/v1/order.go`
- Modify: `routes/router.go`
- Test: `application/order_transition_test.go`

- [ ] **Step 1: Write failing use-case tests**

Add `application/order_transition_test.go` with pure validation tests that avoid DB by testing transition request validation helpers:

```go
package application

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestBuildOrderLogRejectsIllegalTransition(t *testing.T) {
	_, err := buildOrderLog(1, 1001, consts.OrderActionShip, consts.OrderTypeUnPaid, consts.OrderTypeShipping, "seller", 2, "ship unpaid")
	if err == nil {
		t.Fatal("expected illegal transition to fail")
	}
	if e.BusinessCode(err) != e.ErrorOrderStatusTransitionInvalid {
		t.Fatalf("expected transition business code, got %d", e.BusinessCode(err))
	}
}

func TestBuildOrderLogCapturesOperator(t *testing.T) {
	log, err := buildOrderLog(1, 1001, consts.OrderActionPay, consts.OrderTypeUnPaid, consts.OrderTypePendingShipping, "buyer", 7, "balance")
	if err != nil {
		t.Fatal(err)
	}
	if log.OperatorType != "buyer" || log.OperatorID != 7 || log.Action != consts.OrderActionPay {
		t.Fatalf("unexpected log: %+v", log)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application -run TestBuildOrderLog -count=1
```

Expected: FAIL because `buildOrderLog` is undefined.

- [ ] **Step 3: Add application transition helpers**

In `application/order.go`, add:

```go
func buildOrderLog(orderID uint, orderNum uint64, action string, fromType, toType uint, operatorType string, operatorID uint, remark string) (*model.OrderLog, error) {
	if err := orderstate.EnsureOrderStatusTransition(fromType, toType); err != nil {
		return nil, err
	}
	return &model.OrderLog{
		OrderID: orderID, OrderNum: orderNum, Action: action,
		FromType: fromType, ToType: toType,
		OperatorType: operatorType, OperatorID: operatorID,
		Remark: remark,
	}, nil
}
```

Import `github.com/YasinDoyle/e-mall/domain/orderstate` in `application/order.go`. Do not import `service` from `application`; `application/architecture_test.go` must keep passing.

- [ ] **Step 4: Move mutations into application use cases**

Implement these use-case methods in `application/order.go`:

```go
func (u *OrderUsecase) CancelUnpaid(ctx context.Context, orderID uint) (interface{}, error)
func (u *OrderUsecase) Ship(ctx context.Context, req *types.OrderShipReq) (interface{}, error)
func (u *OrderUsecase) Receive(ctx context.Context, req *types.OrderReceiveReq) (interface{}, error)
func (u *OrderUsecase) Logs(ctx context.Context, orderID uint) (interface{}, error)
```

Each method must:

- lock the order row with `GetOrderByIdForUpdate`,
- verify buyer/seller/admin ownership,
- call the state transition helper,
- update the order row,
- create an `OrderLog` row in the same transaction,
- publish `domain/event` only after transaction success.

Update `service/order.go` to become compatibility wrappers calling `application.NewOrderUsecase()` for ship/receive/refund/cancel/logs. Keep simple list/detail methods in `service` for now.

- [ ] **Step 5: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application ./domain/orderstate ./service -run 'TestBuildOrderLog|TestCanTransitionOrderStatus|TestOrderStatusIsTerminal' -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add application/order.go application/payment.go application/order_transition_test.go service/order.go repository/db/dao/order.go api/v1/order.go routes/router.go
git commit -m "feat: route order transitions through application layer"
```

---

## Task 4: Local Logistics Timeline

**Files:**
- Modify: `repository/db/model/order.go`
- Create: `repository/db/model/order_logistics.go`
- Create: `repository/db/dao/order_logistics.go`
- Modify: `repository/db/dao/migrate.go`
- Modify: `types/order.go`
- Modify: `application/order.go`
- Test: `service/order_logistics_test.go`

- [ ] **Step 1: Write failing logistics validation tests**

Add `service/order_logistics_test.go`:

```go
package service

import "testing"

func TestNormalizeShipmentInfoRequiresCompanyAndTrackingNo(t *testing.T) {
	_, err := NormalizeShipmentInfo(" SF ", " SF123456789 ")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = NormalizeShipmentInfo("", "SF123"); err == nil {
		t.Fatal("expected missing logistics company to fail")
	}
	if _, err = NormalizeShipmentInfo("SF", ""); err == nil {
		t.Fatal("expected missing tracking number to fail")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./service -run TestNormalizeShipmentInfoRequiresCompanyAndTrackingNo -count=1
```

Expected: FAIL because `NormalizeShipmentInfo` is undefined.

- [ ] **Step 3: Add order fields and local logistics nodes**

Add to `model.Order`:

```go
PaymentChannel   string     `gorm:"size:20;index"`
LogisticsCompany string     `gorm:"size:64"`
ShippedAt        *time.Time
ReceivedAt       *time.Time
CanceledAt       *time.Time
```

Create `model.OrderLogistics`:

```go
type OrderLogistics struct {
	gorm.Model
	OrderID     uint   `gorm:"not null;index"`
	OrderNum    uint64 `gorm:"not null;index"`
	NodeType    string `gorm:"size:32;not null;index"`
	Description string `gorm:"size:255;not null"`
	OccurredAt  int64  `gorm:"not null;index"`
}
```

Create DAO methods:

```go
func (dao *OrderLogisticsDao) Create(node *model.OrderLogistics) error
func (dao *OrderLogisticsDao) ListByOrderID(orderID uint) ([]*model.OrderLogistics, error)
```

Add `NormalizeShipmentInfo(company, trackingNo string) (types.ShipmentInfo, error)` in a pure helper. It trims both fields, requires both, and caps company at 64 chars and tracking number at 64 chars.

- [ ] **Step 4: Update shipping and receipt use cases**

When seller ships:

- require `logistics_company` and `tracking_no`,
- update order `type`, `logistics_company`, `tracking_no`, `shipped_at`,
- create `OrderLog` action `ship`,
- create `OrderLogistics` node `manual_shipped`.

When buyer receives:

- update order `type` and `received_at`,
- create `OrderLog` action `receive`,
- create `OrderLogistics` node `manual_received`,
- generate settlement for the order.

- [ ] **Step 5: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./service ./repository/db/dao -run 'TestNormalizeShipmentInfo|OrderLogistics' -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add repository/db/model/order.go repository/db/model/order_logistics.go repository/db/dao/order_logistics.go repository/db/dao/migrate.go types/order.go application/order.go service/order_logistics_test.go
git commit -m "feat: add local order logistics timeline"
```

---

## Task 5: Product Order Payment Records and Channels

**Status:** Completed for backend implementation and automated verification. Commit is intentionally deferred until user review, per `AGENTS.md`.

**Post-review decisions:**
- External product-order callbacks validate expected provider channel, amount, status, and idempotency.
- Callback processing locks `order_payment` first so duplicate callbacks for the same `payment_no` return idempotently.
- Only one active external pending payment is allowed per order: same-channel repeat requests reuse the existing `payment_no`; cross-channel repeat requests are rejected until the pending payment is resolved.
- Balance payment is rejected while an external pending payment exists, because P2 does not yet close provider-side QR/pay URLs or auto-refund duplicate external payments.
- Provider order creation errors leave the internal payment `pending` so same-channel retries and late provider callbacks can recover.

**Files:**
- Create: `repository/db/model/order_payment.go`
- Create: `repository/db/dao/order_payment.go`
- Modify: `repository/db/dao/migrate.go`
- Modify: `types/payment.go`
- Modify: `application/payment.go`
- Modify: `service/payment_gateway.go`
- Modify: `api/v1/payment_gateway.go`
- Modify: `routes/router.go`
- Test: `application/order_payment_test.go`

- [x] **Step 1: Write failing payment idempotency tests**

Add `application/order_payment_test.go` with pure tests for payment callback decisions:

```go
package application

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestValidateOrderPaymentCallbackRejectsAmountMismatch(t *testing.T) {
	payment := &model.OrderPayment{OrderID: 1, Amount: 88.66, Status: consts.OrderPaymentStatusPending}
	if err := validateOrderPaymentCallback(payment, 88.65); e.BusinessCode(err) != e.ErrorPaymentAmountMismatch {
		t.Fatalf("expected amount mismatch code, got %d (%v)", e.BusinessCode(err), err)
	}
}

func TestValidateOrderPaymentCallbackIsIdempotentForPaid(t *testing.T) {
	payment := &model.OrderPayment{OrderID: 1, Amount: 88.66, Status: consts.OrderPaymentStatusPaid}
	if err := validateOrderPaymentCallback(payment, 88.66); err != nil {
		t.Fatal(err)
	}
}
```

- [x] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application -run TestValidateOrderPaymentCallback -count=1
```

Expected: FAIL because `model.OrderPayment` or helper functions do not exist.

- [x] **Step 3: Add `OrderPayment` model and DAO**

Create `repository/db/model/order_payment.go`:

```go
package model

import (
	"time"
	"gorm.io/gorm"
)

type OrderPayment struct {
	gorm.Model
	OrderID         uint       `gorm:"not null;index"`
	OrderNum        uint64     `gorm:"not null;index"`
	PaymentNo       string     `gorm:"size:64;not null;uniqueIndex"`
	UserID          uint       `gorm:"not null;index"`
	Channel         string     `gorm:"size:20;not null;index"`
	Amount          float64    `gorm:"not null"`
	Status          string     `gorm:"size:20;not null;default:'pending';index"`
	ProviderTradeNo string     `gorm:"size:128;index"`
	PaidAt          *time.Time
	ClosedAt        *time.Time
}
```

DAO methods:

```go
func (dao *OrderPaymentDao) Create(payment *model.OrderPayment) error
func (dao *OrderPaymentDao) GetByPaymentNo(paymentNo string) (*model.OrderPayment, error)
func (dao *OrderPaymentDao) GetPendingByOrderID(orderID uint) (*model.OrderPayment, error)
func (dao *OrderPaymentDao) MarkPaid(paymentNo, providerTradeNo string, paidAt time.Time) (fresh bool, err error)
func (dao *OrderPaymentDao) MarkFailed(paymentNo string) error
func (dao *OrderPaymentDao) ClosePendingByOrderID(orderID uint) error
```

- [x] **Step 4: Split payment entrypoints**

Keep current `/api/v1/paydown` working for balance payment, but add explicit channel routes:

```go
authed.POST("orders/pay/balance", api.OrderBalancePayHandler())
authed.POST("orders/pay/wechat", api.OrderWechatPayHandler())
authed.POST("orders/pay/alipay", api.OrderAlipayPayHandler())
authed.GET("orders/pay/status", api.OrderPaymentStatusHandler())
```

For balance payments, create an `OrderPayment` row with `channel=balance`, pay immediately in the same transaction, update order via state machine, write `OrderLog`, write P1 account flows/settlement, and emit `OrderPaid` after commit.

For wechat/alipay product-order payments, create an `OrderPayment` row and return QR code/pay URL. Callback handler must route by `payment_no` to product-order payment, verify signature, verify amount, verify pending order status, mark payment paid idempotently, then run the same order-paid application flow as balance payment.

- [x] **Step 5: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application ./repository/db/dao -run 'TestValidateOrderPaymentCallback|OrderPayment' -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit after user review**

Deferred by user rule: do not run `git add`/`git commit` until the user reviews and explicitly asks to submit.

```bash
git add repository/db/model/order_payment.go repository/db/dao/order_payment.go repository/db/dao/migrate.go types/payment.go application/payment.go application/order_payment_test.go service/payment_gateway.go api/v1/payment_gateway.go routes/router.go
git commit -m "feat: add product order payment channels"
```

---

## Task 6: Independent After-Sale Workflow

**Files:**
- Create: `consts/after_sale.go`
- Create: `repository/db/model/after_sale.go`
- Create: `repository/db/dao/after_sale.go`
- Modify: `repository/db/dao/migrate.go`
- Create: `types/after_sale.go`
- Modify: `application/order.go`
- Modify: `api/v1/order.go`
- Modify: `routes/router.go`
- Create: `domain/aftersale/state.go`
- Test: `domain/aftersale/state_test.go`

- [ ] **Step 1: Write failing after-sale state tests**

Add `domain/aftersale/state_test.go`:

```go
package aftersale

import (
	"testing"

	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

func TestEnsureTransitionAllowsP2Flow(t *testing.T) {
	tests := [][2]string{
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusSellerApproved},
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusSellerRejected},
		{consts.AfterSaleStatusSellerRejected, consts.AfterSaleStatusPlatformIntervening},
		{consts.AfterSaleStatusSellerApproved, consts.AfterSaleStatusRefunded},
		{consts.AfterSaleStatusPlatformIntervening, consts.AfterSaleStatusRefunded},
		{consts.AfterSaleStatusRequested, consts.AfterSaleStatusClosed},
	}
	for _, tt := range tests {
		if err := EnsureTransition(tt[0], tt[1]); err != nil {
			t.Fatalf("expected %s -> %s to pass: %v", tt[0], tt[1], err)
		}
	}
}

func TestEnsureTransitionRejectsIllegalFlow(t *testing.T) {
	err := EnsureTransition(consts.AfterSaleStatusRefunded, consts.AfterSaleStatusRequested)
	if e.BusinessCode(err) != e.ErrorAfterSaleStatusInvalid {
		t.Fatalf("expected after-sale status error, got %d (%v)", e.BusinessCode(err), err)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./domain/aftersale -run TestEnsureTransition -count=1
```

Expected: FAIL because constants and helper do not exist.

- [ ] **Step 3: Add after-sale constants, model, DAO, DTOs, and domain state helper**

Create constants:

```go
const (
	AfterSaleTypeRefundOnly   = "refund_only"
	AfterSaleTypeReturnRefund = "return_refund"

	AfterSaleStatusRequested            = "requested"
	AfterSaleStatusSellerApproved       = "seller_approved"
	AfterSaleStatusSellerRejected       = "seller_rejected"
	AfterSaleStatusPlatformIntervening  = "platform_intervening"
	AfterSaleStatusRefunded             = "refunded"
	AfterSaleStatusClosed               = "closed"
)
```

Create `model.AfterSale`:

```go
type AfterSale struct {
	gorm.Model
	OrderID       uint    `gorm:"not null;index"`
	OrderNum      uint64  `gorm:"not null;index"`
	BuyerID       uint    `gorm:"not null;index"`
	SellerID      uint    `gorm:"not null;index"`
	Type          string  `gorm:"size:32;not null;index"`
	Status        string  `gorm:"size:32;not null;index"`
	Reason        string  `gorm:"size:255;not null"`
	RefundAmount  float64 `gorm:"not null"`
	SellerReason  string  `gorm:"size:255"`
	PlatformNote  string  `gorm:"size:255"`
	RefundedAt    *int64
	ClosedAt      *int64
}
```

Add DAO create/list/get-for-update/transition methods and DTOs:

```go
type AfterSaleRequestReq struct {
	OrderId uint   `json:"order_id" form:"order_id" binding:"required"`
	Type    string `json:"type" form:"type" binding:"required"`
	Reason  string `json:"reason" form:"reason" binding:"required"`
}

type SellerAfterSaleHandleReq struct {
	AfterSaleID uint   `json:"after_sale_id" form:"after_sale_id" binding:"required"`
	Action      string `json:"action" form:"action" binding:"required"`
	Reason      string `json:"reason" form:"reason"`
}

type AdminAfterSaleHandleReq struct {
	AfterSaleID uint   `json:"after_sale_id" form:"after_sale_id" binding:"required"`
	Action      string `json:"action" form:"action" binding:"required"`
	Note        string `json:"note" form:"note"`
	Key         string `json:"key" form:"key"`
}
```

Create `domain/aftersale/state.go`:

```go
package aftersale

import (
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/utils/e"
)

var allowedTransitions = map[string]map[string]struct{}{
	consts.AfterSaleStatusRequested: {
		consts.AfterSaleStatusSellerApproved: {},
		consts.AfterSaleStatusSellerRejected: {},
		consts.AfterSaleStatusClosed:         {},
	},
	consts.AfterSaleStatusSellerRejected: {
		consts.AfterSaleStatusPlatformIntervening: {},
		consts.AfterSaleStatusClosed:              {},
	},
	consts.AfterSaleStatusSellerApproved: {
		consts.AfterSaleStatusRefunded: {},
		consts.AfterSaleStatusClosed:   {},
	},
	consts.AfterSaleStatusPlatformIntervening: {
		consts.AfterSaleStatusRefunded: {},
		consts.AfterSaleStatusClosed:   {},
	},
}

func EnsureTransition(from, to string) error {
	if from == to {
		return nil
	}
	if next, ok := allowedTransitions[from]; ok {
		if _, ok = next[to]; ok {
			return nil
		}
	}
	return e.NewBusinessError(e.ErrorAfterSaleStatusInvalid)
}
```

- [ ] **Step 4: Add application use cases and routes**

Implement:

```go
func (u *OrderUsecase) RequestAfterSale(ctx context.Context, req *types.AfterSaleRequestReq) (interface{}, error)
func (u *OrderUsecase) SellerHandleAfterSale(ctx context.Context, req *types.SellerAfterSaleHandleReq) (interface{}, error)
func (u *OrderUsecase) AdminHandleAfterSale(ctx context.Context, req *types.AdminAfterSaleHandleReq) (interface{}, error)
func (u *OrderUsecase) ListBuyerAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error)
func (u *OrderUsecase) ListSellerAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error)
func (u *OrderUsecase) ListAdminAfterSales(ctx context.Context, req *types.AfterSaleListReq) (interface{}, error)
```

Routes:

```go
authed.POST("after-sales/request", api.AfterSaleRequestHandler())
authed.GET("after-sales/list", api.AfterSaleListHandler())
authed.GET("boss/after-sales/list", api.SellerAfterSaleListHandler())
authed.POST("boss/after-sales/handle", api.SellerAfterSaleHandleHandler())
admin.GET("after-sales/list", api.AdminAfterSaleListHandler())
admin.POST("after-sales/handle", api.AdminAfterSaleHandleHandler())
```

- [ ] **Step 5: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./domain/aftersale ./application -run 'TestEnsureTransition|AfterSale' -count=1
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add consts/after_sale.go domain/aftersale/state.go domain/aftersale/state_test.go repository/db/model/after_sale.go repository/db/dao/after_sale.go repository/db/dao/migrate.go types/after_sale.go application/order.go api/v1/order.go routes/router.go
git commit -m "feat: add after-sale workflow"
```

---

## Task 7: Refund Reversal and Idempotent Money Flow

**Files:**
- Modify: `application/finance.go`
- Modify: `application/order.go`
- Modify: `repository/db/dao/settlement.go`
- Modify: `repository/db/dao/account_flow.go`
- Test: `application/refund_reversal_test.go`
- Test: `service/settlement_test.go`

- [ ] **Step 1: Write failing refund idempotency tests**

Add `application/refund_reversal_test.go`:

```go
package application

import (
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestBuildRefundFlowNosAreStablePerOrder(t *testing.T) {
	order := &model.Order{OrderNum: 9001, UserID: 7, BossID: 8}
	flows := buildRefundAccountFlows(order, 66.50)
	if len(flows) != 3 {
		t.Fatalf("expected buyer, seller, platform refund flows, got %d", len(flows))
	}
	seen := map[string]bool{}
	for _, flow := range flows {
		if seen[flow.FlowNo] {
			t.Fatalf("duplicate flow no %s", flow.FlowNo)
		}
		seen[flow.FlowNo] = true
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application -run TestBuildRefundFlowNosAreStablePerOrder -count=1
```

Expected: FAIL because `buildRefundAccountFlows` is missing or returns incomplete flows.

- [ ] **Step 3: Make refund reversal explicit**

Update refund handling to:

- lock the order and after-sale rows,
- reject repeated refund when after-sale is already `refunded` or order is `OrderTypeRefunded`,
- credit buyer through original payment channel when channel is `balance`,
- record a buyer refund flow,
- mark seller pending/generated settlement refunded,
- record seller pending reversal and platform commission reversal flows,
- keep all DB writes in one transaction,
- emit `OrderRefunded` after commit.

For wechat/alipay product orders, P2 should create refund request records and expose status, but actual high-security risk control and automatic reconciliation stay in A1/A4. Gateway refund methods must still validate amount and idempotency.

- [ ] **Step 4: Run targeted tests**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./application ./service -run 'Refund|Settlement' -count=1
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add application/finance.go application/order.go application/refund_reversal_test.go repository/db/dao/settlement.go repository/db/dao/account_flow.go service/settlement_test.go
git commit -m "feat: make refund reversal idempotent"
```

---

## Task 8: User Web P2 Order Experience

**Files:**
- Modify: `web/src/api/order.ts`
- Modify: `web/src/views/checkout/PaymentView.vue`
- Modify: `web/src/views/user/OrderListView.vue`
- Modify: `web/src/views/user/OrderDetailView.vue`
- Modify: `web/src/utils/status-labels.ts`
- Modify: `web/src/locales/zh-CN.ts`
- Modify: `web/src/locales/en-US.ts`
- Test: `web/tests/p2-order-acceptance.test.mjs`

- [ ] **Step 1: Write failing frontend acceptance checks**

Add `web/tests/p2-order-acceptance.test.mjs`:

```js
import { describe, it } from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";

const root = new URL("../src/", import.meta.url);
const read = (path) => fs.readFileSync(new URL(path, root), "utf8");

describe("P2 user order experience", () => {
  it("payment view exposes balance, wechat, and alipay payment channels through i18n", () => {
    const file = read("views/checkout/PaymentView.vue");
    assert.match(file, /payment\.channelBalance/);
    assert.match(file, /payment\.channelWechat/);
    assert.match(file, /payment\.channelAlipay/);
    assert.doesNotMatch(file, /微信支付|支付宝支付|余额支付/);
  });

  it("order detail renders operation logs and local logistics timeline", () => {
    const file = read("views/user/OrderDetailView.vue");
    assert.match(file, /getOrderLogs/);
    assert.match(file, /getOrderLogistics/);
    assert.match(file, /orderDetail\.operationLogs/);
    assert.match(file, /orderDetail\.logisticsTimeline/);
  });

  it("after-sale request is visible from order detail", () => {
    const file = read("views/user/OrderDetailView.vue");
    assert.match(file, /requestAfterSale/);
    assert.match(file, /orderDetail\.afterSaleRequest/);
  });
});
```

- [ ] **Step 2: Run test to verify it fails**

Run:

```bash
cd web && node --test tests/p2-order-acceptance.test.mjs
```

Expected: FAIL because the new APIs and i18n keys are missing.

- [ ] **Step 3: Update user web**

Add API functions:

```ts
export const cancelOrder = (data: { order_id: number }) => request.post("/orders/cancel", data);
export const getOrderLogs = (params: { order_id: number }) => request.get("/orders/logs", { params });
export const getOrderLogistics = (params: { order_id: number }) => request.get("/orders/logistics", { params });
export const payOrderByBalance = (data: { order_id: number; key: string }) => request.post("/orders/pay/balance", data);
export const payOrderByWechat = (data: { order_id: number }) => request.post("/orders/pay/wechat", data);
export const payOrderByAlipay = (data: { order_id: number }) => request.post("/orders/pay/alipay", data);
export const getOrderPaymentStatus = (params: { payment_no: string }) => request.get("/orders/pay/status", { params });
export const requestAfterSale = (data: { order_id: number; type: string; reason: string }) => request.post("/after-sales/request", data);
```

Update views to show:

- payment channel segmented control,
- cancel unpaid button,
- operation log timeline,
- local logistics timeline,
- after-sale request dialog for paid/shipped/received orders,
- after-sale status and refund amount.

All visible strings must use `zh-CN` and `en-US` locale keys.

- [ ] **Step 4: Run user web checks**

Run:

```bash
cd web && node --test tests/p2-order-acceptance.test.mjs
cd web && npm run build
```

Expected: both PASS.

- [ ] **Step 5: Commit**

```bash
git add web/src/api/order.ts web/src/views/checkout/PaymentView.vue web/src/views/user/OrderListView.vue web/src/views/user/OrderDetailView.vue web/src/utils/status-labels.ts web/src/locales/zh-CN.ts web/src/locales/en-US.ts web/tests/p2-order-acceptance.test.mjs
git commit -m "feat: expose p2 order flow in user web"
```

---

## Task 9: Seller and Admin P2 Operations UI

**Files:**
- Modify: `web/src/views/seller/SellerOrderListView.vue`
- Modify: `web/src/api/order.ts`
- Modify: `web/src/locales/zh-CN.ts`
- Modify: `web/src/locales/en-US.ts`
- Modify: `web-admin/src/api/index.ts`
- Modify: `web-admin/src/views/order/OrderView.vue`
- Modify: `web-admin/src/locales/zh-CN.ts`
- Modify: `web-admin/src/locales/en-US.ts`
- Test: `web/tests/p2-order-acceptance.test.mjs`
- Test: `web-admin/tests/p2-admin-order-acceptance.test.mjs`

- [ ] **Step 1: Write failing admin/seller acceptance checks**

Add `web-admin/tests/p2-admin-order-acceptance.test.mjs`:

```js
import { describe, it } from "node:test";
import assert from "node:assert/strict";
import fs from "node:fs";

const root = new URL("../src/", import.meta.url);
const read = (path) => fs.readFileSync(new URL(path, root), "utf8");

describe("P2 admin order operations", () => {
  it("admin order page exposes after-sale list and operation logs", () => {
    const file = read("views/order/OrderView.vue");
    assert.match(file, /getAdminAfterSaleList/);
    assert.match(file, /getAdminOrderLogs/);
    assert.match(file, /page\.order\.afterSale/);
    assert.match(file, /page\.order\.operationLogs/);
  });

  it("admin visible P2 text uses i18n keys", () => {
    const file = read("views/order/OrderView.vue");
    assert.doesNotMatch(file, /平台介入|售后处理|操作日志|退款拒绝/);
  });
});
```

- [ ] **Step 2: Run tests to verify they fail**

Run:

```bash
cd web-admin && node --test tests/p2-admin-order-acceptance.test.mjs
```

Expected: FAIL because admin after-sale/log UI is missing.

- [ ] **Step 3: Update seller center**

Update `SellerOrderListView.vue` so seller can:

- enter logistics company and tracking number,
- view after-sale requests for their orders,
- approve/reject buyer after-sale requests with reason,
- see local logistics status and settlement status together.

Add locale keys under `sellerCenter.order.*` and `sellerCenter.afterSale.*`.

- [ ] **Step 4: Update admin order operations**

Update admin APIs and `OrderView.vue` so admin can:

- filter by order status and after-sale status,
- open a row detail panel with order logs and logistics nodes,
- approve platform intervention refund,
- reject/close after-sale with visible reason,
- see stable business error messages from API response.

Add locale keys under `page.order.*`.

- [ ] **Step 5: Run frontend checks**

Run:

```bash
cd web && node --test tests/p2-order-acceptance.test.mjs
cd web && npm run build
cd web-admin && node --test tests/p2-admin-order-acceptance.test.mjs
cd web-admin && npm run build
```

Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add web/src/views/seller/SellerOrderListView.vue web/src/api/order.ts web/src/locales/zh-CN.ts web/src/locales/en-US.ts web-admin/src/api/index.ts web-admin/src/views/order/OrderView.vue web-admin/src/locales/zh-CN.ts web-admin/src/locales/en-US.ts web-admin/tests/p2-admin-order-acceptance.test.mjs
git commit -m "feat: expose p2 seller and admin operations"
```

---

## Task 10: End-to-End P2 Verification and Plan Status

**Files:**
- Modify: `docs/superpowers/plans/2026-07-31-p2-order-fulfillment-loop.md`
- Modify: `.github/prompts/plan-eMallFullDevelopmentPlan.prompt.md`

- [ ] **Step 1: Run backend verification**

Run:

```bash
env GOCACHE=/private/tmp/e-mall-go-cache go test ./...
```

Expected: PASS.

- [ ] **Step 2: Run frontend verification**

Run:

```bash
cd web && node --test tests/post-p1-acceptance.test.mjs tests/p2-order-acceptance.test.mjs
cd web && npm run build
cd web-admin && node --test tests/p2-admin-order-acceptance.test.mjs
cd web-admin && npm run build
```

Expected: all PASS.

- [ ] **Step 3: Run manual acceptance against local services**

Start the backend and both frontends using the existing project workflow. Exercise these scenarios:

```text
Scenario A: normal fulfillment
1. buyer A logs in and creates an order for seller B's approved product.
2. buyer A pays with balance.
3. seller B sees pending-shipment order and ships with logistics company + tracking number.
4. buyer A opens order detail, sees payment log, shipping log, local logistics node, and confirms receipt.
5. admin generates settlement and marks it paid.
6. seller B sees seller account credited and related account flow.

Scenario B: after-sale refund
1. buyer A pays a second order.
2. buyer A requests refund-only after-sale with reason.
3. seller B approves the after-sale, or rejects and admin intervenes.
4. platform refund completes once.
5. refreshing or repeating the refund action does not duplicate account flows, refund amount, or settlement reversal.
6. buyer, seller, and admin all see the final after-sale/refunded state and operation logs.

Scenario C: isolation
1. seller C cannot ship seller B's order.
2. buyer C cannot load buyer A's order logs, logistics nodes, or after-sale records.
3. tab A and tab B with different accounts keep request tokens isolated through the existing Post-P1 session architecture.
```

- [ ] **Step 4: Update status only if acceptance passes**

If all automated and manual checks pass, update:

```markdown
**Status:** Completed.
```

in this P2 plan and change `.github/prompts/plan-eMallFullDevelopmentPlan.prompt.md` P2 snapshot from `Ready to plan` to `Completed`. If any scenario remains unsupported, mark the task `Partial` and list the exact remaining scenario.

- [ ] **Step 5: Commit**

```bash
git add docs/superpowers/plans/2026-07-31-p2-order-fulfillment-loop.md .github/prompts/plan-eMallFullDevelopmentPlan.prompt.md
git commit -m "docs: record p2 fulfillment verification"
```

---

## Execution Order

1. Task 1 creates the state machine so later changes stop using ad hoc numeric checks.
2. Task 2 adds operation logs and the explicit old-data backfill path.
3. Task 3 moves order mutations into application transactions.
4. Task 4 adds local logistics company, tracking number, and local timeline nodes.
5. Task 5 adds product-order payment records and wechat/alipay order payment channels.
6. Task 6 adds independent after-sale records and role-based handling.
7. Task 7 makes refund reversal and settlement/account flows idempotent.
8. Tasks 8-9 expose the completed flow to buyer, seller, and admin users.
9. Task 10 verifies the business flow and updates phase status only when the actual scenarios pass.

## Plan Self-Review

- Spec coverage: P2 order state machine, payment closure, local logistics, after-sale, refund reversal, notifications through events, frontend entrypoints, and verification scenarios are covered. Real logistics API is explicitly excluded and tracked in P2-Follow-up.
- Placeholder scan: no placeholder markers or open-ended implementation instructions remain. Each task includes concrete files, test commands, expected failures, implementation shape, and commit scope.
- Type consistency: order status constants keep existing `order.type` values; new models use `OrderID`, `OrderNum`, `PaymentNo`, `AfterSaleID`, `LogisticsCompany`, and `TrackingNo` consistently across backend and frontend tasks.
