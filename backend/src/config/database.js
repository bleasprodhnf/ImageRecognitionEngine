import mongoose from 'mongoose';
import dotenv from 'dotenv';

dotenv.config();

const MONGODB_URI = process.env.MONGODB_URI || 'mongodb://localhost:27017/image-recognition-platform';

// 数据库连接选项
const options = {
  useNewUrlParser: true,
  useUnifiedTopology: true,
};

// 连接数据库
const connectDB = async () => {
  try {
    await mongoose.connect(MONGODB_URI, options);
    console.log('MongoDB 连接成功');
    return mongoose.connection;
  } catch (error) {
    console.error('MongoDB 连接失败:', error);
    process.exit(1);
  }
};

// 断开数据库连接
const disconnectDB = async () => {
  try {
    await mongoose.disconnect();
    console.log('MongoDB 连接已断开');
  } catch (error) {
    console.error('MongoDB 断开连接失败:', error);
  }
};

export { connectDB, disconnectDB };
