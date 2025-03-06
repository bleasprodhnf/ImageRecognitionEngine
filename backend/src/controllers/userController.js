import User from '../models/User.js';
import SystemLog from '../models/SystemLog.js';
import { generateToken, generateRefreshToken } from '../middleware/auth.js';

// 登录
export const login = async (req, res) => {
  try {
    const { username, password } = req.body;
    
    if (!username || !password) {
      return res.status(400).json({ message: '用户名和密码为必填项' });
    }
    
    // 查找用户
    const user = await User.findOne({ username });
    
    if (!user) {
      return res.status(401).json({ message: '用户名或密码错误' });
    }
    
    // 验证密码
    if (!user.authenticate(password)) {
      return res.status(401).json({ message: '用户名或密码错误' });
    }
    
    // 检查用户状态
    if (!user.isActive) {
      return res.status(403).json({ message: '账户已被禁用' });
    }
    
    // 生成 Token
    const token = generateToken(user._id);
    const refreshToken = generateRefreshToken(user._id);
    
    // 更新用户的 refreshToken 和最后登录时间
    user.refreshToken = refreshToken;
    user.lastLogin = new Date();
    await user.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'auth-service',
      message: '用户登录成功',
      details: { userId: user._id, username: user.username, role: user.role }
    });
    await log.save();
    
    // 返回用户信息（不包含敏感数据）
    const userInfo = {
      id: user._id,
      username: user.username,
      email: user.email,
      role: user.role
    };
    
    res.json({
      message: '登录成功',
      user: userInfo,
      token,
      refreshToken
    });
  } catch (error) {
    console.error('用户登录失败:', error);
    res.status(500).json({ message: '登录失败', error: error.message });
  }
};

// 刷新 Token
export const refreshToken = async (req, res) => {
  try {
    const { refreshToken } = req.body;
    
    if (!refreshToken) {
      return res.status(400).json({ message: 'Refresh Token 为必填项' });
    }
    
    // 查找具有此 refreshToken 的用户
    const user = await User.findOne({ refreshToken });
    
    if (!user) {
      return res.status(401).json({ message: '无效的 Refresh Token' });
    }
    
    // 生成新的 Token
    const token = generateToken(user._id);
    
    res.json({
      message: 'Token 已刷新',
      token
    });
  } catch (error) {
    console.error('刷新 Token 失败:', error);
    res.status(500).json({ message: '刷新 Token 失败', error: error.message });
  }
};

// 获取所有用户
export const getUsers = async (req, res) => {
  try {
    // 获取查询参数
    const { role, isActive, limit, skip } = req.query;
    
    // 构建查询条件
    const query = {};
    
    if (role) {
      query.role = role;
    }
    
    if (isActive !== undefined) {
      query.isActive = isActive === 'true';
    }
    
    // 执行查询并获取总数
    const total = await User.countDocuments(query);
    
    // 分页
    const userQuery = User.find(query)
      .select('-hashedPassword -salt -refreshToken')
      .sort({ createdAt: -1 });
      
    if (limit) {
      userQuery.limit(parseInt(limit));
    }
    
    if (skip) {
      userQuery.skip(parseInt(skip));
    }
    
    const users = await userQuery.exec();
    
    res.json({
      total,
      users
    });
  } catch (error) {
    console.error('获取用户列表失败:', error);
    res.status(500).json({ message: '获取用户列表失败', error: error.message });
  }
};

// 获取单个用户
export const getUserById = async (req, res) => {
  try {
    const user = await User.findById(req.params.id)
      .select('-hashedPassword -salt -refreshToken');
    
    if (!user) {
      return res.status(404).json({ message: '用户不存在' });
    }
    
    res.json(user);
  } catch (error) {
    console.error('获取用户详情失败:', error);
    res.status(500).json({ message: '获取用户详情失败', error: error.message });
  }
};

