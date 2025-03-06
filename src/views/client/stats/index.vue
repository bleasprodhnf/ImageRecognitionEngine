<template>
  <div class="client-stats">
    <el-row :gutter="20">
      <!-- API调用统计 -->
      <el-col :span="24">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>API调用统计</span>
              <div class="header-controls">
                <el-select v-model="apiTimeRange" placeholder="时间范围">
                  <el-option label="最近7天" value="7" />
                  <el-option label="最近30天" value="30" />
                  <el-option label="最近90天" value="90" />
                </el-select>
              </div>
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="apiChartOption" autoresize />
          </div>
          <div class="stats-summary">
            <div class="summary-item">
              <h4>总调用次数</h4>
              <p>{{ apiSummary.total }}</p>
            </div>
            <div class="summary-item">
              <h4>平均响应时间</h4>
              <p>{{ apiSummary.avgResponseTime }}ms</p>
            </div>
            <div class="summary-item">
              <h4>成功率</h4>
              <p>{{ apiSummary.successRate }}%</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 存储空间统计 -->
      <el-col :span="24" class="mt-20">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>存储空间统计</span>
              <div class="header-controls">
                <el-select v-model="storageTimeRange" placeholder="时间范围">
                  <el-option label="最近7天" value="7" />
                  <el-option label="最近30天" value="30" />
                  <el-option label="最近90天" value="90" />
                </el-select>
              </div>
            </div>
          </template>
          <div class="chart-container">
            <v-chart class="chart" :option="storageChartOption" autoresize />
          </div>
          <div class="stats-summary">
            <div class="summary-item">
              <h4>当前使用量</h4>
              <p>{{ storageSummary.current }}GB</p>
            </div>
            <div class="summary-item">
              <h4>平均增长率</h4>
              <p>{{ storageSummary.growthRate }}%/月</p>
            </div>
            <div class="summary-item">
              <h4>预计可用天数</h4>
              <p>{{ storageSummary.remainingDays }}天</p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 使用记录 -->
      <el-col :span="24" class="mt-20">
        <el-card class="stats-card">
          <template #header>
            <div class="card-header">
              <span>使用记录</span>
              <el-button type="primary" size="small" @click="exportRecord">导出记录</el-button>
            </div>
          </template>
          <el-table :data="usageRecords" style="width: 100%">
            <el-table-column prop="date" label="日期" width="180" />
            <el-table-column prop="type" label="类型" width="120" />
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="value" label="数值" width="120" />
          </el-table>
          <div class="pagination-container">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="totalRecords"
              layout="total, prev, pager, next"
              @current-change="handlePageChange"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { ElMessage } from 'element-plus'

use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

// API统计数据
const apiTimeRange = ref('7')
const apiSummary = ref({
  total: 158000,
  avgResponseTime: 150,
  successRate: 99.5
})

// 存储统计数据
const storageTimeRange = ref('7')
const storageSummary = ref({
  current: 128,
  growthRate: 15,
  remainingDays: 180
})

// 使用记录分页
const currentPage = ref(1)
const pageSize = ref(10)
const totalRecords = ref(100)

// API调用图表配置
const apiChartOption = ref({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['调用次数', '响应时间']
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
    data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  },
  yAxis: [
    {
      type: 'value',
      name: '调用次数'
    },
    {
      type: 'value',
      name: '响应时间(ms)',
      splitLine: {
        show: false
      }
    }
  ],
  series: [
    {
      name: '调用次数',
      type: 'line',
      data: [15000, 18000, 22000, 20000, 25000, 21000, 19000],
      smooth: true,
      areaStyle: {}
    },
    {
      name: '响应时间',
      type: 'line',
      yAxisIndex: 1,
      data: [145, 152, 148, 155, 149, 146, 150],
      smooth: true
    }
  ]
})

// 存储使用图表配置
const storageChartOption = ref({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['存储使用量']
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
    data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
  },
  yAxis: {
    type: 'value',
    name: '存储量(GB)'
  },
  series: [
    {
      name: '存储使用量',
      type: 'line',
      data: [110, 115, 118, 122, 125, 127, 128],
      smooth: true,
      areaStyle: {}
    }
  ]
})

// 使用记录数据
const usageRecords = ref([
  {
    date: '2024-01-20 10:15:30',
    type: 'API调用',
    description: '模型推理请求',
    value: '1000次'
  },
  {
    date: '2024-01-20 09:45:20',
    type: '存储空间',
    description: '模型文件上传',
    value: '2.5GB'
  }
])

// 处理时间范围变化
const handleTimeRangeChange = () => {
  // TODO: 根据时间范围更新图表数据
}

// 处理分页变化
const handlePageChange = (page) => {
  currentPage.value = page
  // TODO: 获取对应页码的记录
}

// 导出记录
const exportRecord = () => {
  // TODO: 实现导出记录功能
  ElMessage.success('记录导出成功')
}

onMounted(() => {
  // TODO: 从API获取实际数据
})
</script>

<style scoped>
.client-stats {
  padding: 20px;
}

.stats-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-controls {
  display: flex;
  gap: 10px;
}

.mt-20 {
  margin-top: 20px;
}

.chart-container {
  height: 400px;
  margin: 20px 0;
}

.chart {
  height: 100%;
}

.stats-summary {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
  text-align: center;
}

.summary-item h4 {
  margin-bottom: 10px;
  color: #606266;
}

.summary-item p {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  margin: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>