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
        <el-button size="small" style="margin-top: 8px">更换头像</el-button>
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
        <el-input v-model="form.email" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >保存修改</el-button
        >
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getUserInfo, updateUserInfo, uploadAvatar } from "@/api/user";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const saving = ref(false);
const form = reactive({
  user_name: "",
  nick_name: "",
  email: "",
  avatar: "",
});

async function loadProfile() {
  try {
    const res: any = await getUserInfo();
    const u = res.data;
    form.user_name = u.user_name;
    form.nick_name = u.nick_name;
    form.email = u.email;
    form.avatar = u.avatar;
  } catch {}
}

async function handleSave() {
  saving.value = true;
  try {
    await updateUserInfo({ nick_name: form.nick_name, email: form.email });
    const res: any = await getUserInfo();
    userStore.setUserInfo(res.data);
    ElMessage.success("保存成功");
  } finally {
    saving.value = false;
  }
}

async function handleAvatarUpload(file: File) {
  const fd = new FormData();
  fd.append("file", file);
  try {
    await uploadAvatar(fd);
    const res: any = await getUserInfo();
    form.avatar = res.data.avatar;
    userStore.setUserInfo(res.data);
    ElMessage.success("头像已更新");
  } catch {}
  return false; // 阻止自动上传
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
