<template>
  <el-card>
    <template #header>
      <div
        style="
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <span>{{ t("page.product.title") }}</span>
        <el-radio-group v-model="auditFilter" @change="loadList">
          <el-radio-button :value="undefined">{{ t("common.all") }}</el-radio-button>
          <el-radio-button :value="0">{{ t("status.productAudit.pending") }}</el-radio-button>
          <el-radio-button :value="1">{{ t("status.productAudit.approved") }}</el-radio-button>
          <el-radio-button :value="2">{{ t("status.productAudit.rejected") }}</el-radio-button>
        </el-radio-group>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column :label="t('page.product.image')" width="80">
        <template #default="{ row }">
          <img
            :src="row.img_path"
            style="
              width: 50px;
              height: 50px;
              object-fit: cover;
              border-radius: 4px;
            "
          />
        </template>
      </el-table-column>
      <el-table-column prop="name" :label="t('page.product.name')" />
      <el-table-column prop="price" :label="t('page.product.price')" width="90" />
      <el-table-column prop="boss_name" :label="t('page.product.seller')" width="100" />
      <el-table-column :label="t('page.product.auditStatus')" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.audit_status)">{{
            statusText(row.audit_status)
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.detail')" width="90">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">{{ t("common.detail") }}</el-button>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="240">
        <template #default="{ row }">
          <div class="product-actions">
            <el-button
              v-if="row.audit_status !== 2"
              size="small"
              type="warning"
              @click="audit(row.id, 2)"
              >{{ t("page.product.rejectAction") }}</el-button
            >
            <el-button
              v-if="row.audit_status !== 1"
              size="small"
              type="success"
              @click="audit(row.id, 1)"
              >{{ t("page.product.listAction") }}</el-button
            >
            <el-button size="small" type="danger" @click="handleDelete(row.id)"
              >{{ t("common.delete") }}</el-button
            >
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 16px; justify-content: flex-end"
      @current-change="loadList"
    />

    <el-drawer v-model="detailVisible" :title="t('page.product.productDetail')" size="520px">
      <template v-if="currentProduct">
        <img
          v-if="currentProduct.img_path"
          :src="currentProduct.img_path"
          class="detail-img"
        />
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="t('page.product.productId')">
            {{ currentProduct.id }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.name')">
            {{ currentProduct.name }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.subtitle')">
            {{ currentProduct.title || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.detail')">
            {{ currentProduct.info || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.brand')">
            {{ currentProduct.brand || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.origin')">
            {{ currentProduct.origin || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.specification')">
            {{ currentProduct.specification || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.productionDate')">
            {{ currentProduct.production_date || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.shelfLife')">
            {{ currentProduct.shelf_life || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.guarantees')">
            {{ currentProduct.service_guarantees || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.certificateMeta')">
            {{ currentProduct.certificate_meta || "-" }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.price')">
            ¥{{ currentProduct.price }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.discountPrice')">
            ¥{{ currentProduct.discount_price }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.stock')">
            {{ currentProduct.num }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.seller')">
            {{ currentProduct.boss_name }}（ID: {{ currentProduct.boss_id }}）
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.auditStatus')">
            <el-tag :type="statusType(currentProduct.audit_status)">
              {{ statusText(currentProduct.audit_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('page.product.saleStatus')">
            <el-tag :type="currentProduct.on_sale ? 'success' : 'info'">
              {{ currentProduct.on_sale ? t("status.sale.onSale") : t("status.sale.offSale") }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        <div class="certificate-section">
          <div class="section-title">{{ t("page.product.certificates") }}</div>
          <el-empty
            v-if="!currentProduct.certificates?.length"
            :description="t('page.product.emptyCertificates')"
          />
          <div v-else class="certificate-grid">
            <div
              v-for="certificate in currentProduct.certificates"
              :key="certificate.id"
              class="certificate-card"
            >
              <el-image
                :src="certificate.file_path"
                :preview-src-list="certificateImages"
                fit="cover"
                class="certificate-img"
              />
              <div class="certificate-name">{{ certificate.name || "-" }}</div>
              <div class="muted">
                {{ certificateTypeText(certificate.certificate_type) }}
              </div>
            </div>
          </div>
        </div>
      </template>
    </el-drawer>
  </el-card>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getAdminProductList, auditProduct, deleteProduct } from "@/api";
import { t } from "@/locales";
import { requestAdminPendingCountsRefresh } from "@/utils/adminPending";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const auditFilter = ref<number | undefined>(undefined);
const detailVisible = ref(false);
const currentProduct = ref<any | null>(null);
const certificateImages = computed(() =>
  (currentProduct.value?.certificates ?? [])
    .map((item: any) => item.file_path)
    .filter(Boolean),
);

const statusText = (s: number) =>
  ({
    0: t("status.productAudit.pending"),
    1: t("status.productAudit.approved"),
    2: t("status.productAudit.rejected"),
  })[s] ?? t("common.unknown");
const statusType = (s: number) =>
  (({ 0: "warning", 1: "success", 2: "danger" })[s] ?? "info") as any;
const certificateTypeText = (type: string) =>
  ({
    qualification: t("page.product.certificates"),
    quality_inspection: t("page.product.qualityInspection"),
    authorization: t("page.product.authorization"),
    other: t("page.product.otherCertificate"),
  })[type] ?? t("page.product.otherCertificate");

async function loadList() {
  const res: any = await getAdminProductList({
    page_num: page.value,
    page_size: pageSize,
    ...(auditFilter.value !== undefined
      ? { audit_status: auditFilter.value }
      : {}),
  });
  list.value = res.data?.item ?? [];
  total.value = res.data?.total ?? 0;
}

async function audit(id: number, status: number) {
  await auditProduct({ id, audit_status: status });
  ElMessage.success(status === 1 ? t("status.productAudit.approved") : t("status.productAudit.rejected"));
  requestAdminPendingCountsRefresh();
  loadList();
}

function openDetail(row: any) {
  currentProduct.value = row;
  detailVisible.value = true;
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm(t("page.product.confirmDelete"), t("common.notice"), { type: "warning" });
  await deleteProduct({ id });
  ElMessage.success(t("common.deleteSuccess"));
  loadList();
}

onMounted(loadList);
</script>

<style scoped>
.detail-img {
  width: 100%;
  max-height: 260px;
  object-fit: contain;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  margin-bottom: 16px;
}
.section-title {
  margin: 18px 0 10px;
  font-weight: 600;
}
.certificate-section {
  margin-top: 4px;
}
.product-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}
.product-actions :deep(.el-button + .el-button) {
  margin-left: 0;
}
.certificate-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}
.certificate-card {
  min-width: 0;
}
.certificate-img {
  width: 100%;
  height: 140px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
.certificate-name {
  margin-top: 6px;
  color: #303133;
  font-weight: 500;
}
.muted {
  color: #909399;
  font-size: 12px;
}
</style>
