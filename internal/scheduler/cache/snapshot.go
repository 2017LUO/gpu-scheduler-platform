package cache

import (
	"context"
	"sync"
	"time"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/repo"
)

type SnapshotCache struct {
	mu         sync.RWMutex
	snapshot   *cluster.Snapshot
	loadedAt   time.Time
	ttl        time.Duration
	repository repo.NodeSnapshotRepository
}

func NewSnapshotCache(repository repo.NodeSnapshotRepository, ttl time.Duration) *SnapshotCache {
	if ttl <= 0 {
		ttl = 2 * time.Second
	}
	return &SnapshotCache{
		ttl:        ttl,
		repository: repository,
	}
}

func (c *SnapshotCache) Get(ctx context.Context) (*cluster.Snapshot, error) {
	c.mu.RLock()
	if c.snapshot != nil && time.Since(c.loadedAt) < c.ttl {
		s := *c.snapshot
		c.mu.RUnlock()
		return &s, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.snapshot != nil && time.Since(c.loadedAt) < c.ttl {
		s := *c.snapshot
		return &s, nil
	}

	s, err := c.repository.GetLatest(ctx)
	if err != nil {
		return nil, err
	}
	c.snapshot = s
	c.loadedAt = time.Now()
	cp := *s
	return &cp, nil
}

func (c *SnapshotCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.snapshot = nil
	c.loadedAt = time.Time{}
}
