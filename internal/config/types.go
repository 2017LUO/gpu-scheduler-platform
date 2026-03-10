package config

import "time"

type ServiceConfig struct {
	Name    string `yaml:"name"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
}

type HTTPServerConfig struct {
	Addr            string        `yaml:"addr"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	IdleTimeout     time.Duration `yaml:"idle_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HTTPServersConfig struct {
	HTTP HTTPServerConfig `yaml:"http"`
}

type HTTPSServerConfig struct {
	Addr            string        `yaml:"addr"`
	CertFile        string        `yaml:"cert_file"`
	KeyFile         string        `yaml:"key_file"`
	ReadTimeout     time.Duration `yaml:"read_timeout"`
	WriteTimeout    time.Duration `yaml:"write_timeout"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
}

type HTTPSServersConfig struct {
	HTTPS HTTPSServerConfig `yaml:"https"`
}

type JWTConfig struct {
	Issuer string        `yaml:"issuer"`
	Secret string        `yaml:"secret"`
	Expire time.Duration `yaml:"expire"`
}

type SecurityConfig struct {
	EnableAuthN bool      `yaml:"enable_authn"`
	EnableAuthZ bool      `yaml:"enable_authz"`
	JWT         JWTConfig `yaml:"jwt"`
}

type MySQLConfig struct {
	DSN             string        `yaml:"dsn"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

type RedisConfig struct {
	Addr         string        `yaml:"addr"`
	Password     string        `yaml:"password"`
	DB           int           `yaml:"db"`
	DialTimeout  time.Duration `yaml:"dial_timeout"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	PoolSize     int           `yaml:"pool_size"`
	MinIdleConns int           `yaml:"min_idle_conns"`
}

type MetricsConfig struct {
	Enabled bool   `yaml:"enabled"`
	Addr    string `yaml:"addr"`
	Path    string `yaml:"path"`
}

type TracingConfig struct {
	Enabled     bool    `yaml:"enabled"`
	Endpoint    string  `yaml:"endpoint"`
	SampleRatio float64 `yaml:"sample_ratio"`
}

type PProfConfig struct {
	Enabled    bool   `yaml:"enabled"`
	Addr       string `yaml:"addr"`
	PathPrefix string `yaml:"path_prefix"`
}

type ObservabilityConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Tracing TracingConfig `yaml:"tracing"`
	PProf   PProfConfig   `yaml:"pprof"`
}

type LoggingConfig struct {
	Level            string   `yaml:"level"`
	Format           string   `yaml:"format"`
	OutputPaths      []string `yaml:"output_paths"`
	ErrorOutputPaths []string `yaml:"error_output_paths"`
}

type KubernetesConfig struct {
	InCluster  bool   `yaml:"in_cluster"`
	Kubeconfig string `yaml:"kubeconfig"`
	QPS        int    `yaml:"qps"`
	Burst      int    `yaml:"burst"`
}

type LeaderElectionConfig struct {
	Enabled        bool          `yaml:"enabled"`
	LeaseName      string        `yaml:"lease_name"`
	LeaseNamespace string        `yaml:"lease_namespace"`
	LeaseDuration  time.Duration `yaml:"lease_duration"`
	RenewDeadline  time.Duration `yaml:"renew_deadline"`
	RetryPeriod    time.Duration `yaml:"retry_period"`
}

type APIFeaturesConfig struct {
	EnableQueueAPI   bool `yaml:"enable_queue_api"`
	EnablePolicyAPI  bool `yaml:"enable_policy_api"`
	EnableQuotaAPI   bool `yaml:"enable_quota_api"`
	EnableTenantAPI  bool `yaml:"enable_tenant_api"`
	EnableClusterAPI bool `yaml:"enable_cluster_api"`
}

type SchedulerCoreConfig struct {
	ScheduleInterval     time.Duration `yaml:"schedule_interval"`
	MaxConcurrentCycles  int           `yaml:"max_concurrent_cycles"`
	PendingBatchSize     int           `yaml:"pending_batch_size"`
	ReservationTTL       time.Duration `yaml:"reservation_ttl"`
	BindTimeout          time.Duration `yaml:"bind_timeout"`
	EnablePreemption     bool          `yaml:"enable_preemption"`
	EnableGangScheduling bool          `yaml:"enable_gang_scheduling"`
	EnableTopologyAware  bool          `yaml:"enable_topology_aware"`
	EnableMIG            bool          `yaml:"enable_mig"`
}

