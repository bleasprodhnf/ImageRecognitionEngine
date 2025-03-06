// 监控控制器
import { AlertRule, AlertHistory } from '../models/Alert.js';
import SystemLog from '../models/SystemLog.js';
import os from 'os';
import mongoose from 'mongoose';

// 生成系统资源使用数据
const generateResourceData = () => {
  // 获取实际系统数据
  const cpuUsage = Math.random() * 100; // 在实际应用中，应该使用如 node-os-utils 等库来获取真实 CPU 使用率
  const totalMem = os.totalmem();
  const freeMem = os.freemem();
  const usedMem = totalMem - freeMem;
  
  // 模拟磁盘数据（在实际应用中应使用 fs 或其他磁盘监控库）
  const totalDisk = 512000;
  const usedDisk = Math.random() * totalDisk;
  const freeDisk = totalDisk - usedDisk;
  
  return {
    cpu: {
      usage: cpuUsage,
      cores: os.cpus().length,
      temperature: 40 + Math.random() * 20
    },
    memory: {
      total: Math.floor(totalMem / (1024 * 1024)),
      used: Math.floor(usedMem / (1024 * 1024)),
      free: Math.floor(freeMem / (1024 * 1024))
    },
    disk: {
      total: totalDisk,
      used: usedDisk,
      free: freeDisk
    },
    network: {
      inbound: Math.random() * 1000,
      outbound: Math.random() * 1000
    },
    timestamp: new Date()
  };
};

// 获取系统资源使用情况
export const getSystemResources = async (req, res) => {
  try {
    // 生成资源数据
    const resourceData = generateResourceData();
    
    // 检查是否触发告警
    await checkAlerts(resourceData);
    
    res.json(resourceData);
  } catch (error) {
    console.error('获取系统资源数据失败:', error);
    res.status(500).json({ message: '获取系统资源数据失败', error: error.message });
  }
};

// 获取告警规则列表
export const getAlertRules = async (req, res) => {
  try {
    const rules = await AlertRule.find().sort({ createdAt: -1 });
    res.json(rules);
  } catch (error) {
    console.error('获取告警规则失败:', error);
    res.status(500).json({ message: '获取告警规则失败', error: error.message });
  }
};

// 创建告警规则
export const createAlertRule = async (req, res) => {
  try {
    const { name, metric, threshold, condition, enabled, notificationChannels } = req.body;
    
    if (!name || !metric || !threshold || !condition) {
      return res.status(400).json({ message: '名称、指标、阈值和条件为必填项' });
    }
    
    const newRule = new AlertRule({
      name,
      metric,
      threshold,
      condition,
      enabled: enabled !== undefined ? enabled : true,
      notificationChannels: notificationChannels || [],
      createdBy: req.userId
    });
    
    await newRule.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'monitoring-service',
      message: `创建了新的告警规则: ${name}`,
      details: { ruleId: newRule._id, ruleName: name }
    });
    await log.save();
    
    res.status(201).json(newRule);
  } catch (error) {
    console.error('创建告警规则失败:', error);
    res.status(500).json({ message: '创建告警规则失败', error: error.message });
  }
};

// 更新告警规则
export const updateAlertRule = async (req, res) => {
  try {
    const { name, metric, threshold, condition, enabled, notificationChannels } = req.body;
    
    // 查找规则
    const rule = await AlertRule.findById(req.params.id);
    
    if (!rule) {
      return res.status(404).json({ message: '告警规则不存在' });
    }
    
    // 更新规则
    rule.name = name || rule.name;
    rule.metric = metric || rule.metric;
    rule.threshold = threshold !== undefined ? threshold : rule.threshold;
    rule.condition = condition || rule.condition;
    rule.enabled = enabled !== undefined ? enabled : rule.enabled;
    rule.notificationChannels = notificationChannels || rule.notificationChannels;
    
    await rule.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'monitoring-service',
      message: `更新了告警规则: ${rule.name}`,
      details: { ruleId: rule._id, ruleName: rule.name }
    });
    await log.save();
    
    res.json(rule);
  } catch (error) {
    console.error('更新告警规则失败:', error);
    res.status(500).json({ message: '更新告警规则失败', error: error.message });
  }
};

