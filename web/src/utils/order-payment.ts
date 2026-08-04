import type { Router } from "vue-router";
import type { Order } from "@/types";
import { activeUserSessionStorageKey } from "@/utils/session";

export function buildPendingPaymentOrder(order: Partial<Order> & Record<string, any>) {
  const orderID = Number(order.order_id ?? order.id ?? 0);
  return {
    ...order,
    order_id: orderID,
    product_id: Number(order.product_id ?? 0),
    boss_id: Number(order.boss_id ?? 0),
    num: Number(order.num ?? 1),
    money: Number(order.money ?? order.discount_price ?? order.price ?? 0),
    order_num: order.order_num,
    resumed_from_order: true,
  };
}

export function startUnpaidOrderPayment(
  router: Router,
  order: Partial<Order> & Record<string, any>,
) {
  sessionStorage.setItem(
    activeUserSessionStorageKey("pending_orders"),
    JSON.stringify([buildPendingPaymentOrder(order)]),
  );
  router.push("/payment");
}
