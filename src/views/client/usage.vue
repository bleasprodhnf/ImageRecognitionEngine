<template>
  <div class="usage-statistics">
    <!-- API调用统计 -->
    <el-card class="usage-card">
      <template #header>
        <div class="card-header">
          <span>API调用统计</span>
          <div class="header-controls">
            <el-date-picker
              v-model="apiDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              @change="updateApiChart"
            />
          </div>
        </div>
      </template>
      <div class="chart-container">
        <div class="chart-wrapper" ref="apiChartRef"></div>
        <div class="stats-summary">
          <div class="stat-item">
            <div class="stat-label">今日调用次数</div>
            <div class="stat-value">{{ todayStats.totalCalls.toLocaleString() }}</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">成功率</div>
            <div class="stat-value">{{ ((todayStats.successCalls / todayStats.totalCalls) * 100).toFixed(2) }}%</div>
          </div>
          <div class="stat-item">
            <div class="stat-label">平均响应时间</div>
            <div class="stat-value">{{ realtimeStatus.avgResponseTime }}ms</div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 存储空间使用统计 -->
    <el-card class="usage-card">
      <template #header>
        <div class="card-header">
          <span>存储空间使用统计</span>
          <div class="header-controls">
            <el-date-picker
              v-model="storageDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              @change="updateStorageChart"
            />
          </div>
        </div>
      </template>
      <div class="chart-container">
        <div class="chart-wrapper" ref="storageChartRef"></div>
        <div class="storage-info">
          <el-progress
            type="dashboard"
            :percentage="(currentStorage / storageLimit) * 100"
            :color="storageProgressColor"
          >
            <template #default="{ percentage }">
              <div class="storage-text">
                <span class="current-storage">{{ currentStorage }}GB</span>
                <span class="total-storage">共{{ storageLimit }}GB</span>
              </div>
            </template>
          </el-progress>
        </div>
      </div>
    </el-card>

    <!-- 使用记录 -->
    <el-card class="usage-card">
      <template #header>
        <div class="card-header">
          <span>最近调用记录</span>
          <el-button type="primary" link @click="refreshUsageRecords">刷新</el-button>
        </div>
      </template>
      <el-table :data="usageRecords" style="width: 100%">
        <el-table-column prop="time" label="时间" width="180" />
        <el-table-column prop="type" label="操作类型" width="120" />
        <el-table-column prop="details" label="详情" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="table-pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import * as echarts from 'echarts'
import { serviceMonitor, resourceStats } from '@/api/customer'

// 图表实例
const apiChartRef = ref(null)
const storageChartRef = ref(null)
let apiChart = null
let storageChart = null

// 日期范围
const apiDateRange = ref([])
const storageDateRange = ref([])

// 实时状态和今日统计
const realtimeStatus = ref(serviceMonitor.status)
const todayStats = ref(serviceMonitor.today)

// 存储空间信息
const currentStorage = ref(950)
const storageLimit = ref(1024)
const storageProgressColor = computed(() => {
  const percentage = (currentStorage.value / storageLimit.value) * 100
  if (percentage < 70) return '#67C23A'
  if (percentage < 90) return '#E6A23C'
  return '#F56C6C'
})

// 使用记录分页
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const usageRecords = ref([])

// 初始化图表
const initCharts = () => {
  // API调用图表
  apiChart = echarts.init(apiChartRef.value)
  const apiOption = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: resourceStats.apiUsage.map(item => item.date)
    },
    yAxis: {
      type: 'value',
      name: '调用次数'
    },
    series: [
      {
        data: resourceStats.apiUsage.map(item => item.calls),
        type: 'line',
        smooth: true,
        areaStyle: {
          opacity: 0.3
        }
      }
    ]
  }
  apiChart.setOption(apiOption)

  // 存储空间图表
  storageChart = echarts.init(storageChartRef.value)
  const storageOption = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: resourceStats.storageUsage.map(item => item.date)
    },
    yAxis: {
      type: 'value',
      name: '存储空间(GB)'
    },
    series: [
      {
        data: resourceStats.storageUsage.map(item => item.usage),
        type: 'line',
        smooth: true,
        areaStyle: {
          opacity: 0.3
        }
      }
    ]
  }
  storageChart.setOption(storageOption)
}

// 更新图表
const updateApiChart = () => {
  // 实际项目中这里需要根据日期范围请求后端数据
  console.log('更新API调用图表', apiDateRange.value)
}

const updateStorageChart = () => {
  // 实际项目中这里需要根据日期范围请求后端数据
  console.log('更新存储空间图表', storageDateRange.value)
}

// 刷新使用记录
const refreshUsageRecords = () => {
  // 模拟数据
  usageRecords.value = [
    {
      time: '2024-01-20 15:30:00',
      type: 'API调用',
      details: '模型预测请求',
      status: 'success'
    },
    {
      time: '2024-01-20 15:28:30',
      type: '文件上传',
      details: '上传训练数据',
      status: 'success'
    },
    {
      time: '2024-01-20 15:25:00',
      type: 'API调用',
      details: '数据处理请求',
      status: 'failed'
    }
  ]
  total.value = 100 // 模拟总记录数
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  refreshUsageRecords()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  refreshUsageRecords()
}

onMounted(() => {
  initCharts()
  refreshUsageRecords()

  // 监听窗口大小变化，重绘图表
  window.addEventListener('resize', () => {
    apiChart?.resize()
    storageChart?.resize()
  })
})
</script>

<style scoped>
.usage-statistics {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.usage-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  padding: 20px 0;
}

.chart-wrapper {
  height: 300px;
  width: 100%;
}

.stats-summary {
  display: flex;
  justify-content: space-around;
  margin-top: 20px;
  flex-wrap: wrap;
  gap: 20px;
}

.stat-item {
  text-align: center;
  flex: 1;
  min-width: 200px;
}

.stat-label {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1890ff;
}

@media screen and (max-width: 768px) {
  .usage-statistics {
    padding: 0 10px;
  }
  
  .stat-item {
    min-width: 150px;
  }
}
.storage-info {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.storage-text {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.current-storage {
  font-size: 20px;
  font-weight: bold;
  color: #303133;
}

.total-storage {
  font-size: 14px;
  color: #909399;
  margin-top: 4px;
}

.table-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>