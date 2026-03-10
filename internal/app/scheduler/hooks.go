package scheduler

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

func (a *App) registerLifecycleHooks() {
	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		a.Logger.Info("scheduler lifecycle start hook completed")
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
		return shutdownHTTPServerWithConfig(a)
	})
}

func shutdownHTTPServerWithConfig(a *App) error {
	timeout := a.metricsHTTPConfig().ShutdownTimeout
	if timeout <= 0 {
		timeout = 20
	}
	return a.shutdownHTTP(timeout)
}
