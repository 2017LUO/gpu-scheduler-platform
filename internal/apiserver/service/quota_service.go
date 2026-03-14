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

type QuotaService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewQuotaService(repos *repoimpl.Repos, lg *zap.Logger) *QuotaService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &QuotaService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *QuotaService) Upsert(ctx context.Context, req dto.UpsertQuotaRequest) (*model.GPUQuota, error) {
	if req.TenantID == "" {
		return nil, repoimpl.ErrInvalidArgument
	}
	if _, err := s.repos.Tenants.Get(ctx, req.TenantID); err != nil {
		return nil, err
	}

	now := s.nowFunc()
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	existing, err := s.repos.GPUQuotas.GetByTenantAndNamespace(ctx, req.TenantID, req.Namespace)
	id := newID("quota")
	createdAt := now
	if err == nil {
		id = existing.ID
		createdAt = existing.CreatedAt
	} else if !errors.Is(err, repoimpl.ErrNotFound) {
		return nil, err
	}

	m := &model.GPUQuota{
		ID:              id,
		TenantID:        req.TenantID,
		Namespace:       req.Namespace,
		MaxGPUCount:     req.MaxGPUCount,
		MaxRunningJobs:  req.MaxRunningJobs,
		MaxQueuedJobs:   req.MaxQueuedJobs,
		MaxGPUMemoryMiB: req.MaxGPUMemoryMiB,
		Enabled:         enabled,
		CreatedAt:       createdAt,
		UpdatedAt:       now,
	}

	if err := s.repos.GPUQuotas.Upsert(ctx, m); err != nil {
		return nil, err
	}
	return s.repos.GPUQuotas.GetByTenantAndNamespace(ctx, req.TenantID, req.Namespace)
}

func (s *QuotaService) Get(ctx context.Context, tenantID, namespace string) (*model.GPUQuota, error) {
	return s.repos.GPUQuotas.GetByTenantAndNamespace(ctx, tenantID, namespace)
}

func (s *QuotaService) ListByTenant(ctx context.Context, tenantID string) ([]model.GPUQuota, error) {
	return s.repos.GPUQuotas.ListByTenant(ctx, tenantID)
}

func (s *QuotaService) Delete(ctx context.Context, tenantID, namespace string) error {
	return s.repos.GPUQuotas.DeleteByTenantAndNamespace(ctx, tenantID, namespace)
}
