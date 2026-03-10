package scheduler

import (
	"time"

	appcfg "gpu-scheduler-platform/internal/config"
)

func (a *App) metricsHTTPConfig() appcfg.HTTPServerConfig {
	addr := a.Config.Observability.Metrics.Addr
	if addr == "" {
		addr = ":9091"
	}

	return appcfg.HTTPServerConfig{
		Addr:            addr,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 20 * time.Second,
	}
}
