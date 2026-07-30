<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>{{ t("page.coupon.title") }}</span>
        <el-button type="primary" @click="openCreate">{{ t("page.coupon.create") }}</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" :label="t('page.coupon.name')" min-width="160" />
      <el-table-column :label="t('page.coupon.type')" width="110">
        <template #default="{ row }">
          <el-tag>{{ row.coupon_type === 1 ? t("page.coupon.fullReduction") : t("page.coupon.discountCoupon") }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.coupon.discount')" width="120">
        <template #default="{ row }">
          {{ discountText(row) }}
        </template>
      </el-table-column>
      <el-table-column prop="min_amount" :label="t('page.coupon.threshold')" width="110" />
      <el-table-column prop="stock" :label="t('page.coupon.stock')" width="90" />
      <el-table-column :label="t('common.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="isActive(row) ? 'success' : 'info'">{{ isActive(row) ? t("page.coupon.claimable") : t("page.coupon.offline") }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.coupon.expireAt')" min-width="180">
        <template #default="{ row }">{{ formatTime(row.expire_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="110">
        <template #default="{ row }">
          <el-button v-if="isActive(row)" size="small" type="warning" @click="offline(row.id)">{{ t("page.coupon.offlineAction") }}</el-button>
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="t('page.coupon.create')" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item :label="t('page.coupon.name')">
          <el-input v-model="form.name" :placeholder="t('page.coupon.namePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('page.coupon.type')">
          <el-radio-group v-model="form.coupon_type">
            <el-radio-button :value="1">{{ t("page.coupon.fullReduction") }}</el-radio-button>
            <el-radio-button :value="2">{{ t("page.coupon.discountCoupon") }}</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 1" :label="t('page.coupon.discountAmount')">
          <div class="field-line">
            <el-input-number
              v-model="form.discount_amount"
              :min="0"
              :precision="2"
              :step="1"
            />
            <span class="field-hint">{{ t("page.coupon.yuan") }}</span>
          </div>
        </el-form-item>
        <el-form-item v-else :label="t('page.coupon.discountPercent')">
          <div class="field-line">
            <el-input-number
              v-model="form.discount_percent"
              :min="1"
              :max="100"
              :precision="0"
              :step="1"
            />
            <span class="field-hint">
              {{ t("page.coupon.percentHint", {
                percent: form.discount_percent,
                fold: percentToFold(form.discount_percent),
              }) }}
            </span>
          </div>
        </el-form-item>
        <el-form-item :label="t('page.coupon.threshold')">
          <el-input-number v-model="form.min_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item :label="t('page.coupon.stock')">
          <el-input-number v-model="form.stock" :min="-1" />
        </el-form-item>
        <el-form-item :label="t('page.coupon.expireAt')">
          <el-date-picker
            v-model="form.expire_at"
            type="datetime"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            :placeholder="t('page.coupon.expirePlaceholder')"
            style="width: 100%"
          />
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
import { createAdminCoupon, getAdminCouponList, offlineAdminCoupon } from "@/api";
import { t } from "@/locales";

const list = ref<any[]>([]);
const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const form = reactive({
  name: "",
  coupon_type: 1,
  discount_amount: 0,
  discount_percent: 90,
  min_amount: 0,
  stock: 100,
  expire_at: "",
});

function isActive(row: any) {
  return new Date(row.expire_at).getTime() > Date.now() && (row.stock > 0 || row.stock === -1);
}

function formatTime(value: string) {
  return value ? new Date(value).toLocaleString() : "-";
}

function percentToFold(percent: number) {
  return (Number(percent || 0) / 10).toFixed(1).replace(/\.0$/, "");
}

function discountText(row: any) {
  if (row.coupon_type === 1) {
    return t("page.coupon.moneyDiscount", { amount: Number(row.discount || 0).toFixed(2) });
  }
  return t("page.coupon.percentDiscount", {
    percent: (Number(row.discount || 0) * 100).toFixed(0),
    fold: (Number(row.discount || 0) * 10).toFixed(1).replace(/\.0$/, ""),
  });
}

function openCreate() {
  Object.assign(form, {
    name: "",
    coupon_type: 1,
    discount_amount: 0,
    discount_percent: 90,
    min_amount: 0,
    stock: 100,
    expire_at: "",
  });
  dialogVisible.value = true;
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAdminCouponList();
    list.value = res.data?.item ?? [];
  } finally {
    loading.value = false;
  }
}

async function save() {
  if (!form.name.trim()) return ElMessage.warning(t("page.coupon.nameRequired"));
  if (!form.expire_at) return ElMessage.warning(t("page.coupon.expireRequired"));
  if (form.coupon_type === 1 && form.discount_amount <= 0) {
    return ElMessage.warning(t("page.coupon.discountAmountRequired"));
  }
  if (
    form.coupon_type === 2 &&
    (form.discount_percent <= 0 || form.discount_percent > 100)
  ) {
    return ElMessage.warning(t("page.coupon.discountPercentRequired"));
  }
  const discount =
    form.coupon_type === 1
      ? form.discount_amount
      : Number((form.discount_percent / 100).toFixed(2));
  saving.value = true;
  try {
    await createAdminCoupon({
      name: form.name,
      coupon_type: form.coupon_type,
      discount,
      min_amount: form.min_amount,
      stock: form.stock,
      expire_at: form.expire_at,
    });
    ElMessage.success(t("common.createSuccess"));
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function offline(id: number) {
  await ElMessageBox.confirm(t("page.coupon.offlineConfirm"), t("common.notice"), { type: "warning" });
  await offlineAdminCoupon({ id });
  ElMessage.success(t("page.coupon.offlineSuccess"));
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
.muted {
  color: #909399;
  font-size: 12px;
}
.field-line {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}
.field-hint {
  color: #909399;
  font-size: 13px;
}
</style>
