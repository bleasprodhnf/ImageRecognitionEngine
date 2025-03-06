<template>
  <div class="role-management">
    <div class="action-bar">
      <a-button type="primary" @click="handleAddRole">新增角色</a-button>
    </div>

    <div class="role-list">
      <a-row :gutter="24">
        <a-col :span="6" v-for="role in roles" :key="role.id">
          <a-card class="role-card">
            <template #title>
              <div class="role-title">
                <span>{{ role.name }}</span>
                <a-tag :color="role.isSystem ? 'blue' : 'green'">
                  {{ role.isSystem ? '系统角色' : '自定义角色' }}
                </a-tag>
              </div>
            </template>
            <div class="role-content">
              <p class="role-desc">{{ role.description }}</p>
              <div class="role-info">
                <div>用户数量：{{ role.userCount }}</div>
                <div>创建时间：{{ role.createTime }}</div>
              </div>
              <div class="actions">
                <a-space>
                  <a-button type="link" @click="handleEditRole(role)" :disabled="role.isSystem">
                    编辑
                  </a-button>
                  <a-button type="link" @click="handleConfigPermissions(role)">
                    权限配置
                  </a-button>
                  <a-popconfirm
                    v-if="!role.isSystem"
                    title="确定要删除该角色吗？"
                    @confirm="handleDeleteRole(role)"
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
      v-model:visible="roleModalVisible"
      :title="modalType === 'add' ? '新增角色' : '编辑角色'"
      @ok="handleRoleModalOk"
    >
      <a-form :model="roleForm" :rules="rules" ref="formRef" layout="vertical">
        <a-form-item label="角色名称" name="name">
          <a-input v-model:value="roleForm.name" placeholder="请输入角色名称" />
        </a-form-item>
        <a-form-item label="角色描述" name="description">
          <a-textarea
            v-model:value="roleForm.description"
            placeholder="请输入角色描述"
            :rows="4"
          />
        </a-form-item>
      </a-form>
    </a-modal>

    <a-modal
      v-model:visible="permissionModalVisible"
      title="权限配置"
      width="800px"
      @ok="handlePermissionModalOk"
    >
      <div class="permission-tree">
        <a-tree
          v-model:checkedKeys="checkedPermissions"
          :treeData="permissionTree"
          checkable
          :defaultExpandAll="true"
        >
          <template #title="{ title, type }">
            <span>
              <FileOutlined v-if="type === 'menu'" />
              <ToolOutlined v-else-if="type === 'action'" />
              {{ title }}
            </span>
          </template>
        </a-tree>
      </div>
      <div class="permission-stats">
        <a-alert
          type="info"
          :message="`已选择 ${checkedPermissions.length} 项权限`"
          style="margin-top: 16px"
        />
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { message } from 'ant-design-vue'
import { FileOutlined, ToolOutlined } from '@ant-design/icons-vue'

const roles = ref([])
const roleModalVisible = ref(false)
const permissionModalVisible = ref(false)
const modalType = ref('add')
const formRef = ref()

const roleForm = reactive({
  name: '',
  description: ''
})

const rules = {
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
  description: [{ required: true, message: '请输入角色描述', trigger: 'blur' }]
}

const checkedPermissions = ref([])
const permissionTree = ref([
  {
    title: '系统管理',
    key: 'system',
    type: 'menu',
    children: [
      {
        title: '用户管理',
        key: 'system:user',
        type: 'menu',
        children: [
          {
            title: '查看用户',
            key: 'system:user:view',
            type: 'action'
          },
          {
            title: '编辑用户',
            key: 'system:user:edit',
            type: 'action'
          }
        ]
      },
      {
        title: '角色管理',
        key: 'system:role',
        type: 'menu',
        children: [
          {
            title: '查看角色',
            key: 'system:role:view',
            type: 'action'
          },
          {
            title: '编辑角色',
            key: 'system:role:edit',
            type: 'action'
          }
        ]
      }
    ]
  }
])

const handleAddRole = () => {
  modalType.value = 'add'
  roleModalVisible.value = true
  roleForm.name = ''
  roleForm.description = ''
}

const handleEditRole = (role) => {
  modalType.value = 'edit'
  roleModalVisible.value = true
  roleForm.name = role.name
  roleForm.description = role.description
}

const handleRoleModalOk = async () => {
  try {
    await formRef.value.validate()
    // TODO: 实现角色保存接口
    message.success(modalType.value === 'add' ? '角色创建成功' : '角色更新成功')
    roleModalVisible.value = false
  } catch (error) {
    // 表单验证失败
  }
}

const handleConfigPermissions = (role) => {
  permissionModalVisible.value = true
  // TODO: 加载角色当前权限
  checkedPermissions.value = []
}

const handlePermissionModalOk = async () => {
  try {
    // TODO: 实现权限保存接口
    message.success('权限配置已更新')
    permissionModalVisible.value = false
  } catch (error) {
    message.error('权限配置保存失败')
  }
}

const handleDeleteRole = async (role) => {
  try {
    // TODO: 实现角色删除接口
    message.success('角色删除成功')
  } catch (error) {
    message.error('角色删除失败')
  }
}
</script>

<style scoped>
.role-management {
  padding: 24px;
}

.action-bar {
  margin-bottom: 24px;
}

.role-card {
  margin-bottom: 24px;
}

.role-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.role-content {
  .role-desc {
    margin-bottom: 16px;
    color: rgba(0, 0, 0, 0.45);
  }

  .role-info {
    margin-bottom: 16px;
    font-size: 14px;
    color: rgba(0, 0, 0, 0.65);
  }
}

.permission-tree {
  max-height: 400px;
  overflow-y: auto;
}
</style>