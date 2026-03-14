package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func ApplyAPIServerEnvOverrides(cfg *APIServerConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonMySQLEnv(&cfg.MySQL)
	applyCommonRedisEnv(&cfg.Redis)
	applyCommonLoggingEnv(&cfg.Logging)
	if err := applyCommonObservabilityEnv(&cfg.Observability); err != nil {
		return err
	}

	overrideString("API_HTTP_ADDR", &cfg.Server.HTTP.Addr)
	overrideDuration("API_HTTP_READ_TIMEOUT", &cfg.Server.HTTP.ReadTimeout)
	overrideDuration("API_HTTP_WRITE_TIMEOUT", &cfg.Server.HTTP.WriteTimeout)
	overrideDuration("API_HTTP_IDLE_TIMEOUT", &cfg.Server.HTTP.IdleTimeout)
	overrideDuration("API_HTTP_SHUTDOWN_TIMEOUT", &cfg.Server.HTTP.ShutdownTimeout)

	overrideBool("API_ENABLE_AUTHN", &cfg.Security.EnableAuthN)
	overrideBool("API_ENABLE_AUTHZ", &cfg.Security.EnableAuthZ)
	overrideString("API_JWT_ISSUER", &cfg.Security.JWT.Issuer)
	overrideString("API_JWT_SECRET", &cfg.Security.JWT.Secret)
	overrideDuration("API_JWT_EXPIRE", &cfg.Security.JWT.Expire)

	overrideBool("API_ENABLE_QUEUE_API", &cfg.Features.EnableQueueAPI)
	overrideBool("API_ENABLE_POLICY_API", &cfg.Features.EnablePolicyAPI)
	overrideBool("API_ENABLE_QUOTA_API", &cfg.Features.EnableQuotaAPI)
	overrideBool("API_ENABLE_TENANT_API", &cfg.Features.EnableTenantAPI)
	overrideBool("API_ENABLE_CLUSTER_API", &cfg.Features.EnableClusterAPI)

	return nil
}

func ApplySchedulerEnvOverrides(cfg *SchedulerConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonMySQLEnv(&cfg.MySQL)
	applyCommonRedisEnv(&cfg.Redis)
	applyCommonLoggingEnv(&cfg.Logging)
	if err := applyCommonObservabilityEnv(&cfg.Observability); err != nil {
		return err
	}
	applyCommonKubernetesEnv(&cfg.Kubernetes)
	applyLeaderElectionEnv(&cfg.LeaderElection)

	overrideDuration("SCHEDULER_INTERVAL", &cfg.Scheduler.ScheduleInterval)
	overrideInt("SCHEDULER_MAX_CONCURRENT_CYCLES", &cfg.Scheduler.MaxConcurrentCycles)
	overrideInt("SCHEDULER_PENDING_BATCH_SIZE", &cfg.Scheduler.PendingBatchSize)
	overrideDuration("SCHEDULER_RESERVATION_TTL", &cfg.Scheduler.ReservationTTL)
	overrideDuration("SCHEDULER_BIND_TIMEOUT", &cfg.Scheduler.BindTimeout)
	overrideBool("SCHEDULER_ENABLE_PREEMPTION", &cfg.Scheduler.EnablePreemption)
	overrideBool("SCHEDULER_ENABLE_GANG_SCHEDULING", &cfg.Scheduler.EnableGangScheduling)
	overrideBool("SCHEDULER_ENABLE_TOPOLOGY_AWARE", &cfg.Scheduler.EnableTopologyAware)
	overrideBool("SCHEDULER_ENABLE_MIG", &cfg.Scheduler.EnableMIG)

	return nil
}

