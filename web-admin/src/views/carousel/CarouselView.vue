<template>
  <el-card>
    <template #header>
      <div
        style="
          display: flex;
          justify-content: space-between;
          align-items: center;
        "
      >
        <span>轮播图管理</span>
        <el-button type="primary" @click="dialogVisible = true"
          >新增轮播图</el-button
        >
      </div>
    </template>

    <el-row :gutter="16">
      <el-col
        :span="6"
        v-for="item in list"
        :key="item.id"
        style="margin-bottom: 16px"
      >
        <el-card shadow="hover">
          <img
            :src="item.img_path"
            style="
              width: 100%;
              height: 120px;
              object-fit: cover;
              border-radius: 4px;
            "
          />
          <div
            style="
              display: flex;
              justify-content: space-between;
              align-items: center;
              margin-top: 8px;
            "
          >
            <span style="font-size: 12px; color: #999"
              >商品ID: {{ item.product_id || "-" }}</span
            >
            <el-button size="small" type="danger" @click="handleDelete(item.id)"
              >删除</el-button
            >
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="dialogVisible" title="新增轮播图" width="440px">
      <el-form :model="form" label-width="90px">
        <el-form-item label="图片地址">
          <el-input v-model="form.img_path" placeholder="输入图片 URL" />
        </el-form-item>
        <el-form-item label="关联商品ID">
          <el-input-number v-model="form.product_id" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave"
          >保存</el-button
        >
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getCarouselList, createCarousel, deleteCarousel } from "@/api";

const list = ref<any[]>([]);
const dialogVisible = ref(false);
const saving = ref(false);
const form = reactive({ img_path: "", product_id: 0 });

async function loadList() {
  const res: any = await getCarouselList();
  list.value = res.data?.item ?? [];
}

async function handleSave() {
  if (!form.img_path.trim()) return ElMessage.warning("请输入图片地址");
  saving.value = true;
  try {
    await createCarousel({
      img_path: form.img_path,
      product_id: form.product_id || undefined,
    });
    ElMessage.success("创建成功");
    dialogVisible.value = false;
    form.img_path = "";
    form.product_id = 0;
    loadList();
  } finally {
    saving.value = false;
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该轮播图？", "提示", { type: "warning" });
  await deleteCarousel({ id });
  ElMessage.success("删除成功");
  loadList();
}

onMounted(loadList);
</script>
