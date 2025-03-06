// 模型控制器
import Model from '../models/Model.js';
import TrainingTask from '../models/TrainingTask.js';
import SystemLog from '../models/SystemLog.js';
import mongoose from 'mongoose';

// 获取所有模型
export const getAllModels = async (req, res) => {
  try {
    const models = await Model.find().sort({ updatedAt: -1 });
    res.json(models);
  } catch (error) {
    console.error('获取模型列表失败:', error);
    res.status(500).json({ message: '获取模型列表失败', error: error.message });
  }
};

// 获取单个模型
export const getModelById = async (req, res) => {
  try {
    const model = await Model.findById(req.params.id);
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    res.json(model);
  } catch (error) {
    console.error('获取模型详情失败:', error);
    res.status(500).json({ message: '获取模型详情失败', error: error.message });
  }
};

// 创建新模型
export const createModel = async (req, res) => {
  try {
    const { name, description, type } = req.body;
    
    if (!name || !type) {
      return res.status(400).json({ message: '模型名称和类型为必填项' });
    }
    
    const newModel = new Model({
      name,
      description,
      type,
      createdBy: req.userId,
      versions: [
        {
          version: '1.0.0',
          status: 'active',
          createdAt: new Date()
        }
      ]
    });
    
    await newModel.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'model-service',
      message: `创建了新模型: ${name}`,
      details: { modelId: newModel._id, modelName: name }
    });
    await log.save();
    
    res.status(201).json(newModel);
  } catch (error) {
    console.error('创建模型失败:', error);
    res.status(500).json({ message: '创建模型失败', error: error.message });
  }
};

// 更新模型
export const updateModel = async (req, res) => {
  try {
    const { name, description, type } = req.body;
    
    // 检查模型是否存在
    const model = await Model.findById(req.params.id);
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    // 更新模型信息
    model.name = name || model.name;
    model.description = description || model.description;
    model.type = type || model.type;
    
    await model.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'model-service',
      message: `更新了模型: ${model.name}`,
      details: { modelId: model._id, modelName: model.name }
    });
    await log.save();
    
    res.json(model);
  } catch (error) {
    console.error('更新模型失败:', error);
    res.status(500).json({ message: '更新模型失败', error: error.message });
  }
};

// 删除模型
export const deleteModel = async (req, res) => {
  try {
    // 查找并删除模型
    const model = await Model.findByIdAndDelete(req.params.id);
    
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    // 删除相关的训练任务
    await TrainingTask.deleteMany({ modelId: req.params.id });
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'model-service',
      message: `删除了模型: ${model.name}`,
      details: { modelId: model._id, modelName: model.name }
    });
    await log.save();
    
    res.json({ message: '模型已删除', model });
  } catch (error) {
    console.error('删除模型失败:', error);
    res.status(500).json({ message: '删除模型失败', error: error.message });
  }
};

// 模型版本回滚
export const rollbackModel = async (req, res) => {
  try {
    const modelId = req.params.id;
    const versionId = req.params.versionId;
    
    // 查找模型
    const model = await Model.findById(modelId);
    
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    // 查找指定版本
    const versionIndex = model.versions.findIndex(v => v._id.toString() === versionId);
    
    if (versionIndex === -1) {
      return res.status(404).json({ message: '模型版本不存在' });
    }
    
    const targetVersion = model.versions[versionIndex];
    
    // 更新当前版本
    model.currentVersion = targetVersion.version;
    
    // 更新所有版本状态
    model.versions.forEach(v => {
      v.status = v._id.toString() === versionId ? 'active' : 'inactive';
    });
    
    await model.save();
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'model-service',
      message: `回滚了模型版本: ${model.name} -> ${targetVersion.version}`,
      details: { 
        modelId: model._id, 
        modelName: model.name,
        versionId: targetVersion._id,
        version: targetVersion.version
      }
    });
    await log.save();
    
    res.json({
      message: '模型版本已回滚',
      model
    });
  } catch (error) {
    console.error('模型版本回滚失败:', error);
    res.status(500).json({ message: '模型版本回滚失败', error: error.message });
  }
};

// 获取模型性能数据
export const getModelPerformance = async (req, res) => {
  try {
    const model = await Model.findById(req.params.id);
    
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    // 查找当前活动版本
    const activeVersion = model.versions.find(v => v.status === 'active');
    
    if (!activeVersion) {
      return res.status(404).json({ message: '没有活动的模型版本' });
    }
    
    // 如果版本已经有性能数据，则直接返回
    if (activeVersion.accuracy && 
        activeVersion.precision && 
        activeVersion.recall && 
        activeVersion.f1Score) {
      return res.json({
        accuracy: activeVersion.accuracy,
        precision: activeVersion.precision,
        recall: activeVersion.recall,
        f1Score: activeVersion.f1Score,
        inferenceTime: activeVersion.inferenceTime || Math.random() * 100 + 50,
        memoryUsage: Math.random() * 1000 + 500,
        lastUpdated: new Date()
      });
    }
    
    // 否则生成模拟性能数据
    const performanceData = {
      accuracy: Math.random() * 0.3 + 0.7,  // 0.7-1.0 之间的随机值
      precision: Math.random() * 0.3 + 0.7,
      recall: Math.random() * 0.3 + 0.7,
      f1Score: Math.random() * 0.3 + 0.7,
      inferenceTime: Math.random() * 100 + 50,  // 50-150ms 之间的随机值
      memoryUsage: Math.random() * 1000 + 500,  // 500-1500MB 之间的随机值
      lastUpdated: new Date()
    };
    
    res.json(performanceData);
  } catch (error) {
    console.error('获取模型性能数据失败:', error);
    res.status(500).json({ message: '获取模型性能数据失败', error: error.message });
  }
};