func ApplyControllerEnvOverrides(cfg *ControllerAppConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonMySQLEnv(&cfg.MySQL)
	applyCommonRedisEnv(&cfg.Redis)
	applyCommonLoggingEnv(&cfg.Logging)
	if err := applyCommonObservabilityEnv(&cfg.Observability); err != nil {
		return err
	}
	applyCommonKubernetesEnv(&cfg.Kubernetes)
	applyLeaderElectionEnv(&cfg.LeaderElection)

	overrideInt("CONTROLLER_WORKERS_GPUJOB", &cfg.Controller.Workers.GPUJob)
	overrideInt("CONTROLLER_WORKERS_NODE", &cfg.Controller.Workers.Node)
	overrideInt("CONTROLLER_WORKERS_POD", &cfg.Controller.Workers.Pod)
	overrideInt("CONTROLLER_WORKERS_QUOTA", &cfg.Controller.Workers.Quota)
	overrideInt("CONTROLLER_WORKERS_POLICY", &cfg.Controller.Workers.Policy)
	overrideDuration("CONTROLLER_RESYNC_PERIOD", &cfg.Controller.ResyncPeriod)
	overrideDuration("CONTROLLER_RECONCILE_TIMEOUT", &cfg.Controller.ReconcileTimeout)
	overrideBool("CONTROLLER_ENABLE_ALLOCATION_RECOVERY", &cfg.Controller.EnableAllocationRecovery)
	overrideBool("CONTROLLER_ENABLE_JOB_STATUS_SYNC", &cfg.Controller.EnableJobStatusSync)

	return nil
}

func ApplyWebhookEnvOverrides(cfg *WebhookAppConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonLoggingEnv(&cfg.Logging)
	if err := applyCommonObservabilityEnv(&cfg.Observability); err != nil {
		return err
	}
	applyCommonKubernetesEnv(&cfg.Kubernetes)

	overrideString("WEBHOOK_ADDR", &cfg.Server.HTTPS.Addr)
	overrideString("WEBHOOK_CERT_FILE", &cfg.Server.HTTPS.CertFile)
	overrideString("WEBHOOK_KEY_FILE", &cfg.Server.HTTPS.KeyFile)
	overrideDuration("WEBHOOK_READ_TIMEOUT", &cfg.Server.HTTPS.ReadTimeout)
	overrideDuration("WEBHOOK_WRITE_TIMEOUT", &cfg.Server.HTTPS.WriteTimeout)
	overrideDuration("WEBHOOK_SHUTDOWN_TIMEOUT", &cfg.Server.HTTPS.ShutdownTimeout)
	overrideBool("WEBHOOK_ENABLE_MUTATING", &cfg.Webhook.EnableMutating)
	overrideBool("WEBHOOK_ENABLE_VALIDATING", &cfg.Webhook.EnableValidating)
	overrideString("WEBHOOK_FAILURE_POLICY", &cfg.Webhook.FailurePolicy)
	overrideString("WEBHOOK_SIDE_EFFECTS", &cfg.Webhook.SideEffects)

	return nil
}

func ApplyAgentEnvOverrides(cfg *AgentConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonLoggingEnv(&cfg.Logging)
	if err := applyCommonObservabilityEnv(&cfg.Observability); err != nil {
		return err
	}
	applyCommonKubernetesEnv(&cfg.Kubernetes)

	overrideString("AGENT_NODE_NAME", &cfg.Agent.NodeName)
	overrideDuration("AGENT_HEARTBEAT_INTERVAL", &cfg.Agent.HeartbeatInterval)
	overrideDuration("AGENT_COLLECT_INTERVAL", &cfg.Agent.CollectInterval)
	overrideDuration("AGENT_REPORT_TIMEOUT", &cfg.Agent.ReportTimeout)
	overrideBool("AGENT_ENABLE_DCGM", &cfg.Agent.EnableDCGM)
	overrideBool("AGENT_ENABLE_NVIDIA_SMI", &cfg.Agent.EnableNvidiaSMI)
	overrideBool("AGENT_ENABLE_MIG_DISCOVERY", &cfg.Agent.EnableMIGDiscovery)
	overrideBool("AGENT_ENABLE_TOPOLOGY_DISCOVERY", &cfg.Agent.EnableTopologyDiscovery)

	overrideString("AGENT_REPORTER_MODE", &cfg.Reporter.Mode)
	overrideString("AGENT_REPORTER_HTTP_ENDPOINT", &cfg.Reporter.HTTP.Endpoint)
	overrideDuration("AGENT_REPORTER_HTTP_TIMEOUT", &cfg.Reporter.HTTP.Timeout)
	overrideString("AGENT_REPORTER_GRPC_ENDPOINT", &cfg.Reporter.GRPC.Endpoint)
	overrideDuration("AGENT_REPORTER_GRPC_TIMEOUT", &cfg.Reporter.GRPC.Timeout)

	return nil
}

