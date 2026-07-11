<template>
  <div class="navbar">
    <div class="navbar-left">
      <RouterLink to="/" class="logo">🛒 E-Mall</RouterLink>
      <el-input
        v-model="searchKeyword"
        placeholder="搜索商品..."
        style="width: 300px; margin-left: 20px"
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="Search" @click="handleSearch" />
        </template>
      </el-input>
    </div>

    <div class="navbar-right">
      <RouterLink to="/flash-sale" class="nav-link">秒杀专场</RouterLink>

      <template v-if="userStore.isLoggedIn">
        <RouterLink to="/cart" class="nav-link">
          <el-badge :value="userStore.cartCount" :hidden="!userStore.cartCount">
            <el-icon size="20"><ShoppingCart /></el-icon>
          </el-badge>
        </RouterLink>

        <el-dropdown @command="handleCommand">
          <div class="user-avatar">
            <el-avatar :size="32" :src="userStore.userInfo?.avatar || ''" />
            <span style="margin-left: 6px">{{
              userStore.userInfo?.nick_name || userStore.userInfo?.user_name
            }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="orders">我的订单</el-dropdown-item>
              <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
              <el-dropdown-item command="wallet">我的钱包</el-dropdown-item>
              <el-dropdown-item divided command="logout"
                >退出登录</el-dropdown-item
              >
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>

      <template v-else>
        <RouterLink to="/login">
          <el-button type="primary" size="small">登录</el-button>
        </RouterLink>
        <RouterLink to="/register" style="margin-left: 8px">
          <el-button size="small">注册</el-button>
        </RouterLink>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { Search, ShoppingCart } from "@element-plus/icons-vue";
import { useUserStore } from "@/stores/user";
import { getCartList } from "@/api/cart";
import { getUserInfo } from "@/api/user";

const router = useRouter();
const userStore = useUserStore();
const searchKeyword = ref("");

function handleSearch() {
  if (searchKeyword.value.trim()) {
    router.push({
      path: "/search",
      query: { info: searchKeyword.value.trim() },
    });
  }
}

function handleCommand(command: string) {
  if (command === "logout") {
    userStore.logout();
    router.push("/login");
  } else {
    router.push(`/user/${command}`);
  }
}

async function syncUserSession() {
  if (!userStore.isLoggedIn) return;
  try {
    const [userRes, cartRes]: any[] = await Promise.all([
      getUserInfo(),
      getCartList(),
    ]);
    userStore.setUserInfo(userRes.data);
    userStore.setCartCount(cartRes.data?.item?.length ?? 0);
  } catch {
    userStore.logout();
    router.push("/login");
  }
}

onMounted(syncUserSession);

watch(
  () => userStore.token,
  () => {
    syncUserSession();
  },
);
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
  padding: 0 24px;
  background: #fff;
}
.navbar-left,
.navbar-right {
  display: flex;
  align-items: center;
  gap: 16px;
}
.logo {
  font-size: 20px;
  font-weight: bold;
  color: #409eff;
  text-decoration: none;
}
.nav-link {
  color: #333;
  text-decoration: none;
  font-size: 14px;
}
.nav-link:hover {
  color: #409eff;
}
.user-avatar {
  display: flex;
  align-items: center;
  cursor: pointer;
}
</style>
