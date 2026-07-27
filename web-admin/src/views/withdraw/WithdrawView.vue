<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>商家提现</span>
        <div class="filters">
          <el-input-number
            v-model="sellerId"
            :min="0"
            :step="1"
            placeholder="商家ID"
            controls-position="right"
            @change="reload"
          />
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
          <el-button :loading="loading" @click="loadList">刷新</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="商家" min-width="220">
        <template #default="{ row }">
          <div class="seller-name">{{ row.shop_name || row.user_name || "-" }}</div>
          <div class="muted">商家ID: {{ row.seller_id }} · {{ row.nick_name || "-" }}</div>
        </template>
      </el-table-column>
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
      <el-table-column label="审核结果" min-width="150" show-overflow-tooltip>
        <template #default="{ row }">{{ row.audit_reason || "-" }}</template>
      </el-table-column>
      <el-table-column label="时间" min-width="220">
        <template #default="{ row }">
          <div>申请：{{ formatTime(row.created_at) }}</div>
          <div class="muted">审核：{{ formatTime(row.audited_at) }}</div>
          <div class="muted">打款：{{ formatTime(row.paid_at) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">流水</el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="success"
            @click="approve(row)"
          >
            通过
          </el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="warning"
            @click="openReject(row)"
          >
            拒绝
          </el-button>
          <el-button
            v-if="row.status === 'approved'"
            size="small"
            type="primary"
            @click="markPaid(row)"
          >
            已打款
          </el-button>
          <el-button
            v-if="row.status === 'approved'"
            size="small"
            type="danger"
            @click="openFailed(row)"
          >
            打款失败
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

    <el-dialog v-model="reasonVisible" :title="reasonDialogTitle" width="520px">
      <el-form label-width="88px">
        <el-form-item label="提现单">
          <span>#{{ currentRow?.id }} · ¥{{ money(currentRow?.amount || 0) }}</span>
        </el-form-item>
        <el-form-item label="原因" required>
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
        <el-button @click="reasonVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="submitReason">
          确认
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="detailVisible" title="提现流水" width="760px">
      <el-table :data="flows" size="small">
        <el-table-column prop="flow_no" label="流水号" min-width="190" />
        <el-table-column prop="flow_type" label="类型" width="180" />
        <el-table-column prop="direction" label="方向" width="80" />
        <el-table-column label="金额" width="100">
          <template #default="{ row }">¥{{ money(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="140" />
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
  await ElMessageBox.confirm(`确认通过提现吗？`, "提示", { type: "warning" });
  saving.value = true;
  try {
    await auditAdminSellerWithdraw({ id: row.id, status: "approved" });
    ElMessage.success("已通过");
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
  reasonMode.value === "reject" ? "拒绝提现" : "标记打款失败",
);
const reasonPlaceholder = computed(() =>
  reasonMode.value === "reject" ? "请输入拒绝原因" : "请输入失败原因",
);

async function submitReason() {
  if (!reason.value.trim()) {
    return ElMessage.warning("请填写原因");
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
      ElMessage.success("已拒绝");
    } else {
      await markAdminSellerWithdrawPaid({
        id: currentRow.value.id,
        status: "failed",
        reason: reason.value.trim(),
      });
      ElMessage.success("已标记失败");
    }
    reasonVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function markPaid(row: any) {
  await ElMessageBox.confirm(`确认已打款提现吗？`, "提示", { type: "warning" });
  saving.value = true;
  try {
    await markAdminSellerWithdrawPaid({ id: row.id, status: "paid" });
    ElMessage.success("已标记打款");
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
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
