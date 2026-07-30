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
        <span>{{ t("address.title") }}</span>
        <el-button type="primary" @click="openCreate">{{ t("address.add") }}</el-button>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="3" animated />
    <el-empty v-else-if="!list.length" :description="t('address.empty')" />

    <template v-else>
      <div v-for="addr in list" :key="addr.id" class="addr-card">
        <div class="addr-info">
          <b>{{ addr.name }}</b>
          <span style="margin-left: 12px; color: #666">{{ addr.phone }}</span>
        </div>
        <div class="addr-detail">{{ addr.address }}</div>
        <div class="addr-actions">
          <el-button size="small" @click="openEdit(addr)">{{ t("address.edit") }}</el-button>
          <el-button size="small" type="danger" @click="handleDelete(addr.id)"
            >{{ t("address.delete") }}</el-button
          >
        </div>
      </div>
    </template>

    <el-dialog
      v-model="dialogVisible"
      :title="form.id ? t('address.editTitle') : t('address.addTitle')"
      width="460px"
    >
      <el-form :model="form" label-width="70px">
        <el-form-item :label="t('address.name')"
          ><el-input v-model="form.name"
        /></el-form-item>
        <el-form-item :label="t('address.phone')"
          ><el-input v-model="form.phone"
        /></el-form-item>
        <el-form-item :label="t('address.address')"
          ><el-input v-model="form.address" type="textarea" :rows="3"
        /></el-form-item>
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
import { useI18n } from "vue-i18n";
import {
  getAddressList,
  createAddress,
  updateAddress,
  deleteAddress,
} from "@/api/address";

const list = ref<any[]>([]);
const { t } = useI18n();
const dialogVisible = ref(false);
const saving = ref(false);
const loading = ref(false);
const form = reactive({ id: 0, name: "", phone: "", address: "" });

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAddressList();
    list.value = res.data?.item ?? [];
  } catch {
    list.value = [];
  } finally {
    loading.value = false;
  }
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
  if (!form.name.trim() || !form.phone.trim() || !form.address.trim())
    return ElMessage.warning(t("address.incomplete"));
  saving.value = true;
  try {
    if (form.id) {
      await updateAddress({
        id: form.id,
        name: form.name.trim(),
        phone: form.phone.trim(),
        address: form.address.trim(),
      });
    } else {
      await createAddress({
        name: form.name.trim(),
        phone: form.phone.trim(),
        address: form.address.trim(),
      });
    }
    ElMessage.success(t("address.saveSuccess"));
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  try {
    await ElMessageBox.confirm(t("address.deleteConfirm"), t("dialog.warningTitle"), {
      type: "warning",
    });
    await deleteAddress({ id });
    ElMessage.success(t("address.deleteSuccess"));
    loadList();
  } catch {}
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
