package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
)

type AllocationRepo struct {
	db *gorm.DB
}

func NewAllocationRepo(db *gorm.DB) (*AllocationRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &AllocationRepo{db: db}, nil
}

func (r *AllocationRepo) Create(ctx context.Context, m *model.Allocation) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.JobID == "" || m.TenantID == "" || m.NodeName == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create allocation: %w", err)
	}
	return nil
}

func (r *AllocationRepo) GetByJobID(ctx context.Context, jobID string) (*model.Allocation, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Allocation
	if err := r.db.WithContext(ctx).First(&m, "job_id = ?", jobID).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *AllocationRepo) ListCommittedByNode(ctx context.Context, nodeName string, page PageQuery) ([]model.Allocation, error) {
	if r == nil || r.db == nil || nodeName == "" {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.Allocation
	if err := r.db.WithContext(ctx).
		Where("node_name = ? AND status = ?", nodeName, "COMMITTED").
		Order("created_at ASC, id ASC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list committed allocations by node: %w", err)
	}
	return out, nil
}

func (r *AllocationRepo) UpdateStatus(ctx context.Context, id, status string, message *string) error {
	if r == nil || r.db == nil || id == "" || status == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Allocation{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     status,
			"message":    message,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("update allocation status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AllocationRepo) MarkCommitted(ctx context.Context, id string, t time.Time) error {
	if r == nil || r.db == nil || id == "" || t.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Allocation{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "COMMITTED",
			"committed_at": t,
			"updated_at":   time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark allocation committed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AllocationRepo) MarkReleased(ctx context.Context, id string, t time.Time, message *string) error {
	if r == nil || r.db == nil || id == "" || t.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Allocation{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      "RELEASED",
			"released_at": t,
			"message":     message,
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark allocation released: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *AllocationRepo) DeleteByJobID(ctx context.Context, jobID string) error {
	if r == nil || r.db == nil || jobID == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Allocation{}, "job_id = ?", jobID)
	if res.Error != nil {
		return fmt.Errorf("delete allocation by job id: %w", res.Error)
	}
	return nil
}
