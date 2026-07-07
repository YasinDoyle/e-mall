import request from "@/utils/request";

export const getFlashSaleList = () => request.get("/flash_sale/list");

export const getFlashSaleDetail = (params: { id: number }) =>
  request.get("/flash_sale/show", { params });

export const doFlashSale = (data: {
  product_id: number;
  boss_id: number;
  address_id: number;
  key: string;
  num: number;
  money: number;
}) => request.post("/flash_sale/skill", data);

export const getMoney = () => request.post("/money", {});
