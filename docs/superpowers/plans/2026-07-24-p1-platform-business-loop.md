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

- [x] Add seller center entry in logged-in user dropdown.
- [x] Add seller application form.
- [x] Add seller status page for pending/rejected states.
- [x] Add own product list with audit status and on/off sale controls.
- [x] Add product publishing form using existing multipart product APIs.
- [x] Run `cd web && npm run build`.

### Task 6: Admin Seller Audit Frontend

**Files:**
- Create: `web-admin/src/api/seller.ts`
- Create: `web-admin/src/views/seller/SellerView.vue`
- Modify: `web-admin/src/router/index.ts`
- Modify: admin sidebar layout file found in `web-admin/src`
- Test: `npm run build` in `web-admin/`

- [x] Add seller management menu.
- [x] Add seller list filters by pending/approved/rejected/banned.
- [x] Add approve/reject actions.
- [x] Require reject reason in dialog.
- [x] Run `cd web-admin && npm run build`.

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

- [x] Add platform commission config model.
- [x] On payment success, create immutable buyer, seller, and platform account flow rows.
- [x] Put seller income into pending settlement.
- [x] Keep first version commission as global percentage.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./...`.

### Task 8: Settlement Admin Operations

**Files:**
- Create: `api/v1/settlement.go`
- Modify: `routes/router.go`
- Create: `web-admin/src/views/settlement/SettlementView.vue`
- Test: backend and admin builds.

- [x] Add admin settlement list.
- [x] Add generate settlement action for approved seller and completed orders.
- [x] Add mark settlement paid action.
- [x] Add settlement detail account flow list.
- [x] Run `env GOCACHE=/private/tmp/e-mall-go-cache go test ./...`.
- [x] Run `cd web-admin && npm run build`.

---

## Execution Order

Implement Task 0 before continuing frontend or fund-flow work so API callers receive stable business codes. Tasks 1-4 make the seller identity boundary real. Then implement Tasks 5-6 so users and admins can operate it. Implement Tasks 7-8 only after seller product publishing is stable, because commission and settlement depend on reliable order and seller ownership data.

---

## Post-P1 Follow-Up Before Part 1.1 Deepening

P1 focuses on making the platform business loop usable. After P1 is accepted, plan these narrow follow-up upgrades as independent work before entering the heavier Part 1.1 phases.

**Scope boundary with `.github/prompts/plan-eMallFullDevelopmentPlan.prompt.md`:**

- Keep Post-P1 focused on usability gaps and first architecture boundaries exposed by P1.
- Do not duplicate P2 order-state/payment/refund refactors; Post-P1 may subscribe to those state changes, but the state machine remains P2 scope.
- Do not duplicate C1/C2/A4 production architecture work; reliable MQ, local message table, dead-letter queues, pressure testing, tracing, and compensation platforms remain C1/C2/A4 scope.
- Do not duplicate A2/A3 search and recommendation depth; product metadata added here should only support seller/admin review and future indexing.
- Do not duplicate A1 payment risk and third-party payout security; seller withdrawal here is the first account/application boundary, while channel security, audit hardening, and risk control remain A1 scope.

### Post-P1 Task A: Realtime Notification MVP and Subscription Boundary

**Status:** Completed. Implemented persistent notification records, user/admin notification APIs, unread count, mark-read operations, SSE unread stream, polling fallback, and user/admin notification views. Follow-up acceptance fix on 2026-07-29: SSE is the default primary channel; navbar/admin layout stop steady unread-count polling while SSE is active, and polling only starts when SSE is disabled or connection setup fails.

**Goal:** Replace refresh-only status discovery with a usable notification inbox and browser subscription boundary, without taking on A4 reliable messaging infrastructure yet.

- Add a `notification` model/table with recipient, scene, title/content, payload, read status, and timestamps.
- Create notifications from P1/P2-visible business state changes: seller audit approved/rejected, product audit result, order paid/shipped/refunded, settlement generated/paid.
- Add user/admin notification APIs: list, unread count, mark read, mark all read.
- Add SSE or WebSocket browser subscription endpoint for logged-in users and admins.
- In this MVP, notification creation stores durable rows and signals active SSE connections through an in-process notification hub; browser clients receive unread-count updates and then load notification list/detail through the notification APIs.
- User web subscribes to own notifications and refreshes seller/product/order state when relevant notifications arrive.
- Admin web subscribes to pending-work notifications and shows badges/toasts for seller/product/refund/settlement queues.
- Provide polling fallback for environments where persistent connections are unavailable; do not run periodic unread-count polling at the same time as a healthy SSE subscription.
- Keep RabbitMQ, local message table, retry, dead-letter queues, and cross-service event reliability in A4.

### Post-P1 Task B: Runtime Configuration for Frontend and Fixed Parameters

**Status:** Completed. Added frontend-local config stores, boot-time config loading, configurable app/admin title text, default notification page size, notification polling interval, upload limit, and feature flags. Removed the backend public config API to keep frontend display defaults out of the Go service boundary.

**Goal:** Move site name, app title, brand text, static defaults, polling intervals, and feature flags out of hardcoded frontend constants.

- Add frontend runtime config loading from local defaults and Vite environment variables.
- Include site name, admin site name, logo text, default page size, notification polling interval, upload limits, and enabled feature flags.
- Load config before app mount and expose it through a Pinia app config store.
- Keep environment-specific values out of component code.
- Add a backend `/api/v1/capabilities/public` style endpoint later only when frontend state must reflect real backend capabilities, such as P2 payment methods or upload constraints.

### Post-P1 Task C: Internationalization Architecture

**Status:** Completed. User web and admin web now ship `zh-CN` and `en-US` locale resources, language switchers, dynamic Element Plus locale, locale headers, and backend localized business errors. Follow-up acceptance fixes moved user web checkout/payment, product detail/search/list, order detail/success, registration, not-found, seller/user center pages, and admin business pages behind locale resources. `web/tests/post-p1-acceptance.test.mjs` now locks the migrated user-facing slices and asserts admin page components stay free of hardcoded Chinese copy.

**Goal:** Introduce a real i18n boundary instead of scattering fixed Chinese strings through frontend and backend.

- Add frontend i18n using Vue I18n with `zh-CN` as the initial locale.
- Add `en-US` locale files for user web and admin web covering migrated menus, buttons, validation messages, status labels, and API error messages.
- Add a language switcher in the user web navbar and admin layout, with selected locale persisted locally.
- Make Element Plus locale switch dynamically with the selected frontend locale.
- Send axios and SSE `X-Locale` / `Accept-Language` headers based on the selected locale, not a hardcoded default.
- Add backend `en-US` business error messages while keeping `utils/e` error codes and `msg_key` stable.
- Move Element Plus locale setup, menus, buttons, validation messages, and status labels into locale files.
- Add backend message keys for business errors and API-facing labels; keep `utils/e` codes stable while allowing locale-specific messages.
- Decide locale propagation strategy: request header, user preference, or query fallback.
- Add tests or manual QA cases proving stable error codes independent of displayed language, and proving user/admin language switching updates frontend labels plus backend business error display messages.
- Acceptance fix added `web/tests/post-p1-acceptance.test.mjs` to lock the user-visible i18n coverage slices already migrated, including the 2026-07-29 user-reported pages.

### Post-P1 Task D: Product Detail and Audit Information Enrichment

**Status:** Completed. Added product audit enrichment fields, product certificate attachments, seller publish/edit form support, product detail display, and admin audit detail review.

**Goal:** Add the product information needed for seller publishing and admin audit detail pages, while leaving search indexing and category-specific commerce modeling to later phases.

- Extend product model/DTOs for review-facing fields such as brand, origin, specification, production date, shelf life, service guarantees, and certificate metadata.
- Add product certificate/image attachments for qualification certificate, quality inspection certificate, authorization certificate, and other category-specific proofs.
- Update seller product publishing/editing to upload and manage the first certificate attachment set.
- Update Admin product detail/audit page to review full product info and certificates before approval.
- Keep current product list views concise, with a separate detail entry point for full review.
- Leave ES indexing, search filters, recommendation features, and category-specific certificate rules to A2/A3 or dedicated later plans.

### Post-P1 Task E: Multi-Account Frontend Session Architecture

**Status:** Completed. User web and admin web keep the shared account session pool in `localStorage`, while the active account pointer is tab-scoped in `sessionStorage`. A fresh tab/window no longer auto-logs in from the last active account; it starts from a pending unauthenticated tab session. Normal logout clears only the current tab's active session and keeps the shared account pool, so tab B logout does not remove account A from tab A. User web and admin web expose explicit saved-account switch/add-account UI. Acceptance tests cover fresh-tab isolation, current-tab-only logout preservation, account switcher UI, and request-token selection isolation across two simulated tab-scoped sessions.

**Goal:** Support multiple accounts in the same browser without later logins overwriting earlier sessions.

- Audit current token, user info, cart count, seller profile, and app cache storage keys in user web and admin web.
- Replace single global `localStorage` session keys with account-scoped session namespaces, for example by user ID, session ID, or selected workspace profile.
- Decide supported UX: account switcher in one tab, isolated sessions per tab/window, or both.
- Support multiple tabs operating different active accounts at the same time: keep the shared account session pool in `localStorage`, but move the active account pointer to per-tab `sessionStorage`.
- When a new tab opens without a tab-scoped active account, start from a pending unauthenticated tab session; do not inherit the last active account from another tab.
- Ensure request interceptors attach the token for the active account only.
- Scope Pinia persistence, seller profile cache, cart state, and notification subscriptions to the active account.
- Add logout behavior that clears only the current tab's active session unless the user chooses to remove an account from the shared session pool or clear all sessions.
- Add regression tests or manual QA cases for logging in as account A and account B in the same browser without state replacement, and for tab A using account A while tab B uses account B without request-token crossover.
- Acceptance fix added `web/tests/post-p1-acceptance.test.mjs` for fresh-tab isolation and current-tab-only logout preservation.

### Post-P1 Task F: Seller Account and Withdrawal MVP

**Status:** Completed.

**Goal:** Separate seller operating funds from user coin wallet and support auditable withdrawal applications, without taking on third-party payout/risk-control hardening yet.

- Add a `seller_account` model/table for seller funds: available balance, frozen balance, total income, total withdrawn, and timestamps.
- Treat settlement paid income as seller account credit, not as user coin wallet balance.
- Add seller account flow records for settlement credit, withdrawal freeze, withdrawal paid, withdrawal rejected/unfrozen, and manual adjustment.
- Add a `seller_withdraw` model/table for withdrawal applications: seller ID, amount, status, payee account info, audit reason, operator, requested time, audited time, and paid time.
- Add seller-side APIs and UI for viewing seller account balance, withdrawal history, and creating withdrawal requests.
- When seller creates a withdrawal request, validate available balance and move amount from available to frozen.
- Add admin withdrawal review and payout UI: approve/reject, mark paid/failed, and show related account flows.
- On approval/paid, reduce frozen balance and increase total withdrawn; on reject/failure, move frozen amount back to available.
- Keep user coin wallet and seller account independent; do not require seller payment password to change seller operating funds.
- Leave third-party payout channel integration, payout risk scoring, high-risk audit hardening, and reconciliation automation to A1/A4.
- Add tests for insufficient balance, duplicate payout prevention, reject unfreeze, paid state transition, and account flow consistency.

### Post-P1 Task G: Backend Application Layer and Domain Event Boundary

**Status:** Completed.

**Goal:** Stop direct service-to-service orchestration by introducing an `application` layer for complete use cases and a lightweight domain-event boundary for side effects, while keeping reliable MQ/outbox work in A4.

**Architecture decision:**

- Use方案 B for strong-consistency flows: `api/v1` calls an application/usecase object, and that usecase opens the transaction and coordinates DAOs/domain helpers.
- Use方案 C for side effects: core business code publishes domain events, and event handlers create notifications or sync product indexes.
- Keep domain services pure where possible: validation, amount calculation, and state transitions should not call other services or HTTP handlers.
- Keep `service` methods as endpoint-facing compatibility wrappers during migration, but forbid new direct `GetXxxSrv()` calls from one service to another.
- Keep A4 reliable messaging out of Post-P1: the event publisher is in-process and synchronous for now, with an interface that can later be backed by outbox/RabbitMQ.

**Migration plan:**

- Add `application` package for cross-module use cases that currently require service orchestration.
- Add `domain/event` package with event names, payloads, a `Publisher` interface, and an in-process publisher implementation.
- Move notification writes behind event handlers instead of calling `NotificationSrv` from seller/order/payment/product/admin/settlement services.
- Move product index sync/delete behind product domain events instead of calling `ProductIndexSrv` from product/admin services.
- Move settlement/seller-account coordination for payment, refund, and settlement-paid flows into application usecases; the usecase controls the GORM transaction.
- Add architecture regression tests that fail if `service/*.go` adds direct `GetXxxSrv()` calls, except for the service's own constructor or explicitly grandfathered migration shims.
- Document the new backend layering rule in `AGENTS.md` after the first migration lands.

**Initial target scope:**

- Notification and product-index side effects must use domain events in this task.
- Payment success, refund approval, and settlement paid must stop calling other services directly and instead go through application/usecase orchestration.
- Existing HTTP routes and frontend API contracts must remain stable.
