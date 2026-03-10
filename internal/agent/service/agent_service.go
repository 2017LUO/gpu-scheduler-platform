package service

import (
	"context"
	"time"

	"gpu-scheduler-platform/internal/agent/collector"
	"gpu-scheduler-platform/internal/agent/discovery"
	agentreporter "gpu-scheduler-platform/internal/agent/reporter"
	appcfg "gpu-scheduler-platform/internal/config"

	"go.uber.org/zap"
)

type Service struct {
	cfg       *appcfg.AgentConfig
	logger    *zap.Logger
	discovery *discovery.DeviceDiscovery
	nodeMeta  *discovery.NodeMetaResolver
	nvidiaSMI *collector.NvidiaSMICollector
	dcgm      *collector.DCGMCollector
	mig       *collector.MIGCollector
	topology  *collector.TopologyCollector
	podGPU    *collector.PodGPUUsageCollector
	reporter  agentreporter.Reporter
	heartbeat *agentreporter.Heartbeat
}

type AgentReport struct {
	NodeName    string                 `json:"node_name"`
	Timestamp   time.Time              `json:"timestamp"`
	GPUs        []collector.GPUInfo    `json:"gpus"`
	MIGs        []collector.MIGInfo    `json:"migs,omitempty"`
	Topology    []collector.GPULink    `json:"topology,omitempty"`
	PodBindings []collector.PodGPUInfo `json:"pod_bindings,omitempty"`
}

func NewService(
	cfg *appcfg.AgentConfig,
	logger *zap.Logger,
	discovery *discovery.DeviceDiscovery,
	nodeMeta *discovery.NodeMetaResolver,
	nvidiaSMI *collector.NvidiaSMICollector,
	dcgm *collector.DCGMCollector,
	mig *collector.MIGCollector,
	topology *collector.TopologyCollector,
	podGPU *collector.PodGPUUsageCollector,
	reporter agentreporter.Reporter,
	heartbeat *agentreporter.Heartbeat,
) *Service {
	return &Service{
		cfg:       cfg,
		logger:    logger,
		discovery: discovery,
		nodeMeta:  nodeMeta,
		nvidiaSMI: nvidiaSMI,
		dcgm:      dcgm,
		mig:       mig,
		topology:  topology,
		podGPU:    podGPU,
		reporter:  reporter,
		heartbeat: heartbeat,
	}
}

func (s *Service) Run(ctx context.Context) error {
	collectTicker := time.NewTicker(s.cfg.Agent.CollectInterval)
	defer collectTicker.Stop()

	hbTicker := time.NewTicker(s.cfg.Agent.HeartbeatInterval)
	defer hbTicker.Stop()

	s.logger.Info("agent service started",
		zap.String("node_name", s.cfg.Agent.NodeName),
		zap.Duration("collect_interval", s.cfg.Agent.CollectInterval),
		zap.Duration("heartbeat_interval", s.cfg.Agent.HeartbeatInterval),
	)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("agent service stopped")
			return nil

		case <-collectTicker.C:
			if err := s.collectAndReport(ctx); err != nil {
				s.logger.Warn("collect and report failed", zap.Error(err))
			}

		case <-hbTicker.C:
			if err := s.heartbeat.Send(ctx, s.cfg.Agent.NodeName); err != nil {
				s.logger.Debug("heartbeat send failed", zap.Error(err))
			}
		}
	}
}

func (s *Service) collectAndReport(ctx context.Context) error {
	nodeName := s.nodeMeta.ResolveNodeName(s.cfg.Agent.NodeName)

	var gpus []collector.GPUInfo
	var err error

	switch {
	case s.cfg.Agent.EnableDCGM:
		gpus, err = s.dcgm.Collect(ctx)
	default:
		gpus, err = s.nvidiaSMI.Collect(ctx)
	}
	if err != nil {
		return err
	}

	report := AgentReport{
		NodeName:  nodeName,
		Timestamp: time.Now().UTC(),
		GPUs:      gpus,
	}

	if s.cfg.Agent.EnableMIGDiscovery {
		migs, err := s.mig.Collect(ctx)
		if err == nil {
			report.MIGs = migs
		}
	}

	if s.cfg.Agent.EnableTopologyDiscovery {
		links, err := s.topology.Collect(ctx)
		if err == nil {
			report.Topology = links
		}
	}

	podBindings, err := s.podGPU.Collect(ctx)
	if err == nil {
		report.PodBindings = podBindings
	}

	return s.reporter.Report(ctx, report)
}
