package metrics

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Registry：统一 Prometheus registry 封装。
// 说明：
// 1. 只保留这一套 metrics registry
// 2. app 层通过 Handler() 暴露 /metrics
// 3. 业务指标通过 MustRegister / RegisterCollector 注册
type Registry struct {
	prom    *prometheus.Registry
	handler http.Handler

	once sync.Once

	agent *AgentMetrics
}

// NewRegistry 创建统一指标注册中心，并默认注册 Go / Process 指标。
func NewRegistry() *Registry {
	r := prometheus.NewRegistry()

	r.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &Registry{
		prom:    r,
		handler: promhttp.HandlerFor(r, promhttp.HandlerOpts{}),
	}
}

// PrometheusRegistry 暴露底层 registry，供少数高级场景使用。
func (r *Registry) PrometheusRegistry() *prometheus.Registry {
	if r == nil {
		return nil
	}
	return r.prom
}

// Handler 返回 /metrics handler。
func (r *Registry) Handler() http.Handler {
	if r == nil || r.handler == nil {
		return http.NotFoundHandler()
	}
	return r.handler
}

// MustRegister 注册 collector，失败时 panic。
// 适合启动阶段固定注册。
func (r *Registry) MustRegister(cs ...prometheus.Collector) {
	if r == nil || r.prom == nil {
		panic("metrics registry is nil")
	}
	r.prom.MustRegister(cs...)
}

// RegisterCollector 注册 collector，返回 error。
// 适合动态注册或测试使用。
func (r *Registry) RegisterCollector(c prometheus.Collector) error {
	if r == nil || r.prom == nil {
		return fmt.Errorf("metrics registry is nil")
	}
	return r.prom.Register(c)
}

// Agent 返回 agent 指标集合，按需懒加载，只注册一次。
func (r *Registry) Agent() *AgentMetrics {
	if r == nil {
		return nil
	}
	r.once.Do(func() {
		r.agent = newAgentMetrics()
		r.MustRegister(r.agent.Collectors()...)
	})
	return r.agent
}
