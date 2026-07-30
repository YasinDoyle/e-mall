import axios from "axios";
import { ElMessage } from "element-plus";
import { ApiErrorCode, resolveApiErrorMessage } from "@/utils/api-error";
import {
  clearActiveAdminSession,
  getActiveAdminRefreshToken,
  getActiveAdminToken,
  setActiveAdminTokens,
} from "@/utils/session";
import { getCurrentLocale, t } from "@/locales";

declare module "axios" {
  export interface AxiosRequestConfig {
    silentError?: boolean;
  }
}

const request = axios.create({
  baseURL: "/api/v1",
  timeout: 10000,
});

request.interceptors.request.use((config) => {
  const token = getActiveAdminToken();
  const refreshToken = getActiveAdminRefreshToken();
  if (token) config.headers["access_token"] = token;
  if (refreshToken) config.headers["refresh_token"] = refreshToken;
  const locale = getCurrentLocale();
  config.headers["X-Locale"] = locale;
  config.headers["Accept-Language"] = locale;
  return config;
});

request.interceptors.response.use(
  (response) => {
    const data = response.data;
    if (data?.status !== undefined && data.status !== ApiErrorCode.SUCCESS) {
      const message = resolveApiErrorMessage(data);
      if (!response.config.silentError) {
        ElMessage.error(message);
      }
      if (data.status === ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_FAIL) {
        clearActiveAdminSession();
        import("@/router").then(({ default: router }) => {
          router.push("/login");
        });
      }
      return Promise.reject(new Error(message));
    }
    // 自动刷新 token
    if (response.headers["access_token"]) {
      setActiveAdminTokens(
        response.headers["access_token"],
        getActiveAdminRefreshToken(),
      );
    }
    if (response.headers["refresh_token"]) {
      setActiveAdminTokens(
        getActiveAdminToken(),
        response.headers["refresh_token"],
      );
    }
    return data;
  },
  (error) => {
    if (error.response?.status === 401) {
      clearActiveAdminSession();
      import("@/router").then(({ default: router }) => router.push("/login"));
      ElMessage.error(t("common.loginExpired", "登录已过期"));
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
