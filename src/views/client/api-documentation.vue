<template>
  <div class="api-documentation-container">
    <el-card class="doc-card">
      <template #header>
        <div class="card-header">
          <span>API 文档</span>
          <el-button type="primary" size="small" @click="goBack">
            <el-icon class="el-icon--left"><Back /></el-icon>返回API设置
          </el-button>
        </div>
      </template>
      
      <div class="doc-content">
        <div class="doc-navigation">
          <el-affix :offset="80">
            <el-card class="nav-card">
              <h3>目录导航</h3>
              <el-menu default-active="1">
                <el-menu-item index="1" @click="scrollToSection('overview')">API 概述</el-menu-item>
                <el-menu-item index="2" @click="scrollToSection('auth')">认证方式</el-menu-item>
                <el-menu-item index="3" @click="scrollToSection('base-endpoint')">基础端点</el-menu-item>
                <el-menu-item index="4" @click="scrollToSection('endpoints')">API 端点</el-menu-item>
                <el-menu-item index="5" @click="scrollToSection('error-codes')">错误码</el-menu-item>
                <el-menu-item index="6" @click="scrollToSection('limits')">使用限制</el-menu-item>
                <el-menu-item index="7" @click="scrollToSection('sdk')">SDK和代码示例</el-menu-item>
                <el-menu-item index="8" @click="scrollToSection('tester')">API 测试工具</el-menu-item>
              </el-menu>
            </el-card>
          </el-affix>
        </div>
        
        <div class="doc-main-content">
          <section id="overview" class="doc-section">
            <h2>API 概述</h2>
            <p>本API提供图像识别和分析功能，支持多种识别模式和参数配置。您可以通过REST API方式调用我们的服务，获取高精度的图像识别结果。</p>
          </section>
          
          <el-divider></el-divider>
          
          <section id="auth" class="doc-section">
            <h2>认证方式</h2>
            <p>所有API请求都需要进行身份验证。您需要在请求头中包含以下两个字段：</p>
            <div class="code-block">
              <pre><code>X-App-ID: 您的应用ID
X-API-Key: 您的API密钥</code></pre>
              <el-button size="small" type="primary" class="copy-btn" @click="copyCode('auth-code')">
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
            <p>您可以在<router-link to="/client/api-settings">API设置</router-link>页面查看和管理您的应用ID和API密钥。</p>
          </section>
          
          <el-divider></el-divider>
          
          <section id="base-endpoint" class="doc-section">
            <h2>基础端点</h2>
            <p>API的基础URL为：<code>https://api.example.com/v1</code></p>
          </section>
          
          <el-divider></el-divider>
          
          <section id="endpoints" class="doc-section">
            <h2>API 端点</h2>
            
            <el-collapse accordion>
              <el-collapse-item title="图像识别 /recognize" name="1">
                <h3>请求方法</h3>
                <p><code>POST /v1/recognize</code></p>
                
                <h3>请求参数</h3>
                <el-table :data="recognizeParams" style="width: 100%" border stripe>
                  <el-table-column prop="param" label="参数名" width="150" />
                  <el-table-column prop="type" label="类型" width="100" />
                  <el-table-column prop="required" label="必填" width="80" />
                  <el-table-column prop="description" label="说明" />
                </el-table>
                
                <h3>请求示例</h3>
                <div class="code-block" id="recognize-url-example">
                  <pre><code>// 使用图片URL
POST /v1/recognize
Content-Type: application/json
X-App-ID: your_app_id
X-API-Key: your_api_key

