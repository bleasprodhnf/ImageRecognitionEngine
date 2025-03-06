<template>
  <div class="system-stats-container">
    <h2>系统运行统计</h2>
    <el-card class="box-card">
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>系统运行时间</span>
                </div>
              </template>
              <div class="stats-value">{{ uptime }}天</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>总请求量</span>
                </div>
              </template>
              <div class="stats-value">{{ totalRequests }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>系统负载</span>
                </div>
              </template>
              <div class="stats-value">{{ systemLoad }}%</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>网络流量</span>
                </div>
              </template>
              <div class="stats-value">{{ networkTraffic }}MB/s</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <div class="stats-charts">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="请求统计" name="requests">
            <div class="chart-container" ref="requestsChart"></div>
          </el-tab-pane>
          <el-tab-pane label="性能统计" name="performance">
            <div class="chart-container" ref="performanceChart"></div>
          </el-tab-pane>
          <el-tab-pane label="错误统计" name="errors">
            <div class="chart-container" ref="errorsChart"></div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { systemMonitor } from '@/api/mock'
import * as echarts from 'echarts'

// 系统统计数据
const uptime = ref(30) // 模拟30天运行时间
const totalRequests = ref(15000) // 模拟总请求量
const systemLoad = ref(systemMonitor.cpu.usage) // 使用CPU使用率作为系统负载
const networkTraffic = ref(systemMonitor.network.upload + systemMonitor.network.download) // 网络流量
const activeTab = ref('requests')

// 图表引用
const requestsChart = ref(null)
const performanceChart = ref(null)
const errorsChart = ref(null)

// 初始化请求统计图表
const initRequestsChart = () => {
  const chart = echarts.init(requestsChart.value)
  const option = {
    title: {
      text: '近7天请求量统计'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      data: [820, 932, 901, 934, 1290, 1330, 1320],
      type: 'line',
      smooth: true
    }]
  }
  chart.setOption(option)
}

// 初始化性能统计图表
const initPerformanceChart = () => {
  const chart = echarts.init(performanceChart.value)
  const option = {
    title: {
      text: '系统性能指标'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['CPU使用率', '内存使用率', '磁盘使用率']
    },
    xAxis: {
      type: 'category',
      data: ['00:00', '04:00', '08:00', '12:00', '16:00', '20:00', '24:00']
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [
      {
        name: 'CPU使用率',
        type: 'line',
        data: [45, 52, 65, 70, 55, 48, 42]
      },
      {
        name: '内存使用率',
        type: 'line',
        data: [50, 55, 58, 62, 60, 57, 52]
      },
      {
        name: '磁盘使用率',
        type: 'line',
        data: [48, 48, 49, 50, 50, 51, 51]
      }
    ]
  }
  chart.setOption(option)
}

// 初始化错误统计图表
const initErrorsChart = () => {
  const chart = echarts.init(errorsChart.value)
  const option = {
    title: {
      text: '系统错误统计'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['错误数', '错误率']
    },
    xAxis: {
      type: 'category',
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: [
      {
        type: 'value',
        name: '错误数',
        position: 'left'
      },
      {
        type: 'value',
        name: '错误率',
        position: 'right',
        axisLabel: {
          formatter: '{value}%'
        }
      }
    ],
    series: [
      {
        name: '错误数',
        type: 'bar',
        data: [25, 18, 22, 15, 20, 16, 19]
      },
      {
        name: '错误率',
        type: 'line',
        yAxisIndex: 1,
        data: [2.5, 1.8, 2.2, 1.5, 2.0, 1.6, 1.9]
      }
    ]
  }
  chart.setOption(option)
}

onMounted(() => {
  // 初始化图表
  initRequestsChart()
  initPerformanceChart()
  initErrorsChart()
})
</script>

<style scoped>
.system-stats-container {
  padding: 20px;
}

.box-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
  text-align: center;
  margin-top: 10px;
}

.chart-container {
  height: 400px;
  margin-top: 20px;
}
</style>