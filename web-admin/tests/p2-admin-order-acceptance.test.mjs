import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { describe, it } from "node:test";

const root = resolve(new URL("../src/", import.meta.url).pathname);
const read = (file) => readFileSync(resolve(root, file), "utf8");

describe("P2 admin order operations", () => {
  it("order page exposes detail logs logistics and after-sale handling", () => {
    const source = read("views/order/OrderView.vue");

    assert.match(source, /getAdminOrderDetail/);
    assert.match(source, /getAdminOrderLogs/);
    assert.match(source, /getAdminAfterSaleList/);
    assert.match(source, /handleAdminAfterSale/);
    assert.match(source, /detailDrawerVisible/);
    assert.match(source, /afterSaleIntervene/);
    assert.match(source, /afterSaleRefund/);
    assert.match(source, /afterSaleClose/);
  });

  it("admin locale files contain the P2 order labels", () => {
    const zh = read("locales/zh-CN.ts");
    const en = read("locales/en-US.ts");

    for (const key of [
      "afterSaleStatus:",
      "afterSaleList:",
      "operationLogs:",
      "logisticsTimeline:",
      "afterSaleHandleSuccess:",
      "afterSaleIntervene:",
      "afterSaleRefund:",
      "afterSaleClose:",
      "afterSaleCloseNoteRequired:",
      "paymentChannel:",
      "logisticsCompany:",
    ]) {
      assert.match(zh, new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")));
      assert.match(en, new RegExp(key.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")));
    }
  });

  it("admin order badge includes independent pending after-sales", () => {
    const source = read("components/AdminLayout.vue");

    assert.match(source, /getAdminAfterSaleList/);
    assert.match(source, /afterSaleRequested/);
    assert.match(source, /afterSaleApproved/);
    assert.match(source, /afterSaleRejected/);
    assert.match(source, /refundOrderCount/);
  });

  it("admin after-sale refund does not ask for the buyer fund key", () => {
    const source = read("views/order/OrderView.vue");

    assert.doesNotMatch(source, /afterSaleRefundKey/);
    assert.doesNotMatch(source, /afterSaleRefundKeyPlaceholder/);
    assert.doesNotMatch(source, /key:\s*afterSaleRefundKey\.value/);
  });

  it("admin after-sale detail opens by order id instead of after-sale id", () => {
    const source = read("views/order/OrderView.vue");

    assert.match(source, /openAfterSaleOrderDetail/);
    assert.match(source, /openAfterSaleOrderDetail\(row\)/);
    assert.match(source, /openDetail\(Number\(row\.order_id\)\)/);
  });

  it("admin can only close unresolved or rejected after-sales, not seller-approved refunds", () => {
    const source = read("views/order/OrderView.vue");

    assert.match(source, /status === "requested" \|\| status === "seller_rejected" \|\| status === "platform_intervening"/);
    assert.match(source, /afterSaleCloseNoteRequired/);
  });
});
