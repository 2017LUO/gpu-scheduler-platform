package discovery

import "os"

type NodeMetaResolver struct{}

func NewNodeMetaResolver() *NodeMetaResolver {
	return &NodeMetaResolver{}
}

func (r *NodeMetaResolver) ResolveNodeName(cfgNodeName string) string {
	if cfgNodeName != "" {
		return cfgNodeName
	}
	if v := os.Getenv("NODE_NAME"); v != "" {
		return v
	}
	if h, _ := os.Hostname(); h != "" {
		return h
	}
	return "unknown-node"
}
