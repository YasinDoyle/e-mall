<template>
  <div class="payment-wrap">
    <h2 class="page-title">{{ t("payment.title") }}</h2>

    <el-card class="section-card">
      <template #header>{{ t("payment.orderInfo") }}</template>
      <el-empty v-if="!pendingOrders.length" :description="t('payment.noPendingOrders')" />
      <div v-for="item in pendingOrders" :key="item.order_id" class="order-row">
        <span>{{ productName(item) }} × {{ item.num }}</span>
        <span class="price"
          >¥{{ (unitPriceValue(item) * item.num).toFixed(2) }}</span
        >
      </div>
      <el-divider v-if="pendingOrders.length" style="margin: 12px 0" />
      <div v-if="pendingOrders.length" class="total-row">
        <span>{{ t("payment.totalAmount") }}</span>
        <span class="total-price">¥{{ totalPrice }}</span>
      </div>
    </el-card>

    <el-card class="section-card">
      <template #header>{{ t("payment.paymentMethod") }}</template>
      <el-radio-group
        v-model="payMethod"
        style="display: flex; flex-direction: column; gap: 12px"
      >
        <el-radio value="balance">
          <span>{{ t("payment.balancePay") }}</span>
          <span style="color: #999; font-size: 12px; margin-left: 8px">
            ({{ balanceText }})
          </span>
        </el-radio>
        <el-radio value="wechat">{{ t("payment.wechatRechargePay") }}</el-radio>
        <el-radio value="alipay">{{ t("payment.alipayRechargePay") }}</el-radio>
      </el-radio-group>
    </el-card>

    <el-card v-if="payMethod === 'balance'" class="section-card">
      <template #header>{{ t("payment.payPassword") }}</template>
      <div class="payment-key-row">
        <el-input
          v-model="payPassword"
          type="password"
          :placeholder="t('payment.payPasswordPlaceholder')"
          maxlength="6"
          show-password
        />
        <el-button
          :loading="balanceLoading"
          :disabled="payPassword.length !== 6"
          @click="loadBalance"
        >
          {{ t("payment.refreshBalance") }}
        </el-button>
      </div>
    </el-card>

    <el-card v-else class="section-card">
      <template #header>
        {{ payMethod === "wechat" ? t("payment.wechatRecharge") : t("payment.alipayRecharge") }}
      </template>
      <div class="recharge-panel">
        <div class="recharge-row">
          <span>{{ t("payment.suggestedRecharge") }}</span>
          <el-input-number
            v-model="rechargeAmount"
            :min="0.01"
            :precision="2"
            :step="10"
          />
        </div>
        <el-button
          type="primary"
          :loading="rechargeLoading"
          @click="startRecharge"
        >
          {{ rechargeOrderNum ? t("payment.rechargeAgain") : t("payment.startRecharge") }}
        </el-button>
      </div>

      <div v-if="rechargeOrderNum" class="recharge-result">
        <div class="result-row">
          <span>{{ t("payment.rechargeNo") }}</span>
          <b>{{ rechargeOrderNum }}</b>
        </div>
        <div class="result-row">
          <span>{{ t("payment.paymentStatus") }}</span>
          <el-tag :type="rechargeStatusTag">{{ rechargeStatusText }}</el-tag>
        </div>
        <div v-if="wechatQrUrl" class="qr-wrap">
          <img :src="qrImageUrl" :alt="t('payment.rechargeQrAlt')" />
          <div class="qr-link">{{ wechatQrUrl }}</div>
        </div>
        <div v-if="pendingCredit > 0" class="apply-box">
          <el-input
            v-model="payPassword"
            type="password"
            maxlength="6"
            show-password
            :placeholder="t('payment.confirmCreditPassword')"
          />
          <el-button
            type="success"
            :loading="applyingCredit"
            :disabled="payPassword.length !== 6"
            @click="handleApplyCredit"
          >
            {{ t("payment.confirmCredit") }}
          </el-button>
        </div>
      </div>
    </el-card>

    <div class="pay-footer">
      <el-button size="large" @click="$router.push('/cart')">{{ t("payment.cancel") }}</el-button>
      <el-button
        type="primary"
        size="large"
        :loading="paying"
        :disabled="!canPay"
        @click="handlePay"
      >
        {{ t("payment.payNow", { amount: totalPrice }) }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { payOrder } from "@/api/order";
import { getMoney } from "@/api/flashSale";
import {
  alipayRecharge,
  applyPendingCredit,
  getRechargeStatus,
  wechatRecharge,
} from "@/api/recharge";
import { useUserStore } from "@/stores/user";
import { activeUserSessionStorageKey } from "@/utils/session";
import { t } from "@/locales";

const router = useRouter();
const userStore = useUserStore();

const pendingOrders = ref<any[]>(
  JSON.parse(
    sessionStorage.getItem(activeUserSessionStorageKey("pending_orders")) ??
      "[]",
  ),
);
const payMethod = ref("balance");
const payPassword = ref("");
const paying = ref(false);
const balance = ref("0.00");
const balanceLoaded = ref(false);
const balanceLoading = ref(false);
const rechargeAmount = ref(0);
const rechargeLoading = ref(false);
const rechargeOrderNum = ref("");
const rechargeStatus = ref("");
const wechatQrUrl = ref("");
const pendingCredit = ref(0);
const applyingCredit = ref(false);
let pollTimer: number | undefined;

const totalPrice = computed(() =>
  pendingOrders.value
    .reduce((s, i) => s + unitPriceValue(i) * i.num, 0)
    .toFixed(2),
);

const balanceText = computed(() =>
  balanceLoaded.value
    ? t("payment.currentBalance", { balance: `¥${balance.value}` })
    : t("payment.balanceAfterPassword"),
);

const rechargeStatusText = computed(
  () =>
    ({
      pending: t("payment.pending"),
      paid: t("payment.paidPendingCredit"),
      credited: t("payment.credited"),
      failed: t("payment.failed"),
    })[rechargeStatus.value] ?? t("payment.pendingStart"),
);

const rechargeStatusTag = computed(
  () =>
    ({
      pending: "warning",
      paid: "success",
      credited: "success",
      failed: "danger",
    })[rechargeStatus.value] ?? "info",
);

const qrImageUrl = computed(
  () =>
    `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(
      wechatQrUrl.value,
    )}`,
);

const canPay = computed(
  () =>
    payMethod.value === "balance" &&
    payPassword.value.length === 6 &&
    pendingOrders.value.length > 0,
);

function productName(item: any) {
  return item.name || item.product_name || t("common.unknown");
}

function unitPriceValue(item: any) {
  return Number(item.money ?? item.discount_price ?? item.price ?? 0);
}

function syncPaidCartCount(paidCount: number) {
  if (paidCount <= 0) return;
  userStore.setCartCount(Math.max(0, userStore.cartCount - paidCount));
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

function stopPolling() {
  if (pollTimer) {
    window.clearInterval(pollTimer);
    pollTimer = undefined;
  }
}

function startPolling() {
  stopPolling();
  pollTimer = window.setInterval(pollRechargeStatus, 2000);
}

async function startRecharge() {
  const amount = rechargeAmount.value || Number(totalPrice.value);
  if (amount <= 0) return ElMessage.warning(t("payment.invalidRechargeAmount"));
  rechargeLoading.value = true;
  try {
    const api = payMethod.value === "wechat" ? wechatRecharge : alipayRecharge;
    const res: any = await api({ amount });
    rechargeOrderNum.value = res.data?.order_num ?? "";
    rechargeStatus.value = "pending";
    pendingCredit.value = 0;
    wechatQrUrl.value = res.data?.qr_code_url ?? "";
    if (res.data?.pay_url) {
      window.open(res.data.pay_url, "_blank");
    }
    if (rechargeOrderNum.value) startPolling();
  } finally {
    rechargeLoading.value = false;
  }
}

async function pollRechargeStatus() {
  if (!rechargeOrderNum.value) return;
  const res: any = await getRechargeStatus({
    order_num: rechargeOrderNum.value,
  });
  rechargeStatus.value = res.data?.status ?? "";
  pendingCredit.value = Number(res.data?.pending_credit ?? pendingCredit.value);
  if (["paid", "credited", "failed"].includes(rechargeStatus.value)) {
    stopPolling();
  }
  if (rechargeStatus.value === "credited") {
    if (payPassword.value.length === 6) {
      await loadBalance();
    }
    payMethod.value = "balance";
  }
}

async function handleApplyCredit() {
  if (payPassword.value.length !== 6) return ElMessage.warning(t("payment.payPasswordRequired"));
  applyingCredit.value = true;
  try {
    await applyPendingCredit({ key: payPassword.value });
    ElMessage.success(t("payment.creditedSuccess"));
    pendingCredit.value = 0;
    rechargeStatus.value = "credited";
    payMethod.value = "balance";
    await loadBalance();
  } finally {
    applyingCredit.value = false;
  }
}

async function handlePay() {
  if (payPassword.value.length !== 6) return ElMessage.warning(t("payment.payPasswordRequired"));
  if (!pendingOrders.value.length)
    return ElMessage.error(t("payment.orderMissing"));

  paying.value = true;
  const paidOrderIds: number[] = [];
  try {
    // Pay each order individually because each cart item maps to one order.
    for (const item of pendingOrders.value) {
      await payOrder({
        order_id: item.order_id,
        money: unitPriceValue(item),
        product_id: item.product_id,
        boss_id: item.boss_id,
        num: item.num,
        key: payPassword.value,
      });
      paidOrderIds.push(item.order_id);
    }
    sessionStorage.removeItem(activeUserSessionStorageKey("pending_orders"));
    sessionStorage.setItem(
      activeUserSessionStorageKey("paid_order_ids"),
      JSON.stringify(paidOrderIds),
    );
    syncPaidCartCount(paidOrderIds.length);
    ElMessage.success(t("payment.paymentSuccess"));
    router.push("/order/success");
  } catch (error: any) {
    if (paidOrderIds.length) {
      syncPaidCartCount(paidOrderIds.length);
      pendingOrders.value = pendingOrders.value.filter(
        (item) => !paidOrderIds.includes(item.order_id),
      );
      if (pendingOrders.value.length) {
        sessionStorage.setItem(
          activeUserSessionStorageKey("pending_orders"),
          JSON.stringify(pendingOrders.value),
        );
      } else {
        sessionStorage.removeItem(activeUserSessionStorageKey("pending_orders"));
      }
      ElMessage.error(t("payment.partialPaymentFailed"));
      return;
    }
    ElMessage.error(error?.message || t("payment.paymentFailed"));
  } finally {
    paying.value = false;
  }
}

watch(
  totalPrice,
  (value) => {
    rechargeAmount.value = Number(value || 0);
  },
  { immediate: true },
);

watch(payMethod, () => {
  stopPolling();
  rechargeOrderNum.value = "";
  rechargeStatus.value = "";
  wechatQrUrl.value = "";
  pendingCredit.value = 0;
});

onUnmounted(stopPolling);
</script>

<style scoped>
.payment-wrap {
  max-width: 600px;
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
.order-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 14px;
}
.price {
  color: #f56c6c;
}
.total-row {
  display: flex;
  justify-content: space-between;
  font-size: 15px;
  font-weight: 600;
}
.total-price {
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
}
.pay-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
.payment-key-row {
  display: flex;
  max-width: 420px;
  gap: 10px;
}
.payment-key-row .el-input {
  flex: 1;
}
.recharge-panel {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 14px;
}
.recharge-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.recharge-result {
  border-top: 1px solid #eee;
  padding-top: 14px;
}
.result-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}
.qr-wrap {
  text-align: center;
  margin: 14px 0;
}
.qr-wrap img {
  width: 220px;
  height: 220px;
}
.qr-link {
  word-break: break-all;
  color: #999;
  font-size: 12px;
  margin-top: 8px;
}
.apply-box {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}
.apply-box .el-input {
  flex: 1;
}
</style>
