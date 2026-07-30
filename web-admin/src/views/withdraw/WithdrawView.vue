<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.withdraw.title") }}</span>
        <div class="filters">
          <el-input-number
            v-model="sellerId"
            :min="0"
            :step="1"
            :placeholder="t('page.withdraw.sellerId')"
            controls-position="right"
            @change="reload"
          />
          <el-select
            v-model="statusFilter"
            clearable
            :placeholder="t('page.withdraw.statusPlaceholder')"
            style="width: 140px"
            @change="reload"
          >
            <el-option :label="t('status.withdraw.pending')" value="pending" />
            <el-option :label="t('status.withdraw.approved')" value="approved" />
            <el-option :label="t('status.withdraw.rejected')" value="rejected" />
            <el-option :label="t('status.withdraw.paid')" value="paid" />
            <el-option :label="t('status.withdraw.failed')" value="failed" />
          </el-select>
          <el-button :loading="loading" @click="loadList">{{ t("common.refresh") }}</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="t('page.withdraw.seller')" min-width="220">
        <template #default="{ row }">
          <div class="seller-name">{{ row.shop_name || row.user_name || "-" }}</div>
          <div class="muted">{{ t("page.withdraw.sellerLine", { id: row.seller_id }) }} / {{ row.nick_name || "-" }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.amount')" width="110">
        <template #default="{ row }">¥{{ money(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('page.withdraw.payeeInfo')" min-width="220">
        <template #default="{ row }">
          <div>{{ row.payee_name }}</div>
          <div class="muted">{{ row.payee_account }} / {{ row.payee_channel }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.withdraw.auditResult')" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">{{ row.audit_reason || "-" }}</template>
      </el-table-column>
      <el-table-column :label="t('common.time')" min-width="220">
        <template #default="{ row }">
          <div>{{ t("page.withdraw.requestedAt", { time: formatTime(row.created_at) }) }}</div>
          <div class="muted">{{ t("page.withdraw.auditedAt", { time: formatTime(row.audited_at) }) }}</div>
          <div class="muted">{{ t("page.withdraw.paidAt", { time: formatTime(row.paid_at) }) }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="280" fixed="right">
        <template #default="{ row }">
          <div class="withdraw-actions">
            <el-button size="small" @click="openDetail(row)">{{ t("common.flow") }}</el-button>
            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="success"
              @click="approve(row)"
            >
              {{ t("page.withdraw.approve") }}
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="warning"
              @click="openReject(row)"
            >
              {{ t("page.withdraw.reject") }}
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              size="small"
              type="primary"
              @click="markPaid(row)"
            >
              {{ t("page.withdraw.markPaid") }}
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              size="small"
              type="danger"
              @click="openFailed(row)"
            >
              {{ t("page.withdraw.markFailed") }}
            </el-button>
          </div>
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

    <el-dialog v-model="reasonVisible" :title="reasonDialogTitle" width="520px">
      <el-form label-width="88px">
        <el-form-item :label="t('page.withdraw.withdrawOrder')">
          <span>#{{ currentRow?.id }} / ¥{{ money(currentRow?.amount || 0) }}</span>
        </el-form-item>
        <el-form-item :label="t('page.withdraw.reason')" required>
          <el-input
            v-model="reason"
            type="textarea"
            :rows="4"
            maxlength="255"
            show-word-limit
            :placeholder="reasonPlaceholder"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reasonVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="saving" @click="submitReason">
          {{ t("common.confirm") }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" :title="t('page.withdraw.flowTitle')" width="760px">
      <el-table :data="flows" size="small">
        <el-table-column prop="flow_no" :label="t('page.withdraw.flowNo')" min-width="190" />
        <el-table-column prop="flow_type" :label="t('common.type')" width="180" />
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
import { computed, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  auditAdminSellerWithdraw,
  getAdminSellerWithdrawDetail,
  getAdminSellerWithdrawList,
  markAdminSellerWithdrawPaid,
} from "@/api/withdraw";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const list = ref<any[]>([]);
const flows = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const saving = ref(false);
const detailVisible = ref(false);
const reasonVisible = ref(false);
const currentRow = ref<any | null>(null);
const reasonMode = ref<"reject" | "failed">("reject");
const reason = ref("");
const sellerId = ref<number | undefined>();
const statusFilter = ref<string | undefined>();

function money(value: number) {
  return Number(value || 0).toFixed(2);
}

function formatTime(value: number) {
  if (!value) return "-";
  return new Date(value * 1000).toLocaleString();
}

function statusTag(status: string) {
  return (
    {
      pending: "warning",
      approved: "primary",
      rejected: "danger",
      paid: "success",
      failed: "info",
    } as any
  )[status] ?? "info";
}

function statusText(status: string) {
  return (
    {
      pending: t("status.withdraw.pending"),
      approved: t("status.withdraw.approved"),
      rejected: t("status.withdraw.rejected"),
      paid: t("status.withdraw.paid"),
      failed: t("status.withdraw.failed"),
    } as any
  )[status] ?? t("common.unknown");
}

function reload() {
  page.value = 1;
  loadList();
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAdminSellerWithdrawList({
      page_num: page.value,
      page_size: pageSize,
      ...(sellerId.value ? { seller_id: sellerId.value } : {}),
      ...(statusFilter.value ? { status: statusFilter.value } : {}),
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function approve(row: any) {
  await ElMessageBox.confirm(t("page.withdraw.approveConfirm"), t("common.notice"), { type: "warning" });
  saving.value = true;
  try {
    await auditAdminSellerWithdraw({ id: row.id, status: "approved" });
    ElMessage.success(t("page.withdraw.approveSuccess"));
    requestAdminPendingCountsRefresh();
    loadList();
  } finally {
    saving.value = false;
  }
}

function openReject(row: any) {
  currentRow.value = row;
  reasonMode.value = "reject";
  reason.value = "";
  reasonVisible.value = true;
}

function openFailed(row: any) {
  currentRow.value = row;
  reasonMode.value = "failed";
  reason.value = "";
  reasonVisible.value = true;
}

const reasonDialogTitle = computed(() =>
  reasonMode.value === "reject" ? t("page.withdraw.rejectTitle") : t("page.withdraw.failedTitle"),
);
const reasonPlaceholder = computed(() =>
  reasonMode.value === "reject" ? t("page.withdraw.rejectPlaceholder") : t("page.withdraw.failedPlaceholder"),
);

async function submitReason() {
  if (!reason.value.trim()) {
    return ElMessage.warning(t("page.withdraw.reasonRequired"));
  }
  if (!currentRow.value) return;
  saving.value = true;
  try {
    if (reasonMode.value === "reject") {
      await auditAdminSellerWithdraw({
        id: currentRow.value.id,
        status: "rejected",
        reason: reason.value.trim(),
      });
      ElMessage.success(t("page.withdraw.rejectSuccess"));
      requestAdminPendingCountsRefresh();
    } else {
      await markAdminSellerWithdrawPaid({
        id: currentRow.value.id,
        status: "failed",
        reason: reason.value.trim(),
      });
      ElMessage.success(t("page.withdraw.failedSuccess"));
    }
    reasonVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function markPaid(row: any) {
  await ElMessageBox.confirm(t("page.withdraw.paidConfirm"), t("common.notice"), { type: "warning" });
  saving.value = true;
  try {
    await markAdminSellerWithdrawPaid({ id: row.id, status: "paid" });
    ElMessage.success(t("page.withdraw.paidSuccess"));
    requestAdminPendingCountsRefresh();
    loadList();
  } finally {
    saving.value = false;
  }
}

async function openDetail(row: any) {
  const res: any = await getAdminSellerWithdrawDetail({ id: row.id });
  flows.value = res.data?.item ?? [];
  detailVisible.value = true;
}

onMounted(loadList);
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
.seller-name {
  color: #303133;
  font-weight: 600;
}
.muted {
  color: #909399;
  font-size: 12px;
}
.withdraw-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.withdraw-actions :deep(.el-button + .el-button) {
  margin-left: 0;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
