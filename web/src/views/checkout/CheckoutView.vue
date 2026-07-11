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
      <div v-for="item in checkoutItems" :key="item.id" class="order-item">
        <img :src="item.product_img" class="order-img" />
        <div class="order-info">
          <div>{{ item.product_name }}</div>
          <div class="order-price">¥{{ item.price }} × {{ item.num }}</div>
        </div>
        <div class="order-subtotal">
          ¥{{ (Number(item.price) * item.num).toFixed(2) }}
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
        <span style="color: #67c23a">暂无优惠</span>
      </div>
      <el-divider style="margin: 10px 0" />
      <div class="summary-row total">
        <span>应付金额</span>
        <span class="total-price">¥{{ totalPrice }}</span>
      </div>
    </el-card>

    <div class="checkout-footer">
      <el-button size="large" @click="$router.push('/cart')"
        >返回购物车</el-button
      >
      <el-button
        type="primary"
        size="large"
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { Location } from "@element-plus/icons-vue";
import { getAddressList, createAddress } from "@/api/address";
import { createOrder } from "@/api/order";

const router = useRouter();

const checkoutItems = ref<any[]>(
  JSON.parse(sessionStorage.getItem("checkout_items") ?? "[]"),
);
const addresses = ref<any[]>([]);
const selectedAddress = ref<any>(null);
const showAddressDialog = ref(false);
const newAddr = reactive({ name: "", phone: "", address: "" });

const totalPrice = computed(() =>
  checkoutItems.value
    .reduce((s, i) => s + Number(i.price) * i.num, 0)
    .toFixed(2),
);

async function loadAddresses() {
  try {
    const res: any = await getAddressList();
    addresses.value = res.data?.item ?? [];
    if (addresses.value.length) selectedAddress.value = addresses.value[0];
  } catch {}
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
  try {
    const orderNums: number[] = [];
    const orderIds: number[] = [];
    for (const item of checkoutItems.value) {
      const orderNum = Date.now() + Math.floor(Math.random() * 1000);
      const res: any = await createOrder({
        product_id: item.product_id,
        num: item.num,
        address_id: selectedAddress.value.id,
        boss_id: item.boss_id,
        money: Math.round(Number(item.price) * item.num),
        order_num: orderNum,
      } as any);
      orderNums.push(orderNum);
      orderIds.push(res.data?.id);
    }
    // 将订单信息传给支付页
    sessionStorage.setItem(
      "pending_orders",
      JSON.stringify(
        checkoutItems.value.map((item, i) => ({
          ...item,
          order_id: orderIds[i],
          order_num: orderNums[i],
          address_id: selectedAddress.value.id,
        })),
      ),
    );
    sessionStorage.removeItem("checkout_items");
    router.push("/payment");
  } catch (err) {
    ElMessage.error("下单失败，请重试");
  }
}

onMounted(loadAddresses);
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
.checkout-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 16px;
}
</style>
