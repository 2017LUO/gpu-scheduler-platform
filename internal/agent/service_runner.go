package agent

import (
	"context"

	"gpu-scheduler-platform/internal/agent/collector"
	"gpu-scheduler-platform/internal/agent/discovery"
	agentreporter "gpu-scheduler-platform/internal/agent/reporter"
	agentservice "gpu-scheduler-platform/internal/agent/service"
	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
)

type ServiceRunner struct {
	service *agentservice.Service
}

func NewServiceRunner(cfg *appcfg.AgentConfig, lg *zap.Logger, _ any) *ServiceRunner {
	var reporter agentreporter.Reporter
	switch cfg.Reporter.Mode {
	case "grpc":
		reporter = agentreporter.NewGRPCReporter(cfg.Reporter.GRPC.Endpoint, cfg.Reporter.GRPC.Timeout)
	default:
		reporter = agentreporter.NewHTTPReporter(cfg.Reporter.HTTP.Endpoint, cfg.Reporter.HTTP.Timeout)
	}

	svc := agentservice.NewService(
		cfg,
		lg,
		discovery.NewDeviceDiscovery(),
		discovery.NewNodeMetaResolver(),
		collector.NewNvidiaSMICollector(),
		collector.NewDCGMCollector(),
		collector.NewMIGCollector(),
		collector.NewTopologyCollector(),
		collector.NewPodGPUUsageCollector(),
		reporter,
		agentreporter.NewHeartbeat(reporter),
	)

	return &ServiceRunner{
		service: svc,
	}
}

func (r *ServiceRunner) Run(ctx context.Context) error {
	if r == nil || r.service == nil {
		return nil
	}
	return r.service.Run(ctx)
}
