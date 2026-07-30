<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>{{ t("sellerCenter.apply.title") }}</span>
        <el-button :loading="loading" @click="loadProfile">
          {{ t("sellerCenter.apply.refreshStatus") }}
        </el-button>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="4" animated />

    <template v-else>
      <el-alert
        v-if="!profile"
        :title="t('sellerCenter.apply.noProfile')"
        type="info"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 0"
        :title="t('sellerCenter.apply.pending')"
        type="warning"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 1"
        :title="t('sellerCenter.apply.approved')"
        type="success"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 2"
        :title="t('sellerCenter.apply.rejected', { reason: profile.reject_reason || t('sellerCenter.apply.noRejectReason') })"
        type="error"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 3"
        :title="t('sellerCenter.apply.banned')"
        type="error"
        :closable="false"
        show-icon
      />

      <div v-if="profile" class="profile-box">
        <div class="profile-title">{{ profile.shop_name }}</div>
        <div class="profile-desc">
          {{ profile.description || t("sellerCenter.apply.noDescription") }}
        </div>
        <div class="profile-meta">
          {{ t("sellerCenter.apply.currentStatus") }}
          <el-tag :type="statusTag(profile.status)">
            {{ sellerStatusLabel(profile.status) }}
          </el-tag>
        </div>
      </div>

      <el-divider />

      <el-form
        v-if="canSubmit"
        label-width="90px"
        :model="form"
        class="apply-form"
      >
        <el-form-item :label="t('sellerCenter.apply.shopName')" required>
          <el-input
            v-model="form.shop_name"
            maxlength="80"
            show-word-limit
            :placeholder="t('sellerCenter.apply.shopNamePlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('sellerCenter.apply.description')">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="5"
            maxlength="500"
            show-word-limit
            :placeholder="t('sellerCenter.apply.descriptionPlaceholder')"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">
            {{
              profile?.status === 2
                ? t("sellerCenter.apply.resubmit")
                : t("sellerCenter.apply.submit")
            }}
          </el-button>
        </el-form-item>
      </el-form>

      <div v-else-if="profile?.status === 1" class="actions">
        <el-button type="success" @click="$router.push('/seller/account')">
          {{ t("sellerCenter.apply.viewAccount") }}
        </el-button>
        <el-button type="primary" @click="$router.push('/seller/products/new')">
          {{ t("sellerCenter.apply.publishProduct") }}
        </el-button>
        <el-button @click="$router.push('/seller/products')">
          {{ t("sellerCenter.apply.viewProducts") }}
        </el-button>
      </div>
    </template>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import { applySeller } from "@/api/seller";
import { useSellerStore } from "@/stores/seller";

const loading = ref(false);
const submitting = ref(false);
const sellerStore = useSellerStore();
const { t } = useI18n();
const profile = computed(() => sellerStore.profile);
const form = ref({ shop_name: "", description: "" });

const canSubmit = computed(
  () => !profile.value || profile.value.status === 2,
);

function statusTag(status: number) {
  return ({ 0: "warning", 1: "success", 2: "danger", 3: "info" } as any)[
    status
  ] ?? "info";
}

function sellerStatusLabel(status: number) {
  return (
    {
      0: t("status.seller.pending"),
      1: t("status.seller.approved"),
      2: t("status.seller.rejected"),
      3: t("status.seller.banned"),
    } as Record<number, string>
  )[status] ?? t("common.unknown");
}

async function loadProfile() {
  loading.value = true;
  try {
    await sellerStore.loadProfile({ force: true, silentError: true });
    form.value = {
      shop_name: profile.value?.shop_name ?? "",
      description: profile.value?.description ?? "",
    };
  } finally {
    loading.value = false;
  }
}

async function submit() {
  if (!form.value.shop_name.trim()) {
    return ElMessage.warning(t("sellerCenter.apply.shopNamePlaceholder"));
  }
  submitting.value = true;
  try {
    const res: any = await applySeller({
      shop_name: form.value.shop_name.trim(),
      description: form.value.description.trim(),
    });
    sellerStore.setProfile(res.data);
    ElMessage.success(t("sellerCenter.apply.submitSuccess"));
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  if (!sellerStore.loaded) {
    await sellerStore.loadProfile({ silentError: true });
  }
  form.value = {
    shop_name: profile.value?.shop_name ?? "",
    description: profile.value?.description ?? "",
  };
});
</script>

<style scoped>
.header,
.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.profile-box {
  margin-top: 16px;
  padding: 14px 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
}
.profile-title {
  font-weight: 600;
  color: #303133;
}
.profile-desc,
.profile-meta {
  margin-top: 8px;
  color: #606266;
  font-size: 14px;
}
.apply-form {
  max-width: 640px;
}
</style>
