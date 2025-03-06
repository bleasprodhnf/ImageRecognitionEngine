import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const clientToken = ref(localStorage.getItem('client-token') || '')
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)
  const isClientAuthenticated = computed(() => !!clientToken.value)
  const userRole = computed(() => user.value?.role || '')

  // 方法
  function setToken(newToken) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function setClientToken(newToken) {
    clientToken.value = newToken
    localStorage.setItem('client-token', newToken)
  }

  function setUser(userData) {
    user.value = userData
  }

  async function login(credentials) {
    loading.value = true
    error.value = null
    
    try {
      // 这里应该调用实际的登录API
      // const response = await api.login(credentials)
      
      // 模拟登录成功
      const mockResponse = {
        token: 'mock-admin-token-' + Date.now(),
        user: {
          id: 1,
          username: credentials.username,
          role: 'admin'
        }
      }
      
      setToken(mockResponse.token)
      setUser(mockResponse.user)
      
      return mockResponse
    } catch (err) {
      error.value = err.message || '登录失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function clientLogin(credentials) {
    loading.value = true
    error.value = null
    
    try {
      // 这里应该调用实际的客户端登录API
      // const response = await api.clientLogin(credentials)
      
      // 模拟登录成功
      const mockResponse = {
        token: 'mock-client-token-' + Date.now(),
        user: {
          id: 100,
          username: credentials.username,
          role: 'client'
        }
      }
      
      setClientToken(mockResponse.token)
      setUser(mockResponse.user)
      
      return mockResponse
    } catch (err) {
      error.value = err.message || '登录失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    router.push('/admin/login')
  }

  function clientLogout() {
    clientToken.value = ''
    user.value = null
    localStorage.removeItem('client-token')
    router.push('/client/login')
  }

  return {
    // 状态
    token,
    clientToken,
    user,
    loading,
    error,
    
    // 计算属性
    isAuthenticated,
    isClientAuthenticated,
    userRole,
    
    // 方法
    login,
    clientLogin,
    logout,
    clientLogout,
    setUser
  }
})
