<template>
  <div class="versions-container">
    <h2>模型版本管理</h2>
    <el-card class="box-card">
      <div class="table-operations">
        <el-button type="primary" @click="handleAdd">添加版本</el-button>
        <el-button type="info" @click="refreshData" :loading="modelStore.loading">刷新</el-button>
      </div>
      <el-table :data="modelStore.versions" style="width: 100%" v-loading="modelStore.loading">
        <el-table-column prop="version" label="版本号" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="releaseDate" label="发布日期" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="accuracy" label="准确率">
          <template #default="scope">
            <span>{{ (scope.row.accuracy * 100).toFixed(2) }}%</span>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button type="success" size="small" @click="handleDeploy(scope.row)">部署</el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 版本编辑对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="50%"
      :close-on-click-modal="false"
    >
      <el-form :model="versionForm" :rules="rules" ref="versionFormRef" label-width="100px">
        <el-form-item label="版本号" prop="version">
          <el-input v-model="versionForm.version" placeholder="请输入版本号" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="versionForm.description" type="textarea" placeholder="请输入版本描述" />
        </el-form-item>
        <el-form-item label="发布日期" prop="releaseDate">
          <el-date-picker
            v-model="versionForm.releaseDate"
            type="date"
            placeholder="选择发布日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="versionForm.status" placeholder="请选择状态">
            <el-option label="生产中" value="production" />
            <el-option label="测试中" value="testing" />
            <el-option label="开发中" value="development" />
          </el-select>
        </el-form-item>
        <el-form-item label="准确率" prop="accuracy">
          <el-input-number
            v-model="versionForm.accuracy"
            :precision="4"
            :step="0.01"
            :min="0"
            :max="1"
          />
        </el-form-item>
        <el-form-item label="参数配置">
          <el-collapse>
            <el-collapse-item title="训练参数" name="1">
              <el-form-item label="批处理大小" prop="parameters.batchSize">
                <el-input-number v-model="versionForm.parameters.batchSize" :min="1" :max="512" />
              </el-form-item>
              <el-form-item label="学习率" prop="parameters.learningRate">
                <el-input-number 
                  v-model="versionForm.parameters.learningRate" 
                  :precision="6" 
                  :step="0.0001" 
                  :min="0.000001" 
                  :max="1" 
                />
              </el-form-item>
              <el-form-item label="训练轮数" prop="parameters.epochs">
                <el-input-number v-model="versionForm.parameters.epochs" :min="1" :max="1000" />
              </el-form-item>
            </el-collapse-item>
          </el-collapse>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useModelStore } from '../../stores/model'

const modelStore = useModelStore()
const dialogVisible = ref(false)
const dialogTitle = ref('')
const versionFormRef = ref(null)
const submitLoading = ref(false)

// 表单数据
const versionForm = reactive({
  id: null,
  version: '',
  description: '',
  releaseDate: '',
  status: 'development',
  accuracy: 0.8,
  parameters: {
    batchSize: 32,
    learningRate: 0.001,
    epochs: 100
  }
})

// 表单验证规则
const rules = {
  version: [{ required: true, message: '请输入版本号', trigger: 'blur' }],
  description: [{ required: true, message: '请输入版本描述', trigger: 'blur' }],
  releaseDate: [{ required: true, message: '请选择发布日期', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  accuracy: [{ required: true, message: '请输入准确率', trigger: 'change' }]
}

onMounted(async () => {
  // 初始化版本列表数据
  await refreshData()
})

// 刷新数据
const refreshData = async () => {
  try {
    await modelStore.fetchVersions()
  } catch (error) {
    ElMessage.error('获取模型版本列表失败')
  }
}

// 添加版本
const handleAdd = () => {
  dialogTitle.value = '添加版本'
  Object.assign(versionForm, {
    id: null,
    version: '',
    description: '',
    releaseDate: new Date().toISOString().split('T')[0],
    status: 'development',
    accuracy: 0.8,
    parameters: {
      batchSize: 32,
      learningRate: 0.001,
      epochs: 100
    }
  })
  dialogVisible.value = true
}

// 编辑版本
const handleEdit = (row) => {
  dialogTitle.value = '编辑版本'
  Object.assign(versionForm, JSON.parse(JSON.stringify(row)))
  dialogVisible.value = true
}

// 部署版本
const handleDeploy = (row) => {
  ElMessageBox.confirm(
    `确定要部署版本 ${row.version} 到生产环境吗？`,
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      // 这里应该调用实际的部署API
      // 模拟部署逻辑
      const index = modelStore.versions.findIndex(item => item.id === row.id)
      if (index > -1) {
        // 将所有版本设置为非生产
        modelStore.versions.forEach(v => {
          if (v.status === 'production') {
            v.status = 'testing'
          }
        })
        // 将当前版本设置为生产
        modelStore.versions[index].status = 'production'
        // 设置为当前版本
        modelStore.setCurrentVersion(modelStore.versions[index])
        ElMessage.success('部署成功')
      }
    } catch (error) {
      ElMessage.error('部署失败')
    }
  }).catch(() => {
    ElMessage.info('已取消部署')
  })
}

// 删除版本
const handleDelete = (row) => {
  ElMessageBox.confirm(
    `确定要删除版本 ${row.version} 吗？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      // 这里应该调用实际的删除API
      // 模拟删除逻辑
      const index = modelStore.versions.findIndex(item => item.id === row.id)
      if (index > -1) {
        modelStore.versions.splice(index, 1)
        ElMessage.success('删除成功')
      }
    } catch (error) {
      ElMessage.error('删除失败')
    }
  }).catch(() => {
    ElMessage.info('已取消删除')
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!versionFormRef.value) return
  
  try {
    const valid = await versionFormRef.value.validate()
    if (!valid) return
    
    submitLoading.value = true
    
    if (versionForm.id) {
      // 编辑模式
      // 这里应该调用实际的更新API
      // 模拟更新逻辑
      const index = modelStore.versions.findIndex(item => item.id === versionForm.id)
      if (index > -1) {
        modelStore.versions[index] = { ...versionForm }
        ElMessage.success('更新成功')
      }
    } else {
      // 添加模式
      // 这里应该调用实际的添加API
      // 模拟添加逻辑
      const newVersion = {
        ...versionForm,
        id: modelStore.versions.length > 0 
          ? Math.max(...modelStore.versions.map(v => v.id)) + 1 
          : 1
      }
      modelStore.versions.push(newVersion)
      ElMessage.success('添加成功')
    }
    
    dialogVisible.value = false
  } catch (error) {
    ElMessage.error('提交失败')
  } finally {
    submitLoading.value = false
  }
}

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    'production': 'success',
    'testing': 'warning',
    'development': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    'production': '生产中',
    'testing': '测试中',
    'development': '开发中'
  }
  return statusMap[status] || status
}
</script>

<style scoped>
.versions-container {
  padding: 20px;
}

.table-operations {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>