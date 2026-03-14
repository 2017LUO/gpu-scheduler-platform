package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type JobRetryRepo struct {
	db *gorm.DB
}

func NewJobRetryRepo(db *gorm.DB) (*JobRetryRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &JobRetryRepo{db: db}, nil
}

func (r *JobRetryRepo) Create(ctx context.Context, m *model.JobRetry) error {
	if r == nil || r.db == nil || m == nil || m.JobID == "" || m.RetryNo <= 0 || m.Reason == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create job retry: %w", err)
	}
	return nil
}

func (r *JobRetryRepo) ListByJob(ctx context.Context, jobID string) ([]model.JobRetry, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	var out []model.JobRetry
	if err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("retry_no ASC, created_at ASC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list job retries by job: %w", err)
	}
	return out, nil
}
