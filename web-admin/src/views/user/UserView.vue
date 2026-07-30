<template>
  <el-card>
    <template #header>{{ t("page.user.title") }}</template>
    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_name" :label="t('page.user.username')" />
      <el-table-column prop="nick_name" :label="t('page.user.nickname')" />
      <el-table-column prop="email" :label="t('page.user.email')" />
      <el-table-column :label="t('common.status')" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === "active" ? t("page.user.active") : t("page.user.banned") }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('page.user.admin')" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.is_admin" type="warning">{{ t("page.user.yes") }}</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="120">
        <template #default="{ row }">
          <el-button
            size="small"
            :type="row.status === 'active' ? 'danger' : 'success'"
            @click="toggleBan(row)"
          >
            {{ row.status === "active" ? t("page.user.ban") : t("page.user.unban") }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      v-model:current-page="page"
      :page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      style="margin-top: 16px; justify-content: flex-end"
      @current-change="loadList"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { getUserList, banUser } from "@/api";
import { t } from "@/locales";

const list = ref<any[]>([]);
const page = ref(1);
const pageSize = 15;
const total = ref(0);

async function loadList() {
  const res: any = await getUserList({
    page_num: page.value,
    page_size: pageSize,
  });
  list.value = res.data?.item ?? [];
  total.value = res.data?.total ?? 0;
}

async function toggleBan(row: any) {
  const banned = row.status === "active";
  await banUser({ id: row.id, banned });
  ElMessage.success(banned ? t("page.user.banSuccess") : t("page.user.unbanSuccess"));
  loadList();
}

onMounted(loadList);
</script>
