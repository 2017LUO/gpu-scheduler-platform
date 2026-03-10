package controller

import (
	"context"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
)

func (a *App) metricsHTTPConfig() appcfg.HTTPServerConfig {
	addr := a.Config.Observability.Metrics.Addr
	if addr == "" {
		addr = ":9092"
	}
	return appcfg.HTTPServerConfig{
		Addr:            addr,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 20 * time.Second,
	}
}

func (a *App) registerLifecycleHooks() {
	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		a.Logger.Info("controller lifecycle start hook completed")
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if err := a.TracingCloser.Shutdown(ctx); err != nil {
			a.Logger.Warn("shutdown tracing failed", zap.Error(err))
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if err := a.Redis.Close(); err != nil {
			return fmt.Errorf("close redis: %w", err)
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		sqlDB, err := a.DB.DB()
		if err != nil {
			return fmt.Errorf("get sql db: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("close mysql: %w", err)
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, a.metricsHTTPConfig().ShutdownTimeout)
	})
}

func (a *App) runWithLeaderElection(ctx context.Context) error {
	if a.LeaderElector == nil {
		return fmt.Errorf("leader elector is nil")
	}
	if a.Manager == nil {
		return fmt.Errorf("manager is nil")
	}

	a.Logger.Info("controller starting",
		zap.Bool("leader_election_enabled", a.Config.LeaderElection.Enabled),
		zap.String("lease_name", a.Config.LeaderElection.LeaseName),
		zap.String("lease_namespace", a.Config.LeaderElection.LeaseNamespace),
	)

	return a.LeaderElector.Run(ctx, func(runCtx context.Context) {
		a.Logger.Info("controller became leader")
		if err := a.Manager.Run(runCtx); err != nil {
			a.Logger.Error("controller manager exited with error", zap.Error(err))
		}
	})
}
