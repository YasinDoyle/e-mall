<template>
  <div class="detail-wrap" v-if="item">
    <div class="detail-header">
      <el-button link @click="$router.push('/flash-sale')">返回秒杀列表</el-button>
      <span>距本场结束 {{ countdownText }}</span>
    </div>

    <el-row :gutter="32">
      <el-col :span="12">
        <div class="flash-badge-large">秒杀专场</div>
        <el-card class="info-card">
          <h2 class="item-title">{{ titleText }}</h2>

          <div class="price-row">
            <span class="flash-price">¥{{ moneyText }}</span>
            <el-tag type="danger" effect="dark" style="margin-left: 12px"
              >限时特惠</el-tag
            >
          </div>

          <div class="stock-section">
            <span class="stock-label">剩余库存</span>
            <el-progress
              :percentage="stockPercent"
              :stroke-width="12"
              status="exception"
              style="margin-top: 8px"
            />
            <span class="stock-num">{{ remainStock ?? item.num }} 件</span>
          </div>

          <el-divider />

          <!-- 选择收货地址 -->
          <el-form-item label="收货地址" style="margin-bottom: 12px">
            <el-select
              v-model="selectedAddressId"
              placeholder="请选择收货地址"
              style="width: 100%"
              :loading="addressLoading"
            >
              <el-option
                v-for="addr in addresses"
                :key="addr.id"
                :value="addr.id"
                :label="`${addr.name} ${addr.phone} - ${addr.address}`"
              />
            </el-select>
            <el-button
              v-if="userStore.isLoggedIn && !addresses.length && !addressLoading"
              link
              type="primary"
              @click="$router.push('/user/addresses')"
            >
              去新增地址
            </el-button>
          </el-form-item>

          <!-- 支付密码 -->
          <el-form-item label="支付密码">
            <el-input
              v-model="payKey"
              type="password"
              placeholder="请输入6位支付密码"
              maxlength="6"
              show-password
            />
          </el-form-item>

          <el-button
            type="danger"
            size="large"
            style="width: 100%"
            :loading="grabbing"
            :disabled="!canGrab"
            @click="handleGrab"
          >
            <template v-if="grabbed">已抢购成功</template>
            <template v-else-if="remainStock !== null && remainStock <= 0"
              >已售罄</template
            >
            <template v-else-if="!userStore.isLoggedIn">登录后抢购</template>
            <template v-else>立即抢购</template>
          </el-button>

          <el-button
            v-if="grabbed"
            size="large"
            style="width: 100%; margin-top: 10px"
            @click="$router.push('/user/orders')"
          >
            查看我的订单
          </el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>活动说明</template>
          <ul class="rules">
            <li>每人限购 1 件，先到先得</li>
            <li>抢购时会直接按秒杀价扣减余额</li>
            <li>请提前确认收货地址和支付密码</li>
            <li>重复提交会被系统拦截</li>
          </ul>
        </el-card>
      </el-col>
    </el-row>
  </div>

  <div v-else-if="loading" style="text-align: center; padding: 80px">
    <el-icon class="is-loading" size="48"><Loading /></el-icon>
  </div>

  <el-result v-else icon="warning" title="活动不存在或已结束">
    <template #extra>
      <el-button type="primary" @click="$router.push('/flash-sale')"
        >返回秒杀列表</el-button
      >
    </template>
  </el-result>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Loading } from "@element-plus/icons-vue";
import { getFlashSaleDetail, doFlashSale } from "@/api/flashSale";
import { getAddressList } from "@/api/address";
import { useUserStore } from "@/stores/user";

const route = useRoute();
const router = useRouter();
const userStore = useUserStore();

const item = ref<any>(null);
const loading = ref(true);
const grabbing = ref(false);
const grabbed = ref(false);
const remainStock = ref<number | null>(null);
const addresses = ref<any[]>([]);
const selectedAddressId = ref<number | undefined>(undefined);
const payKey = ref("");
const addressLoading = ref(false);
const now = ref(Date.now());
let timer: number | undefined;

const stockPercent = computed(() => {
  const stock = remainStock.value ?? stockNum.value;
  if (stock <= 0) return 100;
  if (stock <= 5) return 85;
  if (stock <= 20) return 55;
  return 25;
});

