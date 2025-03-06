<template>
  <div class="image-labels-container">
    <el-card class="label-card">
      <template #header>
        <div class="card-header">
          <span>自定义图片标识库</span>
          <el-button type="primary" size="small" @click="showUploadDialog">上传新图片</el-button>
        </div>
      </template>
      
      <div class="filter-section">
        <el-input
          v-model="searchQuery"
          placeholder="搜索图片或标签"
          prefix-icon="Search"
          clearable
          class="search-input"
        />
        <el-select v-model="labelTypeFilter" placeholder="标签类型" clearable class="ml-10">
          <el-option label="全部" value="" />
          <el-option label="人物" value="person" />
          <el-option label="物品" value="object" />
          <el-option label="动物" value="animal" />
          <el-option label="场景" value="scene" />
        </el-select>
      </div>
      
      <div class="labeled-images">
        <el-empty v-if="filteredImages.length === 0" description="暂无标记图片" />
        <el-row v-else :gutter="20">
          <el-col 
            v-for="image in filteredImages" 
            :key="image.id" 
            :xs="24" 
            :sm="12" 
            :md="8" 
            :lg="6" 
            class="mb-20"
          >
            <div class="labeled-image-card">
              <div class="image-container">
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
                <div class="image-overlay" @click="editLabels(image)">
                  <el-icon><Edit /></el-icon>
                  <span>编辑标签</span>
                </div>
              </div>
              <div class="image-info">
                <div class="image-name">{{ image.name }}</div>
                <div class="image-labels">
                  <el-tag 
                    v-for="label in image.labels" 
                    :key="label.id"
                    :type="getTagType(label.type)"
                    size="small"
                    class="label-tag"
                  >
                    {{ label.name }}
                  </el-tag>
                </div>
                <div class="image-actions">
                  <el-button size="small" @click="editLabels(image)">编辑标签</el-button>
                  <el-button size="small" type="danger" @click="deleteImage(image)">删除</el-button>
                </div>
              </div>
            </div>
          </el-col>
        </el-row>
        
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
      </div>
    </el-card>
    
    <!-- 上传图片对话框 -->
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
    
    <!-- 编辑标签对话框 -->
    <el-dialog v-model="labelDialogVisible" title="编辑图片标签" width="600px">
      <div v-if="currentImage" class="edit-labels-container">
        <div class="image-preview">
          <el-image :src="currentImage.url" fit="contain" />
        </div>
        
        <div class="recognition-results">
          <h4>AI识别结果</h4>
          <div class="recognition-content">
            <p>{{ currentImage.aiDescription || '正在分析图片内容...' }}</p>
          </div>
        </div>
        
        <div class="label-editor">
          <h4>自定义标签</h4>
          <p class="hint-text">根据AI识别结果，您可以为以下类别添加自定义标签：</p>
          
          <div v-for="category in availableCategories" :key="category.type" class="label-category">
            <div class="category-header">
              <h5>{{ getCategoryName(category.type) }}</h5>
              <el-button 
                type="primary" 
                size="small" 
                plain 
                @click="addNewLabel(category.type)"
                :disabled="!category.available"
              >
                添加标签
              </el-button>
            </div>
            
            <div class="category-labels">
              <el-tag
                v-for="label in getLabelsForCategory(category.type)"
                :key="label.id"
                closable
                :type="getTagType(category.type)"
                @close="removeLabel(label.id)"
                class="label-tag"
              >
                {{ label.name }}
              </el-tag>
              <div v-if="!getLabelsForCategory(category.type).length" class="no-labels">
                {{ category.available ? '暂无标签' : '此类别不适用于当前图片' }}
              </div>
            </div>
          </div>
        </div>
        
        <div class="dialog-footer">
          <el-button @click="labelDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveLabels">保存</el-button>
        </div>
      </div>
    </el-dialog>
    
    <!-- 添加新标签对话框 -->
    <el-dialog v-model="newLabelDialogVisible" title="添加新标签" width="400px" append-to-body>
      <el-form :model="newLabelForm" label-width="80px">
        <el-form-item label="标签名称" required>
          <el-input v-model="newLabelForm.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="标签类型">
          <el-tag :type="getTagType(newLabelForm.type)">{{ getCategoryName(newLabelForm.type) }}</el-tag>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="newLabelDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmAddLabel">确定</el-button>
        </span>
      </template>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Picture, Edit, UploadFilled } from '@element-plus/icons-vue'

// 搜索和筛选
const searchQuery = ref('')
const labelTypeFilter = ref('')
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 对话框控制
const uploadDialogVisible = ref(false)
const labelDialogVisible = ref(false)
const newLabelDialogVisible = ref(false)
const showViewer = ref(false)
const previewUrl = ref('')

