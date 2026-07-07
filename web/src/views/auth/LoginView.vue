<template>
  <div class="login-page">
    <el-card class="login-card">
      <h2 class="title">登录</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="form.user_name" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
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
            登录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="footer-links">
        没有账号？<RouterLink to="/register">立即注册</RouterLink>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter, useRoute } from "vue-router";
import type { FormInstance } from "element-plus";
import { userLogin } from "@/api/user";
import { getUserInfo } from "@/api/user";
import { useUserStore } from "@/stores/user";

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();

const formRef = ref<FormInstance>();
const loading = ref(false);
const form = reactive({ user_name: "", password: "" });

const rules = {
  user_name: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

async function handleLogin() {
  await formRef.value?.validate();
  loading.value = true;
  try {
    const res: any = await userLogin(form);
    // 后端返回字段为 access_token / refresh_token
    userStore.setToken(res.data.access_token);
    localStorage.setItem("refreshToken", res.data.refresh_token);
    const infoRes: any = await getUserInfo();
    userStore.setUserInfo(infoRes.data);
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
