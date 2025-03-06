<template>
  <div class="system-monitor">
    <el-card class="monitor-header" v-loading="loading">
      <template #header>
        <div class="card-header">
          <h3>系统状态监控</h3>
          <div class="actions">
            <el-tooltip content="刷新数据" placement="top">
              <el-button type="primary" circle @click="refreshData" :loading="refreshing">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </el-tooltip>
            <el-switch
              v-model="autoRefresh"
              active-text="自动刷新"
              inactive-text="手动刷新"
              @change="handleAutoRefreshChange"
            />
          </div>
        </div>
      </template>
      <el-alert
        v-if="error"
        :title="error"
        type="error"
        show-icon
        :closable="false"
        style="margin-bottom: 15px"
      />
      <el-row :gutter="20">
        <!-- CPU状态 -->
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card">
            <div class="status-title">CPU</div>
            <el-progress
              type="dashboard"
              :percentage="Math.round(systemStatus.cpu.usage || 0)"
              :color="getStatusColor(systemStatus.cpu.status)"
            />
            <div class="status-details">
              <div class="detail-item">
                <span class="label">使用率:</span>
                <span class="value">{{ Math.round(systemStatus.cpu.usage || 0) }}%</span>
              </div>
              <div class="detail-item">
                <span class="label">核心数:</span>
                <span class="value">{{ systemStatus.cpu.cores || 0 }}</span>
              </div>
            </div>
          </div>
        </el-col>
        
        <!-- 内存状态 -->
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card">
            <div class="status-title">内存</div>
            <el-progress
              type="dashboard"
              :percentage="Math.round(systemStatus.memory.usage || 0)"
              :color="getStatusColor(systemStatus.memory.status)"
            />
            <div class="status-details">
              <div class="detail-item">
                <span class="label">已用:</span>
                <span class="value">{{ systemStatus.memory.used }}</span>
              </div>
              <div class="detail-item">
                <span class="label">可用:</span>
                <span class="value">{{ systemStatus.memory.free }}</span>
              </div>
            </div>
          </div>
        </el-col>
        
        <!-- 磁盘状态 -->
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card">
            <div class="status-title">磁盘</div>
            <el-progress
              type="dashboard"
              :percentage="Math.round(systemStatus.disk.usage || 0)"
              :color="getStatusColor(systemStatus.disk.status)"
            />
            <div class="status-details">
              <div class="detail-item">
                <span class="label">已用:</span>
                <span class="value">{{ systemStatus.disk.used }}</span>
              </div>
              <div class="detail-item">
                <span class="label">可用:</span>
                <span class="value">{{ systemStatus.disk.free }}</span>
              </div>
            </div>
          </div>
        </el-col>
        
        <!-- 网络状态 -->
        <el-col :xs="24" :sm="12" :md="6">
          <div class="status-card">
            <div class="status-title">网络</div>
            <div class="network-status">
              <div class="network-item">
                <el-icon :size="24" color="#409EFF"><Upload /></el-icon>
                <div class="network-value">{{ systemStatus.network.out }}</div>
                <div class="network-label">上传</div>
              </div>
              <div class="network-item">
                <el-icon :size="24" color="#67C23A"><Download /></el-icon>
                <div class="network-value">{{ systemStatus.network.in }}</div>
                <div class="network-label">下载</div>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
      
      <div class="system-uptime">
        <el-icon><Timer /></el-icon>
        <span>系统运行时间: {{ systemStatus.uptime || '未知' }}</span>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { monitorService } from '../services/monitorService'

// 组件属性
const props = defineProps({
  // 自动刷新间隔（毫秒）
  refreshInterval: {
    type: Number,
    default: 30000
  },
  // 最大重试次数
  maxRetries: {
    type: Number,
    default: 3
  },
  // 重试延迟（毫秒）
  retryDelay: {
    type: Number,
    default: 5000
  }
})

// 组件事件
const emit = defineEmits(['status-change', 'error'])

// 状态变量
const loading = ref(true)
const refreshing = ref(false)
const error = ref(null)
const autoRefresh = ref(true)
const systemStatus = ref({
  cpu: { usage: 0, cores: 0, status: 'normal' },
  memory: { total: '0 B', used: '0 B', free: '0 B', usage: 0, status: 'normal' },
  disk: { total: '0 B', used: '0 B', free: '0 B', usage: 0, status: 'normal' },
  network: { in: '0 B/s', out: '0 B/s' },
  uptime: '0分钟'
})

// 获取状态颜色
const getStatusColor = (status) => {
  const colors = {
    normal: '#67C23A',  // 绿色
    warning: '#E6A23C', // 黄色
    danger: '#F56C6C'   // 红色
  }
  return colors[status] || colors.normal
}

// 刷新数据
const refreshData = async () => {
  if (refreshing.value) return
  
  refreshing.value = true
  error.value = null
  
  try {
    const data = await monitorService.getSystemStatus()
    systemStatus.value = data
    emit('status-change', data)
  } catch (err) {
    error.value = err.message || '获取系统状态失败'
    emit('error', error.value)
    console.error('获取系统状态失败:', err)
  } finally {
    refreshing.value = false
    loading.value = false
  }
}

// 处理自动刷新变更
const handleAutoRefreshChange = (value) => {
  if (value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

// 开始自动刷新
const startAutoRefresh = () => {
  monitorService.startAutoRefresh(
    (result) => {
      if (result.success) {
        systemStatus.value = result.data
        error.value = null
        emit('status-change', result.data)
      } else {
        error.value = result.error
        emit('error', result.error)
      }
      loading.value = false
    },
    props.maxRetries,
    props.retryDelay
  )
}

// 停止自动刷新
const stopAutoRefresh = () => {
  monitorService.stopAutoRefresh()
}

// 组件挂载时
onMounted(() => {
  if (autoRefresh.value) {
    startAutoRefresh()
  } else {
    refreshData()
  }
})

// 组件卸载时
onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.system-monitor {
  margin-bottom: 20px;
}

.monitor-header {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-card {
  padding: 15px;
  text-align: center;
  border-radius: 4px;
  background-color: #f5f7fa;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 15px;
}

.status-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
  color: #303133;
}

.status-details {
  width: 100%;
  margin-top: 10px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  margin: 5px 0;
}

.label {
  color: #606266;
}

.value {
  font-weight: bold;
  color: #303133;
}

.network-status {
  display: flex;
  justify-content: space-around;
  width: 100%;
  margin: 15px 0;
}

.network-item {
  text-align: center;
}

.network-value {
  font-size: 16px;
  font-weight: bold;
  margin: 5px 0;
  color: #303133;
}

.network-label {
  font-size: 12px;
  color: #606266;
}

.system-uptime {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 15px;
  color: #606266;
}

.system-uptime .el-icon {
  margin-right: 5px;
}
</style>