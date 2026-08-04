<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>{{ t("sellerCenter.order.title") }}</span>
        <el-button :loading="loading" @click="reloadAll">{{ t("common.refresh") }}</el-button>
      </div>
    </template>

    <el-alert
      v-if="pendingAfterSaleCount > 0"
      class="pending-after-sale-alert"
      type="warning"
      show-icon
      :closable="false"
      :title="t('sellerCenter.order.pendingAfterSaleNotice', { count: pendingAfterSaleCount })"
    />

    <div class="summary-grid">
      <div class="summary-item account">
        <span class="summary-label">{{ t("sellerCenter.order.availableBalance") }}</span>
        <b>¥{{ money(accountSummary.available_balance) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.order.pendingSettlement") }}</span>
        <b>¥{{ money(settlementSummary.pending_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.order.generatedSettlement") }}</span>
        <b>¥{{ money(settlementSummary.generated_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.order.paidSettlement") }}</span>
        <b>¥{{ money(settlementSummary.paid_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.order.refundedSettlement") }}</span>
        <b>¥{{ money(settlementSummary.refunded_amount) }}</b>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="t('common.all')" name="0" />
      <el-tab-pane :label="t('sellerCenter.order.pendingShipment')" name="2" />
      <el-tab-pane :label="t('sellerCenter.order.shipped')" name="3" />
      <el-tab-pane :label="t('sellerCenter.order.completed')" name="4" />
      <el-tab-pane :label="t('sellerCenter.order.refunding')" name="5" />
      <el-tab-pane :label="t('sellerCenter.order.refunded')" name="6" />
    </el-tabs>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="order_num" :label="t('sellerCenter.order.orderNo')" min-width="150" />
      <el-table-column :label="t('sellerCenter.order.product')" min-width="240">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div class="name">{{ row.name || t("sellerCenter.order.product") }}</div>
              <div class="muted">{{ t("sellerCenter.order.buyerId", { id: row.user_id }) }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.order.tradeInfo')" min-width="210">
        <template #default="{ row }">
          <div>{{ t("sellerCenter.order.orderAmount", { amount: totalAmount(row) }) }}</div>
          <div class="muted">
            {{
              t("sellerCenter.order.commissionIncome", {
                commission: money(row.commission_amount),
                income: money(row.settlement_amount),
              })
            }}
          </div>
          <el-tag size="small" :type="settlementTag(row.settlement_status)">
            {{ settlementText(row.settlement_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="num" :label="t('sellerCenter.order.quantity')" width="80" />
      <el-table-column :label="t('sellerCenter.order.status')" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.type)">{{ statusText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.order.address')" min-width="260">
        <template #default="{ row }">
          <div>{{ row.address_name }} {{ row.address_phone }}</div>
          <div class="muted">{{ row.address }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="tracking_no" :label="t('sellerCenter.order.trackingNo')" min-width="150">
        <template #default="{ row }">{{ row.tracking_no || "-" }}</template>
      </el-table-column>
      <el-table-column prop="logistics_company" :label="t('sellerCenter.order.logisticsCompany')" min-width="130">
        <template #default="{ row }">{{ row.logistics_company || "-" }}</template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.order.actions')" width="180" fixed="right">
        <template #default="{ row }">
          <el-button v-if="row.type === 2" size="small" type="primary" @click="openShip(row)">
            {{ t("sellerCenter.order.ship") }}
          </el-button>
          <el-button
            v-if="row.type >= 2 && row.type <= 6"
            size="small"
            :type="afterSaleActionType(row)"
            @click="openAfterSale(row)"
          >
            {{ t("sellerCenter.order.afterSale") }}
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
      @current-change="loadList"
    />

    <el-dialog v-model="shipDialogVisible" :title="t('sellerCenter.order.shipDialog')" width="420px">
      <el-form label-width="90px">
        <el-form-item :label="t('sellerCenter.order.orderNo')">
          <span>{{ currentOrder?.order_num }}</span>
        </el-form-item>
        <el-form-item :label="t('sellerCenter.order.logisticsCompany')">
          <el-input v-model="logisticsCompany" maxlength="64" :placeholder="t('sellerCenter.order.logisticsCompanyPlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('sellerCenter.order.trackingNo')">
          <el-input v-model="trackingNo" maxlength="64" :placeholder="t('sellerCenter.order.trackingNoPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="shipping" @click="handleShip">{{ t("sellerCenter.order.confirmShip") }}</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="afterSaleDialogVisible" :title="t('sellerCenter.order.afterSaleDialog')" width="760px">
      <el-descriptions :column="2" border class="detail-grid">
        <el-descriptions-item :label="t('sellerCenter.order.orderNo')">
          {{ currentOrder?.order_num || t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('sellerCenter.order.logisticsCompany')">
          {{ currentOrder?.logistics_company || t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('sellerCenter.order.trackingNo')">
          {{ currentOrder?.tracking_no || t("common.none") }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('sellerCenter.order.status')">
          {{ settlementText(currentOrder?.settlement_status) }}
        </el-descriptions-item>
      </el-descriptions>

      <el-empty v-if="!afterSaleList.length && !afterSaleLoading" :description="t('sellerCenter.order.afterSaleEmpty')" />
      <el-skeleton v-if="afterSaleLoading" :rows="4" animated />
      <div v-else class="after-sale-list">
        <div v-for="item in afterSaleList" :key="item.id" class="after-sale-item">
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
          <div v-if="item.seller_reason" class="after-sale-line muted">
            {{ t("sellerCenter.order.afterSaleSellerReason") }}: {{ item.seller_reason }}
          </div>
          <div v-if="item.platform_note" class="after-sale-line muted">
            {{ t("sellerCenter.order.afterSalePlatformNote") }}: {{ item.platform_note }}
          </div>
          <div v-if="item.status === 'requested'" class="after-sale-actions">
            <el-button size="small" type="primary" :loading="handlingAfterSaleId === item.id" @click="handleAfterSale(item, 'approve')">
              {{ t("sellerCenter.order.afterSaleApprove") }}
            </el-button>
            <el-button size="small" type="danger" :loading="handlingAfterSaleId === item.id" @click="handleAfterSaleReject(item)">
              {{ t("sellerCenter.order.afterSaleReject") }}
            </el-button>
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="afterSaleDialogVisible = false">{{ t("common.cancel") }}</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { useI18n } from "vue-i18n";
import {
  getSellerAfterSaleList,
  getSellerOrderList,
  getSellerSettlementSummary,
  handleSellerAfterSale,
  shipOrder,
} from "@/api/order";
import { getSellerAccountSummary } from "@/api/seller";
import { useSellerStore } from "@/stores/seller";
import {
  afterSaleStatusTagType,
  afterSaleStatusText,
  afterSaleTypeText,
  orderStatusText,
} from "@/utils/status-labels";

const sellerStore = useSellerStore();
const { t } = useI18n();

const accountSummary = ref({
  available_balance: 0,
  frozen_balance: 0,
  total_income: 0,
  total_withdrawn: 0,
});
const settlementSummary = ref({
  available_amount: 0,
  pending_amount: 0,
  generated_amount: 0,
  paid_amount: 0,
  refunded_amount: 0,
});
const activeTab = ref("0");
const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 10;
const total = ref(0);
const loading = ref(false);
const shipDialogVisible = ref(false);
const afterSaleDialogVisible = ref(false);
const shipping = ref(false);
const afterSaleLoading = ref(false);
const handlingAfterSaleId = ref(0);
const pendingAfterSaleCount = ref(0);
const activeAfterSaleByOrderId = ref<Record<number, boolean>>({});
const currentOrder = ref<any>(null);
const afterSaleList = ref<any[]>([]);
const logisticsCompany = ref("");
const trackingNo = ref("");

const statusText = orderStatusText;
const statusTagType = (type: number) =>
  (({
    1: "warning",
    2: "primary",
    3: "warning",
    4: "success",
    5: "danger",
    6: "info",
    7: "info",
  })[type] ?? "info") as any;

function totalAmount(order: any) {
  return (Number(order.money || 0) * Number(order.num || 0)).toFixed(2);
}

function money(value: number) {
  return Number(value || 0).toFixed(2);
}

function settlementText(status?: string) {
  return (
    {
      pending: t("sellerCenter.order.settlementPending"),
      generated: t("sellerCenter.order.settlementGenerated"),
      paid: t("sellerCenter.order.settlementPaid"),
      refunded: t("sellerCenter.order.settlementRefunded"),
    } as any
  )[status || ""] ?? t("sellerCenter.order.settlementNone");
}

function settlementTag(status?: string) {
  return (
    {
      pending: "warning",
      generated: "primary",
      paid: "success",
      refunded: "info",
    } as any
  )[status || ""] ?? "info";
}

function afterSaleStatusTag(status?: string) {
  return afterSaleStatusTagType(status || "");
}

function afterSaleActionType(row: any) {
  return activeAfterSaleByOrderId.value[Number(row.id)] ? "warning" : "";
}

async function loadSummary() {
  const [accountRes, settlementRes]: any[] = await Promise.all([
    getSellerAccountSummary(),
    getSellerSettlementSummary(),
  ]);
  accountSummary.value = {
    available_balance: Number(accountRes.data?.available_balance || 0),
    frozen_balance: Number(accountRes.data?.frozen_balance || 0),
    total_income: Number(accountRes.data?.total_income || 0),
    total_withdrawn: Number(accountRes.data?.total_withdrawn || 0),
  };
  settlementSummary.value = {
    available_amount: Number(accountRes.data?.available_balance || 0),
    pending_amount: Number(settlementRes.data?.pending_amount || 0),
    generated_amount: Number(settlementRes.data?.generated_amount || 0),
    paid_amount: Number(settlementRes.data?.paid_amount || 0),
    refunded_amount: Number(settlementRes.data?.refunded_amount || 0),
  };
}

async function loadPendingAfterSaleCount() {
  const res: any = await getSellerAfterSaleList({
    status: "requested",
    page_num: 1,
    page_size: 1,
  });
  pendingAfterSaleCount.value = Number(res.data?.total ?? (res.data?.item ?? []).length);
}

async function loadActiveAfterSaleMap() {
  if (!list.value.length) {
    activeAfterSaleByOrderId.value = {};
    return;
  }
  const activeStatuses = [
    "requested",
    "seller_approved",
    "seller_rejected",
    "platform_intervening",
  ];
  const results: any[] = await Promise.all(
    activeStatuses.map((status) =>
      getSellerAfterSaleList({ status, page_num: 1, page_size: 100 }).catch(() => null),
    ),
  );
  const visibleOrderIds = new Set(list.value.map((order) => Number(order.id)));
  const next: Record<number, boolean> = {};
  for (const result of results) {
    for (const item of result?.data?.item ?? []) {
      const orderId = Number(item.order_id);
      if (visibleOrderIds.has(orderId)) {
        next[orderId] = true;
      }
    }
  }
  activeAfterSaleByOrderId.value = next;
}

async function loadList() {
  loading.value = true;
  try {
    const params: any = { page_num: page.value, page_size: pageSize };
    if (activeTab.value !== "0") params.type = Number(activeTab.value);
    const res: any = await getSellerOrderList(params);
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
    await loadActiveAfterSaleMap();
  } finally {
    loading.value = false;
  }
}

async function reloadAll() {
  await Promise.all([loadSummary(), loadList(), loadPendingAfterSaleCount()]);
}

function handleTabChange() {
  page.value = 1;
  loadList();
}

function openShip(row: any) {
  currentOrder.value = row;
  logisticsCompany.value = row.logistics_company || "";
  trackingNo.value = row.tracking_no || "";
  shipDialogVisible.value = true;
}

async function loadAfterSales(orderId: number) {
  afterSaleLoading.value = true;
  try {
    const res: any = await getSellerAfterSaleList({
      order_id: orderId,
      page_num: 1,
      page_size: 20,
    });
    afterSaleList.value = res.data?.item ?? [];
  } finally {
    afterSaleLoading.value = false;
  }
}

async function openAfterSale(row: any) {
  currentOrder.value = row;
  afterSaleDialogVisible.value = true;
  await loadAfterSales(row.id);
}

async function handleAfterSale(item: any, action: "approve" | "reject") {
  handlingAfterSaleId.value = item.id;
  try {
    await handleSellerAfterSale({
      after_sale_id: item.id,
      action,
    });
    ElMessage.success(t("sellerCenter.order.afterSaleHandleSuccess"));
    await loadAfterSales(currentOrder.value.id);
    await reloadAll();
  } finally {
    handlingAfterSaleId.value = 0;
  }
}

async function handleAfterSaleReject(item: any) {
  const result = await ElMessageBox.prompt(
    t("sellerCenter.order.afterSaleReasonPlaceholder"),
    t("sellerCenter.order.afterSaleReject"),
    {
      inputType: "textarea",
      inputPlaceholder: t("sellerCenter.order.afterSaleReasonPlaceholder"),
      inputValidator: (v: string) => !!v.trim(),
      inputErrorMessage: t("sellerCenter.order.afterSaleReasonRequired"),
    },
  ).catch(() => null);
  if (!result) return;
  handlingAfterSaleId.value = item.id;
  try {
    await handleSellerAfterSale({
      after_sale_id: item.id,
      action: "reject",
      reason: String(result.value || "").trim(),
    });
    ElMessage.success(t("sellerCenter.order.afterSaleHandleSuccess"));
    await loadAfterSales(currentOrder.value.id);
    await reloadAll();
  } finally {
    handlingAfterSaleId.value = 0;
  }
}

async function handleShip() {
  if (!currentOrder.value) return;
  if (!logisticsCompany.value.trim()) {
    return ElMessage.warning(t("sellerCenter.order.logisticsCompanyPlaceholder"));
  }
  if (!trackingNo.value.trim()) {
    return ElMessage.warning(t("sellerCenter.order.trackingNoPlaceholder"));
  }
  shipping.value = true;
  try {
    await shipOrder({
      order_id: currentOrder.value.id,
      logistics_company: logisticsCompany.value.trim(),
      tracking_no: trackingNo.value.trim(),
    });
    ElMessage.success(t("sellerCenter.order.shippedSuccess"));
    shipDialogVisible.value = false;
    await reloadAll();
  } finally {
    shipping.value = false;
  }
}

onMounted(async () => {
  await sellerStore.loadProfile({ silentError: true });
  await reloadAll();
});
</script>

<style scoped>
.header,
.product-cell,
.after-sale-head,
.after-sale-head-right {
  display: flex;
  align-items: center;
}
.header {
  justify-content: space-between;
  gap: 10px;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 12px;
}
.pending-after-sale-alert {
  margin-bottom: 12px;
}
.summary-item {
  border: 1px solid #ebeef5;
  border-radius: 6px;
  padding: 12px;
}
.summary-item b {
  display: block;
  margin-top: 6px;
  font-size: 18px;
  color: #303133;
}
.summary-label {
  color: #909399;
  font-size: 12px;
}
.summary-item.account {
  border-color: #79bbff;
  background: #ecf5ff;
}
.product-img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}
.name {
  color: #303133;
  font-weight: 500;
}
.muted {
  color: #909399;
  font-size: 12px;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
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
.detail-grid {
  margin-bottom: 16px;
}
</style>
