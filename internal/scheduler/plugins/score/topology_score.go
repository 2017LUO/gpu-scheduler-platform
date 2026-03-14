package score

import (
	"encoding/json"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type TopologyScore struct{}

func NewTopologyScore() *TopologyScore { return &TopologyScore{} }

func (p *TopologyScore) Name() string { return "TopologyScore" }

func (p *TopologyScore) Score(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) (int64, *schedframework.Status) {
	if job == nil || node == nil {
		return 0, schedframework.NewStatus(schedframework.CodeError, "job or node is nil")
	}
	if len(node.TopologyJSON) == 0 {
		return 0, nil
	}

	var topo map[string]any
	if err := json.Unmarshal(node.TopologyJSON, &topo); err != nil {
		return 0, nil
	}

	score := int64(0)
	if v, ok := topo["nvlink"]; ok && v != nil {
		score += 3000
	}
	if v, ok := topo["pcie_switches"]; ok {
		if n, ok2 := v.(float64); ok2 {
			score -= int64(n * 100)
		}
	}
	return score, nil
}
