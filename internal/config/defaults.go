package config

import "time"

func defaultServiceConfig(name string) ServiceConfig {
	return ServiceConfig{
		Name:    name,
		Env:     "dev",
		Version: "v0.1.0",
	}
}

func defaultMySQLConfig() MySQLConfig {
	return MySQLConfig{
		DSN:             "root:123456@tcp(127.0.0.1:3306)/gpu_scheduler?charset=utf8mb4&parseTime=True&loc=Local",
		MaxOpenConns:    50,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

func defaultRedisConfig() RedisConfig {
	return RedisConfig{
		Addr:         "127.0.0.1:6379",
		Password:     "",
		DB:           0,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     20,
		MinIdleConns: 5,
	}
}

func defaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		Level:            "info",
		Format:           "console",
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func defaultKubernetesConfig() KubernetesConfig {
	return KubernetesConfig{
		InCluster:  false,
		Kubeconfig: "",
		QPS:        50,
		Burst:      100,
	}
}

func defaultLeaderElectionConfig(name string) LeaderElectionConfig {
	return LeaderElectionConfig{
		Enabled:        true,
		LeaseName:      name,
		LeaseNamespace: "gpu-system",
		LeaseDuration:  15 * time.Second,
		RenewDeadline:  10 * time.Second,
		RetryPeriod:    2 * time.Second,
	}
}

func defaultObservabilityConfig() ObservabilityConfig {
	return ObservabilityConfig{
		Metrics: MetricsConfig{
			Enabled: true,
			Path:    "/metrics",
		},
		Tracing: TracingConfig{
			Enabled:     false,
			Endpoint:    "",
			SampleRatio: 0.1,
		},
		PProf: PProfConfig{
			Enabled:    true,
			PathPrefix: "/debug/pprof",
		},
	}
}

func ApplyAPIServerDefaults(cfg *APIServerConfig) {
	if cfg == nil {
		return
	}

	if cfg.Service.Name == "" {
		cfg.Service = defaultServiceConfig("gpu-api-server")
	}
	if cfg.Service.Env == "" {
		cfg.Service.Env = "dev"
	}
	if cfg.Service.Version == "" {
		cfg.Service.Version = "v0.1.0"
	}

	if cfg.Server.HTTP.Addr == "" {
		cfg.Server.HTTP.Addr = ":8080"
	}
	if cfg.Server.HTTP.ReadTimeout == 0 {
		cfg.Server.HTTP.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.HTTP.WriteTimeout == 0 {
		cfg.Server.HTTP.WriteTimeout = 15 * time.Second
	}
	if cfg.Server.HTTP.IdleTimeout == 0 {
		cfg.Server.HTTP.IdleTimeout = 60 * time.Second
	}
	if cfg.Server.HTTP.ShutdownTimeout == 0 {
		cfg.Server.HTTP.ShutdownTimeout = 20 * time.Second
	}

	if cfg.Security.JWT.Issuer == "" {
		cfg.Security.JWT.Issuer = "gpu-scheduler-platform"
	}
	if cfg.Security.JWT.Expire == 0 {
		cfg.Security.JWT.Expire = 24 * time.Hour
	}

	applyMySQLDefaults(&cfg.MySQL)
	applyRedisDefaults(&cfg.Redis)
	applyObservabilityDefaults(&cfg.Observability)
	applyLoggingDefaults(&cfg.Logging)

	cfg.Features.EnableQueueAPI = defaultBool(cfg.Features.EnableQueueAPI, true)
	cfg.Features.EnablePolicyAPI = defaultBool(cfg.Features.EnablePolicyAPI, true)
	cfg.Features.EnableQuotaAPI = defaultBool(cfg.Features.EnableQuotaAPI, true)
	cfg.Features.EnableTenantAPI = defaultBool(cfg.Features.EnableTenantAPI, true)
	cfg.Features.EnableClusterAPI = defaultBool(cfg.Features.EnableClusterAPI, true)
}

func ApplySchedulerDefaults(cfg *SchedulerConfig) {
	if cfg == nil {
		return
	}

	if cfg.Service.Name == "" {
		cfg.Service = defaultServiceConfig("gpu-scheduler")
	}
	if cfg.Service.Env == "" {
		cfg.Service.Env = "dev"
	}
	if cfg.Service.Version == "" {
		cfg.Service.Version = "v0.1.0"
	}

	if cfg.Scheduler.ScheduleInterval == 0 {
		cfg.Scheduler.ScheduleInterval = 2 * time.Second
	}
	if cfg.Scheduler.MaxConcurrentCycles <= 0 {
		cfg.Scheduler.MaxConcurrentCycles = 4
	}
	if cfg.Scheduler.PendingBatchSize <= 0 {
		cfg.Scheduler.PendingBatchSize = 50
	}
	if cfg.Scheduler.ReservationTTL == 0 {
		cfg.Scheduler.ReservationTTL = 30 * time.Second
	}
	if cfg.Scheduler.BindTimeout == 0 {
		cfg.Scheduler.BindTimeout = 15 * time.Second
	}
	cfg.Scheduler.EnablePreemption = defaultBool(cfg.Scheduler.EnablePreemption, true)
	cfg.Scheduler.EnableGangScheduling = defaultBool(cfg.Scheduler.EnableGangScheduling, false)
	cfg.Scheduler.EnableTopologyAware = defaultBool(cfg.Scheduler.EnableTopologyAware, true)
	cfg.Scheduler.EnableMIG = defaultBool(cfg.Scheduler.EnableMIG, false)

	if cfg.Observability.Metrics.Addr == "" {
		cfg.Observability.Metrics.Addr = ":9091"
	}
	if cfg.Observability.PProf.Addr == "" {
		cfg.Observability.PProf.Addr = ":6061"
	}

	applyLeaderElectionDefaults(&cfg.LeaderElection, "gpu-scheduler-leader")
	applyMySQLDefaults(&cfg.MySQL)
	applyRedisDefaults(&cfg.Redis)
	applyKubernetesDefaults(&cfg.Kubernetes)
	applyObservabilityDefaults(&cfg.Observability)
	applyLoggingDefaults(&cfg.Logging)
}

func ApplyControllerDefaults(cfg *ControllerAppConfig) {
	if cfg == nil {
		return
	}

	if cfg.Service.Name == "" {
		cfg.Service = defaultServiceConfig("gpu-controller")
	}
	if cfg.Service.Env == "" {
		cfg.Service.Env = "dev"
	}
	if cfg.Service.Version == "" {
		cfg.Service.Version = "v0.1.0"
	}

	if cfg.Controller.Workers.GPUJob <= 0 {
		cfg.Controller.Workers.GPUJob = 4
	}
	if cfg.Controller.Workers.Node <= 0 {
		cfg.Controller.Workers.Node = 2
	}
	if cfg.Controller.Workers.Pod <= 0 {
		cfg.Controller.Workers.Pod = 4
	}
	if cfg.Controller.Workers.Quota <= 0 {
		cfg.Controller.Workers.Quota = 2
	}
	if cfg.Controller.Workers.Policy <= 0 {
		cfg.Controller.Workers.Policy = 2
	}
	if cfg.Controller.ResyncPeriod == 0 {
		cfg.Controller.ResyncPeriod = 30 * time.Second
	}
	if cfg.Controller.ReconcileTimeout == 0 {
		cfg.Controller.ReconcileTimeout = 20 * time.Second
	}
	cfg.Controller.EnableAllocationRecovery = defaultBool(cfg.Controller.EnableAllocationRecovery, true)
	cfg.Controller.EnableJobStatusSync = defaultBool(cfg.Controller.EnableJobStatusSync, true)

	if cfg.Observability.Metrics.Addr == "" {
		cfg.Observability.Metrics.Addr = ":9092"
	}
	if cfg.Observability.PProf.Addr == "" {
		cfg.Observability.PProf.Addr = ":6062"
	}

	applyLeaderElectionDefaults(&cfg.LeaderElection, "gpu-controller-leader")
	applyMySQLDefaults(&cfg.MySQL)
	applyRedisDefaults(&cfg.Redis)
	applyKubernetesDefaults(&cfg.Kubernetes)
	applyObservabilityDefaults(&cfg.Observability)
	applyLoggingDefaults(&cfg.Logging)
}

func ApplyWebhookDefaults(cfg *WebhookAppConfig) {
	if cfg == nil {
		return
	}

	if cfg.Service.Name == "" {
		cfg.Service = defaultServiceConfig("gpu-webhook")
	}
	if cfg.Service.Env == "" {
		cfg.Service.Env = "dev"
	}
	if cfg.Service.Version == "" {
		cfg.Service.Version = "v0.1.0"
	}

	if cfg.Server.HTTPS.Addr == "" {
		cfg.Server.HTTPS.Addr = ":9443"
	}
	if cfg.Server.HTTPS.ReadTimeout == 0 {
		cfg.Server.HTTPS.ReadTimeout = 10 * time.Second
	}
	if cfg.Server.HTTPS.WriteTimeout == 0 {
		cfg.Server.HTTPS.WriteTimeout = 10 * time.Second
	}
	if cfg.Server.HTTPS.ShutdownTimeout == 0 {
		cfg.Server.HTTPS.ShutdownTimeout = 20 * time.Second
	}

	cfg.Webhook.EnableMutating = defaultBool(cfg.Webhook.EnableMutating, true)
	cfg.Webhook.EnableValidating = defaultBool(cfg.Webhook.EnableValidating, true)
	if cfg.Webhook.FailurePolicy == "" {
		cfg.Webhook.FailurePolicy = "Fail"
	}
	if cfg.Webhook.SideEffects == "" {
		cfg.Webhook.SideEffects = "None"
	}

	if cfg.Observability.Metrics.Addr == "" {
		cfg.Observability.Metrics.Addr = ":9093"
	}
	if cfg.Observability.PProf.Addr == "" {
		cfg.Observability.PProf.Addr = ":6063"
	}
	cfg.Observability.PProf.Enabled = defaultBool(cfg.Observability.PProf.Enabled, false)

	applyKubernetesDefaults(&cfg.Kubernetes)
	applyObservabilityDefaults(&cfg.Observability)
	applyLoggingDefaults(&cfg.Logging)
}

func ApplyAgentDefaults(cfg *AgentConfig) {
	if cfg == nil {
		return
	}

	if cfg.Service.Name == "" {
		cfg.Service = defaultServiceConfig("gpu-agent")
	}
	if cfg.Service.Env == "" {
		cfg.Service.Env = "dev"
	}
	if cfg.Service.Version == "" {
		cfg.Service.Version = "v0.1.0"
	}

	if cfg.Agent.HeartbeatInterval == 0 {
		cfg.Agent.HeartbeatInterval = 10 * time.Second
	}
	if cfg.Agent.CollectInterval == 0 {
		cfg.Agent.CollectInterval = 10 * time.Second
	}
	if cfg.Agent.ReportTimeout == 0 {
		cfg.Agent.ReportTimeout = 5 * time.Second
	}
	cfg.Agent.EnableDCGM = defaultBool(cfg.Agent.EnableDCGM, false)
	cfg.Agent.EnableNvidiaSMI = defaultBool(cfg.Agent.EnableNvidiaSMI, true)
	cfg.Agent.EnableMIGDiscovery = defaultBool(cfg.Agent.EnableMIGDiscovery, true)
	cfg.Agent.EnableTopologyDiscovery = defaultBool(cfg.Agent.EnableTopologyDiscovery, true)

	if cfg.Reporter.Mode == "" {
		cfg.Reporter.Mode = "http"
	}
	if cfg.Reporter.HTTP.Timeout == 0 {
		cfg.Reporter.HTTP.Timeout = 5 * time.Second
	}
	if cfg.Reporter.GRPC.Timeout == 0 {
		cfg.Reporter.GRPC.Timeout = 5 * time.Second
	}

	if cfg.Observability.Metrics.Addr == "" {
		cfg.Observability.Metrics.Addr = ":9094"
	}
	if cfg.Observability.PProf.Addr == "" {
		cfg.Observability.PProf.Addr = ":6064"
	}
	cfg.Observability.PProf.Enabled = defaultBool(cfg.Observability.PProf.Enabled, false)

	applyKubernetesDefaults(&cfg.Kubernetes)
	applyObservabilityDefaults(&cfg.Observability)
	applyLoggingDefaults(&cfg.Logging)
}

func applyMySQLDefaults(cfg *MySQLConfig) {
	if cfg == nil {
		return
	}
	def := defaultMySQLConfig()
	if cfg.DSN == "" {
		cfg.DSN = def.DSN
	}
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = def.MaxOpenConns
	}
	if cfg.MaxIdleConns < 0 {
		cfg.MaxIdleConns = 0
	} else if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = def.MaxIdleConns
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = def.ConnMaxLifetime
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = def.ConnMaxIdleTime
	}
}

func applyRedisDefaults(cfg *RedisConfig) {
	if cfg == nil {
		return
	}
	def := defaultRedisConfig()
	if cfg.Addr == "" {
		cfg.Addr = def.Addr
	}
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = def.DialTimeout
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = def.ReadTimeout
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = def.WriteTimeout
	}
	if cfg.PoolSize <= 0 {
		cfg.PoolSize = def.PoolSize
	}
	if cfg.MinIdleConns < 0 {
		cfg.MinIdleConns = 0
	} else if cfg.MinIdleConns == 0 {
		cfg.MinIdleConns = def.MinIdleConns
	}
}

