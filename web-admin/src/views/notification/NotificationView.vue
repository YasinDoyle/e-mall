<template>
  <el-card>
    <template #header>
      <div class="page-head">
        <span>{{ t("page.notification.title") }}</span>
        <div class="actions">
          <el-checkbox v-model="unreadOnly">{{ t("page.notification.unreadOnly") }}</el-checkbox>
          <el-button size="small" :disabled="!list.length" @click="handleMarkAllRead">
            {{ t("page.notification.markAllRead") }}
          </el-button>
        </div>
      </div>
    </template>

    <el-table v-loading="loading" :data="list" style="width: 100%">
      <el-table-column :label="t('common.status')" width="90">
        <template #default="{ row }">
          <el-tag :type="row.read ? 'info' : 'danger'">
            {{ row.read ? t("status.notification.read") : t("status.notification.unread") }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="scene" :label="t('page.notification.scene')" width="140" />
      <el-table-column prop="title" :label="t('page.notification.titleColumn')" width="180" />
      <el-table-column prop="content" :label="t('page.notification.content')" show-overflow-tooltip />
      <el-table-column :label="t('common.time')" width="180">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.actions')" width="110">
        <template #default="{ row }">
          <el-button size="small" :disabled="row.read" @click="handleMarkRead(row)">
            {{ t("page.notification.markRead") }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pager">
      <el-pagination
        layout="prev, pager, next, total"
        :total="total"
        :page-size="pageSize"
        v-model:current-page="page"
        @current-change="loadList"
      />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import {
  getAdminNotificationList,
  markAdminNotificationRead,
  markAllAdminNotificationRead,
  type NotificationItem,
} from "@/api/notification";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";
import { t } from "@/locales";

const list = ref<NotificationItem[]>([]);
const total = ref(0);
const page = ref(1);
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
const pageSize = appConfig.config.default_page_size;
const loading = ref(false);
const unreadOnly = ref(false);

function formatTime(value: number) {
  if (!value) return "-";
  return new Date(value * 1000).toLocaleString();
}

async function loadList() {
  loading.value = true;
  try {
    const res: any = await getAdminNotificationList({
      page_num: page.value,
      page_size: pageSize,
      unread_only: unreadOnly.value,
    });
    list.value = res.data?.item ?? [];
    total.value = res.data?.total ?? 0;
  } finally {
    loading.value = false;
  }
}

async function handleMarkRead(row: NotificationItem) {
  if (row.read) return;
  await markAdminNotificationRead({ ids: [row.id] });
  row.read = true;
  notificationStore.decrementUnreadCount();
  if (unreadOnly.value) {
    await loadList();
  }
}

async function handleMarkAllRead() {
  await markAllAdminNotificationRead();
  notificationStore.clearUnreadCount();
  await loadList();
}

watch(unreadOnly, () => {
  page.value = 1;
  loadList();
});

onMounted(loadList);
</script>

<style scoped>
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.pager {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
