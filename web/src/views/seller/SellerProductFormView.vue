<template>
  <el-card>
    <template #header>{{ isEdit ? "编辑商品" : "发布商品" }}</template>
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
      <el-divider content-position="left">审核资料</el-divider>
      <el-form-item label="品牌">
        <el-input
          v-model="form.brand"
          maxlength="100"
          placeholder="请输入品牌"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="产地">
        <el-input
          v-model="form.origin"
          maxlength="120"
          placeholder="请输入产地"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="规格">
        <el-input
          v-model="form.specification"
          maxlength="120"
          placeholder="如 500g/件、L 码、套装等"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="生产日期">
        <el-date-picker
          v-model="form.production_date"
          type="date"
          value-format="YYYY-MM-DD"
          placeholder="请选择生产日期"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="保质期">
        <el-input
          v-model="form.shelf_life"
          maxlength="80"
          placeholder="如 12个月、30天"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="服务保障">
        <el-input
          v-model="form.service_guarantees"
          maxlength="255"
          placeholder="如 七天无理由、正品保障"
          :disabled="!sellerStore.isApproved"
        />
      </el-form-item>
      <el-form-item label="证书说明">
        <el-input
          v-model="form.certificate_meta"
          type="textarea"
          :rows="3"
          maxlength="500"
          placeholder="填写资质证书、授权文件或质检报告摘要"
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
      <el-form-item label="商品图片" :required="!isEdit">
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
      <el-form-item v-if="isEdit && existingCertificates.length" label="已有证书">
        <div class="existing-certificates">
          <div
            v-for="certificate in existingCertificates"
            :key="certificate.id"
            class="existing-certificate"
          >
            <el-image
              :src="certificate.file_path"
              fit="cover"
              class="existing-certificate-img"
            />
            <span>{{ certificate.name }}</span>
          </div>
        </div>
      </el-form-item>
      <el-form-item label="资质证书">
        <div class="certificate-list">
          <div v-for="(item, index) in certificateItems" :key="item.id" class="certificate-item">
            <el-select
              v-model="item.certificate_type"
              placeholder="证书类型"
              style="width: 180px"
              :disabled="!sellerStore.isApproved"
            >
              <el-option label="资质证书" value="qualification" />
              <el-option label="质检报告" value="quality_inspection" />
              <el-option label="授权证书" value="authorization" />
              <el-option label="其他材料" value="other" />
            </el-select>
            <el-input
              v-model="item.name"
              maxlength="120"
              placeholder="证书名称"
              :disabled="!sellerStore.isApproved"
            />
            <el-upload
              v-model:file-list="item.files"
              :auto-upload="false"
              :limit="1"
              accept="image/*"
              :disabled="!sellerStore.isApproved"
            >
              <el-button :disabled="!sellerStore.isApproved">选择文件</el-button>
            </el-upload>
            <el-button :disabled="!sellerStore.isApproved" @click="removeCertificate(index)">
              删除
            </el-button>
          </div>
          <el-button
            type="primary"
            plain
            :disabled="!sellerStore.isApproved"
            @click="addCertificate"
          >
            {{ isEdit ? "添加替换证书" : "添加证书" }}
          </el-button>
          <span v-if="isEdit" class="muted">编辑时上传新证书会整体替换原资质材料</span>
        </div>
      </el-form-item>
      <el-form-item>
        <el-button
          type="primary"
          :loading="submitting"
          :disabled="!sellerStore.isApproved"
          @click="submit"
        >
          {{ isEdit ? "重新提交审核" : "提交审核" }}
        </el-button>
        <el-button @click="$router.push('/seller/products')">返回列表</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, type UploadUserFile } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import {
  createProduct,
  getCategoryList,
  getProductDetail,
  updateProduct,
} from "@/api/product";
import { useSellerStore } from "@/stores/seller";
import type { ProductCertificate } from "@/types";

const route = useRoute();
const router = useRouter();
const sellerStore = useSellerStore();
const submitting = ref(false);
const categories = ref<any[]>([]);
const fileList = ref<UploadUserFile[]>([]);
const existingCertificates = ref<ProductCertificate[]>([]);
const productID = computed(() => Number(route.params.id || 0));
const isEdit = computed(() => productID.value > 0);
const certificateItems = ref<
  {
    id: number;
    certificate_type: string;
    name: string;
    files: UploadUserFile[];
  }[]
