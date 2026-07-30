<template>
  <el-card>
    <template #header>{{ t("orderList.title") }}</template>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="t('common.all')" name="0" />
      <el-tab-pane :label="t('status.order.unpaid')" name="1" />
      <el-tab-pane :label="t('status.order.pendingShipment')" name="2" />
      <el-tab-pane :label="t('status.order.shipped')" name="3" />
      <el-tab-pane :label="t('status.order.completed')" name="4" />
    </el-tabs>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!orders.length" :description="t('orderList.empty')" />

    <template v-else>
      <div v-for="order in orders" :key="order.id" class="order-card">
        <div class="order-header">
          <span class="order-num">{{ t("orderList.orderNo", { num: order.order_num }) }}</span>
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
            <div class="order-name">{{ order.name || t("orderList.product") }}</div>
            <div class="order-meta">
              {{ t("orderList.quantity", { count: order.num }) }}
              <span v-if="order.tracking_no">
                · {{ t("orderList.trackingNo", { no: order.tracking_no }) }}
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
            >{{ t("orderList.confirmReceive") }}</el-button
          >
          <el-button
            v-if="order.type === 1 || order.type === 4 || order.type === 6"
            size="small"
            type="danger"
            @click="handleDelete(order.id)"
            >{{ t("orderList.delete") }}</el-button
          >
          <el-button
            size="small"
            @click="$router.push(`/user/orders/${order.id}`)"
          >
            {{ t("orderList.detail") }}
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
import { useI18n } from "vue-i18n";
import { getOrderList, deleteOrder, receiveOrder } from "@/api/order";
import { orderStatusText } from "@/utils/status-labels";

const { t } = useI18n();
const activeTab = ref("0");
const orders = ref<any[]>([]);
const page = ref(1);
const pageSize = 10;
const total = ref(0);
const loading = ref(false);

const statusText = orderStatusText;
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
  ({ 1: t("orderList.refundRequested"), 2: t("status.refund.refunded") })[status] ??
  t("orderList.refundProcessing");

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
    await ElMessageBox.confirm(t("orderList.receiveConfirm"), t("dialog.warningTitle"), {
      type: "warning",
    });
    await receiveOrder({ order_id: id });
    ElMessage.success(t("orderList.receiveSuccess"));
    loadOrders();
  } catch {}
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t("orderList.deleteConfirm"), t("dialog.warningTitle"), {
      type: "warning",
    });
    await deleteOrder({ order_id: id });
    ElMessage.success(t("orderList.deleteSuccess"));
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
