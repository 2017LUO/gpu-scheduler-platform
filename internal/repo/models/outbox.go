package models

import "time"

type Outbox struct {
	ID          uint64     `gorm:"column:id;primaryKey;autoIncrement"`
	Topic       string     `gorm:"column:topic;type:varchar(128);not null;index:idx_outbox_topic"`
	EventKey    string     `gorm:"column:event_key;type:varchar(128);not null;default:''"`
	PayloadJSON string     `gorm:"column:payload_json;type:json;not null"`
	Status      string     `gorm:"column:status;type:varchar(32);not null;index:idx_outbox_status"`
	AvailableAt time.Time  `gorm:"column:available_at;not null;index:idx_outbox_available"`
	ProcessedAt *time.Time `gorm:"column:processed_at"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (Outbox) TableName() string {
	return "outbox"
}
