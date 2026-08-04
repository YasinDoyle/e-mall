<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.order.title") }}</span>
        <div class="filters">
          <el-select v-model="typeFilter" :placeholder="t('page.order.orderStatus')" clearable style="width: 150px" @change="reloadOrders">
            <el-option v-for="item in orderTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="refundFilter" :placeholder="t('page.order.refundStatus')" clearable style="width: 140px" @change="reloadOrders">
            <el-option v-for="item in refundTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="afterSaleFilter" :placeholder="t('page.order.afterSaleStatus')" clearable style="width: 160px" @change="loadAfterSales">
            <el-option v-for="item in afterSaleStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-button :loading="loadingOrders" @click="reloadAll">{{ t("common.refresh") }}</el-button>
        </div>
      </div>
    </template>

    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span>{{ t("page.order.afterSale") }}</span>
          <el-button link type="primary" @click="loadAfterSales">{{ t("common.refresh") }}</el-button>
        </div>
      </template>
      <el-table :data="afterSales" style="width: 100%" v-loading="loadingAfterSales">
        <el-table-column prop="order_num" :label="t('page.order.orderNo')" min-width="160" />
        <el-table-column :label="t('page.order.afterSaleType')" width="130">
          <template #default="{ row }">{{ afterSaleTypeText(row.type) }}</template>
        </el-table-column>
        <el-table-column :label="t('page.order.afterSaleStatus')" width="130">
          <template #default="{ row }">
            <el-tag :type="afterSaleStatusTag(row.status)">{{ afterSaleStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" :label="t('page.order.afterSaleReason')" min-width="180" show-overflow-tooltip />
        <el-table-column :label="t('page.order.afterSaleRefundAmount')" width="120">
          <template #default="{ row }">¥{{ Number(row.refund_amount || 0).toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="platform_note" :label="t('page.order.afterSaleNote')" min-width="150" show-overflow-tooltip />
        <el-table-column :label="t('common.actions')" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openAfterSaleOrderDetail(row)">{{ t("page.order.detail") }}</el-button>
            <el-button
              v-if="canIntervene(row.status)"
              size="small"
              type="warning"
              :loading="handlingAfterSaleId === row.id"
              @click="openAfterSaleHandle(row, 'intervene')"
            >
              {{ t("page.order.afterSaleIntervene") }}
            </el-button>
            <el-button
              v-if="canRefund(row.status)"
              size="small"
              type="primary"
              :loading="handlingAfterSaleId === row.id"
              @click="openAfterSaleHandle(row, 'refund')"
            >
              {{ t("page.order.afterSaleRefund") }}
            </el-button>
            <el-button
              v-if="canClose(row.status)"
              size="small"
              type="info"
              :loading="handlingAfterSaleId === row.id"
              @click="openAfterSaleHandle(row, 'close')"
            >
              {{ t("page.order.afterSaleClose") }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-table :data="orders" style="width: 100%" v-loading="loadingOrders">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_num" :label="t('page.order.orderNo')" min-width="170" />
      <el-table-column :label="t('page.order.product')" min-width="220">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div>{{ row.name || "-" }}</div>
              <div class="muted">{{ t("page.order.productId", { id: row.product_id }) }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="user_id" :label="t('page.order.buyerId')" width="90" />
      <el-table-column prop="boss_id" :label="t('page.order.sellerId')" width="90" />
      <el-table-column :label="t('common.amount')" width="110">
        <template #default="{ row }">¥{{ totalAmount(row).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column :label="t('page.order.orderStatus')" width="130">
        <template #default="{ row }">
          <el-tag :type="orderTypeTag(row.type)">{{ orderTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.order.refundStatus')" width="110">
        <template #default="{ row }">
          <el-tag :type="refundTag(row.refund_status)">{{ refundText(row.refund_status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="refund_reason" :label="t('page.order.refundReason')" min-width="160" show-overflow-tooltip />
      <el-table-column :label="t('common.actions')" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row.id)">{{ t("page.order.detail") }}</el-button>
          <el-button
            v-if="row.refund_status === 1"
            size="small"
            type="primary"
            @click="approveRefund(row)"
          >
            {{ t("page.order.refundAudit") }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      class="pager"
      @current-change="loadOrders"
    />

    <el-drawer v-model="detailDrawerVisible" :title="t('page.order.detail')" size="720px">
      <el-skeleton v-if="detailLoading" :rows="5" animated />
      <template v-else-if="detailOrder">
        <el-descriptions :column="2" border class="detail-grid">
          <el-descriptions-item :label="t('page.order.orderNo')">{{ detailOrder.order_num }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.paymentChannel')">{{ paymentChannelText(detailOrder.payment_channel) }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.trackingNo')">{{ detailOrder.tracking_no || t("common.none") }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.logisticsCompany')">{{ detailOrder.logistics_company || t("common.none") }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.quantity')">{{ detailOrder.num }} {{ t("page.order.quantityUnit") }}</el-descriptions-item>
          <el-descriptions-item :label="t('common.amount')">¥{{ totalAmount(detailOrder).toFixed(2) }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.orderStatus')">{{ orderTypeText(detailOrder.type) }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.refundStatus')">{{ refundText(detailOrder.refund_status) }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.createdAt')">{{ formatTime(detailOrder.created_at) }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.shippedAt')">{{ formatTime(detailOrder.shipped_at) || t("common.none") }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.receivedAt')">{{ formatTime(detailOrder.received_at) || t("common.none") }}</el-descriptions-item>
          <el-descriptions-item :label="t('page.order.canceledAt')">{{ formatTime(detailOrder.canceled_at) || t("common.none") }}</el-descriptions-item>
        </el-descriptions>

        <el-card class="section-card" shadow="never">
          <template #header>{{ t("page.order.operationLogs") }}</template>
          <el-timeline v-if="detailLogs.length">
            <el-timeline-item v-for="item in detailLogs" :key="item.id" :timestamp="formatTime(item.created_at)" type="primary">
              <div class="timeline-item">
                <b>{{ orderActionText(item.action) }}</b>
                <span class="muted">{{ operatorTypeText(item.operator_type) }}<template v-if="item.remark"> · {{ item.remark }}</template></span>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="t('page.order.noOperationLog')" />
        </el-card>

        <el-card class="section-card" shadow="never">
          <template #header>{{ t("page.order.logisticsTimeline") }}</template>
          <el-timeline v-if="detailLogistics.length">
            <el-timeline-item v-for="item in detailLogistics" :key="item.key" :timestamp="item.time" :type="item.done ? 'primary' : 'info'" :hollow="!item.done">
              <div class="timeline-item">
                <b>{{ item.title }}</b>
                <span v-if="item.meta" class="muted">{{ item.meta }}</span>
              </div>
            </el-timeline-item>
          </el-timeline>
          <el-empty v-else :description="t('page.order.noLogistics')" />
        </el-card>

        <el-card class="section-card" shadow="never">
          <template #header>{{ t("page.order.afterSaleList") }}</template>
          <el-empty v-if="!detailAfterSales.length" :description="t('page.order.noAfterSale')" />
          <div v-else class="after-sale-list">
            <div v-for="item in detailAfterSales" :key="item.id" class="after-sale-item">
              <div class="after-sale-head">
                <div>
                  <b>{{ afterSaleTypeText(item.type) }}</b>
                  <span class="muted"> · {{ afterSaleStatusText(item.status) }}</span>
                </div>
                <div class="after-sale-head-right">
                  <el-tag size="small" :type="afterSaleStatusTag(item.status)">{{ afterSaleStatusText(item.status) }}</el-tag>
                  <span class="amount">¥{{ Number(item.refund_amount || 0).toFixed(2) }}</span>
                </div>
              </div>
              <div class="after-sale-line">{{ item.reason }}</div>
              <div v-if="item.seller_reason" class="after-sale-line muted">{{ t("page.order.sellerReason") }}: {{ item.seller_reason }}</div>
              <div v-if="item.platform_note" class="after-sale-line muted">{{ t("page.order.platformNote") }}: {{ item.platform_note }}</div>
            </div>
          </div>
        </el-card>
      </template>
    </el-drawer>

    <el-dialog
      v-model="handleAfterSaleDialogVisible"
      :title="t('page.order.afterSaleHandle')"
      width="520px"
    >
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="t('page.order.orderNo')">
          {{ currentAfterSale?.order_num || t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('page.order.afterSaleType')">
          {{ currentAfterSale ? afterSaleTypeText(currentAfterSale.type) : t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('page.order.afterSaleStatus')">
          {{ currentAfterSale ? afterSaleStatusText(currentAfterSale.status) : t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('page.order.afterSaleRefundAmount')">
          ¥{{ Number(currentAfterSale?.refund_amount || 0).toFixed(2) }}
        </el-descriptions-item>
      </el-descriptions>
      <el-form label-width="90px" class="after-sale-handle-form">
        <el-form-item :label="t('page.order.afterSaleNote')">
          <el-input
            v-model="afterSaleHandleNote"
            type="textarea"
            :rows="4"
            maxlength="255"
            show-word-limit
            :placeholder="t('page.order.afterSaleNotePlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="handleAfterSaleDialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="afterSaleHandleSubmitting" @click="submitAfterSaleHandle">
          {{ currentAfterSaleActionLabel }}
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  approveOrderRefund,
  getAdminAfterSaleList,
  getAdminOrderDetail,
  getAdminOrderList,
  getAdminOrderLogs,
  handleAdminAfterSale,
} from "@/api";
import {
  afterSaleStatusTagType,
  afterSaleStatusText,
  afterSaleTypeText,
  orderActionText,
  orderStatusText,
  paymentChannelText,
  operatorTypeText,
  refundStatusText,
} from "@/utils/status-labels";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const orders = ref<any[]>([]);
const afterSales = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loadingOrders = ref(false);
const loadingAfterSales = ref(false);
const typeFilter = ref<number | undefined>();
const refundFilter = ref<number | undefined>();
const afterSaleFilter = ref<string | undefined>();
const handlingAfterSaleId = ref(0);
const handleAfterSaleDialogVisible = ref(false);
const afterSaleHandleSubmitting = ref(false);
const currentAfterSale = ref<any>(null);
const currentAfterSaleAction = ref<"intervene" | "refund" | "close">("intervene");
const afterSaleHandleNote = ref("");

const detailDrawerVisible = ref(false);
const detailLoading = ref(false);
const detailOrder = ref<any>(null);
const detailLogs = ref<any[]>([]);
const detailAfterSales = ref<any[]>([]);

const orderTypes = computed(() =>
  [1, 2, 3, 4, 5, 6, 7].map((value) => ({ value, label: orderStatusText(value) })),
);
const refundTypes = computed(() =>
  [0, 1, 2].map((value) => ({ value, label: refundStatusText(value) })),
);
const afterSaleStatusOptions = computed(() =>
  [
    "requested",
    "seller_approved",
    "seller_rejected",
    "platform_intervening",
    "refunded",
    "closed",
  ].map((value) => ({ value, label: afterSaleStatusText(value) })),
);
const detailLogistics = computed(() => {
  const order = detailOrder.value;
  if (!order) return [];
  return [
    {
      key: "created",
      title: t("page.order.actionCreate"),
      meta: formatTime(order.created_at),
      time: formatTime(order.created_at),
      done: true,
    },
    {
      key: "paid",
      title: t("page.order.actionPay"),
      meta: order.payment_channel ? paymentChannelText(order.payment_channel) : "",
      time: formatTime(order.paid_at),
      done: Boolean(order.paid_at),
    },
    {
      key: "shipped",
      title: t("page.order.actionShip"),
      meta: order.tracking_no
        ? `${order.logistics_company || t("common.none")} · ${order.tracking_no}`
        : "",
      time: formatTime(order.shipped_at),
      done: Boolean(order.shipped_at),
    },
    {
      key: "received",
      title: t("page.order.actionReceive"),
      meta: formatTime(order.received_at),
      time: formatTime(order.received_at),
      done: Boolean(order.received_at),
    },
    {
      key: "canceled",
      title: t("page.order.actionCancel"),
      meta: formatTime(order.canceled_at),
      time: formatTime(order.canceled_at),
      done: Boolean(order.canceled_at),
    },
  ].filter((item) => item.done || item.key === "created");
});
const currentAfterSaleActionLabel = computed(
  () =>
    ({
      intervene: t("page.order.afterSaleIntervene"),
      refund: t("page.order.afterSaleRefund"),
      close: t("page.order.afterSaleClose"),
    })[currentAfterSaleAction.value],
);

function orderTypeText(type: number) {
  return orderStatusText(type);
}

function orderTypeTag(type: number) {
  return ({ 1: "info", 2: "warning", 3: "primary", 4: "success", 5: "danger", 6: "info", 7: "info" } as any)[type] ?? "info";
}

function refundText(status: number) {
  return refundStatusText(status);
}

function refundTag(status: number) {
  return ({ 0: "info", 1: "warning", 2: "success" } as any)[status] ?? "info";
}

function afterSaleStatusTag(status: string) {
  return afterSaleStatusTagType(status);
}

function totalAmount(row: any) {
  return Number(row.money || 0) * Number(row.num || 1);
}

function formatTime(value: any) {
  if (!value) return "";
  const date = typeof value === "number" ? new Date(value * 1000) : new Date(value);
  return Number.isNaN(date.getTime()) ? "" : date.toLocaleString();
}

function reloadOrders() {
  page.value = 1;
  loadOrders();
}

function reloadAfterSales() {
  loadAfterSales();
}

async function loadOrders() {
  loadingOrders.value = true;
  try {
    const res: any = await getAdminOrderList({
      page_num: page.value,
      page_size: pageSize,
      ...(typeFilter.value !== undefined ? { type: typeFilter.value } : {}),
      ...(refundFilter.value !== undefined ? { refund_status: refundFilter.value } : {}),
    });
    orders.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loadingOrders.value = false;
  }
}

async function loadAfterSales() {
  loadingAfterSales.value = true;
  try {
    const res: any = await getAdminAfterSaleList({
      page_num: 1,
      page_size: 20,
      ...(afterSaleFilter.value ? { status: afterSaleFilter.value } : {}),
    });
    afterSales.value = res.data?.item ?? [];
  } finally {
    loadingAfterSales.value = false;
  }
}

async function reloadAll() {
  await Promise.all([loadOrders(), loadAfterSales()]);
}

async function openDetail(orderId: number) {
  detailDrawerVisible.value = true;
  detailLoading.value = true;
  try {
    const [detailRes, logsRes, afterSaleRes]: any[] = await Promise.all([
      getAdminOrderDetail({ order_id: orderId }),
      getAdminOrderLogs({ order_id: orderId }),
      getAdminAfterSaleList({ order_id: orderId, page_num: 1, page_size: 20 }),
    ]);
    detailOrder.value = detailRes.data ?? null;
    detailLogs.value = logsRes.data?.item ?? logsRes.data ?? [];
    detailAfterSales.value = afterSaleRes.data?.item ?? [];
  } finally {
    detailLoading.value = false;
  }
}

function openAfterSaleOrderDetail(row: any) {
  openDetail(Number(row.order_id));
}

function canIntervene(status: string) {
  return status === "seller_rejected";
}

function canRefund(status: string) {
  return status === "seller_approved" || status === "platform_intervening";
}

function canClose(status: string) {
  return status === "requested" || status === "seller_rejected" || status === "platform_intervening";
}

function openAfterSaleHandle(row: any, action: "intervene" | "refund" | "close") {
  currentAfterSale.value = row;
  currentAfterSaleAction.value = action;
  afterSaleHandleNote.value = row.platform_note || row.seller_reason || "";
  handleAfterSaleDialogVisible.value = true;
}

async function submitAfterSaleHandle() {
  if (!currentAfterSale.value) return;
  if (currentAfterSaleAction.value === "close" && !afterSaleHandleNote.value.trim()) {
    return ElMessage.warning(t("page.order.afterSaleCloseNoteRequired"));
  }
  handlingAfterSaleId.value = currentAfterSale.value.id;
  afterSaleHandleSubmitting.value = true;
  try {
    await handleAdminAfterSale({
      after_sale_id: currentAfterSale.value.id,
      action: currentAfterSaleAction.value,
      note: afterSaleHandleNote.value.trim(),
    });
    ElMessage.success(t("page.order.afterSaleHandleSuccess"));
    handleAfterSaleDialogVisible.value = false;
    await reloadAll();
    if (detailDrawerVisible.value && detailOrder.value?.id === currentAfterSale.value.order_id) {
      await openDetail(currentAfterSale.value.order_id);
    }
  } finally {
    afterSaleHandleSubmitting.value = false;
    handlingAfterSaleId.value = 0;
  }
}

async function approveRefund(row: any) {
  await ElMessageBox.confirm(
    t("page.order.refundPrompt", {
      orderNo: row.order_num,
      amount: totalAmount(row).toFixed(2),
    }),
    t("page.order.refundAudit"),
    {
      type: "warning",
      confirmButtonText: t("page.order.refundConfirm"),
      cancelButtonText: t("common.cancel"),
    },
  );
  await approveOrderRefund({ order_id: row.id });
  ElMessage.success(t("page.order.refundSuccess"));
  requestAdminPendingCountsRefresh();
  await reloadAll();
}

onMounted(reloadAll);
</script>

<style scoped>
.card-header,
.filters,
.product-cell,
.section-header,
.after-sale-head,
.after-sale-head-right,
.timeline-item {
  display: flex;
  align-items: center;
}
.card-header,
.section-header {
  justify-content: space-between;
  gap: 12px;
}
.filters {
  gap: 8px;
  flex-wrap: wrap;
}
.section-card {
  margin-bottom: 16px;
}
.product-cell {
  gap: 10px;
}
.product-img {
  width: 46px;
  height: 46px;
  object-fit: cover;
  border-radius: 4px;
}
.muted {
  color: #909399;
  font-size: 12px;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
.detail-grid {
  margin-bottom: 16px;
}
.after-sale-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.after-sale-item {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  padding: 12px;
}
.after-sale-head {
  justify-content: space-between;
  gap: 12px;
}
.after-sale-head-right {
  gap: 10px;
}
.after-sale-line {
  margin-top: 6px;
}
.amount {
  color: #f56c6c;
  font-weight: 600;
}
</style>
