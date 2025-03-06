// 任务处理服务
import request from '../api/request'
import { errorService } from './errorService'
import { responseService } from './responseService'

// 任务状态枚举
export const TaskStatus = {
  PENDING: 'pending',
  PROCESSING: 'processing',
  COMPLETED: 'completed',
  FAILED: 'failed',
  CANCELLED: 'cancelled'
}

// 任务处理服务
export const taskService = {
  // 存储任务状态回调
  statusCallbacks: new Map(),

  // 存储任务进度回调
  progressCallbacks: new Map(),

  // WebSocket连接实例
  ws: null,
  reconnectAttempts: 0,
  maxReconnectAttempts: 10, // 增加最大重连次数
  reconnectInterval: 3000, // 减少重连间隔
  isConnecting: false,
  connectionTimeout: 10000, // 添加连接超时时间
  lastReconnectTime: null, // 记录上次重连时间
  heartbeatInterval: null, // 心跳检测定时器
  heartbeatTimeout: null, // 心跳超时定时器
  connectionState: 'disconnected', // 连接状态：disconnected, connecting, connected, reconnecting

  // 初始化WebSocket连接
  initWebSocket() {
    if (this.isConnecting) return;
    this.isConnecting = true;
    this.connectionState = 'connecting';

    // 清理现有连接和定时器
    this.cleanupConnection();
    
    // 添加网络状态变化监听（确保不会重复添加）
    window.removeEventListener('online', this.handleNetworkChange);
    window.removeEventListener('offline', this.handleNetworkChange);
    window.addEventListener('online', this.handleNetworkChange);
    window.addEventListener('offline', this.handleNetworkChange);
    
    // 先检查网络状态
    if (!navigator.onLine) {
      console.warn('网络连接不可用，延迟WebSocket连接');
      this.isConnecting = false;
      this.connectionState = 'disconnected';
      
      // 添加网络恢复监听（使用一次性监听器避免重复添加）
      const handleNetworkRestore = () => {
        window.removeEventListener('online', handleNetworkRestore);
        console.log('网络已恢复，尝试建立WebSocket连接');
        // 短暂延迟确保网络稳定
        setTimeout(() => {
          this.initWebSocket();
        }, 1000);
      };
      window.addEventListener('online', handleNetworkRestore);
      
      return;
    }

    try {
      console.log('正在建立WebSocket连接...');
      this.ws = new WebSocket(process.env.VUE_APP_WS_URL || 'ws://localhost:8080/ws');

      this.ws.onopen = () => {
        console.log('WebSocket连接已建立');
        this.isConnecting = false;
        this.reconnectAttempts = 0;
        this.connectionState = 'connected';
        
        // 重新订阅所有活跃任务
        this.resubscribeActiveTasks();
        
        // 启动心跳检测
        this.startHeartbeat();
      }

      this.ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data);
          
          // 处理心跳响应
          if (data.type === 'pong') {
            this.handleHeartbeatResponse();
            return;
          }
          
          switch (data.type) {
            case 'status':
              this.handleStatusUpdate(data.taskId, data.status);
              break;
            case 'progress':
              this.handleProgressUpdate(data.taskId, data.progress);
              break;
            case 'error':
              console.error('服务器错误:', data.message);
              errorService.handleServerError({ response: { data: { message: data.message } } });
              break;
            default:
              console.warn('未知的WebSocket消息类型:', data.type);
          }
        } catch (error) {
          console.error('处理WebSocket消息时出错:', error);
          errorService.handleServerError(error);
        }
      }

      this.ws.onclose = (event) => {
        console.log(`WebSocket连接已关闭，代码: ${event.code}，原因: ${event.reason || '未知'}`);
        this.isConnecting = false;
        this.connectionState = 'disconnected';
        this.cleanupConnection();
        this.handleReconnect();
      }

      this.ws.onerror = (error) => {
        console.error('WebSocket错误:', error);
        this.isConnecting = false;
        this.connectionState = 'disconnected';
        errorService.handleNetworkError(error);
      }
      
      // 设置连接超时处理
      const connectionTimeoutHandler = setTimeout(() => {
        if (this.connectionState !== 'connected') {
          console.error('WebSocket连接超时');
          if (this.ws) {
            this.ws.close();
          }
          this.isConnecting = false;
          this.connectionState = 'disconnected';
          this.handleReconnect();
        }
      }, this.connectionTimeout);
      
      // 连接成功后清除超时处理
      this.ws.addEventListener('open', () => {
        clearTimeout(connectionTimeoutHandler);
      });
      
    } catch (error) {
      console.error('WebSocket连接创建失败:', error);
      this.isConnecting = false;
      this.connectionState = 'disconnected';
      errorService.handleNetworkError(error);
      this.handleReconnect();
    }
  },
  
  // 清理连接和相关定时器
  cleanupConnection() {
    // 清理现有WebSocket连接
    if (this.ws) {
      // 移除所有事件监听器以避免内存泄漏
      this.ws.onopen = null;
      this.ws.onmessage = null;
      this.ws.onclose = null;
      this.ws.onerror = null;
      
      try {
        this.ws.close();
      } catch (e) {
        console.warn('关闭WebSocket连接时出错:', e);
      }
      this.ws = null;
    }
    
    // 清理心跳相关定时器
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
    
    if (this.heartbeatTimeout) {
      clearTimeout(this.heartbeatTimeout);
      this.heartbeatTimeout = null;
    }
    
    // 清理重连相关定时器
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
      this.reconnectTimeoutId = null;
    }
    
    if (this.reconnectTimerId) {
      clearTimeout(this.reconnectTimerId);
      this.reconnectTimerId = null;
    }
  },
  
  // 启动心跳检测
  startHeartbeat() {
    // 清理现有心跳定时器
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
    }
    
    // 设置新的心跳定时器，每30秒发送一次
    this.heartbeatInterval = setInterval(() => {
      this.sendHeartbeat();
    }, 30000);
    
    // 立即发送一次心跳检测
    this.sendHeartbeat();
  },
  
  // 发送心跳包
  sendHeartbeat() {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      console.log('发送心跳检测...');
      
      try {
        // 发送心跳包
        this.ws.send(JSON.stringify({
          type: 'ping',
          timestamp: new Date().toISOString(),
          clientId: this.getClientId() // 添加客户端标识
        }));
        
        // 设置心跳超时处理
        if (this.heartbeatTimeout) {
          clearTimeout(this.heartbeatTimeout);
        }
        
        this.heartbeatTimeout = setTimeout(() => {
          console.warn('心跳检测超时，连接可能已断开');
          // 如果心跳超时，尝试重新连接
          if (this.ws) {
            this.ws.close();
            this.connectionState = 'disconnected';
            this.handleReconnect();
          }
        }, 5000); // 5秒内没有响应则认为连接已断开
      } catch (error) {
        console.error('发送心跳包失败:', error);
        this.connectionState = 'disconnected';
        this.handleReconnect();
      }
    } else if (this.connectionState !== 'reconnecting') {
      // 如果连接已关闭且不在重连过程中，尝试重新连接
      console.warn('心跳检测失败：WebSocket连接已关闭');
      this.connectionState = 'disconnected';
      this.handleReconnect();
    }
  },
  
  // 获取客户端唯一标识
  getClientId() {
    if (!this.clientId) {
      // 生成一个随机的客户端ID
      this.clientId = 'client_' + Math.random().toString(36).substring(2, 15);
    }
    return this.clientId;
  },
  
  // 处理心跳响应
  handleHeartbeatResponse() {
    console.log('收到心跳响应，连接正常');
    // 清除心跳超时定时器
    if (this.heartbeatTimeout) {
      clearTimeout(this.heartbeatTimeout);
      this.heartbeatTimeout = null;
    }
  },

  // 处理重连
  handleReconnect() {
    // 如果已经在重连过程中，避免重复触发
    if (this.connectionState === 'reconnecting') {
      return;
    }
    
    const now = Date.now();
    // 检查是否需要重置重连次数
    if (this.lastReconnectTime && (now - this.lastReconnectTime) > 60000) {
      this.reconnectAttempts = 0;
    }

    // 检查网络状态，如果浏览器报告网络已断开，增加等待时间
    const isOnline = navigator.onLine;
    if (!isOnline) {
      console.warn('网络连接已断开，延长重连等待时间');
      // 网络断开时，设置更长的重连间隔
      setTimeout(() => {
        // 重新检查网络状态
        if (navigator.onLine) {
          console.log('网络已恢复，尝试重新连接');
          this.reconnectAttempts = 0; // 重置重连次数
          this.initWebSocket();
        } else {
          this.handleReconnect(); // 网络仍然断开，继续等待
        }
      }, 30000); // 30秒后再次检查
      return;
    }

    if (this.reconnectAttempts < this.maxReconnectAttempts) {
      this.reconnectAttempts++;
      this.lastReconnectTime = now;
      this.connectionState = 'reconnecting';
      
      // 改进的指数退避策略，增加随机因子避免多客户端同时重连
      const randomFactor = 0.5 * Math.random() + 0.5; // 0.5-1.0之间的随机数
      const delay = this.reconnectInterval * Math.min(Math.pow(1.5, this.reconnectAttempts - 1), 10) * randomFactor;
      console.log(`尝试重新连接 (${this.reconnectAttempts}/${this.maxReconnectAttempts})，延迟: ${Math.round(delay)}ms...`);
      
      // 设置重连超时，避免无限等待
      const reconnectTimeout = setTimeout(() => {
        console.warn('WebSocket重连超时，尝试下一次重连');
        this.connectionState = 'disconnected'; // 重置连接状态
        this.handleReconnect();
      }, this.connectionTimeout);
      
      // 存储超时定时器以便在需要时清理
      this.reconnectTimeoutId = reconnectTimeout;
      
      // 设置延迟重连定时器
      const reconnectTimerId = setTimeout(() => {
        if (this.reconnectTimeoutId) {
          clearTimeout(this.reconnectTimeoutId);
          this.reconnectTimeoutId = null;
        }
        this.initWebSocket();
      }, delay);
      
      // 存储重连定时器以便在需要时清理
      this.reconnectTimerId = reconnectTimerId;
    } else {
      console.error('WebSocket重连次数超过最大限制，暂时停止重连');
      this.connectionState = 'disconnected';
      errorService.handleNetworkError(new Error('WebSocket连接失败，请检查网络连接或刷新页面'));
      
      // 设置一个长时间后的最终重试
      this._finalReconnectTimer = setTimeout(() => {
        console.log('执行最终重连尝试...');
        this.reconnectAttempts = 0;
        this.initWebSocket();
        this._finalReconnectTimer = null;
      }, 300000); // 5分钟后尝试最终重连
      
      // 添加网络状态监听，在网络恢复时立即尝试重连
      const handleNetworkChange = () => {
        if (navigator.onLine && this.connectionState === 'disconnected') {
          console.log('网络已恢复，立即尝试重新连接');
          window.removeEventListener('online', handleNetworkChange);
          this.reconnectAttempts = 0;
          this.initWebSocket();
        }
      };
      window.addEventListener('online', handleNetworkChange);
      
      // 通知用户连接状态
      if (this.statusCallbacks.size > 0 || this.progressCallbacks.size > 0) {
        // 如果有活跃的任务监听，通知用户连接已断开
        const errorMsg = {
          type: 'connection_error',
          message: '与服务器的连接已断开，系统将在后台尝试重新连接',
          reconnectIn: 300
        };
        
        // 通知所有活跃的状态回调
        this.statusCallbacks.forEach((callback) => {
          try {
            callback({
              status: 'connection_lost',
              timestamp: new Date().toISOString(),
              message: errorMsg.message
            });
          } catch (e) {
            console.error('通知连接状态失败:', e);
          }
        });
        
        // 同时通知进度回调
        this.progressCallbacks.forEach((callback) => {
          try {
            callback({
              percentage: -1, // 使用-1表示连接中断
              timestamp: new Date().toISOString(),
              message: errorMsg.message
            });
          } catch (e) {
            console.error('通知进度回调失败:', e);
          }
        });
      }
    }
  },

  // 重新订阅活跃任务
  resubscribeActiveTasks() {
    const activeTaskIds = new Set([...this.statusCallbacks.keys(), ...this.progressCallbacks.keys()]);
    activeTaskIds.forEach(taskId => {
      this.subscribeToTask(taskId);
    });
  },

  // 处理任务状态更新
  handleStatusUpdate(taskId, status) {
    try {
      if (!taskId || !status) {
        console.error('无效的任务状态更新:', { taskId, status })
        return
      }

      const callback = this.statusCallbacks.get(taskId)
      if (callback) {
        // 验证状态数据格式
        const validStatus = {
          ...status,
          timestamp: status.timestamp || new Date().toISOString(),
          message: status.message || ''
        }
        callback(validStatus)
      }

      if ([TaskStatus.COMPLETED, TaskStatus.FAILED, TaskStatus.CANCELLED].includes(status.status)) {
        this.cleanupTask(taskId)
      }
    } catch (error) {
      console.error('处理任务状态更新失败:', error)
      errorService.handleServerError(error)
    }
  },

  // 处理任务进度更新
  handleProgressUpdate(taskId, progress) {
    try {
      if (!taskId || !progress) {
        console.error('无效的任务进度更新:', { taskId, progress })
        return
      }

      const callback = this.progressCallbacks.get(taskId)
      if (callback) {
        // 验证进度数据格式
        const validProgress = {
          ...progress,
          percentage: Math.min(Math.max(progress.percentage || 0, 0), 100),
          timestamp: progress.timestamp || new Date().toISOString()
        }
        callback(validProgress)
      }
    } catch (error) {
      console.error('处理任务进度更新失败:', error)
      errorService.handleServerError(error)
    }
  },

  // 发送任务订阅消息
  subscribeToTask(taskId) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        type: 'subscribe',
        taskId: taskId
      }))
    }
  },

  // 取消任务订阅
  unsubscribeFromTask(taskId) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify({
        type: 'unsubscribe',
        taskId: taskId
      }))
    }
  },

  // 注册任务状态监听
  onTaskStatusChange(taskId, callback) {
    this.statusCallbacks.set(taskId, callback)
    this.subscribeToTask(taskId)
  },

  // 注册任务进度监听
  onTaskProgressChange(taskId, callback) {
    this.progressCallbacks.set(taskId, callback)
    this.subscribeToTask(taskId)
  },

  // 取消任务状态监听
  offTaskStatusChange(taskId) {
    this.statusCallbacks.delete(taskId)
    this.unsubscribeFromTask(taskId)
  },

  // 取消任务进度监听
  offTaskProgressChange(taskId) {
    this.progressCallbacks.delete(taskId)
    this.unsubscribeFromTask(taskId)
  },

  // 获取任务状态
  async getTaskStatus(taskId) {
    try {
      const response = await request({
        url: `/tasks/${taskId}/status`,
        method: 'get'
      })
      return responseService.handleSuccess(response)
    } catch (error) {
      return errorService.handleApiError(error)
    }
  },

  // 获取任务进度
  async getTaskProgress(taskId) {
    try {
      const response = await request({
        url: `/tasks/${taskId}/progress`,
        method: 'get'
      })
      return responseService.handleSuccess(response)
    } catch (error) {
      return errorService.handleApiError(error)
    }
  },

  // 开始任务状态轮询
  startStatusPolling(taskId, interval = 5000) {
    const pollStatus = async () => {
      const result = await this.getTaskStatus(taskId)
      if (result.success) {
        const callback = this.statusCallbacks.get(taskId)
        if (callback) {
          callback(result.data)
        }
        // 如果任务已完成或失败，停止轮询
        if ([TaskStatus.COMPLETED, TaskStatus.FAILED, TaskStatus.CANCELLED].includes(result.data.status)) {
          this.stopStatusPolling(taskId)
        }
      }
    }

    // 立即执行一次
    pollStatus()
    // 设置轮询定时器
    const timerId = setInterval(pollStatus, interval)
    // 存储定时器ID
    this.statusPollingTimers.set(taskId, timerId)
  },

  // 停止任务状态轮询
  stopStatusPolling(taskId) {
    const timerId = this.statusPollingTimers.get(taskId)
    if (timerId) {
      clearInterval(timerId)
      this.statusPollingTimers.delete(taskId)
    }
  },

  // 存储状态轮询定时器
  statusPollingTimers: new Map(),

  // 取消任务
  async cancelTask(taskId) {
    try {
      const response = await request({
        url: `/tasks/${taskId}/cancel`,
        method: 'post'
      })
      return responseService.handleSuccess(response)
    } catch (error) {
      return errorService.handleApiError(error)
    }
  },

  // 清理任务相关资源
  cleanupTask(taskId) {
    this.stopStatusPolling(taskId)
    this.offTaskStatusChange(taskId)
    this.offTaskProgressChange(taskId)
  },
  
  // 处理网络状态变化
  handleNetworkChange = (event) => {
    console.log(`任务服务网络状态变化: ${event.type}`)
    
    if (event.type === 'online') {
      // 网络恢复，尝试重新连接WebSocket
      if (this.connectionState === 'disconnected') {
        console.log('网络已恢复，尝试重新连接WebSocket')
        this.reconnectAttempts = 0 // 重置重连次数
        this.initWebSocket()
      }
    } else if (event.type === 'offline') {
      // 网络断开，记录状态但不做特殊处理
      // WebSocket会自动断开并触发onclose事件
      console.warn('网络已断开，WebSocket连接可能会中断')
    }
  },
  
  // 清理所有资源（用于组件卸载或页面关闭时调用）
  cleanup() {
    // 清理所有任务状态轮询
    this.statusPollingTimers.forEach((timerId, taskId) => {
      clearInterval(timerId)
    })
    this.statusPollingTimers.clear()
    
    // 清理WebSocket连接
    this.cleanupConnection()
    
    // 清理所有回调
    this.statusCallbacks.clear()
    this.progressCallbacks.clear()
    
    // 重置连接状态和相关变量
    this.reconnectAttempts = 0
    this.connectionState = 'disconnected'
    this.isConnecting = false
    this.lastReconnectTime = null
    this.clientId = null
    
    // 移除网络状态监听器
    window.removeEventListener('online', this.handleNetworkChange)
    window.removeEventListener('offline', this.handleNetworkChange)
    
    // 清理可能存在的最终重连定时器
    if (this._finalReconnectTimer) {
      clearTimeout(this._finalReconnectTimer)
      this._finalReconnectTimer = null
    }
    
    // 记录清理完成
    console.log('任务服务资源已完全清理')
    
    return true // 返回清理成功标志
  }
}