import { t } from "@/locales";

export function orderStatusText(type: number) {
  const labels: Record<number, string> = {
    1: "status.order.unpaid",
    2: "status.order.pendingShipment",
    3: "status.order.shipped",
    4: "status.order.completed",
    5: "status.order.refunding",
    6: "status.order.refunded",
    7: "status.order.canceled",
  };
  return t(labels[type] ?? "common.unknown", t("common.unknown", "未知"));
}

export function refundStatusText(status: number) {
  const labels: Record<number, string> = {
    0: "status.refund.none",
    1: "status.refund.requested",
    2: "status.refund.refunded",
  };
  return t(labels[status] ?? "common.unknown", t("common.unknown", "未知"));
}

export function paymentChannelText(channel?: string) {
  const labels: Record<string, string> = {
    balance: "page.order.channelBalance",
    wechat: "page.order.channelWechat",
    alipay: "page.order.channelAlipay",
  };
  return t(labels[channel ?? ""] ?? "common.unknown", t("common.unknown", "未知"));
}

export function orderActionText(action: string) {
  const labels: Record<string, string> = {
    create: "page.order.actionCreate",
    pay: "page.order.actionPay",
    cancel: "page.order.actionCancel",
    ship: "page.order.actionShip",
    receive: "page.order.actionReceive",
    refund_request: "page.order.actionRefundRequest",
    refund_approve: "page.order.actionRefundApprove",
    refund_reject: "page.order.actionRefundReject",
    after_sale: "page.order.actionAfterSale",
  };
  return t(labels[action] ?? "common.unknown", t("common.unknown", "未知"));
}

export function operatorTypeText(type: string) {
  const labels: Record<string, string> = {
    buyer: "page.order.operatorBuyer",
    seller: "page.order.operatorSeller",
    admin: "page.order.operatorAdmin",
    system: "page.order.operatorSystem",
  };
  return t(labels[type] ?? "common.unknown", t("common.unknown", "未知"));
}

export function afterSaleTypeText(type: string) {
  const labels: Record<string, string> = {
    refund_only: "page.order.afterSaleRefundOnly",
    return_refund: "page.order.afterSaleReturnRefund",
  };
  return t(labels[type] ?? "common.unknown", t("common.unknown", "未知"));
}

export function afterSaleStatusText(status: string) {
  const labels: Record<string, string> = {
    requested: "status.afterSale.requested",
    seller_approved: "status.afterSale.sellerApproved",
    seller_rejected: "status.afterSale.sellerRejected",
    platform_intervening: "status.afterSale.platformIntervening",
    refunded: "status.afterSale.refunded",
    closed: "status.afterSale.closed",
  };
  return t(labels[status] ?? "common.unknown", t("common.unknown", "未知"));
}

export function afterSaleStatusTagType(status: string) {
  const labels: Record<string, string> = {
    requested: "warning",
    seller_approved: "primary",
    seller_rejected: "info",
    platform_intervening: "warning",
    refunded: "success",
    closed: "info",
  };
  return labels[status] ?? "info";
}
