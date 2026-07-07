import request from "@/utils/request";

// ===== 认证 =====
export const adminLogin = (data: { user_name: string; password: string }) =>
  request.post("/user/login", data);

// ===== 分类 =====
export const getCategoryList = () => request.get("/category/list");
export const createCategory = (data: { category_name: string }) =>
  request.post("/admin/category/create", data);
export const updateCategory = (data: { id: number; category_name: string }) =>
  request.post("/admin/category/update", data);
export const deleteCategory = (data: { id: number }) =>
  request.post("/admin/category/delete", data);

// ===== 轮播图 =====
export const getCarouselList = () => request.get("/carousels");
export const createCarousel = (data: {
  img_path: string;
  product_id?: number;
}) => request.post("/admin/carousel/create", data);
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
  request.post("/product/delete", data);
