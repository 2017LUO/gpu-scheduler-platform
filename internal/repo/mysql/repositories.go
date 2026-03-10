package mysql

import (
	"gpu-scheduler-platform/internal/repo"

	"gorm.io/gorm"
)

type Repositories struct {
	TxManager     repo.TxManager
	Jobs          repo.JobRepository
	JobEvents     repo.JobEventRepository
	NodeSnapshots repo.NodeSnapshotRepository
	Allocations   repo.AllocationRepository
	Reservations  repo.ReservationRepository
	Bindings      repo.BindingRepository
	Quotas        repo.QuotaRepository
	Tenants       repo.TenantRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		TxManager:     NewTxManager(db),
		Jobs:          NewGPUJobRepo(db),
		JobEvents:     NewJobEventRepo(db),
		NodeSnapshots: NewNodeSnapshotRepo(db),
		Allocations:   NewAllocationRepo(db),
		Reservations:  NewReservationRepo(db),
		Bindings:      NewBindingRepo(db),
		Quotas:        NewGPUQuotaRepo(db),
		Tenants:       NewTenantRepo(db),
	}
}
