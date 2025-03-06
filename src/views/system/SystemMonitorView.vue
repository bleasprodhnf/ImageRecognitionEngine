<template>
  <div class="system-monitor-container">
    <div class="monitor-header">
      <h1>系统监控面板</h1>
      <div class="service-status" :class="serviceStatusClass">
        服务状态: {{ serviceStatusText }}
        <span v-if="retryInfo" class="retry-info">{{ retryInfo }}</span>
      </div>
    </div>

    <div v-if="error" class="error-message">
      <i class="error-icon">⚠️</i>
      <span>{{ error }}</span>
      <button v-if="canRetry" @click="retryConnection" class="retry-button">重试连接</button>
    </div>

    <div class="monitor-grid" v-if="!error || showDataWhileError">
      <!-- CPU 监控卡片 -->
      <div class="monitor-card" :class="getStatusClass(systemStatus?.cpu?.status)">
        <div class="card-header">
          <h3>CPU 使用率</h3>
          <div class="status-indicator" :class="getStatusClass(systemStatus?.cpu?.status)"></div>
        </div>
        <div class="card-content">
          <div class="metric-value">{{ systemStatus?.cpu?.usage || 0 }}%</div>
          <div class="metric-detail">{{ systemStatus?.cpu?.cores || 0 }} 核心</div>
          <div class="progress-bar">
            <div class="progress" :style="{width: `${systemStatus?.cpu?.usage || 0}%`}"></div>
          </div>
        </div>
      </div>

      <!-- 内存监控卡片 -->
      <div class="monitor-card" :class="getStatusClass(systemStatus?.memory?.status)">
        <div class="card-header">
          <h3>内存使用率</h3>
          <div class="status-indicator" :class="getStatusClass(systemStatus?.memory?.status)"></div>
        </div>
        <div class="card-content">
          <div class="metric-value">{{ systemStatus?.memory?.usage?.toFixed(2) || 0 }}%</div>
          <div class="metric-detail">
            已用: {{ systemStatus?.memory?.used || '0 B' }} / 总计: {{ systemStatus?.memory?.total || '0 B' }}
          </div>
          <div class="progress-bar">
            <div class="progress" :style="{width: `${systemStatus?.memory?.usage || 0}%`}"></div>
          </div>
        </div>
      </div>

      <!-- 磁盘监控卡片 -->
      <div class="monitor-card" :class="getStatusClass(systemStatus?.disk?.status)">
        <div class="card-header">
          <h3>磁盘使用率</h3>
          <div class="status-indicator" :class="getStatusClass(systemStatus?.disk?.status)"></div>
        </div>
        <div class="card-content">
          <div class="metric-value">{{ systemStatus?.disk?.usage?.toFixed(2) || 0 }}%</div>
          <div class="metric-detail">
            已用: {{ systemStatus?.disk?.used || '0 B' }} / 总计: {{ systemStatus?.disk?.total || '0 B' }}
          </div>
          <div class="progress-bar">
            <div class="progress" :style="{width: `${systemStatus?.disk?.usage || 0}%`}"></div>
          </div>
        </div>
      </div>

      <!-- 网络监控卡片 -->
      <div class="monitor-card">
        <div class="card-header">
          <h3>网络流量</h3>
        </div>
        <div class="card-content">
          <div class="network-metrics">
            <div class="network-metric">
              <span class="network-label">入站:</span>
              <span class="network-value">{{ systemStatus?.network?.in || '0 B/s' }}</span>
            </div>
            <div class="network-metric">
              <span class="network-label">出站:</span>
              <span class="network-value">{{ systemStatus?.network?.out || '0 B/s' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 系统运行时间卡片 -->
      <div class="monitor-card">
        <div class="card-header">
          <h3>系统运行时间</h3>
        </div>
        <div class="card-content">
          <div class="metric-value uptime">{{ systemStatus?.uptime || '0分钟' }}</div>
        </div>
      </div>
    </div>

    <div class="monitor-controls">
      <button @click="toggleAutoRefresh" :class="{active: isAutoRefreshEnabled}">
        {{ isAutoRefreshEnabled ? '停止自动刷新' : '开始自动刷新' }}
      </button>
      <button @click="refreshNow" :disabled="isRefreshing">
        {{ isRefreshing ? '刷新中...' : '立即刷新' }}
      </button>
      <div class="last-update-time" v-if="lastUpdateTime">
        上次更新: {{ lastUpdateTime }}
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { monitorService, MonitorServiceStatus } from '../../services/monitorService'

export default {
  name: 'SystemMonitorView',
  setup() {
    const systemStatus = ref(null)
    const error = ref(null)
    const isRefreshing = ref(false)
    const isAutoRefreshEnabled = ref(true)
    const lastUpdateTime = ref(null)
    const serviceStatus = ref(MonitorServiceStatus.ACTIVE)
    const retryInfo = ref(null)
    const showDataWhileError = ref(true) // 即使有错误也显示最后的数据

    // 计算服务状态文本
    const serviceStatusText = computed(() => {
      switch (serviceStatus.value) {
        case MonitorServiceStatus.ACTIVE:
          return '正常'
        case MonitorServiceStatus.PAUSED:
          return '已暂停'
        case MonitorServiceStatus.ERROR:
          return '错误'
        case MonitorServiceStatus.RECONNECTING:
          return '重新连接中'
        default:
          return '未知'
      }
    })

    // 计算服务状态样式类
    const serviceStatusClass = computed(() => {
      switch (serviceStatus.value) {
        case MonitorServiceStatus.ACTIVE:
          return 'status-normal'
        case MonitorServiceStatus.PAUSED:
          return 'status-paused'
        case MonitorServiceStatus.ERROR:
          return 'status-error'
        case MonitorServiceStatus.RECONNECTING:
          return 'status-reconnecting'
        default:
          return ''
      }
    })

    // 是否可以重试连接
    const canRetry = computed(() => {
      return serviceStatus.value === MonitorServiceStatus.ERROR || 
             serviceStatus.value === MonitorServiceStatus.PAUSED
    })

    // 获取状态样式类
    const getStatusClass = (status) => {
      if (!status) return ''
      switch (status) {
        case 'normal': return 'status-normal'
        case 'warning': return 'status-warning'
        case 'danger': return 'status-danger'
        default: return ''
      }
    }

    // 更新系统状态
    const updateSystemStatus = async () => {
      if (isRefreshing.value) return
      isRefreshing.value = true
      error.value = null

      try {
        const result = await monitorService.getSystemStatus()
        systemStatus.value = result
        lastUpdateTime.value = new Date().toLocaleTimeString()
        error.value = null
      } catch (err) {
        error.value = err.message || '获取系统状态失败'
        console.error('获取系统状态失败:', err)
      } finally {
        isRefreshing.value = false
      }
    }

    // 处理自动刷新回调
    const handleAutoRefreshCallback = (result) => {
      if (result.success) {
        systemStatus.value = result.data
        lastUpdateTime.value = new Date().toLocaleTimeString()
        error.value = null
        retryInfo.value = null
        serviceStatus.value = result.serviceStatus || MonitorServiceStatus.ACTIVE
      } else {
        error.value = result.error || '获取系统状态失败'
        if (result.retryIn) {
          retryInfo.value = `${result.retryIn}秒后重试...`
        } else {
          retryInfo.value = null
        }
        serviceStatus.value = result.serviceStatus || MonitorServiceStatus.ERROR
      }
    }

    // 切换自动刷新
    const toggleAutoRefresh = () => {
      isAutoRefreshEnabled.value = !isAutoRefreshEnabled.value
      if (isAutoRefreshEnabled.value) {
        startAutoRefresh()
      } else {
        stopAutoRefresh()
      }
    }

    // 开始自动刷新
    const startAutoRefresh = () => {
      monitorService.startAutoRefresh(handleAutoRefreshCallback)
    }

    // 停止自动刷新
    const stopAutoRefresh = () => {
      monitorService.stopAutoRefresh()
      serviceStatus.value = MonitorServiceStatus.PAUSED
    }

    // 立即刷新
    const refreshNow = () => {
      updateSystemStatus()
    }

    // 重试连接
    const retryConnection = () => {
      if (isAutoRefreshEnabled.value) {
        startAutoRefresh()
      } else {
        refreshNow()
      }
    }

    // 组件挂载时
    onMounted(() => {
      // 初始化监控服务
      monitorService.init()
      
      // 如果启用了自动刷新，开始自动刷新
      if (isAutoRefreshEnabled.value) {
        startAutoRefresh()
      } else {
        // 否则立即获取一次数据
        refreshNow()
      }
    })

    // 组件卸载时
    onUnmounted(() => {
      // 清理监控服务资源
      monitorService.cleanup()
    })

    return {
      systemStatus,
      error,
      isRefreshing,
      isAutoRefreshEnabled,
      lastUpdateTime,
      serviceStatus,
      serviceStatusText,
      serviceStatusClass,
      retryInfo,
      canRetry,
      showDataWhileError,
      getStatusClass,
      toggleAutoRefresh,
      refreshNow,
      retryConnection
    }
  }
}
</script>

<style scoped>
.system-monitor-container {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.monitor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e0e6ed;
}

.monitor-header h1 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.service-status {
  display: flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 4px;
  font-weight: 500;
}

.status-normal {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-warning {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.status-danger {
  background-color: #fef0f0;
  color: #f56c6c;
}

.status-paused {
  background-color: #f4f4f5;
  color: #909399;
}

.status-error {
  background-color: #fef0f0;
  color: #f56c6c;
}

.status-reconnecting {
  background-color: #ecf5ff;
  color: #409eff;
}

.retry-info {
  margin-left: 8px;
  font-size: 12px;
  opacity: 0.8;
}

.error-message {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 20px;
  background-color: #fef0f0;
  border-radius: 4px;
  color: #f56c6c;
}

.error-icon {
  margin-right: 8px;
  font-size: 18px;
}

.retry-button {
  margin-left: auto;
  padding: 6px 12px;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.retry-button:hover {
  background-color: #e64242;
}

.monitor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.monitor-card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  overflow: hidden;
  transition: transform 0.3s, box-shadow 0.3s;
}

.monitor-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.card-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.status-indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.card-content {
  padding: 16px;
}

.metric-value {
  font-size: 28px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #303133;
}

.metric-detail {
  font-size: 14px;
  color: #606266;
  margin-bottom: 12px;
}

.progress-bar {
  height: 8px;
  background-color: #e6e6e6;
  border-radius: 4px;
  overflow: hidden;
}

.progress {
  height: 100%;
  background-color: #67c23a;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.status-normal .progress {
  background-color: #67c23a;
}

.status-warning .progress {
  background-color: #e6a23c;
}

.status-danger .progress {
  background-color: #f56c6c;
}

.network-metrics {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.network-metric {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.network-label {
  color: #606266;
}

.network-value {
  font-weight: 500;
  color: #303133;
}

.uptime {
  text-align: center;
}

.monitor-controls {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e0e6ed;
}

.monitor-controls button {
  padding: 8px 16px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.monitor-controls button:hover {
  background-color: #66b1ff;
}

.monitor-controls button:disabled {
  background-color: #a0cfff;
  cursor: not-allowed;
}

.monitor-controls button.active {
  background-color: #67c23a;
}

.last-update-time {
  margin-left: auto;
  font-size: 14px;
  color: #909399;
}

@media (max-width: 768px) {
  .monitor-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .service-status {
    margin-top: 10px;
  }
  
  .monitor-grid {
    grid-template-columns: 1fr;
  }
  
  .monitor-controls {
    flex-wrap: wrap;
  }
  
  .last-update-time {
    margin-left: 0;
    margin-top: 10px;
    width: 100%;
  }
}
</style>