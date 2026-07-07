import request from "@/utils/request";
import type { OrderCreateReq } from "@/types";

export const getOrderList = (params: { page_num: number; page_size: number }) =>
  request.get("/orders/list", { params });

export const getOrderDetail = (params: { id: number }) =>
  request.get("/orders/show", { params });

export const createOrder = (data: OrderCreateReq) =>
  request.post("/orders/create", data);

export const deleteOrder = (data: { id: number }) =>
  request.post("/orders/delete", data);

export const shipOrder = (data: { id: number }) =>
  request.post("/orders/ship", data);

export const receiveOrder = (data: { id: number }) =>
  request.post("/orders/receive", data);

export const payOrder = (data: {
  order_id: number;
  payment_password: string;
  money: number;
  boss_id: number;
  address_id: number;
  order_num: number;
}) => request.post("/paydown", data);
