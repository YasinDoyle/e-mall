<template>
  <el-container class="layout">
    <el-header height="60px">
      <NavBar />
    </el-header>
    <el-main>
      <RouterView />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import NavBar from "@/components/common/NavBar.vue";
</script>

<style scoped>
.layout {
  min-height: 100vh;
}
.el-header {
  padding: 0;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}
.el-main {
  background: #f5f5f5;
  padding: 20px;
}
</style>
