import request from './request'

// 基础API路径
const BASE_PATH = '/api/v1'

// 图像识别相关接口
export const recognitionApi = {
  // 获取识别历史记录
  getHistory(params) {
    return request({
      url: `${BASE_PATH}/recognition/history`,
      method: 'get',
      params
    })
  },

  // 获取识别准确率统计
  getAccuracyStats(params) {
    return request({
      url: `${BASE_PATH}/recognition/accuracy`,
      method: 'get',
      params
    })
  },

  // 提交图像识别请求
  submitRecognition(data) {
    return request({
      url: `${BASE_PATH}/recognition/submit`,
      method: 'post',
      data
    })
  },

  // 获取识别结果
  getRecognitionResult(requestId) {
    return request({
      url: `${BASE_PATH}/recognition/result/${requestId}`,
      method: 'get'
    })
  },

  // 获取识别配置
  getRecognitionConfig() {
    return request({
      url: `${BASE_PATH}/recognition/config`,
      method: 'get'
    })
  },

  // 更新识别配置
  updateRecognitionConfig(data) {
    return request({
      url: `${BASE_PATH}/recognition/config`,
      method: 'put',
      data
    })
  }
}