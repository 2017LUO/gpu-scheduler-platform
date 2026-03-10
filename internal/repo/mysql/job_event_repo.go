package mysql

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type JobEventRepo struct {
	db *gorm.DB
}

func NewJobEventRepo(db *gorm.DB) *JobEventRepo {
	return &JobEventRepo{db: db}
}

func (r *JobEventRepo) Create(ctx context.Context, e *event.Event) error {
	if r == nil || r.db == nil || e == nil {
		return util.ErrInvalidArgument
	}
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(toEventModel(e)).Error; err != nil {
		return fmt.Errorf("create job event: %w", err)
	}
	return nil
}

func (r *JobEventRepo) ListByJobID(ctx context.Context, jobID string, limit, offset int) ([]event.Event, int64, error) {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return nil, 0, util.ErrInvalidArgument
	}

	dbq := dbFromContext(ctx, r.db).WithContext(ctx).Model(&models.GPUJobEvent{}).Where("job_id = ?", jobID)

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count job events: %w", err)
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	var ms []models.GPUJobEvent
	if err := dbq.Order("occurred_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("list job events: %w", err)
	}

	out := make([]event.Event, 0, len(ms))
	for i := range ms {
		out = append(out, toEventDomain(&ms[i]))
	}
	return out, total, nil
}
