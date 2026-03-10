package config

type FeatureGates struct {
	EnablePreemption         bool
	EnableMIG                bool
	EnableTopologyAware      bool
	EnableFragmentationScore bool
	EnableGangScheduling     bool
	EnableQuotaEnforcement   bool
	EnableAuditLog           bool
	EnableControllerRecovery bool
	EnableAgentReporting     bool
}

func DefaultFeatureGates() FeatureGates {
	return FeatureGates{
		EnablePreemption:         true,
		EnableMIG:                false,
		EnableTopologyAware:      true,
		EnableFragmentationScore: true,
		EnableGangScheduling:     false,
		EnableQuotaEnforcement:   true,
		EnableAuditLog:           true,
		EnableControllerRecovery: true,
		EnableAgentReporting:     true,
	}
}

func FeatureGatesFromScheduler(cfg *SchedulerConfig) FeatureGates {
	g := DefaultFeatureGates()
	if cfg == nil {
		return g
	}

	g.EnablePreemption = cfg.Scheduler.EnablePreemption
	g.EnableMIG = cfg.Scheduler.EnableMIG
	g.EnableTopologyAware = cfg.Scheduler.EnableTopologyAware
	g.EnableGangScheduling = cfg.Scheduler.EnableGangScheduling
	return g
}
