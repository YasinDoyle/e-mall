import request from "@/utils/request";

export interface SellerProfile {
  id: number;
  user_id: number;
  user_name: string;
  shop_name: string;
  description: string;
  status: number;
  status_text: string;
  reject_reason: string;
  approved_at: number;
  created_at: number;
}

export const getSellerProfile = (options?: { silentError?: boolean }) =>
  request.get("/seller/profile", { silentError: options?.silentError });

export const applySeller = (data: {
  shop_name: string;
  description: string;
}) => request.post("/seller/apply", data);
