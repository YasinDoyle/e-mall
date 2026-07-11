<template>
  <div v-if="total > 0" class="pagination-wrap">
    <el-pagination
      :current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      background
      @update:current-page="$emit('update:page', $event)"
      @current-change="$emit('change', $event)"
    />
  </div>
</template>

<script setup lang="ts">
defineProps<{
  page: number;
  pageSize: number;
  total: number;
}>();

defineEmits<{
  (e: "update:page", value: number): void;
  (e: "change", value: number): void;
}>();
</script>

<style scoped>
.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
