package bootstrap

import (
	"fmt"

	"gpu-scheduler-platform/internal/config"
)

func LoadAPIServerConfig(path string) (*config.APIServerConfig, error) {
	cfg, err := config.LoadAPIServerConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load api-server config: %w", err)
	}
	return cfg, nil
}

func LoadSchedulerConfig(path string) (*config.SchedulerConfig, error) {
	cfg, err := config.LoadSchedulerConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load scheduler config: %w", err)
	}
	return cfg, nil
}

func LoadControllerConfig(path string) (*config.ControllerAppConfig, error) {
	cfg, err := config.LoadControllerConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load controller config: %w", err)
	}
	return cfg, nil
}

func LoadWebhookConfig(path string) (*config.WebhookAppConfig, error) {
	cfg, err := config.LoadWebhookConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load webhook config: %w", err)
	}
	return cfg, nil
}

func LoadAgentConfig(path string) (*config.AgentConfig, error) {
	cfg, err := config.LoadAgentConfig(path)
	if err != nil {
		return nil, fmt.Errorf("load agent config: %w", err)
	}
	return cfg, nil
}
