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
        <el-menu-item index="/user/profile">个人资料</el-menu-item>
        <el-menu-item index="/user/orders">我的订单</el-menu-item>
        <el-menu-item index="/user/addresses">收货地址</el-menu-item>
        <el-menu-item index="/user/favorites">我的收藏</el-menu-item>
        <el-menu-item index="/user/coupons">我的优惠券</el-menu-item>
        <el-menu-item index="/user/wallet">我的钱包</el-menu-item>
        <el-menu-item index="/seller/apply">卖家中心</el-menu-item>
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
