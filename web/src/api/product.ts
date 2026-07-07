import request from "@/utils/request";
import type { ProductListReq, ProductSearchReq } from "@/types";

export const getProductList = (params: ProductListReq) =>
  request.get("/product/list", { params });

export const getProductDetail = (params: { id: number }) =>
  request.get("/product/show", { params });

export const searchProducts = (data: ProductSearchReq) =>
  request.post("/product/search", data);

export const getProductImgs = (params: { id: number }) =>
  request.get("/product/imgs/list", { params });

export const createProduct = (formData: FormData) =>
  request.post("/product/create", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });

export const updateProduct = (formData: FormData) =>
  request.post("/product/update", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });

export const deleteProduct = (data: { id: number }) =>
  request.post("/product/delete", data);

export const getCategoryList = () => request.get("/category/list");

export const getCarousels = () => request.get("/carousels");
