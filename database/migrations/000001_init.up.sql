SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS tenants (
id            VARCHAR(64) PRIMARY KEY,
name          VARCHAR(128) NOT NULL,
enabled       TINYINT(1) NOT NULL DEFAULT 1,
description   TEXT NULL,
created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gpu_jobs (
id                           VARCHAR(64) PRIMARY KEY,
tenant_id                    VARCHAR(64) NOT NULL,
namespace                    VARCHAR(128) NOT NULL,
name                         VARCHAR(128) NOT NULL,
queue                        VARCHAR(64) NOT NULL,
priority                     VARCHAR(32) NOT NULL,
status                       VARCHAR(32) NOT NULL,

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

KEY idx_jobs_tenant (tenant_id),
KEY idx_jobs_namespace (namespace),
KEY idx_jobs_queue (queue),
KEY idx_jobs_status (status),
KEY idx_jobs_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gpu_job_events (
      id            VARCHAR(64) PRIMARY KEY,
job_id        VARCHAR(64) NOT NULL,
tenant_id     VARCHAR(64) NOT NULL,
reason        VARCHAR(64) NOT NULL,
message       TEXT NULL,
source        VARCHAR(64) NOT NULL,
occurred_at   TIMESTAMP NOT NULL,
created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_job_events_job (job_id),
KEY idx_job_events_tenant (tenant_id),
KEY idx_job_events_occurred (occurred_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS node_snapshots (
      id              BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
      version         VARCHAR(64) NOT NULL,
node_name       VARCHAR(128) NOT NULL,
node_state      VARCHAR(32) NOT NULL,
schedulable     TINYINT(1) NOT NULL DEFAULT 1,
labels_json     JSON NULL,
annotations_json JSON NULL,
topology_json   JSON NULL,
created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_node_snapshots_version (version),
KEY idx_node_snapshots_node (node_name),
KEY idx_node_snapshots_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS gpu_devices (
   id              VARCHAR(64) PRIMARY KEY,
snapshot_id     BIGINT UNSIGNED NOT NULL,
node_name       VARCHAR(128) NOT NULL,
uuid            VARCHAR(128) NOT NULL,
gpu_index       INT NOT NULL,
model           VARCHAR(128) NOT NULL,
vendor          VARCHAR(64) NOT NULL DEFAULT '',
type            VARCHAR(16) NOT NULL,
memory_mib      BIGINT NOT NULL,
free_memory_mib BIGINT NOT NULL,
healthy         TINYINT(1) NOT NULL DEFAULT 1,
health          VARCHAR(32) NOT NULL,
mig_enabled     TINYINT(1) NOT NULL DEFAULT 0,
mig_profile     VARCHAR(64) NOT NULL DEFAULT '',
labels_json     JSON NULL,
annotations_json JSON NULL,
allocated       TINYINT(1) NOT NULL DEFAULT 0,
reserved        TINYINT(1) NOT NULL DEFAULT 0,
created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_gpu_devices_snapshot (snapshot_id),
KEY idx_gpu_devices_node (node_name),
KEY idx_gpu_devices_uuid (uuid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS allocations (
   id             VARCHAR(64) PRIMARY KEY,
job_id         VARCHAR(64) NOT NULL,
tenant_id      VARCHAR(64) NOT NULL,
node_name      VARCHAR(128) NOT NULL,
gpu_ids_json   JSON NOT NULL,
status         VARCHAR(32) NOT NULL,
message        TEXT NULL,
committed_at   TIMESTAMP NULL DEFAULT NULL,
released_at    TIMESTAMP NULL DEFAULT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_allocations_job (job_id),
KEY idx_allocations_tenant (tenant_id),
KEY idx_allocations_node (node_name),
KEY idx_allocations_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reservations (
    id             VARCHAR(64) PRIMARY KEY,
job_id         VARCHAR(64) NOT NULL,
node_name      VARCHAR(128) NOT NULL,
gpu_ids_json   JSON NOT NULL,
expire_at      TIMESTAMP NOT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

UNIQUE KEY uk_reservations_job (job_id),
KEY idx_reservations_node (node_name),
KEY idx_reservations_expire (expire_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS bindings (
id             VARCHAR(64) PRIMARY KEY,
job_id         VARCHAR(64) NOT NULL,
node_name      VARCHAR(128) NOT NULL,
gpu_ids_json   JSON NOT NULL,
pod_name       VARCHAR(128) NOT NULL DEFAULT '',
namespace      VARCHAR(128) NOT NULL DEFAULT '',
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

UNIQUE KEY uk_bindings_job (job_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

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

KEY idx_gpu_quotas_tenant (tenant_id),
KEY idx_gpu_quotas_namespace (namespace)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS audit_logs (
  id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
  actor          VARCHAR(128) NOT NULL,
action         VARCHAR(128) NOT NULL,
resource_type  VARCHAR(64) NOT NULL,
resource_id    VARCHAR(128) NOT NULL,
request_id     VARCHAR(128) NOT NULL DEFAULT '',
detail_json    JSON NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

KEY idx_audit_actor (actor),
KEY idx_audit_action (action),
KEY idx_audit_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS outbox (
id             BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
topic          VARCHAR(128) NOT NULL,
event_key      VARCHAR(128) NOT NULL DEFAULT '',
payload_json   JSON NOT NULL,
status         VARCHAR(32) NOT NULL,
available_at   TIMESTAMP NOT NULL,
processed_at   TIMESTAMP NULL DEFAULT NULL,
created_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at     TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

KEY idx_outbox_topic (topic),
KEY idx_outbox_status (status),
KEY idx_outbox_available (available_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;