package policy

type FairnessStrategy string

const (
	FairnessNone           FairnessStrategy = "none"
	FairnessTenantWeighted FairnessStrategy = "tenant_weighted"
	FairnessDRF            FairnessStrategy = "drf"
)

type Fairness struct {
	Enabled       bool
	Strategy      FairnessStrategy
	DefaultWeight int
	TenantWeights map[string]int
}

func (f Fairness) WeightOf(tenantID string) int {
	if tenantID == "" {
		if f.DefaultWeight > 0 {
			return f.DefaultWeight
		}
		return 1
	}
	if w, ok := f.TenantWeights[tenantID]; ok && w > 0 {
		return w
	}
	if f.DefaultWeight > 0 {
		return f.DefaultWeight
	}
	return 1
}
