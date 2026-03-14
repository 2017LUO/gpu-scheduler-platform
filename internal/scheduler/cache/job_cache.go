package cache

import (
	model "gpu-scheduler-platform/internal/repo/models"
	"sync"
)

type JobCache struct {
	mu    sync.RWMutex
	items map[string]*model.GPUJob
}

func NewJobCache() *JobCache {
	return &JobCache{
		items: make(map[string]*model.GPUJob),
	}
}

func (c *JobCache) Set(job *model.GPUJob) {
	if job == nil || job.ID == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := *job
	c.items[job.ID] = &cp
}

func (c *JobCache) Get(id string) (*model.GPUJob, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	j, ok := c.items[id]
	if !ok {
		return nil, false
	}
	cp := *j
	return &cp, true
}

func (c *JobCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

func (c *JobCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*model.GPUJob)
}
