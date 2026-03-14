package score

import (
	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type Spread struct{}

func NewSpread() *Spread { return &Spread{} }

func (p *Spread) Name() string { return "Spread" }

func (p *Spread) Score(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) (int64, *schedframework.Status) {
	if job == nil || node == nil {
		return 0, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	score := int64(node.HealthyGPUCount*10_000) + (node.FreeMemoryMiB / 64)
	return score, nil
}