// 当前编辑的图片
const currentImage = ref(null)

// 新标签表单
const newLabelForm = reactive({
  name: '',
  type: ''
})

// 上传配置
const uploadUrl = '/api/client/images/upload'

// 模拟图片数据
const images = ref([
  {
    id: 1,
    name: '团队合影.jpg',
    url: 'https://example.com/images/team.jpg',
    aiDescription: '图片中有5个人站在办公室内，背景是公司logo墙。',
    labels: [
      { id: 1, name: '张三', type: 'person' },
      { id: 2, name: '李四', type: 'person' },
      { id: 3, name: '公司活动', type: 'scene' }
    ],
    availableCategories: ['person', 'scene']
  },
  {
    id: 2,
    name: '产品展示.jpg',
    url: 'https://example.com/images/product.jpg',
    aiDescription: '图片中是一个放在桌面上的黑色智能设备，旁边有一部手机。',
    labels: [
      { id: 4, name: 'AI识别器', type: 'object' },
      { id: 5, name: '产品展示', type: 'scene' }
    ],
    availableCategories: ['object', 'scene']
  },
  {
    id: 3,
    name: '公司宠物.jpg',
    url: 'https://example.com/images/pet.jpg',
    aiDescription: '图片中是一只橘色的猫咪躺在办公桌上。',
    labels: [
      { id: 6, name: '小橘', type: 'animal' },
      { id: 7, name: '办公室', type: 'scene' }
    ],
    availableCategories: ['animal', 'scene']
  }
])

// 过滤图片
const filteredImages = computed(() => {
  let result = [...images.value]
  
  // 搜索过滤
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(img => {
      // 搜索图片名称
      if (img.name.toLowerCase().includes(query)) return true
      
      // 搜索标签名称
      return img.labels.some(label => label.name.toLowerCase().includes(query))
    })
  }
  
  // 标签类型过滤
  if (labelTypeFilter.value) {
    result = result.filter(img => {
      return img.labels.some(label => label.type === labelTypeFilter.value)
    })
  }
  
  return result
})

// 可用的标签类别
const availableCategories = computed(() => {
  if (!currentImage.value) return []
  
  const categories = [
    { type: 'person', name: '人物', available: false },
    { type: 'object', name: '物品', available: false },
    { type: 'animal', name: '动物', available: false },
    { type: 'scene', name: '场景', available: false }
  ]
  
  // 根据AI识别结果设置可用类别
  return categories.map(category => {
    return {
      ...category,
      available: currentImage.value.availableCategories.includes(category.type)
    }
  })
})

// 获取标签类型对应的名称
const getCategoryName = (type) => {
  const map = {
    person: '人物',
    object: '物品',
    animal: '动物',
    scene: '场景'
  }
  return map[type] || type
}

// 获取标签类型对应的Tag类型
const getTagType = (type) => {
  const map = {
    person: '',
    object: 'success',
    animal: 'warning',
    scene: 'info'
  }
  return map[type] || ''
}

// 获取特定类别的标签
const getLabelsForCategory = (type) => {
  if (!currentImage.value) return []
  return currentImage.value.labels.filter(label => label.type === type)
}

// 显示上传对话框
const showUploadDialog = () => {
  uploadDialogVisible.value = true
}

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
  uploadDialogVisible.value = false
  
  // 模拟上传后的AI分析结果
  const newImage = {
    id: Date.now(),
    name: file.name,
    url: URL.createObjectURL(file.raw),
    aiDescription: '正在分析图片内容...',
    labels: [],
    availableCategories: []
  }
  
  // 添加到图片列表
  images.value.unshift(newImage)
  
  // 模拟AI分析过程
  setTimeout(() => {
    // 模拟AI分析结果
    const aiResults = simulateAIAnalysis(file.name)
    const imageIndex = images.value.findIndex(img => img.id === newImage.id)
    
    if (imageIndex !== -1) {
      images.value[imageIndex].aiDescription = aiResults.description
      images.value[imageIndex].availableCategories = aiResults.categories
    }
  }, 1500)
}

