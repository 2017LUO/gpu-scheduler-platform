package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type GPUJobEventRepo struct {
	db *gorm.DB
}

func NewGPUJobEventRepo(db *gorm.DB) (*GPUJobEventRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &GPUJobEventRepo{db: db}, nil
}

func (r *GPUJobEventRepo) Create(ctx context.Context, m *model.GPUJobEvent) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.JobID == "" || m.TenantID == "" || m.Reason == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create gpu job event: %w", err)
	}
	return nil
}

func (r *GPUJobEventRepo) ListByJob(ctx context.Context, jobID string, page PageQuery) ([]model.GPUJobEvent, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.GPUJobEvent
	if err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("occurred_at DESC, created_at DESC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list gpu job events by job: %w", err)
	}
	return out, nil
}
