<template>
  <div class="detail-wrap" v-if="item">
    <el-row :gutter="32">
      <el-col :span="12">
        <div class="flash-badge-large">秒杀专场</div>
        <el-card class="info-card">
          <h2 class="item-title">
            {{ item.title || `秒杀商品 #${item.product_id}` }}
          </h2>

          <div class="price-row">
            <span class="flash-price">¥{{ item.money }}</span>
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
            >
              <el-option
                v-for="addr in addresses"
                :key="addr.id"
                :value="addr.id"
                :label="`${addr.name} ${addr.phone} - ${addr.address}`"
              />
            </el-select>
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
            :disabled="grabbed || (remainStock !== null && remainStock <= 0)"
            @click="handleGrab"
          >
            <template v-if="grabbed">已抢购成功 🎉</template>
            <template v-else-if="remainStock !== null && remainStock <= 0"
              >已售罄</template
            >
            <template v-else>立即抢购</template>
          </el-button>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header>活动说明</template>
          <ul class="rules">
            <li>每人限购 1 件，先到先得</li>
            <li>秒杀商品不支持退换货</li>
            <li>抢购成功后需在 24 小时内完成支付，否则自动取消</li>
            <li>秒杀价格以实际支付为准</li>
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
import { ref, computed, onMounted } from "vue";
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

const stockPercent = computed(() => {
  const stock = remainStock.value ?? item.value?.num ?? 0;
  if (stock <= 0) return 100;
  if (stock <= 5) return 85;
  if (stock <= 20) return 55;
  return 25;
});

async function loadData() {
  try {
    const id = Number(route.params.id);
    const [itemRes, addrRes]: any[] = await Promise.all([
      getFlashSaleDetail({ id }),
      getAddressList(),
    ]);
    item.value = itemRes.data;
    addresses.value = addrRes.data?.item ?? [];
    if (addresses.value.length) selectedAddressId.value = addresses.value[0].id;
  } catch {
    item.value = null;
  } finally {
    loading.value = false;
  }
}

async function handleGrab() {
  if (!userStore.isLoggedIn) return router.push("/login");
  if (!selectedAddressId.value) return ElMessage.warning("请选择收货地址");
  if (!payKey.value) return ElMessage.warning("请输入支付密码");
  if (grabbing.value || grabbed.value) return;

  grabbing.value = true;
  try {
    await doFlashSale({
      product_id: item.value.product_id,
      boss_id: item.value.boss_id,
      address_id: selectedAddressId.value,
      key: payKey.value,
      num: 1,
      money: item.value.money,
    });
    grabbed.value = true;
    if (remainStock.value !== null) remainStock.value--;
    ElMessage.success("抢购成功！订单已生成，请前往我的订单查看");
  } catch {
    ElMessage.error("抢购失败，可能已售罄或您已购买过");
  } finally {
    grabbing.value = false;
  }
}

onMounted(loadData);
</script>

<style scoped>
.detail-wrap {
  max-width: 1000px;
  margin: 0 auto;
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
