<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Monitor,
  Setting,
  DataLine,
  User,
  Connection,
  Picture
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const activeMenu = ref(route.path)

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/client/profile')
      break
    case 'password':
      router.push('/client/profile?tab=password')
      break
    case 'logout':
      localStorage.removeItem('client-token')
      router.push('/client/login')
      break
  }
}
</script>

<template>
  <el-container class="layout-container">
    <el-header height="60px">
      <div class="header-left">
        <h2 class="system-title">图像识别引擎客户系统</h2>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleCommand">
          <span class="user-info">
            <el-avatar :size="32" icon="UserFilled" />
            <span class="username">客户</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人信息</el-dropdown-item>
              <el-dropdown-item command="password">修改密码</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container class="main-container">
      <el-aside width="200px" class="sidebar">
        <el-menu
          :default-active="activeMenu"
          class="sidebar-menu"
          router
        >
          <el-menu-item index="/client/dashboard">
            <el-icon><Monitor /></el-icon>
            <span>控制台</span>
          </el-menu-item>

          <el-menu-item index="/client/profile">
            <el-icon><Setting /></el-icon>
            <span>账户设置</span>
          </el-menu-item>

          <el-menu-item index="/client/usage">
            <el-icon><DataLine /></el-icon>
            <span>使用统计</span>
          </el-menu-item>
          
          <el-menu-item index="/client/api-settings">
            <el-icon><Connection /></el-icon>
            <span>API接口管理</span>
          </el-menu-item>
          
          <el-menu-item index="/client/image-labels">
            <el-icon><Picture /></el-icon>
            <span>图片标识库</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      
      <el-main class="content-main">
        <div class="content-wrapper">
          <router-view />
        </div>
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-container {
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  margin: 0 auto;
  max-width: 1200px;
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  top: 0;
  bottom: 0;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  background-color: #fff;
}

.main-container {
  flex: 1;
  overflow: hidden;
  display: flex;
  margin: 0 auto;
  width: 100%;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  width: 100%;
  position: relative;
  z-index: 1000;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.system-title {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
  color: #606266;
}

.sidebar {
  background-color: #001529;
  height: 100%;
  box-shadow: 2px 0 6px rgba(0,21,41,.35);
}

.sidebar-menu {
  border-right: none;
  height: 100%;
  background-color: #001529;
}

.el-menu {
  border-right: none;
  background-color: #001529;
}

:deep(.el-menu-item) {
  color: rgba(255,255,255,0.75) !important;
  height: 50px;
  line-height: 50px;
  padding: 0 20px;
  transition: all 0.3s ease;
  background-color: #002140 !important;
}

:deep(.el-menu-item:hover) {
  color: #fff !important;
  background-color: #003366 !important;
}

:deep(.el-menu-item.is-active) {
  background-color: #003366 !important;
  color: #fff !important;
  position: relative;
}

:deep(.el-menu-item.is-active::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 3px;
  background: #fff;
}

.content-main {
  background-color: #f0f2f5;
  padding: 20px;
  overflow-y: auto;
}

.content-wrapper {
  margin: 0 auto;
  width: 100%;
  background-color: transparent;
}

/* 响应式布局 */
@media screen and (min-width: 1200px) {
  .content-wrapper {
    min-width: 1024px;
  }
}

@media screen and (min-width: 992px) and (max-width: 1199px) {
  .content-wrapper {
    min-width: 900px;
  }
}

@media screen and (min-width: 768px) and (max-width: 991px) {
  .content-wrapper {
    min-width: 700px;
  }
}

@media screen and (max-width: 767px) {
  .content-wrapper {
    min-width: 320px;
  }
}
</style>