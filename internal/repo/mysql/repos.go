package mysql

import "gorm.io/gorm"

type Repos struct {
	Tenants            *TenantRepo
	Queues             *QueueRepo
	Nodes              *NodeRepo
	NodeSnapshots      *NodeSnapshotRepo
	NodeHeartbeats     *NodeHeartbeatRepo
	GPUJobs            *GPUJobRepo
	GPUJobEvents       *GPUJobEventRepo
	Reservations       *ReservationRepo
	Allocations        *AllocationRepo
	Bindings           *BindingRepo
	GPUQuotas          *GPUQuotaRepo
	GPUPolicies        *GPUPolicyRepo
	SchedulingAttempts *SchedulingAttemptRepo
	JobRetries         *JobRetryRepo
	AuditLogs          *AuditLogRepo
	Outbox             *OutboxRepo
}

func NewRepos(db *gorm.DB) (*Repos, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	tenants, err := NewTenantRepo(db)
	if err != nil {
		return nil, err
	}
	queues, err := NewQueueRepo(db)
	if err != nil {
		return nil, err
	}
	nodes, err := NewNodeRepo(db)
	if err != nil {
		return nil, err
	}
	nodeSnapshots, err := NewNodeSnapshotRepo(db)
	if err != nil {
		return nil, err
	}
	nodeHeartbeats, err := NewNodeHeartbeatRepo(db)
	if err != nil {
		return nil, err
	}
	gpuJobs, err := NewGPUJobRepo(db)
	if err != nil {
		return nil, err
	}
	gpuJobEvents, err := NewGPUJobEventRepo(db)
	if err != nil {
		return nil, err
	}
	reservations, err := NewReservationRepo(db)
	if err != nil {
		return nil, err
	}
	allocations, err := NewAllocationRepo(db)
	if err != nil {
		return nil, err
	}
	bindings, err := NewBindingRepo(db)
	if err != nil {
		return nil, err
	}
	gpuQuotas, err := NewGPUQuotaRepo(db)
	if err != nil {
		return nil, err
	}
	gpuPolicies, err := NewGPUPolicyRepo(db)
	if err != nil {
		return nil, err
	}
	schedulingAttempts, err := NewSchedulingAttemptRepo(db)
	if err != nil {
		return nil, err
	}
	jobRetries, err := NewJobRetryRepo(db)
	if err != nil {
		return nil, err
	}
	auditLogs, err := NewAuditLogRepo(db)
	if err != nil {
		return nil, err
	}
	outbox, err := NewOutboxRepo(db)
	if err != nil {
		return nil, err
	}

	return &Repos{
		Tenants:            tenants,
		Queues:             queues,
		Nodes:              nodes,
		NodeSnapshots:      nodeSnapshots,
		NodeHeartbeats:     nodeHeartbeats,
		GPUJobs:            gpuJobs,
		GPUJobEvents:       gpuJobEvents,
		Reservations:       reservations,
		Allocations:        allocations,
		Bindings:           bindings,
		GPUQuotas:          gpuQuotas,
		GPUPolicies:        gpuPolicies,
		SchedulingAttempts: schedulingAttempts,
		JobRetries:         jobRetries,
		AuditLogs:          auditLogs,
		Outbox:             outbox,
	}, nil
}
