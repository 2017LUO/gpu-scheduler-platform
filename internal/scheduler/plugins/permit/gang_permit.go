package permit

import (
	"encoding/json"

	model "gpu-scheduler-platform/internal/repo/models"
	schedframework "gpu-scheduler-platform/internal/scheduler/framework"
)

type GangPermit struct{}

func NewGangPermit() *GangPermit { return &GangPermit{} }

func (p *GangPermit) Name() string { return "GangPermit" }

func (p *GangPermit) Permit(cs *schedframework.CycleState, job *model.GPUJob, node *model.Node) *schedframework.Status {
	if job == nil {
		return schedframework.NewStatus(schedframework.CodeError, "job is nil")
	}
	if len(job.RunPolicyJSON) == 0 {
		return nil
	}

	var policy map[string]any
	if err := json.Unmarshal(job.RunPolicyJSON, &policy); err != nil {
		return nil
	}
	v, ok := policy["gang_size"]
	if !ok {
		return nil
	}

	switch x := v.(type) {
	case float64:
		if int(x) > 1 {
			return schedframework.NewStatus(schedframework.CodeWait, "gang scheduling requested, waiting for external coordinator")
		}
	case int:
		if x > 1 {
			return schedframework.NewStatus(schedframework.CodeWait, "gang scheduling requested, waiting for external coordinator")
		}
	}
	return nil
}
