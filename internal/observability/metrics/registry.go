package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Registry struct {
	prom *prometheus.Registry
}

func NewRegistry() *Registry {
	return &Registry{
		prom: prometheus.NewRegistry(),
	}
}

func (r *Registry) MustRegister(cs ...prometheus.Collector) {
	if r == nil || r.prom == nil {
		return
	}
	r.prom.MustRegister(cs...)
}

func (r *Registry) Register(cs ...prometheus.Collector) error {
	if r == nil || r.prom == nil {
		return nil
	}
	for _, c := range cs {
		if err := r.prom.Register(c); err != nil {
			return err
		}
	}
	return nil
}

func (r *Registry) Handler() http.Handler {
	if r == nil || r.prom == nil {
		return http.NotFoundHandler()
	}
	return promhttp.HandlerFor(r.prom, promhttp.HandlerOpts{})
}

func (r *Registry) Raw() *prometheus.Registry {
	if r == nil {
		return nil
	}
	return r.prom
}
