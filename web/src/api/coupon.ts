import request from "@/utils/request";

export const getCouponList = () => request.get("/coupons");

export const claimCoupon = (data: { coupon_id: number }) =>
  request.post("/coupon/claim", data);

export const getUserCouponList = () => request.get("/coupon/list");
