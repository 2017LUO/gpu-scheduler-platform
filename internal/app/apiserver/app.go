package apiserver

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
	Config        *appcfg.APIServerConfig
	Logger        *zap.Logger
	DB            *gorm.DB
	Redis         *redis.Client
	Metrics       *obsmetrics.Registry
	HTTPServer    *http.Server
	Lifecycle     *bootstrap.Lifecycle
	TracingCloser bootstrap.TracingCloser
	readyChecks   []ReadyCheck
}

type ReadyCheck struct {
	Name string
	Fn   func(context.Context) error
}

func New(cfg *appcfg.APIServerConfig) (*App, error) {
	if cfg == nil {
		return nil, fmt.Errorf("api-server config is nil")
	}

	lg, err := bootstrap.NewLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	lg = obslog.WithComponent(lg, "api-server")

	db, err := bootstrap.NewMySQL(cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("init mysql: %w", err)
	}

	rdb, err := bootstrap.NewRedis(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("init redis: %w", err)
	}

	tracingCloser, err := bootstrap.InitTracing(cfg.Observability.Tracing)
	if err != nil {
		return nil, fmt.Errorf("init tracing: %w", err)
	}

	metricsReg := obsmetrics.NewRegistry()
	lifecycle := bootstrap.NewLifecycle(lg)

	app := &App{
		Config:        cfg,
		Logger:        lg,
		DB:            db,
		Redis:         rdb,
		Metrics:       metricsReg,
		Lifecycle:     lifecycle,
		TracingCloser: tracingCloser,
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
