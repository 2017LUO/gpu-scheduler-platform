package mysql

import (
	"context"
	"fmt"
	model "gpu-scheduler-platform/internal/repo/models"
	"time"

	"gorm.io/gorm"
)

type ReservationRepo struct {
	db *gorm.DB
}

func NewReservationRepo(db *gorm.DB) (*ReservationRepo, error) {
	if db == nil {
		return nil, ErrNilDB
	}
	return &ReservationRepo{db: db}, nil
}

func (r *ReservationRepo) Create(ctx context.Context, m *model.Reservation) error {
	if r == nil || r.db == nil || m == nil || m.ID == "" || m.JobID == "" || m.NodeName == "" || m.ExpireAt.IsZero() {
		return ErrInvalidArgument
	}
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepo) GetByJobID(ctx context.Context, jobID string) (*model.Reservation, error) {
	if r == nil || r.db == nil || jobID == "" {
		return nil, ErrInvalidArgument
	}
	var m model.Reservation
	if err := r.db.WithContext(ctx).First(&m, "job_id = ?", jobID).Error; err != nil {
		return nil, mapDBError(err)
	}
	return &m, nil
}

func (r *ReservationRepo) ListExpired(ctx context.Context, now time.Time, page PageQuery) ([]model.Reservation, error) {
	if r == nil || r.db == nil || now.IsZero() {
		return nil, ErrInvalidArgument
	}
	page = page.Normalize(100, 1000)

	var out []model.Reservation
	if err := r.db.WithContext(ctx).
		Where("expire_at <= ?", now).
		Order("expire_at ASC").
		Limit(page.Limit).
		Offset(page.Offset).
		Find(&out).Error; err != nil {
		return nil, fmt.Errorf("list expired reservations: %w", err)
	}
	return out, nil
}

func (r *ReservationRepo) DeleteByJobID(ctx context.Context, jobID string) error {
	if r == nil || r.db == nil || jobID == "" {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Reservation{}, "job_id = ?", jobID)
	if res.Error != nil {
		return fmt.Errorf("delete reservation by job id: %w", res.Error)
	}
	return nil
}

func (r *ReservationRepo) DeleteExpired(ctx context.Context, now time.Time) error {
	if r == nil || r.db == nil || now.IsZero() {
		return ErrInvalidArgument
	}
	res := r.db.WithContext(ctx).Delete(&model.Reservation{}, "expire_at <= ?", now)
	if res.Error != nil {
		return fmt.Errorf("delete expired reservations: %w", res.Error)
	}
	return nil
}
