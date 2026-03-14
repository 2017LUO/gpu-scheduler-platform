package apiserver

import (
	"fmt"
	"net/http"

	apihandler "gpu-scheduler-platform/internal/apiserver/handler"
	apirouter "gpu-scheduler-platform/internal/apiserver/router"
	apiservice "gpu-scheduler-platform/internal/apiserver/service"
	"gpu-scheduler-platform/internal/middleware"
	"gpu-scheduler-platform/internal/observability/profiling"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
)

func (a *App) buildMux() (*http.ServeMux, error) {
	root := http.NewServeMux()

	root.Handle("/healthz", a.wrapMiddlewares(http.HandlerFunc(a.handleHealthz)))
	root.Handle("/readyz", a.wrapMiddlewares(http.HandlerFunc(a.handleReadyz)))

	metricsPath := a.Config.Observability.Metrics.Path
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	root.Handle(metricsPath, a.wrapMiddlewares(a.Metrics.Handler()))

	if a.Config.Observability.PProf.Enabled {
		profiling.RegisterPprofRoutes(root, a.Config.Observability.PProf)
	}

	repos, err := repoimpl.NewRepos(a.DB)
	if err != nil {
		return nil, fmt.Errorf("init repos: %w", err)
	}

	services, err := apiservice.NewServices(repos, a.Logger)
	if err != nil {
		return nil, fmt.Errorf("init services: %w", err)
	}

	handlers := apihandler.NewHandlers(services, a.Logger)
	routes := apirouter.NewRoutes(handlers, a.Config.Features)
	routes.Register(root)

	root.Handle("/", a.wrapMiddlewares(http.HandlerFunc(a.handleRoot)))
	return root, nil
}

func (a *App) wrapMiddlewares(h http.Handler) http.Handler {
	metricsPath := a.Config.Observability.Metrics.Path
	if metricsPath == "" {
		metricsPath = "/metrics"
	}

	pprofPrefix := a.Config.Observability.PProf.PathPrefix
	if pprofPrefix == "" {
		pprofPrefix = "/debug/pprof"
	}

	publicExact := []string{
		"/",
		"/healthz",
		"/readyz",
		metricsPath,
	}
	publicPrefixes := []string{
		pprofPrefix,
		"/internal/agent/",
	}

	return middleware.Chain(
		h,
		middleware.RequestID,
		middleware.Recovery(a.Logger),
		middleware.AccessLog(a.Logger, middleware.AccessLogConfig{Enabled: true}),
		middleware.AuthN(middleware.AuthNConfig{
			Enabled:          a.Config.Security.EnableAuthN,
			Manager:          a.JWTManager,
			PublicExactPaths: publicExact,
			PublicPrefixes:   publicPrefixes,
		}),
		middleware.AuthZ(middleware.AuthZConfig{
			Enabled:          a.Config.Security.EnableAuthZ,
			Authorizer:       a.Authorizer,
			PublicExactPaths: publicExact,
			PublicPrefixes:   publicPrefixes,
		}),
		middleware.Audit(middleware.AuditConfig{
			Enabled:         true,
			Recorder:        a.AuditRecorder,
			SkipExactPaths:  publicExact,
			SkipPrefixes:    publicPrefixes,
			OnlyMutatingAPI: true,
		}),
	)
}
