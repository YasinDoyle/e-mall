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
      <h2 class="title">{{ t("common.login") }}</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('auth.username')" prop="user_name">
          <el-input v-model="form.user_name" :placeholder="t('auth.usernamePlaceholder')" />
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
      <div class="footer-links">
        {{ t("auth.noAccount") }}<RouterLink to="/register">{{ t("auth.registerNow") }}</RouterLink>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import { useI18n } from "vue-i18n";
import { userLogin } from "@/api/user";
import { getUserInfo } from "@/api/user";
import { useUserStore } from "@/stores/user";
import { beginUserLoginSession } from "@/utils/session";
import { currentLocale, setLocale, supportedLocales } from "@/locales";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
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
    const res: any = await userLogin(form);
    beginUserLoginSession();
    // 后端返回字段为 access_token / refresh_token
    userStore.setToken(res.data.access_token);
    userStore.setRefreshToken(res.data.refresh_token);
    const infoRes: any = await getUserInfo();
    userStore.setUserInfo(infoRes.data);
    ElMessage.success(t("auth.loginSuccess"));
    const redirect = (route.query.redirect as string) || "/";
    router.push(redirect);
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
  background: #f5f5f5;
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
.title {
  text-align: center;
  margin-bottom: 24px;
}
.footer-links {
  text-align: center;
  margin-top: 12px;
  font-size: 14px;
}
</style>
