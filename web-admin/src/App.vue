<template>
  <el-config-provider :locale="elementPlusLocale">
    <RouterView />
  </el-config-provider>
</template>

<script setup lang="ts">
import { elementPlusLocale } from "@/locales";
</script>
