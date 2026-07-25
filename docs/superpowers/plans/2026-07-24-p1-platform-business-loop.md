# P1 Platform Business Loop Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first platform business loop so a normal account can apply as a seller, be approved by admin, publish products as a seller, and later support commission, settlement, and account flows.

**Architecture:** Keep one unified user account and add seller capability through `SellerProfile`. Seller admission, product publishing permissions, commission calculation, settlements, and account flows are separate modules so each step can be tested independently.

**Tech Stack:** Go, Gin, GORM, MySQL, Vue 3, TypeScript, Element Plus, Pinia.

---

## File Structure

- `consts/seller.go`: seller status constants and labels.
- `repository/db/model/seller.go`: seller profile model.
- `repository/db/dao/seller.go`: seller profile persistence operations.
- `service/seller.go`: user-facing seller application and profile query.
- `api/v1/seller.go`: user-facing seller HTTP handlers.
- `types/seller.go`: request and response DTOs for seller APIs.
- `routes/router.go`: user and admin seller routes.
- `repository/db/dao/migrate.go`: include `SellerProfile` in AutoMigrate.
- `service/admin.go`: admin seller listing and audit service methods.
- `api/v1/admin.go`: admin seller handlers.
- `types/admin.go`: admin seller list/audit DTOs.
- `service/product.go`: block product publishing and sale enablement unless seller is approved.
- `web/src/views/seller/`: seller center pages after backend is ready.
- `web-admin/src/views/seller/`: seller audit management after backend is ready.

---

### Task 0: Business Error Code Foundation

**Files:**
- Create: `utils/e/business_error.go`
- Test: `utils/e/business_error_test.go`
- Modify: `utils/e/code.go`
- Modify: `utils/e/msg.go`
- Modify: `api/v1/common.go`
- Modify: `service/seller.go`
- Modify: `service/product.go`
- Modify: `service/admin.go`
- Test: `service/seller_test.go`
- Test: `service/product_pay_key_test.go`

- [x] Add a backend business error type that carries a stable `utils/e` code and message.
- [x] Make `api/v1.ErrorResponse` preserve business error codes instead of collapsing all service errors to `500`.
- [x] Add P1 seller/product/pay-key error codes and messages.
- [x] Convert seller application, seller audit validation, product publish, product on-sale, and admin product audit guard failures to structured business errors.
- [x] Keep database/system errors as `500` unless they are intentionally mapped to a business code.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./...`.

### Task 1: Seller Profile Backend Foundation

**Files:**
- Create: `consts/seller.go`
- Create: `repository/db/model/seller.go`
- Create: `repository/db/dao/seller.go`
- Create: `types/seller.go`
- Modify: `repository/db/dao/migrate.go`
- Test: `repository/db/model/seller_test.go`
- Test: `service/seller_test.go`

- [x] Add seller status constants: pending, approved, rejected, banned.
- [x] Add `SellerProfile` model with `UserID`, `ShopName`, `Description`, `Status`, `RejectReason`, `ApprovedAt`.
- [x] AutoMigrate `SellerProfile`.
- [x] Add DAO methods: create, get by user, list by status, audit.
- [x] Add tests for seller status helpers and seller request validation.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./repository/db/model ./repository/db/dao ./service`.

### Task 2: User Seller Application APIs

**Files:**
- Create: `service/seller.go`
- Create: `api/v1/seller.go`
- Modify: `routes/router.go`
- Test: `service/seller_test.go`

- [x] Add `POST /api/v1/seller/apply` for logged-in users.
- [x] Add `GET /api/v1/seller/profile` for logged-in users.
- [x] Prevent duplicate seller applications.
- [x] Allow rejected sellers to resubmit by moving status back to pending and clearing reject reason.
- [x] Return approved/pending/rejected status to frontend.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./service ./api/v1 ./routes`.

### Task 3: Admin Seller Audit APIs

**Files:**
- Modify: `service/admin.go`
- Modify: `api/v1/admin.go`
- Modify: `types/admin.go`
- Modify: `routes/router.go`
- Test: `service/seller_test.go`

- [x] Add `GET /api/v1/admin/seller/list` with status filter and pagination.
- [x] Add `POST /api/v1/admin/seller/audit` with approve/reject action.
- [x] Require reject reason when rejecting.
- [x] Set `ApprovedAt` only when approving.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./service ./repository/db/dao`.

### Task 4: Gate Seller Product Publishing

**Files:**
- Modify: `service/product.go`
- Test: `service/product_pay_key_test.go`

- [x] Product create requires an approved seller profile.
- [x] Product sale enablement requires approved seller profile and pay key.
- [x] Admin product approval verifies the product owner is an approved seller.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./service`.

### Task 5: User Seller Center Frontend

**Files:**
- Modify: `web/src/api/product.ts`
- Create: `web/src/api/seller.ts`
- Create: `web/src/views/seller/SellerLayout.vue`
- Create: `web/src/views/seller/SellerApplyView.vue`
- Create: `web/src/views/seller/SellerProductListView.vue`
- Create: `web/src/views/seller/SellerProductFormView.vue`
- Modify: `web/src/router/index.ts`
- Modify: `web/src/components/common/NavBar.vue`
- Test: `npm run build` in `web/`

- [ ] Add seller center entry in logged-in user dropdown.
- [ ] Add seller application form.
- [ ] Add seller status page for pending/rejected states.
- [ ] Add own product list with audit status and on/off sale controls.
- [ ] Add product publishing form using existing multipart product APIs.
- [ ] Run `cd web && npm run build`.

### Task 6: Admin Seller Audit Frontend

**Files:**
- Create: `web-admin/src/api/seller.ts`
- Create: `web-admin/src/views/seller/SellerView.vue`
- Modify: `web-admin/src/router/index.ts`
- Modify: admin sidebar layout file found in `web-admin/src`
- Test: `npm run build` in `web-admin/`

- [ ] Add seller management menu.
- [ ] Add seller list filters by pending/approved/rejected/banned.
- [ ] Add approve/reject actions.
- [ ] Require reject reason in dialog.
- [ ] Run `cd web-admin && npm run build`.

### Task 7: Commission and Account Flow Foundation

**Files:**
- Create: `repository/db/model/account_flow.go`
- Create: `repository/db/model/commission.go`
- Create: `repository/db/model/settlement.go`
- Create: `repository/db/dao/account_flow.go`
- Create: `repository/db/dao/settlement.go`
- Create: `service/settlement.go`
- Modify: `repository/db/dao/migrate.go`
- Modify: payment success flow in `service/payment.go`
- Test: `service/settlement_test.go`

- [ ] Add platform commission config model.
- [ ] On payment success, create immutable buyer, seller, and platform account flow rows.
- [ ] Put seller income into pending settlement.
- [ ] Keep first version commission as global percentage.
- [ ] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./...`.

### Task 8: Settlement Admin Operations

**Files:**
- Create: `api/v1/settlement.go`
- Modify: `routes/router.go`
- Create: `web-admin/src/views/settlement/SettlementView.vue`
- Test: backend and admin builds.

- [ ] Add admin settlement list.
- [ ] Add generate settlement action for approved seller and completed orders.
- [ ] Add mark settlement paid action.
- [ ] Add settlement detail account flow list.
- [ ] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./...`.
- [ ] Run `cd web-admin && npm run build`.

---

## Execution Order

Implement Task 0 before continuing frontend or fund-flow work so API callers receive stable business codes. Tasks 1-4 make the seller identity boundary real. Then implement Tasks 5-6 so users and admins can operate it. Implement Tasks 7-8 only after seller product publishing is stable, because commission and settlement depend on reliable order and seller ownership data.
