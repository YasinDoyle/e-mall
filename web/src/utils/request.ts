import axios from "axios";
import { ElMessage } from "element-plus";
import router from "@/router";

const request = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

// 请求拦截器：自动注入 JWT token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    const refreshToken = localStorage.getItem("refreshToken");
    if (token) {
      config.headers["access_token"] = token;
    }
    if (refreshToken) {
      config.headers["refresh_token"] = refreshToken;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => {
    const data = response.data;
    // 后端业务错误（status 非 200），backend 使用 status 字段
    if (data.status !== undefined && data.status !== 200) {
      ElMessage.error(data.msg || "请求失败");
      return Promise.reject(new Error(data.msg || "请求失败"));
    }
    if (response.headers["access_token"]) {
      localStorage.setItem("token", response.headers["access_token"]);
    }
    if (response.headers["refresh_token"]) {
      localStorage.setItem("refreshToken", response.headers["refresh_token"]);
    }
    return data;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("userInfo");
      router.push("/login");
      ElMessage.error("登录已过期，请重新登录");
    } else {
      ElMessage.error(error.response?.data?.msg || error.message || "网络错误");
    }
    return Promise.reject(error);
  },
);

export default request;
