package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"
	"gpu-scheduler-platform/internal/middleware"
	"gpu-scheduler-platform/internal/observability/profiling"
	"gpu-scheduler-platform/pkg/version"

	"go.uber.org/zap"
)

type healthResponse struct {
	Code      int               `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"request_id,omitempty"`
	Data      healthDataPayload `json:"data"`
}

type healthDataPayload struct {
	Status  string            `json:"status"`
	Time    string            `json:"time"`
	Version version.Info      `json:"version"`
	Details map[string]string `json:"details,omitempty"`
}

func (a *App) metricsHTTPConfig() appcfg.HTTPServerConfig {
	addr := a.Config.Observability.Metrics.Addr
	if addr == "" {
		addr = ":9094"
	}
	return appcfg.HTTPServerConfig{
		Addr:            addr,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		IdleTimeout:     60 * time.Second,
		ShutdownTimeout: 20 * time.Second,
	}
}

func (a *App) buildMux() *http.ServeMux {
	root := http.NewServeMux()

	root.Handle("/healthz", a.wrapMiddlewares(http.HandlerFunc(a.handleHealthz)))
	root.Handle("/readyz", a.wrapMiddlewares(http.HandlerFunc(a.handleReadyz)))

	metricsPath := a.Config.Observability.Metrics.Path
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	root.Handle(metricsPath, a.wrapMiddlewares(a.Metrics.Handler()))

	if a.Config.Observability.PProf.Enabled {
		profiling.Mount(root, a.Config.Observability.PProf.PathPrefix)
	}

	root.Handle("/", a.wrapMiddlewares(http.HandlerFunc(a.handleRoot)))
	return root
}

func (a *App) wrapMiddlewares(h http.Handler) http.Handler {
	return middleware.Chain(
		h,
		middleware.RequestID,
		middleware.Recovery(a.Logger),
		middleware.AccessLog(a.Logger, middleware.AccessLogConfig{Enabled: true}),
	)
}

func (a *App) handleRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: map[string]string{
				"service": a.Config.Service.Name,
				"node":    a.Config.Agent.NodeName,
			},
		},
	})
}

func (a *App) handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
		},
	})
}

func (a *App) handleReadyz(w http.ResponseWriter, r *http.Request) {
	ok, details := a.Readiness(r.Context())

	statusCode := http.StatusOK
	status := "ready"
	message := "ok"
	if !ok {
		statusCode = http.StatusServiceUnavailable
		status = "not_ready"
		message = "dependency not ready"
	}

	writeJSON(w, statusCode, healthResponse{
		Code:      0,
		Message:   message,
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  status,
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: details,
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (a *App) registerLifecycleHooks() {
	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		a.Logger.Info("agent lifecycle start hook completed")
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if err := a.TracingCloser.Shutdown(ctx); err != nil {
			a.Logger.Warn("shutdown tracing failed", zap.Error(err))
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		return bootstrap.ShutdownHTTPServer(a.Logger, a.HTTPServer, a.metricsHTTPConfig().ShutdownTimeout)
	})
}

func bootstrapRunHTTPServer(a *App, ctx context.Context) error {
	return bootstrap.RunHTTPServer(ctx, a.Logger, a.HTTPServer)
}
