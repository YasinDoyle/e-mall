import request from "@/utils/request";

export const getAdminSellerList = (params: {
  page_num: number;
  page_size: number;
  status?: number;
}) => request.get("/admin/seller/list", { params });

export const auditAdminSeller = (data: {
  id: number;
  status: number;
  reject_reason?: string;
}) => request.post("/admin/seller/audit", data);
