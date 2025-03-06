<template>
  <div class="home-container">
    <!-- 顶部横幅 -->
    <div class="hero-section">
      <div class="hero-content">
        <h1>智能图像识别引擎</h1>
        <p class="subtitle">基于先进AI技术的高精度图像识别解决方案</p>
        <div class="hero-buttons">
          <el-button type="primary" size="large" @click="scrollToDemo">在线体验</el-button>
          <el-button size="large" @click="scrollToFeatures">了解功能</el-button>
        </div>
      </div>
      <div class="hero-image">
        <img :src="heroImage" alt="AI图像识别" />
      </div>
    </div>

    <!-- 核心功能展示 -->
    <div class="features-section" ref="featuresSection">
      <h2 class="section-title">核心功能</h2>
      <div class="features-grid">
        <div class="feature-card">
          <el-icon class="feature-icon"><Aim /></el-icon>
          <h3>高精度识别</h3>
          <p>基于深度学习算法，识别准确率高达98%，支持多种物体同时识别</p>
        </div>
        <div class="feature-card">
          <el-icon class="feature-icon"><Lightning /></el-icon>
          <h3>实时处理</h3>
          <p>毫秒级响应速度，支持大规模并发请求，满足各种实时应用场景</p>
        </div>
        <div class="feature-card">
          <el-icon class="feature-icon"><Connection /></el-icon>
          <h3>简单集成</h3>
          <p>提供RESTful API和多语言SDK，轻松集成到各类应用中</p>
        </div>
        <div class="feature-card">
          <el-icon class="feature-icon"><DataAnalysis /></el-icon>
          <h3>数据分析</h3>
          <p>提供详细的识别结果分析和统计报表，助力业务决策</p>
        </div>
      </div>
    </div>

    <!-- 应用场景 -->
    <div class="scenarios-section">
      <h2 class="section-title">应用场景</h2>
      <el-carousel :interval="4000" type="card" height="300px">
        <el-carousel-item v-for="(scenario, index) in applicationScenarios" :key="index">
          <div class="scenario-card">
            <div class="scenario-image">
              <img :src="scenario.image" :alt="scenario.title" />
            </div>
            <div class="scenario-content">
              <h3>{{ scenario.title }}</h3>
              <p>{{ scenario.description }}</p>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 在线体验区域 -->
    <div class="demo-section" ref="demoSection">
      <h2 class="section-title">在线体验</h2>
      <div class="demo-container">
        <div class="upload-area">
          <h3>上传图片进行识别</h3>
          <el-upload
            class="image-uploader"
            action="http://localhost:8080/api/v1/client/recognize"
            :headers="{
              'X-App-ID': 'demo_app_id',
              'X-API-Key': 'demo_api_key'
            }"
            :auto-upload="true"
            :show-file-list="false"
            :before-upload="handleImageChange"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            accept="image/jpeg,image/png,image/gif">
            <img v-if="imageUrl" :src="imageUrl" class="preview-image" />
            <div v-else class="upload-placeholder">
              <el-icon class="upload-icon"><Plus /></el-icon>
              <div class="upload-text">点击上传图片</div>
              <div class="upload-hint">支持jpg、png、gif格式，大小不超过10MB</div>
            </div>
          </el-upload>
          <el-button type="primary" :disabled="!imageUrl" @click="recognizeImage">开始识别</el-button>
        </div>

        <div class="result-area" v-if="recognitionResult">
          <h3>识别结果</h3>
          <div class="result-content">
            <div class="result-image-container">
              <img :src="imageUrl" class="result-image" />
              <div 
                v-for="(item, index) in recognitionResult.results" 
                :key="index"
                class="bounding-box"
                :style="{
                  left: item.bbox[0] + 'px',
                  top: item.bbox[1] + 'px',
                  width: (item.bbox[2] - item.bbox[0]) + 'px',
                  height: (item.bbox[3] - item.bbox[1]) + 'px'
                }">
                <div class="label">{{ item.label }} ({{ (item.confidence * 100).toFixed(0) }}%)</div>
              </div>
            </div>
            <div class="result-details">
              <div class="result-item">
                <span class="result-label">识别ID:</span>
                <span class="result-value">{{ recognitionResult.id }}</span>
              </div>
              <div class="result-item">
                <span class="result-label">处理时间:</span>
                <span class="result-value">{{ recognitionResult.processing_time }}秒</span>
              </div>
              <div class="result-item">
                <span class="result-label">识别结果:</span>
                <span class="result-value">{{ recognitionResult.results.length }}个对象</span>
              </div>
              <el-table :data="recognitionResult.results" style="width: 100%">
                <el-table-column prop="label" label="对象" width="120" />
                <el-table-column prop="confidence" label="置信度" width="120">
                  <template #default="scope">
                    {{ (scope.row.confidence * 100).toFixed(2) }}%
                  </template>
                </el-table-column>
                <el-table-column prop="bbox" label="位置">
                  <template #default="scope">
                    [{{ scope.row.bbox.join(', ') }}]
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 价格方案 -->
    <div class="pricing-section">
      <h2 class="section-title">价格方案</h2>
      <div class="pricing-cards">
        <div class="pricing-card">
          <div class="pricing-header">
            <h3>基础版</h3>
            <div class="price">¥0<span class="price-period">/月</span></div>
          </div>
          <div class="pricing-features">
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>每日100次API调用</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>基础图像识别功能</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-close"><Close /></el-icon>
              <span>批量识别</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-close"><Close /></el-icon>
              <span>技术支持</span>
            </div>
          </div>
          <el-button>免费开始</el-button>
        </div>

        <div class="pricing-card popular">
          <div class="popular-tag">最受欢迎</div>
          <div class="pricing-header">
            <h3>专业版</h3>
            <div class="price">¥299<span class="price-period">/月</span></div>
          </div>
          <div class="pricing-features">
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>每日10,000次API调用</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>高级图像识别功能</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>批量识别</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>邮件技术支持</span>
            </div>
          </div>
          <el-button type="primary">立即购买</el-button>
        </div>

        <div class="pricing-card">
          <div class="pricing-header">
            <h3>企业版</h3>
            <div class="price">¥999<span class="price-period">/月</span></div>
          </div>
          <div class="pricing-features">
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>无限API调用</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>全部高级功能</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>定制化开发</span>
            </div>
            <div class="feature-item">
              <el-icon class="feature-icon-check"><Check /></el-icon>
              <span>7×24小时专属支持</span>
            </div>
          </div>
          <el-button>联系我们</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { Aim, Lightning, Connection, DataAnalysis, Lock, Refresh, Plus, Check, Close } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';
