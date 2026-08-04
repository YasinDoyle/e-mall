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

export function orderStatusTagType(type: number) {
  const labels: Record<number, string> = {
    1: "warning",
    2: "primary",
    3: "warning",
    4: "success",
    5: "danger",
    6: "info",
    7: "info",
  };
  return labels[type] ?? "info";
}

export function refundStatusText(status: number) {
  const labels: Record<number, string> = {
    0: "status.refund.none",
    1: "status.refund.requested",
    2: "status.refund.refunded",
  };
  return t(labels[status] ?? "status.refund.processing", t("common.unknown", "未知"));
}

export function paymentChannelText(channel?: string) {
  const labels: Record<string, string> = {
    balance: "payment.channelBalance",
    wechat: "payment.channelWechat",
    alipay: "payment.channelAlipay",
  };
  return t(labels[channel ?? ""] ?? "common.unknown", t("common.unknown", "未知"));
}

export function orderActionText(action: string) {
  const labels: Record<string, string> = {
    create: "orderAction.create",
    pay: "orderAction.pay",
    cancel: "orderAction.cancel",
    ship: "orderAction.ship",
    receive: "orderAction.receive",
    refund_request: "orderAction.refundRequest",
    refund_approve: "orderAction.refundApprove",
    refund_reject: "orderAction.refundReject",
    after_sale: "orderAction.afterSale",
  };
  return t(labels[action] ?? "common.unknown", t("common.unknown", "未知"));
}

export function operatorTypeText(type: string) {
  const labels: Record<string, string> = {
    buyer: "orderDetail.operatorBuyer",
    seller: "orderDetail.operatorSeller",
    admin: "orderDetail.operatorAdmin",
    system: "orderDetail.operatorSystem",
  };
  return t(labels[type] ?? "common.unknown", t("common.unknown", "未知"));
}

export function afterSaleTypeText(type: string) {
  const labels: Record<string, string> = {
    refund_only: "orderDetail.afterSaleRefundOnly",
    return_refund: "orderDetail.afterSaleReturnRefund",
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
