<template>
  <el-card>
    <template #header>我的订单</template>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="0" />
      <el-tab-pane label="待支付" name="1" />
      <el-tab-pane label="待发货" name="2" />
      <el-tab-pane label="已发货" name="3" />
      <el-tab-pane label="已完成" name="4" />
    </el-tabs>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!orders.length" description="暂无订单" />

    <template v-else>
      <div v-for="order in orders" :key="order.id" class="order-card">
        <div class="order-header">
          <span class="order-num">订单号：{{ order.order_num }}</span>
          <div class="order-tags">
            <el-tag :type="statusTagType(order.type)">
              {{ statusText(order.type) }}
            </el-tag>
            <el-tag v-if="order.refund_status" type="danger">
              {{ refundText(order.refund_status) }}
            </el-tag>
          </div>
        </div>
        <div
          class="order-body"
          @click="$router.push(`/user/orders/${order.id}`)"
        >
          <img :src="order.img_path" class="order-img" />
          <div class="order-info">
            <div class="order-name">{{ order.name || "商品" }}</div>
            <div class="order-meta">
              {{ order.num }} 件
              <span v-if="order.tracking_no">
                · 物流单号：{{ order.tracking_no }}
              </span>
            </div>
          </div>
          <div class="order-money">¥{{ totalAmount(order) }}</div>
        </div>
        <div class="order-actions">
          <el-button
            v-if="order.type === 3"
            size="small"
            type="primary"
            @click="handleReceive(order.id)"
            >确认收货</el-button
          >
          <el-button
            v-if="order.type === 1 || order.type === 4 || order.type === 6"
            size="small"
            type="danger"
            @click="handleDelete(order.id)"
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
    </template>

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
const loading = ref(false);

const statusText = (type: number) =>
  ({
    1: "待支付",
    2: "待发货",
    3: "已发货",
    4: "已完成",
    5: "退款中",
    6: "已退款",
  })[type] ?? "未知";
const statusTagType = (type: number) =>
  (({
    1: "warning",
    2: "primary",
    3: "warning",
    4: "success",
    5: "danger",
    6: "info",
  })[type] ?? "info") as any;
const refundText = (status: number) =>
  ({ 1: "退款申请中", 2: "已退款" })[status] ?? "退款处理中";

function totalAmount(order: any) {
  return (Number(order.money || 0) * Number(order.num || 0)).toFixed(2);
}

async function loadOrders() {
  loading.value = true;
  try {
    const params: any = { page_num: page.value, page_size: pageSize };
    if (activeTab.value !== "0") params.type = Number(activeTab.value);
    const res: any = await getOrderList(params);
    orders.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } catch {
    orders.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

function handleTabChange() {
  page.value = 1;
  loadOrders();
}

async function handleReceive(id: number) {
  try {
    await ElMessageBox.confirm("确认已收到货物？", "提示", { type: "warning" });
    await receiveOrder({ order_id: id });
    ElMessage.success("已确认收货");
    loadOrders();
  } catch {}
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm("确认删除该订单？", "提示", { type: "warning" });
    await deleteOrder({ order_id: id });
    ElMessage.success("订单已删除");
    loadOrders();
  } catch {}
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
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 0;
}
.order-img {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}
.order-info {
  flex: 1;
  min-width: 0;
}
.order-name {
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.order-meta {
  margin-top: 6px;
  color: #999;
  font-size: 12px;
}
.order-money {
  font-size: 18px;
  font-weight: bold;
  color: #f56c6c;
}
.order-tags {
  display: flex;
  gap: 8px;
}
.order-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 10px;
}
</style>
