<template>
  <div class="checkout-wrap">
    <h2 class="page-title">确认订单</h2>

    <!-- 收货地址 -->
    <el-card class="section-card">
      <template #header>
        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
          "
        >
          <span>收货地址</span>
          <el-button link type="primary" @click="showAddressDialog = true">
            {{ addresses.length ? "更换地址" : "新增地址" }}
          </el-button>
        </div>
      </template>
      <div v-if="selectedAddress" class="address-item selected">
        <el-icon><Location /></el-icon>
        <span
          ><b>{{ selectedAddress.name }}</b> {{ selectedAddress.phone }}</span
        >
        <span style="margin-left: 12px; color: #666">{{
          selectedAddress.address
        }}</span>
      </div>
      <el-empty v-else description="请先添加收货地址" :image-size="60">
        <el-button type="primary" size="small" @click="showAddressDialog = true"
          >新增地址</el-button
        >
      </el-empty>
    </el-card>

    <!-- 商品清单 -->
    <el-card class="section-card">
      <template #header>商品清单</template>
      <el-empty v-if="!checkoutItems.length" description="暂无待结算商品" />
      <div v-for="item in checkoutItems" :key="item.id" class="order-item">
        <img :src="item.img_path" class="order-img" />
        <div class="order-info">
          <div>{{ item.name }}</div>
          <div class="order-price">¥{{ unitPrice(item) }} × {{ item.num }}</div>
        </div>
        <div class="order-subtotal">
          ¥{{ (unitPriceValue(item) * item.num).toFixed(2) }}
        </div>
      </div>
    </el-card>

    <!-- 价格汇总 -->
    <el-card class="section-card">
      <div class="summary-row">
        <span>商品总价</span>
        <span>¥{{ totalPrice }}</span>
      </div>
      <div class="summary-row">
        <span>优惠券</span>
        <div class="coupon-row">
          <span :class="{ discount: couponDiscountValue > 0 }">
            {{ selectedCoupon ? `-${couponDiscount}` : "未选择" }}
          </span>
          <el-button link type="primary" @click="showCouponDialog = true">
            {{ selectedCoupon ? "更换" : "选择" }}
          </el-button>
        </div>
      </div>
      <el-divider style="margin: 10px 0" />
      <div class="summary-row total">
        <span>应付金额</span>
        <span class="total-price">¥{{ payablePrice }}</span>
      </div>
    </el-card>

    <div class="checkout-footer">
      <el-button size="large" @click="$router.push('/cart')"
        >返回购物车</el-button
      >
      <el-button
        type="primary"
        size="large"
        :loading="submitting"
        :disabled="!selectedAddress || !checkoutItems.length"
        @click="handleSubmit"
      >
        提交订单
      </el-button>
    </div>

    <!-- 地址弹窗 -->
    <el-dialog v-model="showAddressDialog" title="选择收货地址" width="500px">
      <div
        v-for="addr in addresses"
        :key="addr.id"
        :class="['address-option', { active: selectedAddress?.id === addr.id }]"
        @click="
          selectedAddress = addr;
          showAddressDialog = false;
        "
      >
        <b>{{ addr.name }}</b> {{ addr.phone }}<br />
        <span style="color: #666; font-size: 13px">{{ addr.address }}</span>
      </div>
      <el-divider />
      <el-form :model="newAddr" label-width="70px" size="small">
        <el-form-item label="姓名"
          ><el-input v-model="newAddr.name"
        /></el-form-item>
        <el-form-item label="手机号"
          ><el-input v-model="newAddr.phone"
        /></el-form-item>
        <el-form-item label="地址"
          ><el-input v-model="newAddr.address" type="textarea"
        /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddressDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAddAddress"
          >保存新地址</el-button
        >
      </template>
    </el-dialog>

    <el-dialog v-model="showCouponDialog" title="选择优惠券" width="560px">
      <el-skeleton v-if="couponLoading" :rows="3" animated />
      <el-empty v-else-if="!usableCoupons.length" description="暂无可用优惠券" />
      <div v-else class="coupon-list">
        <div
          v-for="coupon in usableCoupons"
          :key="coupon.id"
          :class="[
            'coupon-option',
            {
              active: selectedCoupon?.id === coupon.id,
              disabled: !eligibleItemForCoupon(coupon),
            },
          ]"
          @click="selectCoupon(coupon)"
        >
          <div>
            <b>{{ coupon.name }}</b>
            <div class="coupon-desc">
              {{ couponText(coupon) }} · {{ couponScopeText(coupon) }}
            </div>
          </div>
          <el-tag v-if="!eligibleItemForCoupon(coupon)" type="info"
            >未达门槛</el-tag
          >
          <el-tag v-else type="success">可用</el-tag>
        </div>
      </div>
      <template #footer>
        <el-button @click="clearCoupon">不使用优惠券</el-button>
        <el-button type="primary" @click="showCouponDialog = false"
          >确定</el-button
        >
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Location } from "@element-plus/icons-vue";
import { getAddressList, createAddress } from "@/api/address";
import { createOrder } from "@/api/order";
import { deleteCart } from "@/api/cart";
import { getUserCouponList } from "@/api/coupon";

