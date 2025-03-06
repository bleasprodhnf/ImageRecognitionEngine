<template>
  <div class="client-dashboard">
    <el-row :gutter="20">
      <!-- API调用统计 -->
      <el-col :span="12">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>API调用统计</span>
            </div>
          </template>
          <div class="api-stats">
            <div class="stat-item">
              <h4>今日调用次数</h4>
              <p>{{ apiStats.today }}</p>
            </div>
            <div class="stat-item">
              <h4>本月累计调用</h4>
              <p>{{ apiStats.monthly }}</p>
            </div>
            <div class="stat-item">
              <h4>剩余可用次数</h4>
              <p>{{ apiStats.remaining }}</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 存储空间使用情况 -->
      <el-col :span="12">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>存储空间使用情况</span>
            </div>
          </template>
          <div class="storage-stats">
            <el-progress type="dashboard" :percentage="storageUsage.percentage" />
            <div class="storage-info">
              <p>已使用：{{ storageUsage.used }}GB</p>
              <p>总容量：{{ storageUsage.total }}GB</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 套餐信息 -->
      <el-col :span="24" class="mt-20">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>当前套餐信息</span>
              <el-button type="primary" size="small" @click="handleUpgrade">升级套餐</el-button>
            </div>
          </template>
          <div class="package-info">
            <h3>{{ packageInfo.name }}</h3>
            <div class="package-details">
              <p><el-icon><Timer /></el-icon> 到期时间：{{ packageInfo.expireDate }}</p>
              <p><el-icon><Connection /></el-icon> 并发数限制：{{ packageInfo.limits.concurrency }}</p>
              <p><el-icon><DataLine /></el-icon> API限制：{{ packageInfo.features[0] }}</p>
              <p><el-icon><Files /></el-icon> 存储限制：{{ packageInfo.limits.storage }}GB</p>
            </div>
            <div class="package-features">
              <h4>套餐特权：</h4>
              <ul>
                <li v-for="(feature, index) in packageInfo.features" :key="index">{{ feature }}</li>
              </ul>
            </div>
            <div class="usage-progress">
              <div class="progress-item">
                <span>API使用量</span>
                <el-progress :percentage="apiUsagePercentage" :format="percentageFormat" />
              </div>
              <div class="progress-item">
                <span>存储空间</span>
                <el-progress :percentage="storageUsage.percentage" :format="percentageFormat" />
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- API调用趋势 -->
      <el-col :span="24" class="mt-20">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>API调用趋势</span>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                @change="handleDateRangeChange"
              />
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="chartOption" autoresize />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Timer, Connection, DataLine, Files } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart, { THEME_KEY } from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

import { packages, serviceMonitor, resourceStats } from '@/api/customer'

// API统计数据
const apiStats = ref({
  today: serviceMonitor.today.totalCalls,
  monthly: resourceStats.apiUsage.reduce((total, item) => total + item.calls, 0),
  remaining: packages[0].limits.storage * 1024 - serviceMonitor.today.totalCalls // 假设当前是企业版
})

// 存储使用情况
const storageUsage = ref({
  percentage: Math.round((resourceStats.storageUsage[resourceStats.storageUsage.length - 1].usage / packages[0].limits.storage) * 100),
  used: resourceStats.storageUsage[resourceStats.storageUsage.length - 1].usage,
  total: packages[0].limits.storage
})

// 当前套餐信息
const packageInfo = ref({
  ...packages[0],
  expireDate: '2024-12-31', // 这个日期应该从后端获取
  features: packages[0].features,
  apiUsage: serviceMonitor.today.totalCalls,
  storageUsage: resourceStats.storageUsage[resourceStats.storageUsage.length - 1].usage
})

const dateRange = ref([])
const apiUsageData = ref([])

// 计算API使用百分比
const apiUsagePercentage = computed(() => {
  return Math.round((apiStats.value.monthly / (apiStats.value.monthly + apiStats.value.remaining)) * 100)
})

// 格式化百分比显示
const percentageFormat = (percentage) => `${percentage}%`

// 图表配置
const chartOption = ref({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['API调用次数']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: []
  },
  yAxis: {
    type: 'value'
  },
  series: [
    {
      name: 'API调用次数',
      type: 'line',
      data: [],
      smooth: true,
      areaStyle: {}
    }
  ]
})

// 处理日期范围变化
const handleDateRangeChange = () => {
  // TODO: 根据日期范围获取API调用数据
  updateChartData()
}

// 更新图表数据
const updateChartData = () => {
  // 模拟数据
  const dates = ['2024-01-14', '2024-01-15', '2024-01-16', '2024-01-17', '2024-01-18', '2024-01-19', '2024-01-20']
  const calls = [95000, 98000, 120000, 115000, 125000, 110000, 105000]
  
  chartOption.value.xAxis.data = dates
  chartOption.value.series[0].data = calls
}

// 处理套餐升级
const handleUpgrade = () => {
  // 显示套餐升级对话框
  ElMessageBox.confirm(
    '升级套餐可获得更多API调用次数和存储空间，是否查看套餐详情？',
    '升级套餐',
    {
      confirmButtonText: '查看详情',
      cancelButtonText: '取消',
      type: 'info'
    }
  ).then(() => {
    // 跳转到套餐页面
    router.push('/client/packages')
  }).catch(() => {
    // 用户取消操作
  })
}

onMounted(() => {
  // 初始化图表数据
  updateChartData()
  // TODO: 从API获取实际数据
})
</script>

<style scoped>
.client-dashboard {
  padding: 20px;
}

.dashboard-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mt-20 {
  margin-top: 20px;
}

.api-stats {
  display: flex;
  justify-content: space-around;
  text-align: center;
}

.stat-item h4 {
  margin-bottom: 10px;
  color: #606266;
}

.stat-item p {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  margin: 0;
}

.storage-stats {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.storage-info {
  margin-top: 20px;
  text-align: center;
}

.package-info {
  padding: 20px;
}

.package-info h3 {
  margin: 0 0 20px 0;
  color: #409EFF;
}

.package-details {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 15px;
}

.package-details p {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.usage-progress {
  margin-top: 20px;
}

.progress-item {
  margin: 15px 0;
}

.progress-item span {
  display: block;
  margin-bottom: 8px;
}

.chart-container {
  height: 400px;
}

.chart {
  height: 100%;
}
</style>