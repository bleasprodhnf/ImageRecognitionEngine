import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { modelVersions, modelParams, modelMonitorData } from '../api/model'

export const useModelStore = defineStore('model', () => {
  // 状态
  const versions = ref([])
  const params = ref([])
  const monitorData = ref({})
  const loading = ref(false)
  const error = ref(null)
  const currentVersion = ref(null)

  // 计算属性
  const productionVersions = computed(() => 
    versions.value.filter(v => v.status === 'production')
  )
  
  const testingVersions = computed(() => 
    versions.value.filter(v => v.status === 'testing')
  )
  
  const developmentVersions = computed(() => 
    versions.value.filter(v => v.status === 'development')
  )

  // 方法
  async function fetchVersions() {
    loading.value = true
    error.value = null
    
    try {
      // 实际项目中应该调用API获取数据
      // const response = await api.getModelVersions()
      
      // 使用模拟数据
      versions.value = modelVersions
      
      if (versions.value.length > 0 && !currentVersion.value) {
        currentVersion.value = versions.value[0]
      }
      
      return versions.value
    } catch (err) {
      error.value = err.message || '获取模型版本失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchParams() {
    loading.value = true
    error.value = null
    
    try {
      // 实际项目中应该调用API获取数据
      // const response = await api.getModelParams()
      
      // 使用模拟数据
      params.value = modelParams
      
      return params.value
    } catch (err) {
      error.value = err.message || '获取模型参数失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function fetchMonitorData() {
    loading.value = true
    error.value = null
    
    try {
      // 实际项目中应该调用API获取数据
      // const response = await api.getModelMonitorData()
      
      // 使用模拟数据
      monitorData.value = modelMonitorData
      
      return monitorData.value
    } catch (err) {
      error.value = err.message || '获取监控数据失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function updateModelParam(id, newValue) {
    loading.value = true
    error.value = null
    
    try {
      // 实际项目中应该调用API更新数据
      // const response = await api.updateModelParam(id, newValue)
      
      // 模拟更新
      const paramIndex = params.value.findIndex(p => p.id === id)
      if (paramIndex !== -1) {
        params.value[paramIndex].value = newValue
      }
      
      return params.value[paramIndex]
    } catch (err) {
      error.value = err.message || '更新模型参数失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  function setCurrentVersion(version) {
    currentVersion.value = version
  }

  return {
    // 状态
    versions,
    params,
    monitorData,
    loading,
    error,
    currentVersion,
    
    // 计算属性
    productionVersions,
    testingVersions,
    developmentVersions,
    
    // 方法
    fetchVersions,
    fetchParams,
    fetchMonitorData,
    updateModelParam,
    setCurrentVersion
  }
})
