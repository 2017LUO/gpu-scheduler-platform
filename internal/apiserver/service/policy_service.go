package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

type PolicyService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewPolicyService(repos *repoimpl.Repos, lg *zap.Logger) *PolicyService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &PolicyService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *PolicyService) Upsert(ctx context.Context, req dto.UpsertPolicyRequest) (*model.GPUPolicy, error) {
	if req.TenantID == "" || req.Name == "" {
		return nil, repoimpl.ErrInvalidArgument
	}

	if _, err := s.repos.Tenants.Get(ctx, req.TenantID); err != nil {
		return nil, err
	}
	if req.Queue != "" {
		if _, err := s.repos.Queues.GetByName(ctx, req.TenantID, req.Queue); err != nil {
			return nil, err
		}
	}

	now := s.nowFunc()
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	existing, err := s.repos.GPUPolicies.GetByTenantAndName(ctx, req.TenantID, req.Name)
	id := newID("policy")
	createdAt := now
	if err == nil {
		id = existing.ID
		createdAt = existing.CreatedAt
	} else if !errors.Is(err, repoimpl.ErrNotFound) {
		return nil, err
	}

	m := &model.GPUPolicy{
		ID:                     id,
		TenantID:               req.TenantID,
		Name:                   req.Name,
		Queue:                  req.Queue,
		Priority:               req.Priority,
		Enabled:                enabled,
		Preemptible:            req.Preemptible,
		RequireHealthy:         req.RequireHealthy,
		RequireMIG:             req.RequireMIG,
		MaxGPUCount:            req.MaxGPUCount,
		RequiredGPUModel:       req.RequiredGPUModel,
		RequiredNodeLabelsJSON: mustJSONAny(req.RequiredNodeLabels),
		SelectorJSON:           mustJSONAny(req.Selector),
		Description:            req.Description,
		CreatedAt:              createdAt,
		UpdatedAt:              now,
	}

	if err := s.repos.GPUPolicies.Upsert(ctx, m); err != nil {
		return nil, err
	}
	return s.repos.GPUPolicies.GetByTenantAndName(ctx, req.TenantID, req.Name)
}

func (s *PolicyService) Get(ctx context.Context, tenantID, name string) (*model.GPUPolicy, error) {
	return s.repos.GPUPolicies.GetByTenantAndName(ctx, tenantID, name)
}

func (s *PolicyService) List(ctx context.Context, tenantID string, enabled *bool, limit, offset int) ([]model.GPUPolicy, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{Limit: limit, Offset: offset}.Normalize(50, 500)
	items, err := s.repos.GPUPolicies.List(ctx, tenantID, enabled, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func (s *PolicyService) Delete(ctx context.Context, tenantID, name string) error {
	return s.repos.GPUPolicies.DeleteByTenantAndName(ctx, tenantID, name)
}

func mustJSONAny(v any) datatypes.JSON {
	if v == nil {
		return datatypes.JSON([]byte("{}"))
	}
	b, err := json.Marshal(v)
	if err != nil || len(b) == 0 {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(b)
}
