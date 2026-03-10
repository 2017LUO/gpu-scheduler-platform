package mysql

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/domain/allocation"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type BindingRepo struct {
	db *gorm.DB
}

func NewBindingRepo(db *gorm.DB) *BindingRepo {
	return &BindingRepo{db: db}
}

func (r *BindingRepo) Create(ctx context.Context, b *allocation.Binding) error {
	if r == nil || r.db == nil || b == nil {
		return util.ErrInvalidArgument
	}
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(toBindingModel(b)).Error; err != nil {
		return fmt.Errorf("create binding: %w", err)
	}
	return nil
}

func (r *BindingRepo) GetByJobID(ctx context.Context, jobID string) (*allocation.Binding, error) {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.Binding
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("job_id = ?", jobID).Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get binding by job id")
	}
	return toBindingDomain(&m), nil
}