import heroImage from '@/assets/images/hero-image.svg'

// 应用场景数据
const applicationScenarios = [
  {
    title: '智慧零售',
    description: '自动识别商品、分析顾客行为，提升购物体验和运营效率',
    image: '/src/assets/images/retail.jpg'
  },
  {
    title: '安防监控',
    description: '实时识别可疑人员和行为，提高安全防范能力',
    image: '/src/assets/images/security.jpg'
  },
  {
    title: '医疗影像',
    description: '辅助医生分析X光、CT等医疗影像，提高诊断准确率',
    image: '/src/assets/images/medical.jpg'
  },
  {
    title: '智能驾驶',
    description: '识别道路标志、行人和车辆，为自动驾驶提供视觉感知能力',
    image: '/src/assets/images/driving.jpg'
  },
  {
    title: '工业质检',
    description: '自动检测产品缺陷，提高生产质量和效率',
    image: '/src/assets/images/industry.jpg'
  }
];

// 图片上传和识别
const imageUrl = ref('');
const recognitionResult = ref(null);
const isRecognizing = ref(false);
const uploadedFile = ref(null);

// 处理图片上传
const handleImageChange = (file) => {
  // 验证文件类型
  const isJPG = file.type === 'image/jpeg';
  const isPNG = file.type === 'image/png';
  const isGIF = file.type === 'image/gif';
  
  if (!isJPG && !isPNG && !isGIF) {
    ElMessage.error('上传图片只能是 JPG/PNG/GIF 格式!');
    return false;
  }
  
  // 验证文件大小
  const isLt10M = file.size / 1024 / 1024 < 10;
  if (!isLt10M) {
    ElMessage.error('上传图片大小不能超过 10MB!');
    return false;
  }
  
  // 保存文件对象以便在上传成功后使用
  uploadedFile.value = file;
  return true;
};

