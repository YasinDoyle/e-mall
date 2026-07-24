<template>
  <div class="register-page">
    <el-card class="register-card">
      <h2 class="title">注册</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="form.user_name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="昵称">
          <el-input v-model="form.nick_name" placeholder="请输入昵称" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <div class="email-code-row">
            <el-input v-model="form.email" placeholder="请输入邮箱" />
            <el-button
              :loading="sendingCode"
              :disabled="codeCountdown > 0"
              @click="handleSendCode"
            >
              {{ codeButtonText }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item label="邮箱验证码" prop="email_code">
          <el-input
            v-model="form.email_code"
            placeholder="请输入6位邮箱验证码"
            maxlength="6"
          />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="password_confirm">
          <el-input
            v-model="form.password_confirm"
            type="password"
            placeholder="请再次输入密码"
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
            注册
          </el-button>
        </el-form-item>
      </el-form>
      <div class="footer-links">
        已有账号？<RouterLink to="/login">立即登录</RouterLink>
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
  codeCountdown.value > 0 ? `${codeCountdown.value}s` : "发送验证码",
);

const rules = {
  user_name: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  email: [
    { required: true, message: "请输入邮箱", trigger: "blur" },
    { type: "email", message: "邮箱格式不正确", trigger: "blur" },
  ],
  email_code: [
    { required: true, message: "请输入邮箱验证码", trigger: "blur" },
    { len: 6, message: "邮箱验证码为6位", trigger: "blur" },
  ],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, message: "密码至少6位", trigger: "blur" },
  ],
  password_confirm: [
    { required: true, message: "请再次输入密码", trigger: "blur" },
    {
      validator: (_rule: unknown, value: string, callback: (error?: Error) => void) => {
        if (value !== form.password) {
          callback(new Error("两次密码输入不一致"));
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
    return ElMessage.warning("请输入邮箱");
  }
  await formRef.value?.validateField("email");
  sendingCode.value = true;
  try {
    await sendRegisterEmailCode({ email: form.email });
    ElMessage.success("验证码已发送，请查收邮箱");
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
    ElMessage.success("注册成功，请登录");
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
