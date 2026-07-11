<template>
  <div v-if="product" class="detail-wrap">
    <el-row :gutter="32">
      <!-- 左：图片区 -->
      <el-col :span="10">
        <el-image :src="mainImg" fit="contain" class="main-img" />
        <div class="thumb-list">
          <img
            v-for="(img, i) in allImgs"
            :key="i"
            :src="img"
            :class="['thumb', { active: mainImg === img }]"
            @click="mainImg = img"
          />
        </div>
      </el-col>

      <!-- 右：信息区 -->
      <el-col :span="14">
        <h2 class="product-name">{{ product.name }}</h2>
        <p class="product-title">{{ product.title }}</p>

        <div class="price-row">
          <span class="price"
            >¥{{ product.discount_price || product.price }}</span
          >
          <span v-if="product.discount_price" class="original"
            >¥{{ product.price }}</span
          >
        </div>

        <el-descriptions :column="2" border size="small" style="margin: 16px 0">
          <el-descriptions-item label="卖家">{{
            product.boss_name
          }}</el-descriptions-item>
          <el-descriptions-item label="库存"
            >{{ product.num }} 件</el-descriptions-item
          >
          <el-descriptions-item label="浏览量">{{
            product.view
          }}</el-descriptions-item>
          <el-descriptions-item label="分类">{{
            product.category_id
          }}</el-descriptions-item>
        </el-descriptions>

        <div class="num-row">
          <span>购买数量：</span>
          <el-input-number v-model="buyNum" :min="1" :max="product.num" />
        </div>

        <div class="action-row">
          <el-button
            type="primary"
            size="large"
            :loading="addingCart"
            @click="handleAddCart"
          >
            加入购物车
          </el-button>
          <el-button
            size="large"
            :type="isFavorite ? 'warning' : 'default'"
            :icon="isFavorite ? 'StarFilled' : 'Star'"
            :loading="togglingFav"
            @click="handleToggleFavorite"
          >
            {{ isFavorite ? "已收藏" : "收藏" }}
          </el-button>
        </div>

        <el-divider />
        <div class="product-info">
          <p style="white-space: pre-wrap">{{ product.info }}</p>
        </div>
      </el-col>
    </el-row>

    <!-- 评价列表占位 -->
    <el-card style="margin-top: 24px">
      <template #header>商品评价</template>
      <el-empty description="暂无评价" />
    </el-card>
  </div>

  <div v-else-if="loading" style="text-align: center; padding: 80px">
    <el-icon class="is-loading" size="48"><Loading /></el-icon>
  </div>

  <el-result v-else icon="warning" title="商品不存在或已下架">
    <template #extra>
      <el-button type="primary" @click="$router.push('/products')"
        >返回商品列表</el-button
      >
    </template>
  </el-result>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Loading } from "@element-plus/icons-vue";
import { getProductDetail, getProductImgs } from "@/api/product";
import { createCart } from "@/api/cart";
import {
  createFavorite,
  deleteFavorite,
  getFavoriteList,
} from "@/api/favorite";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const product = ref<any>(null);
const extraImgs = ref<string[]>([]);
const mainImg = ref("");
const loading = ref(true);
const buyNum = ref(1);
const addingCart = ref(false);
const isFavorite = ref(false);
const togglingFav = ref(false);

const allImgs = computed(() => {
  const cover = product.value?.img_path ? [product.value.img_path] : [];
  return [...cover, ...extraImgs.value];
});

async function loadProduct() {
  const id = Number(route.params.id);
  try {
    const [detailRes, imgRes]: any[] = await Promise.all([
      getProductDetail({ id }),
      getProductImgs({ id }),
    ]);
    product.value = detailRes.data;
    mainImg.value = detailRes.data?.img_path ?? "";
    extraImgs.value = (imgRes.data?.item ?? []).map((i: any) => i.img_path);

    if (userStore.isLoggedIn) {
      checkFavorite(id);
    }
  } catch {
    product.value = null;
  } finally {
    loading.value = false;
  }
}

async function checkFavorite(productId: number) {
  try {
    const res: any = await getFavoriteList({ page_num: 1, page_size: 100 });
    const list = res.data?.item ?? [];
    isFavorite.value = list.some((f: any) => f.product_id === productId);
  } catch {}
}

async function handleAddCart() {
  if (!userStore.isLoggedIn) return router.push("/login");
  addingCart.value = true;
  try {
    await createCart({
      product_id: product.value.id,
      boss_id: product.value.boss_id,
      num: buyNum.value,
      max_num: product.value.num,
    });
    userStore.setCartCount(userStore.cartCount + 1);
    ElMessage.success("已加入购物车");
  } finally {
    addingCart.value = false;
  }
}

async function handleToggleFavorite() {
  if (!userStore.isLoggedIn) return router.push("/login");
  togglingFav.value = true;
  try {
    if (isFavorite.value) {
      await deleteFavorite({ product_id: product.value.id });
      isFavorite.value = false;
      ElMessage.success("已取消收藏");
    } else {
      await createFavorite({
        product_id: product.value.id,
        boss_id: product.value.boss_id,
      });
      isFavorite.value = true;
      ElMessage.success("收藏成功");
    }
  } finally {
    togglingFav.value = false;
  }
}

onMounted(loadProduct);
</script>

<style scoped>
.detail-wrap {
  max-width: 1100px;
  margin: 0 auto;
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}
.main-img {
  width: 100%;
  height: 380px;
  border-radius: 8px;
  border: 1px solid #eee;
}
.thumb-list {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  flex-wrap: wrap;
}
.thumb {
  width: 64px;
  height: 64px;
  object-fit: cover;
  border-radius: 4px;
  border: 2px solid transparent;
  cursor: pointer;
}
.thumb.active {
  border-color: #409eff;
}
.product-name {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 8px;
}
.product-title {
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}
.price-row {
  margin-bottom: 8px;
}
.price {
  font-size: 28px;
  font-weight: bold;
  color: #f56c6c;
}
.original {
  margin-left: 10px;
  color: #999;
  font-size: 16px;
  text-decoration: line-through;
}
.num-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 16px 0;
}
.action-row {
  display: flex;
  gap: 12px;
}
.product-info {
  color: #555;
  font-size: 14px;
  line-height: 1.8;
}
</style>
