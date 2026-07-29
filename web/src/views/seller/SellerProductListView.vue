<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>我的商品</span>
        <div class="header-actions">
          <el-button
            type="success"
            :disabled="!batchOnSaleCandidates.length"
            :loading="batchLoading"
            @click="batchOnSale"
          >
            批量上架
          </el-button>
          <el-button :loading="loading" @click="loadList">刷新</el-button>
          <el-button
            type="primary"
            :disabled="!sellerStore.isApproved"
            @click="$router.push('/seller/products/new')"
          >
            发布商品
          </el-button>
        </div>
      </div>
    </template>

    <el-alert
      v-if="!sellerStore.isApproved"
      title="商家入驻审核通过后才可以发布和上架商品"
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
      <el-table-column label="商品" min-width="260">
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
      <el-table-column label="价格" width="140">
        <template #default="{ row }">
          <div>¥{{ row.discount_price || row.price }}</div>
          <div v-if="row.discount_price" class="muted">原价 ¥{{ row.price }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="num" label="库存" width="90" />
      <el-table-column label="审核" width="110">
        <template #default="{ row }">
          <el-tag :type="auditTag(row.audit_status)">
            {{ auditText(row.audit_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="销售状态" width="110">
        <template #default="{ row }">
          <el-tag :type="row.on_sale ? 'success' : 'info'">
            {{ row.on_sale ? "销售中" : "已下架" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="240" fixed="right">
        <template #default="{ row }">
          <el-button
            size="small"
            :disabled="!sellerStore.isApproved"
            @click="$router.push(`/seller/products/${row.id}/edit`)"
          >
            编辑
          </el-button>
          <el-button
            size="small"
            :type="row.on_sale ? 'warning' : 'success'"
            :disabled="!canToggle(row) || !sellerStore.isApproved"
            @click="toggleSale(row)"
          >
            {{ row.on_sale ? "下架" : "上架" }}
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
import { getSellerProductList, setSellerProductOnSale } from "@/api/product";
import { useSellerStore } from "@/stores/seller";

const sellerStore = useSellerStore();
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
  return ({ 0: "待审核", 1: "已通过", 2: "已拒绝" } as any)[status] ?? "未知";
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
  await ElMessageBox.confirm(
    `确认${next ? "上架" : "下架"}商品「${row.name}」？`,
    "提示",
    { type: "warning" },
  );
  await setSellerProductOnSale({ id: row.id, on_sale: next });
  ElMessage.success(next ? "商品已上架" : "商品已下架");
  loadList();
}

async function batchOnSale() {
  if (!batchOnSaleCandidates.value.length) {
    return ElMessage.warning("请选择已审核通过且未上架的商品");
  }
  await ElMessageBox.confirm(
    `确认批量上架 ${batchOnSaleCandidates.value.length} 个商品？`,
    "提示",
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
      ElMessage.warning(`部分商品上架失败：${failed} 个`);
    } else {
      ElMessage.success("批量上架完成");
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
