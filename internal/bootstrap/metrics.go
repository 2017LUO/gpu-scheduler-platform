package bootstrap

import (
	"net/http"

	appcfg "gpu-scheduler-platform/internal/config"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func BuildMetricsHandler(cfg appcfg.MetricsConfig) http.Handler {
	if !cfg.Enabled {
		return http.NotFoundHandler()
	}
	return promhttp.Handler()
}
