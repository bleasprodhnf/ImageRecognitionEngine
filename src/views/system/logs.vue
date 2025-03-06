<template>
  <div class="logs-container">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>系统日志</span>
          <div class="search-box">
            <el-select v-model="searchType" placeholder="日志类型" style="width: 120px; margin-right: 10px">
              <el-option label="全部" value="" />
              <el-option label="操作日志" value="operation" />
              <el-option label="安全日志" value="security" />
              <el-option label="错误日志" value="error" />
            </el-select>
            <el-input
              v-model="searchKeyword"
              placeholder="搜索用户/操作"
              style="width: 200px; margin-right: 10px"
            />
            <el-button type="primary" @click="handleSearch">搜索</el-button>
          </div>
        </div>
      </template>
      <el-table :data="filteredLogs" style="width: 100%">
        <el-table-column prop="type" label="日志类型">
          <template #default="{ row }">
            <el-tag
              :type="row.type === 'operation' ? 'primary' : row.type === 'security' ? 'success' : 'danger'"
            >
              {{ row.type === 'operation' ? '操作' : row.type === 'security' ? '安全' : '错误' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="user" label="操作用户" />
        <el-table-column prop="action" label="操作内容" />
        <el-table-column prop="ip" label="IP地址" />
        <el-table-column prop="time" label="操作时间" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="totalLogs"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const searchType = ref('')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const totalLogs = ref(100)

// 模拟日志数据
const logs = ref([
  {
    id: 1,
    type: 'operation',
    user: 'admin',
    action: '更新系统配置',
    ip: '192.168.1.100',
    time: '2024-01-20 10:30:00',
    status: 'success'
  },
  {
    id: 2,
    type: 'security',
    user: 'system',
    action: '用户登录',
    ip: '192.168.1.101',
    time: '2024-01-20 10:35:00',
    status: 'success'
  },
  {
    id: 3,
    type: 'error',
    user: 'user1',
    action: '数据库连接失败',
    ip: '192.168.1.102',
    time: '2024-01-20 10:40:00',
    status: 'error'
  }
])

const filteredLogs = computed(() => {
  return logs.value.filter(log => {
    const typeMatch = !searchType.value || log.type === searchType.value
    const keywordMatch = !searchKeyword.value || 
      log.user.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      log.action.toLowerCase().includes(searchKeyword.value.toLowerCase())
    return typeMatch && keywordMatch
  })
})

const handleSearch = () => {
  currentPage.value = 1
  // 实际应用中这里应该调用API获取数据
}

const handleSizeChange = (val) => {
  pageSize.value = val
  // 实际应用中这里应该调用API获取数据
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  // 实际应用中这里应该调用API获取数据
}
</script>

<style scoped>
.logs-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-content {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 300px;
}
</style>