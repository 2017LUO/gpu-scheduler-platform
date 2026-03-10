package models

import "time"

type AuditLog struct {
	ID           uint64    `gorm:"column:id;primaryKey;autoIncrement"`
	Actor        string    `gorm:"column:actor;type:varchar(128);not null"`
	Action       string    `gorm:"column:action;type:varchar(128);not null"`
	ResourceType string    `gorm:"column:resource_type;type:varchar(64);not null"`
	ResourceID   string    `gorm:"column:resource_id;type:varchar(128);not null"`
	RequestID    string    `gorm:"column:request_id;type:varchar(128);not null;default:''"`
	DetailJSON   string    `gorm:"column:detail_json;type:json"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
