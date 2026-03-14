package agent

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gpu-scheduler-platform/internal/bootstrap"
	appcfg "gpu-scheduler-platform/internal/config"
	"gpu-scheduler-platform/internal/middleware"
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
	addr := ":9094"
	if a != nil && a.Config != nil && a.Config.Observability.Server.Addr != "" {
		addr = a.Config.Observability.Server.Addr
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

	metricsPath := "/metrics"
	if a != nil && a.Config != nil && a.Config.Observability.Metrics.Path != "" {
		metricsPath = a.Config.Observability.Metrics.Path
	}

	metricsHandler := http.NotFoundHandler()
	if a != nil && a.Metrics != nil {
		metricsHandler = a.Metrics.Handler()
	}
	root.Handle(metricsPath, a.wrapMiddlewares(metricsHandler))

	if a != nil && a.Config != nil && a.Config.Observability.PProf.Enabled {
		bootstrap.RegisterPprof(root, a.Config.Observability.PProf)
	}

	root.Handle("/", a.wrapMiddlewares(http.HandlerFunc(a.handleRoot)))
	return root
}

func (a *App) wrapMiddlewares(h http.Handler) http.Handler {
	if h == nil {
		h = http.NotFoundHandler()
	}

	lg := zap.NewNop()
	if a != nil && a.Logger != nil {
		lg = a.Logger
	}

	return middleware.Chain(
		h,
		middleware.RequestID,
		middleware.Recovery(lg),
		middleware.AccessLog(lg, middleware.AccessLogConfig{Enabled: true}),
	)
}

func (a *App) handleRoot(w http.ResponseWriter, r *http.Request) {
	serviceName := "gpu-agent"
	nodeName := ""

	if a != nil && a.Config != nil {
		if a.Config.Service.Name != "" {
			serviceName = a.Config.Service.Name
		}
		nodeName = a.Config.Agent.NodeName
	}

	writeJSON(w, http.StatusOK, healthResponse{
		Code:      0,
		Message:   "ok",
		RequestID: middleware.RequestIDFromContext(r.Context()),
		Data: healthDataPayload{
			Status:  "ok",
			Time:    time.Now().UTC().Format(time.RFC3339Nano),
			Version: version.Get(),
			Details: map[string]string{
				"service": serviceName,
				"node":    nodeName,
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
	ok := true
	details := map[string]string{}

	if a != nil {
		ok, details = a.Readiness(r.Context())
	}

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
	if w == nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func (a *App) registerLifecycleHooks() {
	if a == nil || a.Lifecycle == nil {
		return
	}

	lg := zap.NewNop()
	if a.Logger != nil {
		lg = a.Logger
	}

	a.Lifecycle.AppendOnStart(func(ctx context.Context) error {
		lg.Info("agent lifecycle start hook completed")
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if a.TracingCloser == nil {
			return nil
		}
		if err := a.TracingCloser.Shutdown(ctx); err != nil {
			lg.Warn("shutdown tracing failed", zap.Error(err))
		}
		return nil
	})

	a.Lifecycle.AppendOnStop(func(ctx context.Context) error {
		if a.HTTPServer == nil {
			return nil
		}
		return bootstrap.ShutdownHTTPServer(lg, a.HTTPServer, a.metricsHTTPConfig().ShutdownTimeout)
	})
}

func bootstrapRunHTTPServer(a *App, ctx context.Context) error {
	if a == nil || a.HTTPServer == nil {
		return nil
	}

	lg := zap.NewNop()
	if a.Logger != nil {
		lg = a.Logger
	}
	return bootstrap.RunHTTPServer(ctx, lg, a.HTTPServer)
}
