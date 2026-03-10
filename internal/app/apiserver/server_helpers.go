package apiserver

import (
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
)

func shutdownHTTPServer(a *App, timeout time.Duration) error {
	return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, timeout)
}
