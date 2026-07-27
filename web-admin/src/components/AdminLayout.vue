<template>
  <el-container style="min-height: 100vh">
    <el-aside width="200px" style="background: #001529">
      <div class="logo">🛒 Admin</div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="#001529"
        text-color="#c0c4cc"
        active-text-color="#409eff"
      >
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon>数据概览
        </el-menu-item>
        <el-menu-item index="/product">
          <el-icon><Goods /></el-icon>商品管理
        </el-menu-item>
        <el-menu-item index="/seller">
          <el-icon><Shop /></el-icon>商家审核
        </el-menu-item>
        <el-menu-item index="/order">
          <el-icon><Tickets /></el-icon>订单管理
        </el-menu-item>
        <el-menu-item index="/settlement">
          <el-icon><Money /></el-icon>结算管理
        </el-menu-item>
        <el-menu-item index="/withdraw">
          <el-icon><Wallet /></el-icon>提现审核
        </el-menu-item>
        <el-menu-item index="/category">
          <el-icon><Grid /></el-icon>分类管理
        </el-menu-item>
        <el-menu-item index="/coupon">
          <el-icon><Discount /></el-icon>优惠券
        </el-menu-item>
        <el-menu-item index="/flash-sale">
          <el-icon><Timer /></el-icon>秒杀管理
        </el-menu-item>
        <el-menu-item index="/carousel">
          <el-icon><Picture /></el-icon>轮播图
        </el-menu-item>
        <el-menu-item index="/user">
          <el-icon><User /></el-icon>用户管理
        </el-menu-item>
        <el-menu-item index="/notice">
          <el-icon><Bell /></el-icon>公告管理
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header
        style="
          background: #fff;
          display: flex;
          align-items: center;
          justify-content: space-between;
          box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
        "
      >
        <span style="font-size: 16px; font-weight: 600">E-Mall 管理后台</span>
        <el-dropdown @command="handleCommand">
          <span style="cursor: pointer"
            >{{ store.adminInfo?.nick_name || store.adminInfo?.user_name }}
            <el-icon><ArrowDown /></el-icon
          ></span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>
      <el-main style="background: #f0f2f5">
        <RouterView />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
import { useAdminStore } from "@/stores/admin";
import { Money, Shop, Wallet } from "@element-plus/icons-vue";

const router = useRouter();
const store = useAdminStore();

function handleCommand(cmd: string) {
  if (cmd === "logout") {
    store.logout();
    router.push("/login");
  }
}
</script>

<style scoped>
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 18px;
  font-weight: bold;
}
</style>
