import mongoose from 'mongoose';

// 系统日志 Schema
const SystemLogSchema = new mongoose.Schema({
  level: {
    type: String,
    enum: ['info', 'warning', 'error', 'debug'],
    required: true,
  },
  service: {
    type: String,
    required: true,
  },
  message: {
    type: String,
    required: true,
  },
  details: {
    type: Object,
    default: null,
  },
  timestamp: {
    type: Date,
    default: Date.now,
  },
});

// 添加索引以提高查询性能
SystemLogSchema.index({ timestamp: -1 });
SystemLogSchema.index({ level: 1, timestamp: -1 });
SystemLogSchema.index({ service: 1, timestamp: -1 });

export default mongoose.model('SystemLog', SystemLogSchema);
