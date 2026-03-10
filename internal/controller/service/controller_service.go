package service

import (
	"context"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ControllerService struct {
	jobs         repo.JobRepository
	events       repo.JobEventRepository
	reservations repo.ReservationRepository
	allocations  repo.AllocationRepository
	txManager    repo.TxManager
	logger       *zap.Logger
	resyncPeriod time.Duration
}

func NewControllerService(
	jobs repo.JobRepository,
	events repo.JobEventRepository,
	reservations repo.ReservationRepository,
	allocations repo.AllocationRepository,
	txManager repo.TxManager,
	logger *zap.Logger,
	resyncPeriod time.Duration,
) *ControllerService {
	if resyncPeriod <= 0 {
		resyncPeriod = 30 * time.Second
	}
	return &ControllerService{
		jobs:         jobs,
		events:       events,
		reservations: reservations,
		allocations:  allocations,
		txManager:    txManager,
		logger:       logger,
		resyncPeriod: resyncPeriod,
	}
}

func (s *ControllerService) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.resyncPeriod)
	defer ticker.Stop()

	s.logger.Info("controller service started", zap.Duration("resync_period", s.resyncPeriod))

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("controller service stopped")
			return nil
		case <-ticker.C:
			if err := s.runOneCycle(ctx); err != nil {
				s.logger.Error("controller cycle failed", zap.Error(err))
			}
		}
	}
}

func (s *ControllerService) runOneCycle(ctx context.Context) error {
	now := time.Now().UTC()

	// 1. 清理过期 reservation
	deleted, err := s.reservations.DeleteExpired(ctx, now, 200)
	if err != nil {
		return fmt.Errorf("delete expired reservations: %w", err)
	}
	if deleted > 0 && s.logger != nil {
		s.logger.Info("deleted expired reservations", zap.Int64("count", deleted))
	}

	// 2. 对处于 Bound 的任务做最小状态检查占位
	items, err := s.jobs.List(ctx, repo.JobListFilter{
		Status: job.StatusBound,
		Limit:  200,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("list bound jobs: %w", err)
	}

	for _, j := range items {
		if err := s.reconcileBoundJob(ctx, j, now); err != nil && s.logger != nil {
			s.logger.Warn("reconcile bound job failed", zap.String("job_id", j.ID), zap.Error(err))
		}
	}

	return nil
}

func (s *ControllerService) reconcileBoundJob(ctx context.Context, j job.Job, now time.Time) error {
	_, err := s.allocations.GetByJobID(ctx, j.ID)
	if err != nil {
		return nil
	}

	// 这里先不推进到 Running，因为当前版本还没真正接 Pod/Binding。
	// 只做一次最小保活事件记录占位。
	return s.events.Create(ctx, &event.Event{
		ID:         uuid.NewString(),
		JobID:      j.ID,
		TenantID:   j.TenantID,
		Reason:     event.ReasonBound,
		Message:    "controller observed bound allocation",
		Source:     "controller",
		OccurredAt: now,
	})
}