// 删除告警规则
export const deleteAlertRule = async (req, res) => {
  try {
    // 查找并删除规则
    const rule = await AlertRule.findByIdAndDelete(req.params.id);
    
    if (!rule) {
      return res.status(404).json({ message: '告警规则不存在' });
    }
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'monitoring-service',
      message: `删除了告警规则: ${rule.name}`,
      details: { ruleId: rule._id, ruleName: rule.name }
    });
    await log.save();
    
    res.json({ message: '告警规则已删除', rule });
  } catch (error) {
    console.error('删除告警规则失败:', error);
    res.status(500).json({ message: '删除告警规则失败', error: error.message });
  }
};

// 获取告警历史
export const getAlertHistory = async (req, res) => {
  try {
    const { startDate, endDate, limit } = req.query;
    let query = {};
    
    // 根据日期过滤
    if (startDate || endDate) {
      query.timestamp = {};
      
      if (startDate) {
        query.timestamp.$gte = new Date(startDate);
      }
      
      if (endDate) {
        query.timestamp.$lte = new Date(endDate);
      }
    }
    
    // 查询并排序
    let alertHistoryQuery = AlertHistory.find(query).sort({ timestamp: -1 });
    
    // 限制返回数量
    if (limit) {
      alertHistoryQuery = alertHistoryQuery.limit(parseInt(limit));
    }
    
    const alertHistory = await alertHistoryQuery.exec();
    
    res.json(alertHistory);
  } catch (error) {
    console.error('获取告警历史失败:', error);
    res.status(500).json({ message: '获取告警历史失败', error: error.message });
  }
};

// 获取系统健康状态
export const getSystemHealth = async (req, res) => {
  try {
    // 在实际应用中，这些数据应该从健康检查服务收集
    // 这里模拟检查各服务的健康状态
    
    // 检查数据库连接
    let dbStatus = 'healthy';
    try {
      if (mongoose.connection.readyState !== 1) {
        dbStatus = 'down';
      }
    } catch (error) {
      dbStatus = 'down';
    }
    
    const services = [
      {
        name: '认证服务',
        status: 'healthy',
        uptime: Math.floor(Math.random() * 1000000),
        lastCheck: new Date()
      },
      {
        name: '模型服务',
        status: 'healthy',
        uptime: Math.floor(Math.random() * 1000000),
        lastCheck: new Date()
      },
      {
        name: '数据库服务',
        status: dbStatus,
        uptime: Math.floor(Math.random() * 1000000),
        lastCheck: new Date()
      },
      {
        name: '缓存服务',
        status: Math.random() > 0.9 ? 'degraded' : 'healthy',
        uptime: Math.floor(Math.random() * 1000000),
        lastCheck: new Date()
      }
    ];
    
    const overallStatus = services.every(s => s.status === 'healthy') ? 'healthy' : 
                         services.some(s => s.status === 'down') ? 'down' : 'degraded';
    
    // 如果系统状态不健康，记录日志
    if (overallStatus !== 'healthy') {
      const log = new SystemLog({
        level: overallStatus === 'down' ? 'error' : 'warning',
        service: 'monitoring-service',
        message: `系统健康状态: ${overallStatus}`,
        details: { services: services.filter(s => s.status !== 'healthy').map(s => s.name) }
      });
      await log.save();
    }
    
    res.json({
      status: overallStatus,
      services,
      timestamp: new Date()
    });
  } catch (error) {
    console.error('获取系统健康状态失败:', error);
    res.status(500).json({ message: '获取系统健康状态失败', error: error.message });
  }
};

