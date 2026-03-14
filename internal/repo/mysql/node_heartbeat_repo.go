package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type NodeHeartbeatRepo struct {
	db *gorm.DB
}

func NewNodeHeartbeatRepo(db *gorm.DB) (*NodeHeartbeatRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &NodeHeartbeatRepo{db: db}, nil
}

func (r *NodeHeartbeatRepo) Upsert(ctx context.Context, m *model.NodeHeartbeat) error {
	if r == nil || r.db == nil || m == nil || m.NodeName == "" || m.Status == "" || m.LastSeenAt.IsZero() {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "node_name"}},
			DoUpdates: clause.Assignments(map[string]any{
				"status":       m.Status,
				"message":      m.Message,
				"last_seen_at": m.LastSeenAt,
				"updated_at":   time.Now(),
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("upsert node heartbeat: %w", err)
	}
	return nil
}

func (r *NodeHeartbeatRepo) Get(ctx context.Context, nodeName string) (*model.NodeHeartbeat, error) {
	if r == nil || r.db == nil || nodeName == "" {
		return nil, ErrInvalidArgument
	}
	var m model.NodeHeartbeat
	if err := r.db.WithContext(ctx).First(&m, "node_name = ?", nodeName).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *NodeHeartbeatRepo) ListStale(ctx context.Context, before time.Time, page PageQuery) ([]model.NodeHeartbeat, error) {
	if r == nil || r.db == nil || before.IsZero() {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.NodeHeartbeat
	if err := r.db.WithContext(ctx).
		Where("last_seen_at < ?", before).
		Order("last_seen_at ASC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list stale node heartbeats: %w", err)
	}
	return out, nil
}
