import request from "@/utils/request";

export const getFlashSaleList = () => request.get("/flash_sale/list");

export const getFlashSaleDetail = (params: { product_id: number }) =>
  request.get("/flash_sale/show", { params });

export const doFlashSale = (data: {
  flash_sale_id?: number;
  product_id: number;
  boss_id: number;
  address_id: number;
  key: string;
  num: number;
  money: number;
}) => request.post("/flash_sale/skill", data);

export const getMoney = (data?: { key?: string }) =>
  request.post("/money", data ?? {});
