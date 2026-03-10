package apiserver

import (
	"net/http"

	apihandler "gpu-scheduler-platform/internal/apiserver/handler"
	apirouter "gpu-scheduler-platform/internal/apiserver/router"
	apiservice "gpu-scheduler-platform/internal/apiserver/service"
	"gpu-scheduler-platform/internal/middleware"
	"gpu-scheduler-platform/internal/observability/profiling"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"
)

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

	repos := repoimpl.NewRepositories(a.DB)
	services := apiservice.NewServices(repos)
	handlers := apihandler.NewHandlers(services, a.Logger)
	routes := apirouter.NewRoutes(handlers.Job, handlers.InternalAgent)
	routes.Register(root)

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
