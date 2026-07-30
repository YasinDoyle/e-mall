<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>{{ t("sellerCenter.account.title") }}</span>
        <div class="header-actions">
          <el-select
            v-model="statusFilter"
            clearable
            :placeholder="t('sellerCenter.account.withdrawStatus')"
            style="width: 140px"
            @change="reload"
          >
            <el-option :label="t('status.withdraw.pending')" value="pending" />
            <el-option :label="t('status.withdraw.approved')" value="approved" />
            <el-option :label="t('status.withdraw.rejected')" value="rejected" />
            <el-option :label="t('status.withdraw.paid')" value="paid" />
            <el-option :label="t('status.withdraw.failed')" value="failed" />
          </el-select>
          <el-button :loading="loading" @click="reloadAll">{{ t("common.refresh") }}</el-button>
          <el-button
            type="primary"
            :disabled="!sellerStore.isApproved"
            @click="openApply"
          >
            {{ t("sellerCenter.account.applyWithdraw") }}
          </el-button>
        </div>
      </div>
    </template>

    <el-alert
      v-if="!sellerStore.isApproved"
      :title="t('sellerCenter.account.approvedOnly')"
      type="warning"
      :closable="false"
      show-icon
      class="notice"
    />

    <div class="summary-grid">
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.account.availableBalance") }}</span>
        <b>¥{{ money(summary.available_balance) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.account.frozenBalance") }}</span>
        <b>¥{{ money(summary.frozen_balance) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.account.totalIncome") }}</span>
        <b>¥{{ money(summary.total_income) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">{{ t("sellerCenter.account.totalWithdrawn") }}</span>
        <b>¥{{ money(summary.total_withdrawn) }}</b>
      </div>
    </div>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="t('sellerCenter.account.amount')" width="110">
        <template #default="{ row }">¥{{ money(row.amount) }}</template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.account.payeeInfo')" min-width="220">
        <template #default="{ row }">
          <div>{{ row.payee_name }}</div>
          <div class="muted">{{ row.payee_account }} · {{ row.payee_channel }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.account.status')" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ withdrawStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.account.auditResult')" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <span>{{ row.audit_reason || "-" }}</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.account.time')" min-width="190">
        <template #default="{ row }">
          <div>{{ t("sellerCenter.account.requestedAt", { time: formatTime(row.created_at) }) }}</div>
          <div class="muted">
            {{ t("sellerCenter.account.auditedAt", { time: formatTime(row.audited_at) }) }}
          </div>
          <div class="muted">
            {{ t("sellerCenter.account.paidAt", { time: formatTime(row.paid_at) }) }}
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

    <el-dialog v-model="applyVisible" :title="t('sellerCenter.account.applyWithdraw')" width="520px">
      <el-form label-width="88px">
        <el-form-item :label="t('sellerCenter.account.withdrawAmount')">
          <div class="amount-row">
            <el-input-number
              v-model="applyForm.amount"
              :min="0"
              :precision="2"
              :step="100"
              controls-position="right"
              style="width: 220px"
            />
            <el-button :disabled="!summary.available_balance" @click="useAllAmount">
              {{ t("sellerCenter.account.useAll") }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('sellerCenter.account.payeeName')">
          <el-input
            v-model="applyForm.payee_name"
            maxlength="64"
            :placeholder="t('sellerCenter.account.payeeNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('sellerCenter.account.payeeAccount')">
          <el-input
            v-model="applyForm.payee_account"
            maxlength="128"
            :placeholder="t('sellerCenter.account.payeeAccountPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('sellerCenter.account.payeeChannel')">
          <el-select
            v-model="applyForm.payee_channel"
            :placeholder="t('sellerCenter.account.payeeChannelPlaceholder')"
            style="width: 220px"
          >
            <el-option :label="t('sellerCenter.account.channelBank')" value="bank" />
            <el-option :label="t('sellerCenter.account.channelAlipay')" value="alipay" />
            <el-option :label="t('sellerCenter.account.channelWechat')" value="wechat" />
            <el-option :label="t('sellerCenter.account.channelManual')" value="manual" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="applying" @click="submitApply">
          {{ t("sellerCenter.account.submitApply") }}
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import {
  applySellerWithdraw,
  getSellerAccountSummary,
  getSellerWithdrawList,
} from "@/api/seller";
import { useSellerStore } from "@/stores/seller";

const sellerStore = useSellerStore();
const { t } = useI18n();
const loading = ref(false);
const applying = ref(false);
const applyVisible = ref(false);
const statusFilter = ref<string | undefined>();
const page = ref(1);
const pageSize = 10;
const total = ref(0);
const summary = ref({
  available_balance: 0,
  frozen_balance: 0,
  total_income: 0,
  total_withdrawn: 0,
});
const list = ref<any[]>([]);
const applyForm = reactive({
  amount: 0,
  payee_name: "",
  payee_account: "",
  payee_channel: "bank",
});

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

function withdrawStatusLabel(status: string) {
  return (
    {
      pending: t("status.withdraw.pending"),
      approved: t("status.withdraw.approved"),
      rejected: t("status.withdraw.rejected"),
      paid: t("status.withdraw.paid"),
      failed: t("status.withdraw.failed"),
    } as Record<string, string>
  )[status] ?? t("common.unknown");
}

async function loadSummary() {
  try {
    const res: any = await getSellerAccountSummary();
    summary.value = {
      available_balance: Number(res.data?.available_balance || 0),
      frozen_balance: Number(res.data?.frozen_balance || 0),
      total_income: Number(res.data?.total_income || 0),
      total_withdrawn: Number(res.data?.total_withdrawn || 0),
    };
  } catch {
    summary.value = {
      available_balance: 0,
      frozen_balance: 0,
      total_income: 0,
      total_withdrawn: 0,
    };
  }
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getSellerWithdrawList({
      page_num: page.value,
      page_size: pageSize,
      ...(statusFilter.value ? { status: statusFilter.value } : {}),
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function reload() {
  page.value = 1;
  loadList();
}

async function reloadAll() {
  await Promise.all([loadSummary(), loadList()]);
}

function openApply() {
  applyForm.amount = 0;
  applyForm.payee_name = "";
  applyForm.payee_account = "";
  applyForm.payee_channel = "bank";
  applyVisible.value = true;
}

function useAllAmount() {
  applyForm.amount = Number(summary.value.available_balance || 0);
}

async function submitApply() {
  if (!applyForm.amount || applyForm.amount <= 0) {
    return ElMessage.warning(t("sellerCenter.account.amountRequired"));
  }
  if (!applyForm.payee_name.trim() || !applyForm.payee_account.trim()) {
    return ElMessage.warning(t("sellerCenter.account.payeeRequired"));
  }
  applying.value = true;
  try {
    await applySellerWithdraw({
      amount: applyForm.amount,
      payee_name: applyForm.payee_name.trim(),
      payee_account: applyForm.payee_account.trim(),
      payee_channel: applyForm.payee_channel.trim() || "bank",
    });
    ElMessage.success(t("sellerCenter.account.submitSuccess"));
    applyVisible.value = false;
    await reloadAll();
  } finally {
    applying.value = false;
  }
}

onMounted(async () => {
  await sellerStore.loadProfile({ silentError: true });
  await reloadAll();
});
</script>

<style scoped>
.header,
.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}
.header {
  justify-content: space-between;
}
.amount-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 16px;
}
.summary-item {
  padding: 12px 14px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
}
.summary-label {
  display: block;
  color: #909399;
  font-size: 12px;
  margin-bottom: 6px;
}
.notice {
  margin-bottom: 12px;
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
