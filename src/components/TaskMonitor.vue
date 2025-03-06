<template>
  <div class="task-monitor-container">
    <div class="connection-status">
      <h3>WebSocket连接状态</h3>
      <div class="status-badge" :class="connectionStatusClass">
        {{ connectionStatusText }}
      </div>
      <button v-if="canReconnect" @click="reconnect" class="reconnect-button">
        重新连接
      </button>
    </div>

    <div class="task-list" v-if="hasTasks">
      <h3>活跃任务</h3>
      <div v-for="task in tasks" :key="task.id" class="task-item">
        <div class="task-header">
          <span class="task-id">任务ID: {{ task.id }}</span>
          <span class="task-status" :class="getStatusClass(task.status)">
            {{ getStatusText(task.status) }}
          </span>
        </div>
        <div class="task-progress">
          <div class="progress-bar">
            <div class="progress" :style="{width: `${task.progress}%`}"></div>
          </div>
          <span class="progress-text">{{ task.progress }}%</span>
        </div>
        <div class="task-message" v-if="task.message">
          {{ task.message }}
        </div>
        <div class="task-actions" v-if="task.status === 'processing' || task.status === 'pending'">
          <button @click="cancelTask(task.id)" class="cancel-button">
            取消任务
          </button>
        </div>
      </div>
    </div>

    <div class="no-tasks" v-else>
      <p>当前没有活跃任务</p>
      <button @click="createMockTask" class="create-task-button">
        创建模拟任务
      </button>
    </div>

    <div class="connection-info">
      <h3>连接信息</h3>
      <div class="info-item">
        <span class="info-label">重连次数:</span>
        <span class="info-value">{{ reconnectAttempts }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">上次心跳:</span>
        <span class="info-value">{{ lastHeartbeatTime || '无' }}</span>
      </div>
      <div class="info-item">
        <span class="info-label">客户端ID:</span>
        <span class="info-value">{{ clientId || '未生成' }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { taskService, TaskStatus } from '../services/taskService'

export default {
  name: 'TaskMonitor',
  setup() {
    const connectionStatus = ref('disconnected')
    const tasks = ref([])
    const reconnectAttempts = ref(0)
    const lastHeartbeatTime = ref(null)
    const clientId = ref(null)

    // 计算连接状态文本
    const connectionStatusText = computed(() => {
      switch (connectionStatus.value) {
        case 'connected':
          return '已连接'
        case 'connecting':
          return '连接中'
        case 'reconnecting':
          return '重新连接中'
        case 'disconnected':
          return '已断开'
        default:
          return '未知'
      }
    })

    // 计算连接状态样式类
    const connectionStatusClass = computed(() => {
      switch (connectionStatus.value) {
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
    })

    // 是否可以重新连接
    const canReconnect = computed(() => {
      return connectionStatus.value === 'disconnected'
    })

    // 是否有任务
    const hasTasks = computed(() => {
      return tasks.value.length > 0
    })

    // 获取任务状态样式类
    const getStatusClass = (status) => {
      switch (status) {
        case TaskStatus.COMPLETED:
          return 'status-success'
        case TaskStatus.PROCESSING:
          return 'status-processing'
        case TaskStatus.PENDING:
          return 'status-pending'
        case TaskStatus.FAILED:
          return 'status-error'
        case TaskStatus.CANCELLED:
          return 'status-cancelled'
        default:
          return ''
      }
    }

    // 获取任务状态文本
    const getStatusText = (status) => {
      switch (status) {
        case TaskStatus.COMPLETED:
          return '已完成'
        case TaskStatus.PROCESSING:
          return '处理中'
        case TaskStatus.PENDING:
          return '等待中'
        case TaskStatus.FAILED:
          return '失败'
        case TaskStatus.CANCELLED:
          return '已取消'
        default:
          return '未知'
      }
    }

    // 重新连接
    const reconnect = () => {
      taskService.initWebSocket()
    }

    // 取消任务
    const cancelTask = async (taskId) => {
      const result = await taskService.cancelTask(taskId)
      if (result.success) {
        // 更新任务状态
        const taskIndex = tasks.value.findIndex(t => t.id === taskId)
        if (taskIndex !== -1) {
          tasks.value[taskIndex].status = TaskStatus.CANCELLED
          tasks.value[taskIndex].message = '任务已取消'
        }
      }
    }

    // 创建模拟任务
    const createMockTask = () => {
      const taskId = 'mock_' + Date.now()
      const newTask = {
        id: taskId,
        status: TaskStatus.PENDING,
        progress: 0,
        message: '任务初始化中...',
        createdAt: new Date().toISOString()
      }
      
      tasks.value.push(newTask)
      
      // 模拟任务进度更新
      let progress = 0
      const progressInterval = setInterval(() => {
        progress += Math.floor(Math.random() * 10) + 1
        if (progress >= 100) {
          progress = 100
          clearInterval(progressInterval)
          
          // 更新任务状态为完成
          const taskIndex = tasks.value.findIndex(t => t.id === taskId)
          if (taskIndex !== -1) {
            tasks.value[taskIndex].status = TaskStatus.COMPLETED
            tasks.value[taskIndex].message = '任务已完成'
          }
          
          // 模拟任务完成后的清理
          setTimeout(() => {
            const index = tasks.value.findIndex(t => t.id === taskId)
            if (index !== -1) {
              tasks.value.splice(index, 1)
            }
          }, 5000)
        }
        
        // 更新任务进度
        const taskIndex = tasks.value.findIndex(t => t.id === taskId)
        if (taskIndex !== -1) {
          if (tasks.value[taskIndex].status === TaskStatus.PENDING && progress > 5) {
            tasks.value[taskIndex].status = TaskStatus.PROCESSING
            tasks.value[taskIndex].message = '任务处理中...'
          }
          tasks.value[taskIndex].progress = progress
        }
      }, 1000)
      
      // 存储定时器ID以便清理
      mockTaskTimers.set(taskId, progressInterval)
    }
    
    // 存储模拟任务定时器
    const mockTaskTimers = new Map()
    
    // 监听WebSocket连接状态变化
    const updateConnectionStatus = () => {
      connectionStatus.value = taskService.connectionState
      reconnectAttempts.value = taskService.reconnectAttempts
      clientId.value = taskService.getClientId()
    }
    
    // 定期更新连接状态
    let statusUpdateTimer = null
    
    // 组件挂载时
    onMounted(() => {
      // 初始化WebSocket连接
      taskService.initWebSocket()
      
      // 定期更新连接状态
      statusUpdateTimer = setInterval(() => {
        updateConnectionStatus()
      }, 1000)
      
      // 初始更新一次
      updateConnectionStatus()
    })
    
    // 组件卸载时
    onUnmounted(() => {
      // 清理定时器
      if (statusUpdateTimer) {
        clearInterval(statusUpdateTimer)
        statusUpdateTimer = null
      }
      
      // 清理模拟任务定时器
      mockTaskTimers.forEach((timerId) => {
        clearInterval(timerId)
      })
      mockTaskTimers.clear()
      
      // 注意：不要在这里关闭WebSocket连接，因为其他组件可能仍在使用
      // 但应该移除所有事件监听和任务订阅
      tasks.value.forEach(task => {
        if (task.id) {
          taskService.offTaskStatusChange(task.id)
          taskService.offTaskProgressChange(task.id)
        }
      })
    })

    return {
      connectionStatus,
      connectionStatusText,
      connectionStatusClass,
      tasks,
      reconnectAttempts,
      lastHeartbeatTime,
      clientId,
      canReconnect,
      hasTasks,
      getStatusClass,
      getStatusText,
      reconnect,
      cancelTask,
      createMockTask
    }
  }
}
</script>

<style scoped>
.task-monitor-container {
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.connection-status {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #e0e6ed;
}

.connection-status h3 {
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

.status-processing {
  background-color: #ecf5ff;
  color: #409eff;
}

.status-pending {
  background-color: #f4f4f5;
  color: #909399;
}

.status-cancelled {
  background-color: #f4f4f5;
  color: #909399;
}

.reconnect-button {
  margin-left: auto;
  padding: 6px 12px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.reconnect-button:hover {
  background-color: #66b1ff;
}

.task-list {
  margin-bottom: 20px;
}

.task-list h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.task-item {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 15px;
  margin-bottom: 15px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.task-id {
  font-size: 14px;
  color: #606266;
}

.task-status {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 500;
}

.task-progress {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.progress-bar {
  flex-grow: 1;
  height: 8px;
  background-color: #e6e6e6;
  border-radius: 4px;
  overflow: hidden;
  margin-right: 10px;
}

.progress {
  height: 100%;
  background-color: #409eff;
  border-radius: 4px;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 14px;
  color: #606266;
  min-width: 40px;
  text-align: right;
}

.task-message {
  font-size: 14px;
  color: #606266;
  margin-bottom: 10px;
}

.task-actions {
  display: flex;
  justify-content: flex-end;
}

.cancel-button {
  padding: 4px 10px;
  background-color: #f56c6c;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
  font-size: 12px;
}

.cancel-button:hover {
  background-color: #e64242;
}

.no-tasks {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 20px;
  text-align: center;
  margin-bottom: 20px;
}

.no-tasks p {
  margin-bottom: 15px;
  color: #909399;
}

.create-task-button {
  padding: 6px 12px;
  background-color: #409eff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.create-task-button:hover {
  background-color: #66b1ff;
}

.connection-info {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  padding: 15px;
}

.connection-info h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
  font-size: 14px;
}

.info-label {
  color: #606266;
}

.info-value {
  font-weight: 500;
  color: #303133;
}

@media (max-width: 768px) {
  .connection-status {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .status-badge {
    margin: 10px 0;
  }
  
  .reconnect-button {
    margin-left: 0;
  }
}
</style>