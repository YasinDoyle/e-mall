<template>
  <el-card>
    <template #header>用户管理</template>
    <el-table :data="list" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="user_name" label="用户名" />
      <el-table-column prop="nick_name" label="昵称" />
      <el-table-column prop="email" label="邮箱" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
            {{ row.status === "active" ? "正常" : "封禁" }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="管理员" width="90">
        <template #default="{ row }">
          <el-tag v-if="row.is_admin" type="warning">是</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            size="small"
            :type="row.status === 'active' ? 'danger' : 'success'"
            @click="toggleBan(row)"
          >
            {{ row.status === "active" ? "封禁" : "解封" }}
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
  ElMessage.success(banned ? "封禁成功" : "已解封");
  loadList();
}

onMounted(loadList);
</script>