type ControllerWorkersConfig struct {
	GPUJob int `yaml:"gpujob"`
	Node   int `yaml:"node"`
	Pod    int `yaml:"pod"`
	Quota  int `yaml:"quota"`
	Policy int `yaml:"policy"`
}

type ControllerConfig struct {
	Workers                  ControllerWorkersConfig `yaml:"workers"`
	ResyncPeriod             time.Duration           `yaml:"resync_period"`
	ReconcileTimeout         time.Duration           `yaml:"reconcile_timeout"`
	EnableAllocationRecovery bool                    `yaml:"enable_allocation_recovery"`
	EnableJobStatusSync      bool                    `yaml:"enable_job_status_sync"`
}

type WebhookConfig struct {
	EnableMutating   bool   `yaml:"enable_mutating"`
	EnableValidating bool   `yaml:"enable_validating"`
	FailurePolicy    string `yaml:"failure_policy"`
	SideEffects      string `yaml:"side_effects"`
}

type AgentCoreConfig struct {
	NodeName                string        `yaml:"node_name"`
	HeartbeatInterval       time.Duration `yaml:"heartbeat_interval"`
	CollectInterval         time.Duration `yaml:"collect_interval"`
	ReportTimeout           time.Duration `yaml:"report_timeout"`
	EnableDCGM              bool          `yaml:"enable_dcgm"`
	EnableNvidiaSMI         bool          `yaml:"enable_nvidia_smi"`
	EnableMIGDiscovery      bool          `yaml:"enable_mig_discovery"`
	EnableTopologyDiscovery bool          `yaml:"enable_topology_discovery"`
}

type ReporterHTTPConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Timeout  time.Duration `yaml:"timeout"`
}

type ReporterGRPCConfig struct {
	Endpoint string        `yaml:"endpoint"`
	Timeout  time.Duration `yaml:"timeout"`
}

type ReporterConfig struct {
	Mode string             `yaml:"mode"`
	HTTP ReporterHTTPConfig `yaml:"http"`
	GRPC ReporterGRPCConfig `yaml:"grpc"`
}

type APIServerConfig struct {
	Service       ServiceConfig       `yaml:"service"`
	Server        HTTPServersConfig   `yaml:"server"`
	Security      SecurityConfig      `yaml:"security"`
	MySQL         MySQLConfig         `yaml:"mysql"`
	Redis         RedisConfig         `yaml:"redis"`
	Observability ObservabilityConfig `yaml:"observability"`
	Logging       LoggingConfig       `yaml:"logging"`
	Features      APIFeaturesConfig   `yaml:"features"`
}

type SchedulerConfig struct {
	Service        ServiceConfig        `yaml:"service"`
	Scheduler      SchedulerCoreConfig  `yaml:"scheduler"`
	LeaderElection LeaderElectionConfig `yaml:"leader_election"`
	MySQL          MySQLConfig          `yaml:"mysql"`
	Redis          RedisConfig          `yaml:"redis"`
	Kubernetes     KubernetesConfig     `yaml:"kubernetes"`
	Observability  ObservabilityConfig  `yaml:"observability"`
	Logging        LoggingConfig        `yaml:"logging"`
}

type ControllerAppConfig struct {
	Service        ServiceConfig        `yaml:"service"`
	Controller     ControllerConfig     `yaml:"controller"`
	LeaderElection LeaderElectionConfig `yaml:"leader_election"`
	MySQL          MySQLConfig          `yaml:"mysql"`
	Redis          RedisConfig          `yaml:"redis"`
	Kubernetes     KubernetesConfig     `yaml:"kubernetes"`
	Observability  ObservabilityConfig  `yaml:"observability"`
	Logging        LoggingConfig        `yaml:"logging"`
}

type WebhookAppConfig struct {
	Service       ServiceConfig       `yaml:"service"`
	Server        HTTPSServersConfig  `yaml:"server"`
	Webhook       WebhookConfig       `yaml:"webhook"`
	Kubernetes    KubernetesConfig    `yaml:"kubernetes"`
	Observability ObservabilityConfig `yaml:"observability"`
	Logging       LoggingConfig       `yaml:"logging"`
}

type AgentConfig struct {
	Service       ServiceConfig       `yaml:"service"`
	Agent         AgentCoreConfig     `yaml:"agent"`
	Reporter      ReporterConfig      `yaml:"reporter"`
	Kubernetes    KubernetesConfig    `yaml:"kubernetes"`
	Observability ObservabilityConfig `yaml:"observability"`
	Logging       LoggingConfig       `yaml:"logging"`
}
