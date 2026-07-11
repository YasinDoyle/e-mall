<template>
  <div class="payment-wrap">
    <h2 class="page-title">支付订单</h2>

    <el-card class="section-card">
      <template #header>订单信息</template>
      <div v-for="item in pendingOrders" :key="item.order_id" class="order-row">
        <span>{{ item.product_name }} × {{ item.num }}</span>
        <span class="price"
          >¥{{ (Number(item.price) * item.num).toFixed(2) }}</span
        >
      </div>
      <el-divider style="margin: 12px 0" />
      <div class="total-row">
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
            >（当前余额：¥{{ balance }}）</span
          >
        </el-radio>
        <el-radio value="wechat" disabled>微信支付（即将开放）</el-radio>
        <el-radio value="alipay" disabled>支付宝（即将开放）</el-radio>
      </el-radio-group>
    </el-card>

    <el-card v-if="payMethod === 'balance'" class="section-card">
      <template #header>支付密码</template>
      <el-input
        v-model="payPassword"
        type="password"
        placeholder="请输入6位支付密码"
        maxlength="6"
        show-password
        style="max-width: 300px"
      />
    </el-card>

    <div class="pay-footer">
      <el-button size="large" @click="$router.push('/cart')">取消</el-button>
      <el-button
        type="primary"
        size="large"
        :loading="paying"
        :disabled="!payPassword"
        @click="handlePay"
      >
        立即支付 ¥{{ totalPrice }}
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { payOrder } from "@/api/order";
import { getMoney } from "@/api/flashSale";

const router = useRouter();

const pendingOrders = ref<any[]>(
  JSON.parse(sessionStorage.getItem("pending_orders") ?? "[]"),
);
const payMethod = ref("balance");
const payPassword = ref("");
const paying = ref(false);
const balance = ref("0.00");

const totalPrice = computed(() =>
  pendingOrders.value
    .reduce((s, i) => s + Number(i.price) * i.num, 0)
    .toFixed(2),
);

async function loadBalance() {
  try {
    const res: any = await getMoney();
    balance.value = res.data ?? "0.00";
  } catch {}
}

async function handlePay() {
  if (!payPassword.value) return ElMessage.warning("请输入支付密码");
  if (!pendingOrders.value.length)
    return ElMessage.error("订单信息丢失，请重新下单");

  paying.value = true;
  try {
    // 逐笔支付（每个商品一个订单）
    for (const item of pendingOrders.value) {
      await payOrder({
        order_id: item.order_id,
        money: Number(item.price) * item.num,
        boss_id: item.boss_id,
        key: payPassword.value,
      });
    }
    sessionStorage.removeItem("pending_orders");
    ElMessage.success("支付成功！");
    router.push("/order/success");
  } catch {
    ElMessage.error("支付失败，请检查余额或支付密码");
  } finally {
    paying.value = false;
  }
}

onMounted(loadBalance);
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
</style>
