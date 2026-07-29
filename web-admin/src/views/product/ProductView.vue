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
        <span>商品管理</span>
        <el-radio-group v-model="auditFilter" @change="loadList">
          <el-radio-button :value="undefined">全部</el-radio-button>
          <el-radio-button :value="0">待审核</el-radio-button>
          <el-radio-button :value="1">已上架</el-radio-button>
          <el-radio-button :value="2">已拒绝</el-radio-button>
        </el-radio-group>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column label="图片" width="80">
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
      <el-table-column prop="name" label="商品名称" />
      <el-table-column prop="price" label="价格" width="90" />
      <el-table-column prop="boss_name" label="卖家" width="100" />
      <el-table-column label="审核状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.audit_status)">{{
            statusText(row.audit_status)
          }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="详情" width="90">
        <template #default="{ row }">
          <el-button size="small" @click="openDetail(row)">详情</el-button>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="{ row }">
          <el-button
            v-if="row.audit_status !== 1"
            size="small"
            type="success"
            @click="audit(row.id, 1)"
            >上架</el-button
          >
          <el-button
            v-if="row.audit_status !== 2"
            size="small"
            type="warning"
            @click="audit(row.id, 2)"
            >拒绝</el-button
          >
          <el-button size="small" type="danger" @click="handleDelete(row.id)"
            >删除</el-button
          >
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

    <el-drawer v-model="detailVisible" title="商品详情" size="520px">
      <template v-if="currentProduct">
        <img
          v-if="currentProduct.img_path"
          :src="currentProduct.img_path"
          class="detail-img"
        />
        <el-descriptions :column="1" border>
          <el-descriptions-item label="商品ID">
            {{ currentProduct.id }}
          </el-descriptions-item>
          <el-descriptions-item label="商品名称">
            {{ currentProduct.name }}
          </el-descriptions-item>
          <el-descriptions-item label="标题">
            {{ currentProduct.title || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="详情">
            {{ currentProduct.info || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="品牌">
            {{ currentProduct.brand || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="产地">
            {{ currentProduct.origin || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="规格">
            {{ currentProduct.specification || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="生产日期">
            {{ currentProduct.production_date || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="保质期">
            {{ currentProduct.shelf_life || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="服务保障">
            {{ currentProduct.service_guarantees || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="证书说明">
            {{ currentProduct.certificate_meta || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="价格">
            ¥{{ currentProduct.price }}
          </el-descriptions-item>
          <el-descriptions-item label="优惠价">
            ¥{{ currentProduct.discount_price }}
          </el-descriptions-item>
          <el-descriptions-item label="库存">
            {{ currentProduct.num }}
          </el-descriptions-item>
          <el-descriptions-item label="卖家">
            {{ currentProduct.boss_name }}（ID: {{ currentProduct.boss_id }}）
          </el-descriptions-item>
          <el-descriptions-item label="审核状态">
            <el-tag :type="statusType(currentProduct.audit_status)">
              {{ statusText(currentProduct.audit_status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="销售状态">
            <el-tag :type="currentProduct.on_sale ? 'success' : 'info'">
              {{ currentProduct.on_sale ? "销售中" : "未上架" }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        <div class="certificate-section">
          <div class="section-title">资质证书</div>
          <el-empty
            v-if="!currentProduct.certificates?.length"
            description="暂无资质材料"
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
  ({ 0: "待审核", 1: "已上架", 2: "已拒绝" })[s] ?? "未知";
const statusType = (s: number) =>
  (({ 0: "warning", 1: "success", 2: "danger" })[s] ?? "info") as any;
const certificateTypeText = (type: string) =>
  ({
    qualification: "资质证书",
    quality_inspection: "质检报告",
    authorization: "授权证书",
    other: "其他材料",
  })[type] ?? "其他材料";

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
  ElMessage.success(status === 1 ? "已上架" : "已拒绝");
  loadList();
}

function openDetail(row: any) {
  currentProduct.value = row;
  detailVisible.value = true;
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该商品？", "提示", { type: "warning" });
  await deleteProduct({ id });
  ElMessage.success("删除成功");
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
