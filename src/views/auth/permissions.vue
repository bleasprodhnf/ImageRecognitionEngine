<template>
  <div class="permissions-container">
    <h2>权限管理</h2>
    <el-card class="box-card">
      <div class="table-operations">
        <el-button type="primary" @click="handleAdd">添加权限</el-button>
      </div>
      <el-table :data="permissionList" style="width: 100%">
        <el-table-column prop="name" label="权限名称" />
        <el-table-column prop="code" label="权限代码" />
        <el-table-column prop="type" label="权限类型">
          <template #default="{ row }">
            <el-tag :type="row.type === 'menu' ? 'success' : 'info'">
              {{ row.type === 'menu' ? '菜单' : '操作' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 权限编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
    >
      <el-form :model="permissionForm" :rules="rules" ref="permissionFormRef" label-width="100px">
        <el-form-item label="权限名称" prop="name">
          <el-input v-model="permissionForm.name" placeholder="请输入权限名称" />
        </el-form-item>
        <el-form-item label="权限代码" prop="code">
          <el-input v-model="permissionForm.code" placeholder="请输入权限代码" />
        </el-form-item>
        <el-form-item label="权限类型" prop="type">
          <el-select v-model="permissionForm.type" placeholder="请选择权限类型">
            <el-option label="菜单权限" value="menu" />
            <el-option label="操作权限" value="operation" />
          </el-select>
        </el-form-item>
        <el-form-item label="权限描述" prop="description">
          <el-input
            v-model="permissionForm.description"
            type="textarea"
            placeholder="请输入权限描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitPermission">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 权限列表数据
const permissionList = ref([
  {
    name: '系统管理',
    code: 'system',
    type: 'menu',
    description: '系统管理模块'
  },
  {
    name: '用户管理',
    code: 'system:user',
    type: 'menu',
    description: '用户管理功能'
  },
  {
    name: '添加用户',
    code: 'system:user:add',
    type: 'operation',
    description: '添加用户操作'
  }
])

// 权限表单
const dialogVisible = ref(false)
const dialogTitle = ref('')
const permissionForm = ref({
  name: '',
  code: '',
  type: 'menu',
  description: ''
})
const permissionFormRef = ref(null)
const rules = {
  name: [{ required: true, message: '请输入权限名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入权限代码', trigger: 'blur' }],
  type: [{ required: true, message: '请选择权限类型', trigger: 'change' }]
}

// 添加权限
const handleAdd = () => {
  dialogTitle.value = '添加权限'
  permissionForm.value = {
    name: '',
    code: '',
    type: 'menu',
    description: ''
  }
  dialogVisible.value = true
}

// 编辑权限
const handleEdit = (row) => {
  dialogTitle.value = '编辑权限'
  permissionForm.value = { ...row }
  dialogVisible.value = true
}

// 提交权限
const submitPermission = async () => {
  if (!permissionFormRef.value) return
  await permissionFormRef.value.validate((valid) => {
    if (valid) {
      // 这里应该调用API保存权限
      ElMessage.success('保存成功')
      dialogVisible.value = false
    }
  })
}

// 删除权限
const handleDelete = (row) => {
  ElMessageBox.confirm(
    '确定要删除该权限吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    // 这里应该调用API删除权限
    ElMessage.success('删除成功')
  })
}
</script>

<style scoped>
.permissions-container {
  padding: 20px;
}

.table-operations {
  margin-bottom: 20px;
}
</style>