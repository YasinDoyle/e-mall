<template>
  <div class="flash-wrap">
    <div class="flash-header">
      <el-icon color="#f56c6c" size="28"><Lightning /></el-icon>
      <span class="flash-title">秒杀专场</span>
      <el-tag type="danger" effect="dark" style="margin-left: 12px"
        >火热进行中</el-tag
      >
      <span class="countdown">距本场结束 {{ countdownText }}</span>
    </div>

    <el-skeleton v-if="loading" :rows="5" animated />

    <el-empty v-else-if="!list.length" description="暂无秒杀商品">
      <el-button type="primary" @click="$router.push('/products')"
        >逛逛普通商品</el-button
      >
    </el-empty>

    <div v-else class="flash-grid">
      <el-card
        v-for="item in list"
        :key="flashSaleId(item)"
        shadow="hover"
        class="flash-card"
        @click="$router.push(`/flash-sale/${productId(item)}`)"
      >
        <div class="flash-badge">秒杀</div>
        <div class="flash-name">
          {{ titleText(item) }}
        </div>
        <div class="flash-price">
          <span class="price">¥{{ moneyText(item) }}</span>
        </div>
        <div class="flash-stock">
          <el-progress
            :percentage="stockPercent(item)"
            :stroke-width="8"
            status="exception"
            :show-text="false"
          />
          <span class="stock-text">剩余 {{ stockNum(item) }} 件</span>
        </div>
        <el-button
          type="danger"
          size="small"
          style="width: 100%; margin-top: 10px"
          :disabled="stockNum(item) <= 0"
        >
          {{ stockNum(item) > 0 ? "立即抢购" : "已售罄" }}
        </el-button>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from "vue";
import { Lightning } from "@element-plus/icons-vue";
import { getFlashSaleList } from "@/api/flashSale";

const list = ref<any[]>([]);
const loading = ref(true);
const now = ref(Date.now());
let timer: number | undefined;

const countdownText = computed(() => {
  const end = new Date(now.value);
  end.setHours(23, 59, 59, 999);
  const diff = Math.max(0, end.getTime() - now.value);
  const hours = Math.floor(diff / 1000 / 60 / 60);
  const minutes = Math.floor((diff / 1000 / 60) % 60);
  const seconds = Math.floor((diff / 1000) % 60);
  return `${padTime(hours)}:${padTime(minutes)}:${padTime(seconds)}`;
});

function padTime(value: number) {
  return String(value).padStart(2, "0");
}

function flashSaleId(item: any) {
  return item.id ?? item.Id;
}

function productId(item: any) {
  return item.product_id ?? item.ProductId;
}

function stockNum(item: any) {
  return Number(item.num ?? item.Num ?? 0);
}

function titleText(item: any) {
  return item.title ?? item.Title ?? `商品 #${productId(item)}`;
}

function moneyText(item: any) {
  return Number(item.money ?? item.Money ?? 0).toFixed(2);
}

// 简单估算已售比例（没有原始库存字段时固定显示剩余热度）
function stockPercent(item: any) {
  const remain = stockNum(item);
  if (remain <= 0) return 100;
  if (remain <= 5) return 80;
  if (remain <= 20) return 50;
  return 20;
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getFlashSaleList();
    list.value = Array.isArray(res.data) ? res.data : (res.data?.item ?? []);
  } catch {
    list.value = [];
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadList();
  timer = window.setInterval(() => {
    now.value = Date.now();
  }, 1000);
});

onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});
</script>

<style scoped>
.flash-wrap {
  max-width: 1100px;
  margin: 0 auto;
}
.flash-header {
  display: flex;
  align-items: center;
  margin-bottom: 24px;
  gap: 4px;
}
.flash-title {
  font-size: 24px;
  font-weight: bold;
  color: #f56c6c;
}
.countdown {
  margin-left: auto;
  color: #666;
  font-size: 14px;
}
.flash-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}
.flash-card {
  cursor: pointer;
  position: relative;
  transition: transform 0.2s;
}
.flash-card:hover {
  transform: translateY(-4px);
}
.flash-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  background: #f56c6c;
  color: #fff;
  font-size: 11px;
  padding: 2px 6px;
  border-radius: 10px;
}
.flash-name {
  font-size: 15px;
  font-weight: 500;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.flash-price {
  margin-bottom: 10px;
}
.price {
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
}
.flash-stock {
  margin-top: 8px;
}
.stock-text {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  display: block;
}
</style>
