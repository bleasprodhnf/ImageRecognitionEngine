<template>
  <div class="packages-container">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>套餐配置管理</span>
              <el-button type="primary" size="small" @click="handleAdd">新增套餐</el-button>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="8" v-for="pkg in packageList" :key="pkg.id">
              <el-card class="package-card" :class="{ 'enterprise': pkg.name === '企业版' }">
                <div class="package-header">
                  <h2>{{ pkg.name }}</h2>
                  <div class="price">¥{{ pkg.price }}<span>/年</span></div>
                </div>
                <div class="package-content">
                  <p class="description">{{ pkg.description }}</p>
                  <div class="features">
                    <h4>功能特点</h4>
                    <ul>
                      <li v-for="feature in pkg.features" :key="feature">
                        <el-icon><Check /></el-icon>
                        {{ feature }}
                      </li>
                    </ul>
                  </div>
                  <div class="limits">
                    <h4>使用限制</h4>
                    <div class="limit-item">
                      <span>存储空间</span>
                      <span>{{ pkg.limits.storage }}GB</span>
                    </div>
                    <div class="limit-item">
                      <span>并发请求</span>
                      <span>{{ pkg.limits.concurrency }}次/秒</span>
                    </div>
                    <div class="limit-item">
                      <span>可用模型数</span>
                      <span>{{ pkg.limits.models }}个</span>
                    </div>
                  </div>
                  <div class="package-actions">
                    <el-button type="primary" size="small" @click="editPackage(pkg)">编辑</el-button>
                    <el-button type="danger" size="small" @click="deletePackage(pkg)">删除</el-button>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 编辑套餐对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogTitle"
      width="50%"
    >
      <el-form :model="editForm" label-width="120px">
        <el-form-item label="套餐名称">
          <el-input v-model="editForm.name" />
        </el-form-item>
        <el-form-item label="套餐价格">
          <el-input-number v-model="editForm.price" :min="0" />
        </el-form-item>
        <el-form-item label="套餐描述">
          <el-input v-model="editForm.description" type="textarea" />
        </el-form-item>
        <el-form-item label="功能特点">
          <el-tag
            v-for="(feature, index) in editForm.features"
            :key="index"
            closable
            @close="removeFeature(index)"
          >
            {{ feature }}
          </el-tag>
          <el-input
            v-if="inputVisible"
            ref="InputRef"
            v-model="inputValue"
            class="input-new-tag"
            size="small"
            @keyup.enter="addFeature"
            @blur="addFeature"
          />
          <el-button v-else size="small" @click="showInput">+ 添加功能</el-button>
        </el-form-item>
        <el-form-item label="存储空间(GB)">
          <el-input-number v-model="editForm.limits.storage" :min="1" />
        </el-form-item>
        <el-form-item label="并发请求(次/秒)">
          <el-input-number v-model="editForm.limits.concurrency" :min="1" />
        </el-form-item>
        <el-form-item label="可用模型数">
          <el-input-number v-model="editForm.limits.models" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="savePackage">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'
import { Check } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { packages } from '@/api/customer'

const packageList = ref(packages)
const dialogVisible = ref(false)
const dialogTitle = ref('')
const inputVisible = ref(false)
const inputValue = ref('')
const InputRef = ref()

const editForm = ref({
  id: null,
  name: '',
  price: 0,
  description: '',
  features: [],
  limits: {
    storage: 1,
    concurrency: 1,
    models: 1
  }
})

const resetForm = () => {
  editForm.value = {
    id: null,
    name: '',
    price: 0,
    description: '',
    features: [],
    limits: {
      storage: 1,
      concurrency: 1,
      models: 1
    }
  }
}

const handleAdd = () => {
  dialogTitle.value = '新增套餐'
  resetForm()
  dialogVisible.value = true
}

const editPackage = (pkg) => {
  dialogTitle.value = '编辑套餐'
  editForm.value = JSON.parse(JSON.stringify(pkg))
  dialogVisible.value = true
}

const deletePackage = (pkg) => {
  ElMessageBox.confirm(
    '确定要删除该套餐吗？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    const index = packageList.value.findIndex(item => item.id === pkg.id)
    if (index !== -1) {
      packageList.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  })
}

const savePackage = () => {
  if (!editForm.value.name || !editForm.value.description || editForm.value.features.length === 0) {
    ElMessage.warning('请填写完整的套餐信息')
    return
  }

  if (editForm.value.id) {
    // 编辑现有套餐
    const index = packageList.value.findIndex(item => item.id === editForm.value.id)
    if (index !== -1) {
      packageList.value[index] = { ...editForm.value }
      ElMessage.success('套餐更新成功')
    }
  } else {
    // 添加新套餐
    const newPackage = {
      ...editForm.value,
      id: packageList.value.length + 1
    }
    packageList.value.push(newPackage)
    ElMessage.success('套餐添加成功')
  }
  
  dialogVisible.value = false
}

const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    InputRef.value.focus()
  })
}

const addFeature = () => {
  if (inputValue.value) {
    editForm.value.features.push(inputValue.value)
    inputValue.value = ''
  }
  inputVisible.value = false
}

const removeFeature = (index) => {
  editForm.value.features.splice(index, 1)
}
</script>

<style scoped>
.packages-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.package-card {
  margin-bottom: 20px;
  transition: all 0.3s;
}

.package-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.package-card.enterprise {
  border: 2px solid #409EFF;
}

.package-header {
  text-align: center;
  margin-bottom: 20px;
}

.package-header h2 {
  margin: 0;
  color: #303133;
}

.price {
  font-size: 24px;
  color: #409EFF;
  margin-top: 10px;
}

.price span {
  font-size: 14px;
  color: #909399;
}

.description {
  color: #606266;
  margin-bottom: 20px;
}

.features {
  margin-bottom: 20px;
}

.features h4,
.limits h4 {
  margin: 0 0 10px 0;
  color: #303133;
}

.features ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.features li {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  color: #606266;
}

.features li .el-icon {
  color: #67C23A;
  margin-right: 8px;
}

.limit-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  color: #606266;
}

.package-actions {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 10px;
}

.input-new-tag {
  width: 90px;
  margin-left: 10px;
  vertical-align: bottom;
}
</style>