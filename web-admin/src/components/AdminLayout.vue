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
          <el-icon><DataLine /></el-icon>{{ $t("admin.menu.dashboard") }}
        </el-menu-item>
        <el-menu-item index="/product">
          <el-icon><Goods /></el-icon>
          <span class="menu-label">
            <span>{{ $t("admin.menu.product") }}</span>
            <el-badge
              v-if="pendingMenuCounts.product"
              class="pending-menu-badge"
              :value="pendingMenuCounts.product"
              :max="99"
            />
          </span>
        </el-menu-item>
        <el-menu-item index="/seller">
          <el-icon><Shop /></el-icon>
          <span class="menu-label">
            <span>{{ $t("admin.menu.seller") }}</span>
            <el-badge
              v-if="pendingMenuCounts.seller"
              class="pending-menu-badge"
              :value="pendingMenuCounts.seller"
              :max="99"
            />
          </span>
        </el-menu-item>
        <el-menu-item index="/order">
          <el-icon><Tickets /></el-icon>
          <span class="menu-label">
            <span>{{ $t("admin.menu.order") }}</span>
            <el-badge
              v-if="pendingMenuCounts.order"
              class="pending-menu-badge"
              :value="pendingMenuCounts.order"
              :max="99"
            />
          </span>
        </el-menu-item>
        <el-menu-item index="/settlement">
          <el-icon><Money /></el-icon>
          <span class="menu-label">
            <span>{{ $t("admin.menu.settlement") }}</span>
            <el-badge
              v-if="pendingMenuCounts.settlement"
              class="pending-menu-badge"
              :value="pendingMenuCounts.settlement"
              :max="99"
            />
          </span>
        </el-menu-item>
        <el-menu-item index="/withdraw">
          <el-icon><Wallet /></el-icon>
          <span class="menu-label">
            <span>{{ $t("admin.menu.withdraw") }}</span>
            <el-badge
              v-if="pendingMenuCounts.withdraw"
              class="pending-menu-badge"
              :value="pendingMenuCounts.withdraw"
              :max="99"
            />
          </span>
        </el-menu-item>
        <el-menu-item index="/category">
          <el-icon><Grid /></el-icon>{{ $t("admin.menu.category") }}
        </el-menu-item>
        <el-menu-item index="/coupon">
          <el-icon><Discount /></el-icon>{{ $t("admin.menu.coupon") }}
        </el-menu-item>
        <el-menu-item index="/flash-sale">
          <el-icon><Timer /></el-icon>{{ $t("admin.menu.flashSale") }}
        </el-menu-item>
        <el-menu-item index="/carousel">
          <el-icon><Picture /></el-icon>{{ $t("admin.menu.carousel") }}
        </el-menu-item>
        <el-menu-item index="/user">
          <el-icon><User /></el-icon>{{ $t("admin.menu.user") }}
        </el-menu-item>
        <el-menu-item index="/notice">
          <el-icon><Bell /></el-icon>{{ $t("admin.menu.notice") }}
        </el-menu-item>
        <el-menu-item index="/notification">
          <el-icon><Message /></el-icon>{{ $t("admin.menu.notification") }}
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
          <el-badge
            :value="notificationStore.unreadCount"
            :hidden="!notificationStore.unreadCount"
          >
            <el-button :icon="Bell" circle @click="router.push('/notification')" />
          </el-badge>
          <el-dropdown trigger="click" @command="handleCommand">
            <span class="admin-account-trigger">
              <el-avatar :size="28" :icon="UserFilled" />
              <span class="admin-account-name">
                {{ store.adminInfo?.nick_name || store.adminInfo?.user_name }}
              </span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="switchAccount">{{
                  $t("common.accountSwitchTitle")
                }}</el-dropdown-item>
                <el-dropdown-item command="logout">{{ $t("common.logout") }}</el-dropdown-item>
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

  <el-dialog
    v-model="accountSwitchVisible"
    :title="$t('common.accountSwitchTitle')"
    width="420px"
  >
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
              {{ $t("common.currentAccount") }}
            </el-tag>
          </div>
          <div class="account-subtitle">{{ account.subtitle }}</div>
        </div>
        <el-button
          type="primary"
          size="small"
          :disabled="account.active"
          @click="handleSelectAccount(account.id)"
        >
          {{ $t("common.useAccount") }}
        </el-button>
      </el-card>
    </div>
    <template #footer>
      <el-button @click="accountSwitchVisible = false">{{ $t("common.cancel") }}</el-button>
      <el-button type="primary" @click="handleAddAccount">{{ $t("common.addAccount") }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import {
  getAdminOrderList,
  getAdminProductList,
} from "@/api";
import { getAdminSellerList } from "@/api/seller";
import { getAdminSettlementList } from "@/api/settlement";
import { getAdminSellerWithdrawList } from "@/api/withdraw";
import { useAdminStore } from "@/stores/admin";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";
import { ADMIN_PENDING_COUNTS_REFRESH_EVENT } from "@/utils/adminPending";
import {
  activateAdminSession,
  getActiveAdminRefreshToken,
  getActiveAdminToken,
  listSavedAdminSessions,
} from "@/utils/session";
import {
  currentLocale,
  getCurrentLocale,
  setLocale,
  supportedLocales,
} from "@/locales";
import { Bell, Message, Money, Shop, UserFilled, Wallet } from "@element-plus/icons-vue";

const router = useRouter();
const store = useAdminStore();
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
const accountSwitchVisible = ref(false);
const pendingMenuCounts = reactive({
  product: 0,
  seller: 0,
  order: 0,
  settlement: 0,
  withdraw: 0,
});
const selectedLocale = computed({
  get: () => currentLocale.value,
  set: (locale: string) => setLocale(locale),
});
let unreadTimer: number | undefined;
let unreadStream: AbortController | undefined;

function handleCommand(cmd: string) {
  if (cmd === "logout") {
    store.logout();
    notificationStore.clearUnreadCount();
    window.clearInterval(unreadTimer);
    stopNotificationStream();
    router.push("/login");
  } else if (cmd === "switchAccount") {
    accountSwitchVisible.value = true;
  }
}

const savedAccounts = computed(() => listSavedAdminSessions());

function handleSelectAccount(id: string) {
  if (!activateAdminSession(id)) return;
  accountSwitchVisible.value = false;
  window.location.reload();
}

function handleAddAccount() {
  accountSwitchVisible.value = false;
  store.logout();
  notificationStore.clearUnreadCount();
  window.clearInterval(unreadTimer);
  stopNotificationStream();
  router.push({ path: "/login", query: { switch: "1" } });
}

function listTotal(res: any) {
  return Number(res?.data?.total ?? 0);
}

async function refreshPendingMenuCounts() {
  if (!store.isLoggedIn) {
    pendingMenuCounts.product = 0;
    pendingMenuCounts.seller = 0;
    pendingMenuCounts.order = 0;
    pendingMenuCounts.settlement = 0;
    pendingMenuCounts.withdraw = 0;
    return;
  }
  const [product, seller, order, settlement, withdraw] =
    await Promise.allSettled([
      getAdminProductList({ page_num: 1, page_size: 1, audit_status: 0 }),
      getAdminSellerList({ page_num: 1, page_size: 1, status: 0 }),
      getAdminOrderList({ page_num: 1, page_size: 1, refund_status: 1 }),
      getAdminSettlementList({ page_num: 1, page_size: 1, status: "pending" }),
      getAdminSellerWithdrawList({ page_num: 1, page_size: 1, status: "pending" }),
    ]);

  if (product.status === "fulfilled") pendingMenuCounts.product = listTotal(product.value);
  if (seller.status === "fulfilled") pendingMenuCounts.seller = listTotal(seller.value);
  if (order.status === "fulfilled") pendingMenuCounts.order = listTotal(order.value);
  if (settlement.status === "fulfilled") pendingMenuCounts.settlement = listTotal(settlement.value);
  if (withdraw.status === "fulfilled") pendingMenuCounts.withdraw = listTotal(withdraw.value);
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
  window.clearInterval(unreadTimer);
  if (!store.isLoggedIn) {
    notificationStore.clearUnreadCount();
    return;
  }
  if (!appConfig.config.feature_flags.notification_sse) {
    startNotificationPolling();
    return;
  }
  const token = getActiveAdminToken();
  if (!token) return;

  const controller = new AbortController();
  unreadStream = controller;
  try {
    const locale = getCurrentLocale();
    const response = await fetch("/api/v1/admin/notifications/stream", {
      headers: {
        access_token: token,
        refresh_token: getActiveAdminRefreshToken(),
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
  refreshPendingMenuCounts();
  window.addEventListener(
    ADMIN_PENDING_COUNTS_REFRESH_EVENT,
    refreshPendingMenuCounts,
  );
  startNotificationStream();
});
onBeforeUnmount(() => {
  window.removeEventListener(
    ADMIN_PENDING_COUNTS_REFRESH_EVENT,
    refreshPendingMenuCounts,
  );
  window.clearInterval(unreadTimer);
  stopNotificationStream();
});

watch(
  () => store.token,
  () => {
    refreshPendingMenuCounts();
    startNotificationStream();
  },
);

watch(currentLocale, () => {
  startNotificationStream();
});
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
.menu-label {
  display: inline-flex;
  align-items: center;
  justify-content: space-between;
  min-width: 0;
  width: 100%;
  gap: 8px;
}
.pending-menu-badge {
  margin-left: auto;
  line-height: 1;
}
.pending-menu-badge :deep(.el-badge__content) {
  border: 0;
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
  gap: 12px;
}
.account-meta {
  min-width: 0;
}
.account-title {
  display: flex;
  align-items: center;
  gap: 8px;
}
.account-subtitle {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}
.locale-select {
  width: 104px;
}
.admin-account-trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  line-height: 1;
}
.admin-account-name {
  max-width: 140px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
