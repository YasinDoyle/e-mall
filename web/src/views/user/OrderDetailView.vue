<template>
  <el-card v-if="order">
    <template #header>
      <div class="header-row">
        <span>{{ t("orderDetail.title") }}</span>
        <el-tag :type="statusTagType(order.type)">{{ statusText(order.type) }}</el-tag>
      </div>
    </template>

    <el-descriptions :column="2" border class="detail-grid">
      <el-descriptions-item :label="t('orderDetail.orderNo')">{{ order.order_num }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.paymentChannel')">{{ order.payment_channel ? paymentChannelText(order.payment_channel) : t('orderDetail.none') }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.trackingNo')">{{ order.tracking_no || t('orderDetail.none') }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.logisticsCompany')">{{ order.logistics_company || t('orderDetail.none') }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.quantity')">{{ order.num }} {{ t('orderDetail.quantityUnit') }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.amount')"><span class="amount">¥{{ totalAmount(order) }}</span></el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.recipient')">{{ order.address_name || "-" }} {{ order.address_phone || "" }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.address')">{{ order.address || "-" }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.refundStatus')">{{ refundStatusText(order.refund_status || 0) }}</el-descriptions-item>
      <el-descriptions-item :label="t('orderDetail.refundReason')">{{ order.refund_reason || t("orderDetail.none") }}</el-descriptions-item>
    </el-descriptions>

    <el-card class="section-card" shadow="never">
      <template #header>{{ t("orderDetail.operationLogTitle") }}</template>
      <el-timeline v-if="operationLogItems.length">
        <el-timeline-item v-for="item in operationLogItems" :key="item.id" :timestamp="item.time" type="primary">
          <div class="timeline-item">
            <b>{{ item.title }}</b>
            <span v-if="item.meta" class="muted">{{ item.meta }}</span>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else :description="t('orderDetail.noOperationLog')" />
    </el-card>

    <el-card class="section-card" shadow="never">
      <template #header>{{ t("orderDetail.logisticsTitle") }}</template>
      <el-timeline v-if="logisticsTimeline.length">
        <el-timeline-item v-for="item in logisticsTimeline" :key="item.key" :timestamp="item.time" :type="item.done ? 'primary' : 'info'" :hollow="!item.done">
          <div class="timeline-item">
            <b>{{ item.title }}</b>
            <span v-if="item.meta" class="muted">{{ item.meta }}</span>
          </div>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else :description="t('orderDetail.noLogistics')" />
    </el-card>

    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span>{{ t("orderDetail.afterSaleTitle") }}</span>
          <el-button v-if="canRequestAfterSale" size="small" type="primary" @click="openAfterSaleDialog">
            {{ t("orderDetail.requestAfterSale") }}
          </el-button>
        </div>
      </template>
      <el-empty v-if="!afterSales.length" :description="t('orderDetail.noAfterSale')" />
      <div v-else class="after-sale-list">
        <div v-for="item in afterSales" :key="item.id" class="after-sale-item">
          <div class="after-sale-head">
            <div>
              <b>{{ afterSaleTypeText(item.type) }}</b>
              <span class="muted"> · {{ afterSaleStatusText(item.status) }}</span>
            </div>
            <div class="after-sale-head-right">
              <el-tag size="small" :type="afterSaleStatusTagType(item.status)">{{ afterSaleStatusText(item.status) }}</el-tag>
              <span class="amount">{{ t("orderDetail.afterSaleRefundAmount") }}: ¥{{ Number(item.refund_amount || 0).toFixed(2) }}</span>
            </div>
          </div>
          <div class="after-sale-line">{{ item.reason }}</div>
          <div v-if="item.seller_reason" class="after-sale-line muted">{{ t("orderDetail.sellerReason") }}: {{ item.seller_reason }}</div>
          <div v-if="item.platform_note" class="after-sale-line muted">{{ t("orderDetail.platformNote") }}: {{ item.platform_note }}</div>
        </div>
      </div>
    </el-card>

    <div class="order-actions">
      <el-button v-if="order.type === 1" type="primary" @click="startUnpaidOrderPayment(router, order)">{{ t("orderDetail.payUnpaid") }}</el-button>
      <el-button v-if="order.type === 3" type="primary" @click="handleReceive">{{ t("orderDetail.confirmReceive") }}</el-button>
      <el-button v-if="canRequestAfterSale" type="warning" @click="openAfterSaleDialog">{{ t("orderDetail.requestAfterSale") }}</el-button>
      <el-button v-if="order.type === 4" @click="openReviewDialog">{{ t("orderDetail.writeReview") }}</el-button>
      <el-button @click="$router.back()">{{ t("orderDetail.back") }}</el-button>
    </div>
  </el-card>

  <div v-else-if="loading" class="loading-wrap">
    <el-icon class="is-loading" size="40"><Loading /></el-icon>
  </div>

  <el-empty v-else :description="t('orderDetail.orderNotFound')">
    <el-button @click="$router.push('/user/orders')">{{ t("orderDetail.backToOrders") }}</el-button>
  </el-empty>

  <el-dialog v-model="reviewDialogVisible" :title="t('orderDetail.reviewTitle')" width="520px">
    <el-form label-width="70px">
      <el-form-item :label="t('orderDetail.rating')"><el-rate v-model="reviewForm.rating" /></el-form-item>
      <el-form-item :label="t('orderDetail.review')">
        <el-input v-model="reviewForm.content" type="textarea" :rows="4" maxlength="500" show-word-limit :placeholder="t('orderDetail.reviewPlaceholder')" />
      </el-form-item>
      <el-form-item :label="t('orderDetail.image')">
        <el-upload :show-file-list="false" accept="image/*" :before-upload="handleReviewImageUpload">
          <el-button :loading="uploadingImage">{{ t("orderDetail.uploadImage") }}</el-button>
        </el-upload>
        <div v-if="reviewForm.images.length" class="review-upload-list">
          <div v-for="img in reviewForm.images" :key="img" class="review-upload-item">
            <el-image :src="img" fit="cover" />
            <el-button link type="danger" @click="removeReviewImage(img)">{{ t("orderDetail.delete") }}</el-button>
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="reviewDialogVisible = false">{{ t("common.cancel") }}</el-button>
      <el-button type="primary" :loading="submittingReview" @click="submitReview">{{ t("orderDetail.submitReview") }}</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="afterSaleDialogVisible" :title="t('orderDetail.afterSaleRequestTitle')" width="520px">
    <el-form label-width="90px">
      <el-form-item :label="t('orderDetail.afterSaleType')">
        <el-segmented v-model="afterSaleForm.type" :options="afterSaleTypeOptions" />
      </el-form-item>
      <el-form-item :label="t('orderDetail.afterSaleReason')">
        <el-input v-model="afterSaleForm.reason" type="textarea" :rows="4" maxlength="255" show-word-limit :placeholder="t('orderDetail.afterSaleReasonPlaceholder')" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="afterSaleDialogVisible = false">{{ t("common.cancel") }}</el-button>
      <el-button type="primary" :loading="submittingAfterSale" @click="submitAfterSale">{{ t("orderDetail.submitAfterSale") }}</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { Loading } from "@element-plus/icons-vue";
