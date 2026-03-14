package scheduler

import (
	"context"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"

	"go.uber.org/zap"
)

func (a *App) registerLifecycleHooks() {
	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		go func() {
			leader := NewLeaderRunner(a.Config, a.Logger, a.Runner)
			if err := leader.Run(ctx); err != nil {
				a.Logger.Error("leader runner exited", zap.Error(err))
			}
		}()
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
		timeout := a.Config.Server.HTTP.ShutdownTimeout
		if timeout <= 0 {
			timeout = 20 * time.Second
		}
		return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, timeout)
	})
}
