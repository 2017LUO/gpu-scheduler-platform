package score

import (
	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type UtilizationScore struct{}

func NewUtilizationScore() *UtilizationScore { return &UtilizationScore{} }

func (p *UtilizationScore) Name() string { return "UtilizationScore" }

func (p *UtilizationScore) Score(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) (int64, *schedframework.Status) {
	if job == nil || node == nil {
		return 0, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}

	totalGPU := node.GPUCount
	if totalGPU <= 0 {
		return 0, nil
	}
	usedGPU := totalGPU - node.HealthyGPUCount
	if usedGPU < 0 {
		usedGPU = 0
	}

	score := int64(usedGPU*3000) + (node.TotalMemoryMiB-node.FreeMemoryMiB)/256
	return score, nil
}
