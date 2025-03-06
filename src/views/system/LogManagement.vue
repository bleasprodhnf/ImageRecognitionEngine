<template>
  <div class="log-management">
    <div class="filter-section">
      <a-form layout="inline">
        <a-form-item label="日志类型">
          <a-select v-model:value="filters.type" style="width: 120px">
            <a-select-option value="all">全部</a-select-option>
            <a-select-option value="operation">操作日志</a-select-option>
            <a-select-option value="system">系统日志</a-select-option>
            <a-select-option value="error">错误日志</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model:value="filters.dateRange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" @click="handleSearch">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button type="primary" style="margin-left: 8px" @click="handleExport">导出日志</a-button>
        </a-form-item>
      </a-form>
    </div>

    <a-table
      :columns="columns"
      :data-source="logData"
      :loading="loading"
      :pagination="pagination"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'type'">
          <a-tag :color="getTypeColor(record.type)">
            {{ record.type }}
          </a-tag>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'

const filters = reactive({
  type: 'all',
  dateRange: []
})

const loading = ref(false)
const logData = ref([])

const columns = [
  {
    title: '时间',
    dataIndex: 'timestamp',
    key: 'timestamp',
    sorter: true
  },
  {
    title: '类型',
    dataIndex: 'type',
    key: 'type'
  },
  {
    title: '用户',
    dataIndex: 'user',
    key: 'user'
  },
  {
    title: '操作/事件',
    dataIndex: 'action',
    key: 'action'
  },
  {
    title: '详细信息',
    dataIndex: 'details',
    key: 'details',
    ellipsis: true
  },
  {
    title: 'IP地址',
    dataIndex: 'ip',
    key: 'ip'
  }
]

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0,
  showSizeChanger: true,
  showTotal: (total) => `共 ${total} 条记录`
})

const getTypeColor = (type) => {
  const colorMap = {
    operation: 'blue',
    system: 'green',
    error: 'red'
  }
  return colorMap[type] || 'default'
}

const handleSearch = async () => {
  loading.value = true
  try {
    // TODO: 实现日志查询接口
    // const response = await fetchLogs(filters)
    // logData.value = response.data
    // pagination.total = response.total
  } catch (error) {
    message.error('获取日志数据失败')
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  filters.type = 'all'
  filters.dateRange = []
  handleSearch()
}

const handleExport = async () => {
  try {
    // TODO: 实现日志导出接口
    // await exportLogs(filters)
    message.success('日志导出成功')
  } catch (error) {
    message.error('日志导出失败')
  }
}

const handleTableChange = (pag, filters, sorter) => {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  // TODO: 根据排序和分页重新获取数据
  handleSearch()
}

// 初始加载数据
handleSearch()
</script>

<style scoped>
.log-management {
  padding: 24px;
}

.filter-section {
  margin-bottom: 24px;
  padding: 24px;
  background: #fff;
  border-radius: 4px;
}
</style>