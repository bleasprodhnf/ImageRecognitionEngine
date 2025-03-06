<template>
  <div class="images-container">
    <el-row :gutter="20">
      <!-- 图片库统计信息 -->
      <el-col :span="24">
        <el-card class="stat-card">
          <div class="storage-info">
            <el-progress type="dashboard" :percentage="storageUsagePercent" :color="progressColor">
              <template #default="{ percentage }">
                <div class="storage-text">
                  <span class="current-storage">{{ usedStorage }}GB</span>
                  <span class="total-storage">共{{ totalStorage }}GB</span>
                </div>
              </template>
            </el-progress>
            <div class="stat-details">
              <div class="stat-item">
                <div class="stat-label">图片总数</div>
                <div class="stat-value">{{ totalImages }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">本月上传</div>
                <div class="stat-value">{{ monthlyUploads }}</div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 图片管理区域 -->
      <el-col :span="24" class="mt-20">
        <el-card>
          <template #header>
            <div class="card-header">
              <div class="left">
                <el-input
                  v-model="searchQuery"
                  placeholder="搜索图片"
                  prefix-icon="Search"
                  clearable
                  class="search-input"
                />
                <el-select v-model="currentCategory" placeholder="分类" class="ml-10">
                  <el-option label="全部" value="" />
                  <el-option
                    v-for="category in categories"
                    :key="category.id"
                    :label="category.name"
                    :value="category.id"
                  />
                </el-select>
              </div>
              <div class="right">
                <el-button type="primary" @click="showUploadDialog">上传图片</el-button>
                <el-button @click="showCategoryDialog">管理分类</el-button>
              </div>
            </div>
          </template>

          <!-- 图片列表 -->
          <el-row :gutter="20">
            <el-col
              v-for="image in filteredImages"
              :key="image.id"
              :xs="12"
              :sm="8"
              :md="6"
              :lg="4"
              class="mb-20"
            >
              <div class="image-card">
                <el-image
                  :src="image.url"
                  fit="cover"
                  @click="previewImage(image)"
                >
                  <template #error>
                    <div class="image-error">
                      <el-icon><Picture /></el-icon>
                    </div>
                  </template>
                </el-image>
                <div class="image-info">
                  <div class="image-name">{{ image.name }}</div>
                  <div class="image-meta">
                    <span>{{ formatSize(image.size) }}</span>
                    <span>{{ formatDate(image.uploadTime) }}</span>
                  </div>
                  <div class="image-actions">
                    <el-button-group>
                      <el-button size="small" @click="editImage(image)">编辑</el-button>
                      <el-button size="small" type="danger" @click="deleteImage(image)">删除</el-button>
                    </el-button-group>
                  </div>
                </div>
              </div>
            </el-col>
          </el-row>

          <!-- 分页 -->
          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              :page-sizes="[12, 24, 36, 48]"
              layout="total, sizes, prev, pager, next"
              @size-change="handleSizeChange"
              @current-change="handleCurrentChange"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传图片" width="500px">
      <el-upload
        class="upload-area"
        drag
        multiple
        :action="uploadUrl"
        :before-upload="beforeUpload"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            支持 jpg/png 文件，单个文件不超过 5MB
          </div>
        </template>
      </el-upload>
    </el-dialog>

    <!-- 分类管理对话框 -->
    <el-dialog v-model="categoryDialogVisible" title="分类管理" width="400px">
      <div class="category-list">
        <div v-for="category in categories" :key="category.id" class="category-item">
          <span>{{ category.name }}</span>
          <div class="category-actions">
            <el-button size="small" @click="editCategory(category)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteCategory(category)">删除</el-button>
          </div>
        </div>
        <el-button type="primary" class="add-category" @click="addCategory">添加分类</el-button>
      </div>
    </el-dialog>

    <!-- 图片预览 -->
    <el-image-viewer
      v-if="showViewer"
      :url-list="[previewUrl]"
      @close="closeViewer"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Picture, UploadFilled } from '@element-plus/icons-vue'

