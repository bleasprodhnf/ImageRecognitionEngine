// 模型相关的模拟数据接口

// 模型版本列表
export const modelVersions = [
  {
    id: 1,
    version: 'v1.0.0',
    description: '初始版本',
    releaseDate: '2024-01-01',
    status: 'production',
    accuracy: 0.85,
    parameters: {
      batchSize: 32,
      learningRate: 0.001,
      epochs: 100
    }
  },
  {
    id: 2,
    version: 'v1.1.0',
    description: '优化识别准确率',
    releaseDate: '2024-01-15',
    status: 'testing',
    accuracy: 0.89,
    parameters: {
      batchSize: 64,
      learningRate: 0.0005,
      epochs: 150
    }
  },
  {
    id: 3,
    version: 'v1.2.0-beta',
    description: '新增特征提取',
    releaseDate: '2024-01-20',
    status: 'development',
    accuracy: 0.92,
    parameters: {
      batchSize: 128,
      learningRate: 0.0003,
      epochs: 200
    }
  }
]

// 模型参数配置
export const modelParams = [
  {
    id: 1,
    name: 'batchSize',
    value: '64',
    type: 'number',
    description: '批处理大小'
  },
  {
    id: 2,
    name: 'learningRate',
    value: '0.001',
    type: 'number',
    description: '学习率'
  },
  {
    id: 3,
    name: 'epochs',
    value: '100',
    type: 'number',
    description: '训练轮数'
  },
  {
    id: 4,
    name: 'optimizer',
    value: 'adam',
    type: 'string',
    description: '优化器类型'
  }
]

// 模型监控数据
export const modelMonitorData = {
  // 性能数据
  performance: {
    responseTime: [
      { time: '00:00', value: 120 },
      { time: '04:00', value: 135 },
      { time: '08:00', value: 180 },
      { time: '12:00', value: 210 },
      { time: '16:00', value: 190 },
      { time: '20:00', value: 160 },
      { time: '24:00', value: 140 }
    ],
    throughput: [
      { time: '00:00', value: 1000 },
      { time: '04:00', value: 800 },
      { time: '08:00', value: 2500 },
      { time: '12:00', value: 3500 },
      { time: '16:00', value: 3000 },
      { time: '20:00', value: 2000 },
      { time: '24:00', value: 1200 }
    ]
  },
  // 资源使用
  resources: {
    cpu: [
      { time: '00:00', value: 45 },
      { time: '04:00', value: 40 },
      { time: '08:00', value: 65 },
      { time: '12:00', value: 80 },
      { time: '16:00', value: 75 },
      { time: '20:00', value: 60 },
      { time: '24:00', value: 50 }
    ],
    memory: [
      { time: '00:00', value: 50 },
      { time: '04:00', value: 45 },
      { time: '08:00', value: 70 },
      { time: '12:00', value: 85 },
      { time: '16:00', value: 80 },
      { time: '20:00', value: 65 },
      { time: '24:00', value: 55 }
    ]
  },
  // 调用统计
  calls: {
    success: 9750,
    failed: 250,
    total: 10000,
    errorRate: 2.5
  }
}