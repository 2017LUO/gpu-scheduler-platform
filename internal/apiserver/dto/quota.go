package dto

type UpsertQuotaRequest struct {
	TenantID        string `json:"tenant_id"`
	Namespace       string `json:"namespace"`
	MaxGPUCount     int    `json:"max_gpu_count"`
	MaxRunningJobs  int    `json:"max_running_jobs"`
	MaxQueuedJobs   int    `json:"max_queued_jobs"`
	MaxGPUMemoryMiB int64  `json:"max_gpu_memory_mib"`
	Enabled         *bool  `json:"enabled,omitempty"`
}

type QuotaResponse struct {
	ID              string `json:"id"`
	TenantID        string `json:"tenant_id"`
	Namespace       string `json:"namespace"`
	MaxGPUCount     int    `json:"max_gpu_count"`
	MaxRunningJobs  int    `json:"max_running_jobs"`
	MaxQueuedJobs   int    `json:"max_queued_jobs"`
	MaxGPUMemoryMiB int64  `json:"max_gpu_memory_mib"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}
