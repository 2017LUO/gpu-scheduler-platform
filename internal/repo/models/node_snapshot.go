package models

import "time"

type NodeSnapshot struct {
	ID          uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	Version     string `gorm:"column:version;type:varchar(64);not null;index:idx_node_snapshots_version"`
	NodeName    string `gorm:"column:node_name;type:varchar(128);not null;index:idx_node_snapshots_node"`
	NodeState   string `gorm:"column:node_state;type:varchar(32);not null"`
	Schedulable bool   `gorm:"column:schedulable;not null"`

	LabelsJSON      string `gorm:"column:labels_json;type:json"`
	AnnotationsJSON string `gorm:"column:annotations_json;type:json"`
	TopologyJSON    string `gorm:"column:topology_json;type:json"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime;index:idx_node_snapshots_created"`
}

func (NodeSnapshot) TableName() string {
	return "node_snapshots"
}