>([]);
const form = ref({
  name: "",
  category_id: undefined as number | undefined,
  title: "",
  info: "",
  brand: "",
  origin: "",
  specification: "",
  production_date: "",
  shelf_life: "",
  service_guarantees: "",
  certificate_meta: "",
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
  if (!isEdit.value && !fileList.value.length) return "请上传商品图片";
  if (certificateItems.value.some((item) => item.files.length && !item.name.trim())) {
    return "请填写证书名称";
  }
  return "";
}

async function loadProductForEdit() {
  if (!isEdit.value) return;
  const res: any = await getProductDetail({ id: productID.value });
  const product = res.data;
  if (!product || product.boss_id !== sellerStore.profile?.user_id) {
    ElMessage.warning("只能编辑自己发布的商品");
    router.push("/seller/products");
    return;
  }
  form.value = {
    name: product.name ?? "",
    category_id: product.category_id,
    title: product.title ?? "",
    info: product.info ?? "",
    brand: product.brand ?? "",
    origin: product.origin ?? "",
    specification: product.specification ?? "",
    production_date: product.production_date ?? "",
    shelf_life: product.shelf_life ?? "",
    service_guarantees: product.service_guarantees ?? "",
    certificate_meta: product.certificate_meta ?? "",
    price: Number(product.price || 0),
    discount_price: Number(product.discount_price || 0),
    num: Number(product.num || 1),
  };
  existingCertificates.value = product.certificates ?? [];
}

function addCertificate() {
  certificateItems.value.push({
    id: Date.now() + certificateItems.value.length,
    certificate_type: "qualification",
    name: "",
    files: [],
  });
}

function removeCertificate(index: number) {
  certificateItems.value.splice(index, 1);
}

async function submit() {
  if (!sellerStore.isApproved) {
    return ElMessage.warning("商家入驻审核通过后才可以发布商品");
  }
  const message = validate();
  if (message) return ElMessage.warning(message);

  const data = new FormData();
  if (isEdit.value) {
    data.append("id", String(productID.value));
  }
  data.append("name", form.value.name.trim());
  data.append("category_id", String(form.value.category_id));
  data.append("title", form.value.title.trim());
  data.append("info", form.value.info.trim());
  data.append("brand", form.value.brand.trim());
  data.append("origin", form.value.origin.trim());
  data.append("specification", form.value.specification.trim());
  data.append("production_date", form.value.production_date || "");
  data.append("shelf_life", form.value.shelf_life.trim());
  data.append("service_guarantees", form.value.service_guarantees.trim());
  data.append("certificate_meta", form.value.certificate_meta.trim());
  data.append("price", form.value.price.toFixed(2));
  data.append("discount_price", form.value.discount_price.toFixed(2));
  data.append("num", String(form.value.num));
  fileList.value.forEach((file) => {
    if (file.raw) data.append("image", file.raw);
  });
  certificateItems.value.forEach((item) => {
    const file = item.files[0];
    if (!file?.raw) return;
    data.append("certificate", file.raw);
    data.append("certificate_type", item.certificate_type);
    data.append("certificate_name", item.name.trim() || file.name);
  });
  if (isEdit.value) {
    data.append(
      "replace_certificates",
      certificateItems.value.some((item) => !!item.files[0]?.raw) ? "true" : "false",
    );
  }

  submitting.value = true;
  try {
    if (isEdit.value) {
      await updateProduct(data);
      ElMessage.success("商品已重新提交审核");
    } else {
      await createProduct(data);
      ElMessage.success("商品已提交审核");
    }
    router.push("/seller/products");
  } finally {
    submitting.value = false;
  }
}

onMounted(async () => {
  await sellerStore.loadProfile({ silentError: true });
  if (sellerStore.isApproved) {
    await loadCategories();
    await loadProductForEdit();
  }
});
</script>

<style scoped>
.notice {
  margin-bottom: 18px;
}
.product-form {
  max-width: 880px;
}
.certificate-list {
  display: flex;
  width: 100%;
  flex-direction: column;
  gap: 12px;
}
.certificate-item {
  display: grid;
  grid-template-columns: 180px minmax(180px, 1fr) minmax(120px, auto) 64px;
  align-items: center;
  gap: 10px;
}
.existing-certificates {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
.existing-certificate {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #606266;
}
.existing-certificate-img {
  width: 48px;
  height: 48px;
  border: 1px solid #ebeef5;
  border-radius: 4px;
}
.muted {
  color: #909399;
  font-size: 12px;
}
</style>
