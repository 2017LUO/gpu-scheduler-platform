package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
)

type OutboxRepo struct {
	db *gorm.DB
}

func NewOutboxRepo(db *gorm.DB) (*OutboxRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &OutboxRepo{db: db}, nil
}

func (r *OutboxRepo) Create(ctx context.Context, m *model.Outbox) error {
	if r == nil || r.db == nil || m == nil || m.Topic == "" || m.Status == "" || m.AvailableAt.IsZero() {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create outbox: %w", err)
	}
	return nil
}

func (r *OutboxRepo) ListAvailable(ctx context.Context, now time.Time, status string, limit int) ([]model.Outbox, error) {
	if r == nil || r.db == nil || now.IsZero() {
		return nil, ErrInvalidArgument
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	q := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("available_at <= ?", now)

	if status != "" {
		q = q.Where("status = ?", status)
	}

	var out []model.Outbox
	if err := q.Order("available_at ASC, id ASC").Limit(limit).Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list available outbox: %w", err)
	}
	return out, nil
}

func (r *OutboxRepo) MarkProcessed(ctx context.Context, id uint64, processedAt time.Time) error {
	if r == nil || r.db == nil || id == 0 || processedAt.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "PROCESSED",
			"processed_at": processedAt,
			"updated_at":   time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark outbox processed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *OutboxRepo) MarkFailed(ctx context.Context, id uint64, lastErr string, nextAvailableAt time.Time) error {
	if r == nil || r.db == nil || id == 0 || nextAvailableAt.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       "FAILED",
			"last_error":   lastErr,
			"retry_count":  gorm.Expr("retry_count + 1"),
			"available_at": nextAvailableAt,
			"updated_at":   time.Now(),
		})
	if res.Error != nil {
		return fmt.Errorf("mark outbox failed: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}
