import mongoose from 'mongoose';

// 模型版本 Schema
const ModelVersionSchema = new mongoose.Schema({
  version: {
    type: String,
    required: true,
  },
  status: {
    type: String,
    enum: ['active', 'inactive'],
    default: 'active',
  },
  createdAt: {
    type: Date,
    default: Date.now,
  },
  trainingTaskId: {
    type: String,
    default: null,
  },
  accuracy: {
    type: Number,
    default: null,
  },
  precision: {
    type: Number,
    default: null,
  },
  recall: {
    type: Number,
    default: null,
  },
  f1Score: {
    type: Number,
    default: null,
  },
  inferenceTime: {
    type: Number,
    default: null,
  },
});

// 模型 Schema
const ModelSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    default: '',
  },
  type: {
    type: String,
    required: true,
  },
  currentVersion: {
    type: String,
    default: '1.0.0',
  },
  createdBy: {
    type: mongoose.Schema.Types.ObjectId,
    ref: 'User',
  },
  createdAt: {
    type: Date,
    default: Date.now,
  },
  updatedAt: {
    type: Date,
    default: Date.now,
  },
  versions: [ModelVersionSchema],
});

// 更新时间
ModelSchema.pre('findOneAndUpdate', function(next) {
  this.set({ updatedAt: new Date() });
  next();
});

export default mongoose.model('Model', ModelSchema);
