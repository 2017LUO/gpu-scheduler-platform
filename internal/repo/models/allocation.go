package models

import "time"

type Allocation struct {
	ID         string `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_allocations_job"`
	TenantID   string `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_allocations_tenant"`
	NodeName   string `gorm:"column:node_name;type:varchar(128);not null;index:idx_allocations_node"`
	GPUIDsJSON string `gorm:"column:gpu_ids_json;type:json;not null"`
	Status     string `gorm:"column:status;type:varchar(32);not null;index:idx_allocations_status"`
	Message    string `gorm:"column:message;type:text"`

	CommittedAt *time.Time `gorm:"column:committed_at"`
	ReleasedAt  *time.Time `gorm:"column:released_at"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (Allocation) TableName() string {
	return "allocations"
}
