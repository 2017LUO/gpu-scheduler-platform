package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type AllocationRepo struct {
	db *gorm.DB
}

func NewAllocationRepo(db *gorm.DB) *AllocationRepo {
	return &AllocationRepo{db: db}
}

func (r *AllocationRepo) Create(ctx context.Context, a *allocation.Allocation) error {
	if r == nil || r.db == nil || a == nil {
		return util.ErrInvalidArgument
	}
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(toAllocationModel(a)).Error; err != nil {
		return fmt.Errorf("create allocation: %w", err)
	}
	return nil
}

func (r *AllocationRepo) Update(ctx context.Context, a *allocation.Allocation) error {
	if r == nil || r.db == nil || a == nil || strings.TrimSpace(a.ID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.Allocation{}).
		Where("id = ?", a.ID).
		Updates(map[string]any{
			"job_id":       a.JobID,
			"tenant_id":    a.TenantID,
			"node_name":    a.NodeName,
			"gpu_ids_json": mustJSON(a.GPUIDs),
			"status":       string(a.Status),
			"message":      a.Message,
			"committed_at": a.CommittedAt,
			"released_at":  a.ReleasedAt,
		})
	if res.Error != nil {
		return fmt.Errorf("update allocation: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (r *AllocationRepo) UpdateStatus(ctx context.Context, allocationID string, status allocation.Status, message string) error {
	if r == nil || r.db == nil || strings.TrimSpace(allocationID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.Allocation{}).
		Where("id = ?", allocationID).
		Updates(map[string]any{
			"status":  string(status),
			"message": message,
		})
	if res.Error != nil {
		return fmt.Errorf("update allocation status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (r *AllocationRepo) GetByID(ctx context.Context, allocationID string) (*allocation.Allocation, error) {
	if r == nil || r.db == nil || strings.TrimSpace(allocationID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.Allocation
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", allocationID).Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get allocation by id")
	}
	return toAllocationDomain(&m), nil
}

func (r *AllocationRepo) GetByJobID(ctx context.Context, jobID string) (*allocation.Allocation, error) {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.Allocation
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("job_id = ?", jobID).Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get allocation by job id")
	}
	return toAllocationDomain(&m), nil
}

func (r *AllocationRepo) ListByNode(ctx context.Context, nodeName string, limit, offset int) ([]allocation.Allocation, int64, error) {
	if r == nil || r.db == nil || strings.TrimSpace(nodeName) == "" {
		return nil, 0, util.ErrInvalidArgument
	}
	dbq := dbFromContext(ctx, r.db).WithContext(ctx).Model(&models.Allocation{}).Where("node_name = ?", nodeName)

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count allocations by node: %w", err)
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

	var ms []models.Allocation
	if err := dbq.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("list allocations by node: %w", err)
	}

	out := make([]allocation.Allocation, 0, len(ms))
	for i := range ms {
		out = append(out, *toAllocationDomain(&ms[i]))
	}
	return out, total, nil
}

func (r *AllocationRepo) ReleaseByJobID(ctx context.Context, jobID string, releasedAt time.Time, message string) error {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.Allocation{}).
		Where("job_id = ?", jobID).
		Updates(map[string]any{
			"status":      string(allocation.StatusReleased),
			"released_at": releasedAt,
			"message":     message,
		})
	if res.Error != nil {
		return fmt.Errorf("release allocation by job id: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}
