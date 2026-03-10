package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/repo"
	"gpu-scheduler-platform/internal/util"

	"github.com/google/uuid"
)

type JobService struct {
	jobs      repo.JobRepository
	events    repo.JobEventRepository
	tenants   repo.TenantRepository
	txManager repo.TxManager
}

func NewJobService(
	jobs repo.JobRepository,
	events repo.JobEventRepository,
	tenants repo.TenantRepository,
	txManager repo.TxManager,
) *JobService {
	return &JobService{
		jobs:      jobs,
		events:    events,
		tenants:   tenants,
		txManager: txManager,
	}
}

type CreateJobRequest struct {
	TenantID    string
	Namespace   string
	Name        string
	Queue       job.QueueName
	Priority    job.PriorityClass
	Requirement job.Requirement
	Labels      map[string]string
	Annotations map[string]string
}

type ListJobsRequest struct {
	TenantID  string
	Namespace string
	Queue     job.QueueName
	Status    job.Status
	Priority  job.PriorityClass
	Keyword   string
	Limit     int
	Offset    int
}

func (s *JobService) Create(ctx context.Context, req CreateJobRequest) (*job.Job, error) {
	if s == nil || s.jobs == nil || s.events == nil || s.txManager == nil {
		return nil, util.ErrUnavailable
	}

	req.TenantID = strings.TrimSpace(req.TenantID)
	req.Namespace = strings.TrimSpace(req.Namespace)
	req.Name = strings.TrimSpace(req.Name)

	if req.TenantID == "" || req.Namespace == "" || req.Name == "" {
		return nil, util.ErrInvalidArgument
	}
	if !req.Requirement.Valid() {
		return nil, util.ErrInvalidArgument
	}
	if req.Queue == "" {
		req.Queue = job.QueueDefault
	}
	if req.Priority == "" {
		req.Priority = job.PriorityNormal
	}

	if s.tenants != nil {
		ok, err := s.tenants.Exists(ctx, req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("check tenant exists: %w", err)
		}
		if !ok {
			return nil, util.ErrNotFound
		}
	}

	now := time.Now().UTC()
	j := &job.Job{
		ID:          uuid.NewString(),
		TenantID:    req.TenantID,
		Namespace:   req.Namespace,
		Name:        req.Name,
		Queue:       req.Queue,
		Priority:    req.Priority,
		Status:      job.StatusPending,
		Requirement: req.Requirement,
		Labels:      cloneMap(req.Labels),
		Annotations: cloneMap(req.Annotations),
		RetryCount:  0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := s.txManager.WithinTx(ctx, func(txCtx context.Context) error {
		if err := s.jobs.Create(txCtx, j); err != nil {
			return err
		}
		e := &event.Event{
			ID:         uuid.NewString(),
			JobID:      j.ID,
			TenantID:   j.TenantID,
			Reason:     event.ReasonCreated,
			Message:    "job created",
			Source:     "api-server",
			OccurredAt: now,
		}
		return s.events.Create(txCtx, e)
	})
	if err != nil {
		return nil, err
	}

	return j, nil
}

func (s *JobService) GetByID(ctx context.Context, jobID string) (*job.Job, error) {
	if s == nil || s.jobs == nil {
		return nil, util.ErrUnavailable
	}
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return nil, util.ErrInvalidArgument
	}
	return s.jobs.GetByID(ctx, jobID)
}

func (s *JobService) List(ctx context.Context, req ListJobsRequest) ([]job.Job, int64, error) {
	if s == nil || s.jobs == nil {
		return nil, 0, util.ErrUnavailable
	}
	return s.jobs.List(ctx, repo.JobListFilter{
		TenantID:  strings.TrimSpace(req.TenantID),
		Namespace: strings.TrimSpace(req.Namespace),
		Queue:     req.Queue,
		Status:    req.Status,
		Priority:  req.Priority,
		Keyword:   strings.TrimSpace(req.Keyword),
		Limit:     req.Limit,
		Offset:    req.Offset,
	})
}

func (s *JobService) ListEvents(ctx context.Context, jobID string, limit, offset int) ([]event.Event, int64, error) {
	if s == nil || s.events == nil {
		return nil, 0, util.ErrUnavailable
	}
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return nil, 0, util.ErrInvalidArgument
	}
	return s.events.ListByJobID(ctx, jobID, limit, offset)
}

func cloneMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
