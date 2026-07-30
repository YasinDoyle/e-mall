<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>{{ t("sellerCenter.productList.title") }}</span>
        <div class="header-actions">
          <el-button
            type="success"
            :disabled="!batchOnSaleCandidates.length"
            :loading="batchLoading"
            @click="batchOnSale"
          >
            {{ t("sellerCenter.productList.batchOnSale") }}
          </el-button>
          <el-button :loading="loading" @click="loadList">{{ t("common.refresh") }}</el-button>
          <el-button
            type="primary"
            :disabled="!sellerStore.isApproved"
            @click="$router.push('/seller/products/new')"
          >
            {{ t("sellerCenter.productList.publish") }}
          </el-button>
        </div>
      </div>
    </template>

    <el-alert
      v-if="!sellerStore.isApproved"
      :title="t('sellerCenter.productList.approvedOnly')"
      type="warning"
      :closable="false"
      show-icon
      class="notice"
    />

    <el-table
      :data="list"
      style="width: 100%"
      v-loading="loading"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="46" :selectable="canBatchSelect" />
      <el-table-column :label="t('sellerCenter.productList.product')" min-width="260">
        <template #default="{ row }">
          <div class="product-cell">
            <img v-if="row.img_path" :src="row.img_path" class="product-img" />
            <div>
              <div class="name">{{ row.name }}</div>
              <div class="muted">{{ row.title || row.info || "-" }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.productList.price')" width="140">
        <template #default="{ row }">
          <div>¥{{ row.discount_price || row.price }}</div>
          <div v-if="row.discount_price" class="muted">
            {{ t("sellerCenter.productList.originalPrice", { price: row.price }) }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="num" :label="t('sellerCenter.productList.stock')" width="90" />
      <el-table-column :label="t('sellerCenter.productList.audit')" width="110">
        <template #default="{ row }">
          <el-tag :type="auditTag(row.audit_status)">
            {{ auditText(row.audit_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.productList.saleStatus')" width="110">
        <template #default="{ row }">
          <el-tag :type="row.on_sale ? 'success' : 'info'">
            {{
              row.on_sale
                ? t("sellerCenter.productList.onSale")
                : t("sellerCenter.productList.offSale")
            }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('sellerCenter.productList.actions')" width="240" fixed="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :disabled="!sellerStore.isApproved"
            @click="$router.push(`/seller/products/${row.id}/edit`)"
          >
            {{ t("sellerCenter.productList.edit") }}
          </el-button>
          <el-button
            size="small"
            :type="row.on_sale ? 'warning' : 'success'"
            :disabled="!canToggle(row) || !sellerStore.isApproved"
            @click="toggleSale(row)"
          >
            {{
              row.on_sale
                ? t("sellerCenter.productList.offSaleAction")
                : t("sellerCenter.productList.onSaleAction")
            }}
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
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { useI18n } from "vue-i18n";
import { getSellerProductList, setSellerProductOnSale } from "@/api/product";
import { useSellerStore } from "@/stores/seller";

const sellerStore = useSellerStore();
const { t } = useI18n();
const list = ref<any[]>([]);
const selectedRows = ref<any[]>([]);
const page = ref(1);
const pageSize = 10;
const total = ref(0);
const loading = ref(false);
const batchLoading = ref(false);

const batchOnSaleCandidates = computed(() =>
  selectedRows.value.filter((row) => canBatchSelect(row)),
);

function auditText(status: number) {
  return (
    {
      0: t("sellerCenter.productList.auditPending"),
      1: t("sellerCenter.productList.auditApproved"),
      2: t("sellerCenter.productList.auditRejected"),
    } as Record<number, string>
  )[status] ?? t("common.unknown");
}

function auditTag(status: number) {
  return ({ 0: "warning", 1: "success", 2: "danger" } as any)[status] ?? "info";
}

function canToggle(row: any) {
  return row.on_sale || row.audit_status === 1;
}

function canBatchSelect(row: any) {
  return sellerStore.isApproved && row.audit_status === 1 && !row.on_sale;
}

function handleSelectionChange(rows: any[]) {
  selectedRows.value = rows;
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getSellerProductList({
      page_num: page.value,
      page_size: pageSize,
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
    selectedRows.value = [];
  } finally {
    loading.value = false;
  }
}

async function toggleSale(row: any) {
  const next = !row.on_sale;
  const action = next
    ? t("sellerCenter.productList.onSaleAction")
    : t("sellerCenter.productList.offSaleAction");
  await ElMessageBox.confirm(
    t("sellerCenter.productList.toggleConfirm", { action, name: row.name }),
    t("dialog.warningTitle"),
    { type: "warning" },
  );
  await setSellerProductOnSale({ id: row.id, on_sale: next });
  ElMessage.success(t("sellerCenter.productList.toggleSuccess", { action }));
  loadList();
}

async function batchOnSale() {
  if (!batchOnSaleCandidates.value.length) {
    return ElMessage.warning(t("sellerCenter.productList.selectBatchWarning"));
  }
  await ElMessageBox.confirm(
    t("sellerCenter.productList.batchConfirm", {
      count: batchOnSaleCandidates.value.length,
    }),
    t("dialog.warningTitle"),
    { type: "warning" },
  );
  batchLoading.value = true;
  try {
    const results = await Promise.allSettled(
      batchOnSaleCandidates.value.map((row) =>
        setSellerProductOnSale({ id: row.id, on_sale: true }),
      ),
    );
    const failed = results.filter((item) => item.status === "rejected").length;
    if (failed > 0) {
      ElMessage.warning(t("sellerCenter.productList.batchFailed", { count: failed }));
    } else {
      ElMessage.success(t("sellerCenter.productList.batchSuccess"));
    }
    loadList();
  } finally {
    batchLoading.value = false;
  }
}

onMounted(async () => {
  await sellerStore.loadProfile({ silentError: true });
  await loadList();
});
</script>

<style scoped>
.header,
.header-actions,
.product-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}
.header {
  justify-content: space-between;
}
.notice {
  margin-bottom: 12px;
}
.product-img {
  width: 54px;
  height: 54px;
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
