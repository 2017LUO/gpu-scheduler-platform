package models

import "time"

type GPUJobEvent struct {
	ID         string    `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string    `gorm:"column:job_id;type:varchar(64);not null;index:idx_job_events_job"`
	TenantID   string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_job_events_tenant"`
	Reason     string    `gorm:"column:reason;type:varchar(64);not null"`
	Message    string    `gorm:"column:message;type:text"`
	Source     string    `gorm:"column:source;type:varchar(64);not null"`
	OccurredAt time.Time `gorm:"column:occurred_at;not null;index:idx_job_events_occurred"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (GPUJobEvent) TableName() string {
	return "gpu_job_events"
}
