package syncer

import "context"

type AllocationSyncer struct{}

func NewAllocationSyncer() *AllocationSyncer {
	return &AllocationSyncer{}
}

func (s *AllocationSyncer) Sync(ctx context.Context) error {
	_ = ctx
	return nil
}
