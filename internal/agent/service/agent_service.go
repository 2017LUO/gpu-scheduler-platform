package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/agent/collector"
	"gpu-scheduler-platform/internal/agent/discovery"
	agentreporter "gpu-scheduler-platform/internal/agent/reporter"
	appcfg "gpu-scheduler-platform/internal/config"
	obsmetrics "gpu-scheduler-platform/internal/observability/metrics"

	"go.uber.org/zap"
)

var ErrNoGPUCollectorEnabled = errors.New("no gpu collector enabled")

type Service struct {
	cfg     *appcfg.AgentConfig
	logger  *zap.Logger
	metrics *obsmetrics.AgentMetrics

	discovery *discovery.DeviceDiscovery
	nodeMeta  *discovery.NodeMetaResolver
	nvidiaSMI *collector.NvidiaSMICollector
	dcgm      *collector.DCGMCollector
	mig       *collector.MIGCollector
	topology  *collector.TopologyCollector
	podGPU    *collector.PodGPUUsageCollector
	reporter  agentreporter.Reporter
	heartbeat *agentreporter.Heartbeat

	inventory *discovery.DeviceInventory
}

type AgentReport struct {
	NodeName    string                     `json:"node_name"`
	Timestamp   time.Time                  `json:"timestamp"`
	Inventory   *discovery.DeviceInventory `json:"inventory,omitempty"`
	GPUs        []collector.GPUInfo        `json:"gpus"`
	MIGs        []collector.MIGInfo        `json:"migs,omitempty"`
	Topology    []collector.GPULink        `json:"topology,omitempty"`
	PodBindings []collector.PodGPUInfo     `json:"pod_bindings,omitempty"`
}

func (r AgentReport) GetNodeName() string {
	return r.NodeName
}

func NewService(
	cfg *appcfg.AgentConfig,
	logger *zap.Logger,
	metrics *obsmetrics.AgentMetrics,
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
	if logger == nil {
		logger = zap.NewNop()
	}

	return &Service{
		cfg:       cfg,
		logger:    logger,
		metrics:   metrics,
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
	if s == nil {
		return nil
	}
	if s.cfg == nil {
		return fmt.Errorf("agent service config is nil")
	}
	if s.reporter == nil {
		return fmt.Errorf("agent reporter is nil")
	}

	nodeName := s.resolveNodeName()

	// 启动时先做一次静态发现
	if s.discovery != nil {
		start := time.Now()
		inv, err := s.discovery.Discover(ctx, nodeName)
		if s.metrics != nil {
			s.metrics.ObserveDiscovery(time.Since(start), err)
		}
		if err != nil {
			s.logger.Warn("device discovery failed", zap.Error(err))
		} else if inv != nil {
			s.inventory = inv
			if s.metrics != nil {
				s.metrics.SetInventory(inv.GPUCount, inv.TotalMemoryMiB)
			}

			s.logger.Info("device inventory discovered",
				zap.String("node_name", inv.NodeName),
				zap.Int("gpu_count", inv.GPUCount),
				zap.Int64("total_memory_mib", inv.TotalMemoryMiB),
			)
		}
	}

	collectInterval := s.cfg.Agent.CollectInterval
	if collectInterval <= 0 {
		collectInterval = 15 * time.Second
	}
	heartbeatInterval := s.cfg.Agent.HeartbeatInterval
	if heartbeatInterval <= 0 {
		heartbeatInterval = 10 * time.Second
	}

	collectTicker := time.NewTicker(collectInterval)
	defer collectTicker.Stop()

	hbTicker := time.NewTicker(heartbeatInterval)
	defer hbTicker.Stop()

	s.logger.Info("agent service started",
		zap.String("node_name", nodeName),
		zap.Duration("collect_interval", collectInterval),
		zap.Duration("heartbeat_interval", heartbeatInterval),
	)

	// 启动后立即采一次
	if err := s.collectAndReport(ctx); err != nil {
		s.logger.Warn("initial collect and report failed", zap.Error(err))
	}

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
			if s.heartbeat == nil {
				continue
			}
			if err := s.heartbeat.Send(ctx, nodeName); err != nil {
				s.logger.Debug("heartbeat send failed", zap.Error(err))
			}
		}
	}
}

func (s *Service) collectAndReport(ctx context.Context) (retErr error) {
	start := time.Now()

	if s == nil {
		return nil
	}
	if s.cfg == nil {
		return fmt.Errorf("agent service config is nil")
	}
	if s.reporter == nil {
		return fmt.Errorf("agent reporter is nil")
	}

	nodeName := s.resolveNodeName()

	var (
		gpus          []collector.GPUInfo
		migs          []collector.MIGInfo
		links         []collector.GPULink
		podBindings   []collector.PodGPUInfo
		err           error
		collectorName string
	)

	defer func() {
		if s.metrics != nil {
			s.metrics.ObserveCollect(collectorName, time.Since(start), retErr)
		}
	}()

	// 优先 DCGM，失败则自动回退 nvidia-smi
	switch {
	case s.cfg.Agent.EnableDCGM && s.dcgm != nil:
		collectorName = "dcgm"
		gpus, err = s.dcgm.Collect(ctx)
		if err != nil {
			s.logger.Warn("dcgm collect failed, fallback to nvidia-smi", zap.Error(err))
			if s.metrics != nil {
				s.metrics.IncDCGMFallback()
			}

			if s.cfg.Agent.EnableNvidiaSMI && s.nvidiaSMI != nil {
				collectorName = "nvidia_smi"
				gpus, err = s.nvidiaSMI.Collect(ctx)
			}
		}

	case s.cfg.Agent.EnableNvidiaSMI && s.nvidiaSMI != nil:
		collectorName = "nvidia_smi"
		gpus, err = s.nvidiaSMI.Collect(ctx)

	default:
		collectorName = "unknown"
		err = ErrNoGPUCollectorEnabled
	}

	if err != nil {
		retErr = err
		return retErr
	}

	report := AgentReport{
		NodeName:  nodeName,
		Timestamp: time.Now().UTC(),
		Inventory: s.inventory,
		GPUs:      gpus,
	}

	if s.cfg.Agent.EnableMIGDiscovery && s.mig != nil {
		if migs, err = s.mig.Collect(ctx); err == nil {
			report.MIGs = migs
		} else {
			s.logger.Debug("mig collect failed", zap.Error(err))
		}
	}

	if s.cfg.Agent.EnableTopologyDiscovery && s.topology != nil {
		if links, err = s.topology.Collect(ctx); err == nil {
			report.Topology = links
		} else {
			s.logger.Debug("topology collect failed", zap.Error(err))
		}
	}

	if s.podGPU != nil {
		if podBindings, err = s.podGPU.Collect(ctx); err == nil {
			report.PodBindings = podBindings
		} else {
			s.logger.Debug("pod gpu collect failed", zap.Error(err))
		}
	}

	if s.metrics != nil {
		s.metrics.SetCollectSnapshot(
			len(gpus),
			len(migs),
			len(links),
			len(podBindings),
		)
	}

	reportCtx := ctx
	cancel := func() {}
	if s.cfg.Agent.ReportTimeout > 0 {
		reportCtx, cancel = context.WithTimeout(ctx, s.cfg.Agent.ReportTimeout)
	}
	defer cancel()

	retErr = s.reporter.Report(reportCtx, report)
	return retErr
}

func (s *Service) resolveNodeName() string {
	if s == nil || s.cfg == nil {
		return ""
	}
	if s.nodeMeta == nil {
		return s.cfg.Agent.NodeName
	}
	return s.nodeMeta.ResolveNodeName(s.cfg.Agent.NodeName)
}
