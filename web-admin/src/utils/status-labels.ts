import { t } from "@/locales";

export function orderStatusText(type: number) {
  const labels: Record<number, string> = {
    1: "status.order.unpaid",
    2: "status.order.pendingShipment",
    3: "status.order.shipped",
    4: "status.order.completed",
    5: "status.order.refunding",
    6: "status.order.refunded",
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
