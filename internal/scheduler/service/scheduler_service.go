package service

import (
	"context"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/repo"
	"gpu-scheduler-platform/internal/scheduler/algorithm"
	cachepkg "gpu-scheduler-platform/internal/scheduler/cache"
	"gpu-scheduler-platform/internal/util"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SchedulerService struct {
	jobs         repo.JobRepository
	events       repo.JobEventRepository
	snapshots    repo.NodeSnapshotRepository
	allocations  repo.AllocationRepository
	reservations repo.ReservationRepository
	txManager    repo.TxManager

	snapshotCache    *cachepkg.SnapshotCache
	reservationCache *cachepkg.ReservationCache

	fairness  *FairnessService
	placement *PlacementService

	logger         *zap.Logger
	reservationTTL time.Duration
	batchSize      int
}

type SchedulerCycleResult struct {
	Scanned   int
	Scheduled int
	Skipped   int
	Failed    int
}

func NewSchedulerService(
	jobs repo.JobRepository,
	events repo.JobEventRepository,
	snapshots repo.NodeSnapshotRepository,
	allocations repo.AllocationRepository,
	reservations repo.ReservationRepository,
	txManager repo.TxManager,
	snapshotCache *cachepkg.SnapshotCache,
	reservationCache *cachepkg.ReservationCache,
	fairness *FairnessService,
	placement *PlacementService,
	logger *zap.Logger,
	reservationTTL time.Duration,
	batchSize int,
) *SchedulerService {
	if reservationTTL <= 0 {
		reservationTTL = 30 * time.Second
	}
	if batchSize <= 0 {
		batchSize = 50
	}
	return &SchedulerService{
		jobs:             jobs,
		events:           events,
		snapshots:        snapshots,
		allocations:      allocations,
		reservations:     reservations,
		txManager:        txManager,
		snapshotCache:    snapshotCache,
		reservationCache: reservationCache,
		fairness:         fairness,
		placement:        placement,
		logger:           logger,
		reservationTTL:   reservationTTL,
		batchSize:        batchSize,
	}
}

func (s *SchedulerService) RunOneCycle(ctx context.Context) (*SchedulerCycleResult, error) {
	if s == nil || s.jobs == nil || s.snapshots == nil || s.allocations == nil || s.reservations == nil || s.txManager == nil {
		return nil, util.ErrUnavailable
	}

	result := &SchedulerCycleResult{}

	snapshot, err := s.snapshotCache.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("load latest snapshot: %w", err)
	}

	pendingJobs, err := s.jobs.ListPending(ctx, repo.PendingJobFilter{Limit: s.batchSize})
	if err != nil {
		return nil, fmt.Errorf("list pending jobs: %w", err)
	}
	result.Scanned = len(pendingJobs)

	ordered := s.fairness.OrderJobs(pendingJobs)
	now := time.Now().UTC()

	for _, j := range ordered {
		if s.reservationCache.Exists(j.ID, now) {
			result.Skipped++
			continue
		}

		ok, err := s.scheduleOne(ctx, snapshot, j, now)
		if err != nil {
			result.Failed++
			if s.logger != nil {
				s.logger.Warn("schedule job failed", zap.String("job_id", j.ID), zap.Error(err))
			}
			continue
		}
		if !ok {
			result.Skipped++
			continue
		}
		result.Scheduled++
	}

	return result, nil
}

func (s *SchedulerService) scheduleOne(ctx context.Context, snapshot *cluster.Snapshot, j job.Job, now time.Time) (bool, error) {
	decision, err := s.placement.Place(ctx, snapshot, j)
	if err != nil {
		_ = s.createEventBestEffort(ctx, event.Event{
			ID:         uuid.NewString(),
			JobID:      j.ID,
			TenantID:   j.TenantID,
			Reason:     event.ReasonSchedulingFailed,
			Message:    err.Error(),
			Source:     "scheduler",
			OccurredAt: now,
		})
		return false, nil
	}

	reservation := algorithm.BuildReservation(decision, s.reservationTTL, now)
	allocationRecord := algorithm.BuildAllocation(decision, now)

	err = s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.jobs.UpdateStatus(txCtx, j.ID, job.StatusScheduling, "scheduling started"); err != nil {
			return err
		}

		if err := s.reservations.Create(txCtx, reservation); err != nil {
			return err
		}
		if err := s.allocations.Create(txCtx, allocationRecord); err != nil {
			return err
		}

		j.Status = job.StatusBound
		j.Message = "allocation committed"
		j.ScheduledAt = &now
		j.UpdatedAt = now

		if err := s.jobs.Update(txCtx, &j); err != nil {
			return err
		}

		if err := s.events.Create(txCtx, &event.Event{
			ID:         uuid.NewString(),
			JobID:      j.ID,
			TenantID:   j.TenantID,
			Reason:     event.ReasonSchedulingSucceeded,
			Message:    fmt.Sprintf("scheduled to node=%s gpus=%d", decision.Node.Name, len(decision.GPUs)),
			Source:     "scheduler",
			OccurredAt: now,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	s.reservationCache.Put(j.ID, reservation.ExpireAt)
	s.snapshotCache.Invalidate()
	return true, nil
}

func (s *SchedulerService) createEventBestEffort(ctx context.Context, e event.Event) error {
	if s.events == nil {
		return nil
	}
	return s.events.Create(ctx, &e)
}