// 获取模型版本列表
export const getModelVersions = async (req, res) => {
  try {
    const model = await Model.findById(req.params.id);
    
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    res.json(model.versions);
  } catch (error) {
    console.error('获取模型版本列表失败:', error);
    res.status(500).json({ message: '获取模型版本列表失败', error: error.message });
  }
};

// 启动模型训练任务
export const trainModel = async (req, res) => {
  try {
    const model = await Model.findById(req.params.id);
    
    if (!model) {
      return res.status(404).json({ message: '模型不存在' });
    }
    
    const { trainingParams } = req.body;
    
    // 创建新的训练任务
    const trainingTask = new TrainingTask({
      modelId: model._id,
      status: 'pending',
      progress: 0,
      startTime: new Date(),
      params: trainingParams || {},
      createdBy: req.userId,
      logs: [{
        time: new Date(),
        message: '训练任务已创建，等待启动'
      }]
    });
    
    await trainingTask.save();
    
    // 模拟训练进度更新
    // 注意：在实际应用中，这应该由一个后台任务处理
    setTimeout(() => {
      simulateTraining(trainingTask._id);
    }, 1000);
    
    // 记录日志
    const log = new SystemLog({
      level: 'info',
      service: 'model-service',
      message: `启动了模型训练任务: ${model.name}`,
      details: { 
        modelId: model._id, 
        modelName: model.name,
        taskId: trainingTask._id
      }
    });
    await log.save();
    
    res.status(201).json({
      message: '训练任务已创建',
      taskId: trainingTask._id
    });
  } catch (error) {
    console.error('创建训练任务失败:', error);
    res.status(500).json({ message: '创建训练任务失败', error: error.message });
  }
};

// 获取训练任务状态
export const getTrainingStatus = async (req, res) => {
  try {
    const taskId = req.params.taskId;
    const task = await TrainingTask.findById(taskId);
    
    if (!task) {
      return res.status(404).json({ message: '训练任务不存在' });
    }
    
    res.json(task);
  } catch (error) {
    console.error('获取训练任务状态失败:', error);
    res.status(500).json({ message: '获取训练任务状态失败', error: error.message });
  }
};

// 模拟训练进度更新
async function simulateTraining(taskId) {
  try {
    // 查找任务
    let task = await TrainingTask.findById(taskId);
    
    if (!task) {
      console.error(`找不到训练任务: ${taskId}`);
      return;
    }
    
    // 更新任务状态为运行中
    task.status = 'running';
    await task.save();
    
    // 每秒更新训练进度
    const interval = setInterval(async () => {
      try {
        // 重新查询任务以获取最新状态
        task = await TrainingTask.findById(taskId);
        
        if (!task) {
          clearInterval(interval);
          return;
        }
        
        // 增加进度
        task.progress += 5;
        task.logs.push({
          time: new Date(),
          message: `训练进度: ${task.progress}%`
        });
        
        await task.save();
        
        // 如果训练完成
        if (task.progress >= 100) {
          clearInterval(interval);
          
          // 更新任务状态
          task.status = 'completed';
          task.progress = 100;
          task.endTime = new Date();
          await task.save();
          
          // 更新模型版本信息
          const model = await Model.findById(task.modelId);
          
          if (model) {
            const lastVersion = model.currentVersion;
            const versionParts = lastVersion.split('.');
            const newVersion = `${versionParts[0]}.${versionParts[1]}.${parseInt(versionParts[2]) + 1}`;
            
            // 更新模型的当前版本
            model.currentVersion = newVersion;
            
            // 设置所有现有版本为非活动状态
            model.versions.forEach(v => {
              v.status = 'inactive';
            });
            
            // 添加新版本
            const performanceMetrics = {
              accuracy: Math.random() * 0.3 + 0.7,
              precision: Math.random() * 0.3 + 0.7,
              recall: Math.random() * 0.3 + 0.7,
              f1Score: Math.random() * 0.3 + 0.7,
              inferenceTime: Math.random() * 100 + 50
            };
            
            model.versions.push({
              version: newVersion,
              status: 'active',
              createdAt: new Date(),
              trainingTaskId: task._id,
              ...performanceMetrics
            });
            
            await model.save();
            
            // 记录日志
            const log = new SystemLog({
              level: 'info',
              service: 'model-service',
              message: `模型训练完成，创建了新版本: ${model.name} ${newVersion}`,
              details: { 
                modelId: model._id, 
                modelName: model.name,
                version: newVersion,
                taskId: task._id
              }
            });
            await log.save();
          }
        }
      } catch (error) {
        console.error('更新训练进度失败:', error);
        clearInterval(interval);
      }
    }, 1000);
    
  } catch (error) {
    console.error('模拟训练任务失败:', error);
  }
}
