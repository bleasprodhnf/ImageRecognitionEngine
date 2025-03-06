<template>
  <div class="recovery-monitor-container">
    <div class="recovery-status">
      <h3>系统恢复状态</h3>
      <div class="status-badge" :class="recoveryStatusClass">
        {{ recoveryStatusText }}
      </div>
      <button v-if="canTriggerRecovery" @click="triggerRecovery" class="recovery-button">
        手动触发恢复
      </button>
    </div>

    <div class="recovery-info">
      <h3>恢复信息</h3>
      <div class="info-item">
        <span class="info-label">恢复尝试次数:</span>
        <span class="info-value">{{ recoveryStatus.recoveryAttempts }}/{{ recoveryStatus.maxRecoveryAttempts }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">上次恢复时间:</span>
        <span class="info-value">{{ formatTime(recoveryStatus.lastRecoveryTime) }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">监控服务状态:</span>
        <span class="info-value" :class="getMonitorStatusClass(recoveryStatus.monitorStatus)">
          {{ getMonitorStatusText(recoveryStatus.monitorStatus) }}
        </span>
      </div>
      <div class="info-item">
        <span class="info-label">任务服务状态:</span>
        <span class="info-value" :class="getConnectionStatusClass(recoveryStatus.taskConnectionState)">
          {{ getConnectionStatusText(recoveryStatus.taskConnectionState) }}
        </span>
      </div>
      <div class="info-item">
        <span class="info-label">网络状态:</span>
        <span class="info-value" :class="getNetworkStatusClass(recoveryStatus.networkStatus)">
          {{ recoveryStatus.networkStatus === 'online' ? '在线' : '离线' }}
        </span>
      </div>
    </div>

    <div class="recovery-logs">
      <h3>恢复日志</h3>
      <div v-if="logs.length > 0" class="logs-container">
        <div v-for="(log, index) in logs" :key="index" class="log-item" :class="getLogClass(log)">
          <div class="log-time">{{ formatTime(log.timestamp) }}</div>
          <div class="log-message">{{ log.message }}</div>
          <div v-if="log.error" class="log-error">错误: {{ log.error }}</div>
        </div>
      </div>
      <div v-else class="no-logs">
        <p>暂无恢复日志</p>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { recoveryService, RecoveryStatus } from '../services/recoveryService'
import { MonitorServiceStatus } from '../services/monitorService'

export default {
  name: 'SystemRecovery',
  setup() {
    const recoveryStatus = ref({
      status: RecoveryStatus.IDLE,
      recoveryAttempts: 0,
      maxRecoveryAttempts: 5,
      lastRecoveryTime: null,
      monitorStatus: MonitorServiceStatus.ACTIVE,
      taskConnectionState: 'disconnected',
      networkStatus: 'offline'
    })
    
    const logs = ref([])
    const maxLogs = 50 // 最多显示的日志数量
    
    // 计算恢复状态文本
    const recoveryStatusText = computed(() => {
      switch (recoveryStatus.value.status) {
        case RecoveryStatus.IDLE:
          return '空闲'
        case RecoveryStatus.MONITORING:
          return '监控中'
        case RecoveryStatus.RECOVERING:
          return '恢复中'
        case RecoveryStatus.FAILED:
          return '恢复失败'
        default:
          return '未知'
      }
    })
    
    // 计算恢复状态样式类
    const recoveryStatusClass = computed(() => {
      switch (recoveryStatus.value.status) {
        case RecoveryStatus.IDLE:
          return 'status-neutral'
        case RecoveryStatus.MONITORING:
          return 'status-success'
        case RecoveryStatus.RECOVERING:
          return 'status-warning'
        case RecoveryStatus.FAILED:
          return 'status-error'
        default:
          return ''
      }
    })
    
    // 是否可以手动触发恢复
    const canTriggerRecovery = computed(() => {
      return recoveryStatus.value.status !== RecoveryStatus.RECOVERING
    })
    
    // 获取监控服务状态样式类
    const getMonitorStatusClass = (status) => {
      switch (status) {
        case MonitorServiceStatus.ACTIVE:
          return 'status-success'
        case MonitorServiceStatus.RECONNECTING:
          return 'status-warning'
        case MonitorServiceStatus.PAUSED:
        case MonitorServiceStatus.ERROR:
          return 'status-error'
        default:
          return ''
      }
    }
    
    // 获取监控服务状态文本
    const getMonitorStatusText = (status) => {
      switch (status) {
        case MonitorServiceStatus.ACTIVE:
          return '活跃'
        case MonitorServiceStatus.RECONNECTING:
          return '重连中'
        case MonitorServiceStatus.PAUSED:
          return '已暂停'
        case MonitorServiceStatus.ERROR:
          return '错误'
        default:
          return '未知'
      }
    }
    
    // 获取连接状态样式类
    const getConnectionStatusClass = (status) => {
      switch (status) {
        case 'connected':
          return 'status-success'
        case 'connecting':
        case 'reconnecting':
          return 'status-warning'
        case 'disconnected':
          return 'status-error'
        default:
          return ''
      }
    }
    
    // 获取连接状态文本
    const getConnectionStatusText = (status) => {
      switch (status) {
        case 'connected':
          return '已连接'
        case 'connecting':
          return '连接中'
        case 'reconnecting':
          return '重连中'
        case 'disconnected':
          return '已断开'
        default:
          return '未知'
      }
    }
    
    // 获取网络状态样式类
    const getNetworkStatusClass = (status) => {
      return status === 'online' ? 'status-success' : 'status-error'
    }
    
    // 获取日志样式类
    const getLogClass = (log) => {
      if (log.status === RecoveryStatus.FAILED || log.error) {
        return 'log-error'
      } else if (log.status === RecoveryStatus.RECOVERING) {
        return 'log-warning'
      } else if (log.status === RecoveryStatus.MONITORING) {
        return 'log-success'
      }
      return ''
    }
    
    // 格式化时间
    const formatTime = (timestamp) => {
      if (!timestamp) return '无'
      
      const date = new Date(timestamp)
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    }
    
    // 手动触发恢复
    const triggerRecovery = async () => {
      try {
        await recoveryService.triggerRecovery()
        addLog({
          status: RecoveryStatus.RECOVERING,
          message: '手动触发系统恢复',
          timestamp: new Date().toISOString()
        })
      } catch (error) {
        addLog({
          status: RecoveryStatus.FAILED,
          message: '手动触发系统恢复失败',
          timestamp: new Date().toISOString(),
          error: error.message
        })
      }
    }
    
    // 添加日志
    const addLog = (log) => {
      logs.value.unshift(log)
      
      // 限制日志数量
      if (logs.value.length > maxLogs) {
        logs.value = logs.value.slice(0, maxLogs)
      }
    }
    
    // 更新恢复状态
    const updateRecoveryStatus = () => {
      recoveryStatus.value = recoveryService.getStatus()
    }
    
    // 处理状态变化
    const handleStatusChange = (statusData) => {
      updateRecoveryStatus()
      addLog(statusData)
    }
    
    // 定期更新状态
    let statusUpdateTimer = null
    
    // 组件挂载时
    onMounted(() => {
      // 初始化恢复服务
      recoveryService.init()
      
      // 注册状态变化回调
      recoveryService.onStatusChange(handleStatusChange)
      
      // 定期更新状态
      statusUpdateTimer = setInterval(() => {
        updateRecoveryStatus()
      }, 5000)
      
      // 初始更新一次
      updateRecoveryStatus()
      
      // 添加初始日志
      addLog({
        status: RecoveryStatus.MONITORING,
        message: '系统恢复监控已启动',
        timestamp: new Date().toISOString()
      })
    })
    
    // 组件卸载时
    onUnmounted(() => {
      // 取消注册状态变化回调
      recoveryService.offStatusChange(handleStatusChange)
      
      // 清理定时器
      if (statusUpdateTimer) {
        clearInterval(statusUpdateTimer)
        statusUpdateTimer = null
      }
      
      // 注意：不要在这里清理恢复服务，因为其他组件可能仍在使用
    })
    
    return {
      recoveryStatus,
      logs,
      recoveryStatusText,
      recoveryStatusClass,
      canTriggerRecovery,
      getMonitorStatusClass,
      getMonitorStatusText,
      getConnectionStatusClass,
      getConnectionStatusText,
      getNetworkStatusClass,
      getLogClass,
      formatTime,
      triggerRecovery
    }
  }
}
</script>

<style scoped>
.recovery-monitor-container {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.recovery-status {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e0e6ed;
}

.recovery-status h3 {
  margin: 0;
  margin-right: 15px;
  font-size: 18px;
  color: #303133;
}

.status-badge {
  padding: 4px 10px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: 500;
}

.status-success {
  background-color: #f0f9eb;
  color: #67c23a;
}

.status-warning {
  background-color: #fdf6ec;
  color: #e6a23c;
}

.status-error {
  background-color: #fef0f0;
  color: #f56c6c;
}

.status-neutral {
  background-color: #f4f4f5;
  color: #909399;
}

.recovery-button {
  margin-left: auto;
  padding: 6px 12px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.recovery-button:hover {
  background-color: #66b1ff;
}

.recovery-info {
  margin-bottom: 20px;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 15px;
}

.recovery-info h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.info-item {
  display: flex;
  margin-bottom: 10px;
}

.info-label {
  width: 120px;
  font-weight: 500;
  color: #606266;
}

.info-value {
  flex-grow: 1;
  color: #303133;
}

.recovery-logs {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 15px;
}

.recovery-logs h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.logs-container {
  max-height: 300px;
  overflow-y: auto;
}

.log-item {
  padding: 10px;
  border-bottom: 1px solid #ebeef5;
  font-size: 14px;
}

.log-item:last-child {
  border-bottom: none;
}

.log-time {
  font-size: 12px;
  color: #909399;
  margin-bottom: 5px;
}

.log-message {
  color: #606266;
}

.log-error {
  color: #f56c6c;
  margin-top: 5px;
}

.log-error {
  background-color: #fef0f0;
}

.log-warning {
  background-color: #fdf6ec;
}

.log-success {
  background-color: #f0f9eb;
}

.no-logs {
  text-align: center;
  padding: 20px;
  color: #909399;
}
</style>