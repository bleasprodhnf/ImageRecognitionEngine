<template>
  <div class="profile-container">
    <el-tabs v-model="activeTab">
      <!-- 基本信息表单 -->
      <el-tab-pane label="基本信息" name="basic">
        <el-form
          ref="profileFormRef"
          :model="profileForm"
          :rules="profileRules"
          label-width="100px"
        >
          <el-form-item label="用户名" prop="username">
            <el-input v-model="profileForm.username" disabled />
          </el-form-item>
          <el-form-item label="真实姓名" prop="realname">
            <el-input v-model="profileForm.realname" />
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="profileForm.email" />
          </el-form-item>
          <el-form-item label="手机号" prop="phone">
            <el-input v-model="profileForm.phone" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleUpdateProfile" :loading="updating">
              保存修改
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <!-- 修改密码表单 -->
      <el-tab-pane label="修改密码" name="password">
        <el-form
          ref="passwordFormRef"
          :model="passwordForm"
          :rules="passwordRules"
          label-width="100px"
        >
          <el-form-item label="原密码" prop="oldPassword">
            <el-input
              v-model="passwordForm.oldPassword"
              type="password"
              show-password
            />
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input
              v-model="passwordForm.newPassword"
              type="password"
              show-password
            >
              <template #append>
                <el-popover
                  placement="right"
                  :width="200"
                  trigger="hover"
                >
                  <template #reference>
                    <el-icon><InfoFilled /></el-icon>
                  </template>
                  <div class="password-rules">
                    <h4>密码要求：</h4>
                    <ul>
                      <li :class="{ valid: passwordStrength.length }">长度8-20个字符</li>
                      <li :class="{ valid: passwordStrength.upper }">包含大写字母</li>
                      <li :class="{ valid: passwordStrength.lower }">包含小写字母</li>
                      <li :class="{ valid: passwordStrength.number }">包含数字</li>
                      <li :class="{ valid: passwordStrength.special }">包含特殊字符</li>
                    </ul>
                  </div>
                </el-popover>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input
              v-model="passwordForm.confirmPassword"
              type="password"
              show-password
            />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="handleChangePassword" :loading="changing">
              修改密码
            </el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { InfoFilled } from '@element-plus/icons-vue'
import { userInfo, passwordRules as pwdRules, updateProfileResponse, changePasswordResponse } from '@/api/mock'

import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// 当前激活的标签页
const activeTab = ref('basic')

// 根据路由参数设置当前标签页
onMounted(() => {
  const tab = route.query.tab
  if (tab === 'password') {
    activeTab.value = 'password'
  }
})

// 表单引用
const profileFormRef = ref(null)
const passwordFormRef = ref(null)

// 加载状态
const updating = ref(false)
const changing = ref(false)

// 个人信息表单
const profileForm = reactive({
  username: userInfo.username,
  realname: userInfo.realname,
  email: userInfo.email,
  phone: userInfo.phone
})

// 个人信息验证规则
const profileRules = {
  realname: [{ required: true, message: '请输入真实姓名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ]
}

// 密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码强度检查
const passwordStrength = computed(() => {
  const pwd = passwordForm.newPassword
  return {
    length: pwd.length >= pwdRules.minLength && pwd.length <= pwdRules.maxLength,
    upper: /[A-Z]/.test(pwd),
    lower: /[a-z]/.test(pwd),
    number: /\d/.test(pwd),
    special: new RegExp(`[${pwdRules.specialChars}]`).test(pwd)
  }
})

// 密码验证规则
const passwordRules = {
  oldPassword: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        const strength = passwordStrength.value
        if (!strength.length) {
          callback(new Error(`密码长度必须在${pwdRules.minLength}-${pwdRules.maxLength}个字符之间`))
        } else if (!strength.upper) {
          callback(new Error('密码必须包含大写字母'))
        } else if (!strength.lower) {
          callback(new Error('密码必须包含小写字母'))
        } else if (!strength.number) {
          callback(new Error('密码必须包含数字'))
        } else if (!strength.special) {
          callback(new Error('密码必须包含特殊字符'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
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

// 更新个人信息
const handleUpdateProfile = async () => {
  if (!profileFormRef.value) return

  try {
    updating.value = true
    await profileFormRef.value.validate()
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success(updateProfileResponse.message)
  } catch (error) {
    console.error('更新个人信息失败:', error)
    ElMessage.error('更新个人信息失败，请重试')
  } finally {
    updating.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordFormRef.value) return

  try {
    changing.value = true
    await passwordFormRef.value.validate()
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    ElMessage.success(changePasswordResponse.message)
    // 清空表单
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    passwordFormRef.value.resetFields()
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error('修改密码失败，请重试')
  } finally {
    changing.value = false
  }
}
</script>

<style scoped>
.profile-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
}

.password-rules {
  font-size: 14px;
}

.password-rules h4 {
  margin: 0 0 10px;
  color: #606266;
}

.password-rules ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.password-rules li {
  margin: 5px 0;
  color: #909399;
  display: flex;
  align-items: center;
}

.password-rules li::before {
  content: '×';
  color: #f56c6c;
  margin-right: 5px;
}

.password-rules li.valid {
  color: #67c23a;
}

.password-rules li.valid::before {
  content: '✓';
  color: #67c23a;
}
</style>