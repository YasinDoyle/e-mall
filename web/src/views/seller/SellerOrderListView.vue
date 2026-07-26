<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>订单管理</span>
        <el-button :loading="loading" @click="reloadAll">刷新</el-button>
      </div>
    </template>

    <div class="summary-grid">
      <div class="summary-item account">
        <span class="summary-label">商家账户余额</span>
        <b>¥{{ money(summary.available_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">待结算</span>
        <b>¥{{ money(summary.pending_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">可结算</span>
        <b>¥{{ money(summary.generated_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">已打款</span>
        <b>¥{{ money(summary.paid_amount) }}</b>
      </div>
      <div class="summary-item">
        <span class="summary-label">已退款</span>
        <b>¥{{ money(summary.refunded_amount) }}</b>
      </div>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="0" />
      <el-tab-pane label="待发货" name="2" />
      <el-tab-pane label="已发货" name="3" />
      <el-tab-pane label="已完成" name="4" />
      <el-tab-pane label="退款中" name="5" />
      <el-tab-pane label="已退款" name="6" />
    </el-tabs>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="order_num" label="订单号" min-width="150" />
      <el-table-column label="商品" min-width="240">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div class="name">{{ row.name || "商品" }}</div>
              <div class="muted">买家ID：{{ row.user_id }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="交易信息" min-width="210">
        <template #default="{ row }">
          <div>订单 ¥{{ totalAmount(row) }}</div>
          <div class="muted">
            佣金 ¥{{ money(row.commission_amount) }} / 收入 ¥{{ money(row.settlement_amount) }}
          </div>
          <el-tag size="small" :type="settlementTag(row.settlement_status)">
            {{ settlementText(row.settlement_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="num" label="数量" width="80" />
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.type)">
            {{ statusText(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="收货信息" min-width="260">
        <template #default="{ row }">
          <div>{{ row.address_name }} {{ row.address_phone }}</div>
          <div class="muted">{{ row.address }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="tracking_no" label="物流单号" min-width="150">
        <template #default="{ row }">{{ row.tracking_no || "-" }}</template>
      </el-table-column>
      <el-table-column label="操作" width="120" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.type === 2"
            size="small"
            type="primary"
            @click="openShip(row)"
          >
            发货
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

    <el-dialog v-model="shipDialogVisible" title="订单发货" width="420px">
      <el-form label-width="80px">
        <el-form-item label="订单号">
          <span>{{ currentOrder?.order_num }}</span>
        </el-form-item>
        <el-form-item label="物流单号">
          <el-input v-model="trackingNo" maxlength="64" placeholder="请输入物流单号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="shipping" @click="handleShip">
          确认发货
        </el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import {
  getSellerOrderList,
  getSellerSettlementSummary,
  shipOrder,
} from "@/api/order";
import { useSellerStore } from "@/stores/seller";

const sellerStore = useSellerStore();
const summary = ref({
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

const statusText = (type: number) =>
  ({
    1: "待支付",
    2: "待发货",
    3: "已发货",
    4: "已完成",
    5: "退款中",
    6: "已退款",
  })[type] ?? "未知";

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
      pending: "待结算",
      generated: "可结算",
      paid: "已打款",
      refunded: "已退款",
    } as any
  )[status || ""] ?? "未生成";
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
  const res: any = await getSellerSettlementSummary();
  summary.value = {
    available_amount: Number(res.data?.available_amount || 0),
    pending_amount: Number(res.data?.pending_amount || 0),
    generated_amount: Number(res.data?.generated_amount || 0),
    paid_amount: Number(res.data?.paid_amount || 0),
    refunded_amount: Number(res.data?.refunded_amount || 0),
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
  if (!trackingNo.value.trim()) return ElMessage.warning("请输入物流单号");
  shipping.value = true;
  try {
    await shipOrder({
      order_id: currentOrder.value.id,
      tracking_no: trackingNo.value.trim(),
    });
    ElMessage.success("已发货");
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
