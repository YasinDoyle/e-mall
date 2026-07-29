import axios from "axios";
import { ElMessage } from "element-plus";
import router from "@/router";
import { ApiErrorCode, resolveApiErrorMessage } from "@/utils/api-error";
import {
  clearActiveUserSession,
  getActiveUserRefreshToken,
  getActiveUserToken,
  setActiveUserTokens,
} from "@/utils/session";

declare module "axios" {
  export interface AxiosRequestConfig {
    silentError?: boolean;
  }
}

const request = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

// 请求拦截器：自动注入 JWT token
request.interceptors.request.use(
  (config) => {
    const token = getActiveUserToken();
    const refreshToken = getActiveUserRefreshToken();
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
    if (data?.status !== undefined && data.status !== ApiErrorCode.SUCCESS) {
      const message = resolveApiErrorMessage(data);
      if (!response.config.silentError) {
        ElMessage.error(message);
      }
      return Promise.reject(new Error(message));
    }
    if (response.headers["access_token"]) {
      setActiveUserTokens(
        response.headers["access_token"],
        getActiveUserRefreshToken(),
      );
    }
    if (response.headers["refresh_token"]) {
      setActiveUserTokens(getActiveUserToken(), response.headers["refresh_token"]);
    }
    return data;
  },
  (error) => {
    if (error.response?.status === 401) {
      clearActiveUserSession();
      router.push("/login");
      ElMessage.error("登录已过期，请重新登录");
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
