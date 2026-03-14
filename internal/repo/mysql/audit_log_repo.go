package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"

	"gorm.io/gorm"
)

type AuditLogRepo struct {
	db *gorm.DB
}

func NewAuditLogRepo(db *gorm.DB) (*AuditLogRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &AuditLogRepo{db: db}, nil
}

func (r *AuditLogRepo) Create(ctx context.Context, m *model.AuditLog) error {
	if r == nil || r.db == nil || m == nil || m.Actor == "" || m.Action == "" || m.ResourceType == "" || m.ResourceID == "" {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	return nil
}

func (r *AuditLogRepo) ListByResource(ctx context.Context, resourceType, resourceID string, page PageQuery) ([]model.AuditLog, error) {
	if r == nil || r.db == nil || resourceType == "" || resourceID == "" {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.AuditLog
	if err := r.db.WithContext(ctx).
		Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).
		Order("created_at DESC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list audit logs by resource: %w", err)
	}
	return out, nil
}
