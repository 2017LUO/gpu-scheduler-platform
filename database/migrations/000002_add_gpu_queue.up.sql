SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

CREATE TABLE IF NOT EXISTS gpu_queues (
id               VARCHAR(64) PRIMARY KEY,
name             VARCHAR(64) NOT NULL,
priority         INT NOT NULL DEFAULT 0,
weight           INT NOT NULL DEFAULT 1,
enabled          TINYINT(1) NOT NULL DEFAULT 1,
max_queued_jobs  INT NOT NULL DEFAULT 0,
description      TEXT NULL,
created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

UNIQUE KEY uk_gpu_queues_name (name),
KEY idx_gpu_queues_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS = 1;