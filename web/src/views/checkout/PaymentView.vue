<template>
  <div class="payment-wrap">
    <h2 class="page-title">支付订单</h2>

    <el-card class="section-card">
      <template #header>订单信息</template>
      <el-empty v-if="!pendingOrders.length" description="暂无待支付订单" />
      <div v-for="item in pendingOrders" :key="item.order_id" class="order-row">
        <span>{{ productName(item) }} × {{ item.num }}</span>
        <span class="price"
          >¥{{ (unitPriceValue(item) * item.num).toFixed(2) }}</span
        >
      </div>
      <el-divider v-if="pendingOrders.length" style="margin: 12px 0" />
      <div v-if="pendingOrders.length" class="total-row">
        <span>应付总额</span>
        <span class="total-price">¥{{ totalPrice }}</span>
      </div>
    </el-card>

    <el-card class="section-card">
      <template #header>支付方式</template>
      <el-radio-group
        v-model="payMethod"
        style="display: flex; flex-direction: column; gap: 12px"
      >
        <el-radio value="balance">
          <span>余额支付</span>
          <span style="color: #999; font-size: 12px; margin-left: 8px"
            >（当前余额：{{ balanceText }}）</span
          >
        </el-radio>
        <el-radio value="wechat">微信充值后余额支付</el-radio>
        <el-radio value="alipay">支付宝充值后余额支付</el-radio>
      </el-radio-group>
    </el-card>

    <el-card v-if="payMethod === 'balance'" class="section-card">
      <template #header>支付密码</template>
      <div class="payment-key-row">
        <el-input
          v-model="payPassword"
          type="password"
          placeholder="请输入6位支付密码"
          maxlength="6"
          show-password
        />
        <el-button
          :loading="balanceLoading"
          :disabled="payPassword.length !== 6"
          @click="loadBalance"
        >
          刷新余额
        </el-button>
      </div>
    </el-card>

    <el-card v-else class="section-card">
      <template #header>
        {{ payMethod === "wechat" ? "微信充值" : "支付宝充值" }}
      </template>
      <div class="recharge-panel">
        <div class="recharge-row">
          <span>建议充值金额</span>
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
          {{ rechargeOrderNum ? "重新发起充值" : "发起充值" }}
        </el-button>
      </div>

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
            v-model="payPassword"
            type="password"
            maxlength="6"
            show-password
            placeholder="输入6位支付密码确认入账"
          />
          <el-button
            type="success"
            :loading="applyingCredit"
            :disabled="payPassword.length !== 6"
            @click="handleApplyCredit"
          >
            确认入账
          </el-button>
        </div>
      </div>
    </el-card>

    <div class="pay-footer">
      <el-button size="large" @click="$router.push('/cart')">取消</el-button>
      <el-button
        type="primary"
        size="large"
        :loading="paying"
        :disabled="!canPay"
        @click="handlePay"
      >
        立即支付 ¥{{ totalPrice }}
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

const router = useRouter();

const pendingOrders = ref<any[]>(
  JSON.parse(sessionStorage.getItem("pending_orders") ?? "[]"),
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
  balanceLoaded.value ? `¥${balance.value}` : "输入支付密码后刷新",
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

const canPay = computed(
  () =>
    payMethod.value === "balance" &&
    payPassword.value.length === 6 &&
    pendingOrders.value.length > 0,
);

function productName(item: any) {
  return item.name || item.product_name || "商品";
}

function unitPriceValue(item: any) {
  return Number(item.money ?? item.discount_price ?? item.price ?? 0);
}

async function loadBalance() {
  if (payPassword.value.length !== 6) {
    return ElMessage.warning("请输入6位支付密码");
  }
  balanceLoading.value = true;
  try {
    const res: any = await getMoney({ key: payPassword.value });
    balance.value = Number(res.data?.user_money ?? 0).toFixed(2);
    balanceLoaded.value = true;
  } catch (error: any) {
    balanceLoaded.value = false;
    ElMessage.error(error?.message || "余额查询失败，请检查支付密码");
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
  if (amount <= 0) return ElMessage.warning("充值金额不合法");
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
  if (payPassword.value.length !== 6) return ElMessage.warning("请输入6位支付密码");
  applyingCredit.value = true;
  try {
    await applyPendingCredit({ key: payPassword.value });
    ElMessage.success("入账成功，请继续完成余额支付");
    pendingCredit.value = 0;
    rechargeStatus.value = "credited";
    payMethod.value = "balance";
    await loadBalance();
  } finally {
    applyingCredit.value = false;
  }
}

async function handlePay() {
  if (payPassword.value.length !== 6) return ElMessage.warning("请输入6位支付密码");
  if (!pendingOrders.value.length)
    return ElMessage.error("订单信息丢失，请重新下单");

  paying.value = true;
  const paidOrderIds: number[] = [];
  try {
    // 逐笔支付（每个商品一个订单）
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
    sessionStorage.removeItem("pending_orders");
    sessionStorage.setItem("paid_order_ids", JSON.stringify(paidOrderIds));
    ElMessage.success("支付成功！");
    router.push("/order/success");
  } catch (error: any) {
    if (paidOrderIds.length) {
      pendingOrders.value = pendingOrders.value.filter(
        (item) => !paidOrderIds.includes(item.order_id),
      );
      if (pendingOrders.value.length) {
        sessionStorage.setItem(
          "pending_orders",
          JSON.stringify(pendingOrders.value),
        );
      } else {
        sessionStorage.removeItem("pending_orders");
      }
      ElMessage.error("部分订单支付成功，剩余订单支付失败，请重试");
      return;
    }
    ElMessage.error(error?.message || "支付失败，请检查余额或支付密码");
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
