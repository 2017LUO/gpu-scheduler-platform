package service

import (
	"context"
	"time"

	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type HousekeepingService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewHousekeepingService(repos *repoimpl.Repos, lg *zap.Logger) *HousekeepingService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &HousekeepingService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *HousekeepingService) CleanupExpiredReservations(ctx context.Context) error {
	now := s.nowFunc()
	items, err := s.repos.Reservations.ListExpired(ctx, now, repoimpl.PageQuery{Limit: 1000, Offset: 0})
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return nil
	}

	if err := s.repos.Reservations.DeleteExpired(ctx, now); err != nil {
		return err
	}

	s.logger.Info("expired reservations cleaned", zap.Int("count", len(items)))
	return nil
}
