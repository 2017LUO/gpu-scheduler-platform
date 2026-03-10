package algorithm

import (
	"fmt"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/job"
)

func ScheduleOne(snapshot *cluster.Snapshot, j job.Job) (*PlacementDecision, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("snapshot is nil")
	}
	if !j.CanSchedule() {
		return nil, fmt.Errorf("job %s is not schedulable", j.ID)
	}

	decision, ok := SelectNode(snapshot, j)
	if !ok {
		return nil, fmt.Errorf("no feasible node found")
	}
	return decision, nil
}