// 获取系统日志
export const getLogs = async (req, res) => {
  try {
    const { level, service, limit, startDate, endDate } = req.query;
    let query = {};
    
    // 根据日志级别过滤
    if (level) {
      query.level = level;
    }
    
    // 根据服务过滤
    if (service) {
      query.service = service;
    }
    
    // 根据日期过滤
    if (startDate || endDate) {
      query.timestamp = {};
      
      if (startDate) {
        query.timestamp.$gte = new Date(startDate);
      }
      
      if (endDate) {
        query.timestamp.$lte = new Date(endDate);
      }
    }
    
    // 查询并排序
    let logsQuery = SystemLog.find(query).sort({ timestamp: -1 });
    
    // 限制返回数量
    const logLimit = limit ? parseInt(limit) : 100;
    logsQuery = logsQuery.limit(logLimit);
    
    const logs = await logsQuery.exec();
    
    res.json(logs);
  } catch (error) {
    console.error('获取系统日志失败:', error);
    res.status(500).json({ message: '获取系统日志失败', error: error.message });
  }
};

// 检查告警规则
async function checkAlerts(resourceData) {
  try {
    // 仅检查启用的告警规则
    const enabledRules = await AlertRule.find({ enabled: true });
    
    for (const rule of enabledRules) {
      let metricValue;
      
      // 根据规则指标获取对应的资源数据
      if (rule.metric === 'cpu.usage') {
        metricValue = resourceData.cpu.usage;
      } else if (rule.metric === 'memory.used') {
        metricValue = resourceData.memory.used;
      } else if (rule.metric === 'disk.used') {
        metricValue = resourceData.disk.used;
      } else if (rule.metric === 'network.inbound') {
        metricValue = resourceData.network.inbound;
      } else if (rule.metric === 'network.outbound') {
        metricValue = resourceData.network.outbound;
      }
      
      let isAlertTriggered = false;
      
      // 根据条件检查是否触发告警
      if (rule.condition === '>' && metricValue > rule.threshold) {
        isAlertTriggered = true;
      } else if (rule.condition === '>=' && metricValue >= rule.threshold) {
        isAlertTriggered = true;
      } else if (rule.condition === '<' && metricValue < rule.threshold) {
        isAlertTriggered = true;
      } else if (rule.condition === '<=' && metricValue <= rule.threshold) {
        isAlertTriggered = true;
      } else if (rule.condition === '==' && metricValue === rule.threshold) {
        isAlertTriggered = true;
      }
      
      // 如果触发告警，记录并发送通知
      if (isAlertTriggered) {
        const alert = new AlertHistory({
          ruleId: rule._id,
          ruleName: rule.name,
          metric: rule.metric,
          threshold: rule.threshold,
          condition: rule.condition,
          actualValue: metricValue,
          timestamp: new Date()
        });
        
        await alert.save();
        
        // 在实际应用中，这里应该发送通知
        // await sendNotifications(alert, rule.notificationChannels);
        
        // 添加日志
        const log = new SystemLog({
          level: 'error',
          service: 'monitoring-service',
          message: `告警触发: ${rule.name}`,
          details: {
            alertId: alert._id,
            ruleId: rule._id,
            ruleName: rule.name,
            metric: rule.metric,
            threshold: rule.threshold,
            actualValue: metricValue
          }
        });
        await log.save();
      }
    }
  } catch (error) {
    console.error('检查告警规则失败:', error);
  }
}

// 初始化示例告警规则
export const initializeAlertRules = async () => {
  try {
    // 检查是否已有规则
    const count = await AlertRule.countDocuments();
    
    if (count === 0) {
      // 初始化一些默认告警规则
      const defaultRules = [
        {
          name: 'CPU 使用率过高',
          metric: 'cpu.usage',
          threshold: 90,
          condition: '>',
          enabled: true,
          notificationChannels: ['email', 'sms']
        },
        {
          name: '内存不足',
          metric: 'memory.used',
          threshold: 14000,
          condition: '>',
          enabled: true,
          notificationChannels: ['email']
        },
        {
          name: '磁盘空间不足',
          metric: 'disk.used',
          threshold: 450000,
          condition: '>',
          enabled: true,
          notificationChannels: ['email', 'sms']
        }
      ];
      
      await AlertRule.insertMany(defaultRules);
      console.log('初始化默认告警规则完成');
    }
  } catch (error) {
    console.error('初始化告警规则失败:', error);
  }
};
