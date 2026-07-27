import request from "@/utils/request";

export const getAdminSellerWithdrawList = (params: {
  page_num: number;
  page_size: number;
  seller_id?: number;
  status?: string;
}) => request.get("/admin/seller/withdraw/list", { params });

export const auditAdminSellerWithdraw = (data: {
  id: number;
  status: string;
  reason?: string;
}) => request.post("/admin/seller/withdraw/audit", data);

export const markAdminSellerWithdrawPaid = (data: {
  id: number;
  status: string;
  reason?: string;
}) => request.post("/admin/seller/withdraw/paid", data);

export const getAdminSellerWithdrawDetail = (params: { id: number }) =>
  request.get("/admin/seller/withdraw/detail", { params });
