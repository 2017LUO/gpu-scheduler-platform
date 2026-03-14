package service

import (
	"context"
	"fmt"
	"time"

	"gpu-scheduler-platform/internal/apiserver/dto"
	model "gpu-scheduler-platform/internal/repo/models"
	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type TenantService struct {
	repos   *repoimpl.Repos
	logger  *zap.Logger
	nowFunc func() time.Time
}

func NewTenantService(repos *repoimpl.Repos, lg *zap.Logger) *TenantService {
	if lg == nil {
		lg = zap.NewNop()
	}
	return &TenantService{
		repos:   repos,
		logger:  lg,
		nowFunc: func() time.Time { return time.Now().UTC() },
	}
}

func (s *TenantService) Create(ctx context.Context, req dto.CreateTenantRequest) (*model.Tenant, error) {
	if req.Name == "" {
		return nil, repoimpl.ErrInvalidArgument
	}

	now := s.nowFunc()
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}

	m := &model.Tenant{
		ID:          newID("tenant"),
		Name:        req.Name,
		Enabled:     enabled,
		Description: req.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repos.Tenants.Create(ctx, m); err != nil {
		return nil, err
	}

	_ = s.repos.AuditLogs.Create(ctx, &model.AuditLog{
		TenantID:     m.ID,
		Actor:        "api-server",
		Action:       "tenant.create",
		ResourceType: "tenant",
		ResourceID:   m.ID,
		ResourceName: m.Name,
		Status:       "SUCCESS",
		CreatedAt:    now,
	})

	return m, nil
}

func (s *TenantService) Get(ctx context.Context, id string) (*model.Tenant, error) {
	return s.repos.Tenants.Get(ctx, id)
}

func (s *TenantService) List(ctx context.Context, enabled *bool, limit, offset int) ([]model.Tenant, repoimpl.PageQuery, error) {
	page := repoimpl.PageQuery{Limit: limit, Offset: offset}.Normalize(50, 500)
	items, err := s.repos.Tenants.List(ctx, enabled, page)
	if err != nil {
		return nil, page, err
	}
	return items, page, nil
}

func (s *TenantService) Update(ctx context.Context, id string, req dto.UpdateTenantRequest) (*model.Tenant, error) {
	if id == "" {
		return nil, repoimpl.ErrInvalidArgument
	}
	if req.Name == nil && req.Enabled == nil && req.Description == nil {
		return nil, repoimpl.ErrInvalidArgument
	}

	current, err := s.repos.Tenants.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Enabled != nil {
		current.Enabled = *req.Enabled
	}
	if req.Description != nil {
		current.Description = req.Description
	}
	current.UpdatedAt = s.nowFunc()

	if current.Name == "" {
		return nil, fmt.Errorf("tenant name cannot be empty")
	}

	if err := s.repos.Tenants.Update(ctx, current); err != nil {
		return nil, err
	}
	return s.repos.Tenants.Get(ctx, id)
}

func (s *TenantService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return repoimpl.ErrInvalidArgument
	}
	return s.repos.Tenants.Delete(ctx, id)
}
