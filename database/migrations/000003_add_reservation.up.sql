SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE reservations
ADD KEY idx_reservations_expire_node (expire_at, node_name);

ALTER TABLE allocations
ADD KEY idx_allocations_job_status (job_id, status);

ALTER TABLE gpu_jobs
ADD KEY idx_jobs_queue_status_priority_created (queue, status, priority, created_at);

SET FOREIGN_KEY_CHECKS = 1;