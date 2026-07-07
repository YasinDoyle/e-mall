<template>
  <div>
    <div class="search-bar">
      <el-input
        v-model="keyword"
        placeholder="搜索商品"
        @keyup.enter="doSearch"
      >
        <template #append
          ><el-button @click="doSearch">搜索</el-button></template
        >
      </el-input>
    </div>
    <div class="product-grid">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <el-empty
      v-if="!products.length && searched"
      description="未找到相关商品"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import ProductCard from "@/components/common/ProductCard.vue";
import { searchProducts } from "@/api/product";
import type { Product } from "@/types";

const route = useRoute();
const router = useRouter();
const keyword = ref((route.query.info as string) || "");
const products = ref<Product[]>([]);
const searched = ref(false);

async function doSearch() {
  if (!keyword.value.trim()) return;
  router.replace({ query: { info: keyword.value } });
  const res: any = await searchProducts({
    info: keyword.value,
    page_num: 1,
    page_size: 20,
  });
  products.value = res.data?.item ?? [];
  searched.value = true;
}

onMounted(() => {
  if (keyword.value) doSearch();
});
</script>

<style scoped>
.search-bar {
  margin-bottom: 20px;
  max-width: 500px;
}
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}
</style>
