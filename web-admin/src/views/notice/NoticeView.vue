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
        <span>{{ t("page.notice.title") }}</span>
        <el-button type="primary" @click="openCreate">{{ t("page.notice.create") }}</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Text" :label="t('page.notice.content')" show-overflow-tooltip />
      <el-table-column :label="t('common.actions')" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">{{ t("common.edit") }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.ID)"
            >{{ t("common.delete") }}</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? t('page.notice.edit') : t('page.notice.create')"
      width="500px"
    >
      <el-form :model="form" label-width="60px">
        <el-form-item :label="t('page.notice.content')">
          <el-input
            v-model="form.text"
            type="textarea"
            :rows="4"
            :placeholder="t('page.notice.contentPlaceholder')"
          />
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
import { getNoticeList, createNotice, updateNotice, deleteNotice } from "@/api";
import { t } from "@/locales";

const list = ref<any[]>([]);
const dialogVisible = ref(false);
const saving = ref(false);
const form = reactive({ id: 0, text: "" });

async function loadList() {
  const res: any = await getNoticeList();
  list.value = res.data?.item ?? [];
}

function openCreate() {
  form.id = 0;
  form.text = "";
  dialogVisible.value = true;
}

function openEdit(row: any) {
  form.id = row.ID;
  form.text = row.Text;
  dialogVisible.value = true;
}

async function handleSave() {
  if (!form.text.trim()) return ElMessage.warning(t("page.notice.contentRequired"));
  saving.value = true;
  try {
    if (form.id) {
      await updateNotice({ id: form.id, text: form.text });
    } else {
      await createNotice({ text: form.text });
    }
    ElMessage.success(t("common.saveSuccess"));
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm(t("page.notice.deleteConfirm"), t("common.notice"), { type: "warning" });
  await deleteNotice({ id });
  ElMessage.success(t("common.deleteSuccess"));
  loadList();
}

onMounted(loadList);
</script>
