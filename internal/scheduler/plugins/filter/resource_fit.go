package filter

import (
	"fmt"
	"strings"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type ResourceFit struct{}

func NewResourceFit() *ResourceFit { return &ResourceFit{} }

func (p *ResourceFit) Name() string { return "ResourceFit" }

func (p *ResourceFit) Filter(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) *schedframework.Status {
	if job == nil || node == nil {
		return schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if !node.Schedulable {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "node is unschedulable")
	}
	if node.State != "" && !strings.EqualFold(node.State, "READY") {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, fmt.Sprintf("node state is %s", node.State))
	}
	if job.GPUCount > 0 && node.GPUCount < job.GPUCount {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "insufficient gpu count")
	}
	if job.RequireHealthy && node.HealthyGPUCount < job.GPUCount {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "insufficient healthy gpu count")
	}

	requiredMemory := int64(job.GPUCount) * job.GPUMemoryMiB
	if requiredMemory > 0 && node.FreeMemoryMiB < requiredMemory {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "insufficient free gpu memory")
	}
	return nil
}
