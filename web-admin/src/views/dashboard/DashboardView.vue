<template>
  <div class="dashboard-page">
    <el-row :gutter="16" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6" v-for="card in statCards" :key="card.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-body">
            <el-icon :size="34" :color="card.color">
              <component :is="card.icon" />
            </el-icon>
            <div>
              <div class="stat-value">{{ card.value }}</div>
              <div class="stat-label">{{ card.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="chart-card">
      <template #header>
        <div class="card-header">
          <span>{{ t("page.dashboard.trendTitle") }}</span>
          <el-button :loading="loading" size="small" @click="loadData">{{ t("common.refresh") }}</el-button>
        </div>
      </template>
      <v-chart class="trend-chart" :option="chartOption" autoresize />
    </el-card>

    <el-card>
      <template #header>{{ t("page.dashboard.pendingProducts") }}</template>
      <el-table :data="pendingProducts" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="t('page.product.name')" min-width="180" />
        <el-table-column prop="price" :label="t('page.product.price')" width="100" />
        <el-table-column prop="boss_name" :label="t('page.product.seller')" width="120" />
        <el-table-column :label="t('common.actions')" width="160">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="audit(row.id, 1)">{{ t("page.dashboard.approve") }}</el-button>
            <el-button size="small" type="danger" @click="audit(row.id, 2)">{{ t("page.dashboard.reject") }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import VChart from "vue-echarts";
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { LineChart } from "echarts/charts";
import { GridComponent, LegendComponent, TooltipComponent } from "echarts/components";
import {
  auditProduct,
  getAdminProductList,
  getStatsOrders,
  getStatsOverview,
} from "@/api";
import { t } from "@/locales";

use([CanvasRenderer, LineChart, GridComponent, LegendComponent, TooltipComponent]);

const loading = ref(false);
const pendingProducts = ref<any[]>([]);
const overview = ref({
  today_orders: 0,
  total_sales: 0,
  platform_revenue: 0,
  registered_users: 0,
});
const trend = ref({
  dates: [] as string[],
  order_counts: [] as number[],
  sales_amounts: [] as number[],
});

function buildDefaultTrend() {
  const dates: string[] = [];
  const orderCounts: number[] = [];
  const salesAmounts: number[] = [];
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  for (let offset = 6; offset >= 0; offset -= 1) {
    const day = new Date(today);
    day.setDate(today.getDate() - offset);
    const label = [
      day.getFullYear(),
      String(day.getMonth() + 1).padStart(2, "0"),
      String(day.getDate()).padStart(2, "0"),
    ].join("-");
    dates.push(label);
    orderCounts.push(0);
    salesAmounts.push(0);
  }
  return {
    dates,
    order_counts: orderCounts,
    sales_amounts: salesAmounts,
  };
}

function normalizeTrendData(raw: any) {
  if (
    !raw ||
    !Array.isArray(raw.dates) ||
    raw.dates.length === 0 ||
    !Array.isArray(raw.order_counts) ||
    !Array.isArray(raw.sales_amounts)
  ) {
    return buildDefaultTrend();
  }
  return {
    dates: raw.dates,
    order_counts: raw.order_counts,
    sales_amounts: raw.sales_amounts,
  };
}

const statCards = computed(() => [
  {
    label: t("page.dashboard.todayOrders"),
    value: overview.value.today_orders,
    icon: "Tickets",
    color: "#409eff",
  },
  {
    label: t("page.dashboard.totalSales"),
    value: `¥${Number(overview.value.total_sales || 0).toFixed(2)}`,
    icon: "Money",
    color: "#67c23a",
  },
  {
    label: t("page.dashboard.platformRevenue"),
    value: `¥${Number(overview.value.platform_revenue || 0).toFixed(2)}`,
    icon: "Coin",
    color: "#f56c6c",
  },
  {
    label: t("page.dashboard.registeredUsers"),
    value: overview.value.registered_users,
    icon: "User",
    color: "#909399",
  },
  {
    label: t("page.dashboard.pendingProducts"),
    value: pendingProducts.value.length,
    icon: "Goods",
    color: "#e6a23c",
  },
]);

const chartOption = computed(() => ({
  color: ["#409eff", "#67c23a"],
  tooltip: { trigger: "axis" },
  legend: { top: 0, data: [t("page.dashboard.orderCount"), t("page.dashboard.salesAmount")] },
  grid: { top: 42, left: 48, right: 24, bottom: 36 },
  xAxis: {
    type: "category",
    boundaryGap: false,
    data: trend.value.dates,
  },
  yAxis: [
    { type: "value", name: t("page.dashboard.orderCount"), minInterval: 1 },
    { type: "value", name: t("page.dashboard.salesAmount") },
  ],
  series: [
    {
      name: t("page.dashboard.orderCount"),
      type: "line",
      smooth: true,
      data: trend.value.order_counts,
    },
    {
      name: t("page.dashboard.salesAmount"),
      type: "line",
      smooth: true,
      yAxisIndex: 1,
      data: trend.value.sales_amounts,
    },
  ],
}));

async function loadData() {
  loading.value = true;
  try {
    const [overviewRes, trendRes, productRes]: any[] = await Promise.all([
      getStatsOverview(),
      getStatsOrders(),
      getAdminProductList({ page_num: 1, page_size: 100, audit_status: 0 }),
    ]);
    overview.value = overviewRes.data ?? overview.value;
    trend.value = normalizeTrendData(trendRes.data);
    pendingProducts.value = productRes.data?.item ?? [];
  } finally {
    loading.value = false;
  }
}

async function audit(id: number, status: number) {
  await auditProduct({ id, audit_status: status });
  ElMessage.success(t("page.dashboard.auditSuccess", {
    action: status === 1 ? t("page.dashboard.approve") : t("page.dashboard.reject"),
  }));
  loadData();
}

onMounted(loadData);
</script>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.stats-row {
  row-gap: 16px;
}
.stat-card {
  min-height: 96px;
}
.stat-body {
  display: flex;
  align-items: center;
  gap: 12px;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 32px;
}
.stat-label {
  color: #909399;
  font-size: 13px;
}
.chart-card {
  min-height: 360px;
}
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.trend-chart {
  height: 300px;
}
</style>
