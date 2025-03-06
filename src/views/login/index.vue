<template>
  <div class="login-page">
    <!-- 登录页面背景容器，设置全屏显示 -->
    <div class="login-container">
      <!-- 登录卡片容器，控制最大宽度确保在大屏幕上不会过宽 -->
      <el-card class="login-card">
        <h2 class="login-title">系统登录</h2>
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
          <!-- 用户名输入框，宽度通过父容器控制 -->
          <el-form-item prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入用户名"
              prefix-icon="User"
            />
          </el-form-item>
          <!-- 密码输入框，宽度跟随用户名输入框 -->
          <el-form-item prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          <!-- 登录按钮容器，控制按钮宽度为100% -->
          <div class="button-container">
            <el-button type="primary" @click="handleLogin" :loading="authStore.loading" class="login-button">
              登录
            </el-button>
          </div>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { User, Lock } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loginFormRef = ref(null)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    await loginFormRef.value.validate()
    
    // 使用 Pinia 状态管理中的登录方法
    await authStore.login({
      username: loginForm.username,
      password: loginForm.password
    })
    
    ElMessage.success('登录成功')
    router.push('/admin/dashboard')
  } catch (error) {
    console.error('登录失败:', error)
    ElMessage.error(authStore.error || '登录失败，请检查用户名和密码')
  }
}
</script>

<style scoped>
/* 登录页面根容器，设置全屏显示 */
.login-page {
  width: 100vw; /* 设置宽度为视口宽度 */
  height: 100vh; /* 设置高度为视口高度 */
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #e6f3ff 0%, #2d8cf0 100%);
}

/* 登录内容容器，控制整体布局和响应式宽度 */
.login-container {
  width: 100%; /* 容器宽度100%适应父元素 */
  padding: 20px;
  box-sizing: border-box;
}

/* 登录卡片样式，设置最大宽度和响应式宽度 */
.login-card {
  max-width: 400px; /* 设置最大宽度，避免在大屏幕上过宽 */
  width: 100%; /* 在小屏幕上自适应宽度 */
  margin: 0 auto; /* 水平居中 */
  border-radius: 8px;
}

/* 登录标题样式 */
.login-title {
  text-align: center;
  margin-bottom: 30px;
  color: #303133;
  font-size: 24px;
}

/* 按钮容器样式，控制按钮宽度 */
.button-container {
  margin-top: 20px;
  width: 100%; /* 确保按钮容器占满表单宽度 */
}

/* 登录按钮样式，设置为100%宽度 */
.login-button {
  width: 100%; /* 按钮宽度占满容器 */
}

/* 响应式布局：在小屏幕设备上调整内边距 */
@media screen and (max-width: 480px) {
  .login-container {
    padding: 10px; /* 减小内边距 */
  }
}
</style>