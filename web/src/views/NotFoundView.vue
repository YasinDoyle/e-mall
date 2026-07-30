<template>
  <div class="not-found">
    <el-result icon="warning" title="404" :sub-title="t('notFound.message')">
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">{{ t('notFound.backHome') }}</el-button>
        <el-button @click="$router.back()">{{ t('notFound.backPrevious') }}</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
import { t } from "@/locales";
</script>

<style scoped>
.not-found {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
