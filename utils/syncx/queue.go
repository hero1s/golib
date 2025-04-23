package syncx

import (
	"sync"
)

// SafeQueue 定义一个线程安全的队列
type SafeQueue struct {
	queue []interface{}
	mu    sync.Mutex
}

// NewSafeQueue 创建一个新的线程安全队列
func NewSafeQueue() *SafeQueue {
	return &SafeQueue{
		queue: make([]interface{}, 0),
	}
}

// Enqueue 向队列中添加元素
func (q *SafeQueue) Push(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, item)
}

// Dequeue 从队列中移除并返回第一个元素
func (q *SafeQueue) Pop() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return nil, false
	}
	item := q.queue[0]
	q.queue = q.queue[1:]
	return item, true
}

// Peek 查看队列的第一个元素但不移除
func (q *SafeQueue) Peek() (interface{}, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return nil, false
	}
	return q.queue[0], true
}

// Size 返回队列的大小
func (q *SafeQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}
