package dto

type CreateJobRequest struct {
	TenantID    string                  `json:"tenant_id"`
	Namespace   string                  `json:"namespace"`
	Name        string                  `json:"name"`
	Queue       string                  `json:"queue"`
	Priority    string                  `json:"priority"`
	Requirement CreateJobRequirementDTO `json:"requirement"`
	Labels      map[string]string       `json:"labels,omitempty"`
	Annotations map[string]string       `json:"annotations,omitempty"`
}

type CreateJobRequirementDTO struct {
	GPUCount            int               `json:"gpu_count"`
	GPUMemoryMiB        int64             `json:"gpu_memory_mib"`
	GPUModel            string            `json:"gpu_model,omitempty"`
	RequireSameNode     bool              `json:"require_same_node"`
	RequireHealthy      bool              `json:"require_healthy"`
	RequireMIG          bool              `json:"require_mig"`
	MIGProfile          string            `json:"mig_profile,omitempty"`
	RequireNVLink       bool              `json:"require_nvlink"`
	PreferredNodeLabels map[string]string `json:"preferred_node_labels,omitempty"`
	PreferredGPULabels  map[string]string `json:"preferred_gpu_labels,omitempty"`
	Preemptible         bool              `json:"preemptible"`
	Retryable           bool              `json:"retryable"`
	MaxRetry            int               `json:"max_retry"`
	ExpectedDurationSec int64             `json:"expected_duration_sec"`
}

type JobResponse struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	Namespace   string            `json:"namespace"`
	Name        string            `json:"name"`
	Queue       string            `json:"queue"`
	Priority    string            `json:"priority"`
	Status      string            `json:"status"`
	Requirement JobRequirementDTO `json:"requirement"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	RetryCount  int               `json:"retry_count"`
	Message     string            `json:"message,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	ScheduledAt *string           `json:"scheduled_at,omitempty"`
	StartedAt   *string           `json:"started_at,omitempty"`
	FinishedAt  *string           `json:"finished_at,omitempty"`
}

type JobRequirementDTO struct {
	GPUCount            int               `json:"gpu_count"`
	GPUMemoryMiB        int64             `json:"gpu_memory_mib"`
	GPUModel            string            `json:"gpu_model,omitempty"`
	RequireSameNode     bool              `json:"require_same_node"`
	RequireHealthy      bool              `json:"require_healthy"`
	RequireMIG          bool              `json:"require_mig"`
	MIGProfile          string            `json:"mig_profile,omitempty"`
	RequireNVLink       bool              `json:"require_nvlink"`
	PreferredNodeLabels map[string]string `json:"preferred_node_labels,omitempty"`
	PreferredGPULabels  map[string]string `json:"preferred_gpu_labels,omitempty"`
	Preemptible         bool              `json:"preemptible"`
	Retryable           bool              `json:"retryable"`
	MaxRetry            int               `json:"max_retry"`
	ExpectedDurationSec int64             `json:"expected_duration_sec"`
}

type ListJobsResponse = PageResponse[JobResponse]

type JobEventResponse struct {
	ID         string `json:"id"`
	JobID      string `json:"job_id"`
	TenantID   string `json:"tenant_id"`
	Reason     string `json:"reason"`
	Message    string `json:"message,omitempty"`
	Source     string `json:"source"`
	OccurredAt string `json:"occurred_at"`
}

type ListJobEventsResponse = PageResponse[JobEventResponse]
