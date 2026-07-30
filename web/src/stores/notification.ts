import { defineStore } from "pinia";
import { ref } from "vue";
import { getNotificationUnreadCount } from "@/api/notification";

export const useNotificationStore = defineStore("notification", () => {
  const unreadCount = ref(0);

  function setUnreadCount(value: number) {
    unreadCount.value = Math.max(0, Number(value) || 0);
  }

  function decrementUnreadCount(count = 1) {
    setUnreadCount(unreadCount.value - count);
  }

  function clearUnreadCount() {
    unreadCount.value = 0;
  }

  async function refreshUnreadCount() {
    const res: any = await getNotificationUnreadCount();
    setUnreadCount(res.data?.unread_count);
  }

  return {
    unreadCount,
    setUnreadCount,
    decrementUnreadCount,
    clearUnreadCount,
    refreshUnreadCount,
  };
});
