package cache

import "gpu-scheduler-platform/internal/domain/cluster"

type NodeCache struct {
	nodes map[string]cluster.Node
}

func NewNodeCache(snapshot *cluster.Snapshot) *NodeCache {
	out := &NodeCache{
		nodes: make(map[string]cluster.Node),
	}
	if snapshot == nil {
		return out
	}
	for _, n := range snapshot.Nodes {
		out.nodes[n.Name] = n
	}
	return out
}

func (c *NodeCache) Get(nodeName string) (cluster.Node, bool) {
	if c == nil {
		return cluster.Node{}, false
	}
	n, ok := c.nodes[nodeName]
	return n, ok
}

func (c *NodeCache) All() []cluster.Node {
	if c == nil {
		return nil
	}
	out := make([]cluster.Node, 0, len(c.nodes))
	for _, n := range c.nodes {
		out = append(out, n)
	}
	return out
}
