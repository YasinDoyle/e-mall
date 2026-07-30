<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.settlement.title") }}</span>
        <div class="filters">
          <el-input-number
            v-model="sellerID"
            :min="0"
            :step="1"
            :placeholder="t('page.settlement.sellerId')"
            controls-position="right"
          />
          <el-select v-model="statusFilter" clearable :placeholder="t('page.settlement.statusPlaceholder')" style="width: 140px" @change="reload">
            <el-option :label="t('status.settlement.pending')" value="pending" />
            <el-option :label="t('status.settlement.generated')" value="generated" />
            <el-option :label="t('status.settlement.paid')" value="paid" />
            <el-option :label="t('status.settlement.refunded')" value="refunded" />
          </el-select>
          <el-button :loading="loading" @click="reloadAll">{{ t("common.refresh") }}</el-button>
          <el-button :loading="backfilling" @click="backfillAccount">{{ t("page.settlement.backfillAccount") }}</el-button>
          <el-button type="primary" @click="generate">{{ t("page.settlement.generateBySeller") }}</el-button>
        </div>
      </div>
    </template>

    <div class="summary-grid">
      <div class="summary-item">
        <span class="summary-label">{{ t("page.settlement.platformRevenue") }}</span>
        <b>¥{{ money(overview.platform_revenue) }}</b>
      </div>
    </div>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="seller_id" :label="t('page.settlement.sellerId')" width="90" />
      <el-table-column prop="order_num" :label="t('page.settlement.orderNo')" min-width="170" />
      <el-table-column :label="t('page.settlement.grossAmount')" width="110">
        <template #default="{ row }">¥{{ money(row.gross_amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('page.settlement.commission')" width="130">
        <template #default="{ row }">
          ¥{{ money(row.commission_amount) }}
          <span class="muted">({{ percent(row.commission_rate) }})</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.settlement.settlementAmount')" width="120">
        <template #default="{ row }">¥{{ money(row.settlement_amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">{{ t("common.flow") }}</el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="primary"
            @click="generateOne(row)"
          >
            {{ t("page.settlement.generateOne") }}
          </el-button>
          <el-button
            v-if="row.status === 'generated'"
            size="small"
            type="primary"
            @click="markPaid(row)"
          >
            {{ t("page.settlement.markPaid") }}
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

    <el-dialog v-model="detailVisible" :title="t('page.settlement.flowTitle')" width="760px">
      <el-table :data="flows" size="small">
        <el-table-column prop="flow_no" :label="t('page.settlement.flowNo')" min-width="190" />
        <el-table-column prop="flow_type" :label="t('common.type')" width="150" />
        <el-table-column prop="direction" :label="t('common.direction')" width="80" />
        <el-table-column :label="t('common.amount')" width="100">
          <template #default="{ row }">¥{{ money(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="remark" :label="t('common.remark')" min-width="140" />
      </el-table>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  backfillAdminSellerAccount,
  generateAdminSettlement,
  generateOneAdminSettlement,
  getAdminSettlementDetail,
  getAdminSettlementList,
  markAdminSettlementPaid,
} from "@/api/settlement";
import { getStatsOverview } from "@/api";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const list = ref<any[]>([]);
const flows = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const backfilling = ref(false);
const sellerID = ref<number | undefined>();
const statusFilter = ref<string | undefined>();
const detailVisible = ref(false);
const overview = ref({
  platform_revenue: 0,
});

function money(value: number) {
  return Number(value || 0).toFixed(2);
}

function percent(value: number) {
  return `${(Number(value || 0) * 100).toFixed(1)}%`;
}

function statusText(status: string) {
  return (
    {
      pending: t("status.settlement.pending"),
      generated: t("status.settlement.generated"),
      paid: t("status.settlement.paid"),
      refunded: t("status.settlement.refunded"),
    } as any
  )[status] ?? t("common.unknown");
}

function statusTag(status: string) {
  return (
    {
      pending: "warning",
      generated: "primary",
      paid: "success",
      refunded: "info",
    } as any
  )[status] ?? "info";
}

function reload() {
  page.value = 1;
  loadList();
}

async function loadOverview() {
  const res: any = await getStatsOverview();
  overview.value = {
    platform_revenue: Number(res.data?.platform_revenue || 0),
  };
}

async function loadList(withLoading = true) {
  if (withLoading) {
    loading.value = true;
  }
  try {
    const res: any = await getAdminSettlementList({
      page_num: page.value,
      page_size: pageSize,
      ...(sellerID.value ? { seller_id: sellerID.value } : {}),
      ...(statusFilter.value ? { status: statusFilter.value } : {}),
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    if (withLoading) {
      loading.value = false;
    }
  }
}

async function reloadAll() {
  loading.value = true;
  try {
    await Promise.all([loadOverview(), loadList(false)]);
  } finally {
    loading.value = false;
  }
}

async function generate() {
  if (!sellerID.value) {
    return ElMessage.warning(t("page.settlement.sellerIdRequired"));
  }
  const res: any = await generateAdminSettlement({ seller_id: sellerID.value });
  ElMessage.success(t("page.settlement.generateSuccess", { count: res.data?.count ?? 0 }));
  requestAdminPendingCountsRefresh();
  reload();
}

async function backfillAccount() {
  await ElMessageBox.confirm(t("page.settlement.backfillConfirm"), t("common.notice"), {
    type: "warning",
  });
  backfilling.value = true;
  try {
    const res: any = await backfillAdminSellerAccount();
    const data = res.data || {};
    ElMessage.success(
      t("page.settlement.backfillSuccess", {
        settlementCount: data.settlement_count ?? 0,
        sellerCount: data.seller_count ?? 0,
        amount: money(data.amount ?? 0),
      }),
    );
    loadList();
  } finally {
    backfilling.value = false;
  }
}

async function generateOne(row: any) {
  await ElMessageBox.confirm(
    t("page.settlement.generateOneConfirm", { orderNo: row.order_num }),
    t("common.notice"),
    { type: "warning" },
  );
  await generateOneAdminSettlement({ id: row.id });
  ElMessage.success(t("page.settlement.generateOneSuccess"));
  requestAdminPendingCountsRefresh();
  loadList();
}

async function markPaid(row: any) {
  await ElMessageBox.confirm(
    t("page.settlement.markPaidConfirm", { id: row.id }),
    t("common.notice"),
    { type: "warning" },
  );
  await markAdminSettlementPaid({ id: row.id });
  ElMessage.success(t("page.settlement.markPaidSuccess"));
  loadList();
}

async function openDetail(row: any) {
  const res: any = await getAdminSettlementDetail({ id: row.id });
  flows.value = res.data?.item ?? [];
  detailVisible.value = true;
}

onMounted(reloadAll);
</script>

<style scoped>
.card-header,
.filters {
  display: flex;
  align-items: center;
  gap: 10px;
}
.card-header {
  justify-content: space-between;
}
.muted {
  color: #909399;
  font-size: 12px;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}
.summary-item {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  padding: 14px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
.summary-label {
  color: #909399;
  font-size: 13px;
}
</style>
