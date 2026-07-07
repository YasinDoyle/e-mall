import { defineStore } from "pinia";
import { ref, computed } from "vue";

export const useAdminStore = defineStore("admin", () => {
  const token = ref<string>(localStorage.getItem("admin_token") ?? "");
  const adminInfo = ref<{ user_name: string; nick_name: string } | null>(
    JSON.parse(localStorage.getItem("admin_info") ?? "null"),
  );

  const isLoggedIn = computed(() => !!token.value);

  function setToken(t: string) {
    token.value = t;
    localStorage.setItem("admin_token", t);
  }

  function setAdminInfo(info: { user_name: string; nick_name: string }) {
    adminInfo.value = info;
    localStorage.setItem("admin_info", JSON.stringify(info));
  }

  function logout() {
    token.value = "";
    adminInfo.value = null;
    localStorage.removeItem("admin_token");
    localStorage.removeItem("admin_info");
    localStorage.removeItem("admin_refresh_token");
  }

  return { token, adminInfo, isLoggedIn, setToken, setAdminInfo, logout };
});
