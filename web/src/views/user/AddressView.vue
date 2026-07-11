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
        <span>收货地址</span>
        <el-button type="primary" @click="openCreate">新增地址</el-button>
      </div>
    </template>

    <el-empty v-if="!list.length" description="还没有收货地址" />

    <div v-for="addr in list" :key="addr.id" class="addr-card">
      <div class="addr-info">
        <b>{{ addr.name }}</b>
        <span style="margin-left: 12px; color: #666">{{ addr.phone }}</span>
      </div>
      <div class="addr-detail">{{ addr.address }}</div>
      <div class="addr-actions">
        <el-button size="small" @click="openEdit(addr)">编辑</el-button>
        <el-button size="small" type="danger" @click="handleDelete(addr.id)"
          >删除</el-button
        >
      </div>
    </div>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? '编辑地址' : '新增地址'"
      width="460px"
    >
      <el-form :model="form" label-width="70px">
        <el-form-item label="姓名"
          ><el-input v-model="form.name"
        /></el-form-item>
        <el-form-item label="手机号"
          ><el-input v-model="form.phone"
        /></el-form-item>
        <el-form-item label="地址"
          ><el-input v-model="form.address" type="textarea" :rows="3"
        /></el-form-item>
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
  getAddressList,
  createAddress,
  updateAddress,
  deleteAddress,
} from "@/api/address";

const list = ref<any[]>([]);
const dialogVisible = ref(false);
const saving = ref(false);
const form = reactive({ id: 0, name: "", phone: "", address: "" });

async function loadList() {
  const res: any = await getAddressList();
  list.value = res.data?.item ?? [];
}

function openCreate() {
  Object.assign(form, { id: 0, name: "", phone: "", address: "" });
  dialogVisible.value = true;
}

function openEdit(addr: any) {
  Object.assign(form, addr);
  dialogVisible.value = true;
}

async function handleSave() {
  if (!form.name || !form.phone || !form.address)
    return ElMessage.warning("请填写完整信息");
  saving.value = true;
  try {
    if (form.id) {
      await updateAddress({
        id: form.id,
        name: form.name,
        phone: form.phone,
        address: form.address,
      });
    } else {
      await createAddress({
        name: form.name,
        phone: form.phone,
        address: form.address,
      });
    }
    ElMessage.success("保存成功");
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该地址？", "提示", { type: "warning" });
  await deleteAddress({ id });
  ElMessage.success("已删除");
  loadList();
}

onMounted(loadList);
</script>

<style scoped>
.addr-card {
  border: 1px solid #eee;
  border-radius: 6px;
  padding: 12px 16px;
  margin-bottom: 10px;
}
.addr-info {
  font-size: 15px;
  margin-bottom: 4px;
}
.addr-detail {
  color: #666;
  font-size: 13px;
}
.addr-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 8px;
}
</style>
