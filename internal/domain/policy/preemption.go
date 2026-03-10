package policy

type PreemptionPolicy struct {
	Enabled          bool
	AllowCrossTenant bool
	ReclaimBatchSize int
}
