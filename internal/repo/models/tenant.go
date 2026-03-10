package models

import "time"

type Tenant struct {
	ID          string    `gorm:"column:id;type:varchar(64);primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(128);not null"`
	Enabled     bool      `gorm:"column:enabled;not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (Tenant) TableName() string {
	return "tenants"
}
