package mysql

import (
	"context"
	"fmt"
	"strings"

	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/repo"
	"gpu-scheduler-platform/internal/repo/models"
	"gpu-scheduler-platform/internal/util"

	"gorm.io/gorm"
)

type GPUJobRepo struct {
	db *gorm.DB
}

func NewGPUJobRepo(db *gorm.DB) *GPUJobRepo {
	return &GPUJobRepo{db: db}
}

func (r *GPUJobRepo) Create(ctx context.Context, j *job.Job) error {
	if r == nil || r.db == nil || j == nil {
		return util.ErrInvalidArgument
	}
	m := toJobModel(j)
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Create(m).Error; err != nil {
		return fmt.Errorf("create gpu job: %w", err)
	}
	return nil
}

func (r *GPUJobRepo) Update(ctx context.Context, j *job.Job) error {
	if r == nil || r.db == nil || j == nil || strings.TrimSpace(j.ID) == "" {
		return util.ErrInvalidArgument
	}

	m := toJobModel(j)
	err := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.GPUJob{}).
		Where("id = ?", j.ID).
		Updates(map[string]any{
			"tenant_id":                  m.TenantID,
			"namespace":                  m.Namespace,
			"name":                       m.Name,
			"queue":                      m.Queue,
			"priority":                   m.Priority,
			"status":                     m.Status,
			"gpu_count":                  m.GPUCount,
			"gpu_memory_mib":             m.GPUMemoryMiB,
			"gpu_model":                  m.GPUModel,
			"require_same_node":          m.RequireSameNode,
			"require_healthy":            m.RequireHealthy,
			"require_mig":                m.RequireMIG,
			"mig_profile":                m.MIGProfile,
			"require_nvlink":             m.RequireNVLink,
			"preemptible":                m.Preemptible,
			"retryable":                  m.Retryable,
			"max_retry":                  m.MaxRetry,
			"expected_duration_sec":      m.ExpectedDuration,
			"preferred_node_labels_json": m.PreferredNodeLabelsJSON,
			"preferred_gpu_labels_json":  m.PreferredGPULabelsJSON,
			"labels_json":                m.LabelsJSON,
			"annotations_json":           m.AnnotationsJSON,
			"retry_count":                m.RetryCount,
			"message":                    m.Message,
			"scheduled_at":               m.ScheduledAt,
			"started_at":                 m.StartedAt,
			"finished_at":                m.FinishedAt,
		}).Error
	if err != nil {
		return fmt.Errorf("update gpu job: %w", err)
	}
	return nil
}

func (r *GPUJobRepo) UpdateStatus(ctx context.Context, jobID string, status job.Status, message string) error {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.GPUJob{}).
		Where("id = ?", jobID).
		Updates(map[string]any{
			"status":  string(status),
			"message": message,
		})
	if res.Error != nil {
		return fmt.Errorf("update gpu job status: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}

func (r *GPUJobRepo) GetByID(ctx context.Context, jobID string) (*job.Job, error) {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return nil, util.ErrInvalidArgument
	}
	var m models.GPUJob
	if err := dbFromContext(ctx, r.db).WithContext(ctx).Where("id = ?", jobID).Take(&m).Error; err != nil {
		return nil, wrapNotFound(err, "get gpu job by id")
	}
	return toJobDomain(&m), nil
}

func (r *GPUJobRepo) List(ctx context.Context, filter repo.JobListFilter) ([]job.Job, int64, error) {
	if r == nil || r.db == nil {
		return nil, 0, util.ErrInvalidArgument
	}

	dbq := dbFromContext(ctx, r.db).WithContext(ctx).Model(&models.GPUJob{})

	if filter.TenantID != "" {
		dbq = dbq.Where("tenant_id = ?", filter.TenantID)
	}
	if filter.Namespace != "" {
		dbq = dbq.Where("namespace = ?", filter.Namespace)
	}
	if filter.Queue != "" {
		dbq = dbq.Where("queue = ?", string(filter.Queue))
	}
	if filter.Status != "" {
		dbq = dbq.Where("status = ?", string(filter.Status))
	}
	if filter.Priority != "" {
		dbq = dbq.Where("priority = ?", string(filter.Priority))
	}
	if filter.Keyword != "" {
		kw := "%" + filter.Keyword + "%"
		dbq = dbq.Where("name LIKE ? OR message LIKE ?", kw, kw)
	}
	if filter.CreatedAfter != nil {
		dbq = dbq.Where("created_at >= ?", *filter.CreatedAfter)
	}
	if filter.CreatedBefore != nil {
		dbq = dbq.Where("created_at < ?", *filter.CreatedBefore)
	}

	var total int64
	if err := dbq.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count gpu jobs: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	var ms []models.GPUJob
	if err := dbq.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, fmt.Errorf("list gpu jobs: %w", err)
	}

	out := make([]job.Job, 0, len(ms))
	for i := range ms {
		out = append(out, *toJobDomain(&ms[i]))
	}
	return out, total, nil
}

func (r *GPUJobRepo) ListPending(ctx context.Context, filter repo.PendingJobFilter) ([]job.Job, error) {
	if r == nil || r.db == nil {
		return nil, util.ErrInvalidArgument
	}

	dbq := dbFromContext(ctx, r.db).WithContext(ctx).
		Model(&models.GPUJob{}).
		Where("status IN ?", []string{string(job.StatusPending), string(job.StatusQueued)})

	if len(filter.Queues) > 0 {
		qs := make([]string, 0, len(filter.Queues))
		for _, q := range filter.Queues {
			if q != "" {
				qs = append(qs, string(q))
			}
		}
		if len(qs) > 0 {
			dbq = dbq.Where("queue IN ?", qs)
		}
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 500 {
		limit = 500
	}

	var ms []models.GPUJob
	if err := dbq.Order("created_at ASC").Limit(limit).Find(&ms).Error; err != nil {
		return nil, fmt.Errorf("list pending gpu jobs: %w", err)
	}

	out := make([]job.Job, 0, len(ms))
	for i := range ms {
		out = append(out, *toJobDomain(&ms[i]))
	}
	return out, nil
}

func (r *GPUJobRepo) Delete(ctx context.Context, jobID string) error {
	if r == nil || r.db == nil || strings.TrimSpace(jobID) == "" {
		return util.ErrInvalidArgument
	}
	res := dbFromContext(ctx, r.db).WithContext(ctx).Delete(&models.GPUJob{}, "id = ?", jobID)
	if res.Error != nil {
		return fmt.Errorf("delete gpu job: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return util.ErrNotFound
	}
	return nil
}