const router = useRouter();

const checkoutItems = ref<any[]>(
  JSON.parse(sessionStorage.getItem("checkout_items") ?? "[]"),
);
const addresses = ref<any[]>([]);
const selectedAddress = ref<any>(null);
const showAddressDialog = ref(false);
const showCouponDialog = ref(false);
const submitting = ref(false);
const couponLoading = ref(false);
const coupons = ref<any[]>([]);
const selectedCoupon = ref<any>(null);
const newAddr = reactive({ name: "", phone: "", address: "" });

const totalPrice = computed(() =>
  checkoutItems.value
    .reduce((s, i) => s + unitPriceValue(i) * i.num, 0)
    .toFixed(2),
);

const usableCoupons = computed(() =>
  coupons.value.filter(
    (coupon) =>
      coupon.status === 0 && new Date(coupon.expire_at).getTime() > Date.now(),
  ),
);

const couponEligibleItem = computed(() =>
  selectedCoupon.value ? eligibleItemForCoupon(selectedCoupon.value) : null,
);

const couponDiscountValue = computed(() => {
  if (!selectedCoupon.value || !couponEligibleItem.value) return 0;
  const item = couponEligibleItem.value;
  const originUnit = unitPriceValue(item);
  const finalUnit = discountedUnitPrice(originUnit, selectedCoupon.value);
  return Math.max(0, originUnit - finalUnit) * Number(item.num || 0);
});

const couponDiscount = computed(() => `¥${couponDiscountValue.value.toFixed(2)}`);

const payablePrice = computed(() =>
  Math.max(0, Number(totalPrice.value) - couponDiscountValue.value).toFixed(2),
);

function unitPriceValue(item: any) {
  return Number(item.discount_price || item.price || 0);
}

function unitPrice(item: any) {
  return unitPriceValue(item).toFixed(2);
}

function discountedUnitPrice(unitPrice: number, coupon: any) {
  if (!coupon || unitPrice < Number(coupon.min_amount || 0)) return unitPrice;
  if (coupon.coupon_type === 2) {
    return unitPrice * Number(coupon.discount || 1);
  }
  return Math.max(0, unitPrice - Number(coupon.discount || 0));
}

function eligibleItemForCoupon(coupon: any) {
  return checkoutItems.value.find(
    (item) => unitPriceValue(item) >= Number(coupon.min_amount || 0),
  );
}

function couponText(coupon: any) {
  if (coupon.coupon_type === 2) {
    return `${Number(coupon.discount * 10).toFixed(1)}折`;
  }
  return `减 ¥${Number(coupon.discount || 0).toFixed(0)}`;
}

function couponScopeText(coupon: any) {
  const item = eligibleItemForCoupon(coupon);
  if (!item) return `满 ${Number(coupon.min_amount || 0).toFixed(0)} 可用`;
  return `应用到 ${item.name || "首个可用商品"}`;
}

function selectCoupon(coupon: any) {
  if (!eligibleItemForCoupon(coupon)) return;
  selectedCoupon.value = coupon;
}

