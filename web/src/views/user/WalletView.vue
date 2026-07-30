<template>
  <el-card>
    <template #header>{{ t("wallet.title") }}</template>
    <div class="balance-section">
      <div class="balance-label">{{ t("wallet.balance") }}</div>
      <div class="balance-amount">{{ balanceText }}</div>
      <div v-if="!payKeySet" class="set-key-box">
        <el-alert
          :title="t('wallet.setPayPasswordTip')"
          type="warning"
          :closable="false"
          show-icon
        />
        <div class="set-key-form">
          <el-input
            v-model="setKeyForm.key"
            type="password"
            maxlength="6"
            show-password
            :placeholder="t('wallet.payPasswordPlaceholder')"
          />
          <el-input
            v-model="setKeyForm.key_confirm"
            type="password"
            maxlength="6"
            show-password
            :placeholder="t('wallet.payPasswordConfirmPlaceholder')"
          />
          <el-button type="primary" :loading="settingPayKey" @click="handleSetPayKey">
            {{ t("wallet.setPayPassword") }}
          </el-button>
        </div>
      </div>
      <div v-else class="balance-form">
        <el-input
          v-model="payKey"
          type="password"
          maxlength="6"
          show-password
          :placeholder="t('wallet.payPasswordPlaceholder')"
        />
        <el-button
          type="primary"
          :loading="loading"
          :disabled="payKey.length !== 6"
          @click="loadBalance"
        >
          {{ t("wallet.viewBalance") }}
        </el-button>
      </div>
    </div>
    <el-divider />
    <div class="wallet-actions">
      <el-button type="primary" @click="openRechargeDialog">{{ t("wallet.recharge") }}</el-button>
      <el-button :loading="pendingLoading" @click="loadPendingCredit"
        >{{ t("wallet.refreshPending") }}</el-button
      >
      <span v-if="pendingCredit > 0" class="pending-text">
        {{ t("wallet.pendingCredit", { amount: pendingCredit.toFixed(2) }) }}
      </span>
    </div>

    <el-dialog v-model="rechargeVisible" :title="t('wallet.rechargeTitle')" width="520px">
      <el-form label-width="86px">
        <el-form-item :label="t('wallet.rechargeAmount')">
          <el-input-number
            v-model="rechargeAmount"
            :min="0.01"
            :precision="2"
            :step="10"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item :label="t('wallet.rechargeChannel')">
          <el-radio-group v-model="rechargeChannel">
            <el-radio value="wechat">{{ t("wallet.wechat") }}</el-radio>
            <el-radio value="alipay">{{ t("wallet.alipay") }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <div v-if="rechargeOrderNum" class="recharge-result">
        <div class="result-row">
          <span>{{ t("wallet.rechargeOrderNo") }}</span>
          <b>{{ rechargeOrderNum }}</b>
        </div>
        <div class="result-row">
          <span>{{ t("wallet.payStatus") }}</span>
          <el-tag :type="rechargeStatusTag">{{ rechargeStatusText }}</el-tag>
        </div>
        <div v-if="wechatQrUrl" class="qr-wrap">
          <img :src="qrImageUrl" :alt="t('wallet.qrAlt')" />
          <div class="qr-link">{{ wechatQrUrl }}</div>
        </div>
        <div v-if="pendingCredit > 0" class="apply-box">
          <el-input
            v-model="payKey"
            type="password"
            maxlength="6"
            show-password
            :placeholder="t('wallet.confirmCreditPlaceholder')"
          />
          <el-button
            type="success"
            :loading="applyingCredit"
            :disabled="payKey.length !== 6"
            @click="handleApplyCredit"
          >
            {{ t("wallet.confirmCredit") }}
          </el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="rechargeVisible = false">{{ t("wallet.close") }}</el-button>
        <el-button
          type="primary"
          :loading="rechargeLoading"
          @click="startRecharge"
        >
          {{ t("wallet.startRecharge") }}
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import { getMoney, setPayKey } from "@/api/flashSale";
import { getUserInfo } from "@/api/user";
import {
  alipayRecharge,
  applyPendingCredit,
  getPendingCredit,
  getRechargeStatus,
  wechatRecharge,
} from "@/api/recharge";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const { t } = useI18n();
const balance = ref("0.00");
const loaded = ref(false);
const loading = ref(false);
const payKey = ref("");
const settingPayKey = ref(false);
const setKeyForm = ref({ key: "", key_confirm: "" });
const rechargeVisible = ref(false);
const rechargeAmount = ref(100);
const rechargeChannel = ref<"wechat" | "alipay">("wechat");
const rechargeLoading = ref(false);
const rechargeOrderNum = ref("");
const rechargeStatus = ref("");
const wechatQrUrl = ref("");
const pendingCredit = ref(0);
const pendingLoading = ref(false);
const applyingCredit = ref(false);
let pollTimer: number | undefined;

const payKeySet = computed(() => userStore.userInfo?.pay_key_set ?? false);

const balanceText = computed(() =>
  !payKeySet.value
    ? t("wallet.setPayPasswordFirst")
    : loaded.value
      ? `¥${balance.value}`
      : t("wallet.viewAfterPassword"),
);

const rechargeStatusText = computed(
  () =>
    ({
      pending: t("wallet.statusPending"),
      paid: t("wallet.statusPaid"),
      credited: t("wallet.statusCredited"),
      failed: t("wallet.statusFailed"),
    })[rechargeStatus.value] ?? t("wallet.statusInitial"),
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

async function loadBalance() {
  if (!payKeySet.value) {
    return ElMessage.warning(t("wallet.setPayPasswordFirst"));
  }
  if (payKey.value.length !== 6) {
    return ElMessage.warning(t("wallet.payPasswordPlaceholder"));
  }
  loading.value = true;
  try {
    const res: any = await getMoney({ key: payKey.value });
    balance.value = Number(res.data?.user_money ?? 0).toFixed(2);
    loaded.value = true;
  } catch {
    loaded.value = false;
  } finally {
    loading.value = false;
  }
}

async function refreshUserInfo() {
  const res: any = await getUserInfo();
  if (res.data) {
    userStore.setUserInfo(res.data);
  }
}

async function handleSetPayKey() {
  if (setKeyForm.value.key.length !== 6) {
    return ElMessage.warning(t("wallet.payPasswordPlaceholder"));
  }
  if (setKeyForm.value.key !== setKeyForm.value.key_confirm) {
    return ElMessage.warning(t("wallet.passwordMismatch"));
  }
  settingPayKey.value = true;
  try {
    await setPayKey(setKeyForm.value);
    ElMessage.success(t("wallet.passwordSetSuccess"));
    payKey.value = setKeyForm.value.key;
    setKeyForm.value = { key: "", key_confirm: "" };
    await refreshUserInfo();
    await loadBalance();
  } finally {
    settingPayKey.value = false;
  }
}

function openRechargeDialog() {
  rechargeVisible.value = true;
  loadPendingCredit();
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
  if (!rechargeAmount.value || rechargeAmount.value <= 0) {
    return ElMessage.warning(t("wallet.amountRequired"));
  }
  rechargeLoading.value = true;
  try {
    const api =
      rechargeChannel.value === "wechat" ? wechatRecharge : alipayRecharge;
    const res: any = await api({ amount: Number(rechargeAmount.value) });
    rechargeOrderNum.value = res.data?.order_num ?? "";
    rechargeStatus.value = "pending";
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
    ElMessage.success(t("wallet.rechargeCredited"));
    if (payKey.value.length === 6) await loadBalance();
  }
}

async function loadPendingCredit() {
  pendingLoading.value = true;
  try {
    const res: any = await getPendingCredit();
    pendingCredit.value = Number(res.data?.pending ?? 0);
  } finally {
    pendingLoading.value = false;
  }
}

async function handleApplyCredit() {
  if (!payKeySet.value) return ElMessage.warning(t("wallet.setPayPasswordFirst"));
  if (payKey.value.length !== 6) return ElMessage.warning(t("wallet.payPasswordPlaceholder"));
  applyingCredit.value = true;
  try {
    await applyPendingCredit({ key: payKey.value });
    ElMessage.success(t("wallet.creditSuccess"));
    pendingCredit.value = 0;
    rechargeStatus.value = "credited";
    await loadBalance();
  } finally {
    applyingCredit.value = false;
  }
}

watch(rechargeVisible, (visible) => {
  if (!visible) stopPolling();
});

onMounted(async () => {
  await refreshUserInfo();
  await loadPendingCredit();
});
onUnmounted(stopPolling);
</script>

<style scoped>
.balance-section {
  padding: 20px 0;
}
.balance-label {
  color: #999;
  font-size: 14px;
}
.balance-amount {
  font-size: 36px;
  font-weight: bold;
  color: #409eff;
  margin-top: 4px;
}
.balance-form {
  display: flex;
  gap: 10px;
  max-width: 420px;
  margin-top: 18px;
}
.balance-form .el-input {
  flex: 1;
}
.set-key-box {
  max-width: 520px;
  margin-top: 18px;
}
.set-key-form {
  display: grid;
  grid-template-columns: 1fr 1fr auto;
  gap: 10px;
  margin-top: 12px;
}
.wallet-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
.pending-text {
  color: #f56c6c;
  font-weight: 600;
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
