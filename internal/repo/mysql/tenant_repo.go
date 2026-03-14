package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) (*TenantRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &TenantRepo{db: db}, nil
}

func (r *TenantRepo) Create(ctx context.Context, m *model.Tenant) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.Name == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create tenant: %w", err)
	}
	return nil
}

func (r *TenantRepo) Get(ctx context.Context, id string) (*model.Tenant, error) {
	if r == nil || r.db == nil || id == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Tenant
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *TenantRepo) List(ctx context.Context, enabled *bool, page PageQuery) ([]model.Tenant, error) {
	if r == nil || r.db == nil {
		return nil, ErrNilDB
	}
	page = page.Normalize(50, 500)

	q := r.db.WithContext(ctx).Model(&model.Tenant{})
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}

	var out []model.Tenant
	if err := q.Order("created_at DESC").Limit(page.Limit).Offset(page.Offset).Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list tenants: %w", err)
	}
	return out, nil
}

func (r *TenantRepo) Update(ctx context.Context, m *model.Tenant) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Tenant{}).
		Where("id = ?", m.ID).
		Updates(map[string]any{
			"name":        m.Name,
			"enabled":     m.Enabled,
			"description": m.Description,
		})
	if res.Error != nil {
		return fmt.Errorf("update tenant: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *TenantRepo) Delete(ctx context.Context, id string) error {
	if r == nil || r.db == nil || id == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Tenant{}, "id = ?", id)
	if res.Error != nil {
		return fmt.Errorf("delete tenant: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
