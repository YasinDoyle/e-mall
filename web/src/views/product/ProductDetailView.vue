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
          <el-input-number
            v-model="buyNum"
            :min="1"
            :max="product.num"
            :disabled="isOwnProduct"
          />
        </div>

        <div class="action-row">
          <el-button
            type="primary"
            size="large"
            :loading="addingCart"
            :disabled="isOwnProduct"
            @click="handleAddCart"
          >
            {{ isOwnProduct ? "自己的商品" : "加入购物车" }}
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
        <el-descriptions :column="2" border size="small" style="margin-bottom: 16px">
          <el-descriptions-item label="品牌">
            {{ product.brand || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="产地">
            {{ product.origin || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="规格">
            {{ product.specification || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="生产日期">
            {{ product.production_date || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="保质期">
            {{ product.shelf_life || "-" }}
          </el-descriptions-item>
          <el-descriptions-item label="服务保障">
            {{ product.service_guarantees || "-" }}
          </el-descriptions-item>
        </el-descriptions>
        <div class="product-info">
          <p style="white-space: pre-wrap">{{ product.info }}</p>
        </div>
        <div v-if="product.certificates?.length" class="certificate-preview">
          <div class="section-subtitle">资质材料</div>
          <el-image
            v-for="certificate in product.certificates"
            :key="certificate.id"
            :src="certificate.file_path"
            :preview-src-list="certificateImages"
            fit="cover"
            class="certificate-img"
          >
            <template #error>
              <div class="certificate-error">{{ certificate.name }}</div>
            </template>
          </el-image>
        </div>
      </el-col>
    </el-row>

    <el-card style="margin-top: 24px">
      <template #header>
        <div class="review-header">
          <span>商品评价</span>
          <span class="review-total">共 {{ reviewTotal }} 条</span>
        </div>
      </template>
      <el-skeleton v-if="reviewLoading" :rows="3" animated />
      <el-empty v-else-if="!reviews.length" description="暂无评价" />
      <div v-else class="review-list">
        <div v-for="review in reviews" :key="review.id" class="review-item">
          <el-avatar :size="36" :src="review.user_avatar" />
          <div class="review-content">
            <div class="review-line">
              <span class="review-user">{{ review.user_name || "匿名用户" }}</span>
              <el-rate :model-value="review.rating" disabled size="small" />
            </div>
            <p>{{ review.content || "用户未填写评价内容" }}</p>
            <div v-if="reviewImageList(review).length" class="review-images">
              <el-image
                v-for="img in reviewImageList(review)"
                :key="img"
                :src="img"
                :preview-src-list="reviewImageList(review)"
                fit="cover"
                class="review-img"
              />
            </div>
            <div class="review-time">{{ formatTime(review.created_at) }}</div>
          </div>
        </div>
      </div>
      <el-pagination
        v-if="reviewTotal > reviewPageSize"
        v-model:current-page="reviewPage"
        :page-size="reviewPageSize"
        :total="reviewTotal"
        layout="prev, pager, next"
        style="justify-content: center; margin-top: 16px"
        @current-change="loadReviews"
      />
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
import { createCart, getCartList } from "@/api/cart";
import { getReviewList } from "@/api/review";
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
const favoriteId = ref<number | null>(null);
const togglingFav = ref(false);
const reviews = ref<any[]>([]);
const reviewLoading = ref(false);
const reviewPage = ref(1);
const reviewPageSize = 5;
const reviewTotal = ref(0);

const allImgs = computed(() => {
  const cover = product.value?.img_path ? [product.value.img_path] : [];
  return [...cover, ...extraImgs.value];
});

const certificateImages = computed(() =>
  (product.value?.certificates ?? [])
    .map((item: any) => item.file_path)
    .filter(Boolean),
);

const isOwnProduct = computed(
  () =>
    userStore.isLoggedIn &&
    product.value?.boss_id &&
    product.value.boss_id === userStore.userInfo?.id,
);

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
    loadReviews();
  } catch {
    product.value = null;
  } finally {
    loading.value = false;
  }
}

async function loadReviews() {
  const productId = Number(route.params.id);
  reviewLoading.value = true;
  try {
    const res: any = await getReviewList({
      product_id: productId,
      page_num: reviewPage.value,
      page_size: reviewPageSize,
    });
    reviews.value = res.data?.item ?? [];
    reviewTotal.value = res.data?.total ?? 0;
  } catch {
    reviews.value = [];
    reviewTotal.value = 0;
  } finally {
    reviewLoading.value = false;
  }
}

function reviewImageList(review: any) {
  return String(review.images || "")
    .split(",")
    .map((item) => item.trim())
    .filter(Boolean);
}

function formatTime(timestamp: number) {
  if (!timestamp) return "";
  return new Date(timestamp * 1000).toLocaleDateString();
}

async function checkFavorite(productId: number) {
  try {
    const res: any = await getFavoriteList({ page_num: 1, page_size: 100 });
    const list = res.data?.item ?? [];
    const favorite = list.find((f: any) => f.product_id === productId);
    isFavorite.value = !!favorite;
    favoriteId.value = favorite?.id ?? null;
  } catch {}
}

async function handleAddCart() {
  if (!userStore.isLoggedIn) return router.push("/login");
  if (isOwnProduct.value) return ElMessage.warning("不能购买自己发布的商品");
  addingCart.value = true;
  try {
    await createCart({
      product_id: product.value.id,
      boss_id: product.value.boss_id,
      num: buyNum.value,
      max_num: product.value.num,
    });
    const cartRes: any = await getCartList();
    userStore.setCartCount((cartRes.data?.item ?? []).length);
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
      if (!favoriteId.value) {
        await checkFavorite(product.value.id);
      }
      if (!favoriteId.value) return;
      await deleteFavorite({ id: favoriteId.value });
      isFavorite.value = false;
      favoriteId.value = null;
      ElMessage.success("已取消收藏");
    } else {
      await createFavorite({
        product_id: product.value.id,
        boss_id: product.value.boss_id,
      });
      await checkFavorite(product.value.id);
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
.section-subtitle {
  flex: 0 0 100%;
  margin: 16px 0 0;
  color: #303133;
  font-weight: 600;
}
.certificate-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.certificate-img,
.certificate-error {
  width: 92px;
  height: 92px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
.certificate-error {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 8px;
  color: #909399;
  font-size: 12px;
  text-align: center;
}
.review-header,
.review-line {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}
.review-total,
.review-time {
  color: #999;
  font-size: 12px;
}
.review-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.review-item {
  display: flex;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid #f2f2f2;
}
.review-content {
  flex: 1;
  min-width: 0;
}
.review-user {
  font-weight: 500;
}
.review-images {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin: 8px 0;
}
.review-img {
  width: 72px;
  height: 72px;
  border-radius: 4px;
}
</style>
