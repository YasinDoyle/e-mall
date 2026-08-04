import request from "@/utils/request";
import type { AfterSaleRequestReq, OrderCreateReq } from "@/types";

export const getOrderList = (params: {
  page_num: number;
  page_size: number;
  type?: number;
}) =>
  request.get("/orders/list", { params });

export const getSellerOrderList = (params: {
  page_num: number;
  page_size: number;
  type?: number;
}) => request.get("/boss/order/list", { params });

export const getSellerSettlementSummary = () =>
  request.get("/boss/settlement/summary");

export const getSellerAfterSaleList = (params: {
  order_id?: number;
  status?: string;
  type?: string;
  page_num: number;
  page_size: number;
}) => request.get("/boss/after-sales/list", { params });

export const handleSellerAfterSale = (data: {
  after_sale_id: number;
  action: string;
  reason?: string;
}) => request.post("/boss/after-sales/handle", data);

export const getOrderDetail = (params: { order_id: number }) =>
  request.get("/orders/show", { params });

export const createOrder = (data: OrderCreateReq) =>
  request.post("/orders/create", data);

export const deleteOrder = (data: { order_id: number }) =>
  request.post("/orders/delete", data);

export const cancelOrder = (data: { order_id: number }) =>
  request.post("/orders/cancel", data);

export const shipOrder = (data: {
  order_id: number;
  logistics_company?: string;
  tracking_no?: string;
}) =>
  request.post("/orders/ship", data);

export const receiveOrder = (data: { order_id: number }) =>
  request.post("/orders/receive", data);

export const payOrderByBalance = (data: {
  order_id: number;
  money: number;
  product_id: number;
  boss_id: number;
  num?: number;
  key: string;
}) => request.post("/orders/pay/balance", data);

export const payOrderByWechat = (data: { order_id: number }) =>
  request.post("/orders/pay/wechat", data);

export const payOrderByAlipay = (data: { order_id: number }) =>
  request.post("/orders/pay/alipay", data);

export const getOrderPaymentStatus = (params: { payment_no: string }) =>
  request.get("/orders/pay/status", { params });

export const getOrderLogs = (params: { order_id: number }) =>
  request.get("/orders/logs", { params });

export const requestOrderRefund = (data: {
  order_id: number;
  reason: string;
}) => request.post("/orders/refund/request", data);

export const requestAfterSale = (data: AfterSaleRequestReq) =>
  request.post("/after-sales/request", data);

export const getAfterSaleList = (params: {
  order_id?: number;
  status?: string;
  type?: string;
  page_num: number;
  page_size: number;
}) => request.get("/after-sales/list", { params });
