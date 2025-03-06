import request from './request'

// 基础API路径
const BASE_PATH = '/api/v1'

// 模型版本相关接口
export const modelApi = {
  // 获取模型版本列表
  getVersions() {
    return request({
      url: `${BASE_PATH}/models/versions`,
      method: 'get'
    })
  },

  // 获取模型版本详情
  getVersionDetail(versionId) {
    return request({
      url: `${BASE_PATH}/models/versions/${versionId}`,
      method: 'get'
    })
  },

  // 创建新的模型版本
  createVersion(data) {
    return request({
      url: `${BASE_PATH}/models/versions`,
      method: 'post',
      data
    })
  },

  // 更新模型版本
  updateVersion(versionId, data) {
    return request({
      url: `${BASE_PATH}/models/versions/${versionId}`,
      method: 'put',
      data
    })
  },

  // 删除模型版本
  deleteVersion(versionId) {
    return request({
      url: `${BASE_PATH}/models/versions/${versionId}`,
      method: 'delete'
    })
  },

  // 获取模型性能监控数据
  getPerformanceMetrics(versionId, params) {
    return request({
      url: `${BASE_PATH}/models/${versionId}/metrics`,
      method: 'get',
      params
    })
  },

  // 获取模型参数配置
  getModelParams(versionId) {
    return request({
      url: `${BASE_PATH}/models/${versionId}/params`,
      method: 'get'
    })
  },

  // 更新模型参数配置
  updateModelParams(versionId, data) {
    return request({
      url: `${BASE_PATH}/models/${versionId}/params`,
      method: 'put',
      data
    })
  }
}