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
	applyCommonObservabilityEnv(&cfg.Observability)

	overrideString("API_SERVER_ADDR", &cfg.Server.HTTP.Addr)
	overrideDuration("API_SERVER_READ_TIMEOUT", &cfg.Server.HTTP.ReadTimeout)
	overrideDuration("API_SERVER_WRITE_TIMEOUT", &cfg.Server.HTTP.WriteTimeout)
	overrideDuration("API_SERVER_IDLE_TIMEOUT", &cfg.Server.HTTP.IdleTimeout)
	overrideDuration("API_SERVER_SHUTDOWN_TIMEOUT", &cfg.Server.HTTP.ShutdownTimeout)

	overrideBool("API_ENABLE_AUTHN", &cfg.Security.EnableAuthN)
	overrideBool("API_ENABLE_AUTHZ", &cfg.Security.EnableAuthZ)
	overrideString("API_JWT_ISSUER", &cfg.Security.JWT.Issuer)
	overrideString("API_JWT_SECRET", &cfg.Security.JWT.Secret)
	overrideDuration("API_JWT_EXPIRE", &cfg.Security.JWT.Expire)

	overrideBool("API_FEATURE_QUEUE", &cfg.Features.EnableQueueAPI)
	overrideBool("API_FEATURE_POLICY", &cfg.Features.EnablePolicyAPI)
	overrideBool("API_FEATURE_QUOTA", &cfg.Features.EnableQuotaAPI)
	overrideBool("API_FEATURE_TENANT", &cfg.Features.EnableTenantAPI)
	overrideBool("API_FEATURE_CLUSTER", &cfg.Features.EnableClusterAPI)

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
	applyCommonObservabilityEnv(&cfg.Observability)
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
	applyCommonObservabilityEnv(&cfg.Observability)
	applyCommonKubernetesEnv(&cfg.Kubernetes)
	applyLeaderElectionEnv(&cfg.LeaderElection)
	return nil
}

func ApplyWebhookEnvOverrides(cfg *WebhookAppConfig) error {
	if cfg == nil {
		return nil
	}
	applyCommonServiceEnv(&cfg.Service)
	applyCommonLoggingEnv(&cfg.Logging)
	applyCommonObservabilityEnv(&cfg.Observability)
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
	applyCommonObservabilityEnv(&cfg.Observability)
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

func applyCommonObservabilityEnv(cfg *ObservabilityConfig) {
	overrideBool("METRICS_ENABLED", &cfg.Metrics.Enabled)
	overrideString("METRICS_ADDR", &cfg.Metrics.Addr)
	overrideString("METRICS_PATH", &cfg.Metrics.Path)

	overrideBool("TRACING_ENABLED", &cfg.Tracing.Enabled)
	overrideString("TRACING_ENDPOINT", &cfg.Tracing.Endpoint)
	overrideFloat64("TRACING_SAMPLE_RATIO", &cfg.Tracing.SampleRatio)

	overrideBool("PPROF_ENABLED", &cfg.PProf.Enabled)
	overrideString("PPROF_ADDR", &cfg.PProf.Addr)
	overrideString("PPROF_PATH_PREFIX", &cfg.PProf.PathPrefix)
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

func overrideString(env string, target *string) {
	if target == nil {
		return
	}
	if v, ok := os.LookupEnv(env); ok {
		*target = strings.TrimSpace(v)
	}
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
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}
