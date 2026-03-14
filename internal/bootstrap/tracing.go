package bootstrap

import (
	"context"
	"fmt"
	"strings"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"
	obstracing "gpu-scheduler-platform/internal/observability/tracing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	noop "go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type TracingCloser interface {
	Shutdown(ctx context.Context) error
}

type noopTracingCloser struct{}

func (noopTracingCloser) Shutdown(context.Context) error { return nil }

type tracingCloser struct {
	tp *sdktrace.TracerProvider
}

func (c tracingCloser) Shutdown(ctx context.Context) error {
	if c.tp == nil {
		return nil
	}
	return c.tp.Shutdown(ctx)
}

// InitTracing 使用顶层 service 配置和 observability.tracing 配置初始化全局 tracing。
func InitTracing(serviceCfg appcfg.ServiceConfig, tracingCfg appcfg.TracingConfig) (TracingCloser, error) {
	if !tracingCfg.Enabled {
		otel.SetTracerProvider(noop.NewTracerProvider())
		obstracing.SetGlobalPropagators()
		return noopTracingCloser{}, nil
	}

	if strings.TrimSpace(tracingCfg.Endpoint) == "" {
		return nil, fmt.Errorf("tracing enabled but endpoint is empty")
	}

	obsCfg := obstracing.Config{
		Enabled:     tracingCfg.Enabled,
		Endpoint:    tracingCfg.Endpoint,
		ServiceName: serviceCfg.Name,
		ServiceEnv:  serviceCfg.Env,
		Version:     serviceCfg.Version,
		SampleRatio: tracingCfg.SampleRatio,
		Insecure:    tracingCfg.Insecure,
		Headers:     tracingCfg.Headers,
	}

	res, err := obstracing.NewResource(obsCfg)
	if err != nil {
		return nil, fmt.Errorf("build tracing resource: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	expOpts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(tracingCfg.Endpoint),
	}

	switch {
	case tracingCfg.Insecure:
		expOpts = append(expOpts, otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
	case tracingCfg.TLS.Enabled:
		tlsCfg, err := BuildClientTLSConfig(tracingCfg.TLS)
		if err != nil {
			return nil, fmt.Errorf("build tracing tls config: %w", err)
		}
		expOpts = append(expOpts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(tlsCfg)))
	default:
		expOpts = append(expOpts, otlptracegrpc.WithTLSCredentials(credentials.NewTLS(nil)))
	}

	if len(tracingCfg.Headers) > 0 {
		expOpts = append(expOpts, otlptracegrpc.WithHeaders(tracingCfg.Headers))
	}

	exporter, err := otlptracegrpc.New(ctx, expOpts...)
	if err != nil {
		return nil, fmt.Errorf("create otlp trace exporter: %w", err)
	}

	tp := obstracing.NewTracerProvider(res, exporter, tracingCfg.SampleRatio)
	otel.SetTracerProvider(tp)
	obstracing.SetGlobalPropagators()

	return tracingCloser{tp: tp}, nil
}
