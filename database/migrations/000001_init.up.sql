SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- =========================================================
-- 0. 租户
-- =========================================================
CREATE TABLE IF NOT EXISTS tenants (
id            VARCHAR(64) PRIMARY KEY,
name          VARCHAR(128) NOT NULL,
enabled       TINYINT(1) NOT NULL DEFAULT 1,
description   TEXT NULL,
created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

KEY idx_tenants_enabled (enabled),
KEY idx_tenants_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 1. 队列定义
-- 可做权重、公平性、启停控制
-- tenant_id = '' 表示平台级公共队列
-- =========================================================
CREATE TABLE IF NOT EXISTS queues (
id             VARCHAR(64) PRIMARY KEY,
name           VARCHAR(64) NOT NULL,
tenant_id      VARCHAR(64) NOT NULL DEFAULT '',
weight         INT NOT NULL DEFAULT 1,
priority       INT NOT NULL DEFAULT 0,
enabled        TINYINT(1) NOT NULL DEFAULT 1,
description    TEXT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_queues_name_tenant (name, tenant_id),
KEY idx_queues_tenant (tenant_id),
KEY idx_queues_enabled (enabled),
KEY idx_queues_priority (priority)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 2. 节点当前状态表（当前视图）
-- 调度器优先查它，不必每次都回放历史快照
-- =========================================================
CREATE TABLE IF NOT EXISTS nodes (
node_name            VARCHAR(128) PRIMARY KEY,
cluster_name         VARCHAR(128) NOT NULL DEFAULT '',
source               VARCHAR(32) NOT NULL DEFAULT 'agent',
state                VARCHAR(32) NOT NULL,
schedulable          TINYINT(1) NOT NULL DEFAULT 1,

gpu_count            INT NOT NULL DEFAULT 0,
healthy_gpu_count    INT NOT NULL DEFAULT 0,
total_memory_mib     BIGINT NOT NULL DEFAULT 0,
free_memory_mib      BIGINT NOT NULL DEFAULT 0,

labels_json          JSON NULL,
annotations_json     JSON NULL,
topology_json        JSON NULL,

last_report_time     TIMESTAMP NULL DEFAULT NULL,
last_heartbeat_time  TIMESTAMP NULL DEFAULT NULL,

created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

KEY idx_nodes_state (state),
KEY idx_nodes_schedulable (schedulable),
KEY idx_nodes_cluster (cluster_name),
KEY idx_nodes_last_report_time (last_report_time),
KEY idx_nodes_last_heartbeat_time (last_heartbeat_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 3. GPU 作业
-- 建议 status 统一使用：
-- PENDING / QUEUED / SCHEDULING / RESERVED / ALLOCATED /
-- BINDING / RUNNING / SUCCEEDED / FAILED / CANCELLED / TIMEOUT
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_jobs (
id                           VARCHAR(64) PRIMARY KEY,
tenant_id                    VARCHAR(64) NOT NULL,

namespace                    VARCHAR(128) NOT NULL,
name                         VARCHAR(128) NOT NULL,
queue                        VARCHAR(64) NOT NULL,
priority                     VARCHAR(32) NOT NULL,
status                       VARCHAR(32) NOT NULL,

submitter                    VARCHAR(128) NOT NULL DEFAULT '',
scheduler_name               VARCHAR(64) NOT NULL DEFAULT 'default',

gpu_count                    INT NOT NULL,
gpu_memory_mib               BIGINT NOT NULL,
gpu_model                    VARCHAR(128) NOT NULL DEFAULT '',
require_same_node            TINYINT(1) NOT NULL DEFAULT 0,
require_healthy              TINYINT(1) NOT NULL DEFAULT 1,
require_mig                  TINYINT(1) NOT NULL DEFAULT 0,
mig_profile                  VARCHAR(64) NOT NULL DEFAULT '',
require_nvlink               TINYINT(1) NOT NULL DEFAULT 0,

preemptible                  TINYINT(1) NOT NULL DEFAULT 0,
retryable                    TINYINT(1) NOT NULL DEFAULT 1,
max_retry                    INT NOT NULL DEFAULT 0,
expected_duration_sec        BIGINT NOT NULL DEFAULT 0,

run_policy_json              JSON NULL,
preferred_node_labels_json   JSON NULL,
preferred_gpu_labels_json    JSON NULL,
labels_json                  JSON NULL,
annotations_json             JSON NULL,

retry_count                  INT NOT NULL DEFAULT 0,
message                      TEXT NULL,

scheduled_at                 TIMESTAMP NULL DEFAULT NULL,
started_at                   TIMESTAMP NULL DEFAULT NULL,
finished_at                  TIMESTAMP NULL DEFAULT NULL,

created_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at                   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_gpu_jobs_ns_name (namespace, name),
KEY idx_jobs_tenant (tenant_id),
KEY idx_jobs_namespace (namespace),
KEY idx_jobs_queue (queue),
KEY idx_jobs_status (status),
KEY idx_jobs_created (created_at),
KEY idx_jobs_queue_status_created (queue, status, created_at),
KEY idx_jobs_tenant_status_created (tenant_id, status, created_at),
KEY idx_jobs_scheduler (scheduler_name),

CONSTRAINT fk_gpu_jobs_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 4. GPU 作业事件
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_job_events (
id             VARCHAR(64) PRIMARY KEY,
job_id         VARCHAR(64) NOT NULL,
tenant_id      VARCHAR(64) NOT NULL,
reason         VARCHAR(64) NOT NULL,
message        TEXT NULL,
source         VARCHAR(64) NOT NULL,
occurred_at    TIMESTAMP NOT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_job_events_job (job_id),
KEY idx_job_events_tenant (tenant_id),
KEY idx_job_events_occurred (occurred_at),
KEY idx_job_events_reason (reason),

CONSTRAINT fk_job_events_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE,
CONSTRAINT fk_job_events_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 5. 节点快照头表（历史）
-- 每次 agent report 到达，可插入一条
-- =========================================================
CREATE TABLE IF NOT EXISTS node_snapshots (
id                 BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
version            VARCHAR(64) NOT NULL,
agent_version      VARCHAR(64) NOT NULL DEFAULT '',
cluster_name       VARCHAR(128) NOT NULL DEFAULT '',
node_name          VARCHAR(128) NOT NULL,
source             VARCHAR(32) NOT NULL DEFAULT 'agent',
node_state         VARCHAR(32) NOT NULL,
schedulable        TINYINT(1) NOT NULL DEFAULT 1,

gpu_count          INT NOT NULL DEFAULT 0,
healthy_gpu_count  INT NOT NULL DEFAULT 0,
total_memory_mib   BIGINT NOT NULL DEFAULT 0,
free_memory_mib    BIGINT NOT NULL DEFAULT 0,

labels_json        JSON NULL,
annotations_json   JSON NULL,
topology_json      JSON NULL,

report_time        TIMESTAMP NOT NULL,
created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_node_snapshots_version (version),
KEY idx_node_snapshots_node (node_name),
KEY idx_node_snapshots_report_time (report_time),
KEY idx_node_snapshots_created (created_at),
KEY idx_node_snapshots_node_report (node_name, report_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 6. GPU 设备快照子表
-- 注意：这是快照表，不是全局唯一设备主表
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_devices (
id                   BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
snapshot_id          BIGINT UNSIGNED NOT NULL,
node_name            VARCHAR(128) NOT NULL,
uuid                 VARCHAR(128) NOT NULL,
gpu_index            INT NOT NULL,
model                VARCHAR(128) NOT NULL,
vendor               VARCHAR(64) NOT NULL DEFAULT 'nvidia',
type                 VARCHAR(16) NOT NULL,
memory_mib           BIGINT NOT NULL,
free_memory_mib      BIGINT NOT NULL,
healthy              TINYINT(1) NOT NULL DEFAULT 1,
health               VARCHAR(32) NOT NULL,
mig_enabled          TINYINT(1) NOT NULL DEFAULT 0,
mig_profile          VARCHAR(64) NOT NULL DEFAULT '',
utilization_gpu      DOUBLE NOT NULL DEFAULT 0,
utilization_memory   DOUBLE NOT NULL DEFAULT 0,
temperature          DOUBLE NOT NULL DEFAULT 0,
power_watts          DOUBLE NOT NULL DEFAULT 0,
labels_json          JSON NULL,
annotations_json     JSON NULL,
allocated            TINYINT(1) NOT NULL DEFAULT 0,
reserved             TINYINT(1) NOT NULL DEFAULT 0,
created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_gpu_devices_snapshot (snapshot_id),
KEY idx_gpu_devices_node (node_name),
KEY idx_gpu_devices_uuid (uuid),
KEY idx_gpu_devices_node_uuid (node_name, uuid),
KEY idx_gpu_devices_snapshot_node (snapshot_id, node_name),
KEY idx_gpu_devices_snapshot_health (snapshot_id, healthy),
KEY idx_gpu_devices_snapshot_allocated (snapshot_id, allocated),

CONSTRAINT fk_gpu_devices_snapshot
FOREIGN KEY (snapshot_id) REFERENCES node_snapshots(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 7. MIG 实例快照子表
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_mig_devices (
id                BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
snapshot_id       BIGINT UNSIGNED NOT NULL,
node_name         VARCHAR(128) NOT NULL,
parent_gpu_uuid   VARCHAR(128) NOT NULL,
mig_uuid          VARCHAR(128) NOT NULL,
profile           VARCHAR(64) NOT NULL,
memory_mib        BIGINT NOT NULL DEFAULT 0,
healthy           TINYINT(1) NOT NULL DEFAULT 1,
allocated         TINYINT(1) NOT NULL DEFAULT 0,
reserved          TINYINT(1) NOT NULL DEFAULT 0,
created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_gpu_mig_snapshot (snapshot_id),
KEY idx_gpu_mig_node (node_name),
KEY idx_gpu_mig_parent (parent_gpu_uuid),
KEY idx_gpu_mig_uuid (mig_uuid),

CONSTRAINT fk_gpu_mig_devices_snapshot
FOREIGN KEY (snapshot_id) REFERENCES node_snapshots(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 8. 运行态 Pod-GPU 绑定快照
-- 这是 agent 看到的运行时状态，不是平台承诺绑定
-- =========================================================
CREATE TABLE IF NOT EXISTS pod_gpu_bindings_runtime (
id                BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
snapshot_id       BIGINT UNSIGNED NOT NULL,
node_name         VARCHAR(128) NOT NULL,
namespace         VARCHAR(128) NOT NULL,
pod_name          VARCHAR(128) NOT NULL,
gpu_ids_json      JSON NOT NULL,
created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_runtime_binding_snapshot (snapshot_id),
KEY idx_runtime_binding_node (node_name),
KEY idx_runtime_binding_pod (namespace, pod_name),

CONSTRAINT fk_runtime_bindings_snapshot
FOREIGN KEY (snapshot_id) REFERENCES node_snapshots(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 9. 节点心跳最新状态表
-- 每次 heartbeat 到达时 upsert 即可
-- =========================================================
CREATE TABLE IF NOT EXISTS node_heartbeats (
node_name         VARCHAR(128) PRIMARY KEY,
status            VARCHAR(32) NOT NULL,
message           VARCHAR(255) NOT NULL DEFAULT '',
last_seen_at      TIMESTAMP NOT NULL,
updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

KEY idx_node_heartbeats_last_seen (last_seen_at),
KEY idx_node_heartbeats_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 10. 预留记录
-- 调度器先 reserve，超时未提交则过期回收
-- =========================================================
CREATE TABLE IF NOT EXISTS reservations (
id             VARCHAR(64) PRIMARY KEY,
job_id         VARCHAR(64) NOT NULL,
node_name      VARCHAR(128) NOT NULL,
gpu_ids_json   JSON NOT NULL,
expire_at      TIMESTAMP NOT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

UNIQUE KEY uk_reservations_job (job_id),
KEY idx_reservations_node (node_name),
KEY idx_reservations_expire (expire_at),
KEY idx_reservations_expire_node (expire_at, node_name),

CONSTRAINT fk_reservations_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 11. 调度分配结果
-- reservation -> allocation -> binding
-- =========================================================
CREATE TABLE IF NOT EXISTS allocations (
id               VARCHAR(64) PRIMARY KEY,
reservation_id   VARCHAR(64) NOT NULL DEFAULT '',
job_id           VARCHAR(64) NOT NULL,
tenant_id        VARCHAR(64) NOT NULL,
node_name        VARCHAR(128) NOT NULL,
gpu_ids_json     JSON NOT NULL,
status           VARCHAR(32) NOT NULL,
message          TEXT NULL,
committed_at     TIMESTAMP NULL DEFAULT NULL,
released_at      TIMESTAMP NULL DEFAULT NULL,
created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_allocations_job (job_id),
KEY idx_allocations_tenant (tenant_id),
KEY idx_allocations_node (node_name),
KEY idx_allocations_status (status),
KEY idx_allocations_node_status (node_name, status),
KEY idx_allocations_reservation (reservation_id),

CONSTRAINT fk_allocations_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE,
CONSTRAINT fk_allocations_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 12. 最终绑定结果
-- 这是平台视角的计划/提交绑定
-- =========================================================
CREATE TABLE IF NOT EXISTS bindings (
id              VARCHAR(64) PRIMARY KEY,
allocation_id   VARCHAR(64) NOT NULL DEFAULT '',
job_id          VARCHAR(64) NOT NULL,
node_name       VARCHAR(128) NOT NULL,
gpu_ids_json    JSON NOT NULL,
pod_name        VARCHAR(128) NOT NULL DEFAULT '',
namespace       VARCHAR(128) NOT NULL DEFAULT '',
created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

UNIQUE KEY uk_bindings_job (job_id),
KEY idx_bindings_node (node_name),
KEY idx_bindings_ns_pod (namespace, pod_name),
KEY idx_bindings_allocation (allocation_id),

CONSTRAINT fk_bindings_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 13. 租户 GPU 配额
-- namespace = '' 表示租户全局默认配额
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_quotas (
id                 VARCHAR(64) PRIMARY KEY,
tenant_id          VARCHAR(64) NOT NULL,
namespace          VARCHAR(128) NOT NULL DEFAULT '',
max_gpu_count      INT NOT NULL DEFAULT 0,
max_running_jobs   INT NOT NULL DEFAULT 0,
max_queued_jobs    INT NOT NULL DEFAULT 0,
max_gpu_memory_mib BIGINT NOT NULL DEFAULT 0,
enabled            TINYINT(1) NOT NULL DEFAULT 1,
created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_gpu_quotas_tenant_namespace (tenant_id, namespace),
KEY idx_gpu_quotas_tenant (tenant_id),
KEY idx_gpu_quotas_namespace (namespace),

CONSTRAINT fk_gpu_quotas_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 14. 调度尝试记录
-- 用于解释：为什么这个 job 没调度上、候选节点有哪些、打分如何
-- =========================================================
CREATE TABLE IF NOT EXISTS scheduling_attempts (
id                  BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
job_id               VARCHAR(64) NOT NULL,
tenant_id            VARCHAR(64) NOT NULL,
attempt_no           INT NOT NULL,
phase                VARCHAR(32) NOT NULL,
selected_node        VARCHAR(128) NOT NULL DEFAULT '',
candidate_nodes_json JSON NULL,
scores_json          JSON NULL,
filter_reasons_json  JSON NULL,
result               VARCHAR(32) NOT NULL,
message              TEXT NULL,
created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_sched_attempts_job (job_id),
KEY idx_sched_attempts_tenant (tenant_id),
KEY idx_sched_attempts_result (result),
KEY idx_sched_attempts_created (created_at),
KEY idx_sched_attempts_job_attempt (job_id, attempt_no),

CONSTRAINT fk_sched_attempts_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE,
CONSTRAINT fk_sched_attempts_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 15. 重试记录
-- 记录 job 的每次 retry 原因
-- =========================================================
CREATE TABLE IF NOT EXISTS job_retries (
id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
job_id         VARCHAR(64) NOT NULL,
retry_no       INT NOT NULL,
reason         VARCHAR(64) NOT NULL,
message        TEXT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_job_retries_job (job_id),
KEY idx_job_retries_created (created_at),
KEY idx_job_retries_job_retry_no (job_id, retry_no),

CONSTRAINT fk_job_retries_job
FOREIGN KEY (job_id) REFERENCES gpu_jobs(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 16. 审计日志
-- =========================================================
CREATE TABLE IF NOT EXISTS audit_logs (
id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
tenant_id      VARCHAR(64) NOT NULL DEFAULT '',
actor          VARCHAR(128) NOT NULL,
action         VARCHAR(128) NOT NULL,
resource_type  VARCHAR(64) NOT NULL,
resource_id    VARCHAR(128) NOT NULL,
resource_name  VARCHAR(128) NOT NULL DEFAULT '',
status         VARCHAR(32) NOT NULL DEFAULT '',
request_id     VARCHAR(128) NOT NULL DEFAULT '',
detail_json    JSON NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_audit_tenant (tenant_id),
KEY idx_audit_actor (actor),
KEY idx_audit_action (action),
KEY idx_audit_created (created_at),
KEY idx_audit_resource (resource_type, resource_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 17. Outbox 事件表
-- 用于事务内写事件，异步投递 MQ / webhook / 下游系统
-- =========================================================
CREATE TABLE IF NOT EXISTS outbox (
id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
topic          VARCHAR(128) NOT NULL,
event_key      VARCHAR(128) NOT NULL DEFAULT '',
payload_json   JSON NOT NULL,
status         VARCHAR(32) NOT NULL,
retry_count    INT NOT NULL DEFAULT 0,
last_error     TEXT NULL,
available_at   TIMESTAMP NOT NULL,
processed_at   TIMESTAMP NULL DEFAULT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

KEY idx_outbox_topic (topic),
KEY idx_outbox_status (status),
KEY idx_outbox_available (available_at),
KEY idx_outbox_status_available (status, available_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- =========================================================
-- 18.
--
-- =========================================================
CREATE TABLE IF NOT EXISTS gpu_policies (
id                          VARCHAR(64) PRIMARY KEY,
tenant_id                   VARCHAR(64)  NOT NULL,
name                        VARCHAR(128) NOT NULL,
queue                       VARCHAR(64)  NOT NULL DEFAULT '',
priority                    INT          NOT NULL DEFAULT 0,
enabled                     TINYINT(1)   NOT NULL DEFAULT 1,
preemptible                 TINYINT(1)   NOT NULL DEFAULT 0,
require_healthy             TINYINT(1)   NOT NULL DEFAULT 1,
require_mig                 TINYINT(1)   NOT NULL DEFAULT 0,
max_gpu_count               INT          NOT NULL DEFAULT 0,
required_gpu_model          VARCHAR(128) NOT NULL DEFAULT '',
required_node_labels_json   JSON NULL,
selector_json               JSON NULL,
description                 TEXT NULL,
created_at                  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at                  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_gpu_policies_tenant_name (tenant_id, name),
KEY idx_gpu_policies_tenant (tenant_id),
KEY idx_gpu_policies_name (name),
KEY idx_gpu_policies_queue (queue),
KEY idx_gpu_policies_enabled (enabled),

CONSTRAINT fk_gpu_policies_tenant
FOREIGN KEY (tenant_id) REFERENCES tenants(id)
ON UPDATE CASCADE
ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;