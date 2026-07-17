<template>
  <el-card>
    <template #header>
      <div class="card-header">
        <span>优惠券管理</span>
        <el-button type="primary" @click="openCreate">新增优惠券</el-button>
      </div>
    </template>

    <el-table :data="list" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="名称" min-width="160" />
      <el-table-column label="类型" width="110">
        <template #default="{ row }">
          <el-tag>{{ row.coupon_type === 1 ? "满减券" : "折扣券" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="优惠" width="120">
        <template #default="{ row }">
          {{ discountText(row) }}
        </template>
      </el-table-column>
      <el-table-column prop="min_amount" label="门槛" width="110" />
      <el-table-column prop="stock" label="库存" width="90" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="isActive(row) ? 'success' : 'info'">{{ isActive(row) ? "可领取" : "已下线" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="过期时间" min-width="180">
        <template #default="{ row }">{{ formatTime(row.expire_at) }}</template>
      </el-table-column>
      <el-table-column label="操作" width="110">
        <template #default="{ row }">
          <el-button v-if="isActive(row)" size="small" type="warning" @click="offline(row.id)">下线</el-button>
          <span v-else class="muted">-</span>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" title="新增优惠券" width="520px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="名称">
          <el-input v-model="form.name" placeholder="例如：新人满减券" />
        </el-form-item>
        <el-form-item label="类型">
          <el-radio-group v-model="form.coupon_type">
            <el-radio-button :value="1">满减券</el-radio-button>
            <el-radio-button :value="2">折扣券</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="form.coupon_type === 1" label="减免金额">
          <div class="field-line">
            <el-input-number
              v-model="form.discount_amount"
              :min="0"
              :precision="2"
              :step="1"
            />
            <span class="field-hint">元</span>
          </div>
        </el-form-item>
        <el-form-item v-else label="折扣比例">
          <div class="field-line">
            <el-input-number
              v-model="form.discount_percent"
              :min="1"
              :max="100"
              :precision="0"
              :step="1"
            />
            <span class="field-hint">
              支付 {{ form.discount_percent }}%，约
              {{ percentToFold(form.discount_percent) }} 折
            </span>
          </div>
        </el-form-item>
        <el-form-item label="使用门槛">
          <el-input-number v-model="form.min_amount" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="库存">
          <el-input-number v-model="form.stock" :min="-1" />
        </el-form-item>
        <el-form-item label="过期时间">
          <el-date-picker
            v-model="form.expire_at"
            type="datetime"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
            placeholder="选择过期时间"
            style="width: 100%"
          />
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
import { createAdminCoupon, getAdminCouponList, offlineAdminCoupon } from "@/api";

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
    return `减 ¥${Number(row.discount || 0).toFixed(2)}`;
  }
  return `${(Number(row.discount || 0) * 100).toFixed(0)}%，${(
    Number(row.discount || 0) * 10
  )
    .toFixed(1)
    .replace(/\.0$/, "")} 折`;
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
  if (!form.name.trim()) return ElMessage.warning("请输入优惠券名称");
  if (!form.expire_at) return ElMessage.warning("请选择过期时间");
  if (form.coupon_type === 1 && form.discount_amount <= 0) {
    return ElMessage.warning("请输入减免金额");
  }
  if (
    form.coupon_type === 2 &&
    (form.discount_percent <= 0 || form.discount_percent > 100)
  ) {
    return ElMessage.warning("请输入 1-100 的折扣比例");
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
    ElMessage.success("创建成功");
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function offline(id: number) {
  await ElMessageBox.confirm("确认下线该优惠券？", "提示", { type: "warning" });
  await offlineAdminCoupon({ id });
  ElMessage.success("已下线");
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
