package cache

import "gpu-scheduler-platform/internal/domain/job"

type JobCache struct {
	items map[string]job.Job
}

func NewJobCache(items []job.Job) *JobCache {
	out := &JobCache{
		items: make(map[string]job.Job, len(items)),
	}
	for _, j := range items {
		out.items[j.ID] = j
	}
	return out
}

func (c *JobCache) Get(id string) (job.Job, bool) {
	if c == nil {
		return job.Job{}, false
	}
	j, ok := c.items[id]
	return j, ok
}

func (c *JobCache) All() []job.Job {
	if c == nil {
		return nil
	}
	out := make([]job.Job, 0, len(c.items))
	for _, j := range c.items {
		out = append(out, j)
	}
	return out
}
