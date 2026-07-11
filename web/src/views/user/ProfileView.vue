<template>
  <el-card>
    <template #header>个人资料</template>

    <div class="avatar-section">
      <el-avatar :size="80" :src="form.avatar" />
      <el-upload
        :show-file-list="false"
        accept="image/*"
        :before-upload="handleAvatarUpload"
      >
        <el-button size="small" :loading="uploading" style="margin-top: 8px"
          >更换头像</el-button
        >
      </el-upload>
    </div>

    <el-form
      :model="form"
      label-width="80px"
      style="max-width: 400px; margin-top: 20px"
    >
      <el-form-item label="用户名">
        <el-input :value="form.user_name" disabled />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="form.nick_name" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input :value="form.email || '未绑定'" disabled>
          <template #append>
            <el-button @click="openEmailDialog">
              {{ form.email ? "更换" : "绑定" }}
            </el-button>
          </template>
        </el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >保存修改</el-button
        >
      </el-form-item>
    </el-form>

    <el-dialog v-model="emailDialogVisible" title="绑定邮箱" width="420px">
      <el-form :model="emailForm" label-width="86px">
        <el-form-item label="新邮箱">
          <el-input v-model="emailForm.email" placeholder="请输入邮箱" />
        </el-form-item>
        <el-form-item label="登录密码">
          <el-input
            v-model="emailForm.password"
            type="password"
            show-password
            placeholder="用于生成邮箱验证链接"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="emailDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="sendingEmail" @click="handleSendEmail"
          >发送验证邮件</el-button
        >
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import {
  getUserInfo,
  sendEmail,
  updateUserInfo,
  uploadAvatar,
} from "@/api/user";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const saving = ref(false);
const uploading = ref(false);
const emailDialogVisible = ref(false);
const sendingEmail = ref(false);
const form = reactive({
  user_name: "",
  nick_name: "",
  email: "",
  avatar: "",
});
const emailForm = reactive({ email: "", password: "" });

async function loadProfile() {
  try {
    const res: any = await getUserInfo();
    const u = res.data;
    form.user_name = u.user_name;
    form.nick_name = u.nick_name || u.nickname || u.user_name;
    form.email = u.email;
    form.avatar = u.avatar;
  } catch {}
}

async function handleSave() {
  if (!form.nick_name.trim()) return ElMessage.warning("昵称不能为空");
  saving.value = true;
  try {
    await updateUserInfo({ nick_name: form.nick_name.trim() });
    const res: any = await getUserInfo();
    userStore.setUserInfo(res.data);
    form.nick_name = res.data.nick_name || res.data.nickname || res.data.user_name;
    ElMessage.success("保存成功");
  } finally {
    saving.value = false;
  }
}

async function handleAvatarUpload(file: File) {
  if (!file.type.startsWith("image/")) {
    ElMessage.warning("请选择图片文件");
    return false;
  }
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning("头像大小不能超过 2MB");
    return false;
  }
  const fd = new FormData();
  fd.append("file", file);
  uploading.value = true;
  try {
    await uploadAvatar(fd);
    const res: any = await getUserInfo();
    form.avatar = res.data.avatar;
    userStore.setUserInfo(res.data);
    ElMessage.success("头像已更新");
  } catch {}
  finally {
    uploading.value = false;
  }
  return false; // 阻止自动上传
}

function openEmailDialog() {
  emailForm.email = form.email;
  emailForm.password = "";
  emailDialogVisible.value = true;
}

async function handleSendEmail() {
  if (!emailForm.email || !emailForm.password) {
    return ElMessage.warning("请填写邮箱和登录密码");
  }
  sendingEmail.value = true;
  try {
    await sendEmail({
      email: emailForm.email,
      password: emailForm.password,
      operation_type: 1,
    });
    emailDialogVisible.value = false;
    ElMessage.success("验证邮件已发送，请前往邮箱确认");
  } finally {
    sendingEmail.value = false;
  }
}

onMounted(loadProfile);
</script>

<style scoped>
.avatar-section {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}
</style>
