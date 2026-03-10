SET FOREIGN_KEY_CHECKS = 0;

ALTER TABLE gpu_jobs
DROP KEY idx_jobs_queue_status_priority_created;

ALTER TABLE allocations
DROP KEY idx_allocations_job_status;

ALTER TABLE reservations
DROP KEY idx_reservations_expire_node;

SET FOREIGN_KEY_CHECKS = 1;