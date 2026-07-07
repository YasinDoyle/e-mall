<template>
  <el-card class="product-card" @click="$router.push(`/product/${product.id}`)">
    <img :src="product.img_path" class="product-img" />
    <div class="product-name">{{ product.name }}</div>
    <div class="product-price">
      <span class="price">¥{{ product.discount_price || product.price }}</span>
      <span v-if="product.discount_price" class="original"
        >¥{{ product.price }}</span
      >
    </div>
    <div class="product-seller">{{ product.boss_name }}</div>
  </el-card>
</template>

<script setup lang="ts">
import type { Product } from "@/types";

defineProps<{ product: Product }>();
</script>

<style scoped>
.product-card {
  cursor: pointer;
  transition: box-shadow 0.2s;
}
.product-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
.product-img {
  width: 100%;
  height: 160px;
  object-fit: cover;
  border-radius: 4px;
}
.product-name {
  margin-top: 8px;
  font-size: 14px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.product-price {
  margin-top: 6px;
}
.price {
  color: #f56c6c;
  font-weight: bold;
  font-size: 16px;
}
.original {
  margin-left: 6px;
  color: #999;
  font-size: 12px;
  text-decoration: line-through;
}
.product-seller {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}
</style>
