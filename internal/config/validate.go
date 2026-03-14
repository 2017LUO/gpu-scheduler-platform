package config

import (
	"errors"
	"fmt"
	"strings"
)

func ValidateAPIServerConfig(cfg *APIServerConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if err := validateService(cfg.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if err := validateHTTPServer(cfg.Server.HTTP, true); err != nil {
		return fmt.Errorf("server.http: %w", err)
	}
	if err := validateSecurity(cfg.Security); err != nil {
		return fmt.Errorf("security: %w", err)
	}
	if err := validateMySQL(cfg.MySQL); err != nil {
		return fmt.Errorf("mysql: %w", err)
	}
	if err := validateRedis(cfg.Redis); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if err := validateObservability(cfg.Observability); err != nil {
		return fmt.Errorf("observability: %w", err)
	}
	if err := validateLogging(cfg.Logging); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func ValidateSchedulerConfig(cfg *SchedulerConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if err := validateService(cfg.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if err := validateSchedulerCore(cfg.Scheduler); err != nil {
		return fmt.Errorf("scheduler: %w", err)
	}
	if err := validateLeaderElection(cfg.LeaderElection); err != nil {
		return fmt.Errorf("leader_election: %w", err)
	}
	if err := validateMySQL(cfg.MySQL); err != nil {
		return fmt.Errorf("mysql: %w", err)
	}
	if err := validateRedis(cfg.Redis); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if err := validateKubernetes(cfg.Kubernetes); err != nil {
		return fmt.Errorf("kubernetes: %w", err)
	}
	if err := validateObservability(cfg.Observability); err != nil {
		return fmt.Errorf("observability: %w", err)
	}
	if err := validateLogging(cfg.Logging); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func ValidateControllerConfig(cfg *ControllerAppConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if err := validateService(cfg.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if err := validateController(cfg.Controller); err != nil {
		return fmt.Errorf("controller: %w", err)
	}
	if err := validateLeaderElection(cfg.LeaderElection); err != nil {
		return fmt.Errorf("leader_election: %w", err)
	}
	if err := validateMySQL(cfg.MySQL); err != nil {
		return fmt.Errorf("mysql: %w", err)
	}
	if err := validateRedis(cfg.Redis); err != nil {
		return fmt.Errorf("redis: %w", err)
	}
	if err := validateKubernetes(cfg.Kubernetes); err != nil {
		return fmt.Errorf("kubernetes: %w", err)
	}
	if err := validateObservability(cfg.Observability); err != nil {
		return fmt.Errorf("observability: %w", err)
	}
	if err := validateLogging(cfg.Logging); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func ValidateWebhookConfig(cfg *WebhookAppConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if err := validateService(cfg.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if err := validateHTTPSServer(cfg.Server.HTTPS); err != nil {
		return fmt.Errorf("server.https: %w", err)
	}
	if err := validateWebhook(cfg.Webhook); err != nil {
		return fmt.Errorf("webhook: %w", err)
	}
	if err := validateKubernetes(cfg.Kubernetes); err != nil {
		return fmt.Errorf("kubernetes: %w", err)
	}
	if err := validateObservability(cfg.Observability); err != nil {
		return fmt.Errorf("observability: %w", err)
	}
	if err := validateLogging(cfg.Logging); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func ValidateAgentConfig(cfg *AgentConfig) error {
	if cfg == nil {
		return errors.New("config is nil")
	}
	if err := validateService(cfg.Service); err != nil {
		return fmt.Errorf("service: %w", err)
	}
	if err := validateAgent(cfg.Agent); err != nil {
		return fmt.Errorf("agent: %w", err)
	}
	if err := validateReporter(cfg.Reporter); err != nil {
		return fmt.Errorf("reporter: %w", err)
	}
	if err := validateKubernetes(cfg.Kubernetes); err != nil {
		return fmt.Errorf("kubernetes: %w", err)
	}
	if err := validateObservability(cfg.Observability); err != nil {
		return fmt.Errorf("observability: %w", err)
	}
	if err := validateLogging(cfg.Logging); err != nil {
		return fmt.Errorf("logging: %w", err)
	}
	return nil
}

func validateService(cfg ServiceConfig) error {
	if strings.TrimSpace(cfg.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(cfg.Env) == "" {
		return errors.New("env is required")
	}
	return nil
}

func validateHTTPServer(cfg HTTPServerConfig, addrRequired bool) error {
	if addrRequired && strings.TrimSpace(cfg.Addr) == "" {
		return errors.New("addr is required")
	}
	if cfg.ReadTimeout < 0 || cfg.WriteTimeout < 0 || cfg.IdleTimeout < 0 || cfg.ShutdownTimeout < 0 {
		return errors.New("timeouts must be >= 0")
	}
	return nil
}

func validateHTTPSServer(cfg HTTPSServerConfig) error {
	if strings.TrimSpace(cfg.Addr) == "" {
		return errors.New("addr is required")
	}
	if cfg.ReadTimeout < 0 || cfg.WriteTimeout < 0 || cfg.ShutdownTimeout < 0 {
		return errors.New("timeouts must be >= 0")
	}
	return nil
}

func validateSecurity(cfg SecurityConfig) error {
	if cfg.EnableAuthN {
		if strings.TrimSpace(cfg.JWT.Issuer) == "" {
			return errors.New("jwt issuer is required when authn is enabled")
		}
		if strings.TrimSpace(cfg.JWT.Secret) == "" {
			return errors.New("jwt secret is required when authn is enabled")
		}
		if cfg.JWT.Expire <= 0 {
			return errors.New("jwt expire must be > 0 when authn is enabled")
		}
	}
	return nil
}

func validateMySQL(cfg MySQLConfig) error {
	if strings.TrimSpace(cfg.DSN) == "" {
		return errors.New("dsn is required")
	}
	if cfg.MaxOpenConns <= 0 {
		return errors.New("max_open_conns must be > 0")
	}
	if cfg.MaxIdleConns < 0 {
		return errors.New("max_idle_conns must be >= 0")
	}
	if cfg.ConnMaxLifetime < 0 || cfg.ConnMaxIdleTime < 0 {
		return errors.New("connection durations must be >= 0")
	}
	return nil
}

func validateRedis(cfg RedisConfig) error {
	if strings.TrimSpace(cfg.Addr) == "" {
		return errors.New("addr is required")
	}
	if cfg.DB < 0 {
		return errors.New("db must be >= 0")
	}
	if cfg.DialTimeout < 0 || cfg.ReadTimeout < 0 || cfg.WriteTimeout < 0 {
		return errors.New("timeouts must be >= 0")
	}
	if cfg.PoolSize <= 0 {
		return errors.New("pool_size must be > 0")
	}
	if cfg.MinIdleConns < 0 {
		return errors.New("min_idle_conns must be >= 0")
	}
	return nil
}

func validateObservability(cfg ObservabilityConfig) error {
	if cfg.Server.Enabled && strings.TrimSpace(cfg.Server.Addr) == "" {
		return errors.New("observability.server.addr is required when observability server is enabled")
	}

	if cfg.Metrics.Enabled {
		if strings.TrimSpace(cfg.Metrics.Path) == "" {
			return errors.New("metrics path is required when metrics is enabled")
		}
		if !strings.HasPrefix(cfg.Metrics.Path, "/") {
			return errors.New("metrics path must start with '/'")
		}
	}

	if cfg.Tracing.SampleRatio < 0 || cfg.Tracing.SampleRatio > 1 {
		return errors.New("tracing sample_ratio must be in [0,1]")
	}

	if cfg.Tracing.Enabled {
		if strings.TrimSpace(cfg.Tracing.Endpoint) == "" {
			return errors.New("tracing endpoint is required when tracing is enabled")
		}
		if cfg.Tracing.Insecure && cfg.Tracing.TLS.Enabled {
			return errors.New("tracing insecure and tracing tls.enabled cannot both be true")
		}
		if cfg.Tracing.TLS.Enabled {
			if (cfg.Tracing.TLS.CertFile == "") != (cfg.Tracing.TLS.KeyFile == "") {
				return errors.New("tracing tls cert_file and key_file must be configured together")
			}
		}
	}

	if cfg.PProf.Enabled {
		if strings.TrimSpace(cfg.PProf.PathPrefix) == "" {
			return errors.New("pprof path_prefix is required when pprof is enabled")
		}
		if !strings.HasPrefix(cfg.PProf.PathPrefix, "/") {
			return errors.New("pprof path_prefix must start with '/'")
		}
	}

	return nil
}

func validateLogging(cfg LoggingConfig) error {
	level := strings.ToLower(strings.TrimSpace(cfg.Level))
	switch level {
	case "debug", "info", "warn", "error", "dpanic", "panic", "fatal":
	default:
		return fmt.Errorf("unsupported level %q", cfg.Level)
	}

	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	switch format {
	case "json", "console":
	default:
		return fmt.Errorf("unsupported format %q", cfg.Format)
	}

	if len(cfg.OutputPaths) == 0 {
		return errors.New("output_paths is required")
	}
	if len(cfg.ErrorOutputPaths) == 0 {
		return errors.New("error_output_paths is required")
	}
	return nil
}

func validateKubernetes(cfg KubernetesConfig) error {
	if cfg.QPS <= 0 {
		return errors.New("qps must be > 0")
	}
	if cfg.Burst <= 0 {
		return errors.New("burst must be > 0")
	}
	return nil
}

func validateLeaderElection(cfg LeaderElectionConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if strings.TrimSpace(cfg.LeaseName) == "" {
		return errors.New("lease_name is required when leader election is enabled")
	}
	if strings.TrimSpace(cfg.LeaseNamespace) == "" {
		return errors.New("lease_namespace is required when leader election is enabled")
	}
	if cfg.LeaseDuration <= 0 {
		return errors.New("lease_duration must be > 0")
	}
	if cfg.RenewDeadline <= 0 {
		return errors.New("renew_deadline must be > 0")
	}
	if cfg.RetryPeriod <= 0 {
		return errors.New("retry_period must be > 0")
	}
	if cfg.RenewDeadline >= cfg.LeaseDuration {
		return errors.New("renew_deadline must be less than lease_duration")
	}
	return nil
}

func validateSchedulerCore(cfg SchedulerCoreConfig) error {
	if cfg.ScheduleInterval <= 0 {
		return errors.New("schedule_interval must be > 0")
	}
	if cfg.MaxConcurrentCycles <= 0 {
		return errors.New("max_concurrent_cycles must be > 0")
	}
	if cfg.PendingBatchSize <= 0 {
		return errors.New("pending_batch_size must be > 0")
	}
	if cfg.ReservationTTL <= 0 {
		return errors.New("reservation_ttl must be > 0")
	}
	if cfg.BindTimeout <= 0 {
		return errors.New("bind_timeout must be > 0")
	}
	return nil
}

func validateController(cfg ControllerConfig) error {
	if cfg.Workers.GPUJob <= 0 {
		return errors.New("controller.workers.gpujob must be > 0")
	}
	if cfg.Workers.Node <= 0 {
		return errors.New("controller.workers.node must be > 0")
	}
	if cfg.Workers.Pod <= 0 {
		return errors.New("controller.workers.pod must be > 0")
	}
	if cfg.Workers.Quota <= 0 {
		return errors.New("controller.workers.quota must be > 0")
	}
	if cfg.Workers.Policy <= 0 {
		return errors.New("controller.workers.policy must be > 0")
	}
	if cfg.ResyncPeriod <= 0 {
		return errors.New("resync_period must be > 0")
	}
	if cfg.ReconcileTimeout <= 0 {
		return errors.New("reconcile_timeout must be > 0")
	}
	return nil
}

func validateWebhook(cfg WebhookConfig) error {
	fp := strings.ToLower(strings.TrimSpace(cfg.FailurePolicy))
	switch fp {
	case "fail", "ignore":
	default:
		return fmt.Errorf("unsupported failure_policy %q", cfg.FailurePolicy)
	}

	if strings.TrimSpace(cfg.SideEffects) == "" {
		return errors.New("side_effects is required")
	}
	return nil
}

func validateAgent(cfg AgentCoreConfig) error {
	if cfg.HeartbeatInterval <= 0 {
		return errors.New("heartbeat_interval must be > 0")
	}
	if cfg.CollectInterval <= 0 {
		return errors.New("collect_interval must be > 0")
	}
	if cfg.ReportTimeout <= 0 {
		return errors.New("report_timeout must be > 0")
	}
	if !cfg.EnableDCGM && !cfg.EnableNvidiaSMI {
		return errors.New("at least one of enable_dcgm or enable_nvidia_smi must be true")
	}
	return nil
}

func validateReporter(cfg ReporterConfig) error {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	switch mode {
	case "http":
		if strings.TrimSpace(cfg.HTTP.Endpoint) == "" {
			return errors.New("reporter.http.endpoint is required when reporter mode is http")
		}
		if cfg.HTTP.Timeout <= 0 {
			return errors.New("reporter.http.timeout must be > 0")
		}
	case "grpc":
		if strings.TrimSpace(cfg.GRPC.Endpoint) == "" {
			return errors.New("reporter.grpc.endpoint is required when reporter mode is grpc")
		}
		if cfg.GRPC.Timeout <= 0 {
			return errors.New("reporter.grpc.timeout must be > 0")
		}
	default:
		return fmt.Errorf("unsupported reporter mode %q", cfg.Mode)
	}
	return nil
}
