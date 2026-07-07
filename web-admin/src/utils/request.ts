import axios from "axios";
import { ElMessage } from "element-plus";

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
    if (data.status !== undefined && data.status !== 200) {
      ElMessage.error(data.msg || "请求失败");
      return Promise.reject(new Error(data.msg));
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
      ElMessage.error(error.response?.data?.msg || error.message || "网络错误");
    }
    return Promise.reject(error);
  },
);

export default request;
