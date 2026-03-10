package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Reporter interface {
	Report(ctx context.Context, payload any) error
}

type HTTPReporter struct {
	endpoint string
	timeout  time.Duration
	client   *http.Client
}

func NewHTTPReporter(endpoint string, timeout time.Duration) *HTTPReporter {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HTTPReporter{
		endpoint: endpoint,
		timeout:  timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (r *HTTPReporter) Report(ctx context.Context, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal report payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build report request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("send report request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("report endpoint returned status %d", resp.StatusCode)
	}
	return nil
}
