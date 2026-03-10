package webhook

import (
	"context"
	"fmt"

	"gpu-scheduler-platform/internal/bootstrap"

	"go.uber.org/zap"
)

func (a *App) registerLifecycleHooks() {
	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		a.Logger.Info("webhook lifecycle start hook completed")
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if err := a.TracingCloser.Shutdown(ctx); err != nil {
			a.Logger.Warn("shutdown tracing failed", zap.Error(err))
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, a.Config.Server.HTTPS.ShutdownTimeout)
	})
}

func (a *App) TLSFilesReady() error {
	if a.Config == nil {
		return fmt.Errorf("config is nil")
	}
	if a.Config.Server.HTTPS.CertFile == "" {
		return fmt.Errorf("cert_file is required")
	}
	if a.Config.Server.HTTPS.KeyFile == "" {
		return fmt.Errorf("key_file is required")
	}
	return nil
}
