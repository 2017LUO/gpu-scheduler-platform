package score

import (
	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type FragmentationScore struct{}

func NewFragmentationScore() *FragmentationScore { return &FragmentationScore{} }

func (p *FragmentationScore) Name() string { return "FragmentationScore" }

func (p *FragmentationScore) Score(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) (int64, *schedframework.Status) {
	if job == nil || node == nil {
		return 0, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}

	freeAfterGPU := node.HealthyGPUCount - job.GPUCount
	freeAfterMem := node.FreeMemoryMiB - int64(job.GPUCount)*job.GPUMemoryMiB
	if freeAfterGPU < 0 {
		return 0, nil
	}
	if freeAfterMem < 0 {
		return 0, nil
	}

	score := int64(500_000) - int64(freeAfterGPU*5000) - (freeAfterMem / 128)
	if score < 0 {
		score = 0
	}
	return score, nil
}
