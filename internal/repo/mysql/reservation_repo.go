package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type ReservationRepo struct {
	db *gorm.DB
}

func NewReservationRepo(db *gorm.DB) *ReservationRepo {
	return &ReservationRepo{db: db}
}

func (r *ReservationRepo) Create(ctx context.Context, rv *allocation.Reservation) error {
	if r == nil || r.db == nil || rv == nil {
		return util.ErrInvalidArgument
	}
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(toReservationModel(rv)).Error; err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	return nil
}

func (r *ReservationRepo) GetByJobID(ctx context.Context, jobID string) (*allocation.Reservation, error) {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.Reservation
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("job_id = ?", jobID).Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get reservation by job id")
	}
	return toReservationDomain(&m), nil
}

func (r *ReservationRepo) DeleteByJobID(ctx context.Context, jobID string) error {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).Delete(&models.Reservation{}, "job_id = ?", jobID)
	if res.Error != nil {
		return fmt.Errorf("delete reservation by job id: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (r *ReservationRepo) DeleteExpired(ctx context.Context, now time.Time, limit int) (int64, error) {
	if r == nil || r.db == nil {
		return 0, util.ErrInvalidArgument
	}
	if limit <= 0 {
		limit = 100
	}
	sub := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.Reservation{}).
		Select("id").
		Where("expire_at < ?", now).
		Limit(limit)

	res := dbFromContext(ctx, r.db).WithContext(ctx).Delete(&models.Reservation{}, "id IN (?)", sub)
	if res.Error != nil {
		return 0, fmt.Errorf("delete expired reservations: %w", res.Error)
	}
	return res.RowsAffected, nil
}
