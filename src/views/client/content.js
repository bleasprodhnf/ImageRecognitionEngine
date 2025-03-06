import { videoConfig } from '../../config/video.js'

// 导出视频配置
export const { videoInWhiteList } = videoConfig

// 其他内容相关配置
export const contentConfig = {
  maxFileSize: 10, // MB
  allowedTypes: ['image/jpeg', 'image/png', 'image/gif'],
  uploadPath: '/api/v1/upload',
  processingTimeout: 30000 // ms
}