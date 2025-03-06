// 错误处理服务

// 错误类型枚举
export const ErrorTypes = {
  AUTH_ERROR: 'AUTH_ERROR',
  PERMISSION_ERROR: 'PERMISSION_ERROR',
  RATE_LIMIT_ERROR: 'RATE_LIMIT_ERROR',
  NOT_FOUND_ERROR: 'NOT_FOUND_ERROR',
  SERVER_ERROR: 'SERVER_ERROR',
  NETWORK_ERROR: 'NETWORK_ERROR',
  VALIDATION_ERROR: 'VALIDATION_ERROR',
  BUSINESS_ERROR: 'BUSINESS_ERROR',
  TIMEOUT_ERROR: 'TIMEOUT_ERROR'
}

// 错误处理服务
export const errorService = {
  // 处理API认证错误
  handleAuthError(error) {
    const message = error.response?.data?.message;
    if (message?.includes('API')) {
      this.showError('API认证失败，请检查AppID和APIKey')
      // 清除失效的认证信息
      localStorage.removeItem('appId')
      localStorage.removeItem('apiKey')
    } else {
      this.showError('登录已过期，请重新登录')
      localStorage.removeItem('token')
      // 重定向到登录页
      window.location.href = '/login'
    }
    return { type: ErrorTypes.AUTH_ERROR, message: message || '认证失败' }
  },

  // 处理权限错误
  handlePermissionError(error) {
    this.showError('权限不足，无法访问该资源')
    return { type: ErrorTypes.PERMISSION_ERROR, message: error.response?.data?.message || '权限不足' }
  },

  // 处理请求频率限制
  handleRateLimitError(error) {
    this.showError('请求频率超过限制，请稍后再试')
    return { type: ErrorTypes.RATE_LIMIT_ERROR, message: error.response?.data?.message || '请求频率超限' }
  },

  // 处理资源不存在错误
  handleNotFoundError(error) {
    this.showError('请求的资源不存在')
    return { type: ErrorTypes.NOT_FOUND_ERROR, message: error.response?.data?.message || '资源不存在' }
  },

  // 处理服务器错误
  handleServerError(error) {
    const errorMessage = error.response?.data?.message || '服务器内部错误，请稍后重试'
    this.showError(errorMessage)
    return {
      type: ErrorTypes.SERVER_ERROR,
      message: errorMessage,
      code: error.response?.status || 500,
      details: error.response?.data?.details || null
    }
  },

  // 处理网络错误
  handleNetworkError(error) {
    const isTimeout = error.code === 'ECONNABORTED' || error.message?.includes('timeout')
    const errorMessage = isTimeout ? '请求超时，请稍后重试' : '网络连接错误，请检查网络设置'
    this.showError(errorMessage)
    return {
      type: isTimeout ? ErrorTypes.TIMEOUT_ERROR : ErrorTypes.NETWORK_ERROR,
      message: errorMessage,
      code: error.code || 'NETWORK_ERROR',
      isTimeout: isTimeout
    }
  },

  // 处理数据验证错误
  handleValidationError(error) {
    const message = error.response?.data?.message || '输入数据验证失败';
    this.showError(message)
    return { 
      type: ErrorTypes.VALIDATION_ERROR, 
      message,
      details: error.response?.data?.details || null
    }
  },

  // 处理业务逻辑错误
  handleBusinessError(error) {
    const message = error.response?.data?.message || '业务处理失败';
    this.showError(message)
    return { 
      type: ErrorTypes.BUSINESS_ERROR, 
      message,
      details: error.response?.data?.details || null
    }
  },

  // 统一错误处理入口
  handleError(error) {
    if (error.response) {
      const status = error.response.status
      switch (status) {
        case 401:
          return this.handleAuthError(error)
        case 403:
          return this.handlePermissionError(error)
        case 404:
          return this.handleNotFoundError(error)
        case 429:
          return this.handleRateLimitError(error)
        case 422:
          return this.handleValidationError(error)
        case 500:
          return this.handleServerError(error)
        default:
          if (error.response.data?.businessError) {
            return this.handleBusinessError(error)
          }
          return this.handleServerError(error)
      }
    } else if (error.code === 'ECONNABORTED' || error.message?.includes('timeout')) {
      return this.handleNetworkError(error)
    } else {
      return this.handleNetworkError(error)
    }
  },

  // 显示错误提示
  showError(message) {
    console.error(message)
    // 使用全局消息组件（如果存在）
    if (window.$message) {
      window.$message.error(message)
    } else if (window.ElMessage) {
      window.ElMessage.error(message)
    }
  }
}