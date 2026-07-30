<template>
  <div class="register-page">
    <el-card class="register-card">
      <h2 class="title">{{ t("register.title") }}</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item :label="t('register.username')" prop="user_name">
          <el-input v-model="form.user_name" :placeholder="t('register.usernamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('register.nickname')">
          <el-input v-model="form.nick_name" :placeholder="t('register.nicknamePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('register.email')" prop="email">
          <div class="email-code-row">
            <el-input v-model="form.email" :placeholder="t('register.emailPlaceholder')" />
            <el-button
              :loading="sendingCode"
              :disabled="codeCountdown > 0"
              @click="handleSendCode"
            >
              {{ codeButtonText }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item :label="t('register.emailCode')" prop="email_code">
          <el-input
            v-model="form.email_code"
            :placeholder="t('register.emailCodePlaceholder')"
            maxlength="6"
          />
        </el-form-item>
        <el-form-item :label="t('register.password')" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            :placeholder="t('register.passwordPlaceholder')"
            show-password
          />
        </el-form-item>
        <el-form-item :label="t('register.confirmPassword')" prop="password_confirm">
          <el-input
            v-model="form.password_confirm"
            type="password"
            :placeholder="t('register.confirmPasswordPlaceholder')"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            style="width: 100%"
            :loading="loading"
            @click="handleRegister"
          >
            {{ t("register.register") }}
          </el-button>
        </el-form-item>
      </el-form>
      <div class="footer-links">
        {{ t("register.existingAccount") }}<RouterLink to="/login">{{ t("register.loginNow") }}</RouterLink>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onUnmounted, ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import { sendRegisterEmailCode, userRegister } from "@/api/user";
import { t } from "@/locales";

const router = useRouter();
const formRef = ref<FormInstance>();
const loading = ref(false);
const sendingCode = ref(false);
const codeCountdown = ref(0);
let countdownTimer: number | undefined;
const form = reactive({
  user_name: "",
  nick_name: "",
  email: "",
  email_code: "",
  password: "",
  password_confirm: "",
});

const codeButtonText = computed(() =>
  codeCountdown.value > 0 ? `${codeCountdown.value}${t("register.countdownSuffix")}` : t("register.sendCode"),
);

const rules = {
  user_name: [{ required: true, message: t("register.usernameRequired"), trigger: "blur" }],
  email: [
    { required: true, message: t("register.emailRequired"), trigger: "blur" },
    { type: "email", message: t("register.emailFormatInvalid"), trigger: "blur" },
  ],
  email_code: [
    { required: true, message: t("register.emailCodeRequired"), trigger: "blur" },
    { len: 6, message: t("register.emailCodeLength"), trigger: "blur" },
  ],
  password: [
    { required: true, message: t("register.passwordRequired"), trigger: "blur" },
    { min: 6, message: t("register.passwordMinLength"), trigger: "blur" },
  ],
  password_confirm: [
    { required: true, message: t("register.confirmPasswordRequired"), trigger: "blur" },
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== form.password) {
          callback(new Error(t("register.passwordMismatch")));
          return;
        }
        callback();
      },
      trigger: "blur",
    },
  ],
};

function startCountdown() {
  codeCountdown.value = 60;
  if (countdownTimer) window.clearInterval(countdownTimer);
  countdownTimer = window.setInterval(() => {
    codeCountdown.value -= 1;
    if (codeCountdown.value <= 0 && countdownTimer) {
      window.clearInterval(countdownTimer);
      countdownTimer = undefined;
    }
  }, 1000);
}

async function handleSendCode() {
  if (!form.email) {
    return ElMessage.warning(t("register.emailEmpty"));
  }
  await formRef.value?.validateField("email");
  sendingCode.value = true;
  try {
    await sendRegisterEmailCode({ email: form.email });
    ElMessage.success(t("register.codeSent"));
    startCountdown();
  } finally {
    sendingCode.value = false;
  }
}

async function handleRegister() {
  await formRef.value?.validate();
  loading.value = true;
  try {
    await userRegister({
      ...form,
      nick_name: form.nick_name || form.user_name,
    });
    ElMessage.success(t("register.registerSuccess"));
    router.push("/login");
  } finally {
    loading.value = false;
  }
}

onUnmounted(() => {
  if (countdownTimer) window.clearInterval(countdownTimer);
});
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
.register-card {
  width: 400px;
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
.email-code-row {
  display: flex;
  gap: 8px;
}
.email-code-row .el-input {
  flex: 1;
}
</style>
