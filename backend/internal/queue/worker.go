package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Worker 任务处理器
type Worker struct {
	queue              *Queue
	notificationService *NotificationService
	handlers           map[TaskType]TaskHandler
	concurrent         int
	ctx                context.Context
	cancel             context.CancelFunc
	wg                 sync.WaitGroup
}

// NewWorker 创建新的任务处理器
func NewWorker(queue *Queue, notificationService *NotificationService, concurrent int) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		queue:              queue,
		notificationService: notificationService,
		handlers:           make(map[TaskType]TaskHandler),
		concurrent:         concurrent,
		ctx:                ctx,
		cancel:             cancel,
	}
}

// TaskHandler 任务处理函数类型
type TaskHandler func(context.Context, *Task) error

// RegisterHandler 注册任务处理函数
func (w *Worker) RegisterHandler(taskType TaskType, handler TaskHandler) {
	w.handlers[taskType] = handler
}

// Start 启动工作器
func (w *Worker) Start() {
	for i := 0; i < w.concurrent; i++ {
		w.wg.Add(1)
		go w.process()
	}
}

// Stop 停止工作器
func (w *Worker) Stop() {
	w.cancel()
	w.wg.Wait()
}

// process 处理任务的工作循环
func (w *Worker) process() {
	defer w.wg.Done()

	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			// 轮询所有任务类型
			for taskType := range w.handlers {
				task, err := w.queue.Pop(w.ctx, taskType)
				if err != nil {
					log.Printf("Error popping task: %v", err)
					continue
				}
				if task == nil {
					continue
				}

				// 更新任务状态为处理中
				err = w.queue.UpdateTaskStatus(w.ctx, task.ID, "processing")
				if err != nil {
					log.Printf("Error updating task status: %v", err)
				}

				// 执行任务处理
				handler := w.handlers[task.Type]
				err = handler(w.ctx, task)

				// 更新任务状态
				status := "completed"
				notificationType := NotificationTypeTaskComplete
				message := "任务处理成功"

				if err != nil {
					status = "failed"
					notificationType = NotificationTypeTaskFailed
					message = fmt.Sprintf("任务处理失败: %v", err)
					log.Printf("Error processing task %s: %v", task.ID, err)
				}

				err = w.queue.UpdateTaskStatus(w.ctx, task.ID, status)
				if err != nil {
					log.Printf("Error updating task status: %v", err)
				}

				// 发送任务处理结果通知
				err = w.notificationService.SendNotification(w.ctx, task.ID, notificationType, message)
				if err != nil {
					log.Printf("Error sending notification: %v", err)
				}
			}

			// 避免过于频繁的轮询
			time.Sleep(100 * time.Millisecond)
		}
	}
}

// HandleImageRecognition 处理图像识别任务的示例处理器
func HandleImageRecognition(ctx context.Context, task *Task) error {
	// 解析任务数据
	var data struct {
		ImageURL string `json:"image_url"`
		ModelID  string `json:"model_id"`
	}

	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("unmarshal task data error: %v", err)
	}

	// TODO: 实现图像识别逻辑
	// 1. 下载图像
	// 2. 加载模型
	// 3. 执行识别
	// 4. 保存结果

	return nil
}

// HandleModelTraining 处理模型训练任务的示例处理器
func HandleModelTraining(ctx context.Context, task *Task) error {
	// 解析任务数据
	var data struct {
		DatasetID string `json:"dataset_id"`
		ModelID   string `json:"model_id"`
		Params    map[string]interface{} `json:"params"`
	}

	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("unmarshal task data error: %v", err)
	}

	// TODO: 实现模型训练逻辑
	// 1. 加载数据集
	// 2. 准备训练环境
	// 3. 执行训练
	// 4. 保存模型

	return nil
}