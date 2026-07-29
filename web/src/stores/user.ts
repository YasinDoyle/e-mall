import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { UserInfo } from "@/types";
import {
  clearActiveUserSession,
  getActiveUserInfo,
  getActiveUserToken,
  setActiveUserInfo,
  setActiveUserTokens,
} from "@/utils/session";

export const useUserStore = defineStore("user", () => {
  const token = ref<string>(getActiveUserToken());
  const userInfo = ref<UserInfo | null>(getActiveUserInfo());

  const isLoggedIn = computed(() => !!token.value);
  const cartCount = ref<number>(0);

  function setToken(t: string) {
    token.value = t;
    setActiveUserTokens(t);
  }

  function setRefreshToken(t: string) {
    setActiveUserTokens(token.value, t);
  }

  function setUserInfo(info: UserInfo) {
    const normalized = {
      ...info,
      nick_name: info.nick_name || info.nickname || info.user_name,
    };
    userInfo.value = normalized;
    setActiveUserInfo(normalized);
  }

  function logout() {
    token.value = "";
    userInfo.value = null;
    cartCount.value = 0;
    clearActiveUserSession();
  }

  function setCartCount(n: number) {
    cartCount.value = n;
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    cartCount,
    setToken,
    setRefreshToken,
    setUserInfo,
    logout,
    setCartCount,
  };
});
