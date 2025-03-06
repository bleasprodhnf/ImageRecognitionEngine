<template>
  <div class="error-page">
    <div class="error-container">
      <el-result
        :icon="errorIcon"
        :title="errorTitle"
        :sub-title="errorMessage"
      >
        <template #extra>
          <el-button type="primary" @click="handleRetry">重试</el-button>
          <el-button @click="handleBack">返回</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ErrorTypes } from '@/services/errorService'

export default {
  name: 'ErrorPage',
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    const errorType = ref(route.query.type || ErrorTypes.SERVER_ERROR)
    const errorMessage = ref(route.query.message || '发生未知错误')

    // 根据错误类型计算显示的图标和标题
    const errorIcon = computed(() => {
      switch (errorType.value) {
        case ErrorTypes.AUTH_ERROR:
          return 'warning'
        case ErrorTypes.PERMISSION_ERROR:
          return 'error'
        case ErrorTypes.RATE_LIMIT_ERROR:
          return 'warning'
        case ErrorTypes.NOT_FOUND_ERROR:
          return 'info'
        default:
          return 'error'
      }
    })

    const errorTitle = computed(() => {
      switch (errorType.value) {
        case ErrorTypes.AUTH_ERROR:
          return '认证失败'
        case ErrorTypes.PERMISSION_ERROR:
          return '权限不足'
        case ErrorTypes.RATE_LIMIT_ERROR:
          return '请求频率限制'
        case ErrorTypes.NOT_FOUND_ERROR:
          return '资源不存在'
        case ErrorTypes.NETWORK_ERROR:
          return '网络错误'
        default:
          return '系统错误'
      }
    })

    // 重试按钮处理函数
    const handleRetry = () => {
      if (route.query.from) {
        router.push(route.query.from)
      } else {
        router.go(-1)
      }
    }

    // 返回按钮处理函数
    const handleBack = () => {
      router.push('/')
    }

    return {
      errorIcon,
      errorTitle,
      errorMessage,
      handleRetry,
      handleBack
    }
  }
}
</script>

<style scoped>
.error-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f7fa;
}

.error-container {
  max-width: 800px;
  padding: 40px;
  text-align: center;
}
</style>