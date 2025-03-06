import jwt from 'jsonwebtoken';
import rateLimit from 'express-rate-limit';
import User from '../models/User.js';

// JWT 配置
const JWT_SECRET = process.env.JWT_SECRET || 'your-secret-key';
const JWT_EXPIRES_IN = '1h';
const JWT_REFRESH_EXPIRES_IN = '7d';

// 生成 JWT
const generateToken = (userId) => {
  return jwt.sign({ userId }, JWT_SECRET, { expiresIn: JWT_EXPIRES_IN });
};

// 生成 Refresh Token
const generateRefreshToken = (userId) => {
  return jwt.sign({ userId }, JWT_SECRET, { expiresIn: JWT_REFRESH_EXPIRES_IN });
};

// 验证 Token 中间件
const verifyToken = async (req, res, next) => {
  try {
    // 从请求头中获取 token
    const authHeader = req.headers['authorization'];
    const token = authHeader && authHeader.split(' ')[1]; // Bearer TOKEN
    
    if (!token) {
      return res.status(401).json({ message: '无效的访问凭证' });
    }
    
    // 验证 token
    const decoded = jwt.verify(token, JWT_SECRET);
    
    if (!decoded || !decoded.userId) {
      return res.status(401).json({ message: '无效的访问凭证' });
    }
    
    // 查找用户
    const user = await User.findById(decoded.userId);
    
    if (!user) {
      return res.status(401).json({ message: '用户不存在' });
    }
    
    if (!user.isActive) {
      return res.status(403).json({ message: '账户已被禁用' });
    }
    
    // 将用户信息添加到请求对象中
    req.user = {
      id: user._id,
      role: user.role,
      username: user.username
    };
    
    next();
  } catch (error) {
    if (error.name === 'JsonWebTokenError') {
      return res.status(401).json({ message: '无效的访问凭证' });
    }
    if (error.name === 'TokenExpiredError') {
      return res.status(401).json({ message: '凭证已过期，请重新登录' });
    }
    console.error('验证令牌错误:', error);
    res.status(500).json({ message: '内部服务器错误' });
  }
};

// 管理员权限检查中间件
const requireAdmin = (req, res, next) => {
  if (!req.user || req.user.role !== 'admin') {
    return res.status(403).json({ message: '没有足够的权限操作' });
  }
  next();
};

// 解析 Token 函数
const decodeToken = (token) => {
  try {
    return jwt.verify(token, JWT_SECRET);
  } catch (error) {
    return null;
  }
};

// API 请求限流
const apiLimiter = rateLimit({
  windowMs: 15 * 60 * 1000, // 15 minutes
  max: 100, // 每个 IP 每 15 分钟最多 100 次请求
  message: '请求过于频繁，请稍后再试',
});

export {
  generateToken,
  generateRefreshToken,
  verifyToken,
  requireAdmin,
  decodeToken,
  apiLimiter,
};
