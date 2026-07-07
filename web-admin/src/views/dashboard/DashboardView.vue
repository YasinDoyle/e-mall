<template>
  <div>
    <el-row :gutter="16" style="margin-bottom: 20px">
      <el-col :span="6" v-for="card in statCards" :key="card.label">
        <el-card shadow="hover">
          <div style="display: flex; align-items: center; gap: 12px">
            <el-icon :size="36" :color="card.color"
              ><component :is="card.icon"
            /></el-icon>
            <div>
              <div style="font-size: 24px; font-weight: bold">
                {{ card.value }}
              </div>
              <div style="color: #999; font-size: 13px">{{ card.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card>
      <template #header>待审核商品</template>
      <el-table :data="pendingProducts" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="商品名称" />
        <el-table-column prop="price" label="价格" width="100" />
        <el-table-column prop="boss_name" label="卖家" width="120" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="audit(row.id, 1)"
              >上架</el-button
            >
            <el-button size="small" type="danger" @click="audit(row.id, 2)"
              >拒绝</el-button
            >
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getAdminProductList, auditProduct, getUserList } from "@/api";

const pendingProducts = ref<any[]>([]);
const statCards = ref([
  { label: "待审核商品", value: 0, icon: "Goods", color: "#e6a23c" },
  { label: "注册用户", value: 0, icon: "User", color: "#409eff" },
]);

async function loadData() {
  try {
    const [productRes, userRes]: any[] = await Promise.all([
      getAdminProductList({ page_num: 1, page_size: 100, audit_status: 0 }),
      getUserList({ page_num: 1, page_size: 1 }),
    ]);
    pendingProducts.value = productRes.data?.item ?? [];
    statCards.value[0].value = productRes.data?.total ?? 0;
    statCards.value[1].value = userRes.data?.total ?? 0;
  } catch (err) {
    console.error(err);
  }
}

async function audit(id: number, status: number) {
  await auditProduct({ id, audit_status: status });
  ElMessage.success(status === 1 ? "已上架" : "已拒绝");
  loadData();
}

onMounted(loadData);
</script>
