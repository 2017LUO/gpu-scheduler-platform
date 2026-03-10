package dto

type AgentReportRequest struct {
	NodeName    string            `json:"node_name"`
	Timestamp   string            `json:"timestamp"`
	GPUs        []AgentGPUInfo    `json:"gpus"`
	MIGs        []AgentMIGInfo    `json:"migs,omitempty"`
	Topology    []AgentGPULink    `json:"topology,omitempty"`
	PodBindings []AgentPodGPUInfo `json:"pod_bindings,omitempty"`
}

type AgentGPUInfo struct {
	ID            string `json:"id"`
	UUID          string `json:"uuid"`
	Index         int    `json:"index"`
	Model         string `json:"model"`
	MemoryMiB     int64  `json:"memory_mib"`
	FreeMemoryMiB int64  `json:"free_memory_mib"`
	Healthy       bool   `json:"healthy"`
	Health        string `json:"health"`
}

type AgentMIGInfo struct {
	ID         string `json:"id"`
	ParentUUID string `json:"parent_uuid"`
	Profile    string `json:"profile"`
	MemoryMiB  int64  `json:"memory_mib"`
}

type AgentGPULink struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type AgentPodGPUInfo struct {
	PodName   string   `json:"pod_name"`
	Namespace string   `json:"namespace"`
	GPUIDs    []string `json:"gpu_ids"`
}

type AgentHeartbeatRequest struct {
	Type     string `json:"type"`
	NodeName string `json:"node_name"`
	TS       string `json:"ts"`
}

type AgentReportResponse struct {
	Accepted        bool   `json:"accepted"`
	SnapshotVersion string `json:"snapshot_version"`
	NodeName        string `json:"node_name"`
	GPUCount        int    `json:"gpu_count"`
}
