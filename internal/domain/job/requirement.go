package job

import "time"

type Requirement struct {
	GPUCount        int
	GPUMemoryMiB    int64
	GPUModel        string
	RequireSameNode bool
	RequireHealthy  bool
	RequireMIG      bool
	MIGProfile      string

	// 拓扑约束
	RequireNVLink       bool
	PreferredNodeLabels map[string]string
	PreferredGPULabels  map[string]string

	// 调度行为
	Preemptible      bool
	Retryable        bool
	MaxRetry         int
	ExpectedDuration time.Duration
}

func (r Requirement) Valid() bool {
	if r.GPUCount <= 0 {
		return false
	}
	if r.GPUMemoryMiB < 0 {
		return false
	}
	if r.MaxRetry < 0 {
		return false
	}
	return true
}
