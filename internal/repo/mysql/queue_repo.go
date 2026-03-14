package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type QueueRepo struct {
	db *gorm.DB
}

func NewQueueRepo(db *gorm.DB) (*QueueRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &QueueRepo{db: db}, nil
}

func (r *QueueRepo) Create(ctx context.Context, m *model.Queue) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.Name == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create queue: %w", err)
	}
	return nil
}

func (r *QueueRepo) Upsert(ctx context.Context, m *model.Queue) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.Name == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "name"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"name":        m.Name,
				"tenant_id":   m.TenantID,
				"weight":      m.Weight,
				"priority":    m.Priority,
				"enabled":     m.Enabled,
				"description": m.Description,
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("upsert queue: %w", err)
	}
	return nil
}

func (r *QueueRepo) GetByID(ctx context.Context, id string) (*model.Queue, error) {
	if r == nil || r.db == nil || id == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Queue
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *QueueRepo) GetByName(ctx context.Context, tenantID, name string) (*model.Queue, error) {
	if r == nil || r.db == nil || name == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Queue
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, name).
		First(&m).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *QueueRepo) List(ctx context.Context, tenantID string, enabled *bool, page PageQuery) ([]model.Queue, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(50, 500)

	q := r.db.WithContext(ctx).Model(&model.Queue{})
	if tenantID != "" {
		q = q.Where("tenant_id = ?", tenantID)
	}
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}

	var out []model.Queue
	if err := q.Order("priority DESC, created_at DESC").Limit(page.Limit).Offset(page.Offset).Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list queues: %w", err)
	}
	return out, nil
}

func (r *QueueRepo) DeleteByTenantAndName(ctx context.Context, tenantID, name string) error {
	if r == nil || r.db == nil || name == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Queue{}, "tenant_id = ? AND name = ?", tenantID, name)
	if res.Error != nil {
		return fmt.Errorf("delete queue: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