function clearCoupon() {
  selectedCoupon.value = null;
  showCouponDialog.value = false;
}

async function loadAddresses() {
  try {
    const res: any = await getAddressList();
    addresses.value = res.data?.item ?? [];
    if (addresses.value.length) selectedAddress.value = addresses.value[0];
  } catch {}
}

async function loadCoupons() {
  couponLoading.value = true;
  try {
    const res: any = await getUserCouponList();
    coupons.value = res.data?.item ?? [];
  } catch {
    coupons.value = [];
  } finally {
    couponLoading.value = false;
  }
}

async function handleAddAddress() {
  if (!newAddr.name || !newAddr.phone || !newAddr.address) {
    return ElMessage.warning("请填写完整地址信息");
  }
  await createAddress({ ...newAddr });
  ElMessage.success("地址已保存");
  newAddr.name = "";
  newAddr.phone = "";
  newAddr.address = "";
  await loadAddresses();
  showAddressDialog.value = false;
}

async function handleSubmit() {
  if (!selectedAddress.value) return ElMessage.warning("请选择收货地址");
  if (!checkoutItems.value.length) return ElMessage.warning("订单为空");

  // 每件购物车商品单独创建一个订单（当前后端一个订单对应一个商品）
  submitting.value = true;
  try {
    const pendingOrders: any[] = [];
    for (const item of checkoutItems.value) {
      const useCoupon =
        selectedCoupon.value && couponEligibleItem.value?.id === item.id;
      const res: any = await createOrder({
        product_id: item.product_id,
        num: item.num,
        address_id: selectedAddress.value.id,
        boss_id: item.boss_id,
        money: Math.round(unitPriceValue(item)),
        coupon_id: useCoupon ? selectedCoupon.value.coupon_id : undefined,
      } as any);
      pendingOrders.push({
        ...item,
        money: res.data?.money,
        order_id: res.data?.id,
        order_num: res.data?.order_num,
        address_id: selectedAddress.value.id,
      });
    }
    await Promise.all(checkoutItems.value.map((item) => deleteCart({ id: item.id })));
    // 将订单信息传给支付页
    sessionStorage.setItem("pending_orders", JSON.stringify(pendingOrders));
    sessionStorage.removeItem("checkout_items");
    router.push("/payment");
  } catch (err) {
    ElMessage.error("下单失败，请重试");
  } finally {
    submitting.value = false;
  }
}

onMounted(() => {
  loadAddresses();
  loadCoupons();
});
</script>

<style scoped>
.checkout-wrap {
  max-width: 800px;
  margin: 0 auto;
}
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 16px;
}
.section-card {
  margin-bottom: 16px;
}
.address-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px;
  border-radius: 6px;
  background: #f0f7ff;
}
.address-option {
  padding: 10px 12px;
  border: 1px solid #eee;
  border-radius: 6px;
  margin-bottom: 8px;
  cursor: pointer;
}
.address-option.active {
  border-color: #409eff;
  background: #ecf5ff;
}
.order-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 0;
  border-bottom: 1px solid #f5f5f5;
}
.order-img {
  width: 50px;
  height: 50px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}
.order-info {
  flex: 1;
  font-size: 14px;
}
.order-price {
  color: #999;
  font-size: 13px;
  margin-top: 4px;
}
.order-subtotal {
  font-weight: 500;
  color: #333;
}
.summary-row {
  display: flex;
  justify-content: space-between;
  padding: 6px 0;
  font-size: 14px;
}
.summary-row.total {
  font-size: 16px;
  font-weight: 600;
}
.total-price {
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
}
.coupon-row {
  display: flex;
  align-items: center;
  gap: 8px;
}
.discount {
  color: #67c23a;
}
.coupon-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.coupon-option {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  border: 1px solid #eee;
  border-radius: 6px;
  cursor: pointer;
}
.coupon-option.active {
  border-color: #409eff;
  background: #ecf5ff;
}
.coupon-option.disabled {
  cursor: not-allowed;
  opacity: 0.55;
}
.coupon-desc {
  margin-top: 4px;
  color: #666;
  font-size: 12px;
}
.checkout-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}
</style>
