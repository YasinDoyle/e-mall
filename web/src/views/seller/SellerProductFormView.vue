<template>
  <el-card>
    <template #header>发布商品</template>
    <el-alert
      v-if="!sellerStore.isApproved"
      title="商家入驻审核通过后才可以发布商品"
      type="warning"
      :closable="false"
      show-icon
      class="notice"
    />
    <el-alert
      v-else
      title="商品发布后进入待审核状态，审核通过并上架后才会展示给买家"
      type="info"
      :closable="false"
      show-icon
      class="notice"
    />

    <el-form label-width="96px" :model="form" class="product-form">
      <el-form-item label="商品名称" required>
        <el-input
          v-model="form.name"
          maxlength="80"
          placeholder="请输入商品名称"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="分类" required>
        <el-select
          v-model="form.category_id"
          placeholder="请选择分类"
          style="width: 260px"
          :disabled="!sellerStore.isApproved"
        >
          <el-option
            v-for="item in categories"
            :key="item.id"
            :label="item.category_name"
            :value="item.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="标题" required>
        <el-input
          v-model="form.title"
          maxlength="120"
          placeholder="请输入商品标题"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="详情">
        <el-input
          v-model="form.info"
          type="textarea"
          :rows="5"
          maxlength="1000"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="价格" required>
        <el-input-number
          v-model="form.price"
          :min="0.01"
          :precision="2"
          :step="1"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="优惠价" required>
        <el-input-number
          v-model="form.discount_price"
          :min="0.01"
          :precision="2"
          :step="1"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="库存" required>
        <el-input-number
          v-model="form.num"
          :min="1"
          :step="1"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="商品图片" required>
        <el-upload
          v-model:file-list="fileList"
          list-type="picture-card"
          :auto-upload="false"
          :limit="6"
          accept="image/*"
          :disabled="!sellerStore.isApproved"
        >
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!sellerStore.isApproved"
          @click="submit"
        >
          提交审核
        </el-button>
        <el-button @click="$router.push('/seller/products')">返回列表</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, type UploadUserFile } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import { createProduct, getCategoryList } from "@/api/product";
import { useSellerStore } from "@/stores/seller";

const router = useRouter();
const sellerStore = useSellerStore();
const submitting = ref(false);
const categories = ref<any[]>([]);
const fileList = ref<UploadUserFile[]>([]);
const form = ref({
  name: "",
  category_id: undefined as number | undefined,
  title: "",
  info: "",
  price: 0,
  discount_price: 0,
  num: 1,
});

async function loadCategories() {
  const res: any = await getCategoryList();
  categories.value = res.data?.item ?? [];
}

function validate() {
  if (!form.value.name.trim()) return "请输入商品名称";
  if (!form.value.category_id) return "请选择分类";
  if (!form.value.title.trim()) return "请输入商品标题";
  if (form.value.price <= 0 || form.value.discount_price <= 0) {
    return "请输入有效价格";
  }
  if (!fileList.value.length) return "请上传商品图片";
  return "";
}

async function submit() {
  if (!sellerStore.isApproved) {
    return ElMessage.warning("商家入驻审核通过后才可以发布商品");
  }
  const message = validate();
  if (message) return ElMessage.warning(message);

  const data = new FormData();
  data.append("name", form.value.name.trim());
  data.append("category_id", String(form.value.category_id));
  data.append("title", form.value.title.trim());
  data.append("info", form.value.info.trim());
  data.append("price", form.value.price.toFixed(2));
  data.append("discount_price", form.value.discount_price.toFixed(2));
  data.append("num", String(form.value.num));
  fileList.value.forEach((file) => {
    if (file.raw) data.append("image", file.raw);
  });

  submitting.value = true;
  try {
    await createProduct(data);
    ElMessage.success("商品已提交审核");
    router.push("/seller/products");
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await sellerStore.loadProfile({ silentError: true });
  if (sellerStore.isApproved) {
    await loadCategories();
  }
});
</script>

<style scoped>
.notice {
  margin-bottom: 18px;
}
.product-form {
  max-width: 760px;
}
</style>
