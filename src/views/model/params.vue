<template>
  <div class="params-container">
    <h2>模型参数配置</h2>
    <el-card class="box-card">
      <div class="table-operations">
        <el-button type="primary" @click="handleAdd">添加参数</el-button>
      </div>
      <el-table :data="paramsList" style="width: 100%">
        <el-table-column prop="name" label="参数名称" />
        <el-table-column prop="value" label="参数值" />
        <el-table-column prop="type" label="参数类型" />
        <el-table-column prop="description" label="描述" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 参数编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="50%"
      :close-on-click-modal="false"
    >
      <el-form :model="paramForm" :rules="rules" ref="paramFormRef" label-width="100px">
        <el-form-item label="参数名称" prop="name">
          <el-input v-model="paramForm.name" placeholder="请输入参数名称" />
        </el-form-item>
        <el-form-item label="参数值" prop="value">
          <el-input v-model="paramForm.value" placeholder="请输入参数值" />
        </el-form-item>
        <el-form-item label="参数类型" prop="type">
          <el-select v-model="paramForm.type" placeholder="请选择参数类型">
            <el-option label="数值" value="number" />
            <el-option label="文本" value="string" />
            <el-option label="布尔值" value="boolean" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="paramForm.description" type="textarea" placeholder="请输入参数描述" />
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
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { modelParams } from '@/api/model'

const paramsList = ref([])
const dialogVisible = ref(false)
const dialogTitle = ref('')
const paramFormRef = ref(null)

// 表单数据
const paramForm = ref({
  id: null,
  name: '',
  value: '',
  type: 'string',
  description: ''
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入参数名称', trigger: 'blur' }],
  value: [{ required: true, message: '请输入参数值', trigger: 'blur' }],
  type: [{ required: true, message: '请选择参数类型', trigger: 'change' }],
  description: [{ required: true, message: '请输入参数描述', trigger: 'blur' }]
}

onMounted(() => {
  // 初始化参数列表数据
  paramsList.value = modelParams
})

// 添加参数
const handleAdd = () => {
  dialogTitle.value = '添加参数'
  paramForm.value = {
    id: null,
    name: '',
    value: '',
    type: 'string',
    description: ''
  }
  dialogVisible.value = true
}

// 编辑参数
const handleEdit = (row) => {
  dialogTitle.value = '编辑参数'
  paramForm.value = { ...row }
  dialogVisible.value = true
}

// 删除参数
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除参数 ${row.name} 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    const index = paramsList.value.findIndex(item => item.id === row.id)
    if (index > -1) {
      paramsList.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 提交表单
const handleSubmit = () => {
  paramFormRef.value?.validate((valid) => {
    if (valid) {
      if (paramForm.value.id) {
        // 编辑模式
        const index = paramsList.value.findIndex(item => item.id === paramForm.value.id)
        if (index > -1) {
          paramsList.value[index] = { ...paramForm.value }
          ElMessage.success('更新成功')
        }
      } else {
        // 添加模式
        const newParam = {
          ...paramForm.value,
          id: paramsList.value.length + 1
        }
        paramsList.value.push(newParam)
        ElMessage.success('添加成功')
      }
      dialogVisible.value = false
    }
  })
}
</script>

<style scoped>
.params-container {
  padding: 20px;
}

.table-operations {
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>