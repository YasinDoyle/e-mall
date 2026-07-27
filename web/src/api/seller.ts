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

export interface SellerAccountSummary {
  seller_id: number;
  available_balance: number;
  frozen_balance: number;
  total_income: number;
  total_withdrawn: number;
}

export interface SellerWithdraw {
  id: number;
  seller_id: number;
  user_name: string;
  nick_name: string;
  shop_name: string;
  amount: number;
  status: string;
  status_text: string;
  payee_name: string;
  payee_account: string;
  payee_channel: string;
  audit_reason: string;
  audit_operator_id: number;
  audit_operator_name: string;
  paid_operator_id: number;
  paid_operator_name: string;
  created_at: number;
  audited_at: number;
  paid_at: number;
}

export const getSellerProfile = (options?: { silentError?: boolean }) =>
  request.get("/seller/profile", { silentError: options?.silentError });

export const applySeller = (data: {
  shop_name: string;
  description: string;
}) => request.post("/seller/apply", data);

export const getSellerAccountSummary = () =>
  request.get("/seller/account/summary");

export const getSellerWithdrawList = (params: {
  page_num: number;
  page_size: number;
  seller_id?: number;
  status?: string;
}) => request.get("/seller/withdraw/list", { params });

export const applySellerWithdraw = (data: {
  amount: number;
  payee_name: string;
  payee_account: string;
  payee_channel?: string;
}) => request.post("/seller/withdraw/apply", data);