// 模拟AI分析结果
const simulateAIAnalysis = (filename) => {
  // 根据文件名模拟不同的分析结果
  if (filename.toLowerCase().includes('人') || filename.toLowerCase().includes('合影')) {
    return {
      description: '图片中包含多个人物，可能是在室内环境。',
      categories: ['person', 'scene']
    }
  } else if (filename.toLowerCase().includes('产品') || filename.toLowerCase().includes('设备')) {
    return {
      description: '图片中是一个产品或设备，放置在平面上。',
      categories: ['object', 'scene']
    }
  } else if (filename.toLowerCase().includes('动物') || filename.toLowerCase().includes('宠物')) {
    return {
      description: '图片中有一只动物，可能是宠物。',
      categories: ['animal', 'scene']
    }
  } else {
    // 随机生成分析结果
    const categories = []
    const types = ['person', 'object', 'animal', 'scene']
    const randomCount = Math.floor(Math.random() * 3) + 1
    
    for (let i = 0; i < randomCount; i++) {
      const randomType = types[Math.floor(Math.random() * types.length)]
      if (!categories.includes(randomType)) {
        categories.push(randomType)
      }
    }
    
    // 场景几乎总是可用的
    if (!categories.includes('scene')) {
      categories.push('scene')
    }
    
    return {
      description: '图片内容分析完成，请根据识别结果添加自定义标签。',
      categories
    }
  }
}

// 上传失败处理
const handleUploadError = () => {
  ElMessage.error('上传失败，请重试')
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

// 编辑标签
const editLabels = (image) => {
  currentImage.value = JSON.parse(JSON.stringify(image))
  labelDialogVisible.value = true
}

// 添加新标签
const addNewLabel = (type) => {
  newLabelForm.type = type
  newLabelForm.name = ''
  newLabelDialogVisible.value = true
}

// 确认添加标签
const confirmAddLabel = () => {
  if (!newLabelForm.name.trim()) {
    ElMessage.warning('请输入标签名称')
    return
  }
  
  // 添加新标签
  const newLabel = {
    id: Date.now(),
    name: newLabelForm.name.trim(),
    type: newLabelForm.type
  }
  
  currentImage.value.labels.push(newLabel)
  newLabelDialogVisible.value = false
}

// 移除标签
const removeLabel = (labelId) => {
  if (!currentImage.value) return
  
  const index = currentImage.value.labels.findIndex(label => label.id === labelId)
  if (index !== -1) {
    currentImage.value.labels.splice(index, 1)
  }
}

// 保存标签
const saveLabels = () => {
  if (!currentImage.value) return
  
  // 找到原始图片并更新标签
  const index = images.value.findIndex(img => img.id === currentImage.value.id)
  if (index !== -1) {
    images.value[index].labels = [...currentImage.value.labels]
    ElMessage.success('标签已保存')
  }
  
  labelDialogVisible.value = false
}

// 删除图片
const deleteImage = (image) => {
  ElMessageBox.confirm('确定要删除这张图片吗？删除后将无法恢复。', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    const index = images.value.findIndex(img => img.id === image.id)
    if (index !== -1) {
      images.value.splice(index, 1)
      ElMessage.success('图片已删除')
    }
  }).catch(() => {})
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
}

const handleCurrentChange = (val) => {
  currentPage.value = val
}

// 初始化数据
onMounted(() => {
  // 从服务器获取图片和标签数据
  // TODO: 实现API调用
  total.value = images.value.length
})
</script>

<style scoped>
.image-labels-container {
  padding: 20px;
}

.label-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.filter-section {
  display: flex;
  margin-bottom: 20px;
}

.search-input {
  width: 250px;
}

.ml-10 {
  margin-left: 10px;
}

.mb-20 {
  margin-bottom: 20px;
}

.labeled-image-card {
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
  transition: all 0.3s;
}

.labeled-image-card:hover {
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.image-container {
  position: relative;
  height: 150px;
}

.el-image {
  width: 100%;
  height: 100%;
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

.image-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.5);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  opacity: 0;
  transition: opacity 0.3s;
  cursor: pointer;
  color: white;
}

.image-container:hover .image-overlay {
  opacity: 1;
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

.image-labels {
  margin-bottom: 10px;
  min-height: 24px;
}

.label-tag {
  margin-right: 5px;
  margin-bottom: 5px;
}

.image-actions {
  display: flex;
  justify-content: space-between;
}

.upload-area {
  width: 100%;
}

.edit-labels-container {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.image-preview {
  text-align: center;
  max-height: 200px;
  overflow: hidden;
}

.recognition-results {
  background-color: #f5f7fa;
  padding: 15px;
  border-radius: 4px;
}

.recognition-content {
  color: #606266;
}

.label-editor {
  margin-top: 10px;
}

.hint-text {
  color: #909399;
  font-size: 14px;
  margin-bottom: 15px;
}

.label-category {
  margin-bottom: 15px;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.category-header h5 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.category-labels {
  min-height: 32px;
  padding: 5px 0;
}

.no-labels {
  color: #909399;
  font-size: 14px;
  font-style: italic;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>