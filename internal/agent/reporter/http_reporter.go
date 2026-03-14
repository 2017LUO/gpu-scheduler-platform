package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"
)

type HTTPReporter struct {
	endpoint string
	timeout  time.Duration
	client   *http.Client
	metrics  *obsmetrics.AgentMetrics
}

func NewHTTPReporter(endpoint string, timeout time.Duration, metrics *obsmetrics.AgentMetrics) *HTTPReporter {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}

	return &HTTPReporter{
		endpoint: endpoint,
		timeout:  timeout,
		client: &http.Client{
			Timeout: timeout,
		},
		metrics: metrics,
	}
}

func (r *HTTPReporter) Report(ctx context.Context, payload any) (retErr error) {
	start := time.Now()
	defer func() {
		if r.metrics != nil {
			r.metrics.ObserveReport("http", time.Since(start), retErr)
		}
	}()

	if strings.TrimSpace(r.endpoint) == "" {
		retErr = fmt.Errorf("http reporter endpoint is empty")
		return retErr
	}

	body, err := json.Marshal(payload)
	if err != nil {
		retErr = fmt.Errorf("marshal http payload: %w", err)
		return retErr
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(body))
	if err != nil {
		retErr = fmt.Errorf("build http request: %w", err)
		return retErr
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		retErr = fmt.Errorf("do http request: %w", err)
		return retErr
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		retErr = fmt.Errorf("http reporter unexpected status: %s", resp.Status)
		return retErr
	}

	return nil
}
