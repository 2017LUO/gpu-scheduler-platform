package service

import (
	"fmt"

	repoimpl "gpu-scheduler-platform/internal/repo/mysql"

	"go.uber.org/zap"
)

type Services struct {
	Job           *JobService
	InternalAgent *InternalAgentService
	Tenant        *TenantService
	Queue         *QueueService
	Quota         *QuotaService
	Policy        *PolicyService
	Cluster       *ClusterService
}

func NewServices(repos *repoimpl.Repos, lg *zap.Logger) (*Services, error) {
	if repos == nil {
		return nil, fmt.Errorf("repos is nil")
	}
	if lg == nil {
		lg = zap.NewNop()
	}

	return &Services{
		Job:           NewJobService(repos, lg),
		InternalAgent: NewInternalAgentService(repos, lg),
		Tenant:        NewTenantService(repos, lg),
		Queue:         NewQueueService(repos, lg),
		Quota:         NewQuotaService(repos, lg),
		Policy:        NewPolicyService(repos, lg),
		Cluster:       NewClusterService(repos, lg),
	}, nil
}