// 上传成功处理
const handleUploadSuccess = (response) => {
  // 根据后端返回的数据结构调整
  if (response.code === 200) {
    // 使用本地图片URL，因为我们已经上传了图片
    imageUrl.value = URL.createObjectURL(uploadedFile.value);
    // 存储识别结果
    recognitionResult.value = response.data;
    ElMessage.success('上传并识别成功');
  } else {
    ElMessage.error(response.message || '上传失败');
  }
};

// 上传失败处理
const handleUploadError = (err) => {
  console.error('上传错误:', err);
  const errorMessage = err.response?.data?.message || '上传失败，请检查网络连接或稍后重试';
  ElMessage.error(errorMessage);
  imageUrl.value = '';
  recognitionResult.value = null;
};

// 识别图片
const recognizeImage = () => {
  if (!imageUrl.value) {
    ElMessage.warning('请先上传图片');
    return;
  }

  isRecognizing.value = true;
  ElMessage.info('正在识别图片...');

  // 模拟API调用
  setTimeout(() => {
    isRecognizing.value = false;
    
    // 模拟识别结果
    recognitionResult.value = {
      id: 'rec_' + Math.floor(Math.random() * 100000),
      results: [
        {
          label: '人物',
          confidence: 0.96,
          bbox: [50, 30, 200, 300]
        },
        {
          label: '建筑',
          confidence: 0.85,
          bbox: [300, 100, 500, 400]
        },
        {
          label: '树木',
          confidence: 0.78,
          bbox: [450, 200, 550, 350]
        }
      ],
      processing_time: 0.32
    };

    ElMessage.success('识别完成');
  }, 1500);
};

// 滚动到指定区域
const featuresSection = ref(null);
const demoSection = ref(null);

const scrollToFeatures = () => {
  featuresSection.value.scrollIntoView({ behavior: 'smooth' });
};

const scrollToDemo = () => {
  demoSection.value.scrollIntoView({ behavior: 'smooth' });
};
</script>

<style scoped>
.home-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

/* 顶部横幅 */
.hero-section {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 60px 0;
  gap: 40px;
}

.hero-content {
  flex: 1;
}

.hero-content h1 {
  font-size: 2.5rem;
  margin-bottom: 20px;
  color: #303133;
}

.subtitle {
  font-size: 1.2rem;
  color: #606266;
  margin-bottom: 30px;
}

.hero-buttons {
  display: flex;
  gap: 15px;
}

.hero-image {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
}

.hero-image img {
  max-width: 100%;
  height: auto;
}

/* 核心功能展示 */
.section-title {
  text-align: center;
  font-size: 2rem;
  margin-bottom: 40px;
  color: #303133;
}

.features-section {
  padding: 60px 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 30px;
}

.feature-card {
  background-color: #fff;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  text-align: center;
  transition: transform 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-5px);
}

.feature-icon {
  font-size: 2.5rem;
  color: #409EFF;
  margin-bottom: 20px;
}

.feature-card h3 {
  font-size: 1.2rem;
  margin-bottom: 15px;
  color: #303133;
}

.feature-card p {
  color: #606266;
  line-height: 1.6;
}

/* 应用场景 */
.scenarios-section {
  padding: 60px 0;
  background-color: #f5f7fa;
}

.scenario-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.scenario-image {
  height: 180px;
  overflow: hidden;
}

.scenario-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scenario-content {
  padding: 20px;
}

.scenario-content h3 {
  margin-top: 0;
  margin-bottom: 10px;
  color: #303133;
}

.scenario-content p {
  color: #606266;
  margin: 0;
}

