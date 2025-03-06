import express from 'express';
const router = express.Router();
import { verifyToken } from '../middleware/auth.js';
import * as modelController from '../controllers/modelController.js';

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

// 获取所有模型
router.get('/', authenticate, modelController.getAllModels);

// 获取单个模型
router.get('/:id', authenticate, modelController.getModelById);

// 新增模型
router.post('/', authenticate, modelController.createModel);

// 更新模型
router.put('/:id', authenticate, modelController.updateModel);

// 删除模型
router.delete('/:id', authenticate, modelController.deleteModel);

// 模型版本回滚
router.post('/:id/rollback/:versionId', authenticate, modelController.rollbackModel);

// 获取模型性能数据
router.get('/:id/performance', authenticate, modelController.getModelPerformance);

// 获取模型版本列表
router.get('/:id/versions', authenticate, modelController.getModelVersions);

// 模型训练任务
router.post('/:id/train', authenticate, modelController.trainModel);

// 获取训练任务状态
router.get('/train/:taskId', authenticate, modelController.getTrainingStatus);

export default router;
