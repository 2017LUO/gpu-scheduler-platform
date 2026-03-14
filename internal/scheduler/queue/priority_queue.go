package queue

import (
	"container/heap"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"sync"
	"time"
)

type PriorityQueue struct {
	mu       sync.Mutex
	capacity int
	items    jobHeap
	seen     map[string]struct{}
	nowFunc  func() time.Time
}

func NewPriorityQueue(capacity int) *PriorityQueue {
	if capacity <= 0 {
		capacity = 1024
	}
	h := make(jobHeap, 0, capacity)
	heap.Init(&h)

	return &PriorityQueue{
		capacity: capacity,
		items:    h,
		seen:     make(map[string]struct{}),
		nowFunc:  func() time.Time { return time.Now().UTC() },
	}
}

func (q *PriorityQueue) Push(job *model.GPUJob) error {
	if job == nil || job.ID == "" {
		return fmt.Errorf("invalid job")
	}

	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.seen[job.ID]; ok {
		return nil
	}
	if len(q.items) >= q.capacity {
		return fmt.Errorf("queue is full")
	}

	now := q.nowFunc()
	heap.Push(&q.items, &jobItem{
		job:               job,
		effectivePriority: EffectivePriority(job.Priority, job.CreatedAt, now),
		createdAt:         job.CreatedAt,
	})
	q.seen[job.ID] = struct{}{}
	return nil
}

func (q *PriorityQueue) Pop() (*model.GPUJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}
	item := heap.Pop(&q.items).(*jobItem)
	delete(q.seen, item.job.ID)
	return item.job, true
}

func (q *PriorityQueue) Peek() (*model.GPUJob, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return nil, false
	}
	return q.items[0].job, true
}

func (q *PriorityQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

func (q *PriorityQueue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = q.items[:0]
	q.seen = make(map[string]struct{})
}

type jobItem struct {
	job               *model.GPUJob
	effectivePriority int
	createdAt         time.Time
	index             int
}

type jobHeap []*jobItem

func (h jobHeap) Len() int { return len(h) }

func (h jobHeap) Less(i, j int) bool {
	if h[i].effectivePriority != h[j].effectivePriority {
		return h[i].effectivePriority > h[j].effectivePriority
	}
	return h[i].createdAt.Before(h[j].createdAt)
}

func (h jobHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *jobHeap) Push(x any) {
	item := x.(*jobItem)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *jobHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*h = old[:n-1]
	return item
}
