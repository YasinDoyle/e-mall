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
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getAdminProductList, auditProduct, deleteProduct } from "@/api";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const auditFilter = ref<number | undefined>(undefined);

const statusText = (s: number) =>
  ({ 0: "待审核", 1: "已上架", 2: "已拒绝" })[s] ?? "未知";
const statusType = (s: number) =>
  (({ 0: "warning", 1: "success", 2: "danger" })[s] ?? "info") as any;

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

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该商品？", "提示", { type: "warning" });
  await deleteProduct({ id });
  ElMessage.success("删除成功");
  loadList();
}

onMounted(loadList);
</script>
