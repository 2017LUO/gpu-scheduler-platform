package cache

import (
	model "gpu-scheduler-platform/internal/repo/models"
	"sync"
	"time"
)

type Snapshot struct {
	Nodes       []*model.Node
	GeneratedAt time.Time
}

type SnapshotCache struct {
	mu       sync.RWMutex
	snapshot *Snapshot
}

func NewSnapshotCache() *SnapshotCache {
	return &SnapshotCache{}
}

func (c *SnapshotCache) Set(nodes []*model.Node) {
	c.mu.Lock()
	defer c.mu.Unlock()

	out := make([]*model.Node, 0, len(nodes))
	for _, n := range nodes {
		if n == nil {
			continue
		}
		cp := *n
		out = append(out, &cp)
	}

	c.snapshot = &Snapshot{
		Nodes:       out,
		GeneratedAt: time.Now().UTC(),
	}
}

func (c *SnapshotCache) Get() *Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.snapshot == nil {
		return nil
	}
	out := make([]*model.Node, 0, len(c.snapshot.Nodes))
	for _, n := range c.snapshot.Nodes {
		cp := *n
		out = append(out, &cp)
	}

	return &Snapshot{
		Nodes:       out,
		GeneratedAt: c.snapshot.GeneratedAt,
	}
}
