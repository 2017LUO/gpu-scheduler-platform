package framework

import (
	model "gpu-scheduler-platform/internal/repo/models"

	"go.uber.org/zap"
)

type Framework struct {
	registry   *Registry
	logger     *zap.Logger
	preFilters []PreFilterPlugin
	filters    []FilterPlugin
	scores     []ScorePlugin
	reserves   []ReservePlugin
	permits    []PermitPlugin
}

func NewFramework(registry *Registry, lg *zap.Logger) *Framework {
	if lg == nil {
		lg = zap.NewNop()
	}
	if registry == nil {
		registry = NewRegistry()
	}
	return &Framework{
		registry: registry,
		logger:   lg,
	}
}

func (f *Framework) AddPreFilter(p PreFilterPlugin) { f.preFilters = append(f.preFilters, p) }
func (f *Framework) AddFilter(p FilterPlugin)       { f.filters = append(f.filters, p) }
func (f *Framework) AddScore(p ScorePlugin)         { f.scores = append(f.scores, p) }
func (f *Framework) AddReserve(p ReservePlugin)     { f.reserves = append(f.reserves, p) }
func (f *Framework) AddPermit(p PermitPlugin)       { f.permits = append(f.permits, p) }

func (f *Framework) RunPreFilter(cs *CycleState, job *model.GPUJob) *Status {
	for _, p := range f.preFilters {
		if st := p.PreFilter(cs, job); !st.IsSuccess() {
			return st
		}
	}
	return nil
}

func (f *Framework) RunFilter(cs *CycleState, job *model.GPUJob, node *model.Node) *Status {
	for _, p := range f.filters {
		if st := p.Filter(cs, job, node); !st.IsSuccess() {
			return st
		}
	}
	return nil
}

func (f *Framework) RunScore(cs *CycleState, job *model.GPUJob, node *model.Node) (int64, *Status) {
	var total int64
	for _, p := range f.scores {
		score, st := p.Score(cs, job, node)
		if !st.IsSuccess() {
			return 0, st
		}
		total += score
	}
	return total, nil
}

func (f *Framework) RunReserve(cs *CycleState, job *model.GPUJob, node *model.Node) *Status {
	for _, p := range f.reserves {
		if st := p.Reserve(cs, job, node); !st.IsSuccess() {
			for i := len(f.reserves) - 1; i >= 0; i-- {
				f.reserves[i].Unreserve(cs, job, node)
			}
			return st
		}
	}
	return nil
}

func (f *Framework) RunPermit(cs *CycleState, job *model.GPUJob, node *model.Node) *Status {
	for _, p := range f.permits {
		if st := p.Permit(cs, job, node); !st.IsSuccess() {
			return st
		}
	}
	return nil
}
