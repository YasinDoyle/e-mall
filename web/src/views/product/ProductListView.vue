<template>
  <div>
    <div class="filter-bar">
      <span
        :class="['cat-item', { active: !selectedCategory }]"
        @click="selectCategory(undefined)"
        >{{ t("productList.all") }}</span
      >
      <span
        v-for="cat in categories"
        :key="cat.id"
        :class="['cat-item', { active: selectedCategory === cat.id }]"
        @click="selectCategory(cat.id)"
        >{{ cat.category_name }}</span
      >
    </div>
    <div v-loading="loading" class="product-grid">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <el-empty v-if="!loading && !products.length" :description="t('productList.empty')" />
    <Pagination
      v-model:page="page"
      :page-size="pageSize"
      :total="total"
      @change="loadProducts"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import ProductCard from "@/components/common/ProductCard.vue";
import Pagination from "@/components/common/Pagination.vue";
import { getProductList, getCategoryList } from "@/api/product";
import type { Product, Category } from "@/types";
import { t } from "@/locales";

const route = useRoute();
const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const selectedCategory = ref<number | undefined>(undefined);
const page = ref(1);
const pageSize = 16;
const total = ref(0);
const loading = ref(false);

async function loadProducts() {
  loading.value = true;
  try {
    const res: any = await getProductList({
      page_num: page.value,
      page_size: pageSize,
      category_id: selectedCategory.value,
    });
    products.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

function selectCategory(id?: number) {
  selectedCategory.value = id;
  page.value = 1;
  loadProducts();
}

onMounted(async () => {
  if (route.query.category_id) {
    selectedCategory.value = Number(route.query.category_id);
  }
  const catRes: any = await getCategoryList();
  categories.value = catRes.data?.item ?? [];
  loadProducts();
});

watch(
  () => route.query.category_id,
  (val) => {
    selectedCategory.value = val ? Number(val) : undefined;
    page.value = 1;
    loadProducts();
  },
);
</script>

<style scoped>
.filter-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 16px;
}
.cat-item {
  padding: 4px 12px;
  border-radius: 16px;
  background: #f0f0f0;
  cursor: pointer;
  font-size: 14px;
}
.cat-item.active {
  background: #409eff;
  color: #fff;
}
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}
</style>
