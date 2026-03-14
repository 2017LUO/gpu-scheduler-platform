package score

import (
	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type Binpack struct{}

func NewBinpack() *Binpack { return &Binpack{} }

func (p *Binpack) Name() string { return "Binpack" }

func (p *Binpack) Score(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) (int64, *schedframework.Status) {
	if job == nil || node == nil {
		return 0, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}

	requiredMemory := int64(job.GPUCount) * job.GPUMemoryMiB
	freeAfter := node.FreeMemoryMiB - requiredMemory
	gpuAfter := int64(node.HealthyGPUCount - job.GPUCount)

	if freeAfter < 0 {
		freeAfter = 0
	}
	if gpuAfter < 0 {
		gpuAfter = 0
	}

	score := int64(2_000_000) - (freeAfter / 64) - (gpuAfter * 10_000)
	if score < 0 {
		score = 0
	}
	return score, nil
}
