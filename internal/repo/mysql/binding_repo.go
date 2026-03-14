package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type BindingRepo struct {
	db *gorm.DB
}

func NewBindingRepo(db *gorm.DB) (*BindingRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &BindingRepo{db: db}, nil
}

func (r *BindingRepo) Create(ctx context.Context, m *model.Binding) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.JobID == "" || m.NodeName == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create binding: %w", err)
	}
	return nil
}

func (r *BindingRepo) GetByJobID(ctx context.Context, jobID string) (*model.Binding, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Binding
	if err := r.db.WithContext(ctx).First(&m, "job_id = ?", jobID).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *BindingRepo) DeleteByJobID(ctx context.Context, jobID string) error {
	if r == nil || r.db == nil || jobID == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Binding{}, "job_id = ?", jobID)
	if res.Error != nil {
		return fmt.Errorf("delete binding by job id: %w", res.Error)
	}
	return nil
}
