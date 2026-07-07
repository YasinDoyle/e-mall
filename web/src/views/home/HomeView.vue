<template>
  <div>
    <!-- 轮播图 -->
    <el-carousel height="320px" class="banner">
      <el-carousel-item v-for="item in carousels" :key="item.id">
        <img
          :src="item.img_path"
          style="width: 100%; height: 100%; object-fit: cover"
        />
      </el-carousel-item>
    </el-carousel>

    <!-- 分类导航 -->
    <div class="section-title">商品分类</div>
    <div class="category-list">
      <el-tag
        v-for="cat in categories"
        :key="cat.id"
        class="category-tag"
        :type="selectedCategory === cat.id ? '' : 'info'"
        style="cursor: pointer"
        @click="filterByCategory(cat.id)"
      >
        {{ cat.category_name }}
      </el-tag>
    </div>

    <!-- 商品列表 -->
    <div class="section-title">
      热门商品
      <el-button link type="primary" @click="$router.push('/products')"
        >查看全部</el-button
      >
    </div>
    <div class="product-grid">
      <ProductCard
        v-for="product in products"
        :key="product.id"
        :product="product"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import ProductCard from "@/components/common/ProductCard.vue";
import { getCarousels, getCategoryList, getProductList } from "@/api/product";
import type { Carousel, Category, Product } from "@/types";

const router = useRouter();
const carousels = ref<Carousel[]>([]);
const categories = ref<Category[]>([]);
const products = ref<Product[]>([]);
const selectedCategory = ref<number | undefined>(undefined);

async function loadData() {
  try {
    const [carouselRes, categoryRes, productRes]: any[] = await Promise.all([
      getCarousels(),
      getCategoryList(),
      getProductList({ page_num: 1, page_size: 12 }),
    ]);
    carousels.value = carouselRes.data?.item ?? [];
    categories.value = categoryRes.data?.item ?? [];
    products.value = productRes.data?.item ?? [];
  } catch (err) {
    console.error("首页数据加载失败，请确认后端服务是否启动：", err);
  }
}

function filterByCategory(id: number) {
  selectedCategory.value = id;
  router.push({ path: "/products", query: { category_id: id } });
}

onMounted(loadData);
</script>

<style scoped>
.banner {
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 24px;
}
.section-title {
  font-size: 18px;
  font-weight: 600;
  margin: 16px 0 12px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 16px;
}
.category-tag {
  font-size: 14px;
  padding: 6px 12px;
}
.product-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}
</style>
