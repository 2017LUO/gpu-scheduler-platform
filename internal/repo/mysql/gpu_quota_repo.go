package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GPUQuotaRepo struct {
	db *gorm.DB
}

func NewGPUQuotaRepo(db *gorm.DB) (*GPUQuotaRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &GPUQuotaRepo{db: db}, nil
}

func (r *GPUQuotaRepo) Upsert(ctx context.Context, m *model.GPUQuota) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.TenantID == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "namespace"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"max_gpu_count":      m.MaxGPUCount,
				"max_running_jobs":   m.MaxRunningJobs,
				"max_queued_jobs":    m.MaxQueuedJobs,
				"max_gpu_memory_mib": m.MaxGPUMemoryMiB,
				"enabled":            m.Enabled,
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("upsert gpu quota: %w", err)
	}
	return nil
}

func (r *GPUQuotaRepo) GetByTenantAndNamespace(ctx context.Context, tenantID, namespace string) (*model.GPUQuota, error) {
	if r == nil || r.db == nil || tenantID == "" {
		return nil, ErrInvalidArgument
	}
	var m model.GPUQuota
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND namespace = ?", tenantID, namespace).
		First(&m).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *GPUQuotaRepo) ListByTenant(ctx context.Context, tenantID string) ([]model.GPUQuota, error) {
	if r == nil || r.db == nil || tenantID == "" {
		return nil, ErrInvalidArgument
	}
	var out []model.GPUQuota
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("namespace ASC").
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list gpu quotas by tenant: %w", err)
	}
	return out, nil
}

func (r *GPUQuotaRepo) DeleteByTenantAndNamespace(ctx context.Context, tenantID, namespace string) error {
	if r == nil || r.db == nil || tenantID == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.GPUQuota{}, "tenant_id = ? AND namespace = ?", tenantID, namespace)
	if res.Error != nil {
		return fmt.Errorf("delete gpu quota: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
