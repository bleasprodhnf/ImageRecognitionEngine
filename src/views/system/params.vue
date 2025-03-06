<template>
  <div class="system-params">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统参数配置</span>
        </div>
      </template>
      <div class="params-content">
        <el-form :model="form" label-width="120px">
          <el-form-item label="系统名称">
            <el-input v-model="form.systemName" placeholder="请输入系统名称" />
          </el-form-item>
          <el-form-item label="API接口地址">
            <el-input v-model="form.apiUrl" placeholder="请输入API接口地址" />
          </el-form-item>
          <el-form-item label="最大并发请求数">
            <el-input-number v-model="form.maxConcurrent" :min="1" :max="100" />
          </el-form-item>
          <el-form-item label="缓存过期时间(秒)">
            <el-input-number v-model="form.cacheExpiration" :min="0" :max="3600" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveSettings">保存设置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const form = ref({
  systemName: '图像识别引擎管理系统',
  apiUrl: 'http://api.example.com',
  maxConcurrent: 10,
  cacheExpiration: 300
})

const saveSettings = () => {
  // 这里添加保存设置的逻辑
  ElMessage.success('设置保存成功')
}
</script>

<style scoped>
.system-params {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.params-content {
  max-width: 600px;
  margin: 0 auto;
}
</style>