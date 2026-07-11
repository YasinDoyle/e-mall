<template>
  <el-card>
    <template #header>我的钱包</template>
    <div class="balance-section">
      <div class="balance-label">账户余额</div>
      <div class="balance-amount">{{ balanceText }}</div>
      <div class="balance-form">
        <el-input
          v-model="payKey"
          type="password"
          maxlength="6"
          show-password
          placeholder="请输入6位支付密码"
        />
        <el-button
          type="primary"
          :loading="loading"
          :disabled="payKey.length !== 6"
          @click="loadBalance"
        >
          查看余额
        </el-button>
      </div>
    </div>
    <el-divider />
    <div class="wallet-actions">
      <el-button type="primary" @click="openRechargeDialog">充值</el-button>
      <el-button :loading="pendingLoading" @click="loadPendingCredit"
        >刷新待入账</el-button
      >
      <span v-if="pendingCredit > 0" class="pending-text">
        待入账：¥{{ pendingCredit.toFixed(2) }}
      </span>
    </div>

    <el-dialog v-model="rechargeVisible" title="钱包充值" width="520px">
      <el-form label-width="86px">
        <el-form-item label="充值金额">
          <el-input-number
            v-model="rechargeAmount"
            :min="0.01"
            :precision="2"
            :step="10"
            style="width: 220px"
          />
        </el-form-item>
        <el-form-item label="充值方式">
          <el-radio-group v-model="rechargeChannel">
            <el-radio value="wechat">微信</el-radio>
            <el-radio value="alipay">支付宝</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <div v-if="rechargeOrderNum" class="recharge-result">
        <div class="result-row">
          <span>充值单号</span>
          <b>{{ rechargeOrderNum }}</b>
        </div>
        <div class="result-row">
          <span>支付状态</span>
          <el-tag :type="rechargeStatusTag">{{ rechargeStatusText }}</el-tag>
        </div>
        <div v-if="wechatQrUrl" class="qr-wrap">
          <img :src="qrImageUrl" alt="微信支付二维码" />
          <div class="qr-link">{{ wechatQrUrl }}</div>
        </div>
        <div v-if="pendingCredit > 0" class="apply-box">
          <el-input
            v-model="payKey"
            type="password"
            maxlength="6"
            show-password
            placeholder="输入6位支付密码确认入账"
          />
          <el-button
            type="success"
            :loading="applyingCredit"
            :disabled="payKey.length !== 6"
            @click="handleApplyCredit"
          >
            确认入账
          </el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="rechargeVisible = false">关闭</el-button>
        <el-button
          type="primary"
          :loading="rechargeLoading"
          @click="startRecharge"
        >
          发起充值
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { getMoney } from "@/api/flashSale";
import {
  alipayRecharge,
  applyPendingCredit,
  getPendingCredit,
  getRechargeStatus,
  wechatRecharge,
} from "@/api/recharge";

const balance = ref("0.00");
const loaded = ref(false);
const loading = ref(false);
const payKey = ref("");
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

const balanceText = computed(() =>
  loaded.value ? `¥${balance.value}` : "输入支付密码后查看",
);

const rechargeStatusText = computed(
  () =>
    ({
      pending: "待支付",
      paid: "已支付，待入账",
      credited: "已入账",
      failed: "已失败",
    })[rechargeStatus.value] ?? "待发起",
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
  if (payKey.value.length !== 6) {
    return ElMessage.warning("请输入6位支付密码");
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
    return ElMessage.warning("请输入充值金额");
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
    ElMessage.success("充值已入账");
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
  if (payKey.value.length !== 6) return ElMessage.warning("请输入6位支付密码");
  applyingCredit.value = true;
  try {
    await applyPendingCredit({ key: payKey.value });
    ElMessage.success("入账成功");
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

onMounted(loadPendingCredit);
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
