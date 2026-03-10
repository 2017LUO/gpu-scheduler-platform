package config

import "fmt"

type Component string

const (
	ComponentAPIServer  Component = "api-server"
	ComponentScheduler  Component = "scheduler"
	ComponentController Component = "controller"
	ComponentWebhook    Component = "webhook"
	ComponentAgent      Component = "agent"
)

func LoadByComponent(component Component, path string) (any, error) {
	switch component {
	case ComponentAPIServer:
		return LoadAPIServerConfig(path)
	case ComponentScheduler:
		return LoadSchedulerConfig(path)
	case ComponentController:
		return LoadControllerConfig(path)
	case ComponentWebhook:
		return LoadWebhookConfig(path)
	case ComponentAgent:
		return LoadAgentConfig(path)
	default:
		return nil, fmt.Errorf("unsupported component %q", component)
	}
}
