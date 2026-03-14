package bootstrap

import (
	"net/http"

	appcfg "gpu-scheduler-platform/internal/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsRegistry struct {
	Registry *prometheus.Registry
	Handler  http.Handler
}

func NewMetricsRegistry() *MetricsRegistry {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &MetricsRegistry{
		Registry: reg,
		Handler:  promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
	}
}

func BuildMetricsHandler(cfg appcfg.MetricsConfig, mr *MetricsRegistry) http.Handler {
	if !cfg.Enabled {
		return http.NotFoundHandler()
	}
	if mr == nil || mr.Handler == nil {
		return promhttp.Handler()
	}
	return mr.Handler
}
