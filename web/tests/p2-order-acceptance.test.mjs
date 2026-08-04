import { describe, it } from "node:test";
import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

const root = resolve(new URL("../src/", import.meta.url).pathname);
const read = (file) => readFileSync(resolve(root, file), "utf8");

describe("P2 user order experience", () => {
  it("payment api uses the current order-payment endpoints", () => {
    const source = read("api/order.ts");

    assert.match(source, /\/orders\/pay\/balance/);
    assert.match(source, /\/orders\/pay\/wechat/);
    assert.match(source, /\/orders\/pay\/alipay/);
    assert.match(source, /\/orders\/pay\/status/);
    assert.doesNotMatch(source, /\/paydown/);
    assert.equal((source.match(/payOrderByBalance/g) ?? []).length, 1);
  });

  it("payment view exposes queue payment channels and status polling", () => {
    const source = read("views/checkout/PaymentView.vue");

    assert.match(source, /el-segmented/);
    assert.match(source, /payOrderByBalance/);
    assert.match(source, /payOrderByWechat/);
    assert.match(source, /payOrderByAlipay/);
    assert.match(source, /getOrderPaymentStatus/);
    assert.match(source, /pending_orders/);
    assert.match(source, /payment\.quantity/);
    assert.match(source, /emptyQueueHint/);
    assert.match(source, /queueCount/);
    assert.match(source, /methodBalance/);
    assert.match(source, /methodWechat/);
    assert.match(source, /methodAlipay/);
  });

  it("order list exposes cancel-unpaid and payment/refund state labels", () => {
    const source = read("views/user/OrderListView.vue");

    assert.match(source, /cancelOrder/);
    assert.match(source, /cancelUnpaid/);
    assert.match(source, /payUnpaid/);
    assert.match(source, /startUnpaidOrderPayment/);
    assert.match(source, /paymentChannelText/);
    assert.match(source, /refundStateText/);
    assert.match(source, /status\.order\.canceled/);
  });

  it("order detail renders logs, local logistics, and after-sale entrypoints", () => {
    const source = read("views/user/OrderDetailView.vue");

    assert.match(source, /getOrderLogs/);
    assert.match(source, /getAfterSaleList/);
    assert.match(source, /requestAfterSale/);
    assert.match(source, /payUnpaid/);
    assert.match(source, /startUnpaidOrderPayment/);
    assert.match(source, /operationLogTitle/);
    assert.match(source, /logisticsTitle/);
    assert.match(source, /afterSaleTitle/);
    assert.match(source, /afterSaleRefundAmount/);
    assert.match(source, /afterSaleRequestTitle/);
  });

  it("seller order page exposes shipping logistics and after-sale handling", () => {
    const source = read("views/seller/SellerOrderListView.vue");

    assert.match(source, /getSellerAfterSaleList/);
    assert.match(source, /handleSellerAfterSale/);
    assert.match(source, /logisticsCompany/);
    assert.match(source, /afterSaleDialogVisible/);
    assert.match(source, /handleAfterSaleReject/);
    assert.match(source, /afterSaleHandleSuccess/);
    assert.match(source, /sellerCenter\.order\.afterSaleDialog/);
  });

  it("seller order page surfaces pending after-sales before the table", () => {
    const source = read("views/seller/SellerOrderListView.vue");

    assert.match(source, /pendingAfterSaleCount/);
    assert.match(source, /activeAfterSaleByOrderId/);
    assert.match(source, /afterSaleActionType/);
    assert.match(source, /sellerCenter\.order\.pendingAfterSaleNotice/);
    assert.match(source, /status:\s*"requested"/);
  });

  it("wallet balance can be viewed without a password inside the recent auth window", () => {
    const source = read("views/user/WalletView.vue");

    assert.match(source, /canViewBalanceWithoutPassword/);
    assert.match(source, /wallet\.balanceRecentAuthHint/);
    assert.match(source, /getMoney\(\)/);
  });

  it("account switching leaves stale order pages by navigating to home", () => {
    const source = read("components/common/NavBar.vue");

    assert.match(source, /window\.location\.assign\("\/"\)/);
    assert.doesNotMatch(source, /window\.location\.reload\(\)/);
  });

  it("locale files contain the new visible copy keys", () => {
    const zh = read("locales/zh-CN.ts");
    const en = read("locales/en-US.ts");

    for (const key of [
      "emptyQueueHint:",
      "methodBalance:",
      "methodWechat:",
      "methodAlipay:",
      "quantity:",
      "cancelUnpaid:",
      "payUnpaid:",
      "refundState:",
      "operationLogTitle:",
      "logisticsTitle:",
      "afterSaleTitle:",
      "afterSaleRefundAmount:",
      "afterSaleRequestTitle:",
      "logisticsCompany:",
      "afterSaleReasonPlaceholder:",
      "logisticsCreated:",
      "logisticsCanceled:",
      "afterSaleHandleSuccess:",
      "pendingAfterSaleNotice:",
      "balanceRecentAuthHint:",
      "sellerCenter:",
    ]) {
      assert.match(zh, new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")));
      assert.match(en, new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")));
    }
  });
});
