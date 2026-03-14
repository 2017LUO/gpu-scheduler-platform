package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NodeRepo struct {
	db *gorm.DB
}

func NewNodeRepo(db *gorm.DB) (*NodeRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &NodeRepo{db: db}, nil
}

func (r *NodeRepo) Upsert(ctx context.Context, m *model.Node) error {
	if r == nil || r.db == nil || m == nil || m.NodeName == "" || m.State == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "node_name"}},
			DoUpdates: clause.Assignments(map[string]any{
				"cluster_name":        m.ClusterName,
				"source":              m.Source,
				"state":               m.State,
				"schedulable":         m.Schedulable,
				"gpu_count":           m.GPUCount,
				"healthy_gpu_count":   m.HealthyGPUCount,
				"total_memory_mib":    m.TotalMemoryMiB,
				"free_memory_mib":     m.FreeMemoryMiB,
				"labels_json":         m.LabelsJSON,
				"annotations_json":    m.AnnotationsJSON,
				"topology_json":       m.TopologyJSON,
				"last_report_time":    m.LastReportTime,
				"last_heartbeat_time": m.LastHeartbeatTime,
				"updated_at":          time.Now(),
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("upsert node: %w", err)
	}
	return nil
}

func (r *NodeRepo) UpdateHeartbeatTime(ctx context.Context, nodeName string, t time.Time) error {
	if r == nil || r.db == nil || nodeName == "" || t.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Node{}).
		Where("node_name = ?", nodeName).
		Updates(map[string]any{
			"last_heartbeat_time": t,
			"updated_at":          time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("update node heartbeat time: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *NodeRepo) Get(ctx context.Context, nodeName string) (*model.Node, error) {
	if r == nil || r.db == nil || nodeName == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Node
	if err := r.db.WithContext(ctx).First(&m, "node_name = ?", nodeName).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *NodeRepo) List(ctx context.Context, clusterName, state string, schedulable *bool, page PageQuery) ([]model.Node, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(100, 1000)

	q := r.db.WithContext(ctx).Model(&model.Node{})
	if clusterName != "" {
		q = q.Where("cluster_name = ?", clusterName)
	}
	if state != "" {
		q = q.Where("state = ?", state)
	}
	if schedulable != nil {
		q = q.Where("schedulable = ?", *schedulable)
	}

	var out []model.Node
	if err := q.
		Order("free_memory_mib DESC, healthy_gpu_count DESC, node_name ASC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list nodes: %w", err)
	}
	return out, nil
}

func (r *NodeRepo) ListSchedulable(ctx context.Context, state string, page PageQuery) ([]model.Node, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(100, 1000)

	q := r.db.WithContext(ctx).Model(&model.Node{}).Where("schedulable = ?", true)
	if state != "" {
		q = q.Where("state = ?", state)
	}

	var out []model.Node
	if err := q.Order("free_memory_mib DESC, healthy_gpu_count DESC").Limit(page.Limit).Offset(page.Offset).Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list schedulable nodes: %w", err)
	}
	return out, nil
}
