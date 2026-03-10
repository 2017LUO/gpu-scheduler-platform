package repo

import (
	"context"
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/event"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/domain/policy"
)

type Tx interface {
	Commit() error
	Rollback() error
}

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type JobListFilter struct {
	TenantID      string
	Namespace     string
	Queue         job.QueueName
	Status        job.Status
	Priority      job.PriorityClass
	Keyword       string
	Limit         int
	Offset        int
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
}

type PendingJobFilter struct {
	Queues []job.QueueName
	Limit  int
}

type NodeSnapshotListFilter struct {
	NodeName string
	Limit    int
	Offset   int
}

type JobRepository interface {
	Create(ctx context.Context, j *job.Job) error
	Update(ctx context.Context, j *job.Job) error
	UpdateStatus(ctx context.Context, jobID string, status job.Status, message string) error
	GetByID(ctx context.Context, jobID string) (*job.Job, error)
	List(ctx context.Context, filter JobListFilter) ([]job.Job, int64, error)
	ListPending(ctx context.Context, filter PendingJobFilter) ([]job.Job, error)
	Delete(ctx context.Context, jobID string) error
}

type JobEventRepository interface {
	Create(ctx context.Context, e *event.Event) error
	ListByJobID(ctx context.Context, jobID string, limit, offset int) ([]event.Event, int64, error)
}

type NodeSnapshotRepository interface {
	UpsertSnapshot(ctx context.Context, s *cluster.Snapshot) error
	GetLatest(ctx context.Context) (*cluster.Snapshot, error)
	List(ctx context.Context, filter NodeSnapshotListFilter) ([]cluster.Node, int64, error)
}

type AllocationRepository interface {
	Create(ctx context.Context, a *allocation.Allocation) error
	Update(ctx context.Context, a *allocation.Allocation) error
	UpdateStatus(ctx context.Context, allocationID string, status allocation.Status, message string) error
	GetByID(ctx context.Context, allocationID string) (*allocation.Allocation, error)
	GetByJobID(ctx context.Context, jobID string) (*allocation.Allocation, error)
	ListByNode(ctx context.Context, nodeName string, limit, offset int) ([]allocation.Allocation, int64, error)
	ReleaseByJobID(ctx context.Context, jobID string, releasedAt time.Time, message string) error
}

type ReservationRepository interface {
	Create(ctx context.Context, r *allocation.Reservation) error
	GetByJobID(ctx context.Context, jobID string) (*allocation.Reservation, error)
	DeleteByJobID(ctx context.Context, jobID string) error
	DeleteExpired(ctx context.Context, now time.Time, limit int) (int64, error)
}

type BindingRepository interface {
	Create(ctx context.Context, b *allocation.Binding) error
	GetByJobID(ctx context.Context, jobID string) (*allocation.Binding, error)
}

type QuotaRepository interface {
	Upsert(ctx context.Context, q *policy.Quota) error
	GetByTenant(ctx context.Context, tenantID string, namespace string) (*policy.Quota, error)
	List(ctx context.Context, limit, offset int) ([]policy.Quota, int64, error)
	Delete(ctx context.Context, id string) error
}

type TenantRepository interface {
	Exists(ctx context.Context, tenantID string) (bool, error)
}
