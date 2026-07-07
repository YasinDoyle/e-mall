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
        <span>分类管理</span>
        <el-button type="primary" @click="openCreate">新增分类</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="category_name" label="分类名称" />
      <el-table-column label="操作" width="160">
        <template #default="{ row }">
          <el-button size="small" @click="openEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row.id)"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑分类' : '新增分类'"
      width="400px"
    >
      <el-form :model="form" label-width="80px">
        <el-form-item label="分类名称">
          <el-input v-model="form.category_name" />
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
import {
  getCategoryList,
  createCategory,
  updateCategory,
  deleteCategory,
} from "@/api";

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
  if (!form.category_name.trim()) return ElMessage.warning("请输入分类名称");
  saving.value = true;
  try {
    if (form.id) {
      await updateCategory({ id: form.id, category_name: form.category_name });
    } else {
      await createCategory({ category_name: form.category_name });
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该分类？", "提示", { type: "warning" });
  await deleteCategory({ id });
  ElMessage.success("删除成功");
  loadList();
}

onMounted(loadList);
</script>
