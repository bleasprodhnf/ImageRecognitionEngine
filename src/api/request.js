import axios from 'axios'
import { errorService } from '../services/errorService'
import { responseService } from '../services/responseService'

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 获取API认证信息
    const appId = localStorage.getItem('appId')
    const apiKey = localStorage.getItem('apiKey')
    const token = localStorage.getItem('token')
    
    // 如果存在API认证信息，添加到请求头
    if (appId && apiKey) {
      config.headers['X-App-ID'] = appId
      config.headers['X-API-Key'] = apiKey
    } else if (token) {
      // 如果是管理员接口，添加JWT token
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    // 使用responseService处理成功响应
    return responseService.formatResponse(response.data)
  },
  error => {
    // 使用统一的错误处理服务
    const errorResult = errorService.handleError(error)
    return Promise.reject(errorResult)
  }
)

export default request