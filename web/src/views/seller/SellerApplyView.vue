<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>商家入驻</span>
        <el-button :loading="loading" @click="loadProfile">刷新状态</el-button>
      </div>
    </template>

    <el-skeleton v-if="loading" :rows="4" animated />

    <template v-else>
      <el-alert
        v-if="!profile"
        title="当前账号还没有提交商家入驻申请"
        type="info"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 0"
        title="申请已提交，平台审核中"
        type="warning"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 1"
        title="商家能力已开通，可以发布商品"
        type="success"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 2"
        :title="`申请被拒绝：${profile.reject_reason || '未填写原因'}`"
        type="error"
        :closable="false"
        show-icon
      />
      <el-alert
        v-else-if="profile.status === 3"
        title="商家能力已被封禁，请联系平台处理"
        type="error"
        :closable="false"
        show-icon
      />

      <div v-if="profile" class="profile-box">
        <div class="profile-title">{{ profile.shop_name }}</div>
        <div class="profile-desc">{{ profile.description || "暂无店铺简介" }}</div>
        <div class="profile-meta">
          当前状态：
          <el-tag :type="statusTag(profile.status)">
            {{ profile.status_text }}
          </el-tag>
        </div>
      </div>

      <el-divider />

      <el-form
        v-if="canSubmit"
        label-width="90px"
        :model="form"
        class="apply-form"
      >
        <el-form-item label="店铺名称" required>
          <el-input
            v-model="form.shop_name"
            maxlength="80"
            show-word-limit
            placeholder="请输入店铺名称"
          />
        </el-form-item>
        <el-form-item label="店铺简介">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="5"
            maxlength="500"
            show-word-limit
            placeholder="介绍主营品类、服务承诺等"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="submit">
            {{ profile?.status === 2 ? "重新提交申请" : "提交入驻申请" }}
          </el-button>
        </el-form-item>
      </el-form>

      <div v-else-if="profile?.status === 1" class="actions">
        <el-button type="success" @click="$router.push('/seller/account')">
          资金账户
        </el-button>
        <el-button type="primary" @click="$router.push('/seller/products/new')">
          发布商品
        </el-button>
        <el-button @click="$router.push('/seller/products')">查看商品</el-button>
      </div>
    </template>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { applySeller } from "@/api/seller";
import { useSellerStore } from "@/stores/seller";

const loading = ref(false);
const submitting = ref(false);
const sellerStore = useSellerStore();
const profile = computed(() => sellerStore.profile);
const form = ref({ shop_name: "", description: "" });

const canSubmit = computed(
  () => !profile.value || profile.value.status === 2,
);

function statusTag(status: number) {
  return ({ 0: "warning", 1: "success", 2: "danger", 3: "info" } as any)[
    status
  ] ?? "info";
}

async function loadProfile() {
  loading.value = true;
  try {
    await sellerStore.loadProfile({ force: true, silentError: true });
    form.value = {
      shop_name: profile.value?.shop_name ?? "",
      description: profile.value?.description ?? "",
    };
  } finally {
    loading.value = false;
  }
}

async function submit() {
  if (!form.value.shop_name.trim()) {
    return ElMessage.warning("请输入店铺名称");
  }
  submitting.value = true;
  try {
    const res: any = await applySeller({
      shop_name: form.value.shop_name.trim(),
      description: form.value.description.trim(),
    });
    sellerStore.setProfile(res.data);
    ElMessage.success("申请已提交");
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  if (!sellerStore.loaded) {
    await sellerStore.loadProfile({ silentError: true });
  }
  form.value = {
    shop_name: profile.value?.shop_name ?? "",
    description: profile.value?.description ?? "",
  };
});
</script>

<style scoped>
.header,
.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.profile-box {
  margin-top: 16px;
  padding: 14px 16px;
  border: 1px solid #ebeef5;
  border-radius: 8px;
}
.profile-title {
  font-weight: 600;
  color: #303133;
}
.profile-desc,
.profile-meta {
  margin-top: 8px;
  color: #606266;
  font-size: 14px;
}
.apply-form {
  max-width: 640px;
}
</style>
