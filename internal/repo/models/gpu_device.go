package models

import "time"

type GPUDevice struct {
	ID            string `gorm:"column:id;type:varchar(64);primaryKey"`
	SnapshotID    uint64 `gorm:"column:snapshot_id;not null;index:idx_gpu_devices_snapshot"`
	NodeName      string `gorm:"column:node_name;type:varchar(128);not null;index:idx_gpu_devices_node"`
	UUID          string `gorm:"column:uuid;type:varchar(128);not null;index:idx_gpu_devices_uuid"`
	GPUIndex      int    `gorm:"column:gpu_index;not null"`
	Model         string `gorm:"column:model;type:varchar(128);not null"`
	Vendor        string `gorm:"column:vendor;type:varchar(64);not null;default:''"`
	Type          string `gorm:"column:type;type:varchar(16);not null"`
	MemoryMiB     int64  `gorm:"column:memory_mib;not null"`
	FreeMemoryMiB int64  `gorm:"column:free_memory_mib;not null"`

	Healthy    bool   `gorm:"column:healthy;not null"`
	Health     string `gorm:"column:health;type:varchar(32);not null"`
	MIGEnabled bool   `gorm:"column:mig_enabled;not null"`
	MIGProfile string `gorm:"column:mig_profile;type:varchar(64);not null;default:''"`

	LabelsJSON      string `gorm:"column:labels_json;type:json"`
	AnnotationsJSON string `gorm:"column:annotations_json;type:json"`

	Allocated bool `gorm:"column:allocated;not null"`
	Reserved  bool `gorm:"column:reserved;not null"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (GPUDevice) TableName() string {
	return "gpu_devices"
}
