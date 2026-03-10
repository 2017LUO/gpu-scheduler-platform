package scheduler

import (
	"context"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"
	obslog "gpu-scheduler-platform/internal/observability/logging"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
	cachepkg "gpu-scheduler-platform/internal/scheduler/cache"
	svc "gpu-scheduler-platform/internal/scheduler/service"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Runner struct {
	cfg     *appcfg.SchedulerConfig
	logger  *zap.Logger
	service *svc.SchedulerService
}

func NewRunner(
	cfg *appcfg.SchedulerConfig,
	lg *zap.Logger,
	db *gorm.DB,
	rdb *redis.Client,
	_ any,
) *Runner {
	repos := repoimpl.NewRepositories(db)
	snapshotCache := cachepkg.NewSnapshotCache(repos.NodeSnapshots, 2*time.Second)
	reservationCache := cachepkg.NewReservationCache()

	fairness := svc.NewFairnessService()
	placement := svc.NewPlacementService()

	schedulerService := svc.NewSchedulerService(
		repos.Jobs,
		repos.JobEvents,
		repos.NodeSnapshots,
		repos.Allocations,
		repos.Reservations,
		repos.TxManager,
		snapshotCache,
		reservationCache,
		fairness,
		placement,
		lg,
		cfg.Scheduler.ReservationTTL,
		cfg.Scheduler.PendingBatchSize,
	)

	return &Runner{
		cfg:     cfg,
		logger:  obslog.LoggerOrNop(lg),
		service: schedulerService,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.cfg.Scheduler.ScheduleInterval)
	defer ticker.Stop()

	r.logger.Info("scheduler runner started",
		zap.Duration("schedule_interval", r.cfg.Scheduler.ScheduleInterval),
		zap.Int("pending_batch_size", r.cfg.Scheduler.PendingBatchSize),
	)

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("scheduler runner stopped")
			return nil
		case <-ticker.C:
			r.runOneCycle(ctx)
		}
	}
}

func (r *Runner) runOneCycle(ctx context.Context) {
	start := time.Now()

	result, err := r.service.RunOneCycle(ctx)
	if err != nil {
		r.logger.Error("scheduler cycle failed",
			zap.Error(err),
			zap.Duration("latency", time.Since(start)),
		)
		return
	}

	r.logger.Info("scheduler cycle completed",
		zap.Int("scanned", result.Scanned),
		zap.Int("scheduled", result.Scheduled),
		zap.Int("skipped", result.Skipped),
		zap.Int("failed", result.Failed),
		zap.Duration("latency", time.Since(start)),
	)
}