import { getAfterSaleList, getOrderDetail, getOrderLogs, receiveOrder, requestAfterSale } from "@/api/order";
import { createReview, uploadReviewImage } from "@/api/review";
import { afterSaleStatusTagType, afterSaleStatusText, afterSaleTypeText, orderActionText, orderStatusText, orderStatusTagType, operatorTypeText, paymentChannelText, refundStatusText } from "@/utils/status-labels";
import { startUnpaidOrderPayment } from "@/utils/order-payment";
import { t } from "@/locales";

const route = useRoute();
const router = useRouter();
const order = ref<any>(null);
const loading = ref(false);
const operationLogs = ref<any[]>([]);
const afterSales = ref<any[]>([]);
const reviewDialogVisible = ref(false);
const afterSaleDialogVisible = ref(false);
const submittingReview = ref(false);
const submittingAfterSale = ref(false);
const uploadingImage = ref(false);
const reviewForm = ref({ rating: 5, content: "", images: [] as string[] });
const afterSaleForm = ref({ type: "refund_only", reason: "" });

const statusText = orderStatusText;
const statusTagType = orderStatusTagType;

const afterSaleTypeOptions = computed(() => [
  { label: t("orderDetail.afterSaleRefundOnly"), value: "refund_only" },
  { label: t("orderDetail.afterSaleReturnRefund"), value: "return_refund" },
]);

