package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type SchedulingAttemptRepo struct {
	db *gorm.DB
}

func NewSchedulingAttemptRepo(db *gorm.DB) (*SchedulingAttemptRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &SchedulingAttemptRepo{db: db}, nil
}

func (r *SchedulingAttemptRepo) Create(ctx context.Context, m *model.SchedulingAttempt) error {
	if r == nil || r.db == nil || m == nil || m.JobID == "" || m.TenantID == "" || m.AttemptNo <= 0 || m.Result == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create scheduling attempt: %w", err)
	}
	return nil
}

func (r *SchedulingAttemptRepo) ListByJob(ctx context.Context, jobID string, page PageQuery) ([]model.SchedulingAttempt, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.SchedulingAttempt
	if err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("attempt_no ASC, created_at ASC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list scheduling attempts by job: %w", err)
	}
	return out, nil
}
