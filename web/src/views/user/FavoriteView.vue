<template>
  <el-card>
    <template #header>我的收藏</template>
    <el-skeleton v-if="loading" :rows="4" animated />
    <el-empty v-else-if="!list.length" description="收藏夹是空的">
      <el-button type="primary" @click="$router.push('/products')"
        >去逛逛</el-button
      >
    </el-empty>
    <div v-else class="fav-grid">
      <el-card
        v-for="item in list"
        :key="item.id"
        shadow="hover"
        class="fav-card"
        @click="$router.push(`/product/${item.product_id}`)"
      >
        <img :src="item.img_path || item.product_img" class="fav-img" />
        <div class="fav-name">{{ item.name || item.product_name }}</div>
        <div
          style="
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-top: 6px;
          "
        >
          <span style="color: #f56c6c; font-weight: bold"
            >¥{{ item.discount_price || item.price }}</span
          >
          <el-button
            link
            type="danger"
            size="small"
            @click.stop="handleRemove(item.id)"
            >取消收藏</el-button
          >
        </div>
      </el-card>
    </div>
    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="prev, pager, next"
      style="margin-top: 16px; justify-content: center"
      @current-change="loadList"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getFavoriteList, deleteFavorite } from "@/api/favorite";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 12;
const total = ref(0);
const loading = ref(false);

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getFavoriteList({
      page_num: page.value,
      page_size: pageSize,
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } catch {
    list.value = [];
    total.value = 0;
  } finally {
    loading.value = false;
  }
}

async function handleRemove(id: number) {
  await deleteFavorite({ id });
  ElMessage.success("已取消收藏");
  loadList();
}

onMounted(loadList);
</script>

<style scoped>
.fav-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}
.fav-card {
  cursor: pointer;
}
.fav-img {
  width: 100%;
  height: 140px;
  object-fit: cover;
  border-radius: 4px;
}
.fav-name {
  font-size: 13px;
  margin-top: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
