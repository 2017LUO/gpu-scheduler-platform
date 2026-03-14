package agent

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/agent/collector"
	"gpu-scheduler-platform/internal/agent/discovery"
	agentreporter "gpu-scheduler-platform/internal/agent/reporter"
	agentservice "gpu-scheduler-platform/internal/agent/service"
	appcfg "gpu-scheduler-platform/internal/config"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"

	"go.uber.org/zap"
)

type ServiceRunner struct {
	service  *agentservice.Service
	reporter agentreporter.Reporter
}

func NewServiceRunner(
	cfg *appcfg.AgentConfig,
	lg *zap.Logger,
	reg *obsmetrics.Registry,
	_ any,
) (*ServiceRunner, error) {
	if lg == nil {
		lg = zap.NewNop()
	}
	if cfg == nil {
		return nil, fmt.Errorf("agent config is nil")
	}
	if err := ValidateAgentConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid agent config: %w", err)
	}

	var agentM *obsmetrics.AgentMetrics
	if reg != nil {
		agentM = reg.Agent()
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.Reporter.Mode))

	var reporter agentreporter.Reporter
	switch mode {
	case "grpc":
		reporter = agentreporter.NewGRPCReporter(
			cfg.Reporter.GRPC.Endpoint,
			cfg.Reporter.GRPC.Timeout,
			agentM,
		)
	case "http":
		reporter = agentreporter.NewHTTPReporter(
			cfg.Reporter.HTTP.Endpoint,
			cfg.Reporter.HTTP.Timeout,
			agentM,
		)
	default:
		return nil, fmt.Errorf("unsupported reporter mode: %s", cfg.Reporter.Mode)
	}

	svc := agentservice.NewService(
		cfg,
		lg,
		agentM,
		discovery.NewDeviceDiscovery(),
		discovery.NewNodeMetaResolver(),
		collector.NewNvidiaSMICollector(),
		collector.NewDCGMCollector(),
		collector.NewMIGCollector(),
		collector.NewTopologyCollector(),
		collector.NewPodGPUUsageCollector(),
		reporter,
		agentreporter.NewHeartbeat(reporter, agentM),
	)

	if svc == nil {
		return nil, fmt.Errorf("agent service is nil")
	}

	return &ServiceRunner{
		service:  svc,
		reporter: reporter,
	}, nil
}

func (r *ServiceRunner) Run(ctx context.Context) error {
	if r == nil || r.service == nil {
		return nil
	}
	return r.service.Run(ctx)
}

func (r *ServiceRunner) Close() error {
	if r == nil || r.reporter == nil {
		return nil
	}

	type closer interface {
		Close() error
	}

	if c, ok := r.reporter.(closer); ok {
		return c.Close()
	}

	return nil
}

func (r *ServiceRunner) Reporter() agentreporter.Reporter {
	if r == nil {
		return nil
	}
	return r.reporter
}

func (r *ServiceRunner) Service() *agentservice.Service {
	if r == nil {
		return nil
	}
	return r.service
}
