<template>
  <el-card v-if="order">
    <template #header>
      <div
        style="
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <span>{{ t("orderDetail.title") }}</span>
        <el-tag :type="statusTagType(order.type)">{{
          statusText(order.type)
        }}</el-tag>
      </div>
    </template>

    <!-- 物流时间轴 -->
    <el-timeline style="margin-bottom: 20px">
      <el-timeline-item
        v-for="step in timeline"
        :key="step.label"
        :type="step.done ? 'primary' : 'info'"
        :hollow="!step.done"
      >
        {{ step.label }}
        <span v-if="step.desc" class="timeline-desc">{{ step.desc }}</span>
      </el-timeline-item>
    </el-timeline>

    <div class="product-row" @click="$router.push(`/product/${order.product_id}`)">
      <img :src="order.img_path" class="product-img" />
      <div class="product-info">
        <div class="product-name">{{ order.name || t("orderDetail.product") }}</div>
        <div class="product-meta">{{ t("orderDetail.productId", { id: order.product_id }) }}</div>
      </div>
      <div class="product-price">¥{{ totalAmount(order) }}</div>
    </div>

    <el-descriptions :column="2" border>
      <el-descriptions-item :label="t('orderDetail.orderNo')">{{
        order.order_num
      }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.trackingNo')">{{
        order.tracking_no || t("orderDetail.none")
      }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.quantity', { count: '' })"
        >{{ order.num }} 件</el-descriptions-item
      >
      <el-descriptions-item :label="t('orderDetail.amount')">
        <span style="color: #f56c6c; font-weight: bold">
          ¥{{ totalAmount(order) }}
        </span>
      </el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.recipient')">
        {{ order.address_name || "-" }} {{ order.address_phone || "" }}
      </el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.address')">
        {{ order.address || "-" }}
      </el-descriptions-item>
      <el-descriptions-item v-if="order.refund_status" :label="t('orderDetail.afterSaleStatus')">
        {{ refundText(order.refund_status) }}
      </el-descriptions-item>
      <el-descriptions-item v-if="order.refund_reason" :label="t('orderDetail.refundReason')">
        {{ order.refund_reason }}
      </el-descriptions-item>
    </el-descriptions>

    <div class="order-actions">
      <el-button v-if="order.type === 3" type="primary" @click="handleReceive"
        >{{ t("orderDetail.confirmReceive") }}</el-button
      >
      <el-button
        v-if="order.type === 4"
        @click="openReviewDialog"
        >{{ t("orderDetail.writeReview") }}</el-button
      >
      <el-button @click="$router.back()">{{ t("orderDetail.back") }}</el-button>
    </div>

    <el-dialog v-model="reviewDialogVisible" :title="t('orderDetail.reviewTitle')" width="520px">
      <el-form label-width="70px">
        <el-form-item :label="t('orderDetail.rating')">
          <el-rate v-model="reviewForm.rating" />
        </el-form-item>
        <el-form-item :label="t('orderDetail.review')">
          <el-input
            v-model="reviewForm.content"
            type="textarea"
            :rows="4"
            maxlength="500"
            show-word-limit
            :placeholder="t('orderDetail.reviewPlaceholder')"
          />
        </el-form-item>
        <el-form-item :label="t('orderDetail.image')">
          <el-upload
            :show-file-list="false"
            accept="image/*"
            :before-upload="handleReviewImageUpload"
          >
            <el-button :loading="uploadingImage">{{ t("orderDetail.uploadImage") }}</el-button>
          </el-upload>
          <div v-if="reviewForm.images.length" class="review-upload-list">
            <div
              v-for="img in reviewForm.images"
              :key="img"
              class="review-upload-item"
            >
              <el-image :src="img" fit="cover" />
              <el-button link type="danger" @click="removeReviewImage(img)"
                >{{ t("orderDetail.delete") }}</el-button
              >
            </div>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">{{ t("common.cancel") }}</el-button>
        <el-button type="primary" :loading="submittingReview" @click="submitReview"
          >{{ t("orderDetail.submitReview") }}</el-button
        >
      </template>
    </el-dialog>
  </el-card>

  <div v-else-if="loading" style="text-align: center; padding: 60px">
    <el-icon class="is-loading" size="40"><Loading /></el-icon>
  </div>

  <el-empty v-else :description="t('orderDetail.orderNotFound')">
    <el-button @click="$router.push('/user/orders')">{{ t("orderDetail.backToOrders") }}</el-button>
  </el-empty>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { ElMessage } from "element-plus";
