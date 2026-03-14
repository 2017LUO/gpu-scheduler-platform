package tracing

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

// Config 是 tracing 初始化所需的最小配置。
// endpoint 例子：
//   - "127.0.0.1:4317"
//   - "otel-collector.observability.svc.cluster.local:4317"
type Config struct {
	Enabled     bool
	Endpoint    string
	ServiceName string
	ServiceEnv  string
	Version     string
	SampleRatio float64
	Insecure    bool
	Headers     map[string]string
}

// Shutdowner 抽象 tracer provider 的关闭能力。
type Shutdowner interface {
	Shutdown(ctx context.Context) error
}

type noopShutdowner struct{}

func (noopShutdowner) Shutdown(context.Context) error { return nil }

// NewResource 构造统一的 OTel Resource。
// Resource 应在 provider 创建时绑定，后续不能再修改。
func NewResource(cfg Config) (*sdkresource.Resource, error) {
	serviceName := strings.TrimSpace(cfg.ServiceName)
	if serviceName == "" {
		serviceName = "gpu-scheduler-platform"
	}

	attrs := []attribute.KeyValue{
		semconv.ServiceName(serviceName),
	}

	if v := strings.TrimSpace(cfg.Version); v != "" {
		attrs = append(attrs, semconv.ServiceVersion(v))
	}
	if env := strings.TrimSpace(cfg.ServiceEnv); env != "" {
		attrs = append(attrs, attribute.String("deployment.environment.name", env))
	}

	res, err := sdkresource.Merge(
		sdkresource.Default(),
		sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			attrs...,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("merge otel resource: %w", err)
	}
	return res, nil
}

// NewTracer 返回命名 tracer，供业务代码使用。
// 例：tracing.NewTracer("gpu-agent/service")
func NewTracer(name string) trace.Tracer {
	if strings.TrimSpace(name) == "" {
		name = "gpu-scheduler-platform"
	}
	return otel.Tracer(name)
}

// SetGlobalPropagators 设置全局传播器。
// W3C TraceContext + Baggage 是官方常见默认组合。
func SetGlobalPropagators() {
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}

// StartSpan 是一个轻量 helper，便于统一写法。
func StartSpan(
	ctx context.Context,
	tracerName string,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	return NewTracer(tracerName).Start(ctx, spanName, opts...)
}

// RecordError 统一记录错误到 span。
func RecordError(span trace.Span, err error) {
	if span == nil || err == nil {
		return
	}
	span.RecordError(err)
}

// NewTracerProvider 构造 TracerProvider。
// exporter 由 bootstrap 层创建后传入；这里专注 provider 组装。
func NewTracerProvider(
	res *sdkresource.Resource,
	exporter sdktrace.SpanExporter,
	sampleRatio float64,
) *sdktrace.TracerProvider {
	if sampleRatio <= 0 {
		sampleRatio = 0.1
	}
	if sampleRatio > 1 {
		sampleRatio = 1
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRatio))),
		sdktrace.WithBatcher(exporter),
	)
}
