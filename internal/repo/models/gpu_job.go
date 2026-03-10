package models

import "time"

type GPUJob struct {
	ID        string `gorm:"column:id;type:varchar(64);primaryKey"`
	TenantID  string `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_jobs_tenant"`
	Namespace string `gorm:"column:namespace;type:varchar(128);not null;index:idx_jobs_namespace"`
	Name      string `gorm:"column:name;type:varchar(128);not null"`
	Queue     string `gorm:"column:queue;type:varchar(64);not null;index:idx_jobs_queue"`
	Priority  string `gorm:"column:priority;type:varchar(32);not null"`
	Status    string `gorm:"column:status;type:varchar(32);not null;index:idx_jobs_status"`

	GPUCount        int    `gorm:"column:gpu_count;not null"`
	GPUMemoryMiB    int64  `gorm:"column:gpu_memory_mib;not null"`
	GPUModel        string `gorm:"column:gpu_model;type:varchar(128);not null;default:''"`
	RequireSameNode bool   `gorm:"column:require_same_node;not null"`
	RequireHealthy  bool   `gorm:"column:require_healthy;not null"`
	RequireMIG      bool   `gorm:"column:require_mig;not null"`
	MIGProfile      string `gorm:"column:mig_profile;type:varchar(64);not null;default:''"`
	RequireNVLink   bool   `gorm:"column:require_nvlink;not null"`

	Preemptible      bool  `gorm:"column:preemptible;not null"`
	Retryable        bool  `gorm:"column:retryable;not null"`
	MaxRetry         int   `gorm:"column:max_retry;not null"`
	ExpectedDuration int64 `gorm:"column:expected_duration_sec;not null"`

	PreferredNodeLabelsJSON string `gorm:"column:preferred_node_labels_json;type:json"`
	PreferredGPULabelsJSON  string `gorm:"column:preferred_gpu_labels_json;type:json"`
	LabelsJSON              string `gorm:"column:labels_json;type:json"`
	AnnotationsJSON         string `gorm:"column:annotations_json;type:json"`

	RetryCount int    `gorm:"column:retry_count;not null"`
	Message    string `gorm:"column:message;type:text"`

	ScheduledAt *time.Time `gorm:"column:scheduled_at"`
	StartedAt   *time.Time `gorm:"column:started_at"`
	FinishedAt  *time.Time `gorm:"column:finished_at"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (GPUJob) TableName() string {
	return "gpu_jobs"
}
