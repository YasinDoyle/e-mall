<template>
  <div class="flash-wrap">
    <div class="flash-header">
      <el-icon color="#f56c6c" size="28"><Lightning /></el-icon>
      <span class="flash-title">秒杀专场</span>
      <el-tag type="danger" effect="dark" style="margin-left: 12px"
        >火热进行中</el-tag
      >
    </div>

    <el-empty v-if="!list.length && !loading" description="暂无秒杀商品">
      <el-button type="primary" @click="$router.push('/products')"
        >逛逛普通商品</el-button
      >
    </el-empty>

    <div v-else class="flash-grid">
      <el-card
        v-for="item in list"
        :key="item.id"
        shadow="hover"
        class="flash-card"
        @click="$router.push(`/flash-sale/${item.id}`)"
      >
        <div class="flash-badge">秒杀</div>
        <div class="flash-name">
          {{ item.title || `商品 #${item.product_id}` }}
        </div>
        <div class="flash-price">
          <span class="price">¥{{ item.money }}</span>
        </div>
        <div class="flash-stock">
          <el-progress
            :percentage="stockPercent(item)"
            :stroke-width="8"
            status="exception"
            :show-text="false"
          />
          <span class="stock-text">剩余 {{ item.num }} 件</span>
        </div>
        <el-button
          type="danger"
          size="small"
          style="width: 100%; margin-top: 10px"
        >
          立即抢购
        </el-button>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { Lightning } from "@element-plus/icons-vue";
import { getFlashSaleList } from "@/api/flashSale";

const list = ref<any[]>([]);
const loading = ref(true);

// 简单估算已售比例（没有原始库存字段时固定显示剩余热度）
function stockPercent(item: any) {
  const remain = item.num ?? 0;
  if (remain <= 0) return 100;
  if (remain <= 5) return 80;
  if (remain <= 20) return 50;
  return 20;
}

async function loadList() {
  try {
    const res: any = await getFlashSaleList();
    list.value = res.data?.item ?? [];
  } finally {
    loading.value = false;
  }
}

onMounted(loadList);
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
