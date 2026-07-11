<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>订单管理</span>
        <div class="filters">
          <el-select v-model="typeFilter" placeholder="订单状态" clearable style="width: 150px" @change="reload">
            <el-option v-for="item in orderTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="refundFilter" placeholder="退款状态" clearable style="width: 140px" @change="reload">
            <el-option label="无退款" :value="0" />
            <el-option label="申请中" :value="1" />
            <el-option label="已退款" :value="2" />
          </el-select>
          <el-button :loading="loading" @click="loadList">刷新</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_num" label="订单号" min-width="170" />
      <el-table-column label="商品" min-width="220">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div>{{ row.name || "-" }}</div>
              <div class="muted">商品ID: {{ row.product_id }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="user_id" label="买家ID" width="90" />
      <el-table-column prop="boss_id" label="卖家ID" width="90" />
      <el-table-column label="金额" width="110">
        <template #default="{ row }">¥{{ totalAmount(row).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="订单状态" width="130">
        <template #default="{ row }">
          <el-tag :type="orderTypeTag(row.type)">{{ orderTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="退款状态" width="110">
        <template #default="{ row }">
          <el-tag :type="refundTag(row.refund_status)">{{ refundText(row.refund_status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="refund_reason" label="退款原因" min-width="160" show-overflow-tooltip />
      <el-table-column label="操作" width="130" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.refund_status === 1"
            size="small"
            type="primary"
            @click="approveRefund(row)"
          >
            退款审批
          </el-button>
          <span v-else class="muted">-</span>
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
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { approveOrderRefund, getAdminOrderList } from "@/api";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const typeFilter = ref<number | undefined>();
const refundFilter = ref<number | undefined>();

const orderTypes = [
  { value: 1, label: "未支付" },
  { value: 2, label: "待发货" },
  { value: 3, label: "待收货" },
  { value: 4, label: "已完成" },
  { value: 5, label: "退款申请中" },
  { value: 6, label: "已退款" },
];

function orderTypeText(type: number) {
  return orderTypes.find((item) => item.value === type)?.label ?? "未知";
}

function orderTypeTag(type: number) {
  return ({ 1: "info", 2: "warning", 3: "primary", 4: "success", 5: "danger", 6: "info" } as any)[type] ?? "info";
}

function refundText(status: number) {
  return ({ 0: "无退款", 1: "申请中", 2: "已退款" } as any)[status] ?? "未知";
}

function refundTag(status: number) {
  return ({ 0: "info", 1: "warning", 2: "success" } as any)[status] ?? "info";
}

function totalAmount(row: any) {
  return Number(row.money || 0) * Number(row.num || 1);
}

function reload() {
  page.value = 1;
  loadList();
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAdminOrderList({
      page_num: page.value,
      page_size: pageSize,
      ...(typeFilter.value !== undefined ? { type: typeFilter.value } : {}),
      ...(refundFilter.value !== undefined ? { refund_status: refundFilter.value } : {}),
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function approveRefund(row: any) {
  const { value } = await ElMessageBox.prompt(
    `确认给订单 ${row.order_num} 退款 ¥${totalAmount(row).toFixed(2)}，请输入资金密钥`,
    "退款审批",
    {
      inputType: "password",
      inputPattern: /.+/,
      inputErrorMessage: "请输入资金密钥",
      confirmButtonText: "审批退款",
      cancelButtonText: "取消",
    },
  );
  await approveOrderRefund({ order_id: row.id, key: value });
  ElMessage.success("退款审批完成");
  loadList();
}

onMounted(loadList);
</script>

<style scoped>
.card-header,
.filters,
.product-cell {
  display: flex;
  align-items: center;
}
.card-header {
  justify-content: space-between;
  gap: 12px;
}
.filters {
  gap: 8px;
}
.product-cell {
  gap: 10px;
}
.product-img {
  width: 46px;
  height: 46px;
  object-fit: cover;
  border-radius: 4px;
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
