<template>
  <div class="accuracy-container">
    <h2>准确率统计</h2>
    <el-card class="box-card">
      <div class="stats-overview">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>总体准确率</span>
                </div>
              </template>
              <div class="stats-value">0%</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>样本数量</span>
                </div>
              </template>
              <div class="stats-value">0</div>
            </el-card>
          </el-col>
          <el-col :span="8">
            <el-card shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>评估周期</span>
                </div>
              </template>
              <div class="stats-value">0天</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
      <div class="stats-charts">
        <el-tabs v-model="activeTab">
          <el-tab-pane label="准确率趋势" name="trend">
            <div class="chart-placeholder">准确率趋势图表</div>
          </el-tab-pane>
          <el-tab-pane label="分类统计" name="categories">
            <div class="chart-placeholder">分类统计图表</div>
          </el-tab-pane>
          <el-tab-pane label="错误分析" name="errors">
            <div class="chart-placeholder">错误分析图表</div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const activeTab = ref('trend')
</script>

<style scoped>
.accuracy-container {
  padding: 20px;
}

.stats-overview {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats-value {
  font-size: 24px;
  font-weight: bold;
  text-align: center;
  color: #409EFF;
}

.stats-charts {
  margin-top: 20px;
}

.chart-placeholder {
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  border-radius: 4px;
  color: #909399;
}
</style>