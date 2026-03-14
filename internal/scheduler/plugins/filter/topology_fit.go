package filter

import (
	"encoding/json"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type TopologyFit struct{}

func NewTopologyFit() *TopologyFit { return &TopologyFit{} }

func (p *TopologyFit) Name() string { return "TopologyFit" }

func (p *TopologyFit) Filter(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) *schedframework.Status {
	if job == nil || node == nil {
		return schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if !job.RequireNVLink {
		return nil
	}

	if len(node.TopologyJSON) == 0 {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "nvlink required but topology is missing")
	}

	var topo map[string]any
	if err := json.Unmarshal(node.TopologyJSON, &topo); err != nil {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "invalid topology json")
	}
	if _, ok := topo["nvlink"]; !ok {
		return schedframework.NewStatus(schedframework.CodeUnschedulable, "nvlink metadata is missing")
	}
	return nil
}
