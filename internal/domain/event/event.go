package event

import "time"

type Event struct {
	ID         string
	JobID      string
	TenantID   string
	Reason     Reason
	Message    string
	Source     string
	OccurredAt time.Time
}
