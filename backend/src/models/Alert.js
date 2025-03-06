import mongoose from 'mongoose';

// 告警规则 Schema
const AlertRuleSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
  },
  metric: {
    type: String,
    required: true,
  },
  threshold: {
    type: Number,
    required: true,
  },
  condition: {
    type: String,
    enum: ['>', '>=', '<', '<=', '=='],
    required: true,
  },
  enabled: {
    type: Boolean,
    default: true,
  },
  notificationChannels: {
    type: [String],
    default: [],
  },
  createdAt: {
    type: Date,
    default: Date.now,
  },
  updatedAt: {
    type: Date,
    default: null,
  },
  createdBy: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
  },
});

// 告警历史 Schema
const AlertHistorySchema = new mongoose.Schema({
  ruleId: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'AlertRule',
  },
  ruleName: {
    type: String,
    required: true,
  },
  metric: {
    type: String,
    required: true,
  },
  threshold: {
    type: Number,
    required: true,
  },
  condition: {
    type: String,
    required: true,
  },
  actualValue: {
    type: Number,
    required: true,
  },
  timestamp: {
    type: Date,
    default: Date.now,
  },
  resolved: {
    type: Boolean,
    default: false,
  },
  resolvedAt: {
    type: Date,
    default: null,
  },
});

// 更新时间
AlertRuleSchema.pre('findOneAndUpdate', function(next) {
  this.set({ updatedAt: new Date() });
  next();
});

export const AlertRule = mongoose.model('AlertRule', AlertRuleSchema);
export const AlertHistory = mongoose.model('AlertHistory', AlertHistorySchema);
