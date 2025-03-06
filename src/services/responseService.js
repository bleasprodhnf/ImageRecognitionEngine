// API响应处理服务

// 响应状态枚举
export const ResponseStatus = {
  SUCCESS: 200,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  VALIDATION_ERROR: 422,
  RATE_LIMIT: 429,
  SERVER_ERROR: 500
}

// 响应处理服务
export const responseService = {
  // 处理成功响应
  handleSuccess(response) {
    const { code, message, data } = response
    if (code === ResponseStatus.SUCCESS) {
      return {
        success: true,
        data,
        message: message || '操作成功'
      }
    }
    return this.handleError(response)
  },

  // 处理错误响应
  handleError(response) {
    const errorResponse = {
      success: false,
      data: response.data,
      message: response.message || '操作失败',
      code: response.code
    }

    // 添加错误类型标识
    if (response.code) {
      errorResponse.type = this.getErrorType(response.code)
    }

    return errorResponse
  },

  // 获取错误类型
  getErrorType(code) {
    switch (code) {
      case ResponseStatus.UNAUTHORIZED:
        return 'auth'
      case ResponseStatus.FORBIDDEN:
        return 'permission'
      case ResponseStatus.VALIDATION_ERROR:
        return 'validation'
      case ResponseStatus.RATE_LIMIT:
        return 'rateLimit'
      default:
        return 'error'
    }
  },

  // 格式化响应数据
  formatResponse(response) {
    if (!response) {
      return this.handleError({
        code: ResponseStatus.SERVER_ERROR,
        message: '服务器响应异常'
      })
    }

    // 处理分页数据
    if (response.data?.items && response.data?.total !== undefined) {
      return {
        success: true,
        data: {
          list: response.data.items,
          total: response.data.total,
          page: response.data.page || 1,
          pageSize: response.data.pageSize || 20
        },
        message: response.message || '获取数据成功'
      }
    }

    return this.handleSuccess(response)
  }
}