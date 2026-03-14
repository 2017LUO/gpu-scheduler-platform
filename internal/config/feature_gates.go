package config

type FeatureGates map[string]bool

func (f FeatureGates) Enabled(name string) bool {
	if f == nil {
		return false
	}
	return f[name]
}

func APIServerFeatureGates(cfg APIFeaturesConfig) FeatureGates {
	return FeatureGates{
		"queue_api":   cfg.EnableQueueAPI,
		"policy_api":  cfg.EnablePolicyAPI,
		"quota_api":   cfg.EnableQuotaAPI,
		"tenant_api":  cfg.EnableTenantAPI,
		"cluster_api": cfg.EnableClusterAPI,
	}
}

func SchedulerFeatureGates(cfg SchedulerCoreConfig) FeatureGates {
	return FeatureGates{
		"preemption":      cfg.EnablePreemption,
		"gang_scheduling": cfg.EnableGangScheduling,
		"topology_aware":  cfg.EnableTopologyAware,
		"mig":             cfg.EnableMIG,
	}
}

func ControllerFeatureGates(cfg ControllerConfig) FeatureGates {
	return FeatureGates{
		"allocation_recovery": cfg.EnableAllocationRecovery,
		"job_status_sync":     cfg.EnableJobStatusSync,
	}
}

func WebhookFeatureGates(cfg WebhookConfig) FeatureGates {
	return FeatureGates{
		"mutating":   cfg.EnableMutating,
		"validating": cfg.EnableValidating,
	}
}

func AgentFeatureGates(cfg AgentCoreConfig) FeatureGates {
	return FeatureGates{
		"dcgm":               cfg.EnableDCGM,
		"nvidia_smi":         cfg.EnableNvidiaSMI,
		"mig_discovery":      cfg.EnableMIGDiscovery,
		"topology_discovery": cfg.EnableTopologyDiscovery,
	}
}
