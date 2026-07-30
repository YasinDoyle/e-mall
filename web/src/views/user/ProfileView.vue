<template>
  <el-card>
    <template #header>{{ t("profile.title") }}</template>

    <div class="avatar-section">
      <el-avatar :size="80" :src="form.avatar" />
      <el-upload
        :show-file-list="false"
        accept="image/*"
        :before-upload="handleAvatarUpload"
      >
        <el-button size="small" :loading="uploading" style="margin-top: 8px"
          >{{ t("profile.changeAvatar") }}</el-button
        >
      </el-upload>
    </div>

    <el-form
      :model="form"
      label-width="80px"
      style="max-width: 400px; margin-top: 20px"
    >
      <el-form-item :label="t('profile.username')">
        <el-input :value="form.user_name" disabled />
      </el-form-item>
      <el-form-item :label="t('profile.nickname')">
        <el-input v-model="form.nick_name" />
      </el-form-item>
      <el-form-item :label="t('profile.email')">
        <el-input :value="form.email || t('profile.unbound')" disabled />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >{{ t("profile.saveChanges") }}</el-button
        >
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import { getUserInfo, updateUserInfo, uploadAvatar } from "@/api/user";
import { useUserStore } from "@/stores/user";

const userStore = useUserStore();
const { t } = useI18n();
const saving = ref(false);
const uploading = ref(false);
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
    form.nick_name = u.nick_name || u.nickname || u.user_name;
    form.email = u.email;
    form.avatar = u.avatar;
  } catch {}
}

async function handleSave() {
  if (!form.nick_name.trim()) return ElMessage.warning(t("profile.nicknameRequired"));
  saving.value = true;
  try {
    await updateUserInfo({ nick_name: form.nick_name.trim() });
    const res: any = await getUserInfo();
    userStore.setUserInfo(res.data);
    form.nick_name = res.data.nick_name || res.data.nickname || res.data.user_name;
    ElMessage.success(t("profile.saveSuccess"));
  } finally {
    saving.value = false;
  }
}

async function handleAvatarUpload(file: File) {
  if (!file.type.startsWith("image/")) {
    ElMessage.warning(t("profile.selectImage"));
    return false;
  }
  if (file.size > 2 * 1024 * 1024) {
    ElMessage.warning(t("profile.avatarTooLarge"));
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
    ElMessage.success(t("profile.avatarUpdated"));
  } catch {}
  finally {
    uploading.value = false;
  }
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
