import app from './app.js';
import dotenv from 'dotenv';
import { connectDB } from './config/database.js';

// 加载环境变量
dotenv.config();

const PORT = process.env.PORT || 3000;

// 启动服务器
const startServer = async () => {
  try {
    // 连接数据库
    await connectDB();
    
    // 启动服务器
    app.listen(PORT, () => {
      console.log(`Server running on port ${PORT}`);
    });
  } catch (error) {
    console.error(`服务器启动失败: ${error.message}`);
    process.exit(1);
  }
};

// 处理未捕获的异常
process.on('uncaughtException', (error) => {
  console.error('未捕获的异常:', error);
});

process.on('unhandledRejection', (error) => {
  console.error('未处理的 Promise 拒绝:', error);
});

// 启动服务器
startServer();
