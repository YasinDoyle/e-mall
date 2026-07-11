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
          <span>近 7 日订单趋势</span>
          <el-button :loading="loading" size="small" @click="loadData">刷新</el-button>
        </div>
      </template>
      <v-chart class="trend-chart" :option="chartOption" autoresize />
    </el-card>

    <el-card>
      <template #header>待审核商品</template>
      <el-table :data="pendingProducts" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="商品名称" min-width="180" />
        <el-table-column prop="price" label="价格" width="100" />
        <el-table-column prop="boss_name" label="卖家" width="120" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="audit(row.id, 1)">上架</el-button>
            <el-button size="small" type="danger" @click="audit(row.id, 2)">拒绝</el-button>
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

use([CanvasRenderer, LineChart, GridComponent, LegendComponent, TooltipComponent]);

const loading = ref(false);
const pendingProducts = ref<any[]>([]);
const overview = ref({
  today_orders: 0,
  total_sales: 0,
  registered_users: 0,
});
const trend = ref({
  dates: [] as string[],
  order_counts: [] as number[],
  sales_amounts: [] as number[],
});

const statCards = computed(() => [
  {
    label: "今日订单",
    value: overview.value.today_orders,
    icon: "Tickets",
    color: "#409eff",
  },
  {
    label: "累计销售额",
    value: `¥${Number(overview.value.total_sales || 0).toFixed(2)}`,
    icon: "Money",
    color: "#67c23a",
  },
  {
    label: "注册用户",
    value: overview.value.registered_users,
    icon: "User",
    color: "#909399",
  },
  {
    label: "待审核商品",
    value: pendingProducts.value.length,
    icon: "Goods",
    color: "#e6a23c",
  },
]);

const chartOption = computed(() => ({
  color: ["#409eff", "#67c23a"],
  tooltip: { trigger: "axis" },
  legend: { top: 0, data: ["订单数", "销售额"] },
  grid: { top: 42, left: 48, right: 24, bottom: 36 },
  xAxis: {
    type: "category",
    boundaryGap: false,
    data: trend.value.dates,
  },
  yAxis: [
    { type: "value", name: "订单数", minInterval: 1 },
    { type: "value", name: "销售额" },
  ],
  series: [
    {
      name: "订单数",
      type: "line",
      smooth: true,
      data: trend.value.order_counts,
    },
    {
      name: "销售额",
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
    trend.value = trendRes.data ?? trend.value;
    pendingProducts.value = productRes.data?.item ?? [];
  } finally {
    loading.value = false;
  }
}

async function audit(id: number, status: number) {
  await auditProduct({ id, audit_status: status });
  ElMessage.success(status === 1 ? "已上架" : "已拒绝");
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
