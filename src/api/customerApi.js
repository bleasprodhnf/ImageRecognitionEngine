import request from './request'

// 基础API路径
const BASE_PATH = '/api/v1'

// 客户管理相关接口
export const customerApi = {
  // 获取套餐列表
  getPackages() {
    return request({
      url: `${BASE_PATH}/packages`,
      method: 'get'
    })
  },

  // 获取套餐详情
  getPackageDetail(packageId) {
    return request({
      url: `${BASE_PATH}/packages/${packageId}`,
      method: 'get'
    })
  },

  // 获取服务监控数据
  getServiceMonitor() {
    return request({
      url: `${BASE_PATH}/monitor/service`,
      method: 'get'
    })
  },

  // 获取资源使用统计
  getResourceUsage(params) {
    return request({
      url: `${BASE_PATH}/monitor/resources`,
      method: 'get',
      params
    })
  },

  // 获取API调用记录
  getApiCallLogs(params) {
    return request({
      url: `${BASE_PATH}/monitor/api-logs`,
      method: 'get',
      params
    })
  },

  // 获取异常记录列表
  getErrorLogs(params) {
    return request({
      url: `${BASE_PATH}/monitor/errors`,
      method: 'get',
      params
    })
  },

  // 更新客户配置
  updateCustomerConfig(data) {
    return request({
      url: `${BASE_PATH}/customer/config`,
      method: 'put',
      data
    })
  },

  // 获取API密钥
  getApiKeys() {
    return request({
      url: `${BASE_PATH}/customer/api-keys`,
      method: 'get'
    })
  },

  // 生成新的API密钥
  generateApiKey() {
    return request({
      url: `${BASE_PATH}/customer/api-keys`,
      method: 'post'
    })
  },

  // 删除API密钥
  deleteApiKey(keyId) {
    return request({
      url: `${BASE_PATH}/customer/api-keys/${keyId}`,
      method: 'delete'
    })
  }
}