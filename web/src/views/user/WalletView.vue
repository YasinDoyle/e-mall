<template>
  <el-card>
    <template #header>我的钱包</template>
    <div class="balance-section">
      <div class="balance-label">账户余额</div>
      <div class="balance-amount">¥{{ balance }}</div>
    </div>
    <el-divider />
    <div style="color: #999; font-size: 14px">
      充值功能（微信 / 支付宝）开发中，敬请期待
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { getMoney } from "@/api/flashSale";

const balance = ref("0.00");

async function loadBalance() {
  try {
    const res: any = await getMoney();
    balance.value = res.data ?? "0.00";
  } catch {}
}

onMounted(loadBalance);
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
</style>
