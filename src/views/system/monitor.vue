<template>
  <div class="monitor-container">
    <el-row :gutter="20">
      <!-- CPU监控 -->
      <el-col :span="12">
        <el-card class="monitor-card">
          <template #header>
            <div class="card-header">
              <h3>CPU状态</h3>
            </div>
          </template>
          <el-row>
            <el-col :span="12">
              <el-progress
                type="dashboard"
                :percentage="systemMonitor.cpu.usage"
                :color="getProgressColor"
              />
              <div class="metric-label">使用率</div>
            </el-col>
            <el-col :span="12">
              <div class="metric-item">
                <div class="metric-value">{{ systemMonitor.cpu.cores }}</div>
                <div class="metric-label">核心数</div>
              </div>
              <div class="metric-item">
                <div class="metric-value">{{ systemMonitor.cpu.temperature }}°C</div>
                <div class="metric-label">温度</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 内存监控 -->
      <el-col :span="12">
        <el-card class="monitor-card">
          <template #header>
            <div class="card-header">
              <h3>内存状态</h3>
            </div>
          </template>
          <el-row>
            <el-col :span="12">
              <el-progress
                type="dashboard"
                :percentage="Math.round((systemMonitor.memory.used / systemMonitor.memory.total) * 100)"
                :color="getProgressColor"
              />
              <div class="metric-label">使用率</div>
            </el-col>
            <el-col :span="12">
              <div class="metric-item">
                <div class="metric-value">{{ formatMemory(systemMonitor.memory.used) }}</div>
                <div class="metric-label">已用内存</div>
              </div>
              <div class="metric-item">
                <div class="metric-value">{{ formatMemory(systemMonitor.memory.free) }}</div>
                <div class="metric-label">可用内存</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 磁盘监控 -->
      <el-col :span="12">
        <el-card class="monitor-card">
          <template #header>
            <div class="card-header">
              <h3>磁盘状态</h3>
            </div>
          </template>
          <el-row>
            <el-col :span="12">
              <el-progress
                type="dashboard"
                :percentage="Math.round((systemMonitor.disk.used / systemMonitor.disk.total) * 100)"
                :color="getProgressColor"
              />
              <div class="metric-label">使用率</div>
            </el-col>
            <el-col :span="12">
              <div class="metric-item">
                <div class="metric-value">{{ formatStorage(systemMonitor.disk.used) }}</div>
                <div class="metric-label">已用空间</div>
              </div>
              <div class="metric-item">
                <div class="metric-value">{{ formatStorage(systemMonitor.disk.free) }}</div>
                <div class="metric-label">可用空间</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>

      <!-- 网络监控 -->
      <el-col :span="12">
        <el-card class="monitor-card">
          <template #header>
            <div class="card-header">
              <h3>网络状态</h3>
            </div>
          </template>
          <el-row>
            <el-col :span="12">
              <div class="metric-item">
                <el-icon :size="24" color="#409EFF"><Upload /></el-icon>
                <div class="metric-value">{{ systemMonitor.network.upload }} MB/s</div>
                <div class="metric-label">上传速度</div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="metric-item">
                <el-icon :size="24" color="#67C23A"><Download /></el-icon>
                <div class="metric-value">{{ systemMonitor.network.download }} MB/s</div>
                <div class="metric-label">下载速度</div>
              </div>
            </el-col>
          </el-row>
          <div class="connections">
            <el-icon><Connection /></el-icon>
            <span>当前连接数：{{ systemMonitor.network.connections }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { systemMonitor as mockSystemMonitor } from '@/api/mock'

// 系统监控数据
const systemMonitor = ref({
  cpu: {
    usage: 0,
    cores: 0,
    temperature: 0
  },
  memory: {
    total: 0,
    used: 0,
    free: 0
  },
  disk: {
    total: 0,
    used: 0,
    free: 0
  },
  network: {
    upload: 0,
    download: 0,
    connections: 0
  }
})

// 获取进度条颜色
const getProgressColor = (percentage) => {
  if (percentage < 70) return '#67C23A'
  if (percentage < 90) return '#E6A23C'
  return '#F56C6C'
}

// 格式化内存大小
const formatMemory = (bytes) => {
  const gb = bytes / (1024 * 1024 * 1024)
  return `${gb.toFixed(2)} GB`
}

// 格式化存储大小
const formatStorage = (bytes) => {
  const tb = bytes / (1024 * 1024 * 1024 * 1024)
  return `${tb.toFixed(2)} TB`
}

// 获取系统监控数据
const fetchSystemMonitor = async () => {
  try {
    // 使用mock数据
    systemMonitor.value = {
      cpu: mockSystemMonitor.cpu,
      memory: {
        total: mockSystemMonitor.memory.total * 1024 * 1024 * 1024, // 转换为字节
        used: mockSystemMonitor.memory.used * 1024 * 1024 * 1024,
        free: mockSystemMonitor.memory.free * 1024 * 1024 * 1024
      },
      disk: {
        total: mockSystemMonitor.disk.total * 1024 * 1024 * 1024 * 1024, // 转换为字节
        used: mockSystemMonitor.disk.used * 1024 * 1024 * 1024 * 1024,
        free: mockSystemMonitor.disk.free * 1024 * 1024 * 1024 * 1024
      },
      network: mockSystemMonitor.network
    }
  } catch (error) {
    ElMessage.error('获取系统监控数据失败')
    console.error('获取系统监控数据失败:', error)
  } finally {
    // 确保即使出错也会继续定时更新
    setTimeout(fetchSystemMonitor, 30000)
  }
}

onMounted(() => {
  fetchSystemMonitor()
})
</script>

<style scoped>
.monitor-container {
  padding: 20px;
}

.monitor-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.metric-item {
  text-align: center;
  margin: 10px 0;
}

.metric-value {
  font-size: 24px;
  color: #409EFF;
  margin-bottom: 5px;
}

.metric-label {
  font-size: 14px;
  color: #606266;
}
</style>