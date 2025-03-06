<template>
  <div class="admins-container">
    <h2>管理员账号管理</h2>
    <el-card class="box-card">
      <div class="table-operations">
        <el-button type="primary" @click="handleAdd">添加管理员</el-button>
      </div>
      <el-table :data="adminList" style="width: 100%">
        <el-table-column label="头像" width="80">
          <template #default="scope">
            <el-avatar :size="40" :src="scope.row.avatar || '/src/assets/images/avatars/default-avatar.svg'" />
          </template>
        </el-table-column>
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="createTime" label="创建时间" />
        <el-table-column prop="lastLogin" label="最后登录时间" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 添加/编辑管理员对话框 -->
    <el-dialog
      :title="dialogType === 'add' ? '添加管理员' : '编辑管理员'"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="adminForm" label-width="80px">
        <el-form-item label="头像">
          <el-upload
            class="avatar-uploader"
            action="/api/upload"
            :show-file-list="false"
            :before-upload="beforeAvatarUpload"
            :on-success="handleAvatarSuccess"
          >
            <img v-if="adminForm.avatar" :src="adminForm.avatar" class="avatar-preview" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="adminForm.username" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="adminForm.role" style="width: 100%">
            <el-option label="超级管理员" value="super_admin" />
            <el-option label="运维管理员" value="ops_admin" />
            <el-option label="客服管理员" value="service_admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const dialogVisible = ref(false)
const dialogType = ref('add')
const adminFormRef = ref(null)

const adminForm = reactive({
  username: '',
  role: '',
  email: '',
  password: '',
  avatar: ''
})

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('头像只能是 JPG 或 PNG 格式!')
  }
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB!')
  }
  return isJPG && isLt2M
}

const handleAvatarSuccess = (response, uploadFile) => {
  // 实际项目中这里应该使用后端返回的URL
  adminForm.avatar = URL.createObjectURL(uploadFile.raw)
}

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const roles = [
  { label: '超级管理员', value: 'super_admin' },
  { label: '管理员', value: 'admin' },
  { label: '运维人员', value: 'operator' }
]

// 模拟管理员列表数据
const adminList = ref([
  {
    id: 1,
    username: 'admin',
    role: '超级管理员',
    email: 'admin@example.com',
    lastLogin: '2024-01-20 10:30:00',
    status: 'active'
  },
  {
    id: 2,
    username: 'operator',
    role: '运维人员',
    email: 'operator@example.com',
    lastLogin: '2024-01-20 09:15:00',
    status: 'active'
  }
])

const handleAdd = () => {
  dialogType.value = 'add'
  dialogVisible.value = true
  adminForm.username = ''
  adminForm.role = ''
  adminForm.email = ''
  adminForm.password = ''
}

const handleEdit = (row) => {
  dialogType.value = 'edit'
  dialogVisible.value = true
  adminForm.username = row.username
  adminForm.role = row.role
  adminForm.email = row.email
}

const handleStatusChange = (row) => {
  const action = row.status === 'active' ? '禁用' : '启用'
  ElMessageBox.confirm(`确定要${action}该管理员吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // 实际应用中这里应该调用API
    row.status = row.status === 'active' ? 'inactive' : 'active'
    ElMessage.success(`${action}成功`)
  }).catch(() => {})
}

const handleSubmit = async () => {
  if (!adminFormRef.value) return
  
  try {
    await adminFormRef.value.validate()
    // 实际应用中这里应该调用API
    ElMessage.success(dialogType.value === 'add' ? '添加成功' : '编辑成功')
    dialogVisible.value = false
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}
</script>

<style scoped>
.admins-container {
  padding: 20px;
}

.table-operations {
  margin-bottom: 20px;
}

.avatar-uploader {
  display: flex;
  justify-content: center;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.avatar-uploader-icon:hover {
  border-color: #409EFF;
}

.avatar-preview {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.3s;
}

.avatar-preview:hover {
  transform: scale(1.05);
}
</style>