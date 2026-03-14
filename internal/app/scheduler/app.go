package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"
	obslog "gpu-scheduler-platform/internal/observability/logging"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"
	schedcache "gpu-scheduler-platform/internal/scheduler/cache"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
	schedqueue "gpu-scheduler-platform/internal/scheduler/queue"
	schedservice "gpu-scheduler-platform/internal/scheduler/service"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config        *appcfg.SchedulerConfig
	Logger        *zap.Logger
	DB            *gorm.DB
	Redis         *redis.Client
	Metrics       *obsmetrics.Registry
	HTTPServer    *http.Server
	Lifecycle     *bootstrap.Lifecycle
	TracingCloser bootstrap.TracingCloser

	Framework        *schedframework.Framework
	Queue            schedqueue.Interface
	SnapshotCache    *schedcache.SnapshotCache
	NodeCache        *schedcache.NodeCache
	JobCache         *schedcache.JobCache
	ReservationCache *schedcache.ReservationCache
	Service          *schedservice.SchedulerService
	Runner           *Runner

	readyChecks []ReadyCheck
}

type ReadyCheck struct {
	Name string
	Fn   func(context.Context) error
}

func New(cfg *appcfg.SchedulerConfig) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("scheduler config is nil")
	}

	lg, err := bootstrap.NewLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	lg = obslog.WithComponent(lg, "scheduler")

	db, err := bootstrap.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := bootstrap.NewRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	tracingCloser, err := bootstrap.InitTracing(cfg.Service, cfg.Observability.Tracing)
	if err != nil {
		return nil, fmt.Errorf("init tracing: %w", err)
	}

	metricsReg := obsmetrics.NewRegistry()
	lifecycle := bootstrap.NewLifecycle(lg)

	queue := schedqueue.NewPriorityQueue(4096)
	snapshotCache := schedcache.NewSnapshotCache()
	nodeCache := schedcache.NewNodeCache()
	jobCache := schedcache.NewJobCache()
	reservationCache := schedcache.NewReservationCache()

	registry := schedframework.NewRegistry()
	fw := schedframework.NewFramework(registry, lg)

	svc, err := schedservice.NewSchedulerService(
		db,
		lg,
		queue,
		fw,
		snapshotCache,
		nodeCache,
		jobCache,
		reservationCache,
	)
	if err != nil {
		return nil, fmt.Errorf("init scheduler service: %w", err)
	}

	runner := NewRunner(cfg, lg, svc)

	app := &App{
		Config:           cfg,
		Logger:           lg,
		DB:               db,
		Redis:            rdb,
		Metrics:          metricsReg,
		Lifecycle:        lifecycle,
		TracingCloser:    tracingCloser,
		Framework:        fw,
		Queue:            queue,
		SnapshotCache:    snapshotCache,
		NodeCache:        nodeCache,
		JobCache:         jobCache,
		ReservationCache: reservationCache,
		Service:          svc,
		Runner:           runner,
	}

	app.readyChecks = []ReadyCheck{
		{
			Name: "mysql",
			Fn: func(ctx context.Context) error {
				sqlDB, err := app.DB.DB()
				if err != nil {
					return err
				}
				return sqlDB.PingContext(ctx)
			},
		},
		{
			Name: "redis",
			Fn: func(ctx context.Context) error {
				return app.Redis.Ping(ctx).Err()
			},
		},
		{
			Name: "framework",
			Fn: func(context.Context) error {
				if app.Framework == nil {
					return fmt.Errorf("framework is nil")
				}
				return nil
			},
		},
		{
			Name: "queue",
			Fn: func(context.Context) error {
				if app.Queue == nil {
					return fmt.Errorf("queue is nil")
				}
				return nil
			},
		},
	}

	mux := app.buildMux()
	app.HTTPServer = bootstrap.NewHTTPServer(cfg.Server.HTTP, mux)

	app.registerLifecycleHooks()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	if err := a.Lifecycle.Start(ctx); err != nil {
		return err
	}
	return bootstrap.RunHTTPServer(ctx, a.Logger, a.HTTPServer)
}

func (a *App) Stop(ctx context.Context) error {
	return a.Lifecycle.Stop(ctx)
}

func (a *App) Readiness(ctx context.Context) (bool, map[string]string) {
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	details := make(map[string]string, len(a.readyChecks))
	ok := true
	for _, chk := range a.readyChecks {
		if err := chk.Fn(ctx); err != nil {
			ok = false
			details[chk.Name] = err.Error()
		} else {
			details[chk.Name] = "ok"
		}
	}
	return ok, details
}
