// 视频配置
export const videoConfig = {
  // 视频白名单配置
  videoInWhiteList: {
    enabled: true,
    formats: ['mp4', 'webm', 'avi'],
    maxSize: 100, // MB
    maxDuration: 3600 // seconds
  },
  
  // 视频处理配置
  processing: {
    quality: 'high',
    compression: true,
    autoConvert: true
  },
  
  // 视频存储配置
  storage: {
    location: 'local', // local or cloud
    path: '/uploads/videos',
    backup: false
  }
}