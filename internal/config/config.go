package config

import "fmt"

const (
	ComponentAPIServer  = "api-server"
	ComponentScheduler  = "scheduler"
	ComponentController = "controller"
	ComponentWebhook    = "webhook"
	ComponentAgent      = "agent"
)

func LoadByComponent(component, path string) (any, error) {
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
