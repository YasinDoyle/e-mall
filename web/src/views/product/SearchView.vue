<template>
  <div>
    <div class="search-bar">
      <el-input
        v-model="keyword"
        placeholder="搜索商品"
        @keyup.enter="doSearch(true)"
      >
        <template #append
          ><el-button @click="doSearch(true)">搜索</el-button></template
        >
      </el-input>
    </div>
    <div v-loading="loading" class="product-grid">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <el-empty
      v-if="!loading && !products.length && searched"
      description="未找到相关商品"
    />
    <Pagination
      v-model:page="page"
      :page-size="pageSize"
      :total="total"
      @change="doSearch(false)"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import ProductCard from "@/components/common/ProductCard.vue";
import Pagination from "@/components/common/Pagination.vue";
import { searchProducts } from "@/api/product";
import type { Product } from "@/types";

const route = useRoute();
const router = useRouter();
const keyword = ref((route.query.info as string) || "");
const products = ref<Product[]>([]);
const searched = ref(false);
const loading = ref(false);
const page = ref(1);
const pageSize = 20;
const total = ref(0);

async function doSearch(resetPage = false) {
  if (!keyword.value.trim()) return;
  if (resetPage) page.value = 1;
  router.replace({ query: { info: keyword.value, page: page.value } });
  loading.value = true;
  try {
    const res: any = await searchProducts({
      info: keyword.value,
      page_num: page.value,
      page_size: pageSize,
    });
    products.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
    searched.value = true;
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  if (route.query.page) page.value = Number(route.query.page) || 1;
  if (keyword.value) doSearch(false);
});

watch(
  () => route.query.info,
  (val) => {
    const next = (val as string) || "";
    if (next && next !== keyword.value) {
      keyword.value = next;
      page.value = 1;
      doSearch(false);
    }
  },
);
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
