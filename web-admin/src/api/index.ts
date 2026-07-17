import request from "@/utils/request";

// ===== 认证 =====
export const adminLogin = (data: { user_name: string; password: string }) =>
  request.post("/user/login", data);

// ===== 分类 =====
export const getCategoryList = () => request.get("/admin/category/list");
export const createCategory = (data: { category_name: string }) =>
  request.post("/admin/category/create", data);
export const updateCategory = (data: { id: number; category_name: string }) =>
  request.post("/admin/category/update", data);
export const deleteCategory = (data: { id: number }) =>
  request.post("/admin/category/delete", data);

// ===== 轮播图 =====
export const getCarouselList = () => request.get("/admin/carousel/list");
export const createCarousel = (data: {
  img_path: string;
  product_id: number;
}) => request.post("/admin/carousel/create", data);
export const uploadCarouselImage = (formData: FormData) =>
  request.post("/admin/carousel/upload", formData);
export const updateCarousel = (data: {
  id: number;
  img_path: string;
  product_id: number;
}) => request.post("/admin/carousel/update", data);
export const deleteCarousel = (data: { id: number }) =>
  request.post("/admin/carousel/delete", data);

// ===== 公告 =====
export const getNoticeList = () => request.get("/admin/notice/list");
export const createNotice = (data: { text: string }) =>
  request.post("/admin/notice/create", data);
export const updateNotice = (data: { id: number; text: string }) =>
  request.post("/admin/notice/update", data);
export const deleteNotice = (data: { id: number }) =>
  request.post("/admin/notice/delete", data);

// ===== 用户管理 =====
export const getUserList = (params: { page_num: number; page_size: number }) =>
  request.get("/admin/user/list", { params });
export const banUser = (data: { id: number; banned: boolean }) =>
  request.post("/admin/user/ban", data);

// ===== 商品管理 =====
export const getAdminProductList = (params: {
  page_num: number;
  page_size: number;
  audit_status?: number;
}) => request.get("/admin/product/list", { params });
export const auditProduct = (data: { id: number; audit_status: number }) =>
  request.post("/admin/product/audit", data);
export const deleteProduct = (data: { id: number }) =>
  request.post("/admin/product/delete", data);

// ===== 统计 =====
export const getStatsOverview = () => request.get("/admin/stats/overview");
export const getStatsOrders = (params?: {
  start_date?: string;
  end_date?: string;
}) => request.get("/admin/stats/orders", { params });

// ===== 订单管理 =====
export const getAdminOrderList = (params: {
  page_num: number;
  page_size: number;
  type?: number;
  refund_status?: number;
}) => request.get("/admin/orders/list", { params });
export const approveOrderRefund = (data: { order_id: number; key: string }) =>
  request.post("/admin/orders/refund/approve", data);

// ===== 优惠券管理 =====
export const getAdminCouponList = () => request.get("/admin/coupon/list");
export const createAdminCoupon = (data: {
  name: string;
  coupon_type: number;
  discount: number;
  min_amount: number;
  stock: number;
  expire_at: string;
}) => request.post("/admin/coupon/create", data);
export const offlineAdminCoupon = (data: { id: number }) =>
  request.post("/admin/coupon/offline", data);

// ===== 秒杀管理 =====
export const getAdminFlashSaleList = (params: {
  page_num: number;
  page_size: number;
}) => request.get("/admin/flash-sale/list", { params });
export const createAdminFlashSale = (data: {
  product_id: number;
  boss_id: number;
  title: string;
  money: number;
  num: number;
  custom_id?: number;
  custom_name?: string;
}) => request.post("/admin/flash-sale/create", data);
export const updateAdminFlashSale = (
  data: {
    id: number;
    product_id: number;
    boss_id: number;
    title: string;
    money: number;
    num: number;
    custom_id?: number;
    custom_name?: string;
  },
) => request.post("/admin/flash-sale/update", data);
export const deleteAdminFlashSale = (data: { id: number }) =>
  request.post("/admin/flash-sale/delete", data);
