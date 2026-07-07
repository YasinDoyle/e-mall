<template>
  <div>
    <div class="filter-bar">
      <span
        v-for="cat in categories"
        :key="cat.id"
        :class="['cat-item', { active: selectedCategory === cat.id }]"
        @click="
          selectedCategory = cat.id;
          loadProducts();
        "
        >{{ cat.category_name }}</span
      >
      <span
        :class="['cat-item', { active: !selectedCategory }]"
        @click="
          selectedCategory = undefined;
          loadProducts();
        "
        >全部</span
      >
    </div>
    <div class="product-grid">
      <ProductCard v-for="p in products" :key="p.id" :product="p" />
    </div>
    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="prev, pager, next"
      style="margin-top: 20px; justify-content: center"
      @current-change="loadProducts"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useRoute } from "vue-router";
import ProductCard from "@/components/common/ProductCard.vue";
import { getProductList, getCategoryList } from "@/api/product";
import type { Product, Category } from "@/types";

const route = useRoute();
const products = ref<Product[]>([]);
const categories = ref<Category[]>([]);
const selectedCategory = ref<number | undefined>(undefined);
const page = ref(1);
const pageSize = 16;
const total = ref(0);

async function loadProducts() {
  const res: any = await getProductList({
    page_num: page.value,
    page_size: pageSize,
    category_id: selectedCategory.value,
  });
  products.value = res.data?.item ?? [];
  total.value = res.data?.total ?? 0;
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
