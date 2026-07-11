<template>
  <div class="valid-page">
    <el-card style="width: 400px">
      <el-result
        v-if="status === 'success'"
        icon="success"
        title="邮箱验证成功"
        sub-title="您的账号已激活"
      >
        <template #extra>
          <el-button type="primary" @click="$router.push(successPath)"
            >{{ successButtonText }}</el-button
          >
        </template>
      </el-result>
      <el-result
        v-else-if="status === 'error'"
        icon="error"
        title="验证失败"
        :sub-title="errorMsg"
      />
      <div v-else class="loading-wrap">
        <el-icon class="is-loading" size="40"><Loading /></el-icon>
        <p>验证中...</p>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { Loading } from "@element-plus/icons-vue";
import { getUserInfo, validEmail } from "@/api/user";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const userStore = useUserStore();
const status = ref<"pending" | "success" | "error">("pending");
const errorMsg = ref("");
const successPath = ref("/login");
const successButtonText = ref("去登录");

onMounted(async () => {
  const token = route.query.token as string;
  if (!token) {
    status.value = "error";
    errorMsg.value = "无效的验证链接";
    return;
  }
  try {
    const res: any = await validEmail({ token });
    if (userStore.isLoggedIn) {
      const infoRes: any = await getUserInfo().catch(() => null);
      userStore.setUserInfo(infoRes?.data ?? res.data);
      successPath.value = "/user/profile";
      successButtonText.value = "返回个人资料";
    }
    status.value = "success";
  } catch {
    status.value = "error";
    errorMsg.value = "验证链接已过期或无效";
  }
});
</script>

<style scoped>
.valid-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
}
.loading-wrap {
  text-align: center;
  padding: 40px;
}
</style>
