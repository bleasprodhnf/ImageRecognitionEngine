import express from 'express';
import { generateToken, generateRefreshToken, verifyToken, apiLimiter } from './middleware/auth.js';
import morgan from 'morgan';
import modelRoutes from './routes/modelRoutes.js';
import monitorRoutes from './routes/monitorRoutes.js';
import { initializeAlertRules } from './controllers/monitorController.js';

const app = express();

// 中间件
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(morgan('combined'));

// API 请求限流
app.use(apiLimiter);

// 跨域处理
app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');
  res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  if (req.method === 'OPTIONS') {
    return res.sendStatus(200);
  }
  next();
});

// 登录接口
app.post('/login', (req, res) => {
  // TODO: 实现用户认证逻辑
  const userId = '123'; // 示例用户 ID
  const token = generateToken(userId);
  const refreshToken = generateRefreshToken(userId);
  res.json({ token, refreshToken });
});

// Token 刷新接口
app.post('/refresh-token', (req, res) => {
  const { refreshToken } = req.body;
  const decoded = verifyToken(refreshToken);
  if (!decoded) {
    return res.status(401).json({ message: '无效的 Refresh Token' });
  }
  const newToken = generateToken(decoded.userId);
  res.json({ token: newToken });
});

// 受保护的路由示例
app.get('/protected', (req, res) => {
  const token = req.headers['authorization']?.split(' ')[1];
  const decoded = verifyToken(token);
  if (!decoded) {
    return res.status(401).json({ message: '未授权访问' });
  }
  res.json({ message: '访问成功', userId: decoded.userId });
});

// 模型管理路由
app.use('/api/models', modelRoutes);

// 监控告警路由
app.use('/api/monitor', monitorRoutes);

// 错误处理中间件
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    message: '服务器内部错误',
    error: process.env.NODE_ENV === 'development' ? err.message : undefined
  });
});

// 初始化默认告警规则
initializeAlertRules();

export default app;
