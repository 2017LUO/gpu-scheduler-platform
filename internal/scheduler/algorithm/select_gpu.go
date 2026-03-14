package algorithm

import (
	"context"
	"sort"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

func SelectGPU(
	ctx context.Context,
	deps Dependencies,
	cs *schedframework.CycleState,
	job *model.GPUJob,
	node *model.Node,
) ([]string, *schedframework.Status) {
	if job == nil || node == nil {
		return nil, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if deps.SelectNodeGPUs == nil {
		return nil, schedframework.NewStatus(schedframework.CodeError, "gpu selector is not configured")
	}

	items, err := deps.SelectNodeGPUs(ctx, job, node)
	if err != nil {
		return nil, schedframework.AsError(err)
	}
	if len(items) < job.GPUCount {
		return nil, schedframework.NewStatus(schedframework.CodeUnschedulable, "not enough matching gpus on selected node")
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].FreeMemoryMiB != items[j].FreeMemoryMiB {
			return items[i].FreeMemoryMiB > items[j].FreeMemoryMiB
		}
		if items[i].UtilizationGPU != items[j].UtilizationGPU {
			return items[i].UtilizationGPU < items[j].UtilizationGPU
		}
		return items[i].GPUIndex < items[j].GPUIndex
	})

	out := make([]string, 0, job.GPUCount)
	for i := 0; i < job.GPUCount; i++ {
		out = append(out, items[i].UUID)
	}

	cs.Write(schedframework.StateKeySelectedGPUUUIDs, out)
	return out, nil
}
