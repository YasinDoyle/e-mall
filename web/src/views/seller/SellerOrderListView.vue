<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>{{ t("sellerCenter.order.title") }}</span>
        <el-button :loading="loading" @click="reloadAll">{{ t("common.refresh") }}</el-button>
      </div>
    </template>

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
          <el-tag :type="statusTagType(row.type)">
            {{ statusText(row.type) }}
          </el-tag>
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
      <el-table-column :label="t('sellerCenter.order.actions')" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.type === 2"
            size="small"
            type="primary"
            @click="openShip(row)"
          >
            {{ t("sellerCenter.order.ship") }}
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
      <el-form label-width="80px">
        <el-form-item :label="t('sellerCenter.order.orderNo')">
          <span>{{ currentOrder?.order_num }}</span>
        </el-form-item>
        <el-form-item :label="t('sellerCenter.order.trackingNo')">
          <el-input
            v-model="trackingNo"
            maxlength="64"
            :placeholder="t('sellerCenter.order.trackingNoPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="shipping" @click="handleShip">
          {{ t("sellerCenter.order.confirmShip") }}
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import {
  getSellerOrderList,
  getSellerSettlementSummary,
  shipOrder,
} from "@/api/order";
import { getSellerAccountSummary } from "@/api/seller";
import { useSellerStore } from "@/stores/seller";
import { orderStatusText } from "@/utils/status-labels";

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
const shipping = ref(false);
const currentOrder = ref<any>(null);
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

async function loadList() {
  loading.value = true;
  try {
    const params: any = { page_num: page.value, page_size: pageSize };
    if (activeTab.value !== "0") params.type = Number(activeTab.value);
    const res: any = await getSellerOrderList(params);
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function reloadAll() {
  await Promise.all([loadSummary(), loadList()]);
}

function handleTabChange() {
  page.value = 1;
  loadList();
}

function openShip(row: any) {
  currentOrder.value = row;
  trackingNo.value = row.tracking_no || "";
  shipDialogVisible.value = true;
}

async function handleShip() {
  if (!currentOrder.value) return;
  if (!trackingNo.value.trim()) {
    return ElMessage.warning(t("sellerCenter.order.trackingNoPlaceholder"));
  }
  shipping.value = true;
  try {
    await shipOrder({
      order_id: currentOrder.value.id,
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
.product-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}
.header {
  justify-content: space-between;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 12px;
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
</style>
