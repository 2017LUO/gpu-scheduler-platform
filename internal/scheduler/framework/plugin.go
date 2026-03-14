package framework

import model "gpu-scheduler-platform/internal/repo/models"

type Plugin interface {
	Name() string
}

type PreFilterPlugin interface {
	Plugin
	PreFilter(*CycleState, *model.GPUJob) *Status
}

type FilterPlugin interface {
	Plugin
	Filter(*CycleState, *model.GPUJob, *model.Node) *Status
}

type ScorePlugin interface {
	Plugin
	Score(*CycleState, *model.GPUJob, *model.Node) (int64, *Status)
}

type ReservePlugin interface {
	Plugin
	Reserve(*CycleState, *model.GPUJob, *model.Node) *Status
	Unreserve(*CycleState, *model.GPUJob, *model.Node)
}

type PermitPlugin interface {
	Plugin
	Permit(*CycleState, *model.GPUJob, *model.Node) *Status
}