/* 在线体验区域 */
.demo-section {
  padding: 60px 0;
}

.demo-container {
  display: flex;
  gap: 30px;
  flex-wrap: wrap;
}

.upload-area, .result-area {
  flex: 1;
  min-width: 300px;
  background-color: #fff;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.upload-area h3, .result-area h3 {
  margin-top: 0;
  margin-bottom: 20px;
  color: #303133;
}

.image-uploader {
  margin-bottom: 20px;
}

.upload-placeholder {
  width: 100%;
  height: 200px;
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: border-color 0.3s;
}

.upload-placeholder:hover {
  border-color: #409EFF;
}

.upload-icon {
  font-size: 28px;
  color: #8c939d;
}

.upload-text {
  margin-top: 10px;
  color: #606266;
}

.upload-hint {
  margin-top: 5px;
  font-size: 12px;
  color: #909399;
}

.preview-image {
  width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.result-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.result-image-container {
  position: relative;
  margin-bottom: 20px;
}

.result-image {
  width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.bounding-box {
  position: absolute;
  border: 2px solid #409EFF;
  background-color: rgba(64, 158, 255, 0.1);
}

.label {
  position: absolute;
  top: -25px;
  left: 0;
  background-color: #409EFF;
  color: white;
  padding: 2px 6px;
  font-size: 12px;
  border-radius: 4px;
}

.result-details {
  width: 100%;
}

.result-item {
  margin-bottom: 10px;
  display: flex;
}

.result-label {
  font-weight: bold;
  margin-right: 10px;
  color: #606266;
  min-width: 80px;
}

.result-value {
  color: #303133;
}

/* 价格方案 */
.pricing-section {
  padding: 60px 0;
  background-color: #f5f7fa;
}

.pricing-cards {
  display: flex;
  justify-content: center;
  gap: 30px;
  flex-wrap: wrap;
}

.pricing-card {
  flex: 1;
  min-width: 280px;
  max-width: 350px;
  background-color: #fff;
  border-radius: 8px;
  padding: 30px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  position: relative;
  display: flex;
  flex-direction: column;
}

.pricing-card.popular {
  border: 2px solid #409EFF;
  transform: translateY(-10px);
}

.popular-tag {
  position: absolute;
  top: -12px;
  right: 30px;
  background-color: #409EFF;
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: bold;
}

.pricing-header {
  text-align: center;
  margin-bottom: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.pricing-header h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #303133;
  font-size: 1.5rem;
}

.price {
  font-size: 2.5rem;
  font-weight: bold;
  color: #303133;
}

.price-period {
  font-size: 1rem;
  font-weight: normal;
  color: #909399;
}

.pricing-features {
  flex: 1;
  margin-bottom: 30px;
}

.feature-item {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.feature-icon-check {
  color: #67C23A;
  margin-right: 10px;
}

.feature-icon-close {
  color: #F56C6C;
  margin-right: 10px;
}

.pricing-card .el-button {
  width: 100%;
}

/* 响应式布局 */
@media screen and (max-width: 992px) {
  .hero-section {
    flex-direction: column;
    text-align: center;
  }
  
  .hero-content {
    order: 1;
  }
  
  .hero-image {
    order: 0;
    margin-bottom: 30px;
  }
  
  .hero-buttons {
    justify-content: center;
  }
  
  .features-grid {
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  }
}

@media screen and (max-width: 768px) {
  .hero-content h1 {
    font-size: 2rem;
  }
  
  .subtitle {
    font-size: 1rem;
  }
  
  .section-title {
    font-size: 1.8rem;
  }
  
  .demo-container {
    flex-direction: column;
  }
  
  .pricing-card {
    max-width: none;
  }
  
  .pricing-card.popular {
    transform: none;
  }
}

@media screen and (max-width: 480px) {
  .hero-content h1 {
    font-size: 1.8rem;
  }
  
  .hero-buttons {
    flex-direction: column;
    gap: 10px;
  }
  
  .section-title {
    font-size: 1.5rem;
  }
  
  .feature-card {
    padding: 20px;
  }
}
</style>