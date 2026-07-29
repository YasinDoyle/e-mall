<template>
  <div class="navbar">
    <div class="navbar-left">
      <RouterLink to="/" class="logo">🛒 {{ appConfig.config.logo_text }}</RouterLink>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索商品..."
        style="width: 300px; margin-left: 20px"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>
    </div>

    <div class="navbar-right">
      <RouterLink to="/flash-sale" class="nav-link">秒杀专场</RouterLink>

      <template v-if="userStore.isLoggedIn">
        <RouterLink to="/cart" class="nav-link">
          <el-badge :value="userStore.cartCount" :hidden="!userStore.cartCount">
            <el-icon size="20"><ShoppingCart /></el-icon>
          </el-badge>
        </RouterLink>
        <RouterLink to="/user/notifications" class="nav-link">
          <el-badge
            :value="notificationStore.unreadCount"
            :hidden="!notificationStore.unreadCount"
          >
            <el-icon size="20"><Bell /></el-icon>
          </el-badge>
        </RouterLink>

        <el-dropdown @command="handleCommand">
          <div class="user-avatar">
            <el-avatar :size="32" :src="userStore.userInfo?.avatar || ''" />
            <span style="margin-left: 6px">{{
              userStore.userInfo?.nick_name || userStore.userInfo?.user_name
            }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="orders">我的订单</el-dropdown-item>
              <el-dropdown-item command="addresses">收货地址</el-dropdown-item>
              <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
              <el-dropdown-item command="coupons">我的优惠券</el-dropdown-item>
              <el-dropdown-item command="wallet">我的钱包</el-dropdown-item>
              <el-dropdown-item command="notifications">消息通知</el-dropdown-item>
              <el-dropdown-item command="seller">卖家中心</el-dropdown-item>
              <el-dropdown-item divided command="logout"
                >退出登录</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>

      <template v-else>
        <RouterLink to="/login">
          <el-button type="primary" size="small">登录</el-button>
        </RouterLink>
        <RouterLink to="/register" style="margin-left: 8px">
          <el-button size="small">注册</el-button>
        </RouterLink>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { Bell, Search, ShoppingCart } from "@element-plus/icons-vue";
import { useUserStore } from "@/stores/user";
import { useSellerStore } from "@/stores/seller";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";
import { getCartList } from "@/api/cart";
import { getUserInfo } from "@/api/user";
import { getActiveUserRefreshToken, getActiveUserToken } from "@/utils/session";

const router = useRouter();
const userStore = useUserStore();
const sellerStore = useSellerStore();
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
const searchKeyword = ref("");
let unreadTimer: number | undefined;
let unreadStream: AbortController | undefined;

function handleSearch() {
  if (searchKeyword.value.trim()) {
    router.push({
      path: "/search",
      query: { info: searchKeyword.value.trim() },
    });
  }
}

function handleCommand(command: string) {
  if (command === "logout") {
    userStore.logout();
    sellerStore.clearProfile();
    notificationStore.clearUnreadCount();
    router.push("/login");
  } else if (command === "seller") {
    router.push("/seller");
  } else {
    router.push(`/user/${command}`);
  }
}

async function syncUserSession() {
  if (!userStore.isLoggedIn) return;
  try {
    const [userRes, cartRes]: any[] = await Promise.all([
      getUserInfo(),
      getCartList(),
    ]);
    userStore.setUserInfo(userRes.data);
    userStore.setCartCount(cartRes.data?.item?.length ?? 0);
    refreshUnreadCount();
  } catch {
    userStore.logout();
    router.push("/login");
  }
}

async function refreshUnreadCount() {
  if (!userStore.isLoggedIn) {
    notificationStore.clearUnreadCount();
    return;
  }
  try {
    await notificationStore.refreshUnreadCount();
  } catch {
    // Keep the main session alive; the next polling tick or SSE reconnect can recover.
  }
}

function startNotificationPolling() {
  window.clearInterval(unreadTimer);
  if (!userStore.isLoggedIn || !appConfig.config.feature_flags.notification_polling) {
    notificationStore.clearUnreadCount();
    return;
  }
  unreadTimer = window.setInterval(
    refreshUnreadCount,
    appConfig.config.notification_polling_interval_ms,
  );
}

function stopNotificationStream() {
  unreadStream?.abort();
  unreadStream = undefined;
}

async function startNotificationStream() {
  stopNotificationStream();
  if (!userStore.isLoggedIn || !appConfig.config.feature_flags.notification_sse) return;
  const token = getActiveUserToken();
  if (!token) return;

  const controller = new AbortController();
  unreadStream = controller;
  try {
    const response = await fetch("/api/v1/notifications/stream", {
      headers: {
        access_token: token,
        refresh_token: getActiveUserRefreshToken(),
      },
      signal: controller.signal,
    });
    if (!response.body) return;

    const reader = response.body.getReader();
    const decoder = new TextDecoder();
    let buffer = "";
    while (!controller.signal.aborted) {
      const { value, done } = await reader.read();
      if (done) break;
      buffer += decoder.decode(value, { stream: true });
      const events = buffer.split("\n\n");
      buffer = events.pop() ?? "";
      for (const event of events) {
        const dataLine = event
          .split("\n")
          .find((line) => line.startsWith("data:"));
        if (!dataLine) continue;
        const data = JSON.parse(dataLine.replace(/^data:\s*/, ""));
        notificationStore.setUnreadCount(data.unread_count);
      }
    }
  } catch {
    if (!controller.signal.aborted) {
      startNotificationPolling();
    }
  }
}

onMounted(() => {
  syncUserSession();
  startNotificationPolling();
  startNotificationStream();
});

onBeforeUnmount(() => {
  window.clearInterval(unreadTimer);
  stopNotificationStream();
});

watch(
  () => userStore.token,
  () => {
    syncUserSession();
    startNotificationPolling();
    startNotificationStream();
  },
);
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 24px;
  background: #fff;
}
.navbar-left,
.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.logo {
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
  text-decoration: none;
}
.nav-link {
  color: #333;
  text-decoration: none;
  font-size: 14px;
}
.nav-link:hover {
  color: #409eff;
}
.user-avatar {
  display: flex;
  align-items: center;
  cursor: pointer;
}
</style>