const activeAfterSale = computed(
  () => afterSales.value.find((item) => !["refunded", "closed"].includes(item.status)) ?? null,
);
const latestAfterSale = computed(
  () =>
    [...afterSales.value].sort((a, b) => {
      const aTime = Number(a.updated_at ?? a.created_at ?? 0);
      const bTime = Number(b.updated_at ?? b.created_at ?? 0);
      return bTime - aTime || Number(b.id ?? 0) - Number(a.id ?? 0);
    })[0] ?? null,
);

const canRequestAfterSale = computed(
  () => [2, 3, 4].includes(Number(order.value?.type ?? 0)) && !activeAfterSale.value,
);

const operationLogItems = computed(() =>
  operationLogs.value.map((item) => ({
    id: item.id,
    title: orderActionText(item.action),
    meta: `${operatorTypeText(item.operator_type)}${item.remark ? ` · ${item.remark}` : ""}`,
    time: formatTime(item.created_at),
  })),
);

const logisticsTimeline = computed(() => {
  const latestLog = (action: string) =>
    [...operationLogs.value].reverse().find((item) => item.action === action) ?? null;
  const createLog = latestLog("create");
  const payLog = latestLog("pay");
  const shipLog = latestLog("ship");
  const receiveLog = latestLog("receive");
  const cancelLog = latestLog("cancel");
  const formatLogMeta = (log: any, fallback?: string) => {
    const parts = [log ? operatorTypeText(log.operator_type) : "", log?.remark ?? fallback ?? ""].filter(Boolean);
    return parts.join(" · ");
  };
  const list = [
    {
      key: "created",
      title: t("orderDetail.logisticsCreated"),
      meta: formatLogMeta(createLog),
      time: formatTime(createLog?.created_at ?? order.value?.created_at),
      done: true,
    },
    {
      key: "paid",
      title: t("orderDetail.logisticsPaid"),
      meta: [
        payLog ? operatorTypeText(payLog.operator_type) : "",
        paymentChannelText(payLog?.remark || order.value?.payment_channel || ""),
      ]
        .filter(Boolean)
        .join(" · "),
      time: formatTime(payLog?.created_at ?? order.value?.paid_at),
      done: Boolean(payLog || order.value?.paid_at),
    },
    {
      key: "shipped",
      title: t("orderDetail.logisticsShipped"),
      meta: [order.value?.logistics_company || "", order.value?.tracking_no || "", shipLog?.remark || ""].filter(Boolean).join(" · "),
      time: formatTime(shipLog?.created_at ?? order.value?.shipped_at),
      done: Boolean(shipLog || order.value?.shipped_at),
    },
    {
      key: "received",
      title: t("orderDetail.logisticsReceived"),
      meta: formatLogMeta(receiveLog),
      time: formatTime(receiveLog?.created_at ?? order.value?.received_at),
      done: Boolean(receiveLog || order.value?.received_at),
    },
    {
      key: "canceled",
      title: t("orderDetail.logisticsCanceled"),
      meta: formatLogMeta(cancelLog),
      time: formatTime(cancelLog?.created_at ?? order.value?.canceled_at),
      done: Boolean(cancelLog || order.value?.canceled_at),
    },
  ];
  if (latestAfterSale.value) {
    list.push({
      key: "after-sale",
      title: t("orderDetail.afterSaleTitle"),
      meta: [
        t("orderDetail.afterSaleStatus"),
        afterSaleStatusText(latestAfterSale.value.status),
        `${t("orderDetail.afterSaleRefundAmount")}: ¥${Number(latestAfterSale.value.refund_amount || 0).toFixed(2)}`,
        latestAfterSale.value.seller_reason || "",
        latestAfterSale.value.platform_note || "",
      ].filter(Boolean).join(" · "),
      time: formatTime(
        latestAfterSale.value.refunded_at ??
          latestAfterSale.value.closed_at ??
          latestAfterSale.value.updated_at ??
          latestAfterSale.value.created_at,
      ),
      done: true,
    });
  }
  return list.filter((item) => item.done || item.key === "created" || item.key === "after-sale");
});

