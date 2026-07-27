<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>结算管理</span>
        <div class="filters">
          <el-input-number
            v-model="sellerID"
            :min="0"
            :step="1"
            placeholder="商家ID"
            controls-position="right"
          />
          <el-select v-model="statusFilter" clearable placeholder="结算状态" style="width: 140px" @change="reload">
            <el-option label="待结算" value="pending" />
            <el-option label="已生成" value="generated" />
            <el-option label="已打款" value="paid" />
            <el-option label="已退款" value="refunded" />
          </el-select>
          <el-button :loading="loading" @click="reloadAll">刷新</el-button>
          <el-button :loading="backfilling" @click="backfillAccount">回填账户</el-button>
          <el-button type="primary" @click="generate">按商家批量生成</el-button>
        </div>
      </div>
    </template>

    <div class="summary-grid">
      <div class="summary-item">
        <span class="summary-label">平台收益</span>
        <b>¥{{ money(overview.platform_revenue) }}</b>
      </div>
    </div>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="seller_id" label="商家ID" width="90" />
      <el-table-column prop="order_num" label="订单号" min-width="170" />
      <el-table-column label="订单金额" width="110">
        <template #default="{ row }">¥{{ money(row.gross_amount) }}</template>
      </el-table-column>
      <el-table-column label="佣金" width="130">
        <template #default="{ row }">
          ¥{{ money(row.commission_amount) }}
          <span class="muted">({{ percent(row.commission_rate) }})</span>
        </template>
      </el-table-column>
      <el-table-column label="应结算" width="120">
        <template #default="{ row }">¥{{ money(row.settlement_amount) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">流水</el-button>
          <el-button
            v-if="row.status === 'pending'"
            size="small"
            type="primary"
            @click="generateOne(row)"
          >
            生成结算
          </el-button>
          <el-button
            v-if="row.status === 'generated'"
            size="small"
            type="primary"
            @click="markPaid(row)"
          >
            标记打款
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

    <el-dialog v-model="detailVisible" title="结算流水" width="760px">
      <el-table :data="flows" size="small">
        <el-table-column prop="flow_no" label="流水号" min-width="190" />
        <el-table-column prop="flow_type" label="类型" width="150" />
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
      pending: "待结算",
      generated: "已生成",
      paid: "已打款",
      refunded: "已退款",
    } as any
  )[status] ?? "未知";
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
    return ElMessage.warning("请输入商家ID");
  }
  const res: any = await generateAdminSettlement({ seller_id: sellerID.value });
  ElMessage.success(`已生成 ${res.data?.count ?? 0} 条结算单`);
  reload();
}

async function backfillAccount() {
  await ElMessageBox.confirm("确认回填所有已打款但未入账的商家账户吗？", "提示", {
    type: "warning",
  });
  backfilling.value = true;
  try {
    const res: any = await backfillAdminSellerAccount();
    const data = res.data || {};
    ElMessage.success(
      `回填完成：${data.settlement_count ?? 0} 笔结算，${data.seller_count ?? 0} 个商家，¥${money(
        data.amount ?? 0,
      )}`,
    );
    loadList();
  } finally {
    backfilling.value = false;
  }
}

async function generateOne(row: any) {
  await ElMessageBox.confirm(
    `确认生成订单 ${row.order_num} 的单笔结算？`,
    "提示",
    { type: "warning" },
  );
  await generateOneAdminSettlement({ id: row.id });
  ElMessage.success("已生成单笔结算");
  loadList();
}

async function markPaid(row: any) {
  await ElMessageBox.confirm(
    `确认将结算单 ${row.id} 标记为已打款？`,
    "提示",
    { type: "warning" },
  );
  await markAdminSettlementPaid({ id: row.id });
  ElMessage.success("已标记打款");
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
