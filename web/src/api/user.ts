import request from "@/utils/request";
import type { UserLoginReq, UserRegisterReq, UserUpdateReq } from "@/types";

export const userLogin = (data: UserLoginReq) =>
  request.post("/user/login", data);

export const userRegister = (data: UserRegisterReq) =>
  request.post("/user/register", data);

export const getUserInfo = () => request.get("/user/show_info");

export const updateUserInfo = (data: UserUpdateReq) =>
  request.post("/user/update", data);

export const uploadAvatar = (formData: FormData) =>
  request.post("/user/avatar", formData, {
    headers: { "Content-Type": "multipart/form-data" },
  });

export const sendEmail = (data: {
  email: string;
  password?: string;
  operation_type: 1 | 2 | 3;
}) =>
  request.post("/user/send_email", data);

export const validEmail = (params: { token: string }) =>
  request.get("/user/valid_email", { params });

export const userFollow = (data: { id: number }) =>
  request.post("/user/following", data);

export const userUnFollow = (data: { id: number }) =>
  request.post("/user/unfollowing", data);
