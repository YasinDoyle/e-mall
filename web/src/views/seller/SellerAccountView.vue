<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>资金账户</span>
        <div class="header-actions">
          <el-select
            v-model="statusFilter"
            clearable
            placeholder="提现状态"
            style="width: 140px"
            @change="reload"
          >
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
            <el-option label="已打款" value="paid" />
            <el-option label="打款失败" value="failed" />
          </el-select>
          <el-button :loading="loading" @click="reloadAll">刷新</el-button>
          <el-button
            type="primary"
            :disabled="!sellerStore.isApproved"
            @click="openApply"
          >
            申请提现
          </el-button>
        </div>
      </div>
    </template>

    <el-alert
      v-if="!sellerStore.isApproved"
      title="商家审核通过后才可以查看资金账户并申请提现"
      type="warning"
      :closable="false"
      show-icon
      class="notice"
    />

    <div class="summary-grid">
      <div class="summary-item">
        <span class="summary-label">可提现余额</span>
        <b>¥{{ money(summary.available_balance) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">冻结中</span>
        <b>¥{{ money(summary.frozen_balance) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">累计收入</span>
        <b>¥{{ money(summary.total_income) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">累计提现</span>
        <b>¥{{ money(summary.total_withdrawn) }}</b>
      </div>
    </div>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="金额" width="110">
        <template #default="{ row }">¥{{ money(row.amount) }}</template>
      </el-table-column>
      <el-table-column label="收款信息" min-width="220">
        <template #default="{ row }">
          <div>{{ row.payee_name }}</div>
          <div class="muted">{{ row.payee_account }} · {{ row.payee_channel }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ row.status_text }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="审核结果" min-width="180" show-overflow-tooltip>
        <template #default="{ row }">
          <span>{{ row.audit_reason || "-" }}</span>
        </template>
      </el-table-column>
      <el-table-column label="时间" min-width="190">
        <template #default="{ row }">
          <div>申请：{{ formatTime(row.created_at) }}</div>
          <div class="muted">审核：{{ formatTime(row.audited_at) }}</div>
          <div class="muted">打款：{{ formatTime(row.paid_at) }}</div>
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

    <el-dialog v-model="applyVisible" title="申请提现" width="520px">
      <el-form label-width="88px">
        <el-form-item label="提现金额">
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
              全部
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="收款人">
          <el-input v-model="applyForm.payee_name" maxlength="64" placeholder="请输入收款人姓名" />
        </el-form-item>
        <el-form-item label="收款账号">
          <el-input
            v-model="applyForm.payee_account"
            maxlength="128"
            placeholder="请输入银行账号、支付宝或微信收款号"
          />
        </el-form-item>
        <el-form-item label="收款方式">
          <el-select
            v-model="applyForm.payee_channel"
            placeholder="请选择收款方式"
            style="width: 220px"
          >
            <el-option label="银行卡" value="bank" />
            <el-option label="支付宝" value="alipay" />
            <el-option label="微信" value="wechat" />
            <el-option label="人工处理" value="manual" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyVisible = false">取消</el-button>
        <el-button type="primary" :loading="applying" @click="submitApply">
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import {
  applySellerWithdraw,
  getSellerAccountSummary,
  getSellerWithdrawList,
} from "@/api/seller";
import { useSellerStore } from "@/stores/seller";

const sellerStore = useSellerStore();
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
    return ElMessage.warning("请输入提现金额");
  }
  if (!applyForm.payee_name.trim() || !applyForm.payee_account.trim()) {
    return ElMessage.warning("请填写收款信息");
  }
  applying.value = true;
  try {
    await applySellerWithdraw({
      amount: applyForm.amount,
      payee_name: applyForm.payee_name.trim(),
      payee_account: applyForm.payee_account.trim(),
      payee_channel: applyForm.payee_channel.trim() || "bank",
    });
    ElMessage.success("提现申请已提交");
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
