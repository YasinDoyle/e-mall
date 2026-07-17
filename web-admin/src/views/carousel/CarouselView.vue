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
              >商品ID: {{ item.product_id }}</span
            >
            <el-button size="small" type="danger" @click="handleDelete(item.id)"
              >删除</el-button
            >
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog
      v-model="dialogVisible"
      title="新增轮播图"
      width="520px"
      @closed="resetForm"
    >
      <el-form :model="form" label-width="90px">
        <el-form-item label="上传图片">
          <el-upload
            accept="image/*"
            :show-file-list="false"
            :before-upload="handleImageUpload"
          >
            <el-button>选择图片</el-button>
          </el-upload>
          <div v-if="imagePreview" class="image-preview">
            <img :src="imagePreview" />
          </div>
        </el-form-item>
        <el-form-item label="图片地址">
          <el-input
            v-model="form.img_path"
            placeholder="输入图片 URL"
            @input="handleManualUrlInput"
          />
        </el-form-item>
        <el-form-item label="关联商品ID">
          <el-input-number
            v-model="form.product_id"
            :min="1"
            placeholder="请输入商品ID"
          />
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
import { computed, ref, onMounted, reactive } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import {
  getCarouselList,
  createCarousel,
  deleteCarousel,
  uploadCarouselImage,
} from "@/api";

const list = ref<any[]>([]);
const dialogVisible = ref(false);
const saving = ref(false);
const selectedFile = ref<File | null>(null);
const previewUrl = ref("");
const form = reactive<{ img_path: string; product_id?: number }>({
  img_path: "",
  product_id: undefined,
});
const imagePreview = computed(() => previewUrl.value || form.img_path.trim());

async function loadList() {
  const res: any = await getCarouselList();
  list.value = res.data?.item ?? [];
}

async function handleSave() {
  if (!selectedFile.value && !form.img_path.trim()) {
    return ElMessage.warning("请上传图片或输入图片地址");
  }
  if (!form.product_id || form.product_id < 1) {
    return ElMessage.warning("请输入有效的关联商品ID");
  }
  saving.value = true;
  try {
    let imgPath = form.img_path.trim();
    if (selectedFile.value) {
      const formData = new FormData();
      formData.append("file", selectedFile.value);
      const uploadRes: any = await uploadCarouselImage(formData);
      imgPath = uploadRes.data?.url ?? "";
    }
    if (!imgPath) return ElMessage.warning("图片上传失败");

    await createCarousel({
      img_path: imgPath,
      product_id: form.product_id,
    });
    ElMessage.success("创建成功");
    dialogVisible.value = false;
    loadList();
  } finally {
    saving.value = false;
  }
}

function handleImageUpload(file: File) {
  if (!file.type.startsWith("image/")) {
    ElMessage.warning("请选择图片文件");
    return false;
  }
  selectedFile.value = file;
  form.img_path = "";
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
  }
  previewUrl.value = URL.createObjectURL(file);
  return false;
}

function handleManualUrlInput() {
  if (selectedFile.value) {
    selectedFile.value = null;
  }
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
    previewUrl.value = "";
  }
}

async function handleDelete(id: number) {
  await ElMessageBox.confirm("确认删除该轮播图？", "提示", { type: "warning" });
  await deleteCarousel({ id });
  ElMessage.success("删除成功");
  loadList();
}

function resetForm() {
  if (previewUrl.value) {
    URL.revokeObjectURL(previewUrl.value);
  }
  form.img_path = "";
  form.product_id = undefined;
  selectedFile.value = null;
  previewUrl.value = "";
}

onMounted(loadList);
</script>

<style scoped>
.image-preview {
  width: 100%;
  margin-top: 10px;
}
.image-preview img {
  width: 220px;
  height: 100px;
  object-fit: cover;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
</style>
