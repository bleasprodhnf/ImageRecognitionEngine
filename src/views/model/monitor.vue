<template>
  <div class="monitor-container">
    <h2>模型监控</h2>
    <el-card class="box-card">
      <div class="monitor-stats">
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>运行状态</span>
                </div>
              </template>
              <div class="status-value" :style="{ color: modelStatus.color }">{{ modelStatus.text }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>请求数</span>
                </div>
              </template>
              <div class="status-value">{{ requestCount }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>平均响应时间</span>
                </div>
              </template>
              <div class="status-value">{{ avgResponseTime }}ms</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>错误率</span>
                </div>
              </template>
              <div class="status-value" :style="{ color: errorRate > 5 ? '#F56C6C' : '#67C23A' }">{{ errorRate }}%</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <div class="monitor-charts">
        <el-tabs v-model="activeTab" @tab-click="handleTabChange">
          <el-tab-pane label="性能监控" name="performance">
            <div class="chart-container" ref="performanceChart"></div>
          </el-tab-pane>
          <el-tab-pane label="资源使用" name="resources">
            <div class="chart-container" ref="resourcesChart"></div>
          </el-tab-pane>
          <el-tab-pane label="调用统计" name="calls">
            <div class="chart-container" ref="callsChart"></div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { modelMonitorData } from '@/api/model'

// 模型监控数据
const modelStatus = ref({
  text: '正常',
  color: '#67C23A'
})
const requestCount = ref(modelMonitorData.calls.total)
const avgResponseTime = ref(modelMonitorData.performance.responseTime[3].value)
const errorRate = ref(modelMonitorData.calls.errorRate)
const activeTab = ref('performance')

// 图表引用
const performanceChart = ref(null)
const resourcesChart = ref(null)
const callsChart = ref(null)

// 图表实例
let charts = {
  performance: null,
  resources: null,
  calls: null
}

// 初始化性能监控图表
const initPerformanceChart = () => {
  if (!performanceChart.value) return
  
  charts.performance = echarts.init(performanceChart.value)
  const option = {
    title: {
      text: '性能监控',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['响应时间', '吞吐量'],
      top: 30
    },
    xAxis: {
      type: 'category',
      data: modelMonitorData.performance.responseTime.map(item => item.time)
    },
    yAxis: [
      {
        type: 'value',
        name: '响应时间(ms)',
        position: 'left'
      },
      {
        type: 'value',
        name: '吞吐量(次/分)',
        position: 'right'
      }
    ],
    series: [
      {
        name: '响应时间',
        type: 'line',
        data: modelMonitorData.performance.responseTime.map(item => item.value)
      },
      {
        name: '吞吐量',
        type: 'line',
        yAxisIndex: 1,
        data: modelMonitorData.performance.throughput.map(item => item.value)
      }
    ]
  }
  charts.performance.setOption(option)
}

// 初始化资源使用图表
const initResourcesChart = () => {
  if (!resourcesChart.value) return

  charts.resources = echarts.init(resourcesChart.value)
  const option = {
    title: {
      text: '资源使用',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis'
    },
    legend: {
      data: ['CPU使用率', '内存使用率'],
      top: 30
    },
    xAxis: {
      type: 'category',
      data: modelMonitorData.resources.cpu.map(item => item.time)
    },
    yAxis: {
      type: 'value',
      name: '使用率(%)',
      max: 100
    },
    series: [
      {
        name: 'CPU使用率',
        type: 'line',
        data: modelMonitorData.resources.cpu.map(item => item.value)
      },
      {
        name: '内存使用率',
        type: 'line',
        data: modelMonitorData.resources.memory.map(item => item.value)
      }
    ]
  }
  charts.resources.setOption(option)
}

// 初始化调用统计图表
const initCallsChart = () => {
  if (!callsChart.value) return

  charts.calls = echarts.init(callsChart.value)
  const option = {
    title: {
      text: '调用统计',
      left: 'center'
    },
    tooltip: {
      trigger: 'item'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      top: 'center'
    },
    series: [
      {
        name: '调用情况',
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false,
          position: 'center'
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '20',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data: [
          { value: modelMonitorData.calls.success, name: '成功调用' },
          { value: modelMonitorData.calls.failed, name: '失败调用' }
        ]
      }
    ]
  }
  charts.calls.setOption(option)
}

// 处理标签页切换
const handleTabChange = (tab) => {
  nextTick(() => {
    switch (tab.props.name) {
      case 'performance':
        initPerformanceChart()
        break
      case 'resources':
        initResourcesChart()
        break
      case 'calls':
        initCallsChart()
        break
    }
  })
}

// 监听窗口大小变化
const handleResize = () => {
  Object.values(charts).forEach(chart => {
    chart && chart.resize()
  })
}

onMounted(() => {
  // 初始化默认标签页的图表
  initPerformanceChart()
  // 添加窗口大小变化监听
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  // 移除窗口大小变化监听
  window.removeEventListener('resize', handleResize)
  // 销毁图表实例
  Object.values(charts).forEach(chart => {
    chart && chart.dispose()
  })
})
</script>

<style scoped>
.monitor-container {
  padding: 20px;
}

.monitor-stats {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
}

.chart-container {
  height: 400px;
  margin-top: 20px;
}
</style>