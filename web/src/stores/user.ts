import { defineStore } from "pinia";
import { ref, computed } from "vue";
import type { UserInfo } from "@/types";

export const useUserStore = defineStore("user", () => {
  const token = ref<string>(localStorage.getItem("token") ?? "");
  const userInfo = ref<UserInfo | null>(
    JSON.parse(localStorage.getItem("userInfo") ?? "null"),
  );

  const isLoggedIn = computed(() => !!token.value);
  const cartCount = ref<number>(0);

  function setToken(t: string) {
    token.value = t;
    localStorage.setItem("token", t);
  }

  function setUserInfo(info: UserInfo) {
    userInfo.value = info;
    localStorage.setItem("userInfo", JSON.stringify(info));
  }

  function logout() {
    token.value = "";
    userInfo.value = null;
    cartCount.value = 0;
    localStorage.removeItem("token");
    localStorage.removeItem("userInfo");
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
    setUserInfo,
    logout,
    setCartCount,
  };
});
