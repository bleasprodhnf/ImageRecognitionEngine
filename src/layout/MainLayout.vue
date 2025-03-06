<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  Fold,
  Expand,
  Setting,
  User,
  Monitor,
  DataLine,
  Service
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const isCollapse = ref(false)
const activeMenu = ref(route.path)

const toggleSidebar = () => {
  isCollapse.value = !isCollapse.value
}

const handleLogout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<template>
  <el-container class="layout-container">
    <el-header height="60px">
      <div class="header-left">
        <el-button link @click="toggleSidebar">
          <el-icon :size="20">
            <Fold v-if="isCollapse" />
            <Expand v-else />
          </el-icon>
        </el-button>
        <h2 class="system-title">图像识别引擎管理系统</h2>
      </div>
      <div class="header-right">
        <el-dropdown @command="handleLogout">
          <span class="user-info">
            <el-avatar :size="32" icon="UserFilled" />
            <span class="username">管理员</span>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>个人信息</el-dropdown-item>
              <el-dropdown-item>修改密码</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-header>
    
    <el-container class="main-container">
      <el-aside :width="isCollapse ? '64px' : '200px'" class="sidebar">
        <el-menu
          :collapse="isCollapse"
          :collapse-transition="false"
          :default-active="activeMenu"
          class="sidebar-menu"
          router
        >
          <el-sub-menu index="1">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>系统管理</span>
            </template>
            <el-menu-item index="/admin/system/params">系统参数配置</el-menu-item>
            <el-menu-item index="/admin/system/monitor">服务器监控</el-menu-item>
            <el-menu-item index="/admin/system/logs">日志管理</el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="2">
            <template #title>
              <el-icon><User /></el-icon>
              <span>用户权限管理</span>
            </template>
            <el-menu-item index="/admin/auth/admins">管理员管理</el-menu-item>
            <el-menu-item index="/admin/auth/roles">角色管理</el-menu-item>
            <el-menu-item index="/admin/auth/permissions">权限管理</el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="3">
            <template #title>
              <el-icon><Monitor /></el-icon>
              <span>模型管理</span>
            </template>
            <el-menu-item index="/admin/model/versions">模型版本管理</el-menu-item>
            <el-menu-item index="/admin/model/params">训练参数配置</el-menu-item>
            <el-menu-item index="/admin/model/monitor">性能监控</el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="4">
            <template #title>
              <el-icon><DataLine /></el-icon>
              <span>数据统计</span>
            </template>
            <el-menu-item index="/admin/stats/system">系统使用统计</el-menu-item>
            <el-menu-item index="/admin/stats/accuracy">识别准确率分析</el-menu-item>
            <el-menu-item index="/admin/stats/resources">资源使用报表</el-menu-item>
          </el-sub-menu>

          <el-sub-menu index="5">
            <template #title>
              <el-icon><Service /></el-icon>
              <span>客户管理</span>
            </template>
            <el-menu-item index="/admin/customer/accounts">客户账号管理</el-menu-item>
            <el-menu-item index="/admin/customer/packages">套餐配置管理</el-menu-item>
            <el-menu-item index="/admin/customer/service">客户服务监控</el-menu-item>
          </el-sub-menu>
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
  max-width: 1440px;
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
  background-color: #304156;
  transition: width 0.3s;
  height: 100%;
}

.sidebar-menu {
  border-right: none;
  height: 100%;
}

.el-menu {
  border-right: none;
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

/* 大屏幕设备 */
@media screen and (min-width: 1200px) {
  .content-wrapper {
    min-width: 1024px;
  }
}

/* 中等屏幕设备 */
@media screen and (min-width: 992px) and (max-width: 1199px) {
  .content-wrapper {
    min-width: 900px;
  }
}

/* 平板设备 */
@media screen and (min-width: 768px) and (max-width: 991px) {
  .content-wrapper {
    min-width: 700px;
  }
}

/* 手机设备 */
@media screen and (max-width: 767px) {
  .content-wrapper {
    min-width: 320px;
  }
}
</style>