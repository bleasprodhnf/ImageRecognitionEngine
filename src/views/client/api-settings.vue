<template>
  <div class="api-settings-container">
    <el-card class="api-card">
      <template #header>
        <div class="card-header">
          <span>API接口信息</span>
          <el-button type="primary" size="small" @click="regenerateKeys">重新生成密钥</el-button>
        </div>
      </template>
      
      <div class="api-info">
        <div class="info-item">
          <div class="info-label">AppID</div>
          <div class="info-value-container">
            <div class="info-value">{{ apiInfo.appId }}</div>
            <el-button link @click="copyToClipboard(apiInfo.appId)">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </div>
        </div>
        
        <div class="info-item">
          <div class="info-label">API密钥</div>
          <div class="info-value-container">
            <div class="info-value">{{ showApiKey ? apiInfo.apiKey : '••••••••••••••••' }}</div>
            <el-button link @click="toggleApiKeyVisibility">
              <el-icon><View v-if="!showApiKey" /><Hide v-else /></el-icon>
            </el-button>
            <el-button link @click="copyToClipboard(apiInfo.apiKey)">
              <el-icon><CopyDocument /></el-icon>
            </el-button>
          </div>
        </div>
      </div>

      <div class="api-usage">
        <h3>API使用说明</h3>
        <el-tabs type="border-card">
          <el-tab-pane label="请求示例">
            <div class="code-block">
              <pre><code>// 请求示例 (JavaScript)
const response = await fetch('https://api.example.com/v1/recognize', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-App-ID': '{{ apiInfo.appId }}',
    'X-API-Key': '您的API密钥'
  },
  body: JSON.stringify({
    image_url: 'https://example.com/image.jpg',
    // 或者使用 base64 编码的图片数据
    // image_data: 'base64编码的图片数据'
  })
});

const result = await response.json();</code></pre>
            </div>
          </el-tab-pane>
          <el-tab-pane label="响应格式">
            <div class="code-block">
              <pre><code>// 响应示例
{
  "success": true,
  "data": {
    "id": "rec_12345",
    "results": [
      {
        "label": "对象1",
        "confidence": 0.98,
        "bbox": [10, 10, 100, 100]
      },
      {
        "label": "对象2",
        "confidence": 0.85,
        "bbox": [150, 50, 200, 200]
      }
    ],
    "processing_time": 0.45
  }
}</code></pre>
            </div>
          </el-tab-pane>
          <el-tab-pane label="错误码">
            <el-table :data="errorCodes" style="width: 100%">
              <el-table-column prop="code" label="错误码" width="100" />
              <el-table-column prop="message" label="错误信息" width="200" />
              <el-table-column prop="description" label="说明" />
            </el-table>
          </el-tab-pane>
        </el-tabs>
        
        <div class="api-docs">
          <el-button type="primary" @click="viewFullDocs">
            <el-icon class="el-icon--left"><Document /></el-icon>查看完整API文档
          </el-button>
        </div>
      </div>
    </el-card>

    <el-card class="api-card">
      <template #header>
        <div class="card-header">
          <span>API调用统计</span>
          <el-select v-model="timeRange" placeholder="选择时间范围" size="small">
            <el-option label="今日" value="today" />
            <el-option label="本周" value="week" />
            <el-option label="本月" value="month" />
          </el-select>
        </div>
      </template>
      
      <div class="api-stats">
        <div class="stat-item">
          <div class="stat-value">{{ apiStats.totalCalls }}</div>
          <div class="stat-label">总调用次数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ apiStats.successRate }}%</div>
          <div class="stat-label">成功率</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ apiStats.avgResponseTime }}ms</div>
          <div class="stat-label">平均响应时间</div>
        </div>
      </div>
      
      <div class="usage-chart" ref="chartRef">
        <!-- 图表将在这里渲染 -->
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { CopyDocument, View, Hide, Document } from '@element-plus/icons-vue'
import * as echarts from 'echarts'

const router = useRouter()

// API信息
const apiInfo = reactive({
  appId: 'app_' + Math.random().toString(36).substring(2, 10),
  apiKey: Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
})

// 控制API密钥显示
const showApiKey = ref(false)

// 时间范围选择
const timeRange = ref('today')

// API统计数据
const apiStats = reactive({
  totalCalls: 1250,
  successRate: 99.5,
  avgResponseTime: 145
})

// 错误码表
const errorCodes = [
  {
    code: '400',
    message: '请求参数错误',
    description: '请检查请求参数格式是否正确'
  },
  {
    code: '401',
    message: '认证失败',
    description: 'API密钥无效或已过期'
  },
  {
    code: '403',
    message: '权限不足',
    description: '当前账户无权访问该API'
  },
  {
    code: '429',
    message: '请求过于频繁',
    description: '已超过API调用频率限制'
  },
  {
    code: '500',
    message: '服务器错误',
    description: '服务器内部错误，请稍后重试'
  }
]

// 切换API密钥可见性
const toggleApiKeyVisibility = () => {
  showApiKey.value = !showApiKey.value
}

// 复制到剪贴板
const copyToClipboard = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 重新生成密钥
const regenerateKeys = () => {
  ElMessageBox.confirm('重新生成密钥后，原有密钥将立即失效。确定要继续吗？', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 生成新的API密钥
    apiInfo.apiKey = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
    ElMessage.success('API密钥已重新生成')
  }).catch(() => {})
}

// 查看完整API文档
const viewFullDocs = () => {
  router.push('/client/api-documentation')
}

// 初始化图表
const chartRef = ref(null)
const initChart = () => {
  if (!chartRef.value) return
  
  const chart = echarts.init(chartRef.value)
  const option = {
    tooltip: {
      trigger: 'axis'
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['周一', '周二', '周三', '周四', '周五', '周六', '周日']
    },
    yAxis: {
      type: 'value',
      name: 'API调用次数'
    },
    series: [{
      name: 'API调用量',
      type: 'line',
      smooth: true,
      data: [120, 132, 101, 134, 90, 230, 210]
    }]
  }
  
  chart.setOption(option)
  
  // 监听窗口大小变化，调整图表大小
  window.addEventListener('resize', () => {
    chart.resize()
  })
}

onMounted(() => {
  initChart()
})
</script>

<style scoped>
.api-settings-container {
  padding: 20px;
}

.api-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.api-info {
  margin-bottom: 20px;
}

.info-item {
  margin-bottom: 15px;
}

.info-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 5px;
}

.info-value-container {
  display: flex;
  align-items: center;
  gap: 10px;
}

.info-value {
  font-family: monospace;
  background-color: #f5f7fa;
  padding: 8px 12px;
  border-radius: 4px;
  flex-grow: 1;
}

.code-block {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
  margin: 10px 0;
}

.code-block pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.api-docs {
  margin-top: 20px;
  text-align: center;
}

.api-stats {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  margin-top: 5px;
}

.usage-chart {
  height: 300px;
  margin-top: 20px;
}

.limit-list {
  padding-left: 20px;
}

.limit-list li {
  margin-bottom: 10px;
  color: #606266;
}
</style>