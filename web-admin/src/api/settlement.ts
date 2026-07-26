import request from "@/utils/request";

export const getAdminSettlementList = (params: {
  page_num: number;
  page_size: number;
  seller_id?: number;
  status?: string;
}) => request.get("/admin/settlement/list", { params });

export const generateAdminSettlement = (data: { seller_id: number }) =>
  request.post("/admin/settlement/generate", data);

export const generateOneAdminSettlement = (data: { id: number }) =>
  request.post("/admin/settlement/generate_one", data);

export const markAdminSettlementPaid = (data: { id: number }) =>
  request.post("/admin/settlement/paid", data);

export const getAdminSettlementDetail = (params: { id: number }) =>
  request.get("/admin/settlement/detail", { params });
