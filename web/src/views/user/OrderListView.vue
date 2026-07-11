<template>
  <el-card>
    <template #header>我的订单</template>

    <el-tabs v-model="activeTab" @tab-change="loadOrders">
      <el-tab-pane label="全部" name="0" />
      <el-tab-pane label="待支付" name="1" />
      <el-tab-pane label="待发货" name="2" />
      <el-tab-pane label="已发货" name="3" />
      <el-tab-pane label="已完成" name="4" />
    </el-tabs>

    <el-empty v-if="!orders.length" description="暂无订单" />

    <div v-for="order in orders" :key="order.id" class="order-card">
      <div class="order-header">
        <span class="order-num">订单号：{{ order.order_num }}</span>
        <el-tag :type="statusTagType(order.type)">{{
          statusText(order.type)
        }}</el-tag>
      </div>
      <div class="order-body" @click="$router.push(`/user/orders/${order.id}`)">
        <div class="order-money">
          ¥{{ (order.money * order.num).toFixed(2) }}
        </div>
        <div style="color: #999; font-size: 12px">{{ order.num }} 件商品</div>
      </div>
      <div class="order-actions">
        <el-button
          v-if="order.type === 3"
          size="small"
          type="primary"
          @click="handleReceive(order.id)"
          >确认收货</el-button
        >
        <el-button size="small" type="danger" @click="handleDelete(order.id)"
          >删除</el-button
        >
        <el-button
          size="small"
          @click="$router.push(`/user/orders/${order.id}`)"
        >
          查看详情
        </el-button>
      </div>
    </div>

    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="prev, pager, next"
      style="margin-top: 16px; justify-content: center"
      @current-change="loadOrders"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getOrderList, deleteOrder, receiveOrder } from "@/api/order";

const activeTab = ref("0");
const orders = ref<any[]>([]);
const page = ref(1);
const pageSize = 10;
const total = ref(0);

const statusText = (type: number) =>
  ({ 1: "待支付", 2: "待发货", 3: "已发货", 4: "已完成" })[type] ?? "未知";
const statusTagType = (type: number) =>
  (({ 1: "warning", 2: "primary", 3: "warning", 4: "success" })[type] ??
    "info") as any;

async function loadOrders() {
  try {
    const params: any = { page_num: page.value, page_size: pageSize };
    if (activeTab.value !== "0") params.type = Number(activeTab.value);
    const res: any = await getOrderList(params);
    orders.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } catch {}
}

async function handleReceive(id: number) {
  await ElMessageBox.confirm("确认已收到货物？", "提示", { type: "warning" });
  await receiveOrder({ id });
  ElMessage.success("已确认收货");
  loadOrders();
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该订单？", "提示", { type: "warning" });
  await deleteOrder({ id });
  ElMessage.success("订单已删除");
  loadOrders();
}

onMounted(loadOrders);
</script>

<style scoped>
.order-card {
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  padding: 14px 16px;
  margin-bottom: 12px;
}
.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}
.order-num {
  font-size: 13px;
  color: #666;
}
.order-body {
  cursor: pointer;
  padding: 8px 0;
}
.order-money {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}
.order-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 10px;
}
</style>
