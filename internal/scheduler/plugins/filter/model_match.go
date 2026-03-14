package filter

import (
	"strings"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type ModelMatch struct{}

func NewModelMatch() *ModelMatch { return &ModelMatch{} }

func (p *ModelMatch) Name() string { return "ModelMatch" }

func (p *ModelMatch) Filter(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) *schedframework.Status {
	if job == nil || node == nil {
		return schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if strings.TrimSpace(job.GPUModel) == "" {
		return nil
	}

	inventory, ok := schedframework.ReadNodeGPUInventory(cs)
	if !ok {
		return nil
	}
	gpus := inventory[node.NodeName]
	if len(gpus) == 0 {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "node has no gpu inventory")
	}

	matched := 0
	for _, gpu := range gpus {
		if gpu.Allocated || gpu.Reserved {
			continue
		}
		if job.RequireHealthy && !gpu.Healthy {
			continue
		}
		if !strings.EqualFold(strings.TrimSpace(gpu.Model), strings.TrimSpace(job.GPUModel)) {
			continue
		}
		if job.GPUMemoryMiB > 0 && gpu.FreeMemoryMiB < job.GPUMemoryMiB {
			continue
		}
		matched++
	}

	if matched < job.GPUCount {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "not enough gpus matching required model")
	}
	return nil
}
