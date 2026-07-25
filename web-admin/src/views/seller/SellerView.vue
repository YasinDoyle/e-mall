<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>商家审核</span>
        <div class="filters">
          <el-radio-group v-model="statusFilter" @change="reload">
            <el-radio-button :value="undefined">全部</el-radio-button>
            <el-radio-button :value="0">待审核</el-radio-button>
            <el-radio-button :value="1">已通过</el-radio-button>
            <el-radio-button :value="2">已拒绝</el-radio-button>
            <el-radio-button :value="3">已封禁</el-radio-button>
          </el-radio-group>
          <el-button :loading="loading" @click="loadList">刷新</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="店铺" min-width="220">
        <template #default="{ row }">
          <div class="shop-name">{{ row.shop_name }}</div>
          <div class="muted">{{ row.description || "暂无简介" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="申请人" min-width="180">
        <template #default="{ row }">
          <div>{{ row.nick_name || row.user_name }}</div>
          <div class="muted">用户ID: {{ row.user_id }} · {{ row.email || "-" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ row.status_text }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reject_reason" label="拒绝原因" min-width="160" show-overflow-tooltip />
      <el-table-column label="操作" width="220" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status !== 1"
            size="small"
            type="success"
            @click="audit(row, 1)"
          >
            通过
          </el-button>
          <el-button
            v-if="row.status !== 2"
            size="small"
            type="warning"
            @click="openReject(row)"
          >
            拒绝
          </el-button>
          <el-button
            v-if="row.status !== 3"
            size="small"
            type="danger"
            @click="audit(row, 3)"
          >
            封禁
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

    <el-dialog v-model="rejectVisible" title="拒绝商家申请" width="520px">
      <el-form label-width="86px">
        <el-form-item label="店铺">
          {{ rejectRow?.shop_name }}
        </el-form-item>
        <el-form-item label="拒绝原因" required>
          <el-input
            v-model="rejectReason"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
            placeholder="请填写明确原因，便于用户修改后重新提交"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectVisible = false">取消</el-button>
        <el-button type="primary" :loading="auditing" @click="submitReject">
          确认拒绝
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { auditAdminSeller, getAdminSellerList } from "@/api/seller";

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
  const label = status === 1 ? "通过" : "封禁";
  await ElMessageBox.confirm(`确认${label}店铺「${row.shop_name}」？`, "提示", {
    type: "warning",
  });
  auditing.value = true;
  try {
    await auditAdminSeller({ id: row.id, status });
    ElMessage.success(`已${label}`);
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
    return ElMessage.warning("请填写拒绝原因");
  }
  auditing.value = true;
  try {
    await auditAdminSeller({
      id: rejectRow.value.id,
      status: 2,
      reject_reason: rejectReason.value.trim(),
    });
    ElMessage.success("已拒绝");
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
