import request from "@/utils/request";
import type { OrderCreateReq } from "@/types";

export const getOrderList = (params: {
  page_num: number;
  page_size: number;
  type?: number;
}) =>
  request.get("/orders/list", { params });

export const getOrderDetail = (params: { order_id: number }) =>
  request.get("/orders/show", { params });

export const createOrder = (data: OrderCreateReq) =>
  request.post("/orders/create", data);

export const deleteOrder = (data: { order_id: number }) =>
  request.post("/orders/delete", data);

export const shipOrder = (data: { order_id: number; tracking_no?: string }) =>
  request.post("/orders/ship", data);

export const receiveOrder = (data: { order_id: number }) =>
  request.post("/orders/receive", data);

export const payOrder = (data: {
  order_id: number;
  money: number;
  product_id: number;
  boss_id: number;
  num?: number;
  key: string;
}) => request.post("/paydown", data);

export const requestOrderRefund = (data: {
  order_id: number;
  reason: string;
}) => request.post("/orders/refund/request", data);
