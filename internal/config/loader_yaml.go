package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadAPIServerConfig(path string) (*APIServerConfig, error) {
	var cfg APIServerConfig
	if err := loadYAML(path, &cfg); err != nil {
		return nil, err
	}
	ApplyAPIServerDefaults(&cfg)
	if err := ApplyAPIServerEnvOverrides(&cfg); err != nil {
		return nil, err
	}
	if err := ValidateAPIServerConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadSchedulerConfig(path string) (*SchedulerConfig, error) {
	var cfg SchedulerConfig
	if err := loadYAML(path, &cfg); err != nil {
		return nil, err
	}
	ApplySchedulerDefaults(&cfg)
	if err := ApplySchedulerEnvOverrides(&cfg); err != nil {
		return nil, err
	}
	if err := ValidateSchedulerConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadControllerConfig(path string) (*ControllerAppConfig, error) {
	var cfg ControllerAppConfig
	if err := loadYAML(path, &cfg); err != nil {
		return nil, err
	}
	ApplyControllerDefaults(&cfg)
	if err := ApplyControllerEnvOverrides(&cfg); err != nil {
		return nil, err
	}
	if err := ValidateControllerConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadWebhookConfig(path string) (*WebhookAppConfig, error) {
	var cfg WebhookAppConfig
	if err := loadYAML(path, &cfg); err != nil {
		return nil, err
	}
	ApplyWebhookDefaults(&cfg)
	if err := ApplyWebhookEnvOverrides(&cfg); err != nil {
		return nil, err
	}
	if err := ValidateWebhookConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadAgentConfig(path string) (*AgentConfig, error) {
	var cfg AgentConfig
	if err := loadYAML(path, &cfg); err != nil {
		return nil, err
	}
	ApplyAgentDefaults(&cfg)
	if err := ApplyAgentEnvOverrides(&cfg); err != nil {
		return nil, err
	}
	if err := ValidateAgentConfig(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func loadYAML(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file %q: %w", path, err)
	}

	var root yaml.Node
	if err := yaml.Unmarshal(b, &root); err != nil {
		return fmt.Errorf("unmarshal yaml %q: %w", path, err)
	}

	dec := yaml.NewDecoder(newBytesReader(b))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("decode yaml with known fields %q: %w", path, err)
	}
	return nil
}
