package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// AgentMetrics：gpu-agent 业务指标集合。
// 标签控制尽量简单，避免高基数。
type AgentMetrics struct {
	// discovery
	discoveryTotal       prometheus.Counter
	discoveryFailTotal   prometheus.Counter
	discoveryDuration    prometheus.Histogram
	inventoryGPUCount    prometheus.Gauge
	inventoryTotalMemMiB prometheus.Gauge

	// collect
	collectTotal      *prometheus.CounterVec   // labels: collector
	collectFailTotal  *prometheus.CounterVec   // labels: collector
	collectDuration   *prometheus.HistogramVec // labels: collector
	dcgmFallbackTotal prometheus.Counter

	collectGPUCount        prometheus.Gauge
	collectMIGCount        prometheus.Gauge
	collectTopologyCount   prometheus.Gauge
	collectPodBindingCount prometheus.Gauge

	// report
	reportTotal     *prometheus.CounterVec   // labels: mode
	reportFailTotal *prometheus.CounterVec   // labels: mode
	reportDuration  *prometheus.HistogramVec // labels: mode

	// heartbeat
	heartbeatTotal     prometheus.Counter
	heartbeatFailTotal prometheus.Counter
	heartbeatDuration  prometheus.Histogram
}

func newAgentMetrics() *AgentMetrics {
	return &AgentMetrics{
		// discovery
		discoveryTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "discovery_total",
			Help:      "Total number of inventory discovery attempts.",
		}),
		discoveryFailTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "discovery_fail_total",
			Help:      "Total number of failed inventory discovery attempts.",
		}),
		discoveryDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "discovery_duration_seconds",
			Help:      "Inventory discovery duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		}),
		inventoryGPUCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "inventory_gpu_count",
			Help:      "Discovered static GPU count in node inventory.",
		}),
		inventoryTotalMemMiB: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "inventory_total_memory_mib",
			Help:      "Discovered total GPU memory in MiB from node inventory.",
		}),

		// collect
		collectTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_total",
			Help:      "Total number of collect-and-report attempts by GPU collector type.",
		}, []string{"collector"}),
		collectFailTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_fail_total",
			Help:      "Total number of failed collect-and-report attempts by GPU collector type.",
		}, []string{"collector"}),
		collectDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_duration_seconds",
			Help:      "Collect-and-report duration in seconds by GPU collector type.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"collector"}),
		dcgmFallbackTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "dcgm_fallback_total",
			Help:      "Total number of times DCGM collection fell back to nvidia-smi.",
		}),
		collectGPUCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_gpu_count",
			Help:      "Collected GPU count in latest report cycle.",
		}),
		collectMIGCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_mig_count",
			Help:      "Collected MIG instance count in latest report cycle.",
		}),
		collectTopologyCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_topology_link_count",
			Help:      "Collected GPU topology link count in latest report cycle.",
		}),
		collectPodBindingCount: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "collect_pod_binding_count",
			Help:      "Collected pod-to-GPU binding count in latest report cycle.",
		}),

		// report
		reportTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "report_total",
			Help:      "Total number of agent report attempts by reporter mode.",
		}, []string{"mode"}),
		reportFailTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "report_fail_total",
			Help:      "Total number of failed agent report attempts by reporter mode.",
		}, []string{"mode"}),
		reportDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "report_duration_seconds",
			Help:      "Agent report duration in seconds by reporter mode.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"mode"}),

		// heartbeat
		heartbeatTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "heartbeat_total",
			Help:      "Total number of heartbeat send attempts.",
		}),
		heartbeatFailTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "heartbeat_fail_total",
			Help:      "Total number of failed heartbeat sends.",
		}),
		heartbeatDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: "gpu",
			Subsystem: "agent",
			Name:      "heartbeat_duration_seconds",
			Help:      "Heartbeat send duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		}),
	}
}

func (m *AgentMetrics) Collectors() []prometheus.Collector {
	if m == nil {
		return nil
	}
	return []prometheus.Collector{
		m.discoveryTotal,
		m.discoveryFailTotal,
		m.discoveryDuration,
		m.inventoryGPUCount,
		m.inventoryTotalMemMiB,

		m.collectTotal,
		m.collectFailTotal,
		m.collectDuration,
		m.dcgmFallbackTotal,
		m.collectGPUCount,
		m.collectMIGCount,
		m.collectTopologyCount,
		m.collectPodBindingCount,

		m.reportTotal,
		m.reportFailTotal,
		m.reportDuration,

		m.heartbeatTotal,
		m.heartbeatFailTotal,
		m.heartbeatDuration,
	}
}

// -------------------------
// discovery
// -------------------------

func (m *AgentMetrics) ObserveDiscovery(d time.Duration, err error) {
	if m == nil {
		return
	}
	m.discoveryTotal.Inc()
	m.discoveryDuration.Observe(d.Seconds())
	if err != nil {
		m.discoveryFailTotal.Inc()
	}
}

func (m *AgentMetrics) SetInventory(gpuCount int, totalMemoryMiB int64) {
	if m == nil {
		return
	}
	m.inventoryGPUCount.Set(float64(gpuCount))
	m.inventoryTotalMemMiB.Set(float64(totalMemoryMiB))
}

// -------------------------
// collect
// -------------------------

func (m *AgentMetrics) ObserveCollect(collector string, d time.Duration, err error) {
	if m == nil {
		return
	}
	if collector == "" {
		collector = "unknown"
	}
	m.collectTotal.WithLabelValues(collector).Inc()
	m.collectDuration.WithLabelValues(collector).Observe(d.Seconds())
	if err != nil {
		m.collectFailTotal.WithLabelValues(collector).Inc()
	}
}

func (m *AgentMetrics) IncDCGMFallback() {
	if m == nil {
		return
	}
	m.dcgmFallbackTotal.Inc()
}

func (m *AgentMetrics) SetCollectSnapshot(gpuCount, migCount, topologyLinkCount, podBindingCount int) {
	if m == nil {
		return
	}
	m.collectGPUCount.Set(float64(gpuCount))
	m.collectMIGCount.Set(float64(migCount))
	m.collectTopologyCount.Set(float64(topologyLinkCount))
	m.collectPodBindingCount.Set(float64(podBindingCount))
}

// -------------------------
// report
// -------------------------

func (m *AgentMetrics) ObserveReport(mode string, d time.Duration, err error) {
	if m == nil {
		return
	}
	if mode == "" {
		mode = "unknown"
	}
	m.reportTotal.WithLabelValues(mode).Inc()
	m.reportDuration.WithLabelValues(mode).Observe(d.Seconds())
	if err != nil {
		m.reportFailTotal.WithLabelValues(mode).Inc()
	}
}

// -------------------------
// heartbeat
// -------------------------

func (m *AgentMetrics) ObserveHeartbeat(d time.Duration, err error) {
	if m == nil {
		return
	}
	m.heartbeatTotal.Inc()
	m.heartbeatDuration.Observe(d.Seconds())
	if err != nil {
		m.heartbeatFailTotal.Inc()
	}
}
