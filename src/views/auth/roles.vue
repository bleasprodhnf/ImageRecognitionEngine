<template>
  <div class="roles-container">
    <h2>角色管理</h2>
    <el-card class="box-card">
      <div class="table-operations">
        <el-button type="primary" @click="handleAdd">添加角色</el-button>
      </div>
      <el-table :data="roles" style="width: 100%">
        <el-table-column prop="name" label="角色名称" />
        <el-table-column prop="code" label="角色代码" />
        <el-table-column label="权限数量">
          <template #default="{ row }">
            <el-tag type="info">{{ row.permissions.length }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="success" size="small" @click="handlePermissions(scope.row)">权限设置</el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 角色编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="roleForm" :rules="rules" ref="roleFormRef" label-width="100px">
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="roleForm.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="角色代码" prop="code">
          <el-input v-model="roleForm.code" placeholder="请输入角色代码" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitRole">确定</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 权限设置对话框 -->
    <el-dialog
      title="权限设置"
      v-model="permissionDialogVisible"
      width="600px"
    >
      <el-tree
        ref="permissionTreeRef"
        :data="permissionTree"
        show-checkbox
        node-key="code"
        :default-checked-keys="selectedPermissions"
        :props="{
          label: 'name',
          children: 'children'
        }"
      />
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="permissionDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPermissions">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { roles } from '@/api/mock'

// 角色表单
const dialogVisible = ref(false)
const dialogTitle = ref('')
const roleForm = ref({
  name: '',
  code: ''
})
const roleFormRef = ref(null)
const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入角色代码', trigger: 'blur' }]
}

// 权限树
const permissionDialogVisible = ref(false)
const permissionTreeRef = ref(null)
const selectedPermissions = ref([])
const currentRole = ref(null)

// 模拟权限树数据
const permissionTree = [
  {
    name: '系统管理',
    code: 'system',
    children: [
      { name: '用户管理', code: 'system:user' },
      { name: '角色管理', code: 'system:role' },
      { name: '权限管理', code: 'system:permission' },
      { name: '日志管理', code: 'system:log' }
    ]
  },
  {
    name: '客户管理',
    code: 'customer',
    children: [
      { name: '客户列表', code: 'customer:list' },
      { name: '客户服务', code: 'customer:service' }
    ]
  }
]

// 添加角色
const handleAdd = () => {
  dialogTitle.value = '添加角色'
  roleForm.value = {
    name: '',
    code: ''
  }
  dialogVisible.value = true
}

// 编辑角色
const handleEdit = (row) => {
  dialogTitle.value = '编辑角色'
  roleForm.value = { ...row }
  dialogVisible.value = true
}

// 提交角色
const submitRole = async () => {
  if (!roleFormRef.value) return
  await roleFormRef.value.validate((valid) => {
    if (valid) {
      // 这里应该调用API保存角色
      ElMessage.success('保存成功')
      dialogVisible.value = false
    }
  })
}

// 设置权限
const handlePermissions = (row) => {
  currentRole.value = row
  selectedPermissions.value = row.permissions
  permissionDialogVisible.value = true
}

// 提交权限设置
const submitPermissions = () => {
  if (!permissionTreeRef.value || !currentRole.value) return
  const checkedKeys = permissionTreeRef.value.getCheckedKeys()
  const halfCheckedKeys = permissionTreeRef.value.getHalfCheckedKeys()
  const allCheckedKeys = [...checkedKeys, ...halfCheckedKeys]
  
  // 这里应该调用API保存权限设置
  currentRole.value.permissions = allCheckedKeys
  ElMessage.success('权限设置成功')
  permissionDialogVisible.value = false
}

// 删除角色
const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确定要删除该角色吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    // 这里应该调用API删除角色
    ElMessage.success('删除成功')
  })
}
</script>

<style scoped>
.roles-container {
  padding: 20px;
}

.table-operations {
  margin-bottom: 20px;
}
</style>