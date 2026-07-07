import request from "@/utils/request";
import type { CartCreateReq, CartUpdateReq } from "@/types";

export const getCartList = () => request.get("/carts/list");

export const createCart = (data: CartCreateReq) =>
  request.post("/carts/create", data);

export const updateCart = (data: CartUpdateReq) =>
  request.post("/carts/update", data);

export const deleteCart = (data: { id: number }) =>
  request.post("/carts/delete", data);
