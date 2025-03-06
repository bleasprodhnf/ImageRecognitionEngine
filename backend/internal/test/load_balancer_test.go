package test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRoundRobinBalancer 测试轮询负载均衡器
func TestRoundRobinBalancer(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	balancer := NewRoundRobinBalancer(nodes)

	// 验证节点选择的顺序
	for i := 0; i < len(nodes)*2; i++ {
		expectedNode := nodes[i%len(nodes)]
		actualNode := balancer.GetNode()
		assert.Equal(t, expectedNode, actualNode)
	}
}

// TestWeightedRoundRobinBalancer 测试加权轮询负载均衡器
func TestWeightedRoundRobinBalancer(t *testing.T) {
	nodes := map[string]int{
		"node1": 5,
		"node2": 3,
		"node3": 2,
	}
	balancer := NewWeightedRoundRobinBalancer(nodes)

	// 统计节点被选择的次数
	nodeCount := make(map[string]int)
	totalRequests := 1000

	for i := 0; i < totalRequests; i++ {
		node := balancer.GetNode()
		nodeCount[node]++
	}

	// 验证节点选择比例是否符合权重设置
	totalWeight := 10 // 5 + 3 + 2
	for node, weight := range nodes {
		expectedCount := float64(totalRequests) * float64(weight) / float64(totalWeight)
		actualCount := float64(nodeCount[node])
		// 允许5%的误差
		assert.InDelta(t, expectedCount, actualCount, expectedCount*0.05)
	}
}

// TestConcurrentBalancer 测试负载均衡器的并发性能
func TestConcurrentBalancer(t *testing.T) {
	nodes := []string{"node1", "node2", "node3", "node4"}
	balancer := NewRoundRobinBalancer(nodes)

	concurrency := 100
	requestsPerGoroutine := 1000

	var wg sync.WaitGroup
	nodeCount := make(map[string]int)
	var mu sync.Mutex

	// 启动并发请求
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine; j++ {
				node := balancer.GetNode()
				mu.Lock()
				nodeCount[node]++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// 验证负载是否均衡分布
	totalRequests := concurrency * requestsPerGoroutine
	expectedCount := totalRequests / len(nodes)
	variance := float64(totalRequests) * 0.05 // 允许5%的方差

	for _, count := range nodeCount {
		diff := float64(abs(count - expectedCount))
		assert.Less(t, diff, variance)
	}
}

// TestBalancerFailover 测试负载均衡器的故障转移
func TestBalancerFailover(t *testing.T) {
	nodes := []string{"node1", "node2", "node3"}
	balancer := NewRoundRobinBalancer(nodes)

	// 模拟节点故障
	balancer.MarkNodeDown("node2")

	// 验证故障节点不会被选择
	for i := 0; i < 10; i++ {
		node := balancer.GetNode()
		assert.NotEqual(t, "node2", node)
	}

	// 恢复节点
	balancer.MarkNodeUp("node2")

	// 验证恢复的节点重新参与负载均衡
	foundNode2 := false
	for i := 0; i < len(nodes)*3; i++ {
		if balancer.GetNode() == "node2" {
			foundNode2 = true
			break
		}
	}
	assert.True(t, foundNode2)
}