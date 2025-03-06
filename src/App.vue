<template>
  <div id="app">
    <el-config-provider :locale="zhCn">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <keep-alive :include="cachedViews">
            <component :is="Component" />
          </keep-alive>
        </transition>
      </router-view>
    </el-config-provider>
    
    <!-- 全局加载状态 -->
    <div class="global-loading" v-if="globalLoading">
      <el-loading :fullscreen="true" />
    </div>
  </div>
</template>

<script setup>
import { ref, provide, onMounted } from 'vue'
import zhCn from 'element-plus/dist/locale/zh-cn.mjs'
import { useModelStore } from './stores/model'
import { useAuthStore } from './stores/auth'

// 缓存的视图组件
const cachedViews = ref(['Dashboard', 'Profile'])

// 全局加载状态
const globalLoading = ref(false)

// 提供全局加载状态控制
provide('setGlobalLoading', (loading) => {
  globalLoading.value = loading
})

// 获取状态管理实例
const modelStore = useModelStore()
const authStore = useAuthStore()

// 初始化应用
onMounted(async () => {
  // 如果已登录，预加载一些必要数据
  if (authStore.isAuthenticated) {
    try {
      await modelStore.fetchVersions()
    } catch (error) {
      console.error('初始化数据加载失败:', error)
    }
  }
})
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

#app {
  height: 100vh;
  margin: 0;
  padding: 0;
}

/* 页面过渡动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 全局滚动条样式 */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #c0c4cc;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: #909399;
}

/* 全局加载状态 */
.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0.7);
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