func applyKubernetesDefaults(cfg *KubernetesConfig) {
	if cfg == nil {
		return
	}
	def := defaultKubernetesConfig()
	if cfg.QPS <= 0 {
		cfg.QPS = def.QPS
	}
	if cfg.Burst <= 0 {
		cfg.Burst = def.Burst
	}
}

func applyLeaderElectionDefaults(cfg *LeaderElectionConfig, leaseName string) {
	if cfg == nil {
		return
	}
	def := defaultLeaderElectionConfig(leaseName)
	if cfg.LeaseName == "" {
		cfg.LeaseName = def.LeaseName
	}
	if cfg.LeaseNamespace == "" {
		cfg.LeaseNamespace = def.LeaseNamespace
	}
	if cfg.LeaseDuration == 0 {
		cfg.LeaseDuration = def.LeaseDuration
	}
	if cfg.RenewDeadline == 0 {
		cfg.RenewDeadline = def.RenewDeadline
	}
	if cfg.RetryPeriod == 0 {
		cfg.RetryPeriod = def.RetryPeriod
	}
}

func applyObservabilityDefaults(cfg *ObservabilityConfig) {
	if cfg == nil {
		return
	}
	def := defaultObservabilityConfig()

	if cfg.Metrics.Path == "" {
		cfg.Metrics.Path = def.Metrics.Path
	}
	cfg.Metrics.Enabled = defaultBool(cfg.Metrics.Enabled, def.Metrics.Enabled)

	if cfg.Tracing.SampleRatio == 0 {
		cfg.Tracing.SampleRatio = def.Tracing.SampleRatio
	}

	if cfg.PProf.PathPrefix == "" {
		cfg.PProf.PathPrefix = def.PProf.PathPrefix
	}
}

func applyLoggingDefaults(cfg *LoggingConfig) {
	if cfg == nil {
		return
	}
	def := defaultLoggingConfig()
	if cfg.Level == "" {
		cfg.Level = def.Level
	}
	if cfg.Format == "" {
		cfg.Format = def.Format
	}
	if len(cfg.OutputPaths) == 0 {
		cfg.OutputPaths = def.OutputPaths
	}
	if len(cfg.ErrorOutputPaths) == 0 {
		cfg.ErrorOutputPaths = def.ErrorOutputPaths
	}
}

func defaultBool(v bool, def bool) bool {
	if v {
		return true
	}
	return def && !v
}
