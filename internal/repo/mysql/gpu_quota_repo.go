package mysql

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/domain/policy"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type GPUQuotaRepo struct {
	db *gorm.DB
}

func NewGPUQuotaRepo(db *gorm.DB) *GPUQuotaRepo {
	return &GPUQuotaRepo{db: db}
}

func (r *GPUQuotaRepo) Upsert(ctx context.Context, q *policy.Quota) error {
	if r == nil || r.db == nil || q == nil {
		return util.ErrInvalidArgument
	}
	m := toQuotaModel(q)

	var existing models.GPUQuota
	err := dbFromContext(ctx, r.db).WithContext(ctx).
		Where("tenant_id = ? AND namespace = ?", m.TenantID, m.Namespace).
		Take(&existing).Error

	switch {
	case err == nil:
		res := dbFromContext(ctx, r.db).WithContext(ctx).
			Model(&models.GPUQuota{}).
			Where("id = ?", existing.ID).
			Updates(map[string]any{
				"id":                 m.ID,
				"tenant_id":          m.TenantID,
				"namespace":          m.Namespace,
				"max_gpu_count":      m.MaxGPUCount,
				"max_running_jobs":   m.MaxRunningJobs,
				"max_queued_jobs":    m.MaxQueuedJobs,
				"max_gpu_memory_mib": m.MaxGPUMemoryMiB,
				"enabled":            m.Enabled,
			})
		if res.Error != nil {
			return fmt.Errorf("update gpu quota: %w", res.Error)
		}
		return nil

	case isNotFound(err):
		if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(m).Error; err != nil {
			return fmt.Errorf("create gpu quota: %w", err)
		}
		return nil

	default:
		return fmt.Errorf("lookup gpu quota: %w", err)
	}
}

func (r *GPUQuotaRepo) GetByTenant(ctx context.Context, tenantID string, namespace string) (*policy.Quota, error) {
	if r == nil || r.db == nil || strings.TrimSpace(tenantID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.GPUQuota
	if err := dbFromContext(ctx, r.db).WithContext(ctx).
		Where("tenant_id = ? AND namespace = ?", tenantID, namespace).
		Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get gpu quota by tenant")
	}
	return toQuotaDomain(&m), nil
}

func (r *GPUQuotaRepo) List(ctx context.Context, limit, offset int) ([]policy.Quota, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, util.ErrInvalidArgument
	}
	dbq := dbFromContext(ctx, r.db).WithContext(ctx).Model(&models.GPUQuota{})

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count gpu quotas: %w", err)
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

	var ms []models.GPUQuota
	if err := dbq.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("list gpu quotas: %w", err)
	}

	out := make([]policy.Quota, 0, len(ms))
	for i := range ms {
		out = append(out, *toQuotaDomain(&ms[i]))
	}
	return out, total, nil
}

func (r *GPUQuotaRepo) Delete(ctx context.Context, id string) error {
	if r == nil || r.db == nil || strings.TrimSpace(id) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).Delete(&models.GPUQuota{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete gpu quota: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}
