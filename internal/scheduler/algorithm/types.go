package algorithm

import (
	"gpu-scheduler-platform/internal/domain/cluster"
	"gpu-scheduler-platform/internal/domain/job"
)

type PlacementDecision struct {
	Job    job.Job
	Node   cluster.Node
	GPUs   []cluster.GPU
	Reason string
}

type NodeScore struct {
	Node  cluster.Node
	Score int
}
