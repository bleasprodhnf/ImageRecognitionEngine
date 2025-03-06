// 系统监控服务
import request from '../api/request'

// 定义自动刷新间隔（毫秒）
const REFRESH_INTERVAL = 30000
// 定义最大重试间隔（毫秒）
const MAX_RETRY_INTERVAL = 60000
// 定义网络状态检测超时（毫秒）
const NETWORK_TIMEOUT = 8000
// 定义恢复正常刷新间隔的时间（毫秒）
const RESTORE_INTERVAL_TIMEOUT = 300000 // 5分钟
// 定义网络检查间隔（毫秒）
const NETWORK_CHECK_INTERVAL = 60000 // 1分钟

// 系统监控指标类型
export const MonitorMetrics = {
  CPU: 'cpu',
  MEMORY: 'memory',
  DISK: 'disk',
  NETWORK: 'network',
  UPTIME: 'uptime'
}

// 监控服务状态
export const MonitorServiceStatus = {
  ACTIVE: 'active',
  PAUSED: 'paused',
  ERROR: 'error',
  RECONNECTING: 'reconnecting'
}

// 系统监控服务
export const monitorService = {
  // 存储定时器ID
  refreshTimer: null,
  
  // 存储恢复间隔定时器ID
  restoreIntervalTimer: null,
  
  // 存储网络检查定时器ID
  networkCheckTimer: null,

  // 存储服务状态
  serviceStatus: MonitorServiceStatus.ACTIVE,
  
  // 存储网络状态
  networkStatus: {
    isOnline: true,
    lastCheckTime: null,
    checkInProgress: false
  },
  
  // 检查网络连接状态
  async checkNetworkStatus() {
    if (this.networkStatus.checkInProgress) return this.networkStatus.isOnline
    
    this.networkStatus.checkInProgress = true
    
    try {
      // 使用navigator.onLine作为初步检查
      const isOnline = navigator.onLine
      
      if (!isOnline) {
        this.networkStatus.isOnline = false
        return false
      }
      
      // 发送一个小型请求来验证实际连接
      await Promise.race([
        fetch('/api/health-check', { method: 'HEAD', cache: 'no-store' }),
        new Promise((_, reject) => setTimeout(() => reject(new Error('网络检测超时')), NETWORK_TIMEOUT))
      ])
      
      this.networkStatus.isOnline = true
      this.networkStatus.lastCheckTime = Date.now()
      return true
    } catch (error) {
      console.warn('网络连接检测失败:', error.message)
      this.networkStatus.isOnline = false
      return false
    } finally {
      this.networkStatus.checkInProgress = false
    }
  },
  
  // 开始自动刷新
  async startAutoRefresh(callback, maxRetries = 3, retryDelay = 5000) {
    let retryCount = 0
    let isRefreshing = false
    let consecutiveFailures = 0
    const maxConsecutiveFailures = 5
    
    // 保存回调函数，用于网络恢复时重新启动
    this.lastCallback = callback
    
    // 更新服务状态
    this.serviceStatus = MonitorServiceStatus.ACTIVE
    
    const fetchWithRetry = async () => {
      if (isRefreshing) return
      isRefreshing = true
      
      // 先检查网络状态
      const isNetworkAvailable = await this.checkNetworkStatus()
      if (!isNetworkAvailable) {
        console.warn('网络连接不可用，暂停监控刷新')
        this.serviceStatus = MonitorServiceStatus.PAUSED
        callback({ 
          success: false, 
          error: '网络连接不可用，监控数据刷新已暂停',
          networkStatus: 'offline'
        })
        isRefreshing = false
        return
      }
      
      try {
        // 增加超时处理，避免长时间等待
        const data = await Promise.race([
          this.getSystemStatus(),
          new Promise((_, reject) => 
            setTimeout(() => reject(new Error('请求超时，请检查网络连接')), NETWORK_TIMEOUT)
          )
        ])
        retryCount = 0 // 重置重试计数
        consecutiveFailures = 0 // 重置连续失败计数
        this.serviceStatus = MonitorServiceStatus.ACTIVE
        callback({ success: true, data, serviceStatus: this.serviceStatus })
      } catch (error) {
        console.error('自动刷新失败:', error.message)
        consecutiveFailures++
        
        // 智能重试策略
        if (retryCount < maxRetries) {
          retryCount++
          this.serviceStatus = MonitorServiceStatus.RECONNECTING
          
          // 使用指数退避策略计算重试延迟，但设置上限
          const backoffDelay = Math.min(retryDelay * Math.pow(1.5, retryCount - 1), MAX_RETRY_INTERVAL)
          // 添加随机因子避免多客户端同时重试
          const jitter = 0.5 * Math.random() + 0.5 // 0.5-1.0之间的随机数
          const finalDelay = Math.round(backoffDelay * jitter)
          
          console.log(`${finalDelay/1000}秒后进行第${retryCount}/${maxRetries}次重试...`)
          setTimeout(fetchWithRetry, finalDelay)
          callback({ 
            success: false, 
            error: `正在重试 (${retryCount}/${maxRetries})`,
            retryIn: Math.round(finalDelay/1000),
            serviceStatus: this.serviceStatus
          })
        } else {
          console.error('达到最大重试次数，停止当前重试')
          this.serviceStatus = MonitorServiceStatus.ERROR
          callback({ 
            success: false, 
            error: '系统监控数据暂时不可用，请稍后再试',
            serviceStatus: this.serviceStatus 
          })
          // 重置重试次数，允许下次刷新重新尝试
          retryCount = 0
          
          // 如果连续失败次数过多，增加刷新间隔以减轻服务器负担
          if (consecutiveFailures >= maxConsecutiveFailures && this.refreshTimer) {
            console.warn(`检测到${consecutiveFailures}次连续失败，临时增加刷新间隔`)
            clearInterval(this.refreshTimer)
            this.refreshTimer = setInterval(fetchWithRetry, REFRESH_INTERVAL * 2)
            // 设置恢复定时器，一段时间后恢复正常刷新间隔
            if (this.restoreIntervalTimer) {
              clearTimeout(this.restoreIntervalTimer)
            }
            this.restoreIntervalTimer = setTimeout(() => {
              if (this.refreshTimer) {
                clearInterval(this.refreshTimer)
                this.refreshTimer = setInterval(fetchWithRetry, REFRESH_INTERVAL)
                console.log('已恢复正常刷新间隔')
              }
              this.restoreIntervalTimer = null
            }, RESTORE_INTERVAL_TIMEOUT)
          }
        }
      } finally {
        isRefreshing = false
      }
    }

    // 立即执行一次
    await fetchWithRetry()
    
    // 设置定时刷新
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer) // 确保不会创建多个定时器
    }
    this.refreshTimer = setInterval(fetchWithRetry, REFRESH_INTERVAL)
    
    // 添加网络状态变化监听
    window.addEventListener('online', this.handleNetworkChange)
    window.addEventListener('offline', this.handleNetworkChange)
    
    return this.refreshTimer // 返回定时器ID，方便外部管理
  },
  
  // 处理网络状态变化
  handleNetworkChange = async (event) => {
    console.log(`网络状态变化: ${event.type}`)
    
    if (event.type === 'online') {
      // 网络恢复，重新开始刷新
      this.networkStatus.isOnline = true
      this.networkStatus.lastCheckTime = Date.now()
      
      // 验证网络连接的可用性
      try {
        const isNetworkAvailable = await this.checkNetworkStatus()
        if (!isNetworkAvailable) {
          console.warn('网络连接不稳定，等待下一次检查')
          return
        }
        
        if (this.serviceStatus === MonitorServiceStatus.PAUSED || 
            this.serviceStatus === MonitorServiceStatus.ERROR) {
          console.log('网络已恢复，重新开始监控刷新')
          // 如果定时器不存在且有保存的回调函数，重新创建定时器
          if (!this.refreshTimer && this.lastCallback) {
            await this.startAutoRefresh(this.lastCallback)
          } else if (this.refreshTimer) {
            // 如果定时器存在，只需更新状态
            this.serviceStatus = MonitorServiceStatus.ACTIVE
          }
          
          // 通知UI网络已恢复
          if (this.lastCallback) {
            this.lastCallback({
              success: true,
              message: '网络已恢复，监控服务已重新启动',
              networkStatus: 'online',
              serviceStatus: this.serviceStatus
            })
          }
        }
      } catch (error) {
        console.error('网络恢复处理失败:', error)
        this.networkStatus.isOnline = false
      }
    } else if (event.type === 'offline') {
      // 网络断开，暂停刷新
      this.networkStatus.isOnline = false
      this.serviceStatus = MonitorServiceStatus.PAUSED
      console.warn('网络已断开，监控刷新已暂停')
      
      // 清理定时器
      if (this.refreshTimer) {
        clearInterval(this.refreshTimer)
        this.refreshTimer = null
      }
      
      // 通知UI更新状态
      if (this.lastCallback) {
        try {
          this.lastCallback({
            success: false,
            error: '网络连接已断开，监控数据刷新已暂停',
            networkStatus: 'offline',
            serviceStatus: this.serviceStatus
          })
        } catch (error) {
          console.error('通知UI网络状态变化失败:', error)
        }
      }
    }
  },
  
  // 存储最后一次回调函数
  lastCallback: null,

  // 初始化监控服务
  init() {
    // 清理可能存在的资源，确保不会有资源泄漏
    this.cleanup()
    
    // 重置服务状态
    this.serviceStatus = MonitorServiceStatus.ACTIVE
    
    // 重置网络状态
    this.networkStatus = {
      isOnline: navigator.onLine,
      lastCheckTime: null,
      checkInProgress: false
    }
    
    // 添加网络状态变化监听
    window.addEventListener('online', this.handleNetworkChange)
    window.addEventListener('offline', this.handleNetworkChange)
    
    // 启动定期网络检查
    this.startNetworkCheck()
    
    console.log('监控服务已初始化，当前网络状态:', this.networkStatus.isOnline ? '在线' : '离线')
    return this
  },
  
  // 启动定期网络检查
  startNetworkCheck() {
    // 清理现有定时器
    if (this.networkCheckTimer) {
      clearInterval(this.networkCheckTimer)
    }
    
    // 设置定期检查网络状态
    this.networkCheckTimer = setInterval(async () => {
      // 只在服务处于暂停或错误状态时检查
      if (this.serviceStatus === MonitorServiceStatus.PAUSED || 
          this.serviceStatus === MonitorServiceStatus.ERROR) {
        const isNetworkAvailable = await this.checkNetworkStatus()
        
        // 如果网络恢复且服务处于暂停状态，尝试重新启动服务
        if (isNetworkAvailable && 
            (this.serviceStatus === MonitorServiceStatus.PAUSED || 
             this.serviceStatus === MonitorServiceStatus.ERROR)) {
          console.log('网络已恢复，尝试重新启动监控服务')
          if (this.lastCallback) {
            this.startAutoRefresh(this.lastCallback)
          }
        }
      }
    }, NETWORK_CHECK_INTERVAL)
    
    return this.networkCheckTimer
  },
  
  // 停止自动刷新
  stopAutoRefresh() {
    // 清理刷新定时器
    if (this.refreshTimer) {
      clearInterval(this.refreshTimer)
      this.refreshTimer = null
      console.log('已停止自动刷新')
    }
    
    // 重置服务状态
    this.serviceStatus = MonitorServiceStatus.PAUSED
    
    return true
  },
  
  // 清理所有资源（用于组件卸载时调用）
  cleanup() {
    // 停止自动刷新
    this.stopAutoRefresh()
    
    // 清理其他可能的定时器
    if (this.restoreIntervalTimer) {
      clearTimeout(this.restoreIntervalTimer)
      this.restoreIntervalTimer = null
    }
    
    // 清理网络检查定时器
    if (this.networkCheckTimer) {
      clearInterval(this.networkCheckTimer)
      this.networkCheckTimer = null
    }
    
    // 移除所有事件监听器
    window.removeEventListener('online', this.handleNetworkChange)
    window.removeEventListener('offline', this.handleNetworkChange)
    
    // 重置状态
    this.lastCallback = null
    this.serviceStatus = MonitorServiceStatus.PAUSED
    this.networkStatus.isOnline = navigator.onLine
    this.networkStatus.checkInProgress = false
    
    console.log('监控服务资源已完全清理')
  },

  // 获取性能监控数据
  async getPerformanceMetrics(modelId, startTime, endTime, metrics) {
    try {
      const response = await request({
        url: '/admin/models/performance',
        method: 'get',
        params: {
          modelId,
          startTime,
          endTime,
          metrics
        }
      })
      if (!response.data) {
        throw new Error('获取性能监控数据失败：服务器返回空数据')
      }
      return response.data
    } catch (error) {
      console.error('获取性能监控数据失败:', error)
      if (error.response?.status === 404) {
        throw new Error('未找到指定模型的性能数据')
      } else if (error.response?.status === 403) {
        throw new Error('没有权限访问性能监控数据')
      }
      throw error
    }
  },
  
  async getSystemStatus() {
    try {
      const response = await request({
        url: '/admin/system/monitor',
        method: 'get',
        timeout: NETWORK_TIMEOUT // 添加超时设置
      })
      if (!response.data) {
        throw new Error('获取系统状态失败：服务器返回空数据')
      }
      return this.formatSystemStatus(response.data)
    } catch (error) {
      console.error('获取系统状态失败:', error)
      // 增加更详细的错误处理
      if (error.response?.status === 401) {
        throw new Error('未授权访问系统监控数据，请重新登录')
      } else if (error.response?.status === 403) {
        throw new Error('没有权限访问系统监控数据')
      } else if (error.response?.status === 503) {
        throw new Error('系统监控服务暂时不可用，请稍后重试')
      } else if (error.response?.status === 429) {
        throw new Error('请求过于频繁，请稍后再试')
      } else if (error.code === 'ECONNABORTED') {
        throw new Error('请求超时，请检查网络连接')
      } else if (error.code === 'ERR_NETWORK') {
        throw new Error('网络连接异常，请检查网络设置')
      } else if (!navigator.onLine) {
        throw new Error('网络连接已断开，请检查网络设置')
      }
      throw new Error('获取系统状态失败：' + (error.message || '未知错误'))
    }
  },

  // 格式化系统状态数据
  formatSystemStatus(data) {
    return {
      cpu: {
        usage: data.cpu.usage,
        cores: data.cpu.cores,
        status: this.getResourceStatus(data.cpu.usage)
      },
      memory: {
        total: this.formatSize(data.memory.total),
        used: this.formatSize(data.memory.used),
        free: this.formatSize(data.memory.free),
        usage: (data.memory.used / data.memory.total) * 100,
        status: this.getResourceStatus((data.memory.used / data.memory.total) * 100)
      },
      disk: {
        total: this.formatSize(data.disk.total),
        used: this.formatSize(data.disk.used),
        free: this.formatSize(data.disk.free),
        usage: (data.disk.used / data.disk.total) * 100,
        status: this.getResourceStatus((data.disk.used / data.disk.total) * 100)
      },
      network: {
        in: this.formatSpeed(data.network.in),
        out: this.formatSpeed(data.network.out)
      },
      uptime: this.formatUptime(data.uptime)
    }
  },

  // 获取资源状态
  getResourceStatus(usage) {
    if (usage >= 90) return 'danger'
    if (usage >= 70) return 'warning'
    return 'normal'
  },

  // 格式化存储大小
  formatSize(size) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let value = size
    let unitIndex = 0

    while (value >= 1024 && unitIndex < units.length - 1) {
      value /= 1024
      unitIndex++
    }

    return `${value.toFixed(2)} ${units[unitIndex]}`
  },

  // 格式化网络速度
  formatSpeed(speed) {
    const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
    let value = speed
    let unitIndex = 0

    while (value >= 1024 && unitIndex < units.length - 1) {
      value /= 1024
      unitIndex++
    }

    return `${value.toFixed(2)} ${units[unitIndex]}`
  },

  // 格式化运行时间
  formatUptime(seconds) {
    const days = Math.floor(seconds / 86400)
    const hours = Math.floor((seconds % 86400) / 3600)
    const minutes = Math.floor((seconds % 3600) / 60)

    const parts = []
    if (days > 0) parts.push(`${days}天`)
    if (hours > 0) parts.push(`${hours}小时`)
    if (minutes > 0) parts.push(`${minutes}分钟`)

    return parts.join(' ')
  }
}