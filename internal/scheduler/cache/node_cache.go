package cache

import (
	model "gpu-scheduler-platform/internal/repo/models"
	"sync"
)

type NodeCache struct {
	mu    sync.RWMutex
	items map[string]*model.Node
}

func NewNodeCache() *NodeCache {
	return &NodeCache{
		items: make(map[string]*model.Node),
	}
}

func (c *NodeCache) Set(node *model.Node) {
	if node == nil || node.NodeName == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := *node
	c.items[node.NodeName] = &cp
}

func (c *NodeCache) Get(name string) (*model.Node, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	n, ok := c.items[name]
	if !ok {
		return nil, false
	}
	cp := *n
	return &cp, true
}

func (c *NodeCache) List() []*model.Node {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]*model.Node, 0, len(c.items))
	for _, n := range c.items {
		cp := *n
		out = append(out, &cp)
	}
	return out
}

func (c *NodeCache) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*model.Node)
}
