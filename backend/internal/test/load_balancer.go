package test

import (
	"sync"
	"sync/atomic"
)

// RoundRobinBalancer 实现简单的轮询负载均衡器
type RoundRobinBalancer struct {
	nodes    []string
	current  uint32
	statuses map[string]bool
	mu       sync.RWMutex
}

// NewRoundRobinBalancer 创建一个新的轮询负载均衡器
func NewRoundRobinBalancer(nodes []string) *RoundRobinBalancer {
	statuses := make(map[string]bool)
	for _, node := range nodes {
		statuses[node] = true
	}
	return &RoundRobinBalancer{
		nodes:    nodes,
		statuses: statuses,
	}
}

// GetNode 获取下一个可用节点
func (b *RoundRobinBalancer) GetNode() string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// 使用原子操作增加计数器
	current := atomic.AddUint32(&b.current, 1)
	for i := 0; i < len(b.nodes); i++ {
		index := (int(current) + i) % len(b.nodes)
		node := b.nodes[index]
		if b.statuses[node] {
			return node
		}
	}
	// 如果没有可用节点，返回第一个节点
	return b.nodes[0]
}

// MarkNodeDown 标记节点为不可用
func (b *RoundRobinBalancer) MarkNodeDown(node string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.statuses[node] = false
}

// MarkNodeUp 标记节点为可用
func (b *RoundRobinBalancer) MarkNodeUp(node string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.statuses[node] = true
}

// WeightedRoundRobinBalancer 实现加权轮询负载均衡器
type WeightedRoundRobinBalancer struct {
	nodes         map[string]int
	currentWeight map[string]int
	totalWeight   int
	mu            sync.Mutex
}

// NewWeightedRoundRobinBalancer 创建一个新的加权轮询负载均衡器
func NewWeightedRoundRobinBalancer(nodes map[string]int) *WeightedRoundRobinBalancer {
	totalWeight := 0
	currentWeight := make(map[string]int)
	for node, weight := range nodes {
		totalWeight += weight
		currentWeight[node] = 0
	}

	return &WeightedRoundRobinBalancer{
		nodes:         nodes,
		currentWeight: currentWeight,
		totalWeight:   totalWeight,
	}
}

// GetNode 获取下一个节点，基于加权轮询算法
func (b *WeightedRoundRobinBalancer) GetNode() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	var selectedNode string
	maxWeight := -1

	// 更新当前权重并选择最大权重的节点
	for node, weight := range b.nodes {
		b.currentWeight[node] += weight
		if b.currentWeight[node] > maxWeight {
			maxWeight = b.currentWeight[node]
			selectedNode = node
		}
	}

	// 减去总权重
	b.currentWeight[selectedNode] -= b.totalWeight

	return selectedNode
}