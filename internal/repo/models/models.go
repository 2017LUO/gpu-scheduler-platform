package model

import (
	"time"

	"gorm.io/datatypes"
)

//
// 0. tenants
//

type Tenant struct {
	ID          string    `gorm:"column:id;type:varchar(64);primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(128);not null"`
	Enabled     bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1;index:idx_tenants_enabled"`
	Description *string   `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_tenants_created"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Tenant) TableName() string { return "tenants" }

//
// 1. queues
//

type Queue struct {
	ID          string    `gorm:"column:id;type:varchar(64);primaryKey"`
	Name        string    `gorm:"column:name;type:varchar(64);not null;uniqueIndex:uk_queues_name_tenant,priority:1"`
	TenantID    string    `gorm:"column:tenant_id;type:varchar(64);not null;default:'';uniqueIndex:uk_queues_name_tenant,priority:2;index:idx_queues_tenant"`
	Weight      int       `gorm:"column:weight;type:int;not null;default:1"`
	Priority    int       `gorm:"column:priority;type:int;not null;default:0;index:idx_queues_priority"`
	Enabled     bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1;index:idx_queues_enabled"`
	Description *string   `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Queue) TableName() string { return "queues" }

//
// 2. nodes
//

type Node struct {
	NodeName          string         `gorm:"column:node_name;type:varchar(128);primaryKey"`
	ClusterName       string         `gorm:"column:cluster_name;type:varchar(128);not null;default:'';index:idx_nodes_cluster"`
	Source            string         `gorm:"column:source;type:varchar(32);not null;default:'agent'"`
	State             string         `gorm:"column:state;type:varchar(32);not null;index:idx_nodes_state"`
	Schedulable       bool           `gorm:"column:schedulable;type:tinyint(1);not null;default:1;index:idx_nodes_schedulable"`
	GPUCount          int            `gorm:"column:gpu_count;type:int;not null;default:0"`
	HealthyGPUCount   int            `gorm:"column:healthy_gpu_count;type:int;not null;default:0"`
	TotalMemoryMiB    int64          `gorm:"column:total_memory_mib;type:bigint;not null;default:0"`
	FreeMemoryMiB     int64          `gorm:"column:free_memory_mib;type:bigint;not null;default:0"`
	LabelsJSON        datatypes.JSON `gorm:"column:labels_json;type:json"`
	AnnotationsJSON   datatypes.JSON `gorm:"column:annotations_json;type:json"`
	TopologyJSON      datatypes.JSON `gorm:"column:topology_json;type:json"`
	LastReportTime    *time.Time     `gorm:"column:last_report_time;type:timestamp;index:idx_nodes_last_report_time"`
	LastHeartbeatTime *time.Time     `gorm:"column:last_heartbeat_time;type:timestamp;index:idx_nodes_last_heartbeat_time"`
	CreatedAt         time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Node) TableName() string { return "nodes" }

//
// 3. gpu_jobs
//

type GPUJob struct {
	ID                      string         `gorm:"column:id;type:varchar(64);primaryKey"`
	TenantID                string         `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_jobs_tenant;index:idx_jobs_tenant_status_created,priority:1"`
	Namespace               string         `gorm:"column:namespace;type:varchar(128);not null;uniqueIndex:uk_gpu_jobs_ns_name,priority:1;index:idx_jobs_namespace"`
	Name                    string         `gorm:"column:name;type:varchar(128);not null;uniqueIndex:uk_gpu_jobs_ns_name,priority:2"`
	Queue                   string         `gorm:"column:queue;type:varchar(64);not null;index:idx_jobs_queue;index:idx_jobs_queue_status_created,priority:1"`
	Priority                string         `gorm:"column:priority;type:varchar(32);not null"`
	Status                  string         `gorm:"column:status;type:varchar(32);not null;index:idx_jobs_status;index:idx_jobs_queue_status_created,priority:2;index:idx_jobs_tenant_status_created,priority:2"`
	Submitter               string         `gorm:"column:submitter;type:varchar(128);not null;default:''"`
	SchedulerName           string         `gorm:"column:scheduler_name;type:varchar(64);not null;default:'default';index:idx_jobs_scheduler"`
	GPUCount                int            `gorm:"column:gpu_count;type:int;not null"`
	GPUMemoryMiB            int64          `gorm:"column:gpu_memory_mib;type:bigint;not null"`
	GPUModel                string         `gorm:"column:gpu_model;type:varchar(128);not null;default:''"`
	RequireSameNode         bool           `gorm:"column:require_same_node;type:tinyint(1);not null;default:0"`
	RequireHealthy          bool           `gorm:"column:require_healthy;type:tinyint(1);not null;default:1"`
	RequireMIG              bool           `gorm:"column:require_mig;type:tinyint(1);not null;default:0"`
	MIGProfile              string         `gorm:"column:mig_profile;type:varchar(64);not null;default:''"`
	RequireNVLink           bool           `gorm:"column:require_nvlink;type:tinyint(1);not null;default:0"`
	Preemptible             bool           `gorm:"column:preemptible;type:tinyint(1);not null;default:0"`
	Retryable               bool           `gorm:"column:retryable;type:tinyint(1);not null;default:1"`
	MaxRetry                int            `gorm:"column:max_retry;type:int;not null;default:0"`
	ExpectedDurationSec     int64          `gorm:"column:expected_duration_sec;type:bigint;not null;default:0"`
	RunPolicyJSON           datatypes.JSON `gorm:"column:run_policy_json;type:json"`
	PreferredNodeLabelsJSON datatypes.JSON `gorm:"column:preferred_node_labels_json;type:json"`
	PreferredGPULabelsJSON  datatypes.JSON `gorm:"column:preferred_gpu_labels_json;type:json"`
	LabelsJSON              datatypes.JSON `gorm:"column:labels_json;type:json"`
	AnnotationsJSON         datatypes.JSON `gorm:"column:annotations_json;type:json"`
	RetryCount              int            `gorm:"column:retry_count;type:int;not null;default:0"`
	Message                 *string        `gorm:"column:message;type:text"`
	ScheduledAt             *time.Time     `gorm:"column:scheduled_at;type:timestamp"`
	StartedAt               *time.Time     `gorm:"column:started_at;type:timestamp"`
	FinishedAt              *time.Time     `gorm:"column:finished_at;type:timestamp"`
	CreatedAt               time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_jobs_created;index:idx_jobs_queue_status_created,priority:3;index:idx_jobs_tenant_status_created,priority:3"`
	UpdatedAt               time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (GPUJob) TableName() string { return "gpu_jobs" }

//
// 4. gpu_job_events
//

type GPUJobEvent struct {
	ID         string    `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string    `gorm:"column:job_id;type:varchar(64);not null;index:idx_job_events_job"`
	TenantID   string    `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_job_events_tenant"`
	Reason     string    `gorm:"column:reason;type:varchar(64);not null;index:idx_job_events_reason"`
	Message    *string   `gorm:"column:message;type:text"`
	Source     string    `gorm:"column:source;type:varchar(64);not null"`
	OccurredAt time.Time `gorm:"column:occurred_at;type:timestamp;not null;index:idx_job_events_occurred"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Job    GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (GPUJobEvent) TableName() string { return "gpu_job_events" }

//
// 5. node_snapshots
//

type NodeSnapshot struct {
	ID              uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	Version         string         `gorm:"column:version;type:varchar(64);not null;index:idx_node_snapshots_version"`
	AgentVersion    string         `gorm:"column:agent_version;type:varchar(64);not null;default:''"`
	ClusterName     string         `gorm:"column:cluster_name;type:varchar(128);not null;default:''"`
	NodeName        string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_node_snapshots_node;index:idx_node_snapshots_node_report,priority:1"`
	Source          string         `gorm:"column:source;type:varchar(32);not null;default:'agent'"`
	NodeState       string         `gorm:"column:node_state;type:varchar(32);not null"`
	Schedulable     bool           `gorm:"column:schedulable;type:tinyint(1);not null;default:1"`
	GPUCount        int            `gorm:"column:gpu_count;type:int;not null;default:0"`
	HealthyGPUCount int            `gorm:"column:healthy_gpu_count;type:int;not null;default:0"`
	TotalMemoryMiB  int64          `gorm:"column:total_memory_mib;type:bigint;not null;default:0"`
	FreeMemoryMiB   int64          `gorm:"column:free_memory_mib;type:bigint;not null;default:0"`
	LabelsJSON      datatypes.JSON `gorm:"column:labels_json;type:json"`
	AnnotationsJSON datatypes.JSON `gorm:"column:annotations_json;type:json"`
	TopologyJSON    datatypes.JSON `gorm:"column:topology_json;type:json"`
	ReportTime      time.Time      `gorm:"column:report_time;type:timestamp;not null;index:idx_node_snapshots_report_time;index:idx_node_snapshots_node_report,priority:2"`
	CreatedAt       time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_node_snapshots_created"`
}

func (NodeSnapshot) TableName() string { return "node_snapshots" }

//
// 6. gpu_devices
//

type GPUDevice struct {
	ID                uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	SnapshotID        uint64         `gorm:"column:snapshot_id;type:bigint unsigned;not null;index:idx_gpu_devices_snapshot;index:idx_gpu_devices_snapshot_node,priority:1;index:idx_gpu_devices_snapshot_health,priority:1;index:idx_gpu_devices_snapshot_allocated,priority:1"`
	NodeName          string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_gpu_devices_node;index:idx_gpu_devices_node_uuid,priority:1;index:idx_gpu_devices_snapshot_node,priority:2"`
	UUID              string         `gorm:"column:uuid;type:varchar(128);not null;index:idx_gpu_devices_uuid;index:idx_gpu_devices_node_uuid,priority:2"`
	GPUIndex          int            `gorm:"column:gpu_index;type:int;not null"`
	Model             string         `gorm:"column:model;type:varchar(128);not null"`
	Vendor            string         `gorm:"column:vendor;type:varchar(64);not null;default:'nvidia'"`
	Type              string         `gorm:"column:type;type:varchar(16);not null"`
	MemoryMiB         int64          `gorm:"column:memory_mib;type:bigint;not null"`
	FreeMemoryMiB     int64          `gorm:"column:free_memory_mib;type:bigint;not null"`
	Healthy           bool           `gorm:"column:healthy;type:tinyint(1);not null;default:1;index:idx_gpu_devices_snapshot_health,priority:2"`
	Health            string         `gorm:"column:health;type:varchar(32);not null"`
	MIGEnabled        bool           `gorm:"column:mig_enabled;type:tinyint(1);not null;default:0"`
	MIGProfile        string         `gorm:"column:mig_profile;type:varchar(64);not null;default:''"`
	UtilizationGPU    float64        `gorm:"column:utilization_gpu;type:double;not null;default:0"`
	UtilizationMemory float64        `gorm:"column:utilization_memory;type:double;not null;default:0"`
	Temperature       float64        `gorm:"column:temperature;type:double;not null;default:0"`
	PowerWatts        float64        `gorm:"column:power_watts;type:double;not null;default:0"`
	LabelsJSON        datatypes.JSON `gorm:"column:labels_json;type:json"`
	AnnotationsJSON   datatypes.JSON `gorm:"column:annotations_json;type:json"`
	Allocated         bool           `gorm:"column:allocated;type:tinyint(1);not null;default:0;index:idx_gpu_devices_snapshot_allocated,priority:2"`
	Reserved          bool           `gorm:"column:reserved;type:tinyint(1);not null;default:0"`
	CreatedAt         time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Snapshot NodeSnapshot `gorm:"foreignKey:SnapshotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GPUDevice) TableName() string { return "gpu_devices" }

//
// 7. gpu_mig_devices
//

type GPUMIGDevice struct {
	ID            uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	SnapshotID    uint64    `gorm:"column:snapshot_id;type:bigint unsigned;not null;index:idx_gpu_mig_snapshot"`
	NodeName      string    `gorm:"column:node_name;type:varchar(128);not null;index:idx_gpu_mig_node"`
	ParentGPUUUID string    `gorm:"column:parent_gpu_uuid;type:varchar(128);not null;index:idx_gpu_mig_parent"`
	MIGUUID       string    `gorm:"column:mig_uuid;type:varchar(128);not null;index:idx_gpu_mig_uuid"`
	Profile       string    `gorm:"column:profile;type:varchar(64);not null"`
	MemoryMiB     int64     `gorm:"column:memory_mib;type:bigint;not null;default:0"`
	Healthy       bool      `gorm:"column:healthy;type:tinyint(1);not null;default:1"`
	Allocated     bool      `gorm:"column:allocated;type:tinyint(1);not null;default:0"`
	Reserved      bool      `gorm:"column:reserved;type:tinyint(1);not null;default:0"`
	CreatedAt     time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Snapshot NodeSnapshot `gorm:"foreignKey:SnapshotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GPUMIGDevice) TableName() string { return "gpu_mig_devices" }

//
// 8. pod_gpu_bindings_runtime
//

type PodGPUBindingRuntime struct {
	ID         uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	SnapshotID uint64         `gorm:"column:snapshot_id;type:bigint unsigned;not null;index:idx_runtime_binding_snapshot"`
	NodeName   string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_runtime_binding_node"`
	Namespace  string         `gorm:"column:namespace;type:varchar(128);not null;index:idx_runtime_binding_pod,priority:1"`
	PodName    string         `gorm:"column:pod_name;type:varchar(128);not null;index:idx_runtime_binding_pod,priority:2"`
	GPUIDsJSON datatypes.JSON `gorm:"column:gpu_ids_json;type:json;not null"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Snapshot NodeSnapshot `gorm:"foreignKey:SnapshotID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (PodGPUBindingRuntime) TableName() string { return "pod_gpu_bindings_runtime" }

//
// 9. node_heartbeats
//

type NodeHeartbeat struct {
	NodeName   string    `gorm:"column:node_name;type:varchar(128);primaryKey"`
	Status     string    `gorm:"column:status;type:varchar(32);not null;index:idx_node_heartbeats_status"`
	Message    string    `gorm:"column:message;type:varchar(255);not null;default:''"`
	LastSeenAt time.Time `gorm:"column:last_seen_at;type:timestamp;not null;index:idx_node_heartbeats_last_seen"`
	UpdatedAt  time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (NodeHeartbeat) TableName() string { return "node_heartbeats" }

//
// 10. reservations
//

type Reservation struct {
	ID         string         `gorm:"column:id;type:varchar(64);primaryKey"`
	JobID      string         `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_reservations_job;index:idx_reservations_job"`
	NodeName   string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_reservations_node"`
	GPUIDsJSON datatypes.JSON `gorm:"column:gpu_ids_json;type:json;not null"`
	ExpireAt   time.Time      `gorm:"column:expire_at;type:timestamp;not null;index:idx_reservations_expire;index:idx_reservations_expire_node,priority:1"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Job GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Reservation) TableName() string { return "reservations" }

//
// 11. allocations
//

type Allocation struct {
	ID            string         `gorm:"column:id;type:varchar(64);primaryKey"`
	ReservationID string         `gorm:"column:reservation_id;type:varchar(64);not null;default:'';index:idx_allocations_reservation"`
	JobID         string         `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_allocations_job"`
	TenantID      string         `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_allocations_tenant"`
	NodeName      string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_allocations_node;index:idx_allocations_node_status,priority:1"`
	GPUIDsJSON    datatypes.JSON `gorm:"column:gpu_ids_json;type:json;not null"`
	Status        string         `gorm:"column:status;type:varchar(32);not null;index:idx_allocations_status;index:idx_allocations_node_status,priority:2"`
	Message       *string        `gorm:"column:message;type:text"`
	CommittedAt   *time.Time     `gorm:"column:committed_at;type:timestamp"`
	ReleasedAt    *time.Time     `gorm:"column:released_at;type:timestamp"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Job    GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Allocation) TableName() string { return "allocations" }

//
// 12. bindings
//

type Binding struct {
	ID           string         `gorm:"column:id;type:varchar(64);primaryKey"`
	AllocationID string         `gorm:"column:allocation_id;type:varchar(64);not null;default:'';index:idx_bindings_allocation"`
	JobID        string         `gorm:"column:job_id;type:varchar(64);not null;uniqueIndex:uk_bindings_job"`
	NodeName     string         `gorm:"column:node_name;type:varchar(128);not null;index:idx_bindings_node"`
	GPUIDsJSON   datatypes.JSON `gorm:"column:gpu_ids_json;type:json;not null"`
	PodName      string         `gorm:"column:pod_name;type:varchar(128);not null;default:''"`
	Namespace    string         `gorm:"column:namespace;type:varchar(128);not null;default:'';index:idx_bindings_ns_pod,priority:1"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Job GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Binding) TableName() string { return "bindings" }

//
// 13. gpu_quotas
//

type GPUQuota struct {
	ID              string    `gorm:"column:id;type:varchar(64);primaryKey"`
	TenantID        string    `gorm:"column:tenant_id;type:varchar(64);not null;uniqueIndex:uk_gpu_quotas_tenant_namespace,priority:1;index:idx_gpu_quotas_tenant"`
	Namespace       string    `gorm:"column:namespace;type:varchar(128);not null;default:'';uniqueIndex:uk_gpu_quotas_tenant_namespace,priority:2;index:idx_gpu_quotas_namespace"`
	MaxGPUCount     int       `gorm:"column:max_gpu_count;type:int;not null;default:0"`
	MaxRunningJobs  int       `gorm:"column:max_running_jobs;type:int;not null;default:0"`
	MaxQueuedJobs   int       `gorm:"column:max_queued_jobs;type:int;not null;default:0"`
	MaxGPUMemoryMiB int64     `gorm:"column:max_gpu_memory_mib;type:bigint;not null;default:0"`
	Enabled         bool      `gorm:"column:enabled;type:tinyint(1);not null;default:1"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GPUQuota) TableName() string { return "gpu_quotas" }

//
// 14. scheduling_attempts
//

type SchedulingAttempt struct {
	ID                 uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	JobID              string         `gorm:"column:job_id;type:varchar(64);not null;index:idx_sched_attempts_job;index:idx_sched_attempts_job_attempt,priority:1"`
	TenantID           string         `gorm:"column:tenant_id;type:varchar(64);not null;index:idx_sched_attempts_tenant"`
	AttemptNo          int            `gorm:"column:attempt_no;type:int;not null;index:idx_sched_attempts_job_attempt,priority:2"`
	Phase              string         `gorm:"column:phase;type:varchar(32);not null"`
	SelectedNode       string         `gorm:"column:selected_node;type:varchar(128);not null;default:''"`
	CandidateNodesJSON datatypes.JSON `gorm:"column:candidate_nodes_json;type:json"`
	ScoresJSON         datatypes.JSON `gorm:"column:scores_json;type:json"`
	FilterReasonsJSON  datatypes.JSON `gorm:"column:filter_reasons_json;type:json"`
	Result             string         `gorm:"column:result;type:varchar(32);not null;index:idx_sched_attempts_result"`
	Message            *string        `gorm:"column:message;type:text"`
	CreatedAt          time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_sched_attempts_created"`

	Job    GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SchedulingAttempt) TableName() string { return "scheduling_attempts" }

//
// 15. job_retries
//

type JobRetry struct {
	ID        uint64    `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	JobID     string    `gorm:"column:job_id;type:varchar(64);not null;index:idx_job_retries_job;index:idx_job_retries_job_retry_no,priority:1"`
	RetryNo   int       `gorm:"column:retry_no;type:int;not null;index:idx_job_retries_job_retry_no,priority:2"`
	Reason    string    `gorm:"column:reason;type:varchar(64);not null"`
	Message   *string   `gorm:"column:message;type:text"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_job_retries_created"`

	Job GPUJob `gorm:"foreignKey:JobID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (JobRetry) TableName() string { return "job_retries" }

//
// 16. audit_logs
//

type AuditLog struct {
	ID           uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	TenantID     string         `gorm:"column:tenant_id;type:varchar(64);not null;default:'';index:idx_audit_tenant"`
	Actor        string         `gorm:"column:actor;type:varchar(128);not null;index:idx_audit_actor"`
	Action       string         `gorm:"column:action;type:varchar(128);not null;index:idx_audit_action"`
	ResourceType string         `gorm:"column:resource_type;type:varchar(64);not null;index:idx_audit_resource,priority:1"`
	ResourceID   string         `gorm:"column:resource_id;type:varchar(128);not null;index:idx_audit_resource,priority:2"`
	ResourceName string         `gorm:"column:resource_name;type:varchar(128);not null;default:''"`
	Status       string         `gorm:"column:status;type:varchar(32);not null;default:''"`
	RequestID    string         `gorm:"column:request_id;type:varchar(128);not null;default:''"`
	DetailJSON   datatypes.JSON `gorm:"column:detail_json;type:json"`
	CreatedAt    time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;index:idx_audit_created"`
}

func (AuditLog) TableName() string { return "audit_logs" }

//
// 17. outbox
//

type Outbox struct {
	ID          uint64         `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement"`
	Topic       string         `gorm:"column:topic;type:varchar(128);not null;index:idx_outbox_topic"`
	EventKey    string         `gorm:"column:event_key;type:varchar(128);not null;default:''"`
	PayloadJSON datatypes.JSON `gorm:"column:payload_json;type:json;not null"`
	Status      string         `gorm:"column:status;type:varchar(32);not null;index:idx_outbox_status;index:idx_outbox_status_available,priority:1"`
	RetryCount  int            `gorm:"column:retry_count;type:int;not null;default:0"`
	LastError   *string        `gorm:"column:last_error;type:text"`
	AvailableAt time.Time      `gorm:"column:available_at;type:timestamp;not null;index:idx_outbox_available;index:idx_outbox_status_available,priority:2"`
	ProcessedAt *time.Time     `gorm:"column:processed_at;type:timestamp"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
}

func (Outbox) TableName() string { return "outbox" }

type GPUPolicy struct {
	ID                     string         `gorm:"column:id;type:varchar(64);primaryKey"`
	TenantID               string         `gorm:"column:tenant_id;type:varchar(64);not null;uniqueIndex:uk_gpu_policies_tenant_name,priority:1;index:idx_gpu_policies_tenant"`
	Name                   string         `gorm:"column:name;type:varchar(128);not null;uniqueIndex:uk_gpu_policies_tenant_name,priority:2;index:idx_gpu_policies_name"`
	Queue                  string         `gorm:"column:queue;type:varchar(64);not null;default:'';index:idx_gpu_policies_queue"`
	Priority               int            `gorm:"column:priority;type:int;not null;default:0"`
	Enabled                bool           `gorm:"column:enabled;type:tinyint(1);not null;default:1;index:idx_gpu_policies_enabled"`
	Preemptible            bool           `gorm:"column:preemptible;type:tinyint(1);not null;default:0"`
	RequireHealthy         bool           `gorm:"column:require_healthy;type:tinyint(1);not null;default:1"`
	RequireMIG             bool           `gorm:"column:require_mig;type:tinyint(1);not null;default:0"`
	MaxGPUCount            int            `gorm:"column:max_gpu_count;type:int;not null;default:0"`
	RequiredGPUModel       string         `gorm:"column:required_gpu_model;type:varchar(128);not null;default:''"`
	RequiredNodeLabelsJSON datatypes.JSON `gorm:"column:required_node_labels_json;type:json"`
	SelectorJSON           datatypes.JSON `gorm:"column:selector_json;type:json"`
	Description            *string        `gorm:"column:description;type:text"`
	CreatedAt              time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt              time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP"`

	Tenant Tenant `gorm:"foreignKey:TenantID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (GPUPolicy) TableName() string { return "gpu_policies" }
