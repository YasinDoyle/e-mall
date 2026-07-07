<template>
  <div class="login-page">
    <el-card class="login-card">
      <div class="login-title">🛒 E-Mall 管理后台</div>
      <el-form ref="formRef" :model="form" :rules="rules" label-position="top">
        <el-form-item label="用户名" prop="user_name">
          <el-input v-model="form.user_name" placeholder="请输入管理员账号" />
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
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import { adminLogin } from "@/api";
import { useAdminStore } from "@/stores/admin";

const router = useRouter();
const store = useAdminStore();
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
    const res: any = await adminLogin(form);
    // 验证是否为管理员（后端会在非管理员时鉴权失败，这里做前端提示）
    if (!res.data?.user) {
      ElMessage.error("账号或密码错误");
      return;
    }
    store.setToken(res.data.access_token);
    localStorage.setItem("admin_refresh_token", res.data.refresh_token);
    store.setAdminInfo({
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
}
.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  margin-bottom: 24px;
}
</style>
