<template>
  <el-container>
    <el-aside
      width="200px"
      style="background: #fff; border-right: 1px solid #e8e8e8"
    >
      <el-menu
        :default-active="activeMenu"
        router
        style="border: none; height: 100%"
      >
        <el-menu-item index="/user/profile">{{ $t("userMenu.profile") }}</el-menu-item>
        <el-menu-item index="/user/orders">{{ $t("userMenu.orders") }}</el-menu-item>
        <el-menu-item index="/user/addresses">{{ $t("userMenu.addresses") }}</el-menu-item>
        <el-menu-item index="/user/favorites">{{ $t("userMenu.favorites") }}</el-menu-item>
        <el-menu-item index="/user/coupons">{{ $t("userMenu.coupons") }}</el-menu-item>
        <el-menu-item index="/user/wallet">{{ $t("userMenu.wallet") }}</el-menu-item>
        <el-menu-item index="/user/notifications">{{ $t("userMenu.notifications") }}</el-menu-item>
        <el-menu-item index="/seller/apply">{{ $t("userMenu.sellerCenter") }}</el-menu-item>
      </el-menu>
    </el-aside>
    <el-main>
      <RouterView />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const activeMenu = computed(() => {
  if (route.path.startsWith("/user/orders")) return "/user/orders";
  return route.path;
});
</script>
