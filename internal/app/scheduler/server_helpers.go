package scheduler

import (
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
)

func (a *App) shutdownHTTP(timeoutSeconds any) error {
	timeout, ok := timeoutSeconds.(time.Duration)
	if !ok {
		timeout = 20 * time.Second
	}
	return shutdownHTTPServer(a, timeout)
}

func shutdownHTTPServer(a *App, timeout time.Duration) error {
	return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, timeout)
}
