package service

import (
	"context"

	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/job"
	"gpu-scheduler-platform/internal/scheduler/algorithm"
)

type PlacementService struct{}

func NewPlacementService() *PlacementService {
	return &PlacementService{}
}

func (s *PlacementService) Place(ctx context.Context, snapshot *cluster.Snapshot, j job.Job) (*algorithm.PlacementDecision, error) {
	_ = ctx
	return algorithm.ScheduleOne(snapshot, j)
}
