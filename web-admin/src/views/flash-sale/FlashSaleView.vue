<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.flashSale.title") }}</span>
        <el-button type="primary" @click="openCreate">{{ t("page.flashSale.create") }}</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="product_id" :label="t('page.flashSale.productId')" width="100" />
      <el-table-column prop="boss_id" :label="t('page.flashSale.sellerId')" width="100" />
      <el-table-column prop="title" :label="t('page.flashSale.activityTitle')" min-width="180" />
      <el-table-column :label="t('page.flashSale.price')" width="110">
        <template #default="{ row }">¥{{ Number(row.money || 0).toFixed(2) }}</template>
      </el-table-column>
      <el-table-column prop="num" :label="t('page.flashSale.stock')" width="90" />
      <el-table-column prop="custom_name" :label="t('page.flashSale.remark')" min-width="140" />
      <el-table-column :label="t('common.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.num > 0 ? 'success' : 'info'">{{ row.num > 0 ? t("page.flashSale.active") : t("page.flashSale.soldOut") }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="150" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">{{ t("common.edit") }}</el-button>
          <el-button size="small" type="danger" @click="remove(row.id)">{{ t("common.delete") }}</el-button>
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

    <el-dialog v-model="dialogVisible" :title="editingId ? t('page.flashSale.edit') : t('page.flashSale.create')" width="540px">
      <el-form :model="form" label-width="90px">
        <el-form-item :label="t('page.flashSale.productId')">
          <el-input-number v-model="form.product_id" :min="1" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.sellerId')">
          <el-input-number v-model="form.boss_id" :min="1" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.activityTitle')">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.price')">
          <el-input-number v-model="form.money" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.stock')">
          <el-input-number v-model="form.num" :min="0" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.customId')">
          <el-input-number v-model="form.custom_id" :min="0" />
        </el-form-item>
        <el-form-item :label="t('page.flashSale.remark')">
          <el-input v-model="form.custom_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="saving" @click="save">{{ t("common.save") }}</el-button>
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
import { t } from "@/locales";

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
  if (!form.title.trim()) return ElMessage.warning(t("page.flashSale.titleRequired"));
  saving.value = true;
  try {
    if (editingId.value) {
      await updateAdminFlashSale({ id: editingId.value, ...form });
      ElMessage.success(t("common.updateSuccess"));
    } else {
      await createAdminFlashSale({ ...form });
      ElMessage.success(t("common.createSuccess"));
    }
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function remove(id: number) {
  await ElMessageBox.confirm(t("page.flashSale.deleteConfirm"), t("common.notice"), { type: "warning" });
  await deleteAdminFlashSale({ id });
  ElMessage.success(t("common.deleteSuccess"));
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