func applyCommonServiceEnv(cfg *ServiceConfig) {
	overrideString("APP_NAME", &cfg.Name)
	overrideString("APP_ENV", &cfg.Env)
	overrideString("APP_VERSION", &cfg.Version)
}

func applyCommonMySQLEnv(cfg *MySQLConfig) {
	overrideString("MYSQL_DSN", &cfg.DSN)
	overrideInt("MYSQL_MAX_OPEN_CONNS", &cfg.MaxOpenConns)
	overrideInt("MYSQL_MAX_IDLE_CONNS", &cfg.MaxIdleConns)
	overrideDuration("MYSQL_CONN_MAX_LIFETIME", &cfg.ConnMaxLifetime)
	overrideDuration("MYSQL_CONN_MAX_IDLE_TIME", &cfg.ConnMaxIdleTime)
}

func applyCommonRedisEnv(cfg *RedisConfig) {
	overrideString("REDIS_ADDR", &cfg.Addr)
	overrideString("REDIS_PASSWORD", &cfg.Password)
	overrideInt("REDIS_DB", &cfg.DB)
	overrideDuration("REDIS_DIAL_TIMEOUT", &cfg.DialTimeout)
	overrideDuration("REDIS_READ_TIMEOUT", &cfg.ReadTimeout)
	overrideDuration("REDIS_WRITE_TIMEOUT", &cfg.WriteTimeout)
	overrideInt("REDIS_POOL_SIZE", &cfg.PoolSize)
	overrideInt("REDIS_MIN_IDLE_CONNS", &cfg.MinIdleConns)
}

func applyCommonLoggingEnv(cfg *LoggingConfig) {
	overrideString("LOG_LEVEL", &cfg.Level)
	overrideString("LOG_FORMAT", &cfg.Format)

	if v, ok := os.LookupEnv("LOG_OUTPUT_PATHS"); ok && strings.TrimSpace(v) != "" {
		cfg.OutputPaths = splitCSV(v)
	}
	if v, ok := os.LookupEnv("LOG_ERROR_OUTPUT_PATHS"); ok && strings.TrimSpace(v) != "" {
		cfg.ErrorOutputPaths = splitCSV(v)
	}
}

func applyCommonObservabilityEnv(cfg *ObservabilityConfig) error {
	overrideBool("OBS_SERVER_ENABLED", &cfg.Server.Enabled)
	overrideString("OBS_SERVER_ADDR", &cfg.Server.Addr)

	overrideBool("METRICS_ENABLED", &cfg.Metrics.Enabled)
	overrideString("METRICS_PATH", &cfg.Metrics.Path)

	overrideBool("TRACING_ENABLED", &cfg.Tracing.Enabled)
	overrideString("TRACING_ENDPOINT", &cfg.Tracing.Endpoint)
	overrideFloat64("TRACING_SAMPLE_RATIO", &cfg.Tracing.SampleRatio)
	overrideBool("TRACING_INSECURE", &cfg.Tracing.Insecure)

	if v, ok := os.LookupEnv("TRACING_HEADERS"); ok && strings.TrimSpace(v) != "" {
		headers, err := parseEnvMap(v)
		if err != nil {
			return fmt.Errorf("parse TRACING_HEADERS: %w", err)
		}
		cfg.Tracing.Headers = headers
	}

	overrideBool("TRACING_TLS_ENABLED", &cfg.Tracing.TLS.Enabled)
	overrideString("TRACING_TLS_CERT_FILE", &cfg.Tracing.TLS.CertFile)
	overrideString("TRACING_TLS_KEY_FILE", &cfg.Tracing.TLS.KeyFile)
	overrideString("TRACING_TLS_CA_FILE", &cfg.Tracing.TLS.CAFile)
	overrideBool("TRACING_TLS_INSECURE_SKIP_VERIFY", &cfg.Tracing.TLS.InsecureSkipVerify)
	overrideString("TRACING_TLS_SERVER_NAME", &cfg.Tracing.TLS.ServerName)

	overrideBool("PPROF_ENABLED", &cfg.PProf.Enabled)
	overrideString("PPROF_PATH_PREFIX", &cfg.PProf.PathPrefix)

	return nil
}

