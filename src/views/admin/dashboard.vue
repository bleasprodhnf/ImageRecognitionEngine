<template>
  <div class="dashboard-container">
    <h1>管理员仪表盘</h1>
    <el-row :gutter="20">
      <el-col :span="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>系统概览</span>
            </div>
          </template>
          <div class="card-content">
            <div class="stat-item">
              <div class="stat-label">总用户数</div>
              <div class="stat-value">0</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">今日活跃</div>
              <div class="stat-value">0</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>系统状态</span>
            </div>
          </template>
          <div class="card-content">
            <div class="stat-item">
              <div class="stat-label">系统运行时间</div>
              <div class="stat-value">0天</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">CPU使用率</div>
              <div class="stat-value">0%</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <span>最近活动</span>
            </div>
          </template>
          <div class="card-content">
            <el-empty description="暂无活动" />
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
// 后续可以添加实际的数据获取和处理逻辑
</script>

<style scoped>
.dashboard-container {
  padding: 20px;
}

.dashboard-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-content {
  padding: 10px 0;
}

.stat-item {
  margin-bottom: 15px;
}

.stat-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 5px;
}

.stat-value {
  font-size: 24px;
  color: #303133;
  font-weight: bold;
}
</style>