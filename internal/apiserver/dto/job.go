package dto

type CreateJobRequest struct {
	TenantID            string            `json:"tenant_id"`
	Namespace           string            `json:"namespace"`
	Name                string            `json:"name"`
	Queue               string            `json:"queue"`
	Priority            string            `json:"priority"`
	Submitter           string            `json:"submitter"`
	GPUCount            int               `json:"gpu_count"`
	GPUMemoryMiB        int64             `json:"gpu_memory_mib"`
	GPUModel            string            `json:"gpu_model"`
	RequireSameNode     bool              `json:"require_same_node"`
	RequireHealthy      bool              `json:"require_healthy"`
	RequireMIG          bool              `json:"require_mig"`
	ExpectedDurationSec int64             `json:"expected_duration_sec"`
	RunPolicy           map[string]any    `json:"run_policy"`
	PreferredNodeLabels map[string]string `json:"preferred_node_labels"`
	PreferredGPULabels  map[string]string `json:"preferred_gpu_labels"`
	Labels              map[string]string `json:"labels"`
	Annotations         map[string]string `json:"annotations"`
}

type JobResponse struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	Namespace       string  `json:"namespace"`
	Name            string  `json:"name"`
	Queue           string  `json:"queue"`
	Priority        string  `json:"priority"`
	Status          string  `json:"status"`
	Submitter       string  `json:"submitter"`
	GPUCount        int     `json:"gpu_count"`
	GPUMemoryMiB    int64   `json:"gpu_memory_mib"`
	GPUModel        string  `json:"gpu_model"`
	RequireSameNode bool    `json:"require_same_node"`
	RequireHealthy  bool    `json:"require_healthy"`
	RequireMIG      bool    `json:"require_mig"`
	RetryCount      int     `json:"retry_count"`
	Message         *string `json:"message,omitempty"`
	ScheduledAt     *string `json:"scheduled_at,omitempty"`
	StartedAt       *string `json:"started_at,omitempty"`
	FinishedAt      *string `json:"finished_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type JobEventResponse struct {
	ID         string  `json:"id"`
	JobID      string  `json:"job_id"`
	TenantID   string  `json:"tenant_id"`
	Reason     string  `json:"reason"`
	Message    *string `json:"message,omitempty"`
	Source     string  `json:"source"`
	OccurredAt string  `json:"occurred_at"`
	CreatedAt  string  `json:"created_at"`
}
