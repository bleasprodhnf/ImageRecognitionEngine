<template>
  <div class="error-page">
    <div class="error-container">
      <h1 class="error-code">404</h1>
      <h2 class="error-title">页面未找到</h2>
      <p class="error-message">抱歉，您访问的页面不存在。</p>
      <el-button type="primary" @click="goBack">返回上一页</el-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goBack = () => {
  router.back()
}
</script>

<style scoped>
.error-page {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
}

.error-container {
  text-align: center;
  padding: 40px;
}

.error-code {
  font-size: 120px;
  color: #f56c6c;
  margin: 0;
  line-height: 1;
}

.error-title {
  font-size: 24px;
  color: #303133;
  margin: 20px 0;
}

.error-message {
  font-size: 16px;
  color: #606266;
  margin-bottom: 30px;
}
</style>