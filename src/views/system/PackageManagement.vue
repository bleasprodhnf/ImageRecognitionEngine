<template>
  <div class="package-management">
    <div class="action-bar">
      <a-space>
        <a-button type="primary" @click="handleAddPackage">新增套餐</a-button>
        <a-button @click="showCompareModal = true" :disabled="selectedPackages.length < 2">
          套餐对比
        </a-button>
      </a-space>
    </div>

    <div class="package-list">
      <a-row :gutter="24">
        <a-col :span="8" v-for="pkg in packages" :key="pkg.id">
          <a-card class="package-card" :class="{ selected: selectedPackages.includes(pkg.id) }">
            <template #extra>
              <a-checkbox
                :checked="selectedPackages.includes(pkg.id)"
                @change="(e) => handleSelectPackage(pkg.id, e.target.checked)"
              />
            </template>
            <template #title>
              <div class="package-title">
                <span>{{ pkg.name }}</span>
                <a-tag :color="pkg.status === 'active' ? 'success' : 'default'">
                  {{ pkg.status === 'active' ? '已启用' : '已禁用' }}
                </a-tag>
              </div>
            </template>
            <div class="package-content">
              <div class="price">¥{{ pkg.price }}/月</div>
              <div class="features">
                <div v-for="(feature, index) in pkg.features" :key="index" class="feature-item">
                  <CheckOutlined /> {{ feature }}
                </div>
              </div>
              <div class="actions">
                <a-space>
                  <a-button type="link" @click="handleEditPackage(pkg)">编辑</a-button>
                  <a-button
                    type="link"
                    @click="handleToggleStatus(pkg)"
                  >
                    {{ pkg.status === 'active' ? '禁用' : '启用' }}
                  </a-button>
                  <a-popconfirm
                    title="确定要删除该套餐吗？"
                    @confirm="handleDeletePackage(pkg)"
                  >
                    <a-button type="link" danger>删除</a-button>
                  </a-popconfirm>
                </a-space>
              </div>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <a-modal
      v-model:visible="modalVisible"
      :title="modalType === 'add' ? '新增套餐' : '编辑套餐'"
      @ok="handleModalOk"
      width="600px"
    >
      <a-form :model="packageForm" :rules="rules" ref="formRef" layout="vertical">
        <a-form-item label="套餐名称" name="name">
          <a-input v-model:value="packageForm.name" placeholder="请输入套餐名称" />
        </a-form-item>
        <a-form-item label="套餐价格" name="price">
          <a-input-number
            v-model:value="packageForm.price"
            :min="0"
            :precision="2"
            style="width: 100%"
            placeholder="请输入套餐价格"
          />
        </a-form-item>
        <a-form-item label="功能特点" name="features">
          <a-space direction="vertical" style="width: 100%">
            <div v-for="(feature, index) in packageForm.features" :key="index" class="feature-input">
              <a-input v-model:value="packageForm.features[index]" placeholder="请输入功能特点" />
              <MinusCircleOutlined @click="removeFeature(index)" class="remove-icon" />
            </div>
            <a-button type="dashed" block @click="addFeature">
              <PlusOutlined /> 添加功能特点
            </a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="showCompareModal"
      title="套餐对比"
      width="80%"
      :footer="null"
    >
      <div class="compare-table">
        <a-table
          :columns="compareColumns"
          :data-source="getCompareData()"
          :pagination="false"
          bordered
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { CheckOutlined, PlusOutlined, MinusCircleOutlined } from '@ant-design/icons-vue'

const packages = ref([])
const selectedPackages = ref([])
const modalVisible = ref(false)
const modalType = ref('add')
const showCompareModal = ref(false)
const formRef = ref()

const packageForm = reactive({
  name: '',
  price: 0,
  features: ['']
})

const rules = {
  name: [{ required: true, message: '请输入套餐名称', trigger: 'blur' }],
  price: [{ required: true, message: '请输入套餐价格', trigger: 'blur' }],
  features: [{ required: true, message: '请至少添加一个功能特点', trigger: 'change' }]
}

const compareColumns = [
  {
    title: '对比项',
    dataIndex: 'item',
    key: 'item',
    fixed: 'left',
    width: 150
  }
]

const handleSelectPackage = (id, checked) => {
  if (checked) {
    if (selectedPackages.value.length >= 3) {
      message.warning('最多只能同时对比3个套餐')
      return
    }
    selectedPackages.value.push(id)
  } else {
    selectedPackages.value = selectedPackages.value.filter(pid => pid !== id)
  }
}

const handleAddPackage = () => {
  modalType.value = 'add'
  modalVisible.value = true
  packageForm.name = ''
  packageForm.price = 0
  packageForm.features = ['']
}

const handleEditPackage = (pkg) => {
  modalType.value = 'edit'
  modalVisible.value = true
  packageForm.name = pkg.name
  packageForm.price = pkg.price
  packageForm.features = [...pkg.features]
}

const handleDeletePackage = (pkg) => {
  const index = packages.value.findIndex(item => item.id === pkg.id)
  if (index !== -1) {
    packages.value.splice(index, 1)
    message.success('删除成功')
  }
}

const handleToggleStatus = (pkg) => {
  pkg.status = pkg.status === 'active' ? 'inactive' : 'active'
  message.success(`${pkg.status === 'active' ? '启用' : '禁用'}成功`)
}

const handleModalOk = () => {
  formRef.value.validate().then(() => {
    if (modalType.value === 'add') {
      packages.value.push({
        id: packages.value.length + 1,
        name: packageForm.name,
        price: packageForm.price,
        features: [...packageForm.features],
        status: 'active'
      })
      message.success('添加成功')
    } else {
      const pkg = packages.value.find(item => item.name === packageForm.name)
      if (pkg) {
        pkg.price = packageForm.price
        pkg.features = [...packageForm.features]
        message.success('更新成功')
      }
    }
    modalVisible.value = false
  })
}

const addFeature = () => {
  packageForm.features.push('')
}

const removeFeature = (index) => {
  packageForm.features.splice(index, 1)
}

const getCompareData = () => {
  const selectedPkgs = packages.value.filter(pkg => selectedPackages.value.includes(pkg.id))
  selectedPkgs.forEach((pkg, index) => {
    compareColumns[index + 1] = {
      title: pkg.name,
      dataIndex: `pkg${index}`,
      key: `pkg${index}`,
      align: 'center'
    }
  })
  
  const data = [
    {
      item: '价格（元/月）',
      ...selectedPkgs.reduce((acc, pkg, index) => ({ ...acc, [`pkg${index}`]: pkg.price }), {})
    },
    {
      item: '功能特点',
      ...selectedPkgs.reduce((acc, pkg, index) => ({ ...acc, [`pkg${index}`]: pkg.features.join('\n') }), {})
    },
    {
      item: '状态',
      ...selectedPkgs.reduce((acc, pkg, index) => ({ ...acc, [`pkg${index}`]: pkg.status === 'active' ? '已启用' : '已禁用' }), {})
    }
  ]
  
  return data
}
</script>