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
        <span>公告管理</span>
        <el-button type="primary" @click="openCreate">新增公告</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="ID" label="ID" width="80" />
      <el-table-column prop="Text" label="内容" show-overflow-tooltip />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.ID)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑公告' : '新增公告'"
      width="500px"
    >
      <el-form :model="form" label-width="60px">
        <el-form-item label="内容">
          <el-input
            v-model="form.text"
            type="textarea"
            :rows="4"
            placeholder="请输入公告内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >保存</el-button
        >
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getNoticeList, createNotice, updateNotice, deleteNotice } from "@/api";

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
  if (!form.text.trim()) return ElMessage.warning("请输入公告内容");
  saving.value = true;
  try {
    if (form.id) {
      await updateNotice({ id: form.id, text: form.text });
    } else {
      await createNotice({ text: form.text });
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该公告？", "提示", { type: "warning" });
  await deleteNotice({ id });
  ElMessage.success("删除成功");
  loadList();
}

onMounted(loadList);
</script>
