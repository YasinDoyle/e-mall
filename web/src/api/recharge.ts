import request from "@/utils/request";

export const wechatRecharge = (data: { amount: number }) =>
  request.post("/recharge/wechat", data);

export const alipayRecharge = (data: { amount: number }) =>
  request.post("/recharge/alipay", data);

export const getRechargeStatus = (params: { order_num: string }) =>
  request.get("/recharge/status", { params });

export const getPendingCredit = () => request.get("/recharge/pending");

export const applyPendingCredit = (data: { key: string }) =>
  request.post("/recharge/apply", data);
