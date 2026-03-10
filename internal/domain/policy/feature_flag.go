package policy

type FeatureFlags struct {
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