{
  "image_url": "https://example.com/image.jpg",
  "min_confidence": 0.7,
  "max_results": 10
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('recognize-url-example')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
                
                <div class="code-block" id="recognize-base64-example">
                  <pre><code>// 使用Base64编码的图片数据
POST /v1/recognize
Content-Type: application/json
X-App-ID: your_app_id
X-API-Key: your_api_key

{
  "image_data": "base64编码的图片数据",
  "min_confidence": 0.7,
  "max_results": 10
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('recognize-base64-example')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
                
                <h3>响应格式</h3>
                <div class="code-block" id="recognize-response">
                  <pre><code>{
  "success": true,
  "data": {
    "id": "rec_12345",
    "results": [
      {
        "label": "对象1",
        "confidence": 0.98,
        "bbox": [10, 10, 100, 100]
      },
      {
        "label": "对象2",
        "confidence": 0.85,
        "bbox": [150, 50, 200, 200]
      }
    ],
    "processing_time": 0.45
  }
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('recognize-response')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
              </el-collapse-item>
              
              <el-collapse-item title="批量识别 /batch-recognize" name="2">
                <h3>请求方法</h3>
                <p><code>POST /v1/batch-recognize</code></p>
                
                <h3>请求参数</h3>
                <el-table :data="batchRecognizeParams" style="width: 100%" border stripe>
                  <el-table-column prop="param" label="参数名" width="150" />
                  <el-table-column prop="type" label="类型" width="100" />
                  <el-table-column prop="required" label="必填" width="80" />
                  <el-table-column prop="description" label="说明" />
                </el-table>
                
                <h3>请求示例</h3>
                <div class="code-block" id="batch-recognize-example">
                  <pre><code>POST /v1/batch-recognize
Content-Type: application/json
X-App-ID: your_app_id
X-API-Key: your_api_key

{
  "images": [
    { "id": "img1", "url": "https://example.com/image1.jpg" },
    { "id": "img2", "url": "https://example.com/image2.jpg" }
  ],
  "min_confidence": 0.7
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('batch-recognize-example')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
                
                <h3>响应格式</h3>
                <div class="code-block" id="batch-recognize-response">
                  <pre><code>{
  "success": true,
  "data": {
    "batch_id": "batch_12345",
    "results": [
      {
        "image_id": "img1",
        "status": "success",
        "results": [
          {
            "label": "汽车",
            "confidence": 0.92,
            "bbox": [10, 20, 300, 200]
          }
        ]
      },
      {
        "image_id": "img2",
        "status": "success",
        "results": [
          {
            "label": "树木",
            "confidence": 0.88,
            "bbox": [150, 50, 250, 350]
          }
        ]
      }
    ],
    "processing_time": 0.75
  }
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('batch-recognize-response')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
              </el-collapse-item>
              
              <el-collapse-item title="查询识别结果 /results/{id}" name="3">
                <h3>请求方法</h3>
                <p><code>GET /v1/results/{id}</code></p>
                
                <h3>路径参数</h3>
                <el-table :data="resultPathParams" style="width: 100%" border stripe>
                  <el-table-column prop="param" label="参数名" width="150" />
                  <el-table-column prop="type" label="类型" width="100" />
                  <el-table-column prop="required" label="必填" width="80" />
                  <el-table-column prop="description" label="说明" />
                </el-table>
                
                <h3>请求示例</h3>
                <div class="code-block" id="results-example">
                  <pre><code>GET /v1/results/rec_12345
X-App-ID: your_app_id
X-API-Key: your_api_key</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('results-example')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
                
                <h3>响应格式</h3>
                <div class="code-block" id="results-response">
                  <pre><code>{
  "success": true,
  "data": {
    "id": "rec_12345",
    "status": "completed",
    "results": [
      {
        "label": "动物",
        "confidence": 0.94,
        "bbox": [20, 30, 120, 150]
      }
    ],
    "created_at": "2023-06-15T10:30:00Z",
    "processing_time": 0.28
  }
}</code></pre>
                  <el-button size="small" type="primary" class="copy-btn" @click="copyCode('results-response')">
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </div>
              </el-collapse-item>
            </el-collapse>
          </section>
          
          <el-divider></el-divider>
          
          <section id="error-codes" class="doc-section">
            <h2>错误码</h2>
            <el-table :data="errorCodes" style="width: 100%" border stripe>
              <el-table-column prop="code" label="错误码" width="100" />
              <el-table-column prop="message" label="错误信息" width="200" />
              <el-table-column prop="description" label="说明" />
            </el-table>
          </section>
          
          <el-divider></el-divider>
          
          <section id="limits" class="doc-section">
            <h2>使用限制</h2>
            <el-alert
              title="请注意以下使用限制"
              type="info"
              :closable="false"
              show-icon
            >
              <ul class="limit-list">
                <li>每个API密钥的调用频率限制为每分钟60次请求</li>
                <li>单张图片大小不超过10MB</li>
                <li>批量识别接口单次请求最多支持20张图片</li>
                <li>识别结果保存时间为7天</li>
              </ul>
            </el-alert>
          </section>
          
          <el-divider></el-divider>
          
          <section id="sdk" class="doc-section">
            <h2>SDK和代码示例</h2>
            <el-tabs type="border-card">
              <el-tab-pane label="JavaScript">
                <div class="code-block" id="js-example">
                  <pre><code>// 安装依赖
// npm install axios

const axios = require('axios');

async function recognizeImage(imageUrl) {
  try {
    const response = await axios.post('https://api.example.com/v1/recognize', {
      image_url: imageUrl,
      min_confidence: 0.7
    }, {
      headers: {
        'Content-Type': 'application/json',
        'X-App-ID': 'your_app_id',
        'X-API-Key': 'your_api_key'
      }
    });
    
    console.log('识别结果:', response.data);
    return response.data;
  } catch (error) {
    console.error('识别失败:', error.response ? error.response.data : error.message);
    throw error;
  }
}</code></pre>
              </div>
            </el-tab-pane>
            <el-tab-pane label="Python">
              <div class="code-block">
                <pre><code># 安装依赖
# pip install requests

import requests
import json

def recognize_image(image_url):
    url = "https://api.example.com/v1/recognize"
    headers = {
        "Content-Type": "application/json",
        "X-App-ID": "your_app_id",
        "X-API-Key": "your_api_key"
    }
    payload = {
        "image_url": image_url,
        "min_confidence": 0.7
    }
    
    try:
        response = requests.post(url, headers=headers, data=json.dumps(payload))
        response.raise_for_status()
        result = response.json()
        print("识别结果:", result)
        return result
    except requests.exceptions.RequestException as e:
        print("识别失败:", str(e))
        raise</code></pre>
              </div>
            </el-tab-pane>
          </el-tabs>
          </section>
          
          <el-divider></el-divider>
          
          <h2>API 在线测试工具</h2>
          <div class="api-tester" id="api-tester">
            <el-form :model="apiTestForm" label-position="top">
              <el-form-item label="选择API端点">
                <el-select v-model="apiTestForm.endpoint" placeholder="选择要测试的API端点">
                  <el-option label="图像识别 /recognize" value="/recognize" />
                  <el-option label="批量识别 /batch-recognize" value="/batch-recognize" />
                  <el-option label="查询结果 /results/{id}" value="/results" />
                </el-select>
              </el-form-item>
              
              <template v-if="apiTestForm.endpoint === '/recognize'">
                <el-form-item label="图片来源">
                  <el-radio-group v-model="apiTestForm.imageSource">
                    <el-radio label="url">图片URL</el-radio>
                    <el-radio label="upload">上传图片</el-radio>
                  </el-radio-group>
                </el-form-item>
                
                <el-form-item label="图片URL" v-if="apiTestForm.imageSource === 'url'">
                  <el-input v-model="apiTestForm.imageUrl" placeholder="请输入图片URL" />
                </el-form-item>
                
                <el-form-item label="上传图片" v-if="apiTestForm.imageSource === 'upload'">
                  <el-upload
                    class="upload-demo"
                    action="#"
                    :auto-upload="false"
                    :limit="1"
                    list-type="picture"
                  >
                    <template #trigger>
                      <el-button type="primary">选择图片</el-button>
                    </template>
                    <template #tip>
                      <div class="el-upload__tip">仅支持jpg/png文件，且不超过10MB</div>
                    </template>
                  </el-upload>
                </el-form-item>
                
                <el-form-item label="最小置信度">
                  <el-slider v-model="apiTestForm.minConfidence" :min="0" :max="1" :step="0.05" show-stops />
                </el-form-item>
                
                <el-form-item label="最大结果数">
                  <el-input-number v-model="apiTestForm.maxResults" :min="1" :max="50" />
                </el-form-item>
              </template>
              
              <template v-if="apiTestForm.endpoint === '/results'">
                <el-form-item label="识别结果ID">
                  <el-input v-model="apiTestForm.resultId" placeholder="请输入识别结果ID" />
                </el-form-item>
              </template>
              
              <el-form-item>
                <el-button type="primary" @click="testApi">测试API</el-button>
                <el-button @click="resetForm">重置</el-button>
              </el-form-item>
            </el-form>
            
            <div v-if="apiTestResult" class="api-test-result">
              <h3>测试结果</h3>
              <div class="code-block">
                <pre><code>{{ JSON.stringify(apiTestResult, null, 2) }}</code></pre>
                <el-button size="small" type="primary" icon="CopyDocument" @click="copyTestResult">复制</el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { CopyDocument, Back, Connection, View, Hide } from '@element-plus/icons-vue'

const router = useRouter()

// 返回API设置页面
const goBack = () => {
  router.push('/client/api-settings')
}

// 滚动到指定部分
const scrollToSection = (sectionId) => {
  const section = document.getElementById(sectionId)
  if (section) {
    section.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

// 代码复制功能
const copyCode = (id) => {
  const codeElement = document.getElementById(id)
  if (codeElement) {
    const text = codeElement.textContent
    navigator.clipboard.writeText(text)
      .then(() => {
        ElMessage({
          message: '代码已复制到剪贴板',
          type: 'success',
          duration: 2000
        })
      })
      .catch(() => {
        ElMessage({
          message: '复制失败，请手动复制',
          type: 'error',
          duration: 2000
        })
      })
  }
}

// API测试功能
const showApiTester = ref(true)
const apiTestForm = reactive({
  endpoint: '/recognize',
  imageSource: 'url',
  imageUrl: '',
  resultId: '',
  minConfidence: 0.7,
  maxResults: 10
})

// API测试结果
const apiTestResult = ref(null)

// 测试API功能
const testApi = () => {
  // 表单验证
  if (apiTestForm.endpoint === '/recognize' && apiTestForm.imageSource === 'url' && !apiTestForm.imageUrl) {
    ElMessage.warning('请输入图片URL')
    return
  }
  
  if (apiTestForm.endpoint === '/results' && !apiTestForm.resultId) {
    ElMessage.warning('请输入识别结果ID')
    return
  }
  
  // 模拟API调用
  const loading = ElMessage.loading({
    message: '正在处理请求...',
    duration: 0
  })
  
  setTimeout(() => {
    loading.close()
    
    if (apiTestForm.endpoint === '/recognize') {
      apiTestResult.value = {
        success: true,
        data: {
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
            }
          ],
          processing_time: 0.32
        }
      }
    } else if (apiTestForm.endpoint === '/batch-recognize') {
      apiTestResult.value = {
        success: true,
        data: {
          batch_id: 'batch_' + Math.floor(Math.random() * 100000),
          results: [
            {
              image_id: 'img1',
              status: 'success',
              results: [
                {
                  label: '汽车',
                  confidence: 0.92,
                  bbox: [10, 20, 300, 200]
                }
              ]
            },
            {
              image_id: 'img2',
              status: 'success',
              results: [
                {
                  label: '树木',
                  confidence: 0.88,
                  bbox: [150, 50, 250, 350]
                }
              ]
            }
          ],
          processing_time: 0.75
        }
      }
    } else if (apiTestForm.endpoint === '/results') {
      apiTestResult.value = {
        success: true,
        data: {
          id: apiTestForm.resultId || 'rec_12345',
          status: 'completed',
          results: [
            {
              label: '动物',
              confidence: 0.94,
              bbox: [20, 30, 120, 150]
            }
          ],
          created_at: new Date().toISOString(),
          processing_time: 0.28
        }
      }
    }
    
    ElMessage({
      message: 'API测试成功',
      type: 'success',
      duration: 2000
    })
    
    // 滚动到结果区域
    setTimeout(() => {
      const resultElement = document.querySelector('.api-test-result')
      if (resultElement) {
        resultElement.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }, 100)
  }, 1500)
}

// 复制测试结果
const copyTestResult = () => {
  if (apiTestResult.value) {
    navigator.clipboard.writeText(JSON.stringify(apiTestResult.value, null, 2))
      .then(() => {
        ElMessage({
          message: '测试结果已复制到剪贴板',
          type: 'success',
          duration: 2000
        })
      })
      .catch(() => {
        ElMessage({
          message: '复制失败，请手动复制',
          type: 'error',
          duration: 2000
        })
      })
  }
}

// 重置表单
const resetForm = () => {
  apiTestForm.imageUrl = ''
  apiTestForm.resultId = ''
  apiTestForm.minConfidence = 0.7
  apiTestForm.maxResults = 10
  apiTestResult.value = null
  ElMessage({
    message: '表单已重置',
    type: 'info',
    duration: 2000
  })
}

// 图像识别API参数
const recognizeParams = [
  {
    param: 'image_url',
    type: 'String',
    required: '否*',
    description: '图片的URL地址（*image_url和image_data必须提供一个）'
  },
  {
    param: 'image_data',
    type: 'String',
    required: '否*',
    description: 'Base64编码的图片数据（*image_url和image_data必须提供一个）'
  },
  {
    param: 'min_confidence',
    type: 'Number',
    required: '否',
    description: '最小置信度阈值，范围0-1，默认为0.5'
  },
  {
    param: 'max_results',
    type: 'Integer',
    required: '否',
    description: '返回结果的最大数量，默认为10'
  }
]

// 批量识别API参数
const batchRecognizeParams = [
  {
    param: 'images',
    type: 'Array',
    required: '是',
    description: '图片数组，每个元素包含id和url字段'
  },
  {
    param: 'min_confidence',
    type: 'Number',
    required: '否',
    description: '最小置信度阈值，范围0-1，默认为0.5'
  }
]

// 查询结果API路径参数
const resultPathParams = [
  {
    param: 'id',
    type: 'String',
    required: '是',
    description: '识别任务的唯一标识符'
  }
]

// 错误码表
const errorCodes = [
  {
    code: '400',
    message: '请求参数错误',
    description: '请检查请求参数格式是否正确'
  },
  {
    code: '401',
    message: '认证失败',
    description: '请检查API密钥是否正确'
  },
  {
    code: '403',
    message: '权限不足',
    description: '当前API密钥没有权限访问该接口'
  },
  {
    code: '404',
    message: '资源不存在',
    description: '请求的资源不存在'
  },
  {
    code: '429',
    message: '请求过于频繁',
    description: '已超过API调用频率限制'
  },
  {
    code: '500',
    message: '服务器错误',
    description: '服务器内部错误，请稍后重试'
  }
]

// 页面加载完成后的处理
onMounted(() => {
  // 检查URL中是否有锚点，如果有则滚动到对应部分
  const hash = window.location.hash
  if (hash) {
    const sectionId = hash.substring(1)
    setTimeout(() => {
      scrollToSection(sectionId)
    }, 300)
  }
})
</script>

<style scoped>
.api-documentation-container {
  padding: 20px;
  max-width: 1400px;
  margin: 0 auto;
}

.doc-card {
  margin-bottom: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.doc-content {
  padding: 20px 0;
  display: flex;
  gap: 30px;
}

.doc-navigation {
  width: 250px;
  flex-shrink: 0;
}

.nav-card {
  padding: 10px 0;
}

.nav-card h3 {
  padding: 0 20px;
  margin-top: 0;
}

.doc-main-content {
  flex: 1;
  min-width: 0;
}

.doc-section {
  margin-bottom: 30px;
  scroll-margin-top: 100px;
}

.code-block {
  background-color: #282c34;
  color: #abb2bf;
  padding: 15px;
  border-radius: 6px;
  margin: 15px 0;
  position: relative;
  overflow: auto;
}

.code-block pre {
  margin: 0;
  white-space: pre-wrap;
  font-family: 'Courier New', Courier, monospace;
  font-size: 14px;
}

.code-block code {
  font-family: 'Courier New', Courier, monospace;
}

.copy-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  opacity: 0.7;
  transition: opacity 0.2s;
}

.copy-btn:hover {
  opacity: 1;
}

.limit-list {
  padding-left: 20px;
}

.limit-list li {
  margin-bottom: 10px;
}

/* API测试工具样式 */
.api-tester {
  background-color: #f9fafc;
  border-radius: 8px;
  padding: 20px;
  margin-top: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.api-test-result {
  margin-top: 20px;
  border-top: 1px solid #ebeef5;
  padding-top: 20px;
}

.api-test-result h3 {
  margin-bottom: 15px;
  color: #409EFF;
}

.api-test-result .code-block {
  position: relative;
}

.api-test-result .el-button {
  position: absolute;
  top: 10px;
  right: 10px;
}

/* 响应式布局 */
@media (max-width: 992px) {
  .doc-content {
    flex-direction: column;
  }
  
  .doc-navigation {
    width: 100%;
    margin-bottom: 20px;
  }
  
  .nav-card {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .api-documentation-container {
    padding: 10px;
  }
  
  .doc-content {
    padding: 10px 0;
  }
  
  .code-block {
    padding: 10px;
    font-size: 13px;
  }
  
  .api-tester {
    padding: 15px;
  }
  
  .el-table {
    font-size: 13px;
  }
  
  .el-collapse-item__header {
    font-size: 14px;
  }
}
</style>