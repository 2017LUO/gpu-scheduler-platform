package models

import "time"

type GPUQuota struct {
	ID              string    `gorm:"column:id;type:varchar(64);primaryKey"`
	TenantID        string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_gpu_quotas_tenant"`
	Namespace       string    `gorm:"column:namespace;type:varchar(128);not null;default:'';index:idx_gpu_quotas_namespace"`
	MaxGPUCount     int       `gorm:"column:max_gpu_count;not null"`
	MaxRunningJobs  int       `gorm:"column:max_running_jobs;not null"`
	MaxQueuedJobs   int       `gorm:"column:max_queued_jobs;not null"`
	MaxGPUMemoryMiB int64     `gorm:"column:max_gpu_memory_mib;not null"`
	Enabled         bool      `gorm:"column:enabled;not null"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (GPUQuota) TableName() string {
	return "gpu_quotas"
}
