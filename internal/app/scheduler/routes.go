package scheduler

import "net/http"

func (a *App) buildMux() *http.ServeMux {
	root := http.NewServeMux()
	root.Handle("/healthz", http.HandlerFunc(a.handleHealthz))
	root.Handle("/readyz", http.HandlerFunc(a.handleReadyz))

	metricsPath := a.Config.Observability.Metrics.Path
	if metricsPath == "" {
		metricsPath = "/metrics"
	}
	root.Handle(metricsPath, a.Metrics.Handler())

	root.Handle("/", http.HandlerFunc(a.handleRoot))
	return root
}
