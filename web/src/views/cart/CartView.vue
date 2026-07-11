<template>
  <div class="cart-wrap">
    <h2 class="page-title">购物车</h2>

    <el-empty v-if="!cartList.length" description="购物车是空的">
      <el-button type="primary" @click="$router.push('/products')"
        >去逛逛</el-button
      >
    </el-empty>

    <template v-else>
      <el-table
        ref="tableRef"
        :data="cartList"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />
        <el-table-column label="商品" min-width="260">
          <template #default="{ row }">
            <div
              class="cart-product"
              @click="$router.push(`/product/${row.product_id}`)"
            >
              <img :src="row.img_path" class="cart-img" />
              <span>{{ row.name }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="单价" width="110">
          <template #default="{ row }">
            <span class="price">¥{{ unitPrice(row) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="数量" width="160">
          <template #default="{ row }">
            <el-input-number
              v-model="row.num"
              :min="1"
              :max="row.max_num"
              size="small"
              @change="(val: number) => handleNumChange(row, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="小计" width="110">
          <template #default="{ row }">
            <span class="price"
              >¥{{ (unitPriceValue(row) * row.num).toFixed(2) }}</span
            >
          </template>
        </el-table-column>
        <el-table-column label="操作" width="80">
          <template #default="{ row }">
            <el-button link type="danger" @click="handleDelete(row)"
              >删除</el-button
            >
          </template>
        </el-table-column>
      </el-table>

      <div class="cart-footer">
        <div class="footer-left">
          <el-checkbox :model-value="allChecked" @change="toggleAll"
            >全选</el-checkbox
          >
          <el-button
            link
            type="danger"
            :disabled="!selected.length"
            @click="handleBatchDelete"
          >
            删除选中
          </el-button>
        </div>
        <div class="footer-right">
          <span
            >已选 <b>{{ selected.length }}</b> 件，合计：</span
          >
          <span class="total-price">¥{{ totalPrice }}</span>
          <el-button
            type="primary"
            size="large"
            :disabled="!selected.length"
            @click="goCheckout"
          >
            结算 ({{ selected.length }})
          </el-button>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import { getCartList, updateCart, deleteCart } from "@/api/cart";
import { useUserStore } from "@/stores/user";

const router = useRouter();
const userStore = useUserStore();

const cartList = ref<any[]>([]);
const selected = ref<any[]>([]);
const tableRef = ref<any>();

const allChecked = computed(
  () =>
    cartList.value.length > 0 &&
    selected.value.length === cartList.value.length,
);

const totalPrice = computed(() =>
  selected.value
    .reduce((sum, item) => sum + unitPriceValue(item) * item.num, 0)
    .toFixed(2),
);

function unitPriceValue(item: any) {
  return Number(item.discount_price || item.price || 0);
}

function unitPrice(item: any) {
  return unitPriceValue(item).toFixed(2);
}

async function loadCart() {
  try {
    const res: any = await getCartList();
    cartList.value = res.data?.item ?? [];
    selected.value = [];
    userStore.setCartCount(cartList.value.length);
  } catch {}
}

function handleSelectionChange(rows: any[]) {
  selected.value = rows;
}

function toggleAll(val: boolean) {
  if (val) {
    cartList.value.forEach((row) => tableRef.value?.toggleRowSelection(row, true));
  } else {
    tableRef.value?.clearSelection();
  }
}

async function handleNumChange(row: any, val: number) {
  try {
    await updateCart({
      id: row.id,
      num: val,
    });
  } catch {}
}

async function handleDelete(row: any) {
  await ElMessageBox.confirm("确认删除该商品？", "提示", { type: "warning" });
  await deleteCart({ id: row.id });
  ElMessage.success("已删除");
  loadCart();
}

async function handleBatchDelete() {
  await ElMessageBox.confirm(
    `确认删除选中的 ${selected.value.length} 件商品？`,
    "提示",
    { type: "warning" },
  );
  await Promise.all(selected.value.map((item) => deleteCart({ id: item.id })));
  ElMessage.success("已删除");
  loadCart();
}

function goCheckout() {
  if (!selected.value.length) return;
  // 将选中商品存入 sessionStorage，结算页读取
  sessionStorage.setItem("checkout_items", JSON.stringify(selected.value));
  router.push("/checkout");
}

onMounted(loadCart);
</script>

<style scoped>
.cart-wrap {
  max-width: 1000px;
  margin: 0 auto;
  background: #fff;
  padding: 24px;
  border-radius: 8px;
}
.page-title {
  font-size: 20px;
  font-weight: 600;
  margin-bottom: 20px;
}
.cart-product {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
}
.cart-img {
  width: 56px;
  height: 56px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}
.price {
  color: #f56c6c;
  font-weight: 500;
}
.cart-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 16px;
  padding: 16px 0;
  border-top: 1px solid #eee;
}
.footer-left {
  display: flex;
  align-items: center;
  gap: 16px;
}
.footer-right {
  display: flex;
  align-items: center;
  gap: 12px;
}
.total-price {
  font-size: 22px;
  font-weight: bold;
  color: #f56c6c;
}
</style>
