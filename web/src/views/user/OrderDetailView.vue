<template>
  <el-card v-if="order">
    <template #header>
      <div
        style="
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <span>订单详情</span>
        <el-tag :type="statusTagType(order.type)">{{
          statusText(order.type)
        }}</el-tag>
      </div>
    </template>

    <!-- 物流时间轴 -->
    <el-timeline style="margin-bottom: 20px">
      <el-timeline-item
        v-for="step in timeline"
        :key="step.label"
        :type="step.done ? 'primary' : 'info'"
        :hollow="!step.done"
      >
        {{ step.label }}
      </el-timeline-item>
    </el-timeline>

    <el-descriptions :column="2" border>
      <el-descriptions-item label="订单号">{{
        order.order_num
      }}</el-descriptions-item>
      <el-descriptions-item label="商品ID">{{
        order.product_id
      }}</el-descriptions-item>
      <el-descriptions-item label="数量"
        >{{ order.num }} 件</el-descriptions-item
      >
      <el-descriptions-item label="金额">
        <span style="color: #f56c6c; font-weight: bold">
          ¥{{ (order.money * order.num).toFixed(2) }}
        </span>
      </el-descriptions-item>
    </el-descriptions>

    <div class="order-actions">
      <el-button v-if="order.type === 3" type="primary" @click="handleReceive"
        >确认收货</el-button
      >
      <el-button
        v-if="order.type === 4"
        @click="$router.push(`/product/${order.product_id}`)"
        >写评价</el-button
      >
      <el-button @click="$router.back()">返回</el-button>
    </div>
  </el-card>

  <div v-else style="text-align: center; padding: 60px">
    <el-icon class="is-loading" size="40"><Loading /></el-icon>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import { Loading } from "@element-plus/icons-vue";
import { getOrderDetail, receiveOrder } from "@/api/order";

const route = useRoute();
const order = ref<any>(null);

const statusText = (type: number) =>
  ({ 1: "待支付", 2: "待发货", 3: "已发货", 4: "已完成" })[type] ?? "未知";
const statusTagType = (type: number) =>
  (({ 1: "warning", 2: "primary", 3: "warning", 4: "success" })[type] ??
    "info") as any;

const timeline = computed(() => {
  const type = order.value?.type ?? 0;
  return [
    { label: "提交订单", done: type >= 1 },
    { label: "支付成功", done: type >= 2 },
    { label: "商家发货", done: type >= 3 },
    { label: "确认收货", done: type >= 4 },
  ];
});

async function loadOrder() {
  try {
    const res: any = await getOrderDetail({ id: Number(route.params.id) });
    order.value = res.data;
  } catch {}
}

async function handleReceive() {
  await receiveOrder({ id: order.value.id });
  ElMessage.success("已确认收货");
  loadOrder();
}

onMounted(loadOrder);
</script>

<style scoped>
.order-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  justify-content: flex-end;
}
</style>
