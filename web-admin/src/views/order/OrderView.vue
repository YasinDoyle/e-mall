<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.order.title") }}</span>
        <div class="filters">
          <el-select v-model="typeFilter" :placeholder="t('page.order.orderStatus')" clearable style="width: 150px" @change="reload">
            <el-option v-for="item in orderTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="refundFilter" :placeholder="t('page.order.refundStatus')" clearable style="width: 140px" @change="reload">
            <el-option v-for="item in refundTypes" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-button :loading="loading" @click="loadList">{{ t("common.refresh") }}</el-button>
        </div>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="order_num" :label="t('page.order.orderNo')" min-width="170" />
      <el-table-column :label="t('page.order.product')" min-width="220">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div>{{ row.name || "-" }}</div>
              <div class="muted">{{ t("page.order.productId", { id: row.product_id }) }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="user_id" :label="t('page.order.buyerId')" width="90" />
      <el-table-column prop="boss_id" :label="t('page.order.sellerId')" width="90" />
      <el-table-column :label="t('common.amount')" width="110">
        <template #default="{ row }">¥{{ totalAmount(row).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column :label="t('page.order.orderStatus')" width="130">
        <template #default="{ row }">
          <el-tag :type="orderTypeTag(row.type)">{{ orderTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.order.refundStatus')" width="110">
        <template #default="{ row }">
          <el-tag :type="refundTag(row.refund_status)">{{ refundText(row.refund_status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="refund_reason" :label="t('page.order.refundReason')" min-width="160" show-overflow-tooltip />
      <el-table-column :label="t('common.actions')" width="130" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.refund_status === 1"
            size="small"
            type="primary"
            @click="approveRefund(row)"
          >
            {{ t("page.order.refundAudit") }}
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
import { computed, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { approveOrderRefund, getAdminOrderList } from "@/api";
import { orderStatusText, refundStatusText } from "@/utils/status-labels";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const typeFilter = ref<number | undefined>();
const refundFilter = ref<number | undefined>();

const orderTypes = computed(() =>
  [1, 2, 3, 4, 5, 6].map((value) => ({ value, label: orderStatusText(value) })),
);
const refundTypes = computed(() =>
  [0, 1, 2].map((value) => ({ value, label: refundStatusText(value) })),
);

function orderTypeText(type: number) {
  return orderStatusText(type);
}

function orderTypeTag(type: number) {
  return ({ 1: "info", 2: "warning", 3: "primary", 4: "success", 5: "danger", 6: "info" } as any)[type] ?? "info";
}

function refundText(status: number) {
  return refundStatusText(status);
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
    t("page.order.refundPrompt", {
      orderNo: row.order_num,
      amount: totalAmount(row).toFixed(2),
    }),
    t("page.order.refundAudit"),
    {
      inputType: "password",
      inputPattern: /.+/,
      inputErrorMessage: t("page.order.fundKeyRequired"),
      confirmButtonText: t("page.order.refundConfirm"),
      cancelButtonText: t("common.cancel"),
    },
  );
  await approveOrderRefund({ order_id: row.id, key: value });
  ElMessage.success(t("page.order.refundSuccess"));
  requestAdminPendingCountsRefresh();
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