function formatTime(value: any) {
  if (!value) return "";
  const date = typeof value === "number" ? new Date(value * 1000) : new Date(value);
  return Number.isNaN(date.getTime()) ? "" : date.toLocaleString();
}

function totalAmount(row: any) {
  return (Number(row.money || 0) * Number(row.num || 0)).toFixed(2);
}

async function loadOrder() {
  loading.value = true;
  try {
    const orderId = Number(route.params.id);
    const [detailRes, logsRes, afterSaleRes]: any[] = await Promise.all([
      getOrderDetail({ order_id: orderId }),
      getOrderLogs({ order_id: orderId }),
      getAfterSaleList({ order_id: orderId, page_num: 1, page_size: 20 }),
    ]);
    order.value = detailRes.data;
    operationLogs.value = logsRes.data?.item ?? logsRes.data ?? [];
    afterSales.value = afterSaleRes.data?.item ?? afterSaleRes.data ?? [];
  } catch {
    order.value = null;
    operationLogs.value = [];
    afterSales.value = [];
  } finally {
    loading.value = false;
  }
}

async function handleReceive() {
  try {
    await ElMessageBox.confirm(t("orderDetail.receiveConfirm"), t("dialog.warningTitle"), { type: "warning" });
    await receiveOrder({ order_id: order.value.id });
    ElMessage.success(t("orderDetail.receivedSuccess"));
    await loadOrder();
  } catch {}
}

function openReviewDialog() {
  reviewForm.value = { rating: 5, content: "", images: [] };
  reviewDialogVisible.value = true;
}

function openAfterSaleDialog() {
  afterSaleForm.value = { type: "refund_only", reason: "" };
  afterSaleDialogVisible.value = true;
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
    if (res.data?.url) reviewForm.value.images.push(res.data.url);
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

async function submitAfterSale() {
  if (!afterSaleForm.value.reason.trim()) {
    return ElMessage.warning(t("orderDetail.afterSaleReasonRequired"));
  }
  submittingAfterSale.value = true;
  try {
    await requestAfterSale({
      order_id: order.value.id,
      type: afterSaleForm.value.type,
      reason: afterSaleForm.value.reason.trim(),
    });
    ElMessage.success(t("orderDetail.afterSaleSubmitted"));
    afterSaleDialogVisible.value = false;
    await loadOrder();
  } finally {
    submittingAfterSale.value = false;
  }
}

onMounted(loadOrder);
</script>

<style scoped>
.header-row, .section-header { display: flex; align-items: center; justify-content: space-between; }
.detail-grid, .section-card { margin-bottom: 16px; }
.amount { color: #f56c6c; font-weight: 600; }
.timeline-item { display: flex; gap: 8px; flex-wrap: wrap; }
.muted { color: #909399; font-size: 12px; }
.order-actions { display: flex; justify-content: flex-end; gap: 10px; }
.loading-wrap { text-align: center; padding: 60px; }
.review-upload-list { display: flex; gap: 10px; flex-wrap: wrap; margin-top: 10px; }
.review-upload-item { width: 84px; }
.review-upload-item .el-image { width: 84px; height: 84px; border-radius: 4px; display: block; }
.after-sale-list { display: flex; flex-direction: column; gap: 12px; }
.after-sale-item { border: 1px solid #ebeef5; border-radius: 6px; padding: 12px; }
.after-sale-head { display: flex; justify-content: space-between; gap: 12px; }
.after-sale-head-right { display: flex; align-items: center; gap: 10px; }
.after-sale-line { margin-top: 6px; }
</style>