import { Loading } from "@element-plus/icons-vue";
import { getOrderDetail, receiveOrder } from "@/api/order";
import { createReview, uploadReviewImage } from "@/api/review";
import { orderStatusText, refundStatusText } from "@/utils/status-labels";
import { t } from "@/locales";

const route = useRoute();
const order = ref<any>(null);
const loading = ref(false);
const reviewDialogVisible = ref(false);
const submittingReview = ref(false);
const uploadingImage = ref(false);
const reviewForm = ref({
  rating: 5,
  content: "",
  images: [] as string[],
});

const statusText = orderStatusText;
const statusTagType = (type: number) =>
  (({ 1: "warning", 2: "primary", 3: "warning", 4: "success", 5: "danger", 6: "info" })[
    type
  ] ?? "info") as any;
const refundText = refundStatusText;

const timeline = computed(() => {
  const type = order.value?.type ?? 0;
  return [
    { label: t("orderTimeline.submitted", "提交订单"), done: type >= 1 },
    { label: t("orderTimeline.paid", "支付成功"), done: type >= 2 },
    {
      label: t("orderTimeline.shipped", "商家发货"),
      done: type >= 3,
      desc: order.value?.tracking_no ? `物流单号：${order.value.tracking_no}` : "",
    },
    { label: t("orderTimeline.received", "确认收货"), done: type >= 4 },
  ];
});

function totalAmount(row: any) {
  return (Number(row.money || 0) * Number(row.num || 0)).toFixed(2);
}

async function loadOrder() {
  loading.value = true;
  try {
    const res: any = await getOrderDetail({
      order_id: Number(route.params.id),
    });
    order.value = res.data;
  } catch {
    order.value = null;
  } finally {
    loading.value = false;
  }
}

async function handleReceive() {
  try {
    await receiveOrder({ order_id: order.value.id });
    ElMessage.success(t("orderDetail.receivedSuccess"));
    loadOrder();
  } catch {}
}

function openReviewDialog() {
  reviewForm.value = {
    rating: 5,
    content: "",
    images: [],
  };
  reviewDialogVisible.value = true;
}

async function handleReviewImageUpload(file: File) {
  if (!file.type.startsWith("image/")) {
    ElMessage.warning(t("orderDetail.chooseImage"));
    return false;
  }
  if (file.size > 3 * 1024 * 1024) {
    ElMessage.warning(t("orderDetail.reviewImageTooLarge"));
    return false;
  }
  if (reviewForm.value.images.length >= 3) {
    ElMessage.warning(t("orderDetail.reviewImageLimit"));
    return false;
  }

  const formData = new FormData();
  formData.append("file", file);
  uploadingImage.value = true;
  try {
    const res: any = await uploadReviewImage(formData);
    if (res.data?.url) {
      reviewForm.value.images.push(res.data.url);
    }
  } finally {
    uploadingImage.value = false;
  }
  return false;
}

function removeReviewImage(img: string) {
  reviewForm.value.images = reviewForm.value.images.filter((item) => item !== img);
}

async function submitReview() {
  if (!reviewForm.value.rating) return ElMessage.warning(t("orderDetail.chooseRating"));
  submittingReview.value = true;
  try {
    await createReview({
      product_id: order.value.product_id,
      order_id: order.value.id,
      rating: reviewForm.value.rating,
      content: reviewForm.value.content.trim(),
      images: reviewForm.value.images.join(","),
    });
    ElMessage.success(t("orderDetail.reviewSuccess"));
    reviewDialogVisible.value = false;
  } finally {
    submittingReview.value = false;
  }
}

onMounted(loadOrder);
</script>

<style scoped>
.product-row {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border: 1px solid #eee;
  border-radius: 6px;
  margin-bottom: 16px;
  cursor: pointer;
}
.product-img {
  width: 72px;
  height: 72px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}
.product-info {
  flex: 1;
  min-width: 0;
}
.product-name {
  font-weight: 500;
}
.product-meta,
.timeline-desc {
  margin-left: 8px;
  color: #999;
  font-size: 12px;
}
.product-price {
  color: #f56c6c;
  font-weight: 600;
}
.order-actions {
  display: flex;
  gap: 10px;
  margin-top: 20px;
  justify-content: flex-end;
}
.review-upload-list {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  margin-top: 10px;
}
.review-upload-item {
  width: 84px;
}
.review-upload-item .el-image {
  width: 84px;
  height: 84px;
  border-radius: 4px;
  display: block;
}
</style>
