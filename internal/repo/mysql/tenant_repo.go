package mysql

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) Exists(ctx context.Context, tenantID string) (bool, error) {
	if r == nil || r.db == nil || strings.TrimSpace(tenantID) == "" {
		return false, util.ErrInvalidArgument
	}

	var count int64
	if err := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.Tenant{}).
		Where("id = ? AND enabled = ?", tenantID, true).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check tenant exists: %w", err)
	}
	return count > 0, nil
}