const countdownText = computed(() => {
  const end = new Date(now.value);
  end.setHours(23, 59, 59, 999);
  const diff = Math.max(0, end.getTime() - now.value);
  const hours = Math.floor(diff / 1000 / 60 / 60);
  const minutes = Math.floor((diff / 1000 / 60) % 60);
  const seconds = Math.floor((diff / 1000) % 60);
  return `${padTime(hours)}:${padTime(minutes)}:${padTime(seconds)}`;
});

const canGrab = computed(() => {
  const soldOut = remainStock.value !== null && remainStock.value <= 0;
  return (
    !grabbing.value &&
    !grabbed.value &&
    !soldOut &&
    (!userStore.isLoggedIn ||
      (!!selectedAddressId.value && payKey.value.length === 6))
  );
});

const flashSaleId = computed(() => item.value?.id ?? item.value?.Id);
const productId = computed(() => item.value?.product_id ?? item.value?.ProductId);
const bossId = computed(() => item.value?.boss_id ?? item.value?.BossId);
const stockNum = computed(() => Number(item.value?.num ?? item.value?.Num ?? 0));
const titleText = computed(
  () =>
    item.value?.title ??
    item.value?.Title ??
    `秒杀商品 #${productId.value || ""}`,
);
const moneyValue = computed(() => Number(item.value?.money ?? item.value?.Money ?? 0));
const moneyText = computed(() => moneyValue.value.toFixed(2));

function padTime(value: number) {
  return String(value).padStart(2, "0");
}

async function loadData() {
  try {
    const productId = Number(route.params.id);
    const itemRes: any = await getFlashSaleDetail({ product_id: productId });
    item.value = itemRes.data;
    remainStock.value = stockNum.value;
    if (userStore.isLoggedIn) {
      await loadAddresses();
    }
  } catch {
    item.value = null;
  } finally {
    loading.value = false;
  }
}

async function loadAddresses() {
  addressLoading.value = true;
  try {
    const addrRes: any = await getAddressList();
    addresses.value = addrRes.data?.item ?? [];
    if (addresses.value.length) selectedAddressId.value = addresses.value[0].id;
  } catch {
    addresses.value = [];
  } finally {
    addressLoading.value = false;
  }
}

async function handleGrab() {
  if (!userStore.isLoggedIn) {
    return router.push({ path: "/login", query: { redirect: route.fullPath } });
  }
  if (!selectedAddressId.value) return ElMessage.warning("请选择收货地址");
  if (payKey.value.length !== 6) return ElMessage.warning("请输入6位支付密码");
  if (grabbing.value || grabbed.value) return;

  grabbing.value = true;
  try {
    const res: any = await doFlashSale({
      flash_sale_id: flashSaleId.value,
      product_id: productId.value,
      boss_id: bossId.value,
      address_id: selectedAddressId.value,
      key: payKey.value,
      num: 1,
      money: moneyValue.value,
    });
    grabbed.value = true;
    remainStock.value = res.data?.remaining_stock ?? Math.max(0, (remainStock.value ?? 1) - 1);
    ElMessage.success("抢购成功，订单生成中");
  } catch {
    ElMessage.error("抢购失败，可能已售罄或您已购买过");
  } finally {
    grabbing.value = false;
  }
}

onMounted(() => {
  loadData();
  timer = window.setInterval(() => {
    now.value = Date.now();
  }, 1000);
});

onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});
</script>

<style scoped>
.detail-wrap {
  max-width: 1000px;
  margin: 0 auto;
}
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
  color: #666;
  font-size: 14px;
}
.flash-badge-large {
  display: inline-block;
  background: #f56c6c;
  color: #fff;
  padding: 4px 14px;
  border-radius: 20px;
  font-size: 13px;
  margin-bottom: 12px;
}
.info-card {
  padding: 8px;
}
.item-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0 0 12px;
}
.price-row {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}
.flash-price {
  font-size: 32px;
  font-weight: bold;
  color: #f56c6c;
}
.stock-section {
  margin-bottom: 16px;
}
.stock-label {
  font-size: 14px;
  color: #666;
}
.stock-num {
  font-size: 13px;
  color: #999;
  margin-top: 4px;
  display: block;
}
.rules {
  padding-left: 20px;
  color: #666;
  line-height: 2;
  font-size: 14px;
}
</style>
