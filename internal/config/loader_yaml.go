package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadAPIServerConfig(path string) (*APIServerConfig, error) {
	cfg := &APIServerConfig{}
	if err := loadYAML(path, cfg); err != nil {
		return nil, err
	}
	ApplyAPIServerDefaults(cfg)
	if err := ApplyAPIServerEnvOverrides(cfg); err != nil {
		return nil, err
	}
	if err := ValidateAPIServerConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadSchedulerConfig(path string) (*SchedulerConfig, error) {
	cfg := &SchedulerConfig{}
	if err := loadYAML(path, cfg); err != nil {
		return nil, err
	}
	ApplySchedulerDefaults(cfg)
	if err := ApplySchedulerEnvOverrides(cfg); err != nil {
		return nil, err
	}
	if err := ValidateSchedulerConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadControllerConfig(path string) (*ControllerAppConfig, error) {
	cfg := &ControllerAppConfig{}
	if err := loadYAML(path, cfg); err != nil {
		return nil, err
	}
	ApplyControllerDefaults(cfg)
	if err := ApplyControllerEnvOverrides(cfg); err != nil {
		return nil, err
	}
	if err := ValidateControllerConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadWebhookConfig(path string) (*WebhookAppConfig, error) {
	cfg := &WebhookAppConfig{}
	if err := loadYAML(path, cfg); err != nil {
		return nil, err
	}
	ApplyWebhookDefaults(cfg)
	if err := ApplyWebhookEnvOverrides(cfg); err != nil {
		return nil, err
	}
	if err := ValidateWebhookConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	cfg := &AgentConfig{}
	if err := loadYAML(path, cfg); err != nil {
		return nil, err
	}
	ApplyAgentDefaults(cfg)
	if err := ApplyAgentEnvOverrides(cfg); err != nil {
		return nil, err
	}
	if err := ValidateAgentConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func loadYAML(path string, out any) error {
	if out == nil {
		return fmt.Errorf("yaml decode target is nil")
	}
	if path == "" {
		return fmt.Errorf("yaml path is empty")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read yaml file %s: %w", path, err)
	}

	dec := yaml.NewDecoder(newBytesReader(data))
	dec.KnownFields(true)

	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode yaml file %s: %w", path, err)
	}
	return nil
}
