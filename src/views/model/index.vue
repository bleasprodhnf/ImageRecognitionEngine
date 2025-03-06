<template>
  <div class="model-management">
    <!-- 模型版本管理 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>模型版本管理</span>
          <el-button type="primary" size="small">新增版本</el-button>
        </div>
      </template>
      <el-table :data="modelVersions" style="width: 100%">
        <el-table-column prop="version" label="版本号" width="120" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="releaseDate" label="发布日期" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="accuracy" label="准确率" width="100">
          <template #default="{ row }">
            {{ (row.accuracy * 100).toFixed(2) }}%
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link>查看</el-button>
            <el-button type="warning" link>编辑</el-button>
            <el-button type="danger" link>删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 参数配置 -->
    <el-card class="mb-4">
      <template #header>
        <div class="card-header">
          <span>参数配置</span>
          <el-button type="primary" size="small">保存配置</el-button>
        </div>
      </template>
      <el-form :model="currentParams" label-width="120px">
        <el-form-item v-for="param in modelParams" :key="param.id" :label="param.description">
          <el-input v-model="currentParams[param.name]" :type="param.type" />
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 性能监控 -->
    <el-card>
      <template #header>
        <div class="card-header">
          <span>性能监控</span>
          <el-radio-group v-model="monitorTimeRange" size="small">
            <el-radio-button label="day">24小时</el-radio-button>
            <el-radio-button label="week">7天</el-radio-button>
            <el-radio-button label="month">30天</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div class="monitor-charts">
        <el-row :gutter="20">
          <el-col :span="12">
            <div ref="performanceChart" style="height: 300px"></div>
          </el-col>
          <el-col :span="12">
            <div ref="resourceChart" style="height: 300px"></div>
          </el-col>
        </el-row>
        <el-row class="mt-4">
          <el-col :span="24">
            <div class="statistics-cards">
              <el-card shadow="hover" class="stat-card">
                <template #header>总调用次数</template>
                <div class="stat-value">{{ modelMonitorData.calls.total }}</div>
              </el-card>
              <el-card shadow="hover" class="stat-card">
                <template #header>成功调用</template>
                <div class="stat-value success">{{ modelMonitorData.calls.success }}</div>
              </el-card>
              <el-card shadow="hover" class="stat-card">
                <template #header>失败调用</template>
                <div class="stat-value error">{{ modelMonitorData.calls.failed }}</div>
              </el-card>
              <el-card shadow="hover" class="stat-card">
                <template #header>错误率</template>
                <div class="stat-value">{{ modelMonitorData.calls.errorRate }}%</div>
              </el-card>
            </div>
          </el-col>
        </el-row>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import * as echarts from 'echarts'
import { modelVersions, modelParams, modelMonitorData } from '@/api/model'

// 状态管理
const monitorTimeRange = ref('day')
const currentParams = reactive({})

// 初始化参数配置
modelParams.forEach(param => {
  currentParams[param.name] = param.value
})

// 状态标签配置
const getStatusType = (status) => {
  const types = {
    production: 'success',
    testing: 'warning',
    development: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    production: '生产中',
    testing: '测试中',
    development: '开发中'
  }
  return texts[status] || status
}

// 图表引用
const performanceChart = ref(null)
const resourceChart = ref(null)

// 初始化图表
onMounted(() => {
  // 性能图表
  const perfChart = echarts.init(performanceChart.value)
  perfChart.setOption({
    title: { text: '性能监控' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['响应时间', '吞吐量'] },
    xAxis: {
      type: 'category',
      data: modelMonitorData.performance.responseTime.map(item => item.time)
    },
    yAxis: [
      { type: 'value', name: '响应时间(ms)' },
      { type: 'value', name: '吞吐量(次/分)' }
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
  })

  // 资源使用图表
  const resChart = echarts.init(resourceChart.value)
  resChart.setOption({
    title: { text: '资源使用' },
    tooltip: { trigger: 'axis' },
    legend: { data: ['CPU使用率', '内存使用率'] },
    xAxis: {
      type: 'category',
      data: modelMonitorData.resources.cpu.map(item => item.time)
    },
    yAxis: { type: 'value', name: '使用率(%)' },
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
  })

  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    perfChart.resize()
    resChart.resize()
  })
})
</script>

<style scoped>
.model-management {
  padding: 20px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mt-4 {
  margin-top: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.statistics-cards {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}

.stat-value.success {
  color: #67C23A;
}

.stat-value.error {
  color: #F56C6C;
}
</style>