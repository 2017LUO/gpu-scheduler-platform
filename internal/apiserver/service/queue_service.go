package service

import (
	"context"
	"errors"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type QueueService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewQueueService(repos *repoimpl.Repos, lg *zap.Logger) *QueueService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &QueueService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *QueueService) Upsert(ctx context.Context, req dto.UpsertQueueRequest) (*model.Queue, error) {
	if req.Name == "" {
		return nil, repoimpl.ErrInvalidArgument
	}
	if req.TenantID != "" {
		if _, err := s.repos.Tenants.Get(ctx, req.TenantID); err != nil {
			return nil, err
		}
	}

	now := s.nowFunc()
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	weight := req.Weight
	if weight <= 0 {
		weight = 1
	}

	existing, err := s.repos.Queues.GetByName(ctx, req.TenantID, req.Name)
	id := newID("queue")
	createdAt := now
	if err == nil {
		id = existing.ID
		createdAt = existing.CreatedAt
	} else if !errors.Is(err, repoimpl.ErrNotFound) {
		return nil, err
	}

	m := &model.Queue{
		ID:          id,
		Name:        req.Name,
		TenantID:    req.TenantID,
		Weight:      weight,
		Priority:    req.Priority,
		Enabled:     enabled,
		Description: req.Description,
		CreatedAt:   createdAt,
		UpdatedAt:   now,
	}

	if err := s.repos.Queues.Upsert(ctx, m); err != nil {
		return nil, err
	}
	return s.repos.Queues.GetByName(ctx, req.TenantID, req.Name)
}

func (s *QueueService) Get(ctx context.Context, tenantID, name string) (*model.Queue, error) {
	return s.repos.Queues.GetByName(ctx, tenantID, name)
}

func (s *QueueService) List(ctx context.Context, tenantID string, enabled *bool, limit, offset int) ([]model.Queue, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{Limit: limit, Offset: offset}.Normalize(50, 500)
	items, err := s.repos.Queues.List(ctx, tenantID, enabled, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func (s *QueueService) Delete(ctx context.Context, tenantID, name string) error {
	return s.repos.Queues.DeleteByTenantAndName(ctx, tenantID, name)
}
