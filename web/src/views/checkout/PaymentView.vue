<template>
  <div class="payment-wrap">
    <h2 class="page-title">{{ t("payment.title") }}</h2>

    <el-card class="section-card">
      <template #header>{{ t("payment.orderInfo") }}</template>
      <el-empty v-if="!pendingOrders.length" :description="t('payment.emptyQueueHint')" />
      <div v-for="item in pendingOrders" :key="item.order_id" class="order-row">
        <div class="order-main">
          <div class="order-name">{{ productName(item) }}</div>
          <div class="order-meta">
            {{ t("payment.orderNo", { num: item.order_num }) }}
            <span>·</span>
            {{ t("payment.quantity", { count: item.num }) }}
          </div>
        </div>
        <span class="price">¥{{ orderAmount(item).toFixed(2) }}</span>
      </div>
      <el-divider v-if="pendingOrders.length" style="margin: 12px 0" />
      <div v-if="pendingOrders.length" class="queue-summary">
        {{ t("payment.queueCount", { count: pendingOrders.length }) }}
      </div>
      <div v-if="pendingOrders.length" class="total-row">
        <span>{{ t("payment.totalAmount") }}</span>
        <span class="total-price">¥{{ totalPrice }}</span>
      </div>
    </el-card>

    <el-card class="section-card">
      <template #header>{{ t("payment.paymentMethod") }}</template>
      <el-segmented
        v-model="payMethod"
        :options="paymentMethodOptions"
        class="payment-methods"
      />
    </el-card>

    <el-card v-if="payMethod === 'balance'" class="section-card">
      <template #header>{{ t("payment.balanceSection") }}</template>
      <div class="payment-key-row">
        <el-input
          v-model="payPassword"
          type="password"
          :placeholder="t('payment.payPasswordPlaceholder')"
          maxlength="6"
          show-password
        />
        <el-button :loading="balanceLoading" :disabled="payPassword.length !== 6" @click="loadBalance">
          {{ t("payment.refreshBalance") }}
        </el-button>
      </div>
      <div class="balance-row">
        <span>{{ balanceText }}</span>
      </div>
    </el-card>

    <el-card v-else class="section-card">
      <template #header>{{ gatewayTitle }}</template>
      <div class="gateway-panel">
        <div class="gateway-row">
          <span>{{ t("payment.currentGatewayOrder") }}</span>
          <b v-if="currentOrder">{{ t("payment.orderNo", { num: currentOrder.order_num }) }}</b>
          <b v-else>{{ t("payment.none") }}</b>
        </div>
        <div class="gateway-row">
          <span>{{ t("payment.gatewayStatus") }}</span>
          <el-tag :type="gatewayStatusTag">{{ gatewayStatusText }}</el-tag>
        </div>
        <div v-if="currentPaymentNo" class="gateway-row">
          <span>{{ t("payment.paymentNo") }}</span>
          <code>{{ currentPaymentNo }}</code>
        </div>
        <div v-if="currentQrUrl" class="qr-wrap">
          <img :src="qrImageUrl" :alt="t('payment.rechargeQrAlt')" />
          <div class="qr-link">{{ currentQrUrl }}</div>
        </div>
        <div v-if="currentPayUrl" class="gateway-link">
          <el-link :href="currentPayUrl" target="_blank" type="primary">
            {{ t("payment.openGatewayLink") }}
          </el-link>
        </div>
      </div>
    </el-card>

    <div class="pay-footer">
      <el-button size="large" @click="$router.push('/checkout')">
        {{ t("payment.backToCheckout") }}
      </el-button>
      <el-button
        type="primary"
        size="large"
        :loading="paying"
        :disabled="!canPay"
        @click="handlePay"
      >
        {{ payButtonText }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { useUserStore } from "@/stores/user";
import {
  getOrderPaymentStatus,
  payOrderByBalance,
  payOrderByAlipay,
  payOrderByWechat,
} from "@/api/order";
import { getMoney } from "@/api/flashSale";
import { activeUserSessionStorageKey } from "@/utils/session";
import { t } from "@/locales";

const router = useRouter();
const userStore = useUserStore();

type PayMethod = "balance" | "wechat" | "alipay";

const pendingOrders = ref<any[]>(
  JSON.parse(
    sessionStorage.getItem(activeUserSessionStorageKey("pending_orders")) ??
      "[]",
  ),
);
const payMethod = ref<PayMethod>("balance");
const payPassword = ref("");
const paying = ref(false);
const balance = ref("0.00");
const balanceLoaded = ref(false);
const balanceLoading = ref(false);
const currentPaymentNo = ref("");
const currentQrUrl = ref("");
const currentPayUrl = ref("");
const currentStatus = ref("");
const activeGatewayChannel = ref<PayMethod | "">("");
let pollTimer: number | undefined;
const paidOrderIds = ref<number[]>(
  JSON.parse(
    sessionStorage.getItem(activeUserSessionStorageKey("paid_order_ids")) ??
      "[]",
  ),
);

const paymentMethodOptions = computed(() => [
  { label: t("payment.methodBalance"), value: "balance" },
  { label: t("payment.methodWechat"), value: "wechat" },
  { label: t("payment.methodAlipay"), value: "alipay" },
]);

const totalPrice = computed(() =>
  pendingOrders.value
    .reduce((sum, item) => sum + orderAmount(item), 0)
    .toFixed(2),
);

const balanceText = computed(() =>
  balanceLoaded.value
    ? t("payment.currentBalance", { balance: `¥${balance.value}` })
    : t("payment.balanceAfterPassword"),
);

const currentOrder = computed(() => pendingOrders.value[0] ?? null);

const gatewayTitle = computed(() =>
  payMethod.value === "wechat"
    ? t("payment.wechatPaySection")
    : t("payment.alipayPaySection"),
);

const gatewayStatusText = computed(
  () =>
    ({
      pending: t("payment.pending"),
      paid: t("payment.paid"),
      failed: t("payment.failed"),
      closed: t("payment.closed"),
    })[currentStatus.value] ?? t("payment.pending"),
);

const gatewayStatusTag = computed(
  () =>
    ({
      pending: "warning",
      paid: "success",
      failed: "danger",
      closed: "info",
    })[currentStatus.value] ?? "info",
);

const qrImageUrl = computed(
  () =>
    `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(
      currentQrUrl.value,
    )}`,
);

const canPay = computed(
  () =>
    pendingOrders.value.length > 0 &&
    (payMethod.value === "balance" ? payPassword.value.length === 6 : true) &&
    !paying.value,
);

const payButtonText = computed(() => {
  if (!pendingOrders.value.length) {
    return t("payment.none");
  }
  if (payMethod.value === "balance") {
    return t("payment.payBalance");
  }
  if (currentPaymentNo.value && currentStatus.value === "pending") {
    return t("payment.refreshStatus");
  }
  if (currentPaymentNo.value && currentStatus.value === "paid") {
    return t("payment.nextOrderReady");
  }
  return t("payment.payGateway");
});

function productName(item: any) {
  return item.name || item.product_name || t("common.unknown");
}

function orderAmount(item: any) {
  return Number(item.money ?? item.discount_price ?? item.price ?? 0) * Number(item.num || 0);
}

function stopPolling() {
  if (pollTimer) {
    window.clearInterval(pollTimer);
    pollTimer = undefined;
  }
}

function resetGatewayState() {
  stopPolling();
  currentPaymentNo.value = "";
  currentQrUrl.value = "";
  currentPayUrl.value = "";
  currentStatus.value = "";
  activeGatewayChannel.value = "";
}

function persistPendingOrders() {
  if (pendingOrders.value.length) {
    sessionStorage.setItem(
      activeUserSessionStorageKey("pending_orders"),
      JSON.stringify(pendingOrders.value),
    );
  } else {
    sessionStorage.removeItem(activeUserSessionStorageKey("pending_orders"));
  }
  sessionStorage.setItem(
    activeUserSessionStorageKey("paid_order_ids"),
    JSON.stringify(paidOrderIds.value),
  );
}

function syncPaidCartCount(paidCount: number) {
  if (paidCount <= 0) {
    return;
  }
  userStore.setCartCount(Math.max(userStore.cartCount - paidCount, 0));
}

function shouldSyncPaidCartCount(item: any) {
  return !item?.resumed_from_order;
}

async function loadBalance() {
  if (payPassword.value.length !== 6) {
    return ElMessage.warning(t("payment.payPasswordRequired"));
  }
  balanceLoading.value = true;
  try {
    const res: any = await getMoney({ key: payPassword.value });
    balance.value = Number(res.data?.user_money ?? 0).toFixed(2);
    balanceLoaded.value = true;
  } catch (error: any) {
    balanceLoaded.value = false;
    ElMessage.error(error?.message || t("payment.balanceQueryFailed"));
  } finally {
    balanceLoading.value = false;
  }
}

async function payByBalance() {
  if (payPassword.value.length !== 6) {
    return ElMessage.warning(t("payment.payPasswordRequired"));
  }
  paying.value = true;
  const successfulOrderIds: number[] = [];
  let successfulCartSyncCount = 0;
  try {
    for (const item of pendingOrders.value) {
      await payOrderByBalance({
        order_id: item.order_id,
        money: Number(item.money ?? 0),
        product_id: item.product_id,
        boss_id: item.boss_id,
        num: item.num,
        key: payPassword.value,
      });
      successfulOrderIds.push(item.order_id);
      paidOrderIds.value.push(item.order_id);
      if (shouldSyncPaidCartCount(item)) {
        successfulCartSyncCount += 1;
      }
    }
    pendingOrders.value = [];
    persistPendingOrders();
    syncPaidCartCount(successfulCartSyncCount);
    ElMessage.success(t("payment.paymentSuccess"));
    router.push("/order/success");
  } catch (error: any) {
    if (successfulOrderIds.length) {
      pendingOrders.value = pendingOrders.value.filter(
        (item) => !successfulOrderIds.includes(item.order_id),
      );
      persistPendingOrders();
      syncPaidCartCount(successfulCartSyncCount);
      ElMessage.error(t("payment.partialPaymentFailed"));
      return;
    }
    ElMessage.error(error?.message || t("payment.paymentFailed"));
  } finally {
    paying.value = false;
  }
}

async function startGatewayPayment() {
  const item = currentOrder.value;
  if (!item) {
    return;
  }
  resetGatewayState();
  activeGatewayChannel.value = payMethod.value;
  const api = payMethod.value === "wechat" ? payOrderByWechat : payOrderByAlipay;
  const res: any = await api({ order_id: item.order_id });
  currentPaymentNo.value = res.data?.payment_no ?? "";
  currentQrUrl.value = res.data?.qr_code_url ?? "";
  currentPayUrl.value = res.data?.pay_url ?? "";
  currentStatus.value = res.data?.status ?? "pending";
  if (currentPayUrl.value) {
    window.open(currentPayUrl.value, "_blank");
  }
  if (currentPaymentNo.value) {
    pollTimer = window.setInterval(pollGatewayStatus, 2000);
  }
}

async function pollGatewayStatus() {
  if (!currentPaymentNo.value) {
    return;
  }
  const res: any = await getOrderPaymentStatus({
    payment_no: currentPaymentNo.value,
  });
  currentStatus.value = res.data?.status ?? "";
          if (currentStatus.value === "paid") {
            const paid = pendingOrders.value.shift();
            if (paid) {
              paidOrderIds.value.push(paid.order_id);
              if (shouldSyncPaidCartCount(paid)) {
                syncPaidCartCount(1);
              }
            }
            persistPendingOrders();
            stopPolling();
    if (pendingOrders.value.length) {
      ElMessage.success(t("payment.nextOrderReady"));
      await startGatewayPayment();
      return;
    }
    ElMessage.success(t("payment.paymentSuccess"));
    router.push("/order/success");
  } else if (currentStatus.value === "failed" || currentStatus.value === "closed") {
    stopPolling();
    paying.value = false;
    ElMessage.error(t("payment.paymentFailed"));
  }
}

async function handlePay() {
  if (!pendingOrders.value.length) {
    return ElMessage.error(t("payment.orderMissing"));
  }
  if (payMethod.value === "balance") {
    await payByBalance();
    return;
  }
  paying.value = true;
  try {
    if (currentPaymentNo.value && currentStatus.value === "pending") {
      await pollGatewayStatus();
      return;
    }
    await startGatewayPayment();
  } catch (error: any) {
    ElMessage.error(error?.message || t("payment.paymentFailed"));
    paying.value = false;
  }
}

watch(payMethod, (next, prev) => {
  if (
    currentPaymentNo.value &&
    currentStatus.value === "pending" &&
    activeGatewayChannel.value &&
    next !== activeGatewayChannel.value
  ) {
    payMethod.value = prev;
    ElMessage.warning(t("payment.channelLocked"));
    return;
  }
  resetGatewayState();
});

onUnmounted(() => {
  stopPolling();
});

onMounted(() => {
  if (!pendingOrders.value.length) {
    router.replace("/checkout");
  }
});
</script>

<style scoped>
.payment-wrap {
  max-width: 760px;
  margin: 0 auto;
}
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
}
.section-card {
  margin-bottom: 16px;
}
.payment-methods {
  width: 100%;
}
.order-row {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 8px 0;
}
.order-main {
  min-width: 0;
}
.order-name {
  font-weight: 500;
}
.order-meta {
  margin-top: 4px;
  color: #909399;
  font-size: 12px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
.queue-summary {
  margin-bottom: 8px;
  color: #606266;
  font-size: 12px;
}
.price {
  color: #f56c6c;
  white-space: nowrap;
}
.total-row,
.balance-row,
.gateway-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.total-price {
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
}
.payment-key-row {
  display: flex;
  gap: 10px;
  max-width: 460px;
}
.payment-key-row .el-input {
  flex: 1;
}
.balance-row {
  margin-top: 12px;
  color: #606266;
}
.gateway-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.gateway-row code {
  word-break: break-all;
}
.qr-wrap {
  text-align: center;
}
.qr-wrap img {
  width: 220px;
  height: 220px;
}
.qr-link {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
  word-break: break-all;
}
.gateway-link {
  display: flex;
  justify-content: center;
}
.pay-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
