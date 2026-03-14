package dto

type ClusterNodeSummaryResponse struct {
	NodeName          string            `json:"node_name"`
	ClusterName       string            `json:"cluster_name"`
	Source            string            `json:"source"`
	State             string            `json:"state"`
	Schedulable       bool              `json:"schedulable"`
	GPUCount          int               `json:"gpu_count"`
	HealthyGPUCount   int               `json:"healthy_gpu_count"`
	TotalMemoryMiB    int64             `json:"total_memory_mib"`
	FreeMemoryMiB     int64             `json:"free_memory_mib"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	LastReportTime    *string           `json:"last_report_time,omitempty"`
	LastHeartbeatTime *string           `json:"last_heartbeat_time,omitempty"`
	CreatedAt         string            `json:"created_at"`
	UpdatedAt         string            `json:"updated_at"`
}

type ClusterHeartbeatResponse struct {
	Status     string `json:"status"`
	Message    string `json:"message"`
	LastSeenAt string `json:"last_seen_at"`
	UpdatedAt  string `json:"updated_at"`
}

type ClusterSnapshotResponse struct {
	ID              uint64            `json:"id"`
	Version         string            `json:"version"`
	AgentVersion    string            `json:"agent_version"`
	ClusterName     string            `json:"cluster_name"`
	NodeName        string            `json:"node_name"`
	Source          string            `json:"source"`
	NodeState       string            `json:"node_state"`
	Schedulable     bool              `json:"schedulable"`
	GPUCount        int               `json:"gpu_count"`
	HealthyGPUCount int               `json:"healthy_gpu_count"`
	TotalMemoryMiB  int64             `json:"total_memory_mib"`
	FreeMemoryMiB   int64             `json:"free_memory_mib"`
	Labels          map[string]string `json:"labels,omitempty"`
	Annotations     map[string]string `json:"annotations,omitempty"`
	Topology        map[string]any    `json:"topology,omitempty"`
	ReportTime      string            `json:"report_time"`
	CreatedAt       string            `json:"created_at"`
}

type ClusterGPUDeviceResponse struct {
	ID                uint64            `json:"id"`
	SnapshotID        uint64            `json:"snapshot_id"`
	NodeName          string            `json:"node_name"`
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
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
	Allocated         bool              `json:"allocated"`
	Reserved          bool              `json:"reserved"`
	CreatedAt         string            `json:"created_at"`
}

type ClusterMIGDeviceResponse struct {
	ID            uint64 `json:"id"`
	SnapshotID    uint64 `json:"snapshot_id"`
	NodeName      string `json:"node_name"`
	ParentGPUUUID string `json:"parent_gpu_uuid"`
	MIGUUID       string `json:"mig_uuid"`
	Profile       string `json:"profile"`
	MemoryMiB     int64  `json:"memory_mib"`
	Healthy       bool   `json:"healthy"`
	Allocated     bool   `json:"allocated"`
	Reserved      bool   `json:"reserved"`
	CreatedAt     string `json:"created_at"`
}

type ClusterRuntimeBindingResponse struct {
	ID         uint64   `json:"id"`
	SnapshotID uint64   `json:"snapshot_id"`
	NodeName   string   `json:"node_name"`
	Namespace  string   `json:"namespace"`
	PodName    string   `json:"pod_name"`
	GPUIDs     []string `json:"gpu_ids"`
	CreatedAt  string   `json:"created_at"`
}

type ClusterNodeDetailResponse struct {
	Node            ClusterNodeSummaryResponse      `json:"node"`
	Heartbeat       *ClusterHeartbeatResponse       `json:"heartbeat,omitempty"`
	LatestSnapshot  *ClusterSnapshotResponse        `json:"latest_snapshot,omitempty"`
	GPUDevices      []ClusterGPUDeviceResponse      `json:"gpu_devices,omitempty"`
	MIGDevices      []ClusterMIGDeviceResponse      `json:"mig_devices,omitempty"`
	RuntimeBindings []ClusterRuntimeBindingResponse `json:"runtime_bindings,omitempty"`
}
