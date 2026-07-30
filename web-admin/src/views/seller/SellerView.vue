<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.seller.title") }}</span>
        <div class="filters">
          <el-radio-group v-model="statusFilter" @change="reload">
            <el-radio-button :value="undefined">{{ t("common.all") }}</el-radio-button>
            <el-radio-button :value="0">{{ t("status.seller.pending") }}</el-radio-button>
            <el-radio-button :value="1">{{ t("status.seller.approved") }}</el-radio-button>
            <el-radio-button :value="2">{{ t("status.seller.rejected") }}</el-radio-button>
            <el-radio-button :value="3">{{ t("status.seller.banned") }}</el-radio-button>
          </el-radio-group>
          <el-button :loading="loading" @click="loadList">{{ t("common.refresh") }}</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="t('page.seller.shop')" min-width="220">
        <template #default="{ row }">
          <div class="shop-name">{{ row.shop_name }}</div>
          <div class="muted">{{ row.description || t("page.seller.noDescription") }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.seller.applicant')" min-width="180">
        <template #default="{ row }">
          <div>{{ row.nick_name || row.user_name }}</div>
          <div class="muted">{{ t("page.seller.userId", { id: row.user_id }) }} / {{ row.email || "-" }}</div>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.status')" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reject_reason" :label="t('page.seller.rejectReason')" min-width="160" show-overflow-tooltip />
      <el-table-column :label="t('common.actions')" width="220" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status !== 1"
            size="small"
            type="success"
            @click="audit(row, 1)"
          >
            {{ t("page.seller.approve") }}
          </el-button>
          <el-button
            v-if="row.status !== 2"
            size="small"
            type="warning"
            @click="openReject(row)"
          >
            {{ t("page.seller.reject") }}
          </el-button>
          <el-button
            v-if="row.status !== 3"
            size="small"
            type="danger"
            @click="audit(row, 3)"
          >
            {{ t("page.seller.ban") }}
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

    <el-dialog v-model="rejectVisible" :title="t('page.seller.rejectDialogTitle')" width="520px">
      <el-form label-width="86px">
        <el-form-item :label="t('page.seller.shop')">
          {{ rejectRow?.shop_name }}
        </el-form-item>
        <el-form-item :label="t('page.seller.rejectReason')" required>
          <el-input
            v-model="rejectReason"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
            :placeholder="t('page.seller.rejectPlaceholder')"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="auditing" @click="submitReject">
          {{ t("page.seller.confirmReject") }}
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { auditAdminSeller, getAdminSellerList } from "@/api/seller";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const auditing = ref(false);
const statusFilter = ref<number | undefined>(0);
const rejectVisible = ref(false);
const rejectRow = ref<any | null>(null);
const rejectReason = ref("");

function statusTag(status: number) {
  return ({ 0: "warning", 1: "success", 2: "danger", 3: "info" } as any)[
    status
  ] ?? "info";
}

function statusText(status: number) {
  return (
    {
      0: t("status.seller.pending"),
      1: t("status.seller.approved"),
      2: t("status.seller.rejected"),
      3: t("status.seller.banned"),
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
    const res: any = await getAdminSellerList({
      page_num: page.value,
      page_size: pageSize,
      ...(statusFilter.value !== undefined ? { status: statusFilter.value } : {}),
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function audit(row: any, status: number) {
  const label = status === 1 ? t("page.seller.approve") : t("page.seller.ban");
  await ElMessageBox.confirm(t("page.seller.confirmAudit", { action: label, name: row.shop_name }), t("common.notice"), {
    type: "warning",
  });
  auditing.value = true;
  try {
    await auditAdminSeller({ id: row.id, status });
    ElMessage.success(t("page.seller.auditSuccess", { action: label }));
    requestAdminPendingCountsRefresh();
    loadList();
  } finally {
    auditing.value = false;
  }
}

function openReject(row: any) {
  rejectRow.value = row;
  rejectReason.value = "";
  rejectVisible.value = true;
}

async function submitReject() {
  if (!rejectReason.value.trim()) {
    return ElMessage.warning(t("page.seller.rejectReasonRequired"));
  }
  auditing.value = true;
  try {
    await auditAdminSeller({
      id: rejectRow.value.id,
      status: 2,
      reject_reason: rejectReason.value.trim(),
    });
    ElMessage.success(t("status.seller.rejected"));
    requestAdminPendingCountsRefresh();
    rejectVisible.value = false;
    loadList();
  } finally {
    auditing.value = false;
  }
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
.shop-name {
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
