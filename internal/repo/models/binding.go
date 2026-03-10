package models

import "time"

type Binding struct {
	ID         string    `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string    `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_bindings_job"`
	NodeName   string    `gorm:"column:node_name;type:varchar(128);not null"`
	GPUIDsJSON string    `gorm:"column:gpu_ids_json;type:json;not null"`
	PodName    string    `gorm:"column:pod_name;type:varchar(128);not null;default:''"`
	Namespace  string    `gorm:"column:namespace;type:varchar(128);not null;default:''"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (Binding) TableName() string {
	return "bindings"
}
