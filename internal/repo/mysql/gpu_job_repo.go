package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
)

type GPUJobFilter struct {
	TenantID  string
	Namespace string
	Queue     string
	Status    string
}

type GPUJobRepo struct {
	db *gorm.DB
}

func NewGPUJobRepo(db *gorm.DB) (*GPUJobRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &GPUJobRepo{db: db}, nil
}

func (r *GPUJobRepo) Create(ctx context.Context, m *model.GPUJob) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.TenantID == "" || m.Namespace == "" || m.Name == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create gpu job: %w", err)
	}
	return nil
}

func (r *GPUJobRepo) Get(ctx context.Context, id string) (*model.GPUJob, error) {
	if r == nil || r.db == nil || id == "" {
		return nil, ErrInvalidArgument
	}
	var m model.GPUJob
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *GPUJobRepo) List(ctx context.Context, f GPUJobFilter, page PageQuery) ([]model.GPUJob, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(50, 500)

	q := r.db.WithContext(ctx).Model(&model.GPUJob{})
	if f.TenantID != "" {
		q = q.Where("tenant_id = ?", f.TenantID)
	}
	if f.Namespace != "" {
		q = q.Where("namespace = ?", f.Namespace)
	}
	if f.Queue != "" {
		q = q.Where("queue = ?", f.Queue)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}

	var out []model.GPUJob
	if err := q.Order("created_at DESC").Limit(page.Limit).Offset(page.Offset).Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list gpu jobs: %w", err)
	}
	return out, nil
}

func (r *GPUJobRepo) CountByStatus(ctx context.Context, tenantID, status string) (int64, error) {
	if r == nil || r.db == nil || tenantID == "" || status == "" {
		return 0, ErrInvalidArgument
	}
	var n int64
	if err := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("tenant_id = ? AND status = ?", tenantID, status).
		Count(&n).Error; err != nil {
		return 0, fmt.Errorf("count gpu jobs by status: %w", err)
	}
	return n, nil
}

func (r *GPUJobRepo) UpdateStatus(ctx context.Context, id, status string, message *string) error {
	if r == nil || r.db == nil || id == "" || status == "" {
		return ErrInvalidArgument
	}
	updates := map[string]any{
		"status":     status,
		"message":    message,
		"updated_at": time.Now(),
	}
	res := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("id = ?", id).
		Updates(updates)
	if res.Error != nil {
		return fmt.Errorf("update gpu job status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) MarkScheduled(ctx context.Context, id string, at time.Time) error {
	if r == nil || r.db == nil || id == "" || at.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "ALLOCATED",
			"scheduled_at": at,
			"updated_at":   time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark gpu job scheduled: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) MarkRunning(ctx context.Context, id string, at time.Time) error {
	if r == nil || r.db == nil || id == "" || at.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     "RUNNING",
			"started_at": at,
			"updated_at": time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark gpu job running: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) MarkFinished(ctx context.Context, id, status string, at time.Time, message *string) error {
	if r == nil || r.db == nil || id == "" || status == "" || at.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":      status,
			"message":     message,
			"finished_at": at,
			"updated_at":  time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark gpu job finished: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) IncrementRetryCount(ctx context.Context, id string) error {
	if r == nil || r.db == nil || id == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.GPUJob{}).
		Where("id = ?", id).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1"))
	if res.Error != nil {
		return fmt.Errorf("increment gpu job retry count: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) GetByNamespaceAndName(ctx context.Context, namespace, name string) (*model.GPUJob, error) {
	if r == nil || r.db == nil || namespace == "" || name == "" {
		return nil, ErrInvalidArgument
	}

	var m model.GPUJob
	if err := r.db.WithContext(ctx).
		First(&m, "namespace = ? AND name = ?", namespace, name).Error; err != nil {
		return nil, mapDBError(err)
	}

	return &m, nil
}
