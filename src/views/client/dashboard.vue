<template>
  <div class="client-dashboard">
    <!-- 账户信息卡片 -->
    <el-card class="dashboard-card">
      <template #header>
        <div class="card-header">
          <span>账户信息</span>
          <div class="header-actions">
            <el-button type="primary" link @click="showEditProfileDialog">修改信息</el-button>
            <el-button type="primary" link @click="showChangePasswordDialog">修改密码</el-button>
          </div>
        </div>
      </template>
      <div class="account-info">
        <el-descriptions :column="2">
          <el-descriptions-item label="公司名称">{{ customerInfo.name }}</el-descriptions-item>
          <el-descriptions-item label="联系人">{{ customerInfo.contact }}</el-descriptions-item>
          <el-descriptions-item label="联系电话">{{ customerInfo.phone }}</el-descriptions-item>
          <el-descriptions-item label="地址">{{ customerInfo.address }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-card>

    <!-- 套餐信息卡片 -->
    <el-card class="dashboard-card">
      <template #header>
        <div class="card-header">
          <span>套餐信息</span>
          <el-tag :type="customerInfo.status === 'active' ? 'success' : 'danger'">
            {{ customerInfo.status === 'active' ? '使用中' : '已过期' }}
          </el-tag>
        </div>
      </template>
      <div class="package-info">
        <el-descriptions :column="2">
          <el-descriptions-item label="当前套餐">{{ customerInfo.package }}</el-descriptions-item>
          <el-descriptions-item label="到期时间">{{ customerInfo.expireDate }}</el-descriptions-item>
        </el-descriptions>
        <div class="usage-stats">
          <div class="usage-item">
            <div class="usage-title">API调用量</div>
            <el-progress :percentage="(customerInfo.usage.api / packageInfo.limits.apiCalls) * 100" :format="format" />
            <div class="usage-detail">{{ customerInfo.usage.api }} / {{ packageInfo.limits.apiCalls }}次</div>
          </div>
          <div class="usage-item">
            <div class="usage-title">存储空间</div>
            <el-progress :percentage="(customerInfo.usage.storage / packageInfo.limits.storage) * 100" :format="format" />
            <div class="usage-detail">{{ customerInfo.usage.storage }} / {{ packageInfo.limits.storage }}GB</div>
          </div>
        </div>
        <div class="package-features">
          <h4>套餐特权：</h4>
          <ul>
            <li v-for="(feature, index) in packageInfo.features" :key="index">{{ feature }}</li>
          </ul>
        </div>
      </div>
    </el-card>

    <!-- API调用统计图表 -->
    <el-card class="dashboard-card full-width">
      <template #header>
        <div class="card-header">
          <span>API调用统计</span>
          <div class="chart-controls">
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              size="small"
              @change="updateChart"
            />
            <el-radio-group v-model="apiChartPeriod" size="small" @change="updateChart">
              <el-radio-button label="week">周</el-radio-button>
              <el-radio-button label="month">月</el-radio-button>
            </el-radio-group>
          </div>
        </div>
      </template>
      <div class="chart-container" ref="apiChartRef"></div>
    </el-card>
  </div>

  <!-- 修改信息对话框 -->
  <el-dialog v-model="editProfileVisible" title="修改账户信息" width="500px">
    <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
      <el-form-item label="联系人" prop="contact">
        <el-input v-model="editForm.contact" />
      </el-form-item>
      <el-form-item label="联系电话" prop="phone">
        <el-input v-model="editForm.phone" />
      </el-form-item>
      <el-form-item label="地址" prop="address">
        <el-input v-model="editForm.address" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="editProfileVisible = false">取消</el-button>
      <el-button type="primary" @click="handleUpdateProfile">确认</el-button>
    </template>
  </el-dialog>

  <!-- 修改密码对话框 -->
  <el-dialog v-model="changePasswordVisible" title="修改密码" width="500px">
    <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
      <el-form-item label="原密码" prop="oldPassword">
        <el-input v-model="passwordForm.oldPassword" type="password" show-password />
      </el-form-item>
      <el-form-item label="新密码" prop="newPassword">
        <el-input v-model="passwordForm.newPassword" type="password" show-password />
      </el-form-item>
      <el-form-item label="确认密码" prop="confirmPassword">
        <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="changePasswordVisible = false">取消</el-button>
      <el-button type="primary" @click="handleChangePassword">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, onMounted, reactive, onBeforeUnmount } from 'vue'
import { customers } from '@/api/mock'
import { packages } from '@/api/customer'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'

// 客户信息
const customerInfo = ref(customers[0])
const packageInfo = ref(packages.find(p => p.name === customerInfo.value.package))

// API图表相关
const apiChartRef = ref(null)
const apiChartPeriod = ref('week')
const dateRange = ref([])
let apiChart = null

// 修改信息相关
const editProfileVisible = ref(false)
const editFormRef = ref(null)
const editForm = reactive({
  contact: '',
  phone: '',
  address: ''
})

const editRules = {
  contact: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  address: [{ required: true, message: '请输入地址', trigger: 'blur' }]
}

// 修改密码相关
const changePasswordVisible = ref(false)
const passwordFormRef = ref(null)
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 8, message: '密码长度不能小于8位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 显示修改信息对话框
const showEditProfileDialog = () => {
  editForm.contact = customerInfo.value.contact
  editForm.phone = customerInfo.value.phone
  editForm.address = customerInfo.value.address
  editProfileVisible.value = true
}

// 显示修改密码对话框
const showChangePasswordDialog = () => {
  changePasswordVisible.value = true
}

// 处理更新信息
const handleUpdateProfile = async () => {
  if (!editFormRef.value) return
  
  try {
    await editFormRef.value.validate()
    // TODO: 实现更新信息的API调用
    editProfileVisible.value = false
    ElMessage.success('信息更新成功')
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 处理修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    await passwordFormRef.value.validate()
    // TODO: 实现修改密码的API调用
    changePasswordVisible.value = false
    ElMessage.success('密码修改成功')
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 格式化进度条显示
const format = (percentage) => {
  return percentage.toFixed(1) + '%'
}

// 初始化API调用图表
const initApiChart = () => {
  if (!apiChartRef.value) return
  
  apiChart = echarts.init(apiChartRef.value)
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
      data: [820, 932, 901, 934, 1290, 1330, 1320],
      type: 'line',
      smooth: true,
      areaStyle: {
        opacity: 0.3
      }
    }]
  }
  apiChart.setOption(option)
}

// 更新图表数据
const updateChart = () => {
  // TODO: 根据选择的时间范围和周期获取数据并更新图表
}

onMounted(() => {
  initApiChart()
  window.addEventListener('resize', () => {
    apiChart?.resize()
  })
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', () => {
    apiChart?.resize()
  })
  apiChart?.dispose()
})
</script>

<style scoped>
.client-dashboard {
  padding: 20px;
  display: grid;
  gap: 20px;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
}

.dashboard-card {
  height: fit-content;
}

.full-width {
  grid-column: 1 / -1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.chart-controls {
  display: flex;
  gap: 10px;
  align-items: center;
}

.usage-stats {
  margin: 20px 0;
}

.usage-item {
  margin-bottom: 15px;
}

.usage-title {
  margin-bottom: 5px;
  font-weight: 500;
}

.usage-detail {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
}

.package-features {
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
  padding-top: 15px;
}

.package-features h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.package-features ul {
  margin: 0;
  padding-left: 20px;
  color: #606266;
}

.package-features li {
  margin-bottom: 5px;
}

.chart-container {
  height: 400px;
}
</style>