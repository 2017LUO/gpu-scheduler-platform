package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	nodeagentv1 "gpu-scheduler-platform/api/proto/nodeagent/v1"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCReporter struct {
	endpoint string
	timeout  time.Duration
	metrics  *obsmetrics.AgentMetrics

	mu     sync.Mutex
	conn   *grpc.ClientConn
	client nodeagentv1.AgentReporterServiceClient
}

func NewGRPCReporter(endpoint string, timeout time.Duration, metrics *obsmetrics.AgentMetrics) *GRPCReporter {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &GRPCReporter{
		endpoint: endpoint,
		timeout:  timeout,
		metrics:  metrics,
	}
}

func (r *GRPCReporter) Report(ctx context.Context, payload any) (retErr error) {
	start := time.Now()
	defer func() {
		if r.metrics != nil {
			r.metrics.ObserveReport("grpc", time.Since(start), retErr)
		}
	}()

	if r == nil {
		return fmt.Errorf("grpc reporter is nil")
	}
	if strings.TrimSpace(r.endpoint) == "" {
		retErr = fmt.Errorf("grpc reporter endpoint is empty")
		return retErr
	}
	if ctx == nil {
		ctx = context.Background()
	}

	client, err := r.getClient()
	if err != nil {
		retErr = err
		return retErr
	}

	body, err := json.Marshal(payload)
	if err != nil {
		retErr = fmt.Errorf("marshal grpc payload: %w", err)
		return retErr
	}

	callCtx := ctx
	cancel := func() {}
	if _, hasDeadline := ctx.Deadline(); !hasDeadline && r.timeout > 0 {
		callCtx, cancel = context.WithTimeout(ctx, r.timeout)
	}
	defer cancel()

	req := &nodeagentv1.ReportRequest{
		PayloadType: detectPayloadType(payload),
		BodyJson:    string(body),
		ContentType: "application/json",
		SentAtUnix:  time.Now().UTC().Unix(),
		NodeName:    detectNodeName(payload),
	}

	_, err = client.Report(callCtx, req)
	if err != nil {
		retErr = fmt.Errorf("grpc report call failed: %w", err)
		return retErr
	}

	return nil
}

func (r *GRPCReporter) getClient() (nodeagentv1.AgentReporterServiceClient, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		return r.client, nil
	}

	conn, err := grpc.NewClient(
		r.endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("create grpc client for %q: %w", r.endpoint, err)
	}

	r.conn = conn
	r.client = nodeagentv1.NewAgentReporterServiceClient(conn)
	return r.client, nil
}

func (r *GRPCReporter) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.conn == nil {
		return nil
	}

	err := r.conn.Close()
	r.conn = nil
	r.client = nil
	return err
}

func detectPayloadType(payload any) string {
	type payloadTyper interface {
		GetPayloadType() string
	}

	if t, ok := payload.(payloadTyper); ok {
		if typ := strings.TrimSpace(t.GetPayloadType()); typ != "" {
			return typ
		}
	}

	switch v := payload.(type) {
	case map[string]any:
		if typ, ok := v["type"].(string); ok && strings.TrimSpace(typ) != "" {
			return typ
		}
		if typ, ok := v["payload_type"].(string); ok && strings.TrimSpace(typ) != "" {
			return typ
		}
	}

	return "agent_report"
}

func detectNodeName(payload any) string {
	type nodeNamer interface {
		GetNodeName() string
	}

	if n, ok := payload.(nodeNamer); ok {
		return strings.TrimSpace(n.GetNodeName())
	}

	switch v := payload.(type) {
	case map[string]any:
		if s, ok := v["node_name"].(string); ok {
			return strings.TrimSpace(s)
		}
		if s, ok := v["nodeName"].(string); ok {
			return strings.TrimSpace(s)
		}
	}

	return ""
}
