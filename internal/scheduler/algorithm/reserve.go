package algorithm

import (
	"time"

	"gpu-scheduler-platform/internal/domain/allocation"
)

func BuildReservation(decision *PlacementDecision, ttl time.Duration, now time.Time) *allocation.Reservation {
	if decision == nil {
		return nil
	}
	gpuIDs := make([]string, 0, len(decision.GPUs))
	for _, g := range decision.GPUs {
		gpuIDs = append(gpuIDs, g.ID)
	}

	return &allocation.Reservation{
		ID:        decision.Job.ID + "-reservation",
		JobID:     decision.Job.ID,
		NodeName:  decision.Node.Name,
		GPUIDs:    gpuIDs,
		ExpireAt:  now.Add(ttl),
		CreatedAt: now,
	}
}
