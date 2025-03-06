// 客户管理相关的模拟数据接口

// 套餐配置列表
export const packages = [
  {
    id: 1,
    name: '企业版',
    price: 9999,
    description: '适用于大型企业的高级套餐',
    features: [
      '无限API调用次数',
      '优先级技术支持',
      '自定义模型训练',
      '高级数据分析'
    ],
    limits: {
      storage: 1024, // GB
      concurrency: 100,
      models: 10
    }
  },
  {
    id: 2,
    name: '专业版',
    price: 4999,
    description: '适用于中型企业的专业套餐',
    features: [
      '月度100万次API调用',
      '标准技术支持',
      '基础模型训练',
      '基础数据分析'
    ],
    limits: {
      storage: 512, // GB
      concurrency: 50,
      models: 5
    }
  },
  {
    id: 3,
    name: '基础版',
    price: 1999,
    description: '适用于小型企业的基础套餐',
    features: [
      '月度10万次API调用',
      '邮件技术支持',
      '使用预训练模型',
      '基础数据统计'
    ],
    limits: {
      storage: 128, // GB
      concurrency: 20,
      models: 2
    }
  }
]

// 客户服务监控数据
export const serviceMonitor = {
  // 实时状态
  status: {
    activeUsers: 126,
    apiCalls: 1580,
    errorRate: 0.5,
    avgResponseTime: 150 // ms
  },
  // 今日统计
  today: {
    totalCalls: 15800,
    successCalls: 15700,
    failedCalls: 100,
    totalUsers: 280
  },
  // 异常记录
  errors: [
    {
      id: 1,
      customer: '示例公司A',
      type: 'API调用超时',
      time: '2024-01-20 10:15:30',
      status: 'resolved'
    },
    {
      id: 2,
      customer: '示例公司B',
      type: '并发超限',
      time: '2024-01-20 09:45:20',
      status: 'pending'
    },
    {
      id: 3,
      customer: '示例公司C',
      type: '存储空间不足',
      time: '2024-01-20 09:30:15',
      status: 'processing'
    }
  ]
}

// 资源使用统计
export const resourceStats = {
  // API调用统计
  apiUsage: [
    { date: '2024-01-14', calls: 95000 },
    { date: '2024-01-15', calls: 98000 },
    { date: '2024-01-16', calls: 120000 },
    { date: '2024-01-17', calls: 115000 },
    { date: '2024-01-18', calls: 125000 },
    { date: '2024-01-19', calls: 110000 },
    { date: '2024-01-20', calls: 105000 }
  ],
  // 存储使用统计
  storageUsage: [
    { date: '2024-01-14', usage: 850 },
    { date: '2024-01-15', usage: 875 },
    { date: '2024-01-16', usage: 890 },
    { date: '2024-01-17', usage: 920 },
    { date: '2024-01-18', usage: 950 },
    { date: '2024-01-19', usage: 980 },
    { date: '2024-01-20', usage: 1024 }
  ],
  // 客户资源使用排名
  topCustomers: [
    {
      name: '示例公司A',
      apiCalls: 450000,
      storage: 512,
      cost: 9999
    },
    {
      name: '示例公司B',
      apiCalls: 280000,
      storage: 256,
      cost: 4999
    },
    {
      name: '示例公司C',
      apiCalls: 95000,
      storage: 128,
      cost: 1999
    }
  ]
}