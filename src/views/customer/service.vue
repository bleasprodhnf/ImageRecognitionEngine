<template>
  <div class="service-container">
    <el-row :gutter="20">
      <!-- 实时状态 -->
      <el-col :span="24">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>实时监控</span>
              <el-tag type="success">在线监控中</el-tag>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="6">
              <div class="monitor-item">
                <div class="monitor-title">当前在线用户</div>
                <div class="monitor-value primary">{{ monitor.status.activeUsers }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="monitor-item">
                <div class="monitor-title">API调用次数</div>
                <div class="monitor-value success">{{ monitor.status.apiCalls }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="monitor-item">
                <div class="monitor-title">错误率</div>
                <div class="monitor-value warning">{{ monitor.status.errorRate }}%</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="monitor-item">
                <div class="monitor-title">平均响应时间</div>
                <div class="monitor-value info">{{ monitor.status.avgResponseTime }}ms</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 今日统计 -->
      <el-col :span="12" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>今日统计</span>
            </div>
          </template>
          <div class="stats-content">
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="stats-item">
                  <div class="stats-label">总调用次数</div>
                  <div class="stats-value">{{ monitor.today.totalCalls.toLocaleString() }}</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="stats-item">
                  <div class="stats-label">总用户数</div>
                  <div class="stats-value">{{ monitor.today.totalUsers }}</div>
                </div>
              </el-col>
            </el-row>
            <el-row :gutter="20" class="mt-20">
              <el-col :span="12">
                <div class="stats-item">
                  <div class="stats-label">成功调用</div>
                  <div class="stats-value success">{{ monitor.today.successCalls.toLocaleString() }}</div>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="stats-item">
                  <div class="stats-label">失败调用</div>
                  <div class="stats-value danger">{{ monitor.today.failedCalls }}</div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>

      <!-- 异常记录 -->
      <el-col :span="12" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>异常记录</span>
            </div>
          </template>
          <el-table :data="monitor.errors" style="width: 100%">
            <el-table-column prop="customer" label="客户" />
            <el-table-column prop="type" label="异常类型" />
            <el-table-column prop="time" label="发生时间" width="160" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <!-- API调用统计 -->
      <el-col :span="12" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>API调用统计</span>
              <el-select v-model="apiChartPeriod" size="small">
                <el-option label="最近7天" value="7" />
                <el-option label="最近30天" value="30" />
              </el-select>
            </div>
          </template>
          <div class="chart-container">
            <v-chart :option="apiChartOption" autoresize />
          </div>
        </el-card>
      </el-col>

      <!-- 存储使用统计 -->
      <el-col :span="12" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>存储使用统计</span>
              <el-select v-model="storageChartPeriod" size="small">
                <el-option label="最近7天" value="7" />
                <el-option label="最近30天" value="30" />
              </el-select>
            </div>
          </template>
          <div class="chart-container">
            <v-chart :option="storageChartOption" autoresize />
          </div>
        </el-card>
      </el-col>

      <!-- 客户资源使用排名 -->
      <el-col :span="24" class="mt-20">
        <el-card class="box-card">
          <template #header>
            <div class="card-header">
              <span>客户资源使用排名</span>
            </div>
          </template>
          <el-table :data="resourceStats.topCustomers" style="width: 100%">
            <el-table-column prop="name" label="客户名称" />
            <el-table-column label="API调用次数">
              <template #default="{ row }">
                {{ row.apiCalls.toLocaleString() }}
              </template>
            </el-table-column>
            <el-table-column label="存储空间(GB)" prop="storage" />
            <el-table-column label="费用(元)" prop="cost" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import { serviceMonitor, resourceStats } from '@/api/customer'

use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

// 监控数据
const monitor = ref(serviceMonitor)

// 图表周期选择
const apiChartPeriod = ref('7')
const storageChartPeriod = ref('7')

// API调用统计图表配置
const apiChartOption = computed(() => ({
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
      smooth: true
    }
  ]
}))

// 存储使用统计图表配置
const storageChartOption = computed(() => ({
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
      smooth: true
    }
  ]
}))

// 状态类型映射
const getStatusType = (status) => {
  const types = {
    resolved: 'success',
    pending: 'warning',
    processing: 'info'
  }
  return types[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const texts = {
    resolved: '已解决',
    pending: '待处理',
    processing: '处理中'
  }
  return texts[status] || status
}
</script>

<style scoped>
.service-container {
  padding: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.monitor-item {
  text-align: center;
  padding: 20px;
}

.monitor-title {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.monitor-value {
  font-size: 24px;
  font-weight: bold;
}

.monitor-value.primary { color: #409EFF; }
.monitor-value.success { color: #67C23A; }
.monitor-value.warning { color: #E6A23C; }
.monitor-value.info { color: #909399; }
.monitor-value.danger { color: #F56C6C; }

.stats-content {
  padding: 20px 0;
}

.stats-item {
  text-align: center;
}

.stats-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
}

.stats-value {
  font-size: 20px;
  font-weight: bold;
}

.stats-value.success { color: #67C23A; }
.stats-value.danger { color: #F56C6C; }

.chart-container {
  height: 300px;
}
</style>