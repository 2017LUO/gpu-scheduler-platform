package job

import "time"

type Job struct {
	ID          string
	TenantID    string
	Namespace   string
	Name        string
	Queue       QueueName
	Priority    PriorityClass
	Status      Status
	Requirement Requirement

	Labels      map[string]string
	Annotations map[string]string

	RetryCount int
	Message    string

	CreatedAt   time.Time
	UpdatedAt   time.Time
	ScheduledAt *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time
}

func (j Job) CanSchedule() bool {
	switch j.Status {
	case StatusPending, StatusQueued:
		return true
	default:
		return false
	}
}

func (j Job) PriorityScore() int {
	return PriorityValue(j.Priority)
}
