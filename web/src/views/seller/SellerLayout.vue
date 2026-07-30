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
        <el-menu-item index="/seller/apply">{{ $t("sellerCenter.menu.onboarding") }}</el-menu-item>
        <el-menu-item index="/seller/account" :disabled="!sellerStore.isApproved">
          {{ $t("sellerCenter.menu.account") }}
        </el-menu-item>
        <el-menu-item index="/seller/products">{{ $t("sellerCenter.menu.products") }}</el-menu-item>
        <el-menu-item index="/seller/orders" :disabled="!sellerStore.isApproved">
          {{ $t("sellerCenter.menu.orders") }}
        </el-menu-item>
        <el-menu-item index="/seller/products/new" :disabled="!sellerStore.isApproved">
          {{ $t("sellerCenter.menu.publish") }}
        </el-menu-item>
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
import { useSellerStore } from "@/stores/seller";

const route = useRoute();
const sellerStore = useSellerStore();
const activeMenu = computed(() => {
  if (route.path.startsWith("/seller/account")) return "/seller/account";
  if (route.path.startsWith("/seller/products/new")) {
    return "/seller/products/new";
  }
  if (route.path.startsWith("/seller/orders")) return "/seller/orders";
  if (route.path.startsWith("/seller/products")) return "/seller/products";
  return "/seller/apply";
});

sellerStore.loadProfile({ silentError: true });
</script>
