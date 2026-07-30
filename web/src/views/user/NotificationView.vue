<template>
  <div class="notification-page">
    <div class="page-head">
      <h2>{{ t("notificationPage.title") }}</h2>
      <div class="actions">
        <el-checkbox v-model="unreadOnly">{{ t("notificationPage.unreadOnly") }}</el-checkbox>
        <el-button size="small" :disabled="!list.length" @click="handleMarkAllRead">
          {{ t("notificationPage.markAllRead") }}
        </el-button>
      </div>
    </div>

    <el-empty v-if="!loading && !list.length" :description="t('notificationPage.empty')" />
    <el-skeleton v-if="loading" :rows="4" animated />
    <div v-else class="notice-list">
      <div
        v-for="item in list"
        :key="item.id"
        class="notice-item"
        :class="{ unread: !item.read }"
        @click="handleMarkRead(item)"
      >
        <div class="notice-main">
          <div class="notice-title">
            <span>{{ item.title }}</span>
            <el-tag v-if="!item.read" size="small" type="danger">
              {{ t("notificationPage.unread") }}
            </el-tag>
          </div>
          <p>{{ item.content }}</p>
        </div>
        <span class="notice-time">{{ formatTime(item.created_at) }}</span>
      </div>
    </div>

    <Pagination
      v-if="total > pageSize"
      v-model:page="page"
      :page-size="pageSize"
      :total="total"
      @change="loadList"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import Pagination from "@/components/common/Pagination.vue";
import {
  getNotificationList,
  markAllNotificationRead,
  markNotificationRead,
  type NotificationItem,
} from "@/api/notification";
import { useAppConfigStore } from "@/stores/appConfig";
import { useNotificationStore } from "@/stores/notification";

const list = ref<NotificationItem[]>([]);
const total = ref(0);
const page = ref(1);
const appConfig = useAppConfigStore();
const notificationStore = useNotificationStore();
const { t } = useI18n();
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
    const res: any = await getNotificationList({
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

async function handleMarkRead(item: NotificationItem) {
  if (item.read) return;
  await markNotificationRead({ ids: [item.id] });
  item.read = true;
  notificationStore.decrementUnreadCount();
  if (unreadOnly.value) {
    await loadList();
  }
}

async function handleMarkAllRead() {
  await markAllNotificationRead();
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
.notification-page {
  max-width: 960px;
}
.page-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}
.page-head h2 {
  margin: 0;
  font-size: 22px;
}
.actions {
  display: flex;
  align-items: center;
  gap: 12px;
}
.notice-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.notice-item {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 14px 16px;
  border: 1px solid #ebeef5;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
}
.notice-item.unread {
  border-color: #f56c6c;
  background: #fff7f7;
}
.notice-main {
  min-width: 0;
}
.notice-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-weight: 600;
}
.notice-main p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  word-break: break-word;
}
.notice-time {
  flex: 0 0 auto;
  color: #909399;
  font-size: 13px;
}
</style>
