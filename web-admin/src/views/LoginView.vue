<template>
  <div class="login-page">
    <el-card class="login-card">
      <el-select
        v-model="selectedLocale"
        class="locale-select"
        size="small"
        :aria-label="t('common.language')"
      >
        <el-option
          v-for="locale in supportedLocales"
          :key="locale.value"
          :label="locale.label"
          :value="locale.value"
        />
      </el-select>
      <div class="login-title">🛒 {{ t("admin.title") }}</div>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('auth.username')" prop="user_name">
          <el-input v-model="form.user_name" :placeholder="t('auth.adminUsernamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('auth.password')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('auth.passwordPlaceholder')"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%"
            :loading="loading"
            @click="handleLogin"
          >
            {{ t("common.login") }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import { useI18n } from "vue-i18n";
import { adminLogin, getStatsOverview } from "@/api";
import { useAdminStore } from "@/stores/admin";
import { ApiErrorCode } from "@/utils/api-error";
import { beginAdminLoginSession } from "@/utils/session";
import { currentLocale, setLocale, supportedLocales } from "@/locales";

const router = useRouter();
const store = useAdminStore();
const { t } = useI18n();
const selectedLocale = computed({
  get: () => currentLocale.value,
  set: (locale: string) => setLocale(locale),
});
const formRef = ref<FormInstance>();
const loading = ref(false);
const form = reactive({ user_name: "", password: "" });
const rules = computed(() => ({
  user_name: [
    { required: true, message: t("auth.usernamePlaceholder"), trigger: "blur" },
  ],
  password: [
    { required: true, message: t("auth.passwordPlaceholder"), trigger: "blur" },
  ],
}));

async function handleLogin() {
  await formRef.value?.validate();
  loading.value = true;
  try {
    const res: any = await adminLogin(form);
    if (!res.data?.user) {
      ElMessage.error(t("auth.accountInvalid"));
      return;
    }
    beginAdminLoginSession();
    store.setToken(res.data.access_token);
    store.setRefreshToken(res.data.refresh_token);
    try {
      await getStatsOverview();
    } catch (error: any) {
      const status = error?.response?.data?.status;
      store.logout();
      if (status === ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_FAIL) {
        return;
      }
      return;
    }
    store.setAdminInfo({
      id: res.data.user.id,
      user_name: res.data.user.user_name,
      nick_name: res.data.user.nick_name,
    });
    router.push("/");
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f2f5;
}
.login-card {
  width: 400px;
  position: relative;
}
.locale-select {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 104px;
}
.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 24px;
}
</style>
