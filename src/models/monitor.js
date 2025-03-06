// 系统监控数据模型
import { reactive } from 'vue'
import { monitorService } from '../services/monitorService'

// 初始状态
const initialState = {
  loading: false,
  error: null,
  data: {
    cpu: {
      usage: 0,
      cores: 0,
      status: 'normal'
    },
    memory: {
      total: '0 GB',
      used: '0 GB',
      free: '0 GB',
      usage: 0,
      status: 'normal'
    },
    disk: {
      total: '0 GB',
      used: '0 GB',
      free: '0 GB',
      usage: 0,
      status: 'normal'
    },
    network: {
      in: '0 B/s',
      out: '0 B/s'
    },
    uptime: '0分钟'
  },
  refreshInterval: null
}

// 创建响应式状态
export const monitorState = reactive({ ...initialState })

// 监控数据模型
export const monitorModel = {
  // 开始监控
  startMonitoring(interval = 30000) {
    this.fetchSystemStatus()
    monitorState.refreshInterval = setInterval(() => {
      this.fetchSystemStatus()
    }, interval)
  },

  // 停止监控
  stopMonitoring() {
    if (monitorState.refreshInterval) {
      clearInterval(monitorState.refreshInterval)
      monitorState.refreshInterval = null
    }
  },

  // 获取系统状态
  async fetchSystemStatus() {
    try {
      monitorState.loading = true
      monitorState.error = null
      const data = await monitorService.getSystemStatus()
      monitorState.data = data
    } catch (error) {
      monitorState.error = error.message || '获取系统状态失败'
      console.error('获取系统状态失败:', error)
    } finally {
      monitorState.loading = false
    }
  },

  // 重置状态
  resetState() {
    Object.assign(monitorState, initialState)
  }
}