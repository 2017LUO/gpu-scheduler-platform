package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type GPUPolicyRepo struct {
	db *gorm.DB
}

func NewGPUPolicyRepo(db *gorm.DB) (*GPUPolicyRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &GPUPolicyRepo{db: db}, nil
}

func (r *GPUPolicyRepo) Upsert(ctx context.Context, m *model.GPUPolicy) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.TenantID == "" || m.Name == "" {
		return ErrInvalidArgument
	}

	if err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "tenant_id"},
				{Name: "name"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"queue":                     m.Queue,
				"priority":                  m.Priority,
				"enabled":                   m.Enabled,
				"preemptible":               m.Preemptible,
				"require_healthy":           m.RequireHealthy,
				"require_mig":               m.RequireMIG,
				"max_gpu_count":             m.MaxGPUCount,
				"required_gpu_model":        m.RequiredGPUModel,
				"required_node_labels_json": m.RequiredNodeLabelsJSON,
				"selector_json":             m.SelectorJSON,
				"description":               m.Description,
				"updated_at":                m.UpdatedAt,
			}),
		}).
		Create(m).Error; err != nil {
		return fmt.Errorf("upsert gpu policy: %w", err)
	}
	return nil
}

func (r *GPUPolicyRepo) GetByTenantAndName(ctx context.Context, tenantID, name string) (*model.GPUPolicy, error) {
	if r == nil || r.db == nil || tenantID == "" || name == "" {
		return nil, ErrInvalidArgument
	}

	var m model.GPUPolicy
	if err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND name = ?", tenantID, name).
		First(&m).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *GPUPolicyRepo) List(ctx context.Context, tenantID string, enabled *bool, page PageQuery) ([]model.GPUPolicy, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(50, 500)

	q := r.db.WithContext(ctx).Model(&model.GPUPolicy{})
	if tenantID != "" {
		q = q.Where("tenant_id = ?", tenantID)
	}
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}

	var out []model.GPUPolicy
	if err := q.
		Order("priority DESC, created_at DESC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list gpu policies: %w", err)
	}
	return out, nil
}

func (r *GPUPolicyRepo) DeleteByTenantAndName(ctx context.Context, tenantID, name string) error {
	if r == nil || r.db == nil || tenantID == "" || name == "" {
		return ErrInvalidArgument
	}

	res := r.db.WithContext(ctx).
		Delete(&model.GPUPolicy{}, "tenant_id = ? AND name = ?", tenantID, name)
	if res.Error != nil {
		return fmt.Errorf("delete gpu policy: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
