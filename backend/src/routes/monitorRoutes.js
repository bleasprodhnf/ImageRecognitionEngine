import express from 'express';
const router = express.Router();
import { verifyToken } from '../middleware/auth.js';
import * as monitorController from '../controllers/monitorController.js';

// 验证 Token 中间件
const authenticate = (req, res, next) => {
  const token = req.headers['authorization']?.split(' ')[1];
  const decoded = verifyToken(token);
  if (!decoded) {
    return res.status(401).json({ message: '未授权访问' });
  }
  req.userId = decoded.userId;
  next();
};

// 获取系统资源使用情况
router.get('/resources', authenticate, monitorController.getSystemResources);

// 获取告警规则列表
router.get('/alerts/rules', authenticate, monitorController.getAlertRules);

// 创建告警规则
router.post('/alerts/rules', authenticate, monitorController.createAlertRule);

// 更新告警规则
router.put('/alerts/rules/:id', authenticate, monitorController.updateAlertRule);

// 删除告警规则
router.delete('/alerts/rules/:id', authenticate, monitorController.deleteAlertRule);

// 获取告警历史
router.get('/alerts/history', authenticate, monitorController.getAlertHistory);

// 获取系统健康状态
router.get('/health', authenticate, monitorController.getSystemHealth);

// 获取实时日志
router.get('/logs', authenticate, monitorController.getLogs);

export default router;
