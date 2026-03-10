package algorithm

import (
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
)

func BuildAllocation(decision *PlacementDecision, now time.Time) *allocation.Allocation {
	if decision == nil {
		return nil
	}
	gpuIDs := make([]string, 0, len(decision.GPUs))
	for _, g := range decision.GPUs {
		gpuIDs = append(gpuIDs, g.ID)
	}

	return &allocation.Allocation{
		ID:        decision.Job.ID + "-allocation",
		JobID:     decision.Job.ID,
		TenantID:  decision.Job.TenantID,
		NodeName:  decision.Node.Name,
		GPUIDs:    gpuIDs,
		Status:    allocation.StatusCommitted,
		Message:   decision.Reason,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
