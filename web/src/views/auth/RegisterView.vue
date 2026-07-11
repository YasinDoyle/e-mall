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
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码（至少6位）"
            show-password
          />
        </el-form-item>
        <el-form-item label="支付密码（6位数字）" prop="key">
          <el-input
            v-model="form.key"
            placeholder="请设置支付密码"
            maxlength="6"
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
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import type { FormInstance } from "element-plus";
import { userRegister } from "@/api/user";

const router = useRouter();
const formRef = ref<FormInstance>();
const loading = ref(false);
const form = reactive({ user_name: "", nick_name: "", password: "", key: "" });

const rules = {
  user_name: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [
    { required: true, message: "请输入密码", trigger: "blur" },
    { min: 6, message: "密码至少6位", trigger: "blur" },
  ],
  key: [
    { required: true, message: "请设置支付密码", trigger: "blur" },
    { len: 6, message: "支付密码为6位", trigger: "blur" },
  ],
};

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
</style>
