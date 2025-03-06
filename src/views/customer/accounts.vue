<template>
  <div class="accounts-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>客户账户管理</span>
              <el-button type="primary" size="small" @click="handleAdd">添加客户</el-button>
            </div>
          </template>
          <el-table :data="customerList" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="name" label="客户名称" />
            <el-table-column prop="package" label="套餐类型" />
            <el-table-column prop="address" label="地址" />
            <el-table-column prop="contact" label="联系人" />
            <el-table-column prop="phone" label="联系电话" />
            <el-table-column prop="status" label="状态">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'danger'">
                  {{ row.status === 'active' ? '正常' : '停用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="expireDate" label="到期时间" />
            <el-table-column label="使用情况">
              <template #default="{ row }">
                <div>API调用：{{ row.usage.api.toLocaleString() }}次</div>
                <div>存储空间：{{ row.usage.storage }}GB</div>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button-group>
                  <el-button size="small" type="primary" @click="handleEdit(row)">编辑</el-button>
                  <el-button size="small" type="warning" @click="handleRenew(row)">续费</el-button>
                  <el-button 
                    size="small" 
                    :type="row.status === 'active' ? 'danger' : 'success'"
                    @click="handleStatusChange(row)"
                  >
                    {{ row.status === 'active' ? '停用' : '启用' }}
                  </el-button>
                </el-button-group>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑客户对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="客户名称">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="套餐类型">
          <el-select v-model="formData.package" style="width: 100%">
            <el-option label="企业版" value="企业版" />
            <el-option label="专业版" value="专业版" />
            <el-option label="基础版" value="基础版" />
          </el-select>
        </el-form-item>
        <el-form-item label="登录账号">
          <el-input v-model="formData.username" placeholder="请设置登录账号" />
        </el-form-item>
        <el-form-item label="登录密码">
          <el-input v-model="formData.password" type="password" placeholder="请设置登录密码" />
        </el-form-item>
        <el-form-item label="确认密码">
          <el-input v-model="formData.confirmPassword" type="password" placeholder="请再次输入密码" />
        </el-form-item>
        <el-form-item label="详细地址">
          <el-input v-model="formData.address" placeholder="请输入详细地址" />
        </el-form-item>
        <el-form-item label="联系人">
          <el-input v-model="formData.contact" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话">
          <el-input v-model="formData.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="到期时间">
          <el-date-picker
            v-model="formData.expireDate"
            type="date"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 续费对话框 -->
    <el-dialog
      title="续费管理"
      v-model="renewDialogVisible"
      width="400px"
    >
      <el-form :model="renewForm" label-width="100px">
        <el-form-item label="续费时长">
          <el-select v-model="renewForm.duration" style="width: 100%">
            <el-option label="1个月" value="1" />
            <el-option label="3个月" value="3" />
            <el-option label="6个月" value="6" />
            <el-option label="12个月" value="12" />
          </el-select>
        </el-form-item>
        <el-form-item label="续费金额">
          <el-input v-model="renewForm.amount" type="number" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="renewDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleRenewSubmit">确认续费</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { customers } from '@/api/mock'

// 客户列表数据
const customerList = ref(customers)

// 对话框控制
const dialogVisible = ref(false)
const renewDialogVisible = ref(false)
const dialogTitle = ref('')

// 表单数据
const formData = ref({
  id: null,
  name: '',
  package: '',
  username: '',
  password: '',
  confirmPassword: '',
  address: '',
  contact: '',
  phone: '',
  expireDate: ''
})

// 表单验证规则
const formRules = {
  name: [{ required: true, message: '请输入客户名称', trigger: 'blur' }],
  package: [{ required: true, message: '请选择套餐类型', trigger: 'change' }],
  username: [{ required: true, message: '请设置登录账号', trigger: 'blur' }],
  password: [
    { required: true, message: '请设置登录密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== formData.value.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  contact: [{ required: true, message: '请输入联系人姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  expireDate: [{ required: true, message: '请选择到期时间', trigger: 'change' }]
}

// 表单引用
const formRef = ref(null)

// 处理表单提交
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()

    if (!formData.value.name || !formData.value.package || !formData.value.expireDate ||
        !formData.value.username || !formData.value.password) {
      ElMessage.warning('请填写完整信息')
      return
    }

    if (formData.value.password !== formData.value.confirmPassword) {
      ElMessage.warning('两次输入的密码不一致')
      return
    }

    if (formData.value.id) {
      // 编辑现有客户
      const index = customerList.value.findIndex(item => item.id === formData.value.id)
      if (index !== -1) {
        const { confirmPassword, ...submitData } = formData.value
        customerList.value[index] = { ...customerList.value[index], ...submitData }
        ElMessage.success('客户信息已更新')
      }
    } else {
      // 添加新客户
      const { confirmPassword, ...submitData } = formData.value
      const newCustomer = {
        ...submitData,
        id: customerList.value.length + 1,
        status: 'active',
        usage: { api: 0, storage: 0 }
      }
      customerList.value.push(newCustomer)
      ElMessage.success('客户添加成功')
    }
    
    dialogVisible.value = false
    formRef.value.resetFields()
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

// 续费表单数据
const renewForm = ref({
  customerId: null,
  duration: '1',
  amount: 0
})

// 处理续费操作
const handleRenew = (customer) => {
  renewForm.value.customerId = customer.id
  renewForm.value.duration = '1'
  renewForm.value.amount = calculateRenewAmount(1, customer.package)
  renewDialogVisible.value = true
}

// 计算续费金额
const calculateRenewAmount = (duration, packageType) => {
  const basePrice = {
    '企业版': 9999,
    '专业版': 4999,
    '基础版': 1999
  }
  return (basePrice[packageType] || 0) * (duration / 12)
}

// 提交续费
const handleRenewSubmit = () => {
  const customer = customerList.value.find(item => item.id === renewForm.value.customerId)
  if (customer) {
    const months = parseInt(renewForm.value.duration)
    const currentExpireDate = new Date(customer.expireDate)
    const newExpireDate = new Date(currentExpireDate.setMonth(currentExpireDate.getMonth() + months))
    customer.expireDate = newExpireDate.toISOString().split('T')[0]
    ElMessage.success(`续费成功，已延长${months}个月`)
  }
  renewDialogVisible.value = false
}

// 处理添加客户
const handleAdd = () => {
  dialogTitle.value = '添加客户'
  formData.value = {
    id: null,
    name: '',
    package: '',
    username: '',
    password: '',
    confirmPassword: '',
    address: '',
    contact: '',
    phone: '',
    expireDate: ''
  }
  dialogVisible.value = true
}

// 处理编辑客户
const handleEdit = (row) => {
  dialogTitle.value = '编辑客户'
  formData.value = { ...row, confirmPassword: row.password }
  dialogVisible.value = true
}

// 处理状态变更
const handleStatusChange = (row) => {
  row.status = row.status === 'active' ? 'inactive' : 'active'
  ElMessage.success(`客户状态已${row.status === 'active' ? '启用' : '停用'}`)
}

// 监听续费时长变化
const watchRenewDuration = computed(() => {
  const customer = customerList.value.find(item => item.id === renewForm.value.customerId)
  if (customer) {
    const duration = parseInt(renewForm.value.duration)
    renewForm.value.amount = calculateRenewAmount(duration, customer.package)
  }
  return renewForm.value.duration
})
</script>

<style scoped>
.accounts-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.dialog-footer {
  text-align: right;
}

.el-button + .el-button {
  margin-left: 10px;
}
</style>