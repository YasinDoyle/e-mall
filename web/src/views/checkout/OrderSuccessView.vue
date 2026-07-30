<template>
  <div style="text-align: center; padding: 60px">
    <el-result
      icon="success"
      :title="t('orderSuccess.title')"
      :sub-title="t('orderSuccess.subtitle')"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/user/orders')">{{ t('orderSuccess.viewOrder') }}</el-button>
        <el-button @click="$router.push('/')">{{ t('orderSuccess.continueShopping') }}</el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup lang="ts">
import { t } from "@/locales";
</script>
