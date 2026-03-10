package reporter

import (
	"context"
	"time"
)

type Heartbeat struct {
	reporter Reporter
}

func NewHeartbeat(reporter Reporter) *Heartbeat {
	return &Heartbeat{reporter: reporter}
}

func (h *Heartbeat) Send(ctx context.Context, nodeName string) error {
	if h == nil || h.reporter == nil {
		return nil
	}
	payload := map[string]any{
		"type":      "heartbeat",
		"node_name": nodeName,
		"ts":        time.Now().UTC().Format(time.RFC3339Nano),
	}
	return h.reporter.Report(ctx, payload)
}
