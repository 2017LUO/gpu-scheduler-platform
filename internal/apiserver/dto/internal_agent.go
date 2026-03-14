package dto

type AgentHeartbeatRequest struct {
	NodeName string  `json:"node_name"`
	Status   string  `json:"status"`
	Message  *string `json:"message,omitempty"`
	SeenAt   string  `json:"seen_at,omitempty"`
}

type AgentHeartbeatResponse struct {
	NodeName string `json:"node_name"`
	Status   string `json:"status"`
	SeenAt   string `json:"seen_at"`
}

type AgentReportRequest struct {
	Version         string                `json:"version"`
	AgentVersion    string                `json:"agent_version"`
	ClusterName     string                `json:"cluster_name"`
	NodeName        string                `json:"node_name"`
	Source          string                `json:"source"`
	NodeState       string                `json:"node_state"`
	Schedulable     bool                  `json:"schedulable"`
	GPUCount        int                   `json:"gpu_count"`
	HealthyGPUCount int                   `json:"healthy_gpu_count"`
	TotalMemoryMiB  int64                 `json:"total_memory_mib"`
	FreeMemoryMiB   int64                 `json:"free_memory_mib"`
	Labels          map[string]string     `json:"labels"`
	Annotations     map[string]string     `json:"annotations"`
	Topology        map[string]any        `json:"topology"`
	ReportTime      string                `json:"report_time,omitempty"`
	GPUs            []AgentGPUDevice      `json:"gpus"`
	MIGs            []AgentMIGDevice      `json:"migs"`
	RuntimeBindings []AgentRuntimeBinding `json:"runtime_bindings"`
}

type AgentGPUDevice struct {
	UUID              string            `json:"uuid"`
	GPUIndex          int               `json:"gpu_index"`
	Model             string            `json:"model"`
	Vendor            string            `json:"vendor"`
	Type              string            `json:"type"`
	MemoryMiB         int64             `json:"memory_mib"`
	FreeMemoryMiB     int64             `json:"free_memory_mib"`
	Healthy           bool              `json:"healthy"`
	Health            string            `json:"health"`
	MIGEnabled        bool              `json:"mig_enabled"`
	MIGProfile        string            `json:"mig_profile"`
	UtilizationGPU    float64           `json:"utilization_gpu"`
	UtilizationMemory float64           `json:"utilization_memory"`
	Temperature       float64           `json:"temperature"`
	PowerWatts        float64           `json:"power_watts"`
	Allocated         bool              `json:"allocated"`
	Reserved          bool              `json:"reserved"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

type AgentMIGDevice struct {
	ParentGPUUUID string `json:"parent_gpu_uuid"`
	MIGUUID       string `json:"mig_uuid"`
	Profile       string `json:"profile"`
	MemoryMiB     int64  `json:"memory_mib"`
	Healthy       bool   `json:"healthy"`
	Allocated     bool   `json:"allocated"`
	Reserved      bool   `json:"reserved"`
}

type AgentRuntimeBinding struct {
	Namespace string   `json:"namespace"`
	PodName   string   `json:"pod_name"`
	GPUIDs    []string `json:"gpu_ids"`
}

type AgentReportResponse struct {
	NodeName   string `json:"node_name"`
	SnapshotID uint64 `json:"snapshot_id"`
	ReportTime string `json:"report_time"`
}
