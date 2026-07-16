import axios from "axios";
import { ElMessage } from "element-plus";
import { ApiErrorCode, resolveApiErrorMessage } from "@/utils/api-error";

const request = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

request.interceptors.request.use((config) => {
  const token = localStorage.getItem("admin_token");
  const refreshToken = localStorage.getItem("admin_refresh_token");
  if (token) config.headers["access_token"] = token;
  if (refreshToken) config.headers["refresh_token"] = refreshToken;
  return config;
});

request.interceptors.response.use(
  (response) => {
    const data = response.data;
    if (data?.status !== undefined && data.status !== ApiErrorCode.SUCCESS) {
      const message = resolveApiErrorMessage(data);
      ElMessage.error(message);
      if (data.status === ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_FAIL) {
        localStorage.removeItem("admin_token");
        import("@/router").then(({ default: router }) => {
          router.push("/login");
        });
      }
      return Promise.reject(new Error(message));
    }
    // 自动刷新 token
    if (response.headers["access_token"]) {
      localStorage.setItem("admin_token", response.headers["access_token"]);
    }
    if (response.headers["refresh_token"]) {
      localStorage.setItem(
        "admin_refresh_token",
        response.headers["refresh_token"],
      );
    }
    return data;
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("admin_token");
      import("@/router").then(({ default: router }) => router.push("/login"));
      ElMessage.error("登录已过期");
    } else {
      const message = resolveApiErrorMessage(
        error.response?.data,
        error.message || "网络错误",
      );
      ElMessage.error(message);
    }
    return Promise.reject(error);
  },
);

export default request;
