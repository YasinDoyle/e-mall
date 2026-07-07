import request from "@/utils/request";
import type { AddressCreateReq } from "@/types";

export const getAddressList = () => request.get("/addresses/list");

export const getAddressDetail = (params: { id: number }) =>
  request.get("/addresses/show", { params });

export const createAddress = (data: AddressCreateReq) =>
  request.post("/addresses/create", data);

export const updateAddress = (data: AddressCreateReq & { id: number }) =>
  request.post("/addresses/update", data);

export const deleteAddress = (data: { id: number }) =>
  request.post("/addresses/delete", data);
