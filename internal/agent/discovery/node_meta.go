package discovery

import (
	"os"
	"strings"
)

type NodeMeta struct {
	NodeName    string `json:"node_name"`
	Hostname    string `json:"hostname,omitempty"`
	ClusterName string `json:"cluster_name,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	PodName     string `json:"pod_name,omitempty"`
	PodIP       string `json:"pod_ip,omitempty"`
	HostIP      string `json:"host_ip,omitempty"`
}

type NodeMetaResolver struct{}

func NewNodeMetaResolver() *NodeMetaResolver {
	return &NodeMetaResolver{}
}

// ResolveNodeName：兼容当前 service 代码的旧接口。
// 优先级：配置 > NODE_NAME > hostname > unknown-node
func (r *NodeMetaResolver) ResolveNodeName(cfgNodeName string) string {
	meta := r.Resolve(cfgNodeName)
	return meta.NodeName
}

// Resolve：返回完整节点元信息。
// 推荐后续 agent/report 都使用这个方法。
func (r *NodeMetaResolver) Resolve(cfgNodeName string) NodeMeta {
	cfgNodeName = strings.TrimSpace(cfgNodeName)

	hostname := getenvTrim("HOSTNAME")
	if hostname == "" {
		if h, err := os.Hostname(); err == nil {
			hostname = strings.TrimSpace(h)
		}
	}

	nodeName := cfgNodeName
	if nodeName == "" {
		nodeName = getenvTrim("NODE_NAME")
	}
	if nodeName == "" {
		nodeName = hostname
	}
	if nodeName == "" {
		nodeName = "unknown-node"
	}

	return NodeMeta{
		NodeName:    nodeName,
		Hostname:    hostname,
		ClusterName: getenvTrim("CLUSTER_NAME"),
		Namespace:   getenvTrim("POD_NAMESPACE"),
		PodName:     getenvTrim("POD_NAME"),
		PodIP:       getenvTrim("POD_IP"),
		HostIP:      getenvTrim("HOST_IP"),
	}
}

func getenvTrim(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
