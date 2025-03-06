import { createPinia } from 'pinia'

const pinia = createPinia()

export default pinia

// 导出所有存储，方便导入
export * from './auth'
export * from './model'
