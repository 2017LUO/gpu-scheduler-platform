package scheduler

import (
	"context"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"
	schedservice "gpu-scheduler-platform/internal/scheduler/service"

	"go.uber.org/zap"
)

type Runner struct {
	cfg     *appcfg.SchedulerConfig
	logger  *zap.Logger
	service *schedservice.SchedulerService
}

func NewRunner(cfg *appcfg.SchedulerConfig, lg *zap.Logger, svc *schedservice.SchedulerService) *Runner {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &Runner{
		cfg:     cfg,
		logger:  lg,
		service: svc,
	}
}

func (r *Runner) Run(ctx context.Context) error {
	interval := r.cfg.Scheduler.ScheduleInterval
	if interval <= 0 {
		interval = 2 * time.Second
	}
	housekeepingInterval := 30 * time.Second

	scheduleTicker := time.NewTicker(interval)
	housekeepingTicker := time.NewTicker(housekeepingInterval)
	defer scheduleTicker.Stop()
	defer housekeepingTicker.Stop()

	r.logger.Info("scheduler runner started",
		zap.Duration("schedule_interval", interval),
		zap.Duration("housekeeping_interval", housekeepingInterval),
	)

	if err := r.service.WarmUp(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			r.logger.Info("scheduler runner stopped")
			return nil

		case <-scheduleTicker.C:
			if err := r.service.RunOnce(ctx); err != nil {
				r.logger.Warn("scheduler run once failed", zap.Error(err))
			}

		case <-housekeepingTicker.C:
			if err := r.service.RunHousekeeping(ctx); err != nil {
				r.logger.Warn("scheduler housekeeping failed", zap.Error(err))
			}
		}
	}
}