// 创建用户
export const createUser = async (req, res) => {
  try {
    const { username, email, password, role } = req.body;
    
    if (!username || !email || !password) {
      return res.status(400).json({ message: '用户名、邮箱和密码为必填项' });
    }
    
    // 检查用户名是否已存在
    const existingUsername = await User.findOne({ username });
    if (existingUsername) {
      return res.status(400).json({ message: '用户名已存在' });
    }
    
    // 检查邮箱是否已存在
    const existingEmail = await User.findOne({ email });
    if (existingEmail) {
      return res.status(400).json({ message: '邮箱已存在' });
    }
    
    // 创建新用户
    const newUser = new User({
      username,
      email,
      password, // 这里利用虚拟字段 'password' 设置，会自动加盐和哈希
      role: role || 'customer'
    });
    
    await newUser.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'user-service',
      message: '创建了新用户',
      details: { userId: newUser._id, username: newUser.username, role: newUser.role }
    });
    await log.save();
    
    // 返回用户信息（不包含敏感数据）
    const userInfo = {
      id: newUser._id,
      username: newUser.username,
      email: newUser.email,
      role: newUser.role,
      createdAt: newUser.createdAt
    };
    
    res.status(201).json({
      message: '用户创建成功',
      user: userInfo
    });
  } catch (error) {
    console.error('创建用户失败:', error);
    res.status(500).json({ message: '创建用户失败', error: error.message });
  }
};

// 更新用户
export const updateUser = async (req, res) => {
  try {
    const { username, email, password, role, isActive } = req.body;
    
    // 查找用户
    const user = await User.findById(req.params.id);
    
    if (!user) {
      return res.status(404).json({ message: '用户不存在' });
    }
    
    // 检查权限 - 只有管理员可以修改其他用户或更改角色
    const isAdmin = req.user && req.user.role === 'admin';
    const isSelfUpdate = req.user && req.user.id === req.params.id;
    
    if (!isAdmin && !isSelfUpdate) {
      return res.status(403).json({ message: '没有权限修改此用户' });
    }
    
    if (role && role !== user.role && !isAdmin) {
      return res.status(403).json({ message: '只有管理员可以更改用户角色' });
    }
    
    // 更新字段
    if (username && username !== user.username) {
      // 检查用户名是否被其他用户使用
      const existingUsername = await User.findOne({ username, _id: { $ne: user._id } });
      if (existingUsername) {
        return res.status(400).json({ message: '用户名已存在' });
      }
      user.username = username;
    }
    
    if (email && email !== user.email) {
      // 检查邮箱是否被其他用户使用
      const existingEmail = await User.findOne({ email, _id: { $ne: user._id } });
      if (existingEmail) {
        return res.status(400).json({ message: '邮箱已存在' });
      }
      user.email = email;
    }
    
    if (password) {
      user.password = password; // 利用虚拟字段重新设置密码
    }
    
    if (isAdmin) {
      if (role) {
        user.role = role;
      }
      
      if (isActive !== undefined) {
        user.isActive = isActive;
      }
    }
    
    await user.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'user-service',
      message: '更新了用户信息',
      details: { userId: user._id, username: user.username, role: user.role }
    });
    await log.save();
    
    // 返回更新后的用户信息（不包含敏感数据）
    const userInfo = {
      id: user._id,
      username: user.username,
      email: user.email,
      role: user.role,
      isActive: user.isActive,
      updatedAt: user.updatedAt
    };
    
    res.json({
      message: '用户更新成功',
      user: userInfo
    });
  } catch (error) {
    console.error('更新用户失败:', error);
    res.status(500).json({ message: '更新用户失败', error: error.message });
  }
};

// 删除用户
export const deleteUser = async (req, res) => {
  try {
    // 检查权限 - 只有管理员可以删除用户
    const isAdmin = req.user && req.user.role === 'admin';
    
    if (!isAdmin) {
      return res.status(403).json({ message: '只有管理员可以删除用户' });
    }
    
    // 查找并删除用户
    const user = await User.findByIdAndDelete(req.params.id);
    
    if (!user) {
      return res.status(404).json({ message: '用户不存在' });
    }
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'user-service',
      message: '删除了用户',
      details: { userId: user._id, username: user.username, role: user.role }
    });
    await log.save();
    
    res.json({
      message: '用户删除成功',
      user: {
        id: user._id,
        username: user.username,
        email: user.email
      }
    });
  } catch (error) {
    console.error('删除用户失败:', error);
    res.status(500).json({ message: '删除用户失败', error: error.message });
  }
};
