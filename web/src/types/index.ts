// ===== 通用 =====
export interface ApiResponse<T = unknown> {
  status: number;
  data: T;
  msg: string;
  error?: string;
  track_id?: string;
}

export interface PageReq {
  page_num: number;
  page_size: number;
}

export interface DataListResp<T> {
  item: T[];
  total: number;
}

// ===== 用户 =====
export interface UserInfo {
  id: number;
  user_name: string;
  nick_name: string;
  nickname?: string;
  email: string;
  avatar: string;
  money: string;
  status: string;
  pay_key_set?: boolean;
}

export interface UserLoginReq {
  user_name: string;
  password: string;
}

export interface UserRegisterReq {
  user_name: string;
  nick_name?: string;
  email: string;
  email_code: string;
  password: string;
  password_confirm: string;
}

export interface UserUpdateReq {
  nick_name?: string;
}

// ===== 商品 =====
export interface Product {
  id: number;
  name: string;
  category_id: number;
  title: string;
  info: string;
  img_path: string;
  price: string;
  discount_price: string;
  on_sale: boolean;
  num: number;
  boss_id: number;
  boss_name: string;
  boss_avatar: string;
  audit_status?: number;
  brand?: string;
  origin?: string;
  specification?: string;
  production_date?: string;
  shelf_life?: string;
  service_guarantees?: string;
  certificate_meta?: string;
  certificates?: ProductCertificate[];
}

export interface ProductCertificate {
  id: number;
  product_id: number;
  certificate_type: string;
  name: string;
  file_path: string;
  created_at: number;
}

export interface ProductListReq extends PageReq {
  category_id?: number;
}

export interface ProductSearchReq extends PageReq {
  info: string;
}

// ===== 分类 =====
export interface Category {
  id: number;
  category_name: string;
}

// ===== 轮播图 =====
export interface Carousel {
  id: number;
  img_path: string;
  product_id: number;
}

// ===== 购物车 =====
export interface CartItem {
  id: number;
  user_id: number;
  product_id: number;
  boss_id: number;
  num: number;
  max_num: number;
  check: boolean;
  name: string;
  img_path: string;
  price: string;
  discount_price: string;
}

export interface CartCreateReq {
  product_id: number;
  boss_id: number;
  num?: number;
  max_num?: number;
}

export interface CartUpdateReq {
  id: number;
  num: number;
}

// ===== 订单 =====
export interface Order {
  id: number;
  user_id: number;
  product_id: number;
  boss_id: number;
  address_id: number;
  num: number;
  order_num: number;
  type: number;
  money: number;
  refund_status?: number;
  refund_reason?: string;
  tracking_no?: string;
  name?: string;
  img_path?: string;
  discount_price?: string;
  address_name?: string;
  address_phone?: string;
  address?: string;
  payment_channel?: string;
  logistics_company?: string;
  shipped_at?: number;
  received_at?: number;
  canceled_at?: number;
  created_at?: number;
  updated_at?: number;
  paid_at: string;
}

export interface OrderCreateReq {
  product_id: number;
  num: number;
  address_id: number;
  boss_id: number;
  money: number;
  coupon_id?: number;
}

export interface OrderLog {
  id: number;
  order_id: number;
  order_num: number;
  action: string;
  from_type: number;
  to_type: number;
  operator_type: string;
  operator_id: number;
  remark: string;
  created_at: number;
}

export interface AfterSaleRequestReq {
  order_id: number;
  type: string;
  reason: string;
}

export interface AfterSaleItem {
  id: number;
  order_id: number;
  order_num: number;
  buyer_id: number;
  seller_id: number;
  type: string;
  status: string;
  reason: string;
  refund_amount: number;
  seller_reason: string;
  platform_note: string;
  created_at: number;
  updated_at: number;
  refunded_at?: number;
  closed_at?: number;
}

// ===== 收货地址 =====
export interface Address {
  id: number;
  user_id: number;
  name: string;
  phone: string;
  address: string;
}

export interface AddressCreateReq {
  name: string;
  phone: string;
  address: string;
}

// ===== 收藏 =====
export interface Favorite {
  id: number;
  user_id: number;
  product_id: number;
  boss_id: number;
  product_name: string;
  product_img: string;
  price: string;
}

// ===== 秒杀 =====
export interface FlashSaleProduct {
  id: number;
  product_id: number;
  product_name: string;
  product_img: string;
  money: number;
  num: number;
}
