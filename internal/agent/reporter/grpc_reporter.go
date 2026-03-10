package reporter

import (
	"context"
	"fmt"
	"time"
)

type GRPCReporter struct {
	endpoint string
	timeout  time.Duration
}

func NewGRPCReporter(endpoint string, timeout time.Duration) *GRPCReporter {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &GRPCReporter{
		endpoint: endpoint,
		timeout:  timeout,
	}
}

func (r *GRPCReporter) Report(ctx context.Context, payload any) error {
	_ = ctx
	_ = payload
	return fmt.Errorf("grpc reporter not implemented for endpoint %s", r.endpoint)
}
