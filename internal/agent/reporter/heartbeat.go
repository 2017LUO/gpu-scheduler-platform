package reporter

import (
	"context"
	"time"

	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"
)

type Heartbeat struct {
	reporter Reporter
	metrics  *obsmetrics.AgentMetrics
}

func NewHeartbeat(reporter Reporter, metrics *obsmetrics.AgentMetrics) *Heartbeat {
	return &Heartbeat{
		reporter: reporter,
		metrics:  metrics,
	}
}

func (h *Heartbeat) Send(ctx context.Context, nodeName string) (retErr error) {
	start := time.Now()
	defer func() {
		if h.metrics != nil {
			h.metrics.ObserveHeartbeat(time.Since(start), retErr)
		}
	}()

	payload := map[string]any{
		"type":      "heartbeat",
		"node_name": nodeName,
		"timestamp": time.Now().UTC(),
	}

	retErr = h.reporter.Report(ctx, payload)
	return retErr
}
