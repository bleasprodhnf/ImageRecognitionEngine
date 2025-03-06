<template>
  <div class="client-profile">
    <el-row :gutter="20">
      <!-- 基本信息 -->
      <el-col :span="12">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button type="primary" size="small" @click="handleUpdateProfile">保存修改</el-button>
            </div>
          </template>
          <el-form :model="profileForm" label-width="100px">
            <el-form-item label="企业名称">
              <el-input v-model="profileForm.companyName" />
            </el-form-item>
            <el-form-item label="联系人">
              <el-input v-model="profileForm.contactPerson" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="profileForm.phone" />
            </el-form-item>
            <el-form-item label="电子邮箱">
              <el-input v-model="profileForm.email" />
            </el-form-item>
            <el-form-item label="企业地址">
              <el-input v-model="profileForm.address" type="textarea" :rows="2" />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 密码修改 -->
      <el-col :span="12">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>修改密码</span>
              <el-button type="primary" size="small" @click="handleUpdatePassword">确认修改</el-button>
            </div>
          </template>
          <el-form :model="passwordForm" label-width="100px">
            <el-form-item label="当前密码">
              <el-input v-model="passwordForm.currentPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="passwordForm.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item label="确认新密码">
              <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

// 基本信息表单
const profileForm = ref({
  companyName: '示例科技有限公司',
  contactPerson: '张三',
  phone: '13800138000',
  email: 'example@company.com',
  address: '北京市朝阳区xxx街道xxx号'
})

// 密码修改表单
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 处理基本信息更新
const handleUpdateProfile = async () => {
  try {
    // TODO: 调用API更新基本信息
    ElMessage.success('基本信息更新成功')
  } catch (error) {
    ElMessage.error('基本信息更新失败')
  }
}

// 处理密码修改
const handleUpdatePassword = async () => {
  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }

  try {
    // TODO: 调用API修改密码
    ElMessage.success('密码修改成功')
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
  } catch (error) {
    ElMessage.error('密码修改失败')
  }
}
</script>

<style scoped>
.client-profile {
  padding: 20px;
}

.profile-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-form {
  max-width: 100%;
  margin-top: 20px;
}

.el-form-item:last-child {
  margin-bottom: 0;
}
</style>