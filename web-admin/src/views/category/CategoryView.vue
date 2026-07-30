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
        <span>{{ t("page.category.title") }}</span>
        <el-button type="primary" @click="openCreate">{{ t("page.category.create") }}</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="category_name" :label="t('page.category.name')" />
      <el-table-column :label="t('common.actions')" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">{{ t("common.edit") }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)"
            >{{ t("common.delete") }}</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? t('page.category.edit') : t('page.category.create')"
      width="400px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item :label="t('page.category.name')">
          <el-input v-model="form.category_name" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >{{ t("common.save") }}</el-button
        >
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  getCategoryList,
  createCategory,
  updateCategory,
  deleteCategory,
} from "@/api";
import { t } from "@/locales";

const list = ref<any[]>([]);
const dialogVisible = ref(false);
const saving = ref(false);
const form = reactive({ id: 0, category_name: "" });

async function loadList() {
  const res: any = await getCategoryList();
  list.value = res.data?.item ?? [];
}

function openCreate() {
  form.id = 0;
  form.category_name = "";
  dialogVisible.value = true;
}

function openEdit(row: any) {
  form.id = row.id;
  form.category_name = row.category_name;
  dialogVisible.value = true;
}

async function handleSave() {
  if (!form.category_name.trim()) return ElMessage.warning(t("page.category.nameRequired"));
  saving.value = true;
  try {
    if (form.id) {
      await updateCategory({ id: form.id, category_name: form.category_name });
    } else {
      await createCategory({ category_name: form.category_name });
    }
    ElMessage.success(t("common.saveSuccess"));
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm(t("page.category.deleteConfirm"), t("common.notice"), { type: "warning" });
  await deleteCategory({ id });
  ElMessage.success(t("common.deleteSuccess"));
  loadList();
}

onMounted(loadList);
</script>
