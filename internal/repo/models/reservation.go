package models

import "time"

type Reservation struct {
	ID         string    `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string    `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_reservations_job"`
	NodeName   string    `gorm:"column:node_name;type:varchar(128);not null;index:idx_reservations_node"`
	GPUIDsJSON string    `gorm:"column:gpu_ids_json;type:json;not null"`
	ExpireAt   time.Time `gorm:"column:expire_at;not null;index:idx_reservations_expire"`
	CreatedAt  time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (Reservation) TableName() string {
	return "reservations"
}
