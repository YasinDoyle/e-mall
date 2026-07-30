<template>
  <el-card>
    <template #header>{{ t("coupon.title") }}</template>

    <el-tabs v-model="activeTab" @tab-change="loadData">
      <el-tab-pane :label="t('coupon.claimable')" name="claimable" />
      <el-tab-pane :label="t('coupon.usable')" name="usable" />
      <el-tab-pane :label="t('coupon.used')" name="used" />
      <el-tab-pane :label="t('coupon.expired')" name="expired" />
    </el-tabs>

    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!displayCoupons.length" :description="t('coupon.empty')" />

    <div v-else class="coupon-grid">
      <div
        v-for="coupon in displayCoupons"
        :key="couponKey(coupon)"
        :class="['coupon-card', { disabled: isDisabled(coupon) }]"
      >
        <div class="coupon-main">
          <div class="coupon-discount">{{ discountText(coupon) }}</div>
          <div class="coupon-threshold">{{ thresholdText(coupon) }}</div>
        </div>
        <div class="coupon-info">
          <div class="coupon-name">{{ coupon.name }}</div>
          <div class="coupon-expire">
            {{ t("coupon.expireAt", { date: formatDate(coupon.expire_at) }) }}
          </div>
          <el-button
            v-if="activeTab === 'claimable'"
            size="small"
            type="primary"
            :loading="claimingId === coupon.id"
            @click="handleClaim(coupon.id)"
          >
            {{ t("coupon.claim") }}
          </el-button>
          <el-tag v-else :type="tagType(coupon)">
            {{ statusText(coupon) }}
          </el-tag>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { useI18n } from "vue-i18n";
import {
  claimCoupon,
  getCouponList,
  getUserCouponList,
} from "@/api/coupon";

const activeTab = ref("usable");
const { t } = useI18n();
const loading = ref(false);
const claimingId = ref<number | null>(null);
const claimableCoupons = ref<any[]>([]);
const userCoupons = ref<any[]>([]);

const displayCoupons = computed(() => {
  const now = Date.now();
  if (activeTab.value === "claimable") return claimableCoupons.value;
  return userCoupons.value.filter((coupon) => {
    const expired = new Date(coupon.expire_at).getTime() < now;
    if (activeTab.value === "usable") return coupon.status === 0 && !expired;
    if (activeTab.value === "used") return coupon.status === 1;
    return expired;
  });
});

function couponKey(coupon: any) {
  return `${activeTab.value}-${coupon.id}-${coupon.coupon_id ?? ""}`;
}

function discountText(coupon: any) {
  if (coupon.coupon_type === 2) {
    return t("coupon.discount", { value: Number(coupon.discount * 10).toFixed(1) });
  }
  return `¥${Number(coupon.discount || 0).toFixed(0)}`;
}

function thresholdText(coupon: any) {
  return Number(coupon.min_amount || 0) > 0
    ? t("coupon.threshold", { amount: Number(coupon.min_amount).toFixed(0) })
    : t("coupon.noThreshold");
}

function formatDate(value: string) {
  if (!value) return "-";
  return new Date(value).toLocaleDateString();
}

function isDisabled(coupon: any) {
  return activeTab.value !== "claimable" && statusText(coupon) !== t("coupon.usable");
}

function statusText(coupon: any) {
  if (coupon.status === 1) return t("coupon.used");
  if (new Date(coupon.expire_at).getTime() < Date.now()) return t("coupon.expired");
  return t("coupon.usable");
}

function tagType(coupon: any) {
  return statusText(coupon) === t("coupon.usable") ? "success" : "info";
}

async function loadData() {
  loading.value = true;
  try {
    if (activeTab.value === "claimable") {
      const res: any = await getCouponList();
      claimableCoupons.value = res.data?.item ?? [];
    } else {
      const res: any = await getUserCouponList();
      userCoupons.value = res.data?.item ?? [];
    }
  } finally {
    loading.value = false;
  }
}

async function handleClaim(couponId: number) {
  claimingId.value = couponId;
  try {
    await claimCoupon({ coupon_id: couponId });
    ElMessage.success(t("coupon.claimSuccess"));
    await loadData();
  } finally {
    claimingId.value = null;
  }
}

onMounted(loadData);
</script>

<style scoped>
.coupon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 14px;
}
.coupon-card {
  display: flex;
  min-height: 118px;
  border: 1px solid #f3d2d2;
  border-radius: 8px;
  overflow: hidden;
  background: #fff;
}
.coupon-card.disabled {
  filter: grayscale(0.7);
  opacity: 0.65;
}
.coupon-main {
  width: 104px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  background: #f56c6c;
  color: #fff;
}
.coupon-discount {
  font-size: 28px;
  font-weight: 700;
}
.coupon-threshold {
  font-size: 12px;
  margin-top: 6px;
}
.coupon-info {
  flex: 1;
  padding: 14px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  justify-content: center;
}
.coupon-name {
  font-weight: 600;
}
.coupon-expire {
  color: #999;
  font-size: 12px;
}
</style>
