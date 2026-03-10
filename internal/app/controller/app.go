package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"
	obslog "gpu-scheduler-platform/internal/observability/logging"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	Config        *appcfg.ControllerAppConfig
	Logger        *zap.Logger
	DB            *gorm.DB
	Redis         *redis.Client
	K8s           *bootstrap.K8sClients
	Metrics       *obsmetrics.Registry
	HTTPServer    *http.Server
	Lifecycle     *bootstrap.Lifecycle
	TracingCloser bootstrap.TracingCloser
	LeaderElector bootstrap.LeaderElector
	Manager       *Manager
	readyChecks   []ReadyCheck
}

type ReadyCheck struct {
	Name string
	Fn   func(context.Context) error
}

func New(cfg *appcfg.ControllerAppConfig) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("controller config is nil")
	}

	lg, err := bootstrap.NewLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	lg = obslog.WithComponent(lg, "controller")

	db, err := bootstrap.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := bootstrap.NewRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	k8sClients, err := bootstrap.NewKubernetesClients(cfg.Kubernetes)
	if err != nil {
		return nil, fmt.Errorf("init kubernetes clients: %w", err)
	}

	tracingCloser, err := bootstrap.InitTracing(cfg.Observability.Tracing)
	if err != nil {
		return nil, fmt.Errorf("init tracing: %w", err)
	}

	metricsReg := obsmetrics.NewRegistry()
	lifecycle := bootstrap.NewLifecycle(lg)
	leader := bootstrap.NewLeaderElector(cfg.LeaderElection, lg)

	app := &App{
		Config:        cfg,
		Logger:        lg,
		DB:            db,
		Redis:         rdb,
		K8s:           k8sClients,
		Metrics:       metricsReg,
		Lifecycle:     lifecycle,
		TracingCloser: tracingCloser,
		LeaderElector: leader,
	}

	app.Manager = NewManager(cfg, lg, db, rdb, k8sClients)

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
	}

	mux := app.buildMux()
	app.HTTPServer = bootstrap.NewHTTPServer(app.metricsHTTPConfig(), mux)
	app.registerLifecycleHooks()

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	if err := a.Lifecycle.Start(ctx); err != nil {
		return err
	}
	return a.runWithLeaderElection(ctx)
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
