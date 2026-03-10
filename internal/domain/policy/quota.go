package policy

type Quota struct {
	ID              string
	TenantID        string
	Namespace       string
	MaxGPUCount     int
	MaxRunningJobs  int
	MaxQueuedJobs   int
	MaxGPUMemoryMiB int64
	Enabled         bool
}

func (q Quota) Valid() bool {
	if q.MaxGPUCount < 0 || q.MaxRunningJobs < 0 || q.MaxQueuedJobs < 0 || q.MaxGPUMemoryMiB < 0 {
		return false
	}
	return true
}
