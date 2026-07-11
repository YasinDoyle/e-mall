// ===== 通用 =====
export interface ApiResponse<T = unknown> {
  code: number;
  data: T;
  msg: string;
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
}

export interface UserLoginReq {
  user_name: string;
  password: string;
}

export interface UserRegisterReq {
  user_name: string;
  nick_name?: string;
  password: string;
  key: string;
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
