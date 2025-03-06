// 模拟接口数据

// 用户信息
export const userInfo = {
  username: 'admin',
  realname: '系统管理员',
  email: 'admin@example.com',
  phone: '13800138000',
  avatar: '/src/assets/images/avatars/default-avatar.svg',
  roles: ['super_admin'],
  permissions: ['*']
}

// 修改密码响应
export const changePasswordResponse = {
  success: true,
  message: '密码修改成功'
}

// 修改个人信息响应
export const updateProfileResponse = {
  success: true,
  message: '个人信息更新成功'
}

// 密码验证规则
export const passwordRules = {
  minLength: 8,
  maxLength: 20,
  requireUppercase: true,
  requireLowercase: true,
  requireNumbers: true,
  requireSpecialChars: true,
  specialChars: '!@#$%^&*()'
}

// 管理员列表数据
export const admins = [
  {
    id: 1,
    username: 'admin',
    role: '超级管理员',
    avatar: '/src/assets/images/avatars/default-avatar.svg',
    createTime: '2024-01-01 00:00:00',
    lastLogin: '2024-01-15 08:30:00'
  },
  {
    id: 2,
    username: 'operator',
    role: '运维管理员',
    avatar: '/src/assets/images/avatars/operator-avatar.svg',
    createTime: '2024-01-02 00:00:00',
    lastLogin: '2024-01-15 09:15:00'
  }
]

// 套餐配置数据
export const packages = [
  {
    id: 1,
    name: '企业版',
    price: 9999,
    description: '适合大型企业使用的高级套餐',
    features: [
      '所有基础功能',
      '高级数据分析',
      '专属技术支持',
      '自定义模型训练'
    ],
    limits: {
      storage: 1000,
      concurrency: 100,
      models: 10
    }
  },
  {
    id: 2,
    name: '专业版',
    price: 4999,
    description: '适合中型企业使用的专业套餐',
    features: [
      '所有基础功能',
      '数据分析',
      '技术支持',
      '模型训练'
    ],
    limits: {
      storage: 500,
      concurrency: 50,
      models: 5
    }
  },
  {
    id: 3,
    name: '基础版',
    price: 1999,
    description: '适合小型企业使用的基础套餐',
    features: [
      '基础功能',
      '基础分析',
      '邮件支持'
    ],
    limits: {
      storage: 100,
      concurrency: 20,
      models: 2
    }
  }
]

// 客户列表数据
export const customers = [
  {
    id: 1,
    name: '示例科技有限公司',
    package: '企业版',
    username: 'example_corp',
    password: 'encrypted_password',
    address: '北京市朝阳区科技园A座',
    contact: '张经理',
    phone: '13911112222',
    status: 'active',
    expireDate: '2024-12-31',
    usage: {
      api: 15000,
      storage: 850
    }
  },
  {
    id: 2,
    name: '创新软件科技',
    package: '专业版',
    username: 'innovation_soft',
    password: 'encrypted_password',
    address: '上海市浦东新区创新路88号',
    contact: '李工',
    phone: '13822223333',
    status: 'active',
    expireDate: '2024-06-30',
    usage: {
      api: 8000,
      storage: 320
    }
  }
]

// 系统参数配置
export const systemParams = {
  siteName: '图像识别引擎管理系统',
  storageConfig: {
    maxFileSize: 10, // MB
    allowedTypes: ['jpg', 'jpeg', 'png']
  },
  notificationConfig: {
    email: true,
    sms: false
  },
  securityConfig: {
    passwordMinLength: 8,
    loginAttempts: 5,
    lockDuration: 30 // minutes
  }
}

// 系统监控数据
export const systemMonitor = {
  cpu: {
    usage: 45,
    cores: 8,
    temperature: 65
  },
  memory: {
    total: 16384,
    used: 8192,
    free: 8192
  },
  disk: {
    total: 1024,
    used: 512,
    free: 512
  },
  network: {
    upload: 2.5,
    download: 5.8,
    connections: 126
  }
}

// 角色列表
export const roles = [
  {
    id: 1,
    name: '超级管理员',
    code: 'super_admin',
    permissions: ['*']
  },
  {
    id: 2,
    name: '运维管理员',
    code: 'ops_admin',
    permissions: ['system:monitor', 'system:logs']
  },
  {
    id: 3,
    name: '客服管理员',
    code: 'service_admin',
    permissions: ['customer:view', 'customer:service']
  }
]

// 管理员列表
export const adminList = [
  {
    id: 1,
    username: 'admin',
    realname: '系统管理员',
    role: 'super_admin',
    status: 1,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    lastLogin: '2024-01-20 10:30:00'
  },
  {
    id: 2,
    username: 'ops',
    realname: '运维人员',
    role: 'ops_admin',
    status: 1,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    lastLogin: '2024-01-20 09:15:00'
  },
  {
    id: 3,
    username: 'service',
    realname: '客服人员',
    role: 'service_admin',
    status: 1,
    avatar: 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png',
    lastLogin: '2024-01-20 08:45:00'
  }
]

// 系统日志
export const systemLogs = [
  {
    id: 1,
    type: 'operation',
    user: 'admin',
    action: '修改系统配置',
    ip: '192.168.1.100',
    time: '2024-01-20 10:30:00',
    status: 'success'
  },
  {
    id: 2,
    type: 'security',
    user: 'ops',
    action: '登录系统',
    ip: '192.168.1.101',
    time: '2024-01-20 09:15:00',
    status: 'success'
  },
  {
    id: 3,
    type: 'error',
    user: 'service',
    action: '导出报表',
    ip: '192.168.1.102',
    time: '2024-01-20 08:45:00',
    status: 'failed'
  }
]

// 模型版本列表
export const modelVersions = [
  {
    id: 1,
    version: 'v1.0.0',
    description: '初始版本',
    accuracy: 0.85,
    status: 'production',
    createTime: '2024-01-01 00:00:00'
  },
  {
    id: 2,
    version: 'v1.1.0',
    description: '优化识别准确率',
    accuracy: 0.89,
    status: 'testing',
    createTime: '2024-01-15 00:00:00'
  },
  {
    id: 3,
    version: 'v1.2.0-beta',
    description: '新增特征提取',
    accuracy: 0.92,
    status: 'development',
    createTime: '2024-01-20 00:00:00'
  }
]

// 客户账号列表
export const customerAccounts = [
  {
    id: 1,
    name: '示例公司A',
    package: '企业版',
    status: 'active',
    address: '北京市朝阳区科技园区88号',
    contact: '张经理',
    phone: '13911112222',
    expireDate: '2024-12-31',
    usage: {
      api: 8500,
      storage: 128
    }
  },
  {
    id: 2,
    name: '示例公司B',
    package: '专业版',
    status: 'active',
    address: '上海市浦东新区陆家嘴1号',
    contact: '李总监',
    phone: '13822223333',
    expireDate: '2024-06-30',
    usage: {
      api: 3200,
      storage: 64
    }
  },
  {
    id: 3,
    name: '示例公司C',
    package: '基础版',
    status: 'inactive',
    address: '广州市天河区珠江新城66号',
    contact: '王经理',
    phone: '13733334444',
    expireDate: '2024-01-31',
    usage: {
      api: 500,
      storage: 16
    }
  }
]