package test

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// LoadBalancer 简单的负载均衡器实现
type LoadBalancer struct {
	nodes []string
	current int
	mutex sync.Mutex
}

// NewLoadBalancer 创建新的负载均衡器
func NewLoadBalancer(nodes []string) *LoadBalancer {
	return &LoadBalancer{
		nodes: nodes,
		current: 0,
	}
}

// GetNode 获取下一个节点
func (lb *LoadBalancer) GetNode() string {
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	node := lb.nodes[lb.current]
	lb.current = (lb.current + 1) % len(lb.nodes)
	return node
}

// TestLoadBalancing 测试负载均衡性能
func TestLoadBalancing(t *testing.T) {
	// 初始化负载均衡器
	nodes := []string{"node1", "node2", "node3"}
	loadBalancer := NewLoadBalancer(nodes)

	nodeCount := len(nodes)
	requestCount := 1000

	// 记录每个节点处理的请求数
	nodeRequests := make([]int, nodeCount)

	// 模拟请求分发
	for i := 0; i < requestCount; i++ {
		// 获取下一个节点
		nodeIndex := loadBalancer.GetNode()
		// 记录请求数
		for j, node := range nodes {
			if node == nodeIndex {
				nodeRequests[j]++
				break
			}
		}
	}

	// 验证负载是否均衡分布
	expectedRequests := requestCount / nodeCount
	variance := float64(100) // 允许的方差

	for _, count := range nodeRequests {
		diff := float64(abs(count - expectedRequests))
		assert.Less(t, diff, variance)
	}
}

// TestConcurrentImageRecognition 测试图像识别接口的并发性能
func TestConcurrentImageRecognition(t *testing.T) {
	concurrency := 100 // 并发请求数
	totalRequests := 1000 // 总请求数

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 创建通道用于收集响应时间
	responseTimes := make(chan time.Duration, totalRequests)

	// 启动并发请求
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < totalRequests/concurrency; j++ {
				select {
				case <-ctx.Done():
					return
				default:
					start := time.Now()
					// 调用图像识别API
					resp, err := http.Post("http://localhost:8080/api/v1/recognition", "application/json", bytes.NewBuffer([]byte(`{"image_url": "test_image.jpg"}`)))
					if err != nil {
						t.Errorf("Failed to send request: %v", err)
						return
					}
					defer resp.Body.Close()
					responseTimes <- time.Since(start)
				}
			}
		}()
	}

	// 等待所有请求完成
	wg.Wait()
	close(responseTimes)

	// 计算性能指标
	var total time.Duration
	count := 0
	for rt := range responseTimes {
		total += rt
		count++
	}

	avgResponseTime := total / time.Duration(count)
	// 验证平均响应时间是否在可接受范围内
	assert.Less(t, avgResponseTime, 500*time.Millisecond)
}

// TestSystemStability 测试系统稳定性
func TestSystemStability(t *testing.T) {
	duration := 5 * time.Minute
	tickInterval := 1 * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()

	var failureCount int
	var totalRequests int

	for {
		select {
		case <-ctx.Done():
			// 计算系统稳定性指标
			successRate := float64(totalRequests-failureCount) / float64(totalRequests) * 100
			// 验证成功率是否达到要求
			assert.GreaterOrEqual(t, successRate, 99.9)
			return
		case <-ticker.C:
			totalRequests++
			// 发送请求并检查响应
			resp, err := http.Get("http://localhost:8080/api/v1/health")
			if err != nil {
			    failureCount++
			    continue
			}
			if resp.StatusCode != http.StatusOK {
			    failureCount++
			}
			resp.Body.Close()
		}
	}
}

// abs 返回整数的绝对值
func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}