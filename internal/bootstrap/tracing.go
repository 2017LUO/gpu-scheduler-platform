package bootstrap

import (
	"context"
	"fmt"

	appcfg "gpu-scheduler-platform/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type TracingCloser interface {
	Shutdown(ctx context.Context) error
}

type noopTracingCloser struct{}

func (n noopTracingCloser) Shutdown(context.Context) error { return nil }

func InitTracing(cfg appcfg.TracingConfig) (TracingCloser, error) {
	if !cfg.Enabled {
		otel.SetTracerProvider(trace.NewNoopTracerProvider())
		return noopTracingCloser{}, nil
	}

	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("tracing enabled but endpoint is empty")
	}

	otel.SetTracerProvider(trace.NewNoopTracerProvider())
	return noopTracingCloser{}, nil
}
