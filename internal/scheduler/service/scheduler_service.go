package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	"gpu-scheduler-platform/internal/scheduler/algorithm"
	schedcache "gpu-scheduler-platform/internal/scheduler/cache"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
	filterplugin "gpu-scheduler-platform/internal/scheduler/plugins/filter"
	permitplugin "gpu-scheduler-platform/internal/scheduler/plugins/permit"
	reserveplugin "gpu-scheduler-platform/internal/scheduler/plugins/reserve"
	scoreplugin "gpu-scheduler-platform/internal/scheduler/plugins/score"
	schedqueue "gpu-scheduler-platform/internal/scheduler/queue"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SchedulerService struct {
	db               *gorm.DB
	logger           *zap.Logger
	repos            *repoimpl.Repos
	queue            schedqueue.Interface
	framework        *schedframework.Framework
	snapshotCache    *schedcache.SnapshotCache
	nodeCache        *schedcache.NodeCache
	jobCache         *schedcache.JobCache
	reservationCache *schedcache.ReservationCache

	placement    *PlacementService
	fairness     *FairnessService
	preemption   *PreemptionService
	housekeeping *HousekeepingService

	nowFunc          func() time.Time
	reservationTTL   time.Duration
	nodeScoreTopK    int
	enablePreemption bool
}

func NewSchedulerService(
	db *gorm.DB,
	lg *zap.Logger,
	queue schedqueue.Interface,
	fw *schedframework.Framework,
	snapshotCache *schedcache.SnapshotCache,
	nodeCache *schedcache.NodeCache,
	jobCache *schedcache.JobCache,
	reservationCache *schedcache.ReservationCache,
) (*SchedulerService, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}
	if lg == nil {
		lg = zap.NewNop()
	}
	if queue == nil {
		return nil, fmt.Errorf("queue is nil")
	}
	if fw == nil {
		return nil, fmt.Errorf("framework is nil")
	}

	repos, err := repoimpl.NewRepos(db)
	if err != nil {
		return nil, err
	}

	placement := NewPlacementService(repos, lg)
	fairness := NewFairnessService(repos, lg)
	preemption := NewPreemptionService(repos, lg)
	housekeeping := NewHousekeepingService(repos, lg)

	fw.AddFilter(filterplugin.NewResourceFit())
	fw.AddFilter(filterplugin.NewModelMatch())
	fw.AddFilter(filterplugin.NewMIGFit())
	fw.AddFilter(filterplugin.NewTopologyFit())

	fw.AddScore(scoreplugin.NewBinpack())
	fw.AddScore(scoreplugin.NewSpread())
	fw.AddScore(scoreplugin.NewTopologyScore())
	fw.AddScore(scoreplugin.NewFragmentationScore())
	fw.AddScore(scoreplugin.NewUtilizationScore())

	fw.AddReserve(reserveplugin.NewReservation(repos, reservationCache, 45*time.Second, lg))
	fw.AddPermit(permitplugin.NewGangPermit())

	return &SchedulerService{
		db:               db,
		logger:           lg,
		repos:            repos,
		queue:            queue,
		framework:        fw,
		snapshotCache:    snapshotCache,
		nodeCache:        nodeCache,
		jobCache:         jobCache,
		reservationCache: reservationCache,
		placement:        placement,
		fairness:         fairness,
		preemption:       preemption,
		housekeeping:     housekeeping,
		nowFunc:          func() time.Time { return time.Now().UTC() },
		reservationTTL:   45 * time.Second,
		nodeScoreTopK:    5,
		enablePreemption: true,
	}, nil
}

func (s *SchedulerService) WarmUp(ctx context.Context) error {
	if err := s.refreshNodes(ctx); err != nil {
		return err
	}
	if err := s.refreshQueuedJobs(ctx); err != nil {
		return err
	}
	if s.housekeeping != nil {
		if err := s.housekeeping.CleanupExpiredReservations(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (s *SchedulerService) RunHousekeeping(ctx context.Context) error {
	if s.housekeeping == nil {
		return nil
	}
	return s.housekeeping.CleanupExpiredReservations(ctx)
}

func (s *SchedulerService) RunOnce(ctx context.Context) error {
	if err := s.refreshNodes(ctx); err != nil {
		return err
	}
	if err := s.refreshQueuedJobs(ctx); err != nil {
		return err
	}

	job, ok := s.queue.Pop()
	if !ok {
		return nil
	}
	s.jobCache.Set(job)

	if _, err := s.repos.Allocations.GetByJobID(ctx, job.ID); err == nil {
		s.logger.Info("skip already allocated job", zap.String("job_id", job.ID))
		return nil
	} else if !errors.Is(err, repoimpl.ErrNotFound) {
		return err
	}

	snapshot := s.snapshotCache.Get()
	if snapshot == nil || len(snapshot.Nodes) == 0 {
		s.logger.Warn("skip scheduling because snapshot is empty", zap.String("job_id", job.ID))
		return nil
	}

	s.logger.Info("schedule job",
		zap.String("job_id", job.ID),
		zap.String("tenant_id", job.TenantID),
		zap.String("namespace", job.Namespace),
		zap.String("name", job.Name),
		zap.Int("candidate_nodes", len(snapshot.Nodes)),
	)

	outcome, err := algorithm.ScheduleOne(ctx, algorithm.Dependencies{
		DB:                s.db,
		Repos:             s.repos,
		Framework:         s.framework,
		ReservationCache:  s.reservationCache,
		Logger:            s.logger,
		Now:               s.nowFunc,
		ReservationTTL:    s.reservationTTL,
		TopK:              s.nodeScoreTopK,
		LoadNodeInventory: s.placement.LoadNodeInventory,
		SelectNodeGPUs:    s.placement.SelectNodeGPUs,
	}, job, snapshot.Nodes)
	if err != nil {
		return err
	}

	if outcome.NeedsPreemption && s.enablePreemption && s.preemption != nil {
		changed, msg, preemptErr := s.preemption.TryPreempt(ctx, job, snapshot.Nodes)
		if preemptErr != nil {
			return preemptErr
		}
		if changed {
			s.logger.Info("preemption triggered for pending job",
				zap.String("job_id", job.ID),
				zap.String("message", msg),
			)
		}
	}

	return nil
}

func (s *SchedulerService) refreshNodes(ctx context.Context) error {
	items, err := s.repos.Nodes.ListSchedulable(ctx, "READY", repoimpl.PageQuery{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("refresh nodes: %w", err)
	}

	s.nodeCache.Reset()
	nodes := make([]*model.Node, 0, len(items))

	for i := range items {
		item := items[i]
		s.nodeCache.Set(&item)
		cp := item
		nodes = append(nodes, &cp)
	}

	s.snapshotCache.Set(nodes)
	return nil
}

func (s *SchedulerService) refreshQueuedJobs(ctx context.Context) error {
	items, err := s.repos.GPUJobs.List(ctx, repoimpl.GPUJobFilter{
		Status: "QUEUED",
	}, repoimpl.PageQuery{
		Limit:  1000,
		Offset: 0,
	})
	if err != nil {
		return fmt.Errorf("refresh queued jobs: %w", err)
	}

	if s.fairness != nil {
		items = s.fairness.Reorder(ctx, items)
	}

	s.queue.Clear()
	s.jobCache.Reset()

	for i := range items {
		job := items[i]
		if err := s.queue.Push(&job); err != nil {
			s.logger.Warn("push queued job failed", zap.String("job_id", job.ID), zap.Error(err))
			continue
		}
	}
	return nil
}
