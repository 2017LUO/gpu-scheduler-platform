package agent

import (
	"fmt"
	"strings"
	"time"

	appcfg "gpu-scheduler-platform/internal/config"
)

func ValidateAgentConfig(cfg *appcfg.AgentConfig) error {
	if cfg == nil {
		return fmt.Errorf("agent config is nil")
	}

	if cfg.Agent.HeartbeatInterval <= 0 {
		return fmt.Errorf("agent.heartbeat_interval must be > 0")
	}
	if cfg.Agent.CollectInterval <= 0 {
		return fmt.Errorf("agent.collect_interval must be > 0")
	}
	if cfg.Agent.ReportTimeout <= 0 {
		return fmt.Errorf("agent.report_timeout must be > 0")
	}

	if !cfg.Agent.EnableDCGM && !cfg.Agent.EnableNvidiaSMI {
		return fmt.Errorf("at least one gpu collector must be enabled: enable_dcgm or enable_nvidia_smi")
	}

	mode := strings.ToLower(strings.TrimSpace(cfg.Reporter.Mode))
	switch mode {
	case "grpc":
		if strings.TrimSpace(cfg.Reporter.GRPC.Endpoint) == "" {
			return fmt.Errorf("reporter.grpc.endpoint is required when reporter.mode=grpc")
		}
		if cfg.Reporter.GRPC.Timeout <= 0 {
			cfg.Reporter.GRPC.Timeout = 5 * time.Second
		}
	case "http":
		if strings.TrimSpace(cfg.Reporter.HTTP.Endpoint) == "" {
			return fmt.Errorf("reporter.http.endpoint is required when reporter.mode=http")
		}
		if cfg.Reporter.HTTP.Timeout <= 0 {
			cfg.Reporter.HTTP.Timeout = 5 * time.Second
		}
	default:
		return fmt.Errorf("unsupported reporter.mode=%q, allowed values: grpc, http", cfg.Reporter.Mode)
	}

	return nil
}