func applyCommonKubernetesEnv(cfg *KubernetesConfig) {
	overrideBool("K8S_IN_CLUSTER", &cfg.InCluster)
	overrideString("K8S_KUBECONFIG", &cfg.Kubeconfig)
	overrideInt("K8S_QPS", &cfg.QPS)
	overrideInt("K8S_BURST", &cfg.Burst)
}

func applyLeaderElectionEnv(cfg *LeaderElectionConfig) {
	overrideBool("LEADER_ELECTION_ENABLED", &cfg.Enabled)
	overrideString("LEADER_ELECTION_LEASE_NAME", &cfg.LeaseName)
	overrideString("LEADER_ELECTION_LEASE_NAMESPACE", &cfg.LeaseNamespace)
	overrideDuration("LEADER_ELECTION_LEASE_DURATION", &cfg.LeaseDuration)
	overrideDuration("LEADER_ELECTION_RENEW_DEADLINE", &cfg.RenewDeadline)
	overrideDuration("LEADER_ELECTION_RETRY_PERIOD", &cfg.RetryPeriod)
}

func overrideString(env string, target *string) error {
	if target == nil {
		return nil
	}
	if v, ok := os.LookupEnv(env); ok {
		*target = strings.TrimSpace(v)
	}
	return nil
}

func overrideBool(env string, target *bool) error {
	if target == nil {
		return nil
	}
	v, ok := os.LookupEnv(env)
	if !ok {
		return nil
	}
	parsed, err := strconv.ParseBool(strings.TrimSpace(v))
	if err != nil {
		return fmt.Errorf("parse bool env %s=%q: %w", env, v, err)
	}
	*target = parsed
	return nil
}

func overrideInt(env string, target *int) error {
	if target == nil {
		return nil
	}
	v, ok := os.LookupEnv(env)
	if !ok {
		return nil
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return fmt.Errorf("parse int env %s=%q: %w", env, v, err)
	}
	*target = parsed
	return nil
}

func overrideFloat64(env string, target *float64) error {
	if target == nil {
		return nil
	}
	v, ok := os.LookupEnv(env)
	if !ok {
		return nil
	}
	parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return fmt.Errorf("parse float env %s=%q: %w", env, v, err)
	}
	*target = parsed
	return nil
}

func overrideDuration(env string, target *time.Duration) error {
	if target == nil {
		return nil
	}
	v, ok := os.LookupEnv(env)
	if !ok {
		return nil
	}
	parsed, err := time.ParseDuration(strings.TrimSpace(v))
	if err != nil {
		return fmt.Errorf("parse duration env %s=%q: %w", env, v, err)
	}
	*target = parsed
	return nil
}

func splitCSV(s string) []string {
	items := strings.Split(s, ",")
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func parseEnvMap(s string) (map[string]string, error) {
	out := make(map[string]string)
	items := strings.Split(s, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid item %q, expect key=value", item)
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			return nil, fmt.Errorf("empty key in %q", item)
		}
		out[k] = v
	}
	return out, nil
}
