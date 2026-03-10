package job

type Status string

const (
	StatusPending    Status = "Pending"
	StatusQueued     Status = "Queued"
	StatusScheduling Status = "Scheduling"
	StatusReserved   Status = "Reserved"
	StatusBound      Status = "Bound"
	StatusRunning    Status = "Running"
	StatusSucceeded  Status = "Succeeded"
	StatusFailed     Status = "Failed"
	StatusPreempted  Status = "Preempted"
	StatusCancelled  Status = "Cancelled"
)

func (s Status) IsTerminal() bool {
	switch s {
	case StatusSucceeded, StatusFailed, StatusCancelled:
		return true
	default:
		return false
	}
}
