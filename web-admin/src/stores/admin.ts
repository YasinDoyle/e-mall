import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  clearActiveAdminSession,
  getActiveAdminInfo,
  getActiveAdminToken,
  setActiveAdminInfo,
  setActiveAdminTokens,
  type AdminInfo,
} from "@/utils/session";

export const useAdminStore = defineStore("admin", () => {
  const token = ref<string>(getActiveAdminToken());
  const adminInfo = ref<AdminInfo | null>(getActiveAdminInfo());

  const isLoggedIn = computed(() => !!token.value);

  function setToken(t: string) {
    token.value = t;
    setActiveAdminTokens(t);
  }

  function setRefreshToken(t: string) {
    setActiveAdminTokens(token.value, t);
  }

  function setAdminInfo(info: AdminInfo) {
    adminInfo.value = info;
    setActiveAdminInfo(info);
  }

  function logout() {
    token.value = "";
    adminInfo.value = null;
    clearActiveAdminSession();
  }

  return {
    token,
    adminInfo,
    isLoggedIn,
    setToken,
    setRefreshToken,
    setAdminInfo,
    logout,
  };
});
