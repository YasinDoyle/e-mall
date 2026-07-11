<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>秒杀管理</span>
        <el-button type="primary" @click="openCreate">新增秒杀</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="product_id" label="商品ID" width="100" />
      <el-table-column prop="boss_id" label="卖家ID" width="100" />
      <el-table-column prop="title" label="活动标题" min-width="180" />
      <el-table-column label="秒杀价" width="110">
        <template #default="{ row }">¥{{ Number(row.money || 0).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="num" label="库存" width="90" />
      <el-table-column prop="custom_name" label="备注" min-width="140" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.num > 0 ? 'success' : 'info'">{{ row.num > 0 ? "进行中" : "售罄" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="remove(row.id)">删除</el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑秒杀' : '新增秒杀'" width="540px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="商品ID">
          <el-input-number v-model="form.product_id" :min="1" />
        </el-form-item>
        <el-form-item label="卖家ID">
          <el-input-number v-model="form.boss_id" :min="1" />
        </el-form-item>
        <el-form-item label="活动标题">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="秒杀价">
          <el-input-number v-model="form.money" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.num" :min="0" />
        </el-form-item>
        <el-form-item label="扩展ID">
          <el-input-number v-model="form.custom_id" :min="0" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.custom_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  createAdminFlashSale,
  deleteAdminFlashSale,
  getAdminFlashSaleList,
  updateAdminFlashSale,
} from "@/api";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);
const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const form = reactive({
  product_id: 1,
  boss_id: 1,
  title: "",
  money: 0,
  num: 0,
  custom_id: 0,
  custom_name: "",
});

function resetForm(row?: any) {
  editingId.value = row?.id ?? null;
  Object.assign(form, {
    product_id: row?.product_id ?? 1,
    boss_id: row?.boss_id ?? 1,
    title: row?.title ?? "",
    money: row?.money ?? 0,
    num: row?.num ?? 0,
    custom_id: row?.custom_id ?? 0,
    custom_name: row?.custom_name ?? "",
  });
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: any) {
  resetForm(row);
  dialogVisible.value = true;
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAdminFlashSaleList({
      page_num: page.value,
      page_size: pageSize,
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.title.trim()) return ElMessage.warning("请输入活动标题");
  saving.value = true;
  try {
    if (editingId.value) {
      await updateAdminFlashSale({ id: editingId.value, ...form });
      ElMessage.success("更新成功");
    } else {
      await createAdminFlashSale({ ...form });
      ElMessage.success("创建成功");
    }
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function remove(id: number) {
  await ElMessageBox.confirm("确认删除该秒杀活动？", "提示", { type: "warning" });
  await deleteAdminFlashSale({ id });
  ElMessage.success("删除成功");
  loadList();
}

onMounted(loadList);
</script>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.pager {
  margin-top: 16px;
  justify-content: flex-end;
}
</style>
