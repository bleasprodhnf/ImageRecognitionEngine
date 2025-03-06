<template>
  <div class="resources-container">
    <el-row :gutter="20">
      <!-- API调用统计 -->
      <el-col :span="24">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>API调用统计</span>
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                size="small"
              />
            </div>
          </template>
          <div class="chart-container">
            <el-row :gutter="20">
              <el-col :span="16">
                <div ref="apiChartRef" style="height: 300px"></div>
              </el-col>
              <el-col :span="8">
                <div class="stats-info">
                  <h3>调用量排名</h3>
                  <div v-for="(customer, index) in stats.topCustomers" :key="customer.name" class="rank-item">
                    <span class="rank-number">{{ index + 1 }}</span>
                    <span class="rank-name">{{ customer.name }}</span>
                    <span class="rank-value">{{ customer.apiCalls.toLocaleString() }}次</span>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>

      <!-- 存储使用统计 -->
      <el-col :span="24" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>存储使用统计</span>
            </div>
          </template>
          <div class="chart-container">
            <el-row :gutter="20">
              <el-col :span="16">
                <div ref="storageChartRef" style="height: 300px"></div>
              </el-col>
              <el-col :span="8">
                <div class="stats-info">
                  <h3>存储排名</h3>
                  <div v-for="(customer, index) in stats.topCustomers" :key="customer.name" class="rank-item">
                    <span class="rank-number">{{ index + 1 }}</span>
                    <span class="rank-name">{{ customer.name }}</span>
                    <span class="rank-value">{{ customer.storage }}GB</span>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts'
import { resourceStats } from '@/api/customer'

// 日期范围
const dateRange = ref([])

// 统计数据
const stats = ref(resourceStats)

// 图表实例
let apiChart = null
let storageChart = null

// 图表DOM引用
const apiChartRef = ref(null)
const storageChartRef = ref(null)

// 初始化API调用图表
const initApiChart = () => {
  apiChart = echarts.init(apiChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: stats.value.apiUsage.map(item => item.date)
    },
    yAxis: {
      type: 'value',
      name: 'API调用次数'
    },
    series: [
      {
        data: stats.value.apiUsage.map(item => item.calls),
        type: 'line',
        smooth: true,
        areaStyle: {}
      }
    ]
  }
  apiChart.setOption(option)
}

// 初始化存储使用图表
const initStorageChart = () => {
  storageChart = echarts.init(storageChartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: stats.value.storageUsage.map(item => item.date)
    },
    yAxis: {
      type: 'value',
      name: '存储使用量(GB)'
    },
    series: [
      {
        data: stats.value.storageUsage.map(item => item.usage),
        type: 'line',
        smooth: true,
        areaStyle: {}
      }
    ]
  }
  storageChart.setOption(option)
}

// 监听窗口大小变化
const handleResize = () => {
  apiChart?.resize()
  storageChart?.resize()
}

onMounted(() => {
  initApiChart()
  initStorageChart()
  window.addEventListener('resize', handleResize)
})
</script>

<style scoped>
.resources-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mt-20 {
  margin-top: 20px;
}

.chart-container {
  padding: 20px 0;
}

.stats-info {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.stats-info h3 {
  margin: 0 0 20px 0;
  font-size: 16px;
  color: #303133;
}

.rank-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.rank-number {
  width: 24px;
  height: 24px;
  line-height: 24px;
  text-align: center;
  background-color: #e6e8eb;
  border-radius: 50%;
  margin-right: 10px;
  font-size: 14px;
  color: #606266;
}

.rank-name {
  flex: 1;
  color: #606266;
}

.rank-value {
  color: #409eff;
  font-weight: bold;
}
</style>