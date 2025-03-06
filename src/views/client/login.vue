<template>
  <div class="login-page">
    <div class="login-container">
      <el-card class="login-card">
        <template #header>
          <div class="card-header">
            <div class="logo">
              <el-icon :size="40" class="logo-icon"><Picture /></el-icon>
            </div>
            <h2>客户登录</h2>
          </div>
        </template>
        
        <el-form
          ref="loginFormRef"
          :model="loginForm"
          :rules="rules"
          label-position="left"
          label-width="60px"
        >
          <el-form-item label="账号" prop="username">
            <el-input
              v-model="loginForm.username"
              placeholder="请输入账号"
              prefix-icon="User"
            />
          </el-form-item>
          
          <el-form-item label="密码" prop="password">
            <el-input
              v-model="loginForm.password"
              type="password"
              placeholder="请输入密码"
              prefix-icon="Lock"
              show-password
            />
          </el-form-item>
          
          <el-form-item>
            <el-button
              type="primary"
              :loading="loading"
              class="login-button"
              @click="handleLogin"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock, Picture } from '@element-plus/icons-vue'

const router = useRouter()
const loading = ref(false)
const loginFormRef = ref(null)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    loading.value = true
    await loginFormRef.value.validate()
    // TODO: 实现客户登录逻辑
    localStorage.setItem('client-token', 'demo-client-token')
    router.push('/client/dashboard')
  } catch (error) {
    console.error('登录验证失败:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  width: 100vw;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background: linear-gradient(135deg, #e6f3ff 0%, #2d8cf0 100%);
}

.login-container {
  width: 100%;
  padding: 20px;
  box-sizing: border-box;
}

.login-card {
  max-width: 400px;
  width: 100%;
  margin: 0 auto;
}

.card-header {
  text-align: center;
}

.logo {
  margin-bottom: 15px;
}

.card-header h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.login-button {
  width: 100%;
  margin-top: 20px;
}

.logo-icon {
  background: linear-gradient(135deg, #409EFF 0%, #2d8cf0 100%);
  border-radius: 50%;
  padding: 10px;
  transition: transform 0.3s ease;
}

.logo-icon:hover {
  transform: scale(1.1);
}
</style>