<template>
  <div class="navbar">
    <div class="navbar-left">
      <RouterLink to="/" class="logo">🛒 {{ appConfig.config.logo_text }}</RouterLink>
      <el-input
        v-model="searchKeyword"
        :placeholder="$t('nav.searchPlaceholder')"
        style="width: 300px; margin-left: 20px"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>
    </div>

    <div class="navbar-right">
      <el-select
        v-model="selectedLocale"
        class="locale-select"
        size="small"
        :aria-label="$t('common.language')"
      >
        <el-option
          v-for="locale in supportedLocales"
          :key="locale.value"
          :label="locale.label"
          :value="locale.value"
        />
      </el-select>
      <RouterLink to="/flash-sale" class="nav-link">{{ $t("nav.flashSale") }}</RouterLink>

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

        <el-dropdown trigger="click" @command="handleCommand">
          <div class="user-avatar">
            <el-avatar :size="32" :src="userStore.userInfo?.avatar || ''" />
            <span style="margin-left: 6px">{{
              userStore.userInfo?.nick_name || userStore.userInfo?.user_name
            }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">{{ $t("nav.profile") }}</el-dropdown-item>
              <el-dropdown-item command="orders">{{ $t("nav.orders") }}</el-dropdown-item>
              <el-dropdown-item command="addresses">{{ $t("nav.addresses") }}</el-dropdown-item>
              <el-dropdown-item command="favorites">{{ $t("nav.favorites") }}</el-dropdown-item>
              <el-dropdown-item command="coupons">{{ $t("nav.coupons") }}</el-dropdown-item>
              <el-dropdown-item command="wallet">{{ $t("nav.wallet") }}</el-dropdown-item>
              <el-dropdown-item command="notifications">{{ $t("nav.notifications") }}</el-dropdown-item>
              <el-dropdown-item command="seller">{{ $t("nav.sellerCenter") }}</el-dropdown-item>
              <el-dropdown-item command="switchAccount">{{ $t("nav.switchAccount") }}</el-dropdown-item>
              <el-dropdown-item divided command="logout"
                >{{ $t("common.logout") }}</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>

      <template v-else>
        <RouterLink to="/login">
          <el-button type="primary" size="small">{{ $t("common.login") }}</el-button>
        </RouterLink>
        <RouterLink to="/register" style="margin-left: 8px">
          <el-button size="small">{{ $t("common.register") }}</el-button>
        </RouterLink>
      </template>
    </div>
  </div>

  <el-dialog v-model="accountSwitchVisible" :title="$t('common.accountSwitchTitle')" width="420px">
    <div class="account-switch-list">
      <el-empty v-if="!savedAccounts.length" :description="$t('common.noSavedAccounts')" />
      <el-card
        v-for="account in savedAccounts"
        :key="account.id"
        class="account-switch-item"
        :shadow="account.active ? 'always' : 'hover'"
      >
        <div class="account-meta">
          <div class="account-title">
            <b>{{ account.title }}</b>
            <el-tag v-if="account.active" size="small" type="success">
              {{ $t('common.currentAccount') }}
            </el-tag>
          </div>
          <div class="account-subtitle">{{ account.subtitle }}</div>
        </div>
        <el-button type="primary" size="small" :disabled="account.active" @click="handleSelectAccount(account.id)">
          {{ $t('common.useAccount') }}
        </el-button>
      </el-card>
    </div>
    <template #footer>
      <el-button @click="accountSwitchVisible = false">{{ $t('common.cancel') }}</el-button>
      <el-button type="primary" @click="handleAddAccount">{{ $t('common.addAccount') }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { Bell, Search, ShoppingCart } from "@element-plus/icons-vue";
import { useUserStore } from "@/stores/user";
import { useSellerStore } from "@/stores/seller";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";
import { getCartList } from "@/api/cart";
import { getUserInfo } from "@/api/user";
import {
  activateUserSession,
  getActiveUserRefreshToken,
  getActiveUserToken,
  listSavedUserSessions,
} from "@/utils/session";
import {
  currentLocale,
  getCurrentLocale,
  setLocale,
  supportedLocales,
} from "@/locales";

const router = useRouter();
const userStore = useUserStore();
const sellerStore = useSellerStore();
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
const searchKeyword = ref("");
const accountSwitchVisible = ref(false);
const selectedLocale = computed({
  get: () => currentLocale.value,
  set: (locale: string) => setLocale(locale),
});
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
  } else if (command === "switchAccount") {
    accountSwitchVisible.value = true;
  } else if (command === "seller") {
    router.push("/seller");
  } else {
    router.push(`/user/${command}`);
  }
}

const savedAccounts = computed(() => listSavedUserSessions());

function handleSelectAccount(id: string) {
  if (!activateUserSession(id)) return;
  accountSwitchVisible.value = false;
  window.location.assign("/");
}

function handleAddAccount() {
  accountSwitchVisible.value = false;
  userStore.logout();
  sellerStore.clearProfile();
  notificationStore.clearUnreadCount();
  router.push({ path: "/login", query: { switch: "1" } });
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
  window.clearInterval(unreadTimer);
  if (!userStore.isLoggedIn) {
    notificationStore.clearUnreadCount();
    return;
  }
  if (!appConfig.config.feature_flags.notification_sse) {
    startNotificationPolling();
    return;
  }
  const token = getActiveUserToken();
  if (!token) return;

  const controller = new AbortController();
  unreadStream = controller;
  try {
    const locale = getCurrentLocale();
    const response = await fetch("/api/v1/notifications/stream", {
      headers: {
        access_token: token,
        refresh_token: getActiveUserRefreshToken(),
        "X-Locale": locale,
        "Accept-Language": locale,
      },
      signal: controller.signal,
    });
    if (!response.body) {
      startNotificationPolling();
      return;
    }

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
    startNotificationStream();
  },
);

watch(currentLocale, () => {
  startNotificationStream();
});
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
.locale-select {
  width: 104px;
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
.account-switch-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.account-switch-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.account-meta {
  display: flex;
  flex-direction: column;
  gap: 4px;
}
.account-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.account-subtitle {
  color: #909399;
  font-size: 12px;
}
</style>
