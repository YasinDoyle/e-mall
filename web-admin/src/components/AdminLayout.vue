<template>
  <el-container style="min-height: 100vh">
    <el-aside width="200px" style="background: #001529">
      <div class="logo">🛒 {{ appConfig.config.admin_logo_text }}</div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="#001529"
        text-color="#c0c4cc"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>数据概览
        </el-menu-item>
        <el-menu-item index="/product">
          <el-icon><Goods /></el-icon>商品管理
        </el-menu-item>
        <el-menu-item index="/seller">
          <el-icon><Shop /></el-icon>商家审核
        </el-menu-item>
        <el-menu-item index="/order">
          <el-icon><Tickets /></el-icon>订单管理
        </el-menu-item>
        <el-menu-item index="/settlement">
          <el-icon><Money /></el-icon>结算管理
        </el-menu-item>
        <el-menu-item index="/withdraw">
          <el-icon><Wallet /></el-icon>提现审核
        </el-menu-item>
        <el-menu-item index="/category">
          <el-icon><Grid /></el-icon>分类管理
        </el-menu-item>
        <el-menu-item index="/coupon">
          <el-icon><Discount /></el-icon>优惠券
        </el-menu-item>
        <el-menu-item index="/flash-sale">
          <el-icon><Timer /></el-icon>秒杀管理
        </el-menu-item>
        <el-menu-item index="/carousel">
          <el-icon><Picture /></el-icon>轮播图
        </el-menu-item>
        <el-menu-item index="/user">
          <el-icon><User /></el-icon>用户管理
        </el-menu-item>
        <el-menu-item index="/notice">
          <el-icon><Bell /></el-icon>公告管理
        </el-menu-item>
        <el-menu-item index="/notification">
          <el-icon><Message /></el-icon>消息通知
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header
        style="
          background: #fff;
          display: flex;
          align-items: center;
          justify-content: space-between;
          box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
        "
      >
        <span style="font-size: 16px; font-weight: 600">
          {{ appConfig.config.admin_site_name }}
        </span>
        <div class="header-actions">
          <el-badge
            :value="notificationStore.unreadCount"
            :hidden="!notificationStore.unreadCount"
          >
            <el-button :icon="Bell" circle @click="router.push('/notification')" />
          </el-badge>
          <el-dropdown @command="handleCommand">
            <span style="cursor: pointer"
              >{{ store.adminInfo?.nick_name || store.adminInfo?.user_name }}
              <el-icon><ArrowDown /></el-icon
            ></span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main style="background: #f0f2f5">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/stores/admin";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";
import { Bell, Message, Money, Shop, Wallet } from "@element-plus/icons-vue";

const router = useRouter();
const store = useAdminStore();
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
let unreadTimer: number | undefined;
let unreadStream: AbortController | undefined;

function handleCommand(cmd: string) {
  if (cmd === "logout") {
    store.logout();
    notificationStore.clearUnreadCount();
    window.clearInterval(unreadTimer);
    stopNotificationStream();
    router.push("/login");
  }
}

async function refreshUnreadCount() {
  if (!store.isLoggedIn) {
    notificationStore.clearUnreadCount();
    return;
  }
  try {
    await notificationStore.refreshUnreadCount();
  } catch {
    // Keep polling alive; transient notification errors should not disable the fallback.
  }
}

function startNotificationPolling() {
  window.clearInterval(unreadTimer);
  if (!store.isLoggedIn || !appConfig.config.feature_flags.notification_polling) {
    notificationStore.clearUnreadCount();
    return;
  }
  refreshUnreadCount();
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
  if (!store.isLoggedIn || !appConfig.config.feature_flags.notification_sse) return;
  const token = localStorage.getItem("admin_token");
  if (!token) return;

  const controller = new AbortController();
  unreadStream = controller;
  try {
    const response = await fetch("/api/v1/admin/notifications/stream", {
      headers: {
        access_token: token,
        refresh_token: localStorage.getItem("admin_refresh_token") ?? "",
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
  startNotificationPolling();
  startNotificationStream();
});
onBeforeUnmount(() => {
  window.clearInterval(unreadTimer);
  stopNotificationStream();
});

watch(
  () => store.token,
  () => {
    startNotificationPolling();
    startNotificationStream();
  },
);
</script>

<style scoped>
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}
.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}
</style>
