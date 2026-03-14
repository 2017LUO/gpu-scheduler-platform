package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

var ErrQuotaExceeded = errors.New("job exceeds quota")

type JobService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

type ListJobsInput struct {
	TenantID  string
	Namespace string
	Queue     string
	Status    string
	Limit     int
	Offset    int
}

func NewJobService(repos *repoimpl.Repos, lg *zap.Logger) *JobService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &JobService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *JobService) Create(ctx context.Context, req dto.CreateJobRequest) (*model.GPUJob, error) {
	if req.TenantID == "" || req.Namespace == "" || req.Name == "" || req.Queue == "" {
		return nil, repoimpl.ErrInvalidArgument
	}

	tenant, err := s.repos.Tenants.Get(ctx, req.TenantID)
	if err != nil {
		return nil, err
	}
	if !tenant.Enabled {
		return nil, fmt.Errorf("tenant %s is disabled", req.TenantID)
	}

	queue, err := s.repos.Queues.GetByName(ctx, req.TenantID, req.Queue)
	if err != nil {
		return nil, err
	}
	if !queue.Enabled {
		return nil, fmt.Errorf("queue %s is disabled", req.Queue)
	}

	if quota, err := s.repos.GPUQuotas.GetByTenantAndNamespace(ctx, req.TenantID, req.Namespace); err == nil && quota.Enabled {
		if quota.MaxGPUCount > 0 && req.GPUCount > quota.MaxGPUCount {
			return nil, ErrQuotaExceeded
		}
		if quota.MaxQueuedJobs > 0 {
			queued, cntErr := s.repos.GPUJobs.CountByStatus(ctx, req.TenantID, "QUEUED")
			if cntErr == nil && queued >= int64(quota.MaxQueuedJobs) {
				return nil, ErrQuotaExceeded
			}
		}
	}

	now := s.nowFunc()

	job := &model.GPUJob{
		ID:                      newID("job"),
		TenantID:                req.TenantID,
		Namespace:               req.Namespace,
		Name:                    req.Name,
		Queue:                   req.Queue,
		Priority:                defaultString(req.Priority, "NORMAL"),
		Status:                  "QUEUED",
		Submitter:               defaultString(req.Submitter, "unknown"),
		GPUCount:                req.GPUCount,
		GPUMemoryMiB:            req.GPUMemoryMiB,
		GPUModel:                req.GPUModel,
		RequireSameNode:         req.RequireSameNode,
		RequireHealthy:          req.RequireHealthy,
		RequireMIG:              req.RequireMIG,
		ExpectedDurationSec:     req.ExpectedDurationSec,
		RunPolicyJSON:           mustJSON(req.RunPolicy),
		PreferredNodeLabelsJSON: mustJSON(req.PreferredNodeLabels),
		PreferredGPULabelsJSON:  mustJSON(req.PreferredGPULabels),
		LabelsJSON:              mustJSON(req.Labels),
		AnnotationsJSON:         mustJSON(req.Annotations),
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	if err := s.repos.GPUJobs.Create(ctx, job); err != nil {
		return nil, err
	}

	createdMsg := "job created and queued"
	event := &model.GPUJobEvent{
		ID:         newID("evt"),
		JobID:      job.ID,
		TenantID:   job.TenantID,
		Reason:     "JOB_CREATED",
		Message:    &createdMsg,
		Source:     "apiserver",
		OccurredAt: now,
		CreatedAt:  now,
	}
	if err := s.repos.GPUJobEvents.Create(ctx, event); err != nil {
		return nil, err
	}

	_ = s.repos.AuditLogs.Create(ctx, &model.AuditLog{
		TenantID:     job.TenantID,
		Actor:        job.Submitter,
		Action:       "job.create",
		ResourceType: "gpu_job",
		ResourceID:   job.ID,
		ResourceName: job.Name,
		Status:       "SUCCESS",
		CreatedAt:    now,
		DetailJSON:   mustJSON(map[string]any{"namespace": job.Namespace, "queue": job.Queue}),
	})

	_ = s.repos.Outbox.Create(ctx, &model.Outbox{
		Topic:       "job.created",
		EventKey:    job.ID,
		PayloadJSON: mustJSON(map[string]any{"job_id": job.ID, "tenant_id": job.TenantID, "namespace": job.Namespace, "name": job.Name}),
		Status:      "PENDING",
		AvailableAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
	})

	return job, nil
}

func (s *JobService) Get(ctx context.Context, id string) (*model.GPUJob, error) {
	return s.repos.GPUJobs.Get(ctx, id)
}

func (s *JobService) List(ctx context.Context, in ListJobsInput) ([]model.GPUJob, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{
		Limit:  in.Limit,
		Offset: in.Offset,
	}.Normalize(50, 500)

	items, err := s.repos.GPUJobs.List(ctx, repoimpl.GPUJobFilter{
		TenantID:  in.TenantID,
		Namespace: in.Namespace,
		Queue:     in.Queue,
		Status:    in.Status,
	}, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func (s *JobService) ListEvents(ctx context.Context, jobID string, limit, offset int) ([]model.GPUJobEvent, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{
		Limit:  limit,
		Offset: offset,
	}.Normalize(100, 1000)

	items, err := s.repos.GPUJobEvents.ListByJob(ctx, jobID, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func defaultString(v, def string) string {
	if v == "" {
		return def
	}
	return v
}

func mustJSON(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}

func newID(prefix string) string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s-%x-%d", prefix, b[:], time.Now().UTC().UnixNano())
}