// 存储使用统计
const totalStorage = ref(100) // GB
const usedStorage = ref(45) // GB
const totalImages = ref(0)
const monthlyUploads = ref(0)

// 计算存储使用百分比
const storageUsagePercent = computed(() => {
  return Math.round((usedStorage.value / totalStorage.value) * 100)
})

// 进度条颜色
const progressColor = computed(() => {
  const percent = storageUsagePercent.value
  if (percent < 70) return '#67C23A'
  if (percent < 90) return '#E6A23C'
  return '#F56C6C'
})

// 图片列表相关
const searchQuery = ref('')
const currentCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(24)
const total = ref(0)

// 分类数据
const categories = ref([
  { id: 1, name: '未分类' },
  { id: 2, name: '产品图片' },
  { id: 3, name: '标志设计' }
])

// 图片数据
const images = ref([])
const filteredImages = computed(() => {
  let result = [...images.value]
  
  // 搜索过滤
  if (searchQuery.value) {
    result = result.filter(img => 
      img.name.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }
  
  // 分类过滤
  if (currentCategory.value) {
    result = result.filter(img => img.categoryId === currentCategory.value)
  }
  
  return result
})

// 对话框控制
const uploadDialogVisible = ref(false)
const categoryDialogVisible = ref(false)
const showViewer = ref(false)
const previewUrl = ref('')

// 上传配置
const uploadUrl = '/api/client/images/upload'

// 上传前验证
const beforeUpload = (file) => {
  const isImage = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt5M = file.size / 1024 / 1024 < 5

  if (!isImage) {
    ElMessage.error('只能上传 JPG/PNG 格式的图片！')
    return false
  }
  if (!isLt5M) {
    ElMessage.error('图片大小不能超过 5MB！')
    return false
  }
  return true
}

// 上传成功处理
const handleUploadSuccess = (response, file) => {
  ElMessage.success('上传成功')
  // TODO: 更新图片列表
}

// 上传失败处理
const handleUploadError = () => {
  ElMessage.error('上传失败，请重试')
}

// 显示上传对话框
const showUploadDialog = () => {
  uploadDialogVisible.value = true
}

// 显示分类管理对话框
const showCategoryDialog = () => {
  categoryDialogVisible.value = true
}

// 预览图片
const previewImage = (image) => {
  previewUrl.value = image.url
  showViewer.value = true
}

// 关闭预览
const closeViewer = () => {
  showViewer.value = false
}

// 编辑图片
const editImage = (image) => {
  // TODO: 实现图片编辑功能
}

// 删除图片
const deleteImage = (image) => {
  ElMessageBox.confirm('确定要删除这张图片吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    // TODO: 实现删除逻辑
    ElMessage.success('删除成功')
  }).catch(() => {})
}

// 添加分类
const addCategory = () => {
  // TODO: 实现添加分类功能
}

// 编辑分类
const editCategory = (category) => {
  // TODO: 实现编辑分类功能
}

// 删除分类
const deleteCategory = (category) => {
  // TODO: 实现删除分类功能
}

// 格式化文件大小
const formatSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
const formatDate = (date) => {
  return new Date(date).toLocaleDateString()
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  // TODO: 重新加载数据
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  // TODO: 重新加载数据
}

// 初始化数据
onMounted(() => {
  // TODO: 加载初始数据
})
</script>

<style scoped>
.images-container {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.storage-info {
  display: flex;
  justify-content: space-around;
  align-items: center;
  padding: 20px;
}

.storage-text {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.current-storage {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.total-storage {
  font-size: 14px;
  color: #909399;
}

.stat-details {
  display: flex;
  gap: 40px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-input {
  width: 200px;
}

.image-card {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
  transition: all 0.3s;
}

.image-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.el-image {
  width: 100%;
  height: 150px;
  display: block;
}

.image-error {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
}

.image-info {
  padding: 10px;
}

.image-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.image-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 12px;
  color: #909399;
}
</style>