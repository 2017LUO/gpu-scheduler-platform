package allocation

import "time"

type Status string

const (
	StatusPending   Status = "Pending"
	StatusReserved  Status = "Reserved"
	StatusCommitted Status = "Committed"
	StatusReleased  Status = "Released"
	StatusFailed    Status = "Failed"
)

type Allocation struct {
	ID       string
	JobID    string
	TenantID string
	NodeName string
	GPUIDs   []string
	Status   Status
	Message  string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	CommittedAt *time.Time
	ReleasedAt  *time.Time
}
